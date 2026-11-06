# Improve and digest

The owner asks to improve a spec file, or to summarize one. Every spec file counts: items, `SPEC.md` and `ROADMAP.md`.

## Scope

| File | Improve | Digest |
| --- | --- | --- |
| feature, issue, plan | yes; an approved item needs a yes first | yes |
| update, fix | yes, with the items they include | yes, one digest for the group |
| `SPEC.md` | only on the owner's explicit yes naming the file | yes |
| `ROADMAP.md` | never: generated | yes, from `xspec roadmap` |

An approved item: ask first, because an edit voids its approval. On yes, edit it and tell the owner to run `xspec approve <num>` again; until then `xspec check` reports the edit and `xspec start` refuses.

## Improve passes

Run in order, on the file itself:

1. **English.** Rewrite any other language into English. Chat answers stay in the user's language.
2. **Wording.** Fix grammar and terms. Use the words from `notes/glossary.md`.
3. **Bad requirements.** Remove a requirement that is vague, untestable, duplicated, against `SPEC.md`, or a solution stated as a need, with its tasks. Or flag it as `[?: why]`.
4. **Shorten.** Cut padding. Keep every fact.
5. **Expand.** Add an acceptance criterion, edge case or scope limit only when the code or the owner's files support it. Otherwise add a marker.

Update and fix: run the passes on the parent and on each included item.

## Hard rules

- Never delete a `[?: …]` marker without the owner's answer.
- Never invent a requirement.
- Keep `Verify:` lines and the numbers of the requirements that stay.
- Run `xspec check` after the edit. Fix what it reports, except open questions and the edit-after-approval report, which only the owner's re-approval clears.
- Never write the digest to a file. A stored summary changes the item hash and voids approval.
- Never run `xspec approve`.

## Digest

Chat only, in the user's language, at most 10 lines:

- Goal, scope in and out.
- Key requirements, by number. Name each dropped or added requirement.
- Decisions taken and open markers.
- Risks.
- What `done` will verify: the Verify lines.

For `ROADMAP.md`, give Now, Next and History. For `SPEC.md`, give Goal, Scope and Principles.

End with one question: approve, change, or ask?

Show the digest after every draft or improve, before asking for approval, and before `xspec start`.
