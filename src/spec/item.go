package spec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var (
	nameRe    = regexp.MustCompile(`^(\d{3,})-(.+?)(\.md)?$`)
	markerRe  = regexp.MustCompile(`\[\?:[^\]]*\]`)
	boxRe     = regexp.MustCompile(`^\s*- \[([ xX])\] `)
	includeRe = regexp.MustCompile(`^\s*- \[([ xX])\] (\d{3,}) (feature|issue):`)
	reqRe     = regexp.MustCompile(`^\s*- (R\d+):`)
	reqRefRe  = regexp.MustCompile(`R\d+`)
)

var Statuses = []string{"open", "doing", "done"}

// Item is one numbered spec file, or a folder whose main file is <type>.md.
type Item struct {
	Num      int
	Type     Type
	Path     string // file or folder
	File     string // the markdown file holding the status line
	Slug     string
	Archived bool

	Status   string
	Change   string
	Title    string
	Markers  int
	Done     int // checked task boxes, Includes lines excluded
	Total    int
	Includes []int
	Lines    []string
}

func (it *Item) ID() string { return fmt.Sprintf("%03d", it.Num) }

func (it *Item) parse() error {
	data, err := os.ReadFile(it.File)
	if err != nil {
		return err
	}
	it.Lines = strings.Split(string(data), "\n")
	if s, ok := strings.CutPrefix(it.Lines[0], "status: "); ok {
		it.Status = strings.TrimSpace(s)
	}
	for i, line := range it.Lines {
		if c, ok := strings.CutPrefix(line, "Change: "); ok && i == 1 {
			it.Change = c
		}
		if h, ok := strings.CutPrefix(line, "# "); ok && it.Title == "" {
			it.Title = strings.TrimSpace(strings.TrimPrefix(h, it.ID()))
		}
		it.Markers += len(markerRe.FindAllString(line, -1))
		if m := includeRe.FindStringSubmatch(line); m != nil {
			n, _ := strconv.Atoi(m[2])
			it.Includes = append(it.Includes, n)
			continue
		}
		if m := boxRe.FindStringSubmatch(line); m != nil {
			it.Total++
			if m[1] != " " {
				it.Done++
			}
		}
	}
	return nil
}

// Scan reads every numbered item, active and archived, sorted by number.
func Scan(root string) ([]*Item, error) {
	var items []*Item
	for _, t := range Types {
		if !t.Status {
			continue
		}
		for _, base := range []string{filepath.Join(root, Dir), filepath.Join(root, Dir, Archive)} {
			entries, err := os.ReadDir(filepath.Join(base, t.Folder))
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return nil, err
			}
			for _, e := range entries {
				m := nameRe.FindStringSubmatch(e.Name())
				if m == nil || e.IsDir() == (m[3] != "") {
					continue
				}
				n, _ := strconv.Atoi(m[1])
				it := &Item{
					Num: n, Type: t, Slug: m[2],
					Path:     filepath.Join(base, t.Folder, e.Name()),
					Archived: base != filepath.Join(root, Dir),
				}
				it.File = it.Path
				if e.IsDir() {
					it.File = filepath.Join(it.Path, t.Name+".md")
				}
				if err := it.parse(); err != nil {
					return nil, err
				}
				items = append(items, it)
			}
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Num < items[j].Num })
	return items, nil
}

// Find resolves a ref: a number (7, 007) or a slug.
func Find(items []*Item, ref string) (*Item, error) {
	var found []*Item
	n, numErr := strconv.Atoi(ref)
	for _, it := range items {
		if (numErr == nil && it.Num == n) || it.Slug == ref {
			found = append(found, it)
		}
	}
	switch len(found) {
	case 0:
		return nil, fmt.Errorf("no item %q", ref)
	case 1:
		return found[0], nil
	}
	return nil, fmt.Errorf("%q matches %d items; use the number", ref, len(found))
}

func byNum(items []*Item, n int) *Item {
	for _, it := range items {
		if it.Num == n {
			return it
		}
	}
	return nil
}

func Slugify(title string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(title) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	s := strings.TrimSuffix(b.String(), "-")
	if r := []rune(s); len(r) > 50 {
		s = strings.TrimSuffix(string(r[:50]), "-")
	}
	return s
}

// NewOptions are the extras for update and fix.
type NewOptions struct {
	With []int    // items to include
	Bugs []string // fix: small bugs
	Dir  bool     // fix: make a folder with design.md
}

// New writes an item from its template and returns the path created.
func New(root string, t Type, title string, opt NewOptions) (string, error) {
	slug := Slugify(title)
	if slug == "" {
		return "", errors.New("the title has no letters or digits")
	}
	if (len(opt.With) > 0 && t.Name != "update" && t.Name != "fix") ||
		((len(opt.Bugs) > 0 || opt.Dir) && t.Name != "fix") {
		return "", fmt.Errorf("--with is for update and fix; --bug and --dir are for fix")
	}
	repl := map[string]string{"{{title}}": title, "{{includes}}": "", "{{bugs}}": ""}

	if !t.Status {
		path := filepath.Join(root, Dir, t.Folder, slug+".md")
		return path, writeNew(path, templates[t.Name], repl)
	}

	items, err := Scan(root)
	if err != nil {
		return "", err
	}
	num := 1
	if len(items) > 0 {
		num = items[len(items)-1].Num + 1
	}
	repl["{{num}}"] = fmt.Sprintf("%03d", num)

	allowed := map[string]bool{"feature": t.Name == "update", "issue": true}
	for _, n := range opt.With {
		it := byNum(items, n)
		switch {
		case it == nil:
			return "", fmt.Errorf("no item %03d", n)
		case it.Archived:
			return "", fmt.Errorf("%03d is already done and archived", n)
		case !allowed[it.Type.Name]:
			return "", fmt.Errorf("%03d is a %s; %s includes only %s", n, it.Type.Name, t.Name, includable(t.Name))
		}
		repl["{{includes}}"] += fmt.Sprintf("- [ ] %s %s: %s\n", it.ID(), it.Type.Name, it.Title)
	}
	for _, b := range opt.Bugs {
		repl["{{bugs}}"] += "- [ ] " + b + "\n"
	}

	name := repl["{{num}}"] + "-" + slug
	dir := filepath.Join(root, Dir, t.Folder)
	if t.Name == "update" || opt.Dir {
		dir = filepath.Join(dir, name)
		if err := writeNew(filepath.Join(dir, "design.md"), templates["design"], repl); err != nil {
			return "", err
		}
		return dir, writeNew(filepath.Join(dir, t.Name+".md"), templates[t.Name], repl)
	}
	path := filepath.Join(dir, name+".md")
	return path, writeNew(path, templates[t.Name], repl)
}

