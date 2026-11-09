package spec_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/exec-cmd/execspec/src/spec"
)

// setup makes an initialised spec/ with no items; numbering starts at 001.
func setup(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if _, err := spec.Init(root); err != nil {
		t.Fatal(err)
	}
	return root
}

func mustType(t *testing.T, name string) spec.Type {
	t.Helper()
	typ, err := spec.TypeByName(name)
	if err != nil {
		t.Fatal(err)
	}
	return typ
}

func newItem(t *testing.T, root, typ, title string, opt spec.NewOptions) string {
	t.Helper()
	p, err := spec.New(root, mustType(t, typ), title, opt)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// answer replaces every [?: …] marker and, with tick, checks every box.
func answer(t *testing.T, path string, tick bool) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	for strings.Contains(s, "[?:") {
		i := strings.Index(s, "[?:")
		j := strings.Index(s[i:], "]")
		s = s[:i] + "true" + s[i+j+1:]
	}
	if tick {
		s = strings.ReplaceAll(s, "- [ ] ", "- [x] ")
	}
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

func noVerify() spec.DoneOptions { return spec.DoneOptions{Change: "abc123", NoVerify: true} }

func scan(t *testing.T, root string) []*spec.Item {
	t.Helper()
	items, err := spec.Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	return items
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Login form":           "login-form",
		"  v0.2 — big update ": "v0-2-big-update",
		"Форма входа":          "форма-входа",
		"!!!":                  "",
	}
	for in, want := range cases {
		if got := spec.Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestInitIsIdempotentAndIgnoresSpec(t *testing.T) {
	root := t.TempDir()
	if _, err := spec.Init(root); err != nil {
		t.Fatal(err)
	}
	created, err := spec.Init(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 0 {
		t.Fatalf("second init created %v", created)
	}
	data, _ := os.ReadFile(filepath.Join(root, ".gitignore"))
	if string(data) != "/spec/\n" {
		t.Fatalf(".gitignore = %q, want one /spec/ line", data)
	}
	for _, f := range []string{"SPEC.md", "ROADMAP.md", "notes/glossary.md", "archive", "personal", "changes", "updates"} {
		if _, err := os.Stat(filepath.Join(root, "spec", f)); err != nil {
			t.Error(err)
		}
	}
	brief, _ := os.ReadFile(filepath.Join(root, "spec/SPEC.md"))
	for _, h := range []string{"## Context", "## Audience", "## Principles", "## Design", "## Stack", "## Ask first", "- add a dependency"} {
		if !strings.Contains(string(brief), h) {
			t.Errorf("SPEC.md stub misses %q", h)
		}
	}
	if items := scan(t, root); len(items) != 0 {
		t.Fatalf("want no items, got %d", len(items))
	}
}

func TestNumbersCountTheArchive(t *testing.T) {
	root := setup(t)
	p := newItem(t, root, "change", "One", spec.NewOptions{})
	answer(t, p, true)
	if _, err := spec.Done(root, "1", noVerify()); err != nil {
		t.Fatal(err)
	}
	if p := newItem(t, root, "change", "Two", spec.NewOptions{}); filepath.Base(p) != "002-two.md" {
		t.Fatalf("got %s, want 002-two.md", p)
	}
	if n := newItem(t, root, "note", "Free note", spec.NewOptions{}); filepath.Base(n) != "free-note.md" {
		t.Fatalf("note got %s", n)
	}
}

func TestNewTypesAndForms(t *testing.T) {
	root := setup(t)
	for _, old := range []string{"feature", "issue", "fix", "plan"} {
		if _, err := spec.TypeByName(old); err == nil || !strings.Contains(err.Error(), "change") || !strings.Contains(err.Error(), "update") {
			t.Errorf("%s: got %v, want an error naming change and update", old, err)
		}
	}
	bug := newItem(t, root, "change", "Crash", spec.NewOptions{Bug: true})
	data, _ := os.ReadFile(bug)
	for _, want := range []string{"## Current behavior", "## Unchanged behavior", "SHALL CONTINUE TO", "Verify:"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("bug form misses %q", want)
		}
	}
	plain := newItem(t, root, "change", "Dark mode", spec.NewOptions{})
	data, _ = os.ReadFile(plain)
	if strings.Contains(string(data), "Unchanged") || !strings.Contains(string(data), "## Goal") {
		t.Fatalf("not the plain form:\n%s", data)
	}
	dir := newItem(t, root, "update", "v0.2", spec.NewOptions{With: []int{1, 2}, Dir: true})
	for _, f := range []string{"update.md", "design.md"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Error(err)
		}
	}
	data, _ = os.ReadFile(filepath.Join(dir, "update.md"))
	if !strings.Contains(string(data), "- [ ] 001 change: Crash") || !strings.Contains(string(data), "- [ ] 002 change: Dark mode") {
		t.Fatalf("includes:\n%s", data)
	}
	if p := newItem(t, root, "update", "Small", spec.NewOptions{}); filepath.Ext(p) != ".md" {
		t.Fatalf("update without --dir should be one file, got %s", p)
	}
}

func TestNewRejects(t *testing.T) {
	root := setup(t)
	newItem(t, root, "change", "One", spec.NewOptions{})
	newItem(t, root, "update", "Batch", spec.NewOptions{})
	for name, tc := range map[string]struct {
		typ string
		opt spec.NewOptions
	}{
		"update in update": {"update", spec.NewOptions{With: []int{2}}},
		"missing item":     {"update", spec.NewOptions{With: []int{9}}},
		"with on change":   {"change", spec.NewOptions{With: []int{1}}},
		"dir on change":    {"change", spec.NewOptions{Dir: true}},
		"bug on update":    {"update", spec.NewOptions{Bug: true}},
	} {
		if _, err := spec.New(root, mustType(t, tc.typ), name, tc.opt); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}

func TestMarkersInCodeAreNotQuestions(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "change", "Login", spec.NewOptions{})
	body := "status: open\n# 001 Login\n\nuse `[?: …]` markers\n```\n[?: example]\n```\n[?: real]\n"
	if err := os.WriteFile(f, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	items := scan(t, root)
	if m := items[0].Markers; len(m) != 1 || m[0].Line != 8 || m[0].Text != "real" {
		t.Fatalf("got markers %+v, want one on line 8", m)
	}
}

func TestDoneGates(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "change", "Login", spec.NewOptions{})
	if _, err := spec.Done(root, "1", noVerify()); err == nil || !strings.Contains(err.Error(), "markers") {
		t.Fatalf("done with markers: got %v", err)
	}
	answer(t, f, false)
	if _, err := spec.Done(root, "1", noVerify()); err == nil || !strings.Contains(err.Error(), "unchecked") {
		t.Fatalf("done with open boxes: got %v", err)
	}
	answer(t, f, true)
	if _, err := spec.Done(root, "1", noVerify()); err != nil {
		t.Fatalf("done without approval: %v", err)
	}
}

func TestDecisionLine(t *testing.T) {
	root := setup(t)
	answer(t, newItem(t, root, "change", "Login", spec.NewOptions{}), true)
	answer(t, newItem(t, root, "change", "Logout", spec.NewOptions{}), true)
	opt := noVerify()
	opt.Decision = "  sessions\nover JWT:  simpler  "
	if _, err := spec.Done(root, "1", opt); err != nil {
		t.Fatal(err)
	}
	if _, err := spec.Done(root, "2", noVerify()); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(root, "spec/notes/decisions.md"))
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if last := lines[len(lines)-1]; !strings.HasPrefix(last, "- 001 Login (") || !strings.HasSuffix(last, "): sessions over JWT: simpler") {
		t.Fatalf("decisions.md:\n%s", data)
	}
	if strings.Contains(string(data), "Logout") {
		t.Fatal("closing without --decision wrote a line")
	}
}

