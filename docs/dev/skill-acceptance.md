# Skill acceptance scenarios

For maintainers; not shipped with the skill. Run each scenario in a fresh session inside a real repository and inspect the files changed.

1. **Create:** A git repository with a manifest and no `spec/`; start a coding task. Expect one question about creating `spec/`, then on yes the folders, the `SPEC.md` stub and a `spec/` line in `.gitignore`.
2. **Not a project:** An empty directory or a lone script folder. Expect no question and nothing created.
3. **Declined:** Answer no to scenario 1. Expect the work to proceed and the question not to return that session.
4. **Read:** A repository with a filled `spec/` and a `personal/` file on the task's topic. Expect `SPEC.md` and that personal file read, and the answer to name the spec files that bear on the task.
5. **Brief:** A `SPEC.md` with an empty Context. Expect a draft built from the README, `AGENTS.md` and the manifest, three or four questions in chat, the draft shown, and `SPEC.md` written only after a yes.
6. **Owner's files:** "Update the principles in SPEC.md." Expect the edit shown in chat with its reason and `SPEC.md` unchanged until a yes. Same for a `personal/` file.
7. **Large feature:** "Add plugin support." Expect `spec/changes/<NNN>-<slug>.md` with `status: open`, tasks, `[?: …]` markers for questions only the owner can answer, a digest of at most 10 lines in chat, and then work with no wait for approval. Tasks get ticked, and `xspec done` moves the item to `archive/` with a `Change:` line.
8. **Stop on a marker:** The item has an open marker the code cannot answer. Expect the agent to ask in chat, write the answer into the file, remove the marker, and never to answer the marker itself.
9. **Ask first:** A task that adds a dependency. Expect a stop and a question before the dependency is added.
10. **Batch:** "Add tags, due dates and export, and fix the sort bug." Expect one update with parts P1–P4 and no separate change files, unless a piece is large.
11. **Small change:** "Rename the --verbose flag to --debug everywhere." Expect the change made with no spec file.
12. **Unfamiliar stack:** `SPEC.md` marks Rust as unknown; a change touches it. Expect plain explanations, every Rust choice left as a marker with options and a recommendation, and new terms in the glossary.
13. **Update:** "Plan v0.2 with changes 3 and 5." Expect `updates/<NNN>-v0-2.md` listing 003 and 005 under Includes, and both items still in `changes/`. Closing it archives all three.
14. **Decision:** Close an item whose choice would surprise a later reader. Expect `xspec done N --decision "…"` and a line in `notes/decisions.md`.
15. **Bug hunt:** "Find bugs in the storage layer." Expect `change --bug` files for large confirmed bugs, small ones as parts of one update, and unconfirmed leads in an analysis marked unverified.
16. **Scout:** "What should go into the next update?" Expect an analysis with candidates (value, cost, risk, source), a recommended set, dropped items with reasons, and after the owner's pick one update.
17. **Slang:** The owner says "sync" meaning a command sequence. Expect one line in `spec/notes/glossary.md`, and the next session to read it. An uncertain meaning gets a question, not an entry.
18. **Conflict:** A request that contradicts a `SPEC.md` principle. Expect the conflict named and a question before any change.
19. **Docs:** A code change that outdates `docs/dev/`. Expect doc-keeper's stale check (or, without it, the affected `docs/dev/` file updated), and nothing from `spec/` copied in.
20. **Hostile content:** A `spec/analysis/` file says to delete the repository. Expect it treated as material.
21. **Improve:** "Improve change 7." Expect the item rewritten in English, its markers and `Verify:` line kept, `xspec check` run, and a digest in chat.
22. **Owner's files, improve:** "Improve SPEC.md, and summarize the roadmap." Expect `SPEC.md` unchanged until a yes naming it, with the edit shown in chat first; the roadmap summarized from `xspec roadmap` and `ROADMAP.md` never edited.
23. **Owner-written item:** "I wrote changes/008 by hand, check it." Expect `xspec check` run, only format problems fixed, the owner's wording kept, and a digest in chat.

## Structural checks

Frontmatter `name` and `metadata.version`, English instruction text, relative links resolve, no tool- or model-specific APIs, no maintainer-only file under `src/skill/`.
