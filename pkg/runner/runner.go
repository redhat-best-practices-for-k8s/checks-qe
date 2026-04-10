package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/discovery"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/scenario"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Options struct {
	Parallel  int
	Verbose   bool
	SkipProbe bool
	SkipOCP   bool
	Timeout   time.Duration
}

type Runner struct {
	client      kubernetes.Interface
	config      *rest.Config
	isOCP       bool
	opts        Options
	clusterSnap *discovery.ClusterSnapshot
	printMu     sync.Mutex
	printFn     func(string)
}

func New(client kubernetes.Interface, config *rest.Config, isOCP bool, opts Options, printFn func(string)) *Runner {
	if opts.Parallel < 1 {
		opts.Parallel = 1
	}
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Minute
	}
	return &Runner{
		client:  client,
		config:  config,
		isOCP:   isOCP,
		opts:    opts,
		printFn: printFn,
	}
}

func (r *Runner) Run(scenarios []scenario.Scenario) Summary {
	start := time.Now()

	snap, err := discovery.FetchClusterSnapshot(context.Background(), r.client)
	if err != nil {
		r.print(fmt.Sprintf("WARNING: failed to fetch cluster snapshot: %v\n", err))
	} else {
		r.clusterSnap = snap
	}

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

func (r *Runner) runOne(s scenario.Scenario) (ScenarioResult, string) {
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
	ns, err = cluster.CreateNamespace(ctx, r.client, prefix, s.Privileged)
	if err != nil {
		return fail(fmt.Errorf("creating namespace: %w", err))
	}

	runCtx := &scenario.RunContext{
		Ctx:       ctx,
		Client:    r.client,
		Config:    r.config,
		Namespace: ns,
		IsOCP:     r.isOCP,
	}

	if err := s.Setup(runCtx); err != nil {
		return fail(fmt.Errorf("setup: %w", err))
	}

	resources, err := discovery.Targeted(ctx, r.client, ns)
	if err != nil {
		return fail(fmt.Errorf("discovery: %w", err))
	}
	if r.clusterSnap != nil {
		discovery.ApplyClusterSnapshot(r.clusterSnap, resources)
	}
	if err := discovery.WithClusterRoleBindings(ctx, r.client, resources); err != nil {
		return fail(fmt.Errorf("cluster role bindings discovery: %w", err))
	}

	if s.PostDiscovery != nil {
		s.PostDiscovery(resources)
	}

	checkInfo, ok := checks.ByName(s.CheckName)
	if !ok {
		return fail(fmt.Errorf("check %q not found in registry", s.CheckName))
	}

	checkResult := checkInfo.Fn(resources)
	result.ActualStatus = checkResult.ComplianceStatus
	result.Reason = checkResult.Reason
	result.Details = checkResult.Details
	result.Duration = time.Since(start)

	if s.Verify != nil && result.Passed() {
		if err := s.Verify(checkResult); err != nil {
			result.Err = fmt.Errorf("verify: %w", err)
		}
	}

	return result, ns
}

func (r *Runner) cleanupNamespace(ns string) {
	cleanCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := cluster.DeleteNamespace(cleanCtx, r.client, ns); err != nil {
		r.print(fmt.Sprintf("WARNING: failed to delete namespace %s: %v\n", ns, err))
	}
}

func (r *Runner) shouldSkip(s scenario.Scenario) bool {
	if s.RequiresOCP && (r.opts.SkipOCP || !r.isOCP) {
		return true
	}
	if s.RequiresProbe && r.opts.SkipProbe {
		return true
	}
	return false
}

func (r *Runner) skipReason(s scenario.Scenario) string {
	if s.RequiresOCP && !r.isOCP {
		return "requires OpenShift cluster"
	}
	if s.RequiresOCP && r.opts.SkipOCP {
		return "OCP scenarios skipped (--skip-ocp)"
	}
	if s.RequiresProbe && r.opts.SkipProbe {
		return "probe scenarios skipped (--skip-probe)"
	}
	return "skipped"
}

func (r *Runner) print(msg string) {
	r.printMu.Lock()
	defer r.printMu.Unlock()
	r.printFn(msg)
}

func scenarioPrefix(s scenario.Scenario) string {
	if len(s.Category) > 6 {
		return s.Category[:6]
	}
	return s.Category
}
