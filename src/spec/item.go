package spec

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	nameRe    = regexp.MustCompile(`^(\d{3,})-(.+?)(\.md)?$`)
	markerRe  = regexp.MustCompile(`\[\?:[^\]]*\]`)
	codeRe    = regexp.MustCompile("`[^`]*`")
	boxRe     = regexp.MustCompile(`^\s*- \[([ xX])\] `)
	includeRe = regexp.MustCompile(`^\s*- \[([ xX])\] (\d{3,}) (feature|issue):`)
	reqRe     = regexp.MustCompile(`^\s*- (R\d+):`)
	reqRefRe  = regexp.MustCompile(`R\d+`)
	dateRe    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

var Statuses = []string{"open", "approved", "doing", "done"}

// Marker is one open [?: …] question.
type Marker struct {
	Line int // 1-based
	Text string
}

// Item is one numbered spec file, or a folder whose main file is <type>.md.
type Item struct {
	Num      int
	Type     Type
	Path     string // file or folder
	File     string // the markdown file holding the status line
	Slug     string
	Archived bool

	Status   string
	Hash     string // content hash recorded by approve
	Change   string // what closed it: a commit, PR or path
	Date     string // when it was closed, YYYY-MM-DD
	Title    string
	Markers  []Marker
	Done     int // checked task boxes, Includes lines excluded
	Total    int
	Includes []int
	Verify   []string
	Lines    []string
}

func (it *Item) ID() string { return fmt.Sprintf("%03d", it.Num) }

func (it *Item) hasChangeLine() bool {
	return len(it.Lines) > 1 && strings.HasPrefix(it.Lines[1], "Change:")
}

// Changed reports whether an approved item was edited after approval.
func (it *Item) Changed() bool { return it.Status == "approved" && it.Hash != it.hash() }

// hash covers the content without the status and Change lines and without
// task-box state, so ticking tasks never invalidates an approval.
func (it *Item) hash() string {
	body := it.Lines[1:]
	if it.hasChangeLine() {
		body = it.Lines[2:]
	}
	h := sha256.New()
	for _, line := range body {
		if m := boxRe.FindStringSubmatchIndex(line); m != nil {
			line = line[:m[2]] + " " + line[m[3]:]
		}
		fmt.Fprintln(h, line)
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:8]
}

