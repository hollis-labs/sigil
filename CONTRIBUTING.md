# Contributing

How a change gets from your clone into `main`. This is deliberately short: most of what you need is already written somewhere closer to the thing it describes, and this page points at it rather than keeping a second copy that drifts.

## Before your first change

`README.md` has the install and build steps. Run `lefthook install` — a tracked config installs no hooks by itself, and a clone that skips it has no commit or push checks and says nothing about it.

`AGENTS.md` is the fastest orientation to the layout: which package does what, and the boundaries that are not obvious from reading the code. It is written for an agent working in the repo, which makes it unusually direct about the things that bite. The specs live in `docs/01_architecture.md` through `docs/09_sysop-ui-kit.md`.

## The sequence

1. **Branch.** `<type>/<short-slug>`, where the type matches the change — `feat`, `fix`, `docs`, `chore`. Nothing enforces this; it is what the history does.
2. **Change one thing.** A branch carrying two unrelated changes costs the reviewer the ability to accept one and question the other.
3. **Run the checks** when the change is done: `make vet test`, and `cd sysop && go test ./...` if you touched `sysop/`.
4. **Push, and open a pull request** against `main`.

Commit subjects follow the conventional-commit shape — a type, an optional scope, a colon, then the summary. To see what is actually in use rather than trusting this sentence:

```
git log --no-merges -40 --format='%s' | grep -oE '^[a-z]+(\([a-z-]+\))?:' | sort -u
```

## What a pull request should carry

The reviewer was not there when you made the decisions. State what the change does, what it deliberately leaves alone, and the evidence that it works — the commands you ran and what came back, not a claim that it passes. For a renderer change, name the target and show the generated output you inspected.

If a number appears in the description, put the command that produced it beside it.

## The one that cannot be undone

**`sigil generate --clean` deletes files.** It removes everything in the output directory before generating. Point `--output` at a scratch directory while developing, never at one holding hand-written code, and check `--dry-run` first.

## Things that surprise people

- **`sysop/` is its own Go module.** `go test ./...` at the root does not reach it.
- **A component renders in every target or its unsupported status is documented.** A schema added under `internal/components/builtin/` without a case in each renderer validates and then emits nothing.
- **Only `internal/renderer/engine.go` writes to disk.** Renderers return `[]OutputFile`; that is what makes `--dry-run` and `--clean` work.
- **The Sysop UI screens are hand-composed on purpose**, and the demo hooks are hand-written mocks. Regenerating over them breaks them; the `make` demo targets protect the hooks with `rsync --ignore-existing`.
- **`make demo`, `make se-demo` and `make cw-demo` regenerate checked-in apps.** Use them to see a renderer change, not to verify one.
- **Generated code must never depend on Sigil at runtime.**

## What this does not cover

- **Which change is worth making.** There is no roadmap commitment here by design; that conversation happens in issues.
- **Release.** Sigil has no tagged releases yet.
