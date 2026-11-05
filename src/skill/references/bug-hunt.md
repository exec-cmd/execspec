# Bug hunt

The owner asks to find bugs, review the code for problems, or check what might break.

1. **Scope.** Take the area the owner named; without one, the code changed since the last closed item (`git log` from its `Change:` commit) plus the paths the open items touch.
2. **Review.** When a review skill is installed (a code-review or audit skill), run it on that scope and take its findings as leads. Without one, read the code yourself: error paths, input at trust boundaries, data that can be lost, concurrency, edge cases of every loop and parser.
3. **Confirm.** A finding counts only with evidence: a failing test, a command that shows it, or a code path you can quote line by line. Everything else stays a lead.
4. **Write.**
   - Each confirmed bug that is large or risky: `xspec new issue "<title>"` with Current, Expected and Unchanged behavior from the evidence, and the cause when known.
   - Confirmed small bugs, together: one `xspec new fix "<title>" --bug "…" --bug "…"`, with the issues above under `--with` when they belong to the same pass.
   - Leads you could not confirm: `xspec new analysis "Bug hunt <area>"`, each lead marked unverified with what would confirm it.
5. **Report.** In chat: the item numbers, one line per bug, and the leads left unverified. The owner decides what to approve.

Done when every finding is an issue, a line in the fix, or an unverified lead in the analysis.