func (it *Item) parse() error {
	data, err := os.ReadFile(it.File)
	if err != nil {
		return err
	}
	it.Lines = strings.Split(string(data), "\n")
	if s, ok := strings.CutPrefix(it.Lines[0], "status: "); ok {
		f := strings.Fields(s)
		if len(f) > 0 {
			it.Status = f[0]
		}
		if len(f) > 1 {
			it.Hash = f[1]
		}
	}
	if it.hasChangeLine() {
		f := strings.Fields(strings.TrimPrefix(it.Lines[1], "Change:"))
		if n := len(f); n > 0 && dateRe.MatchString(f[n-1]) {
			it.Date, f = f[n-1], f[:n-1]
		}
		it.Change = strings.Join(f, " ")
	}
	fenced := false
	for i, line := range it.Lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		if h, ok := strings.CutPrefix(line, "# "); ok && it.Title == "" {
			it.Title = strings.TrimSpace(strings.TrimPrefix(h, it.ID()))
		}
		if v, ok := strings.CutPrefix(line, "Verify: "); ok && strings.TrimSpace(v) != "" {
			it.Verify = append(it.Verify, strings.TrimSpace(v))
		}
		// Markers quoted as code are examples, not questions.
		for _, m := range markerRe.FindAllString(codeRe.ReplaceAllString(line, ""), -1) {
			it.Markers = append(it.Markers, Marker{i + 1, strings.TrimSpace(m[3 : len(m)-1])})
		}
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

// NewOptions are the extras for some types.
type NewOptions struct {
	With  []int    // update, fix: items to include
	Bugs  []string // fix: small bugs
	Dir   bool     // fix: make a folder with design.md
	Depth string   // feature, issue, plan: "light" or "full" (default)
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
	tmpl := templates[t.Name]
	switch opt.Depth {
	case "", "full":
	case "light":
		var ok bool
		if tmpl, ok = templates[t.Name+"/light"]; !ok {
			return "", errors.New("--depth light is for feature, issue and plan")
		}
	default:
		return "", fmt.Errorf("unknown depth %q, available: light, full", opt.Depth)
	}
	repl := map[string]string{"{{title}}": title, "{{includes}}": "", "{{bugs}}": ""}

	if !t.Status {
		path := filepath.Join(root, Dir, t.Folder, slug+".md")
		return path, writeNew(path, tmpl, repl)
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
	path := filepath.Join(root, Dir, t.Folder, name+".md")
	if t.Name == "update" || opt.Dir {
		path = filepath.Join(root, Dir, t.Folder, name)
		if err := writeNew(filepath.Join(path, "design.md"), templates["design"], repl); err != nil {
			return "", err
		}
		err = writeNew(filepath.Join(path, t.Name+".md"), tmpl, repl)
	} else {
		err = writeNew(path, tmpl, repl)
	}
	if err != nil {
		return "", err
	}
	return path, WriteRoadmap(root)
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

// group resolves ref to the item plus the active items it includes.
func group(root, ref string) ([]*Item, error) {
	items, err := Scan(root)
	if err != nil {
		return nil, err
	}
	it, err := Find(items, ref)
	if err != nil {
		return nil, err
	}
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

func markerProblems(root string, g []*Item) []string {
	var problems []string
	for _, x := range g {
		if n := len(x.Markers); n > 0 {
			problems = append(problems, fmt.Sprintf("%s has %d open [?: …] markers (%s)", x.ID(), n, short(root, x.File)))
		}
	}
	return problems
}

// Approve records the owner's approval on the item and everything it
// includes. It refuses while any of them has [?: …] markers.
func Approve(root, ref string) ([]*Item, error) {
	g, err := group(root, ref)
	if err != nil {
		return nil, err
	}
	problems := markerProblems(root, g)
	for _, x := range g {
		if x.Status != "open" && x.Status != "approved" {
			problems = append(problems, fmt.Sprintf("%s is %s; only open items are approved", x.ID(), x.Status))
		}
	}
	if len(problems) > 0 {
		return nil, gateError("cannot approve", problems)
	}
	for _, x := range g {
		if err := x.write("approved "+x.hash(), ""); err != nil {
			return nil, err
		}
	}
	return g, WriteRoadmap(root)
}

// Start sets doing on the item and everything it includes. Each must be
// approved and unchanged since, or already doing.
func Start(root, ref string) ([]*Item, error) {
	g, err := group(root, ref)
	if err != nil {
		return nil, err
	}
	problems := markerProblems(root, g)
	for _, x := range g {
		switch {
		case x.Status == "doing":
		case x.Status != "approved":
			problems = append(problems, fmt.Sprintf("%s is %s; the owner approves it with `xspec approve %s`", x.ID(), x.Status, x.ID()))
		case x.Changed():
			problems = append(problems, fmt.Sprintf("%s changed after approval; the owner re-approves it with `xspec approve %s`", x.ID(), x.ID()))
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
	return g, WriteRoadmap(root)
}

// DoneOptions control how Done closes an item.
type DoneOptions struct {
	Change   string                 // what closed it; default: git HEAD
	NoVerify bool                   // skip Verify commands
	Run      func(cmd string) error // runs one Verify command in root
}

// Done closes the item and everything it includes, then moves them to the
// archive. It refuses on open markers, unchecked task boxes, an item that was
// never approved, or a failing Verify command.
func Done(root, ref string, opt DoneOptions) ([]*Item, error) {
	g, err := group(root, ref)
	if err != nil {
		return nil, err
	}
	problems := markerProblems(root, g)
	for _, x := range g {
		if x.Status != "doing" && x.Status != "approved" {
			problems = append(problems, fmt.Sprintf("%s is %s; approve and start it first", x.ID(), x.Status))
		}
		if x.Done < x.Total {
			problems = append(problems, fmt.Sprintf("%s has %d of %d boxes unchecked (%s)", x.ID(), x.Total-x.Done, x.Total, short(root, x.File)))
		}
	}
	if len(problems) > 0 {
		return nil, gateError("cannot close", problems)
	}

	if !opt.NoVerify {
		seen := map[string]bool{}
		for _, x := range g {
			for _, cmd := range x.Verify {
				if seen[cmd] {
					continue
				}
				seen[cmd] = true
				if err := opt.Run(cmd); err != nil {
					return nil, fmt.Errorf("cannot close: %s Verify `%s` failed: %w", x.ID(), cmd, err)
				}
			}
		}
	}

	change := opt.Change
	if change == "" {
		if out, err := exec.Command("git", "-C", root, "rev-parse", "--short", "HEAD").Output(); err == nil {
			change = strings.TrimSpace(string(out))
		}
	}
	change = strings.TrimSpace(change + " " + time.Now().Format("2006-01-02"))
	for i, line := range g[0].Lines {
		if includeRe.MatchString(line) {
			g[0].Lines[i] = strings.Replace(line, "- [ ] ", "- [x] ", 1)
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
	return g, WriteRoadmap(root)
}

// Shell runs a Verify command in dir with the terminal attached.
func Shell(dir string) func(string) error {
	return func(command string) error {
		c := exec.Command("sh", "-c", command)
		if runtime.GOOS == "windows" {
			c = exec.Command("cmd", "/C", command)
		}
		c.Dir, c.Stdin, c.Stdout, c.Stderr = dir, os.Stdin, os.Stdout, os.Stderr
		return c.Run()
	}
}

// write sets the status line and, when change is set, the Change line.
func (it *Item) write(status, change string) error {
	if !strings.HasPrefix(it.Lines[0], "status: ") {
		it.Lines = append([]string{""}, it.Lines...)
	}
	it.Lines[0] = "status: " + status
	if change != "" {
		if it.hasChangeLine() {
			it.Lines[1] = "Change: " + change
		} else {
			it.Lines = append(it.Lines[:1], append([]string{"Change: " + change}, it.Lines[1:]...)...)
		}
	}
	it.Status, it.Hash, _ = strings.Cut(status, " ")
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
