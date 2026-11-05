package spec_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/exec-cmd/execspec/src/spec"
)

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

// answer replaces every [?: …] marker and ticks every box in path.
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
		s = s[:i] + "done" + s[i+j+1:]
	}
	if tick {
		s = strings.ReplaceAll(s, "- [ ] ", "- [x] ")
	}
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
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
	root := setup(t)
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
	for _, f := range []string{"SPEC.md", "notes/glossary.md", "archive", "personal", "updates"} {
		if _, err := os.Stat(filepath.Join(root, "spec", f)); err != nil {
			t.Error(err)
		}
	}
}

func TestNumbersCountTheArchive(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "feature", "One", spec.NewOptions{})
	answer(t, f, true)
	if _, err := spec.Done(root, "1", "abc"); err != nil {
		t.Fatal(err)
	}
	p := newItem(t, root, "issue", "Two", spec.NewOptions{})
	if filepath.Base(p) != "002-two.md" {
		t.Fatalf("got %s, want 002-two.md", p)
	}
	n := newItem(t, root, "note", "Free note", spec.NewOptions{})
	if filepath.Base(n) != "free-note.md" {
		t.Fatalf("note got %s", n)
	}
}

func TestNewRejects(t *testing.T) {
	root := setup(t)
	newItem(t, root, "fix", "Typos", spec.NewOptions{})
	for name, tc := range map[string]struct {
		typ string
		opt spec.NewOptions
	}{
		"fix in update": {"update", spec.NewOptions{With: []int{1}}},
		"missing item":  {"update", spec.NewOptions{With: []int{9}}},
		"with on plan":  {"plan", spec.NewOptions{With: []int{1}}},
		"bug on update": {"update", spec.NewOptions{Bugs: []string{"x"}}},
	} {
		if _, err := spec.New(root, mustType(t, tc.typ), name, tc.opt); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}

func TestStartRefusesOnMarkers(t *testing.T) {
	root := setup(t)
	newItem(t, root, "feature", "Login", spec.NewOptions{})
	if _, err := spec.Start(root, "login"); err == nil || !strings.Contains(err.Error(), "markers") {
		t.Fatalf("got %v, want a markers error", err)
	}
}

func TestUpdateCascade(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "feature", "Login", spec.NewOptions{})
	i := newItem(t, root, "issue", "Crash", spec.NewOptions{})
	u := newItem(t, root, "update", "v0.2", spec.NewOptions{With: []int{1, 2}})
	for _, p := range []string{f, i, filepath.Join(u, "update.md")} {
		answer(t, p, false)
	}

	g, err := spec.Start(root, "3")
	if err != nil {
		t.Fatal(err)
	}
	if len(g) != 3 {
		t.Fatalf("start touched %d items, want 3", len(g))
	}

	if _, err := spec.Done(root, "3", "abc123"); err == nil || !strings.Contains(err.Error(), "unchecked") {
		t.Fatalf("got %v, want an unchecked-boxes error", err)
	}
	for _, p := range []string{f, i, filepath.Join(u, "update.md")} {
		answer(t, p, true)
	}
	if _, err := spec.Done(root, "3", "abc123"); err != nil {
		t.Fatal(err)
	}

	items, err := spec.Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, it := range items {
		if !it.Archived || it.Status != "done" || it.Change != "abc123" {
			t.Errorf("%s: archived=%v status=%s change=%s", it.ID(), it.Archived, it.Status, it.Change)
		}
	}
	upd, _ := os.ReadFile(filepath.Join(root, "spec/archive/updates/003-v0-2/update.md"))
	if !strings.Contains(string(upd), "- [x] 001 feature: Login") {
		t.Fatalf("includes not ticked:\n%s", upd)
	}
	if problems, _ := spec.Check(root); len(problems) != 0 {
		t.Fatalf("check after done: %v", problems)
	}
}

func TestCheckFindsProblems(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "feature", "Login", spec.NewOptions{})
	data, _ := os.ReadFile(f)
	s := strings.Replace(string(data), "status: open", "status: wip", 1)
	s = strings.Replace(s, "(R1)", "", 1)
	if err := os.WriteFile(f, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
	problems, err := spec.Check(root)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(problems, "\n")
	for _, want := range []string{"line 1 must be", "R1 has no task"} {
		if !strings.Contains(joined, want) {
			t.Errorf("check missed %q in:\n%s", want, joined)
		}
	}
}