func includable(typ string) string {
	if typ == "update" {
		return "features and issues"
	}
	return "issues"
}

// Includable lists the active items an update or fix may include.
func Includable(items []*Item, typ string) []*Item {
	var out []*Item
	for _, it := range items {
		if !it.Archived && (it.Type.Name == "issue" || (typ == "update" && it.Type.Name == "feature")) {
			out = append(out, it)
		}
	}
	return out
}

func writeNew(path, tmpl string, repl map[string]string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}
	for k, v := range repl {
		tmpl = strings.ReplaceAll(tmpl, k, v)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(tmpl), 0o644)
}

// group is the item plus the active items it includes.
func group(items []*Item, it *Item) ([]*Item, error) {
	if it.Archived {
		return nil, fmt.Errorf("%s is done and archived", it.ID())
	}
	g := []*Item{it}
	for _, n := range it.Includes {
		inc := byNum(items, n)
		if inc == nil {
			return nil, fmt.Errorf("%s includes %03d, which does not exist", it.ID(), n)
		}
		if !inc.Archived {
			g = append(g, inc)
		}
	}
	return g, nil
}

// Start sets doing on the item and everything it includes. It refuses while
// any of them still has [?: …] markers.
func Start(root, ref string) ([]*Item, error) {
	items, err := Scan(root)
	if err != nil {
		return nil, err
	}
	it, err := Find(items, ref)
	if err != nil {
		return nil, err
	}
	g, err := group(items, it)
	if err != nil {
		return nil, err
	}
	var problems []string
	for _, x := range g {
		if x.Markers > 0 {
			problems = append(problems, fmt.Sprintf("%s has %d open [?: …] markers (%s)", x.ID(), x.Markers, short(root, x.File)))
		}
	}
	if len(problems) > 0 {
		return nil, gateError("cannot start", problems)
	}
	for _, x := range g {
		if err := x.write("doing", ""); err != nil {
			return nil, err
		}
	}
	return g, nil
}

// Done closes the item and everything it includes, then moves them to the
// archive. It refuses on open markers or unchecked task boxes.
func Done(root, ref, change string) ([]*Item, error) {
	items, err := Scan(root)
	if err != nil {
		return nil, err
	}
	it, err := Find(items, ref)
	if err != nil {
		return nil, err
	}
	g, err := group(items, it)
	if err != nil {
		return nil, err
	}
	var problems []string
	for _, x := range g {
		if x.Markers > 0 {
			problems = append(problems, fmt.Sprintf("%s has %d open [?: …] markers", x.ID(), x.Markers))
		}
		if x.Done < x.Total {
			problems = append(problems, fmt.Sprintf("%s has %d of %d boxes unchecked (%s)", x.ID(), x.Total-x.Done, x.Total, short(root, x.File)))
		}
	}
	if len(problems) > 0 {
		return nil, gateError("cannot close", problems)
	}

	if change == "" {
		if out, err := exec.Command("git", "-C", root, "rev-parse", "--short", "HEAD").Output(); err == nil {
			change = strings.TrimSpace(string(out))
		}
	}
	for i, line := range it.Lines {
		if includeRe.MatchString(line) {
			it.Lines[i] = strings.Replace(line, "- [ ] ", "- [x] ", 1)
		}
	}
	for _, x := range g {
		if err := x.write("done", change); err != nil {
			return nil, err
		}
		dst := filepath.Join(root, Dir, Archive, x.Type.Folder, filepath.Base(x.Path))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return nil, err
		}
		if _, err := os.Stat(dst); err == nil {
			return nil, fmt.Errorf("%s already exists", dst)
		}
		if err := os.Rename(x.Path, dst); err != nil {
			return nil, err
		}
		x.Path, x.Archived = dst, true
	}
	return g, nil
}

// write sets the status line and, when change is set, the Change line.
func (it *Item) write(status, change string) error {
	if it.Status == "" {
		it.Lines = append([]string{""}, it.Lines...)
	}
	it.Lines[0] = "status: " + status
	it.Status = status
	if change != "" {
		if it.Change != "" {
			it.Lines[1] = "Change: " + change
		} else {
			it.Lines = append(it.Lines[:1], append([]string{"Change: " + change}, it.Lines[1:]...)...)
		}
		it.Change = change
	}
	return os.WriteFile(it.File, []byte(strings.Join(it.Lines, "\n")), 0o644)
}

func gateError(what string, problems []string) error {
	return fmt.Errorf("%s:\n  %s", what, strings.Join(problems, "\n  "))
}

// short is p relative to the project root, for messages.
func short(root, p string) string {
	if r, err := filepath.Rel(root, p); err == nil {
		return r
	}
	return p
}
