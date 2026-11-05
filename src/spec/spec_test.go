package spec_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/exec-cmd/execspec/src/spec"
)

// setup makes a spec/ whose kickoff (001) is already closed, so tests start
// numbering at 002 with an empty board.
func setup(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if _, err := spec.Init(root); err != nil {
		t.Fatal(err)
	}
	kickoff := filepath.Join(root, "spec/plans/001-kickoff.md")
	answer(t, kickoff, true)
	run(t, root, "1")
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

// run approves, starts and closes ref.
func run(t *testing.T, root, ref string) {
	t.Helper()
	if _, err := spec.Approve(root, ref); err != nil {
		t.Fatal(err)
	}
	if _, err := spec.Start(root, ref); err != nil {
		t.Fatal(err)
	}
	if _, err := spec.Done(root, ref, noVerify()); err != nil {
		t.Fatal(err)
	}
}

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
	for _, f := range []string{"SPEC.md", "ROADMAP.md", "notes/glossary.md", "archive", "personal", "plans/001-kickoff.md"} {
		if _, err := os.Stat(filepath.Join(root, "spec", f)); err != nil {
			t.Error(err)
		}
	}
	if items := scan(t, root); len(items) != 1 || items[0].Slug != "kickoff" {
		t.Fatalf("want only the kickoff, got %d items", len(items))
	}
}

