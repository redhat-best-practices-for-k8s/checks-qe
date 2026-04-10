package runner

import (
	"fmt"
	"strings"
	"time"

	checks "github.com/redhat-best-practices-for-k8s/checks"
)

type ScenarioResult struct {
	Name           string
	CheckName      string
	ExpectedStatus string
	ActualStatus   string
	Reason         string
	Details        []checks.ResourceDetail
	Duration       time.Duration
	Err            error
	Skipped        bool
	SkipReason     string
}

func (r *ScenarioResult) Passed() bool {
	if r.Err != nil || r.Skipped {
		return false
	}
	return r.ActualStatus == r.ExpectedStatus
}

type Summary struct {
	Total    int
	Passed   int
	Failed   int
	Skipped  int
	Errors   int
	Duration time.Duration
}

func (s Summary) String() string {
	return fmt.Sprintf("Results: %d passed, %d failed, %d skipped, %d errors (%s)",
		s.Passed, s.Failed, s.Skipped, s.Errors, s.Duration.Round(time.Millisecond))
}

func FormatResult(r *ScenarioResult, verbose bool) string {
	var sb strings.Builder

	if r.Skipped {
		fmt.Fprintf(&sb, "=== SKIP  %s\n", r.Name)
		fmt.Fprintf(&sb, "    %s\n", r.SkipReason)
		return sb.String()
	}

	fmt.Fprintf(&sb, "=== RUN   %s\n", r.Name)
	fmt.Fprintf(&sb, "    check: %s\n", r.CheckName)

	if r.Err != nil {
		fmt.Fprintf(&sb, "--- ERROR (%s)\n", r.Duration.Round(time.Millisecond))
		fmt.Fprintf(&sb, "    %v\n", r.Err)
		return sb.String()
	}

	if r.Passed() {
		fmt.Fprintf(&sb, "--- PASS  (%s)", r.Duration.Round(time.Millisecond))
		if r.ExpectedStatus != checks.StatusCompliant {
			fmt.Fprintf(&sb, "  [expected %s]", r.ExpectedStatus)
		}
		sb.WriteString("\n")
	} else {
		fmt.Fprintf(&sb, "--- FAIL  (%s)  [expected %s, got %s]\n",
			r.Duration.Round(time.Millisecond), r.ExpectedStatus, r.ActualStatus)
		fmt.Fprintf(&sb, "    reason: %q\n", r.Reason)
		if verbose || len(r.Details) > 0 {
			fmt.Fprintf(&sb, "    details: %s\n", formatDetails(r.Details))
		}
	}

	return sb.String()
}

func formatDetails(details []checks.ResourceDetail) string {
	if len(details) == 0 {
		return "(none)"
	}
	var parts []string
	for _, d := range details {
		status := "OK"
		if !d.Compliant {
			status = "FAIL"
		}
		parts = append(parts, fmt.Sprintf("[%s] %s %s/%s: %s",
			status, d.Kind, d.Namespace, d.Name, d.Message))
	}
	return strings.Join(parts, "; ")
}
