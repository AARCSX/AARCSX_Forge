package doctor

import (
	"context"
)

// Runner executes a collection of checks and reports results
type Runner struct {
	checks []Check
}

// NewRunner creates a new check runner with the given checks
func NewRunner(checks ...Check) *Runner {
	return &Runner{
		checks: checks,
	}
}

// Run executes all checks and returns the results
func (r *Runner) Run(ctx context.Context) []CheckResult {
	results := make([]CheckResult, 0, len(r.checks))

	for _, check := range r.checks {
		result := check.Check(ctx)
		results = append(results, result)
	}

	return results
}

// Summary returns a summary of check results
func Summary(results []CheckResult) (passed int, failed int) {
	for _, r := range results {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}
	return passed, failed
}