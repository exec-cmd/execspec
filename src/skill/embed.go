// Package skill embeds the execspec agent skill and installs it.
package skill

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed SKILL.md references
var files embed.FS

// Targets are the skill directories, relative to a project root or home.
var Targets = []string{".claude/skills/execspec", ".agents/skills/execspec"}

// Install writes the skill into every target under base. Without force it
// refuses before writing anything if a target already exists.
func Install(base string, force bool) ([]string, error) {
	var dirs []string
	for _, t := range Targets {
		dir := filepath.Join(base, t)
		if _, err := os.Stat(dir); err == nil && !force {
			return nil, fmt.Errorf("%s already exists; use --force to replace it", dir)
		}
		dirs = append(dirs, dir)
	}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return nil, err
		}
		if err := os.CopyFS(dir, files); err != nil {
			return nil, err
		}
	}
	return dirs, nil
}
