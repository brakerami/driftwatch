// Package rollup aggregates drift findings across multiple containers
// into a summarised report grouped by drift type and severity.
package rollup

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Summary holds aggregated drift statistics for a single run.
type Summary struct {
	TotalContainers int
	DriftedContainers int
	ByType map[string]int
	TopOffenders []ContainerCount
}

// ContainerCount pairs a container name with its finding count.
type ContainerCount struct {
	Name  string
	Count int
}

// String returns a human-readable summary.
func (s Summary) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Containers checked: %d, drifted: %d\n", s.TotalContainers, s.DriftedContainers)
	fmt.Fprintf(&sb, "Findings by type:\n")
	for _, k := range sortedKeys(s.ByType) {
		fmt.Fprintf(&sb, "  %-20s %d\n", k, s.ByType[k])
	}
	if len(s.TopOffenders) > 0 {
		fmt.Fprintf(&sb, "Top offenders:\n")
		for _, c := range s.TopOffenders {
			fmt.Fprintf(&sb, "  %-30s %d finding(s)\n", c.Name, c.Count)
		}
	}
	return sb.String()
}

// Aggregate builds a Summary from a map of container name → findings.
func Aggregate(results map[string][]drift.Finding) Summary {
	s := Summary{
		TotalContainers: len(results),
		ByType:          make(map[string]int),
	}

	counts := make([]ContainerCount, 0, len(results))
	for name, findings := range results {
		if len(findings) == 0 {
			continue
		}
		s.DriftedContainers++
		counts = append(counts, ContainerCount{Name: name, Count: len(findings)})
		for _, f := range findings {
			s.ByType[string(f.Type)]++
		}
	}

	sort.Slice(counts, func(i, j int) bool {
		if counts[i].Count != counts[j].Count {
			return counts[i].Count > counts[j].Count
		}
		return counts[i].Name < counts[j].Name
	})

	const maxOffenders = 5
	if len(counts) > maxOffenders {
		counts = counts[:maxOffenders]
	}
	s.TopOffenders = counts
	return s
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