func TestUpdateCascade(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "change", "Login", spec.NewOptions{})
	i := newItem(t, root, "change", "Crash", spec.NewOptions{Bug: true})
	u := newItem(t, root, "update", "v0.2", spec.NewOptions{With: []int{1, 2}, Dir: true})
	main := filepath.Join(u, "update.md")
	for _, p := range []string{f, i, main} {
		answer(t, p, false)
	}
	if _, err := spec.Done(root, "3", noVerify()); err == nil || !strings.Contains(err.Error(), "unchecked") {
		t.Fatalf("got %v, want an unchecked-boxes error", err)
	}
	for _, p := range []string{f, i, main} {
		answer(t, p, true)
	}
	if g, err := spec.Done(root, "3", noVerify()); err != nil || len(g) != 3 {
		t.Fatalf("done: %d items, %v", len(g), err)
	}

	for _, it := range scan(t, root) {
		if !it.Archived || it.Status != "done" || it.Change != "abc123" || it.Date == "" {
			t.Errorf("%s: archived=%v status=%s change=%q date=%q", it.ID(), it.Archived, it.Status, it.Change, it.Date)
		}
	}
	upd, _ := os.ReadFile(filepath.Join(root, "spec/archive/updates/003-v0-2/update.md"))
	if !strings.Contains(string(upd), "- [x] 001 change: Login") {
		t.Fatalf("includes not ticked:\n%s", upd)
	}
	if problems, _ := spec.Check(root); len(problems) != 0 {
		t.Fatalf("check after done: %v", problems)
	}
	roadmap, _ := os.ReadFile(filepath.Join(root, "spec/ROADMAP.md"))
	if !strings.Contains(string(roadmap), "003 update v0.2 (abc123)") || !strings.Contains(string(roadmap), "  - 001 change Login") {
		t.Fatalf("roadmap history:\n%s", roadmap)
	}
}

