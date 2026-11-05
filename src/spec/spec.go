// Package spec reads and changes the spec/ directory: layout, numbered items,
// statuses, the archive and check rules. It knows nothing about the CLI.
package spec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	Dir     = "spec"
	Archive = "archive"
)

// Type is a kind of spec file. Status types are numbered and go through
// open → doing → done; the others are free files named by slug.
type Type struct {
	Name   string
	Folder string
	Status bool
}

var Types = []Type{
	{"feature", "features", true},
	{"issue", "issues", true},
	{"update", "updates", true},
	{"fix", "fixes", true},
	{"plan", "plans", true},
	{"note", "notes", false},
	{"analysis", "analysis", false},
}

// Folders created by Init besides the type folders.
var extraFolders = []string{"personal", Archive}

func TypeByName(name string) (Type, error) {
	for _, t := range Types {
		if t.Name == name {
			return t, nil
		}
	}
	names := make([]string, len(Types))
	for i, t := range Types {
		names[i] = t.Name
	}
	return Type{}, fmt.Errorf("unknown type %q, available: %s", name, strings.Join(names, ", "))
}

// ProjectRoot is the git top-level of dir, or dir itself outside git.
func ProjectRoot(dir string) string {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return dir
	}
	return strings.TrimSpace(string(out))
}

// FindRoot walks up from dir to the nearest directory holding spec/.
func FindRoot(dir string) (string, error) {
	for d := dir; ; d = filepath.Dir(d) {
		if fi, err := os.Stat(filepath.Join(d, Dir)); err == nil && fi.IsDir() {
			return d, nil
		}
		if filepath.Dir(d) == d {
			return "", errors.New("no spec/ directory found; run `xspec init` first")
		}
	}
}

// Init creates the layout and the SPEC.md stub under root and makes git ignore
// spec/. Existing files are kept, so it is safe to run twice.
func Init(root string) ([]string, error) {
	var created []string
	folders := extraFolders
	for _, t := range Types {
		folders = append(folders, t.Folder)
	}
	for _, f := range folders {
		p := filepath.Join(root, Dir, f)
		if _, err := os.Stat(p); err == nil {
			continue
		}
		if err := os.MkdirAll(p, 0o755); err != nil {
			return created, err
		}
		created = append(created, p)
	}

	stubs := map[string]string{
		"SPEC.md":           fmt.Sprintf(specStub, filepath.Base(root)),
		"notes/glossary.md": glossaryStub,
	}
	for name, body := range stubs {
		p := filepath.Join(root, Dir, name)
		if _, err := os.Stat(p); err == nil {
			continue
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			return created, err
		}
		created = append(created, p)
	}

	added, err := ensureIgnored(root)
	if err != nil {
		return created, err
	}
	if added {
		created = append(created, filepath.Join(root, ".gitignore")+" (+spec/)")
	}
	return created, nil
}

// ignored reports whether git ignores spec/; ok is false outside git.
func ignored(root string) (ignored, ok bool) {
	err := exec.Command("git", "-C", root, "check-ignore", "-q", Dir+"/").Run()
	var exitErr *exec.ExitError
	switch {
	case err == nil:
		return true, true
	case errors.As(err, &exitErr) && exitErr.ExitCode() == 1:
		return false, true
	}
	return false, false
}

func ensureIgnored(root string) (bool, error) {
	if ign, _ := ignored(root); ign {
		return false, nil
	}
	path := filepath.Join(root, ".gitignore")
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		switch strings.TrimSpace(line) {
		case "spec", "spec/", "/spec", "/spec/":
			return false, nil
		}
	}
	if len(data) > 0 && !strings.HasSuffix(string(data), "\n") {
		data = append(data, '\n')
	}
	// Anchored, so nested folders such as src/spec stay tracked.
	data = append(data, "/"+Dir+"/\n"...)
	return true, os.WriteFile(path, data, 0o644)
}
