# Improve and digest

The owner asks to improve a spec file, or to summarize one; the work cycle also shows a digest after every draft. Every spec file counts: items, `SPEC.md` and `ROADMAP.md`.

## Scope

| File | Improve | Digest |
| --- | --- | --- |
| change | yes | yes |
| update | yes, with the changes it includes | yes, one digest for the group |
| `SPEC.md` | only on the owner's explicit yes naming the file | yes |
| `ROADMAP.md` | never: generated | yes, from `xspec roadmap` |

## Improve passes

Run in order, on the file itself:

1. **English.** Rewrite any other language into English. Chat answers stay in the user's language.
2. **Wording.** Fix grammar and terms. Use the words from `notes/glossary.md`.
3. **Bad requirements.** Remove a requirement that is vague, untestable, duplicated, against `SPEC.md`, or a solution stated as a need, with its tasks. Or flag it as `[?: why]`.
4. **Shorten.** Cut padding. Keep every fact.
5. **Expand.** Add an acceptance criterion, edge case or scope limit only when the code or the owner's files support it. Otherwise add a marker.

An update: run the passes on the update and on each included change.

## Hard rules

- Never delete a `[?: …]` marker without the owner's answer.
- Never invent a requirement.
- Keep `Verify:` lines and the numbers of the requirements that stay.
- Run `xspec check` after the edit and fix what it reports, except open questions.
- Never write the digest to a file; it would go stale beside the item.

## Digest

Chat only, in the user's language, at most 10 lines:

- Goal, scope in and out.
- Tasks, or key requirements by number. Name each dropped or added one.
- Decisions taken and open markers.
- Risks.
- What `done` will verify: the Verify lines.

For `ROADMAP.md`, give Now, Next and History. For `SPEC.md`, give Context, Audience and Principles.

After an improve, end with one question: keep, change, or ask?
