package spec

import (
	"fmt"
	"slices"
)

// Check returns every problem found in spec/; none means the spec is healthy.
func Check(root string) ([]string, error) {
	var problems []string
	if ign, ok := ignored(root); ok && !ign {
		problems = append(problems, "spec/ is not ignored by git; run `xspec init` or add it to .gitignore")
	}
	items, err := Scan(root)
	if err != nil {
		return nil, err
	}
	seen := map[int]string{}
	for _, it := range items {
		add := func(format string, a ...any) {
			problems = append(problems, short(root, it.File)+": "+fmt.Sprintf(format, a...))
		}
		if prev, ok := seen[it.Num]; ok {
			add("number %s is also used by %s", it.ID(), prev)
		}
		seen[it.Num] = it.File

		switch {
		case !slices.Contains(Statuses, it.Status):
			add("line 1 must be `status: open` or `status: done`")
		case it.Status == "done" && !it.Archived:
			add("done but not archived; close it with `xspec done %s`", it.ID())
		case it.Status != "done" && it.Archived:
			add("archived but status is %s", it.Status)
		case it.Started() && len(it.Markers) > 0:
			add("work began with %d open [?: …] markers; ask the owner and answer them", len(it.Markers))
		}
		for _, n := range it.Includes {
			if byNum(items, n) == nil {
				add("includes %03d, which does not exist", n)
			}
		}
		for _, r := range uncovered(it.Lines) {
			add("requirement %s has no task referencing it", r)
		}
	}
	return problems, nil
}

// uncovered lists requirement IDs (R1, R2…) that no task box mentions.
func uncovered(lines []string) []string {
	var reqs []string
	refs := map[string]bool{}
	for _, line := range lines {
		if m := reqRe.FindStringSubmatch(line); m != nil {
			reqs = append(reqs, m[1])
		} else if boxRe.MatchString(line) {
			for _, r := range reqRefRe.FindAllString(line, -1) {
				refs[r] = true
			}
		}
	}
	var out []string
	for _, r := range reqs {
		if !refs[r] {
			out = append(out, r)
		}
	}
	return out
}
