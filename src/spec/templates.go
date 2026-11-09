package spec

// Placeholders: {{num}}, {{title}}, {{includes}}.
// Every [?: …] marker stops `xspec done` until someone answers it.

const specStub = `# %s

<!-- Project brief. Fill every section with a few lines; the agent drafts it from the repository and asks you a few questions. -->

## Context
<!-- What the project is and why it exists. -->

## Audience
<!-- Who uses it. -->

## Principles
<!-- One line each: the rules every change follows. -->

## Design
<!-- How it is built in a few lines, or a link to the design doc. -->

## Stack
<!-- Part | Choice | Owner knows it (yes or no); the agent explains the parts you do not know. -->

## Ask first
The agent stops and asks before it does any of these:
- add a dependency
- change a file format or a public interface
- migrate or delete data
- touch secrets, CI or release configuration
`

const glossaryStub = `# Glossary

Project slang, one line per term: **term** - meaning (where it came up)
`

var templates = map[string]string{
	"change": `status: open
# {{num}} {{title}}

## Goal
[?: what changes, and for whom]

## Tasks
- [ ] T1 [?: task]

Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
`,

	"change/bug": `status: open
# {{num}} {{title}}

## Current behavior
- WHEN [?: condition] THEN the system [?: wrong behavior]

## Expected behavior
- WHEN [?: condition] THEN the system SHALL [?: correct behavior]

## Unchanged behavior
- WHEN [?: condition] THEN the system SHALL CONTINUE TO [?: behavior that must not break]

## Tasks
- [ ] T1 regression test reproduces the bug
- [ ] T2 [?: fix]

Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
`,

	"update": `status: open
# {{num}} {{title}}

## Goal
[?: why this batch, and what the user gets]

## Includes
{{includes}}
## Parts

### P1 [?: part name]
[?: why this part]

## Tasks
- [ ] T1 [?: task]

Verify: [?: a command that must pass before done, e.g. go test ./...; or delete this line]
`,

	"analysis": `# {{title}}

## Question

## Sources

## Findings

## Verdict
`,

	"note": `# {{title}}
`,

	// design.md beside update.md, with --dir.
	"design": `# Design: {{title}}
`,
}
