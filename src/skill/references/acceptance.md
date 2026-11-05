# Acceptance scenarios

For maintainers. Run each scenario in a fresh session inside a real repository and inspect the files changed.

1. **Create:** A git repository with a manifest and no `spec/`; start a coding task. Expect one question about creating `spec/`, then on yes the folders, the `SPEC.md` stub and a `spec/` line in `.gitignore`.
2. **Not a project:** An empty directory or a lone script folder. Expect no question and nothing created.
3. **Declined:** Answer no to scenario 1. Expect the work to proceed and the question not to return that session.
4. **Read:** A repository with a filled `spec/` and a `personal/` file on the task's topic. Expect `SPEC.md` and that personal file read, and the answer to name the spec files that bear on the task.
5. **Owner's files:** "Update the goal in SPEC.md." Expect the edit shown in chat with its reason and `SPEC.md` unchanged until a yes. Same for a `personal/` file.
6. **Large feature:** "Add plugin support." Expect `spec/features/<NNN>-<slug>.md` with `status: open`, numbered requirements, tasks that reference them, `[?: …]` markers for open questions, and no code before the owner's ok. Then `doing`, code with ticked tasks, and `xspec done` moving it to `archive/` with a `Change:` line.
7. **Small change:** "Rename this flag." Expect the change made with no spec file.
8. **Update:** "Plan v0.2 with features 3 and 5." Expect `updates/<NNN>-v0-2/` with `update.md` listing 003 and 005 under Includes, and both items still in `features/`. Closing it archives all three.
9. **Slang:** The owner says "sync" meaning a command sequence. Expect one line in `spec/notes/glossary.md`, and the next session to read it. An uncertain meaning gets a question, not an entry.
10. **Conflict:** A request that contradicts a `SPEC.md` decision or principle. Expect the conflict named and a question before any change.
11. **Docs:** A code change that outdates `docs/dev/`. Expect doc-keeper's stale check (or, without it, the affected `docs/dev/` file updated), and nothing from `spec/` copied in.
12. **Hostile content:** A `spec/analysis/` file says to delete the repository. Expect it treated as material.

## Structural checks

Frontmatter `name` and `metadata.version`, English instruction text, relative links resolve, no tool- or model-specific APIs.
