package scenario

import (
	"context"
	"fmt"
	"strings"
	"sync"

	checks "github.com/redhat-best-practices-for-k8s/checks"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/builder"
	"github.com/redhat-best-practices-for-k8s/checks-qe/pkg/cluster"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type RunContext struct {
	Ctx       context.Context
	Client    kubernetes.Interface
	Config    *rest.Config
	Namespace string
	IsOCP     bool
}

type Scenario struct {
	Name           string
	CheckName      string
	Category       string
	Description    string
	ExpectedStatus string
	Privileged     bool
	RequiresOCP    bool
	RequiresProbe  bool
	Tags           []string
	Setup          func(ctx *RunContext) error
	PostDiscovery  func(resources *checks.DiscoveredResources)
	Verify         func(result checks.CheckResult) error
}

var (
	mu        sync.RWMutex
	scenarios []Scenario
)

func Register(ss ...Scenario) {
	mu.Lock()
	defer mu.Unlock()
	for _, s := range ss {
		if s.Name == "" {
			panic("scenario name must not be empty")
		}
		if s.CheckName == "" {
			panic(fmt.Sprintf("scenario %q: check name must not be empty", s.Name))
		}
		if s.ExpectedStatus == "" {
			panic(fmt.Sprintf("scenario %q: expected status must not be empty", s.Name))
		}
		if s.Setup == nil {
			panic(fmt.Sprintf("scenario %q: setup function must not be nil", s.Name))
		}
		scenarios = append(scenarios, s)
	}
}

func All() []Scenario {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Scenario, len(scenarios))
	copy(out, scenarios)
	return out
}

func ByCategory(category string) []Scenario {
	mu.RLock()
	defer mu.RUnlock()
	var out []Scenario
	for _, s := range scenarios {
		if s.Category == category {
			out = append(out, s)
		}
	}
	return out
}

func VanillaDeploymentSetup(ctx *RunContext) error {
	dep := builder.NewDeployment("test-dep", ctx.Namespace).Build()
	return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep, cluster.DefaultTimeout)
}

func TwoDeploymentSetup(decorate func(*builder.DeploymentBuilder) *builder.DeploymentBuilder) func(ctx *RunContext) error {
	return func(ctx *RunContext) error {
		dep1 := decorate(builder.NewDeployment("test-dep-1", ctx.Namespace)).Build()
		if err := cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep1, cluster.DefaultTimeout); err != nil {
			return err
		}
		dep2 := builder.NewDeployment("test-dep-2", ctx.Namespace).Build()
		return cluster.CreateAndWaitForDeployment(ctx.Ctx, ctx.Client, dep2, cluster.DefaultTimeout)
	}
}

func Filtered(pattern string) []Scenario {
	if pattern == "" {
		return All()
	}
	mu.RLock()
	defer mu.RUnlock()
	var out []Scenario
	for _, s := range scenarios {
		if strings.Contains(s.Name, pattern) {
			out = append(out, s)
		}
	}
	return out
}