func TestNumbersCountTheArchive(t *testing.T) {
	root := setup(t)
	if p := newItem(t, root, "issue", "Two", spec.NewOptions{}); filepath.Base(p) != "002-two.md" {
		t.Fatalf("got %s, want 002-two.md", p)
	}
	if n := newItem(t, root, "note", "Free note", spec.NewOptions{}); filepath.Base(n) != "free-note.md" {
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
		"fix in update":    {"update", spec.NewOptions{With: []int{2}}},
		"missing item":     {"update", spec.NewOptions{With: []int{9}}},
		"with on plan":     {"plan", spec.NewOptions{With: []int{2}}},
		"bug on update":    {"update", spec.NewOptions{Bugs: []string{"x"}}},
		"light update":     {"update", spec.NewOptions{Depth: "light"}},
		"unknown depth":    {"feature", spec.NewOptions{Depth: "deep"}},
		"archived kickoff": {"update", spec.NewOptions{With: []int{1}}},
	} {
		if _, err := spec.New(root, mustType(t, tc.typ), name, tc.opt); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}

func TestLightDepth(t *testing.T) {
	root := setup(t)
	p := newItem(t, root, "feature", "Dark mode", spec.NewOptions{Depth: "light"})
	data, _ := os.ReadFile(p)
	if strings.Contains(string(data), "## Requirements") || !strings.Contains(string(data), "## Goal") {
		t.Fatalf("not the light template:\n%s", data)
	}
	if _, err := spec.Start(root, "2"); err == nil {
		t.Fatal("light item started without approval")
	}
}

func TestMarkersInCodeAreNotQuestions(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "feature", "Login", spec.NewOptions{})
	body := "status: open\n# 002 Login\n\nuse `[?: …]` markers\n```\n[?: example]\n```\n[?: real]\n"
	if err := os.WriteFile(f, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	items := scan(t, root)
	if m := items[1].Markers; len(m) != 1 || m[0].Line != 8 || m[0].Text != "real" {
		t.Fatalf("got markers %+v, want one on line 8", m)
	}
}

func TestApprovalGates(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "feature", "Login", spec.NewOptions{})

	if _, err := spec.Approve(root, "2"); err == nil || !strings.Contains(err.Error(), "markers") {
		t.Fatalf("approve with markers: got %v", err)
	}
	answer(t, f, false)
	if _, err := spec.Start(root, "2"); err == nil || !strings.Contains(err.Error(), "approve") {
		t.Fatalf("start before approve: got %v", err)
	}
	if _, err := spec.Approve(root, "2"); err != nil {
		t.Fatal(err)
	}

	// ticking a box keeps the approval, editing the text drops it
	answer(t, f, true)
	if items := scan(t, root); items[1].Changed() {
		t.Fatal("ticking a task box invalidated the approval")
	}
	data, _ := os.ReadFile(f)
	if err := os.WriteFile(f, append(data, "one more line\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := spec.Start(root, "2"); err == nil || !strings.Contains(err.Error(), "changed after approval") {
		t.Fatalf("start after edit: got %v", err)
	}
	if problems, _ := spec.Check(root); !strings.Contains(strings.Join(problems, "\n"), "changed after approval") {
		t.Fatalf("check missed the edit: %v", problems)
	}
	if _, err := spec.Approve(root, "2"); err != nil {
		t.Fatal(err)
	}
	if _, err := spec.Start(root, "2"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateCascade(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "feature", "Login", spec.NewOptions{})
	i := newItem(t, root, "issue", "Crash", spec.NewOptions{})
	u := newItem(t, root, "update", "v0.2", spec.NewOptions{With: []int{2, 3}})
	for _, p := range []string{f, i, filepath.Join(u, "update.md")} {
		answer(t, p, false)
	}

	if g, err := spec.Approve(root, "4"); err != nil || len(g) != 3 {
		t.Fatalf("approve: %d items, %v", len(g), err)
	}
	if g, err := spec.Start(root, "4"); err != nil || len(g) != 3 {
		t.Fatalf("start: %d items, %v", len(g), err)
	}
	if _, err := spec.Done(root, "4", noVerify()); err == nil || !strings.Contains(err.Error(), "unchecked") {
		t.Fatalf("got %v, want an unchecked-boxes error", err)
	}
	for _, p := range []string{f, i, filepath.Join(u, "update.md")} {
		answer(t, p, true)
	}
	if _, err := spec.Done(root, "4", noVerify()); err != nil {
		t.Fatal(err)
	}

	for _, it := range scan(t, root) {
		if !it.Archived || it.Status != "done" {
			t.Errorf("%s: archived=%v status=%s", it.ID(), it.Archived, it.Status)
		}
		if it.Num > 1 && (it.Change != "abc123" || it.Date == "") {
			t.Errorf("%s: change=%q date=%q", it.ID(), it.Change, it.Date)
		}
	}
	upd, _ := os.ReadFile(filepath.Join(root, "spec/archive/updates/004-v0-2/update.md"))
	if !strings.Contains(string(upd), "- [x] 002 feature: Login") {
		t.Fatalf("includes not ticked:\n%s", upd)
	}
	if problems, _ := spec.Check(root); len(problems) != 0 {
		t.Fatalf("check after done: %v", problems)
	}
	roadmap, _ := os.ReadFile(filepath.Join(root, "spec/ROADMAP.md"))
	if !strings.Contains(string(roadmap), "004 update v0.2 (abc123)") || !strings.Contains(string(roadmap), "  - 002 feature Login") {
		t.Fatalf("roadmap history:\n%s", roadmap)
	}
}

func TestVerifyRunsBeforeDone(t *testing.T) {
	root := setup(t)
	f := newItem(t, root, "feature", "Login", spec.NewOptions{})
	answer(t, f, true) // the Verify marker becomes "Verify: true"
	if _, err := spec.Approve(root, "2"); err != nil {
		t.Fatal(err)
	}
	var ran []string
	fail := spec.DoneOptions{Run: func(c string) error { ran = append(ran, c); return errors.New("exit 1") }}
	if _, err := spec.Done(root, "2", fail); err == nil || !strings.Contains(err.Error(), "Verify") {
		t.Fatalf("got %v, want a Verify failure", err)
	}
	if len(ran) != 1 || ran[0] != "true" {
		t.Fatalf("ran %v, want [true]", ran)
	}
	if items := scan(t, root); items[1].Archived {
		t.Fatal("closed despite the failing Verify")
	}
	ok := spec.DoneOptions{Run: func(string) error { return nil }}
	if _, err := spec.Done(root, "2", ok); err != nil {
		t.Fatal(err)
	}
}

func TestNextAndQuestions(t *testing.T) {
	root := setup(t)
	if got := spec.Next(scan(t, root)); got != "nothing to do" {
		t.Fatalf("empty board: %q", got)
	}
	f := newItem(t, root, "feature", "Login", spec.NewOptions{})
	if got := spec.Next(scan(t, root)); !strings.HasPrefix(got, "answer 6 questions in 002") {
		t.Fatalf("open with markers: %q", got)
	}
	q := spec.Questions(root, scan(t, root))
	if !strings.Contains(q, "spec/features/002-login.md:5  the problem") {
		t.Fatalf("questions:\n%s", q)
	}
	answer(t, f, false)
	if got := spec.Next(scan(t, root)); !strings.HasPrefix(got, "the owner reviews and approves 002") {
		t.Fatalf("open: %q", got)
	}
	spec.Approve(root, "2")
	if got := spec.Next(scan(t, root)); !strings.HasPrefix(got, "start 002") {
		t.Fatalf("approved: %q", got)
	}
	spec.Start(root, "2")
	if got := spec.Next(scan(t, root)); !strings.HasPrefix(got, "work on 002 Login: T1") {
		t.Fatalf("doing: %q", got)
	}
}

func TestContext(t *testing.T) {
	root := setup(t)
	newItem(t, root, "feature", "Login", spec.NewOptions{})
	specMD := filepath.Join(root, "spec/SPEC.md")
	if err := os.WriteFile(specMD, []byte("# p\n\n## Principles\nNo network.\n\n## Decisions\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := spec.Context(root, "2")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"=== spec/features/002-login.md ===", "=== spec/SPEC.md § Principles ===\nNo network.", "§ Stack ===", "=== spec/notes/glossary.md ==="} {
		if !strings.Contains(out, want) {
			t.Errorf("context missed %q in:\n%s", want, out)
		}
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
