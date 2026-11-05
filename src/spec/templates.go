package spec

// Placeholders: {{num}}, {{title}}, {{includes}}, {{bugs}}.
// Every [?: …] marker blocks `xspec start` until someone answers it.

const specStub = `# %s

> The owner fills this file. Agents read it and propose edits in chat.

## Goal

## Scope

## Principles

One line each. Plans check themselves against these in their Gates section.

## Decisions
`

const glossaryStub = `# Glossary

Project slang, one line per term: **term** - meaning (where it came up)
`

var templates = map[string]string{
	"feature": `status: open
# {{num}} {{title}}

## Why
[?: the problem, and who gets what when it is solved]

## Requirements
- R1: WHEN [?: condition] THE SYSTEM SHALL [?: behavior]

## Design
[?: the approach and the code it touches]

## Tasks
- [ ] T1 [?: task] (R1)
`,

	"issue": `status: open
# {{num}} {{title}}

## Current behavior
- WHEN [?: condition] THEN the system [?: wrong behavior]

## Expected behavior
- WHEN [?: condition] THEN the system SHALL [?: correct behavior]

## Unchanged behavior
- WHEN [?: condition] THEN the system SHALL CONTINUE TO [?: behavior that must not break]

## Cause
[?: root cause, once found]

## Tasks
- [ ] T1 regression test reproduces the bug
- [ ] T2 [?: fix]
`,

	"update": `status: open
# {{num}} {{title}}

## Goal
[?: why this batch, and what the user gets]

## Includes
{{includes}}
## Also changes
Migrations, config, dependencies outside the included items.

## Tasks
- [ ] T1 [?: task]

## Risks
What may break, and how to roll back.
`,

	"fix": `status: open
# {{num}} {{title}}

## Includes
{{includes}}
## Small bugs
{{bugs}}
## Tasks
- [ ] T1 a regression test for every bug above
`,

	"plan": `status: open
# {{num}} {{title}}

## Context
[?: what prompted this, and the constraints]

## Gates
One line per principle in SPEC.md: kept, or broken and why.

## Decision
[?: the chosen approach]

## Steps
- [ ] [?: step; name the items it creates, e.g. 007]
`,

	"analysis": `# {{title}}

## Question

## Sources

## Findings

## Verdict
`,

	"note": `# {{title}}
`,

	// design.md beside update.md, and fix.md with --dir.
	"design": `# Design: {{title}}
`,
}
