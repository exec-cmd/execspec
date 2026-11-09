# Scout

The owner asks what to do next, what could go into the next update, or for options.

1. **Gather.** Read `SPEC.md` (Context, Principles, Ask first), every `personal/` file, `notes/decisions.md`, `xspec list` and `xspec questions`, `ROADMAP.md` History, and `TODO`/`FIXME` comments in the code. Run [bug-hunt.md](bug-hunt.md) only when the owner asks for bugs too.
2. **List candidates.** Each one: what it is, the value to the user, the cost (hours or days), the risk, and where it came from (an open item, a TODO, the owner's note). Drop anything against a principle, and say that you dropped it.
3. **Recommend.** Pick a set that makes one coherent update, and say why this set and not the next best.
4. **Write.** `xspec new analysis "Next update options"` with the candidates as a table, the recommendation, and the dropped items.
5. **On the owner's pick.** One `xspec new update "<name>"` with a part per picked candidate. Existing changes go under `--with`. Then the work cycle continues from Draft.

Done when the analysis is written and, after the owner picks, the update exists.
