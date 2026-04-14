package runner

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	bpsv1alpha1 "github.com/sebrandon1/bps-operator/api/v1alpha1"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const scanPollInterval = 2 * time.Second

type OperatorRunner struct {
	k8sClient   kubernetes.Interface
	ctrlClient  ctrlclient.Client
	config      *rest.Config
	isOCP       bool
	opts        Options
	printMu     sync.Mutex
	printFn     func(string)
}

func NewOperator(k8sClient kubernetes.Interface, ctrlClient ctrlclient.Client, config *rest.Config, isOCP bool, opts Options, printFn func(string)) *OperatorRunner {
	if opts.Parallel < 1 {
		opts.Parallel = 1
	}
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Minute
	}
	return &OperatorRunner{
		k8sClient:  k8sClient,
		ctrlClient: ctrlClient,
		config:     config,
		isOCP:      isOCP,
		opts:       opts,
		printFn:    printFn,
	}
}

func (r *OperatorRunner) Run(scenarios []scenario.Scenario) Summary {
	start := time.Now()

	results := make([]ScenarioResult, len(scenarios))
	sem := make(chan struct{}, r.opts.Parallel)
	var wg sync.WaitGroup

	for i, s := range scenarios {
		if r.shouldSkip(s) {
			results[i] = ScenarioResult{
				Name:       s.Name,
				CheckName:  s.CheckName,
				Skipped:    true,
				SkipReason: r.skipReason(s),
			}
			r.print(FormatResult(&results[i], r.opts.Verbose))
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, sc scenario.Scenario) {
			defer wg.Done()
			result, nsToClean := r.runOne(sc)
			<-sem
			results[idx] = result
			r.print(FormatResult(&results[idx], r.opts.Verbose))
			if nsToClean != "" {
				r.cleanupNamespace(nsToClean)
			}
		}(i, s)
	}

	wg.Wait()

	summary := Summary{
		Total:    len(scenarios),
		Duration: time.Since(start),
	}
	for _, res := range results {
		switch {
		case res.Skipped:
			summary.Skipped++
		case res.Err != nil:
			summary.Errors++
		case res.Passed():
			summary.Passed++
		default:
			summary.Failed++
		}
	}
	return summary
}

func (r *OperatorRunner) runOne(s scenario.Scenario) (ScenarioResult, string) {
	start := time.Now()
	result := ScenarioResult{
		Name:           s.Name,
		CheckName:      s.CheckName,
		ExpectedStatus: s.ExpectedStatus,
	}
	var ns string
	fail := func(err error) (ScenarioResult, string) {
		result.Err = err
		result.Duration = time.Since(start)
		return result, ns
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.opts.Timeout)
	defer cancel()

	prefix := scenarioPrefix(s)
	var err error
	ns, err = cluster.CreateNamespace(ctx, r.k8sClient, prefix, s.Privileged)
	if err != nil {
		return fail(fmt.Errorf("creating namespace: %w", err))
	}

	runCtx := &scenario.RunContext{
		Ctx:       ctx,
		Client:    r.k8sClient,
		Config:    r.config,
		Namespace: ns,
		IsOCP:     r.isOCP,
	}

	if err := s.Setup(runCtx); err != nil {
		return fail(fmt.Errorf("setup: %w", err))
	}

	scannerName := scannerName(s)
	scanner := &bpsv1alpha1.BestPracticeScanner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      scannerName,
			Namespace: ns,
		},
		Spec: bpsv1alpha1.BestPracticeScannerSpec{
			TargetNamespace: ns,
			Checks:          []string{s.CheckName},
		},
	}

	if err := r.ctrlClient.Create(ctx, scanner); err != nil {
		return fail(fmt.Errorf("creating scanner CR: %w", err))
	}

	if err := r.waitForScanComplete(ctx, scanner); err != nil {
		return fail(fmt.Errorf("waiting for scan: %w", err))
	}

	var resultList bpsv1alpha1.BestPracticeResultList
	if err := r.ctrlClient.List(ctx, &resultList, ctrlclient.InNamespace(ns)); err != nil {
		return fail(fmt.Errorf("listing results: %w", err))
	}

	found := false
	for _, bpr := range resultList.Items {
		if bpr.Spec.CheckName == s.CheckName {
			found = true
			result.ActualStatus = string(bpr.Spec.ComplianceStatus)
			result.Reason = bpr.Spec.Reason
			break
		}
	}

	if !found {
		return fail(fmt.Errorf("no BestPracticeResult found for check %q", s.CheckName))
	}

	result.Duration = time.Since(start)
	return result, ns
}

func (r *OperatorRunner) waitForScanComplete(ctx context.Context, scanner *bpsv1alpha1.BestPracticeScanner) error {
	key := ctrlclient.ObjectKeyFromObject(scanner)
	ticker := time.NewTicker(scanPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for scanner %s to complete", scanner.Name)
		case <-ticker.C:
			if err := r.ctrlClient.Get(ctx, key, scanner); err != nil {
				return err
			}
			switch scanner.Status.Phase {
			case bpsv1alpha1.PhaseCompleted:
				return nil
			case bpsv1alpha1.PhaseError:
				return fmt.Errorf("scanner entered Error phase")
			}
		}
	}
}

func (r *OperatorRunner) shouldSkip(s scenario.Scenario) bool {
	if s.PostDiscovery != nil {
		return true
	}
	if s.RequiresOCP && (r.opts.SkipOCP || !r.isOCP) {
		return true
	}
	if s.RequiresProbe && r.opts.SkipProbe {
		return true
	}
	if s.RequiresDualStack {
		return true
	}
	if s.RequiresMultus {
		return true
	}
	return false
}

func (r *OperatorRunner) skipReason(s scenario.Scenario) string {
	if s.PostDiscovery != nil {
		return "uses PostDiscovery (not supported in operator mode)"
	}
	if s.RequiresOCP && !r.isOCP {
		return "requires OpenShift cluster"
	}
	if s.RequiresOCP && r.opts.SkipOCP {
		return "OCP scenarios skipped (--skip-ocp)"
	}
	if s.RequiresProbe && r.opts.SkipProbe {
		return "probe scenarios skipped (--skip-probe)"
	}
	if s.RequiresDualStack {
		return "requires dual-stack cluster (not detected in operator mode)"
	}
	if s.RequiresMultus {
		return "requires Multus CNI (not detected in operator mode)"
	}
	return "skipped"
}

func (r *OperatorRunner) cleanupNamespace(ns string) {
	cleanCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := cluster.DeleteNamespace(cleanCtx, r.k8sClient, ns); err != nil {
		r.print(fmt.Sprintf("WARNING: failed to delete namespace %s: %v\n", ns, err))
	}
}

func (r *OperatorRunner) print(msg string) {
	r.printMu.Lock()
	defer r.printMu.Unlock()
	r.printFn(msg)
}

func scannerName(s scenario.Scenario) string {
	h := sha256.Sum256([]byte(s.Name))
	return fmt.Sprintf("cqe-%x", h[:4])
}