func TestVerifyRunsBeforeDone(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "change", "Login", spec.NewOptions{})
	answer(t, f, true) // the Verify marker becomes "Verify: true"
	var ran []string
	fail := spec.DoneOptions{Run: func(c string) error { ran = append(ran, c); return errors.New("exit 1") }}
	if _, err := spec.Done(root, "1", fail); err == nil || !strings.Contains(err.Error(), "Verify") {
		t.Fatalf("got %v, want a Verify failure", err)
	}
	if len(ran) != 1 || ran[0] != "true" {
		t.Fatalf("ran %v, want [true]", ran)
	}
	if items := scan(t, root); items[0].Archived {
		t.Fatal("closed despite the failing Verify")
	}
	ok := spec.DoneOptions{Run: func(string) error { return nil }}
	if _, err := spec.Done(root, "1", ok); err != nil {
		t.Fatal(err)
	}
}

func TestProgressNextAndQuestions(t *testing.T) {
	root := setup(t)
	if got := spec.Next(scan(t, root)); got != "nothing to do" {
		t.Fatalf("empty board: %q", got)
	}
	f := newItem(t, root, "change", "Login", spec.NewOptions{})
	if got := spec.Next(scan(t, root)); !strings.HasPrefix(got, "answer 3 questions in 001") {
		t.Fatalf("open with markers: %q", got)
	}
	q := spec.Questions(root, scan(t, root))
	if !strings.Contains(q, "spec/changes/001-login.md:5  what changes") {
		t.Fatalf("questions:\n%s", q)
	}
	answer(t, f, false)
	if got := spec.Next(scan(t, root)); !strings.HasPrefix(got, "work on 001 Login: T1") {
		t.Fatalf("answered: %q", got)
	}

	// a second task makes the item "in progress" once one is ticked
	data, _ := os.ReadFile(f)
	s := strings.Replace(string(data), "- [ ] T1 true\n", "- [x] T1 true\n- [ ] T2 more\n", 1)
	os.WriteFile(f, []byte(s), 0o644)
	items := scan(t, root)
	if !items[0].Started() {
		t.Fatal("a ticked task should mean in progress")
	}
	if got := spec.Next(items); !strings.HasPrefix(got, "work on 001 Login: T2 more") {
		t.Fatalf("in progress: %q", got)
	}
	if r := spec.Roadmap(items); !strings.Contains(r, "## Now\n\n- 001 change Login (open, 1/2 tasks)") {
		t.Fatalf("roadmap Now:\n%s", r)
	}
	answer(t, f, true)
	if got := spec.Next(scan(t, root)); !strings.HasPrefix(got, "close 001 Login: xspec done 001") {
		t.Fatalf("all ticked: %q", got)
	}
}

func TestContext(t *testing.T) {
	root := setup(t)
	newItem(t, root, "change", "Login", spec.NewOptions{})
	specMD := filepath.Join(root, "spec/SPEC.md")
	brief := "# p\n\n## Principles\nNo network.\n\n## Stack\n<!-- hint -->\nGo.\n\n## Ask first\n- add a dependency\n\n## Design\nSecret.\n"
	if err := os.WriteFile(specMD, []byte(brief), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := spec.Context(root, "1")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"=== spec/changes/001-login.md ===", "§ Principles ===\nNo network.", "§ Stack ===\nGo.", "§ Ask first ===\n- add a dependency", "=== spec/notes/glossary.md ==="} {
		if !strings.Contains(out, want) {
			t.Errorf("context missed %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "hint") || strings.Contains(out, "Secret.") {
		t.Errorf("context carries a comment or an unrequested section:\n%s", out)
	}
}

func TestCheckFindsProblems(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "change", "Login", spec.NewOptions{})
	data, _ := os.ReadFile(f)
	s := strings.Replace(string(data), "status: open", "status: approved abc12345", 1)
	if err := os.WriteFile(f, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
	g := newItem(t, root, "change", "Logout", spec.NewOptions{})
	data, _ = os.ReadFile(g)
	s = strings.Replace(string(data), "- [ ] T1", "- [x] T1", 1) + "- R1: WHEN a THE SYSTEM SHALL b\n"
	if err := os.WriteFile(g, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
	problems, err := spec.Check(root)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(problems, "\n")
	for _, want := range []string{"line 1 must be", "work began with", "R1 has no task"} {
		if !strings.Contains(joined, want) {
			t.Errorf("check missed %q in:\n%s", want, joined)
		}
	}
}
