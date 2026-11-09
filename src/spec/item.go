package spec

import (
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
	includeRe = regexp.MustCompile(`^\s*- \[([ xX])\] (\d{3,}) (change):`)
	reqRe     = regexp.MustCompile(`^\s*- (R\d+):`)
	reqRefRe  = regexp.MustCompile(`R\d+`)
	dateRe    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

var Statuses = []string{"open", "done"}

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

// Started reports whether work began: an open item with a checked task.
func (it *Item) Started() bool { return it.Status == "open" && it.Done > 0 }

func (it *Item) parse() error {
	data, err := os.ReadFile(it.File)
	if err != nil {
		return err
	}
	it.Lines = strings.Split(string(data), "\n")
	if s, ok := strings.CutPrefix(it.Lines[0], "status: "); ok {
		it.Status = strings.TrimSpace(s)
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
	With []int // update: changes to include
	Bug  bool  // change: Current / Expected / Unchanged behavior form
	Dir  bool  // update: make a folder with design.md
}

// New writes an item from its template and returns the path created.
func New(root string, t Type, title string, opt NewOptions) (string, error) {
	slug := Slugify(title)
	if slug == "" {
		return "", errors.New("the title has no letters or digits")
	}
	if ((len(opt.With) > 0 || opt.Dir) && t.Name != "update") || (opt.Bug && t.Name != "change") {
		return "", errors.New("--with and --dir are for update; --bug is for change")
	}
	tmpl := templates[t.Name]
	if opt.Bug {
		tmpl = templates["change/bug"]
	}
	repl := map[string]string{"{{title}}": title, "{{includes}}": ""}

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

	for _, n := range opt.With {
		it := byNum(items, n)
		switch {
		case it == nil:
			return "", fmt.Errorf("no item %03d", n)
		case it.Archived:
			return "", fmt.Errorf("%03d is already done and archived", n)
		case it.Type.Name != "change":
			return "", fmt.Errorf("%03d is an %s; an update includes only changes", n, it.Type.Name)
		}
		repl["{{includes}}"] += fmt.Sprintf("- [ ] %s %s: %s\n", it.ID(), it.Type.Name, it.Title)
	}

	name := repl["{{num}}"] + "-" + slug
	path := filepath.Join(root, Dir, t.Folder, name+".md")
	if opt.Dir {
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

// Includable lists the active changes an update may include.
func Includable(items []*Item) []*Item {
	var out []*Item
	for _, it := range items {
		if !it.Archived && it.Type.Name == "change" {
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

// DoneOptions control how Done closes an item.
type DoneOptions struct {
	Change   string                 // what closed it; default: git HEAD
	Decision string                 // optional line for notes/decisions.md
	NoVerify bool                   // skip Verify commands
	Run      func(cmd string) error // runs one Verify command in root
}

// Done closes the item and everything it includes, then moves them to the
// archive. It refuses on open markers, unchecked task boxes or a failing
// Verify command.
func Done(root, ref string, opt DoneOptions) ([]*Item, error) {
	g, err := group(root, ref)
	if err != nil {
		return nil, err
	}
	problems := markerProblems(root, g)
	for _, x := range g {
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
	today := time.Now().Format("2006-01-02")
	change = strings.TrimSpace(change + " " + today)
	if d := strings.Join(strings.Fields(opt.Decision), " "); d != "" {
		if err := appendDecision(root, g[0], d, today); err != nil {
			return nil, err
		}
	}
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

// appendDecision adds one line to notes/decisions.md, creating the file.
func appendDecision(root string, it *Item, text, date string) error {
	p := filepath.Join(root, Dir, "notes", "decisions.md")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	line := fmt.Sprintf("- %s %s (%s): %s\n", it.ID(), it.Title, date, text)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		line = "# Decisions\n\nOne line per closed item: what was chosen and why.\n\n" + line
	}
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
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
	it.Status = status
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
