---
name: code-linting
description: Run and fix golangci-lint findings in the dusupay-sdk-go repo the way CI checks them. It separates findings in changed files from pre-existing ones, fixes code properly instead of suppressing, asks before adding any //nolint, and can run the exact CI-pinned linter version. Use this whenever the user asks to lint, run golangci-lint, fix lint/vet/gofmt/wsl/funlen/lll/gosec/bodyclose/errcheck warnings, clean up code style, or when CI's Lint step fails or "lint passes locally but fails in CI", even if they just say "make the linter happy" or paste a linter error.
---

# code-linting

The linter is part of the CI gate (`.github/workflows/tests.yml`, step "Lint"), so a finding
here means a red build. Config lives in `.golangci.yml` (golangci-lint v2 format: ~20 linters
including wsl_v5, nlreturn, funlen, lll, gosec and bodyclose, with relaxed rules for `_test.go`).

The aim is to make the code actually better, not to make the linter quiet. A `//nolint` hides
the finding for everyone afterwards, so it needs a real reason and the user's agreement.

## 1. Run it

```bash
bash .claude/skills/code-linting/scripts/lint.sh
```

The script runs `golangci-lint run ./...` (like CI) and groups the output into **findings in
changed files** (vs HEAD, plus untracked `.go` files) and **pre-existing findings**.

**Version mismatch.** CI pins a version in `tests.yml` (currently v2.8.0). If the local binary
differs, the script prints a WARNING. Different versions can disagree, because newer releases add
checks and change defaults. Run the exact CI version when:
- the user says CI fails but local passes, or the other way round
- you're about to report "lint is clean" as the final word before a commit or PR
- a finding looks like it might be version-specific

```bash
bash .claude/skills/code-linting/scripts/lint.sh --ci
```

This uses `go run ...golangci-lint@<pinned>`. The first run downloads it (about a minute), and
later runs are cached.

If the user pasted a CI log, start from that log. It's the ground truth for what CI saw.

**Windows line endings.** If gofmt flags many files at `1:1` ("File is not properly formatted")
that nobody touched, the cause is almost certainly CRLF line endings from Git's `core.autocrlf`.
CI runs on Linux with LF and won't see these findings. Don't reformat or re-save those files,
because that would produce a huge noise diff. Tell the user. The real fix is a `.gitattributes`
with `*.go text eol=lf`, which is a repo-config change, so propose it rather than adding it.

## 2. Triage

- **Findings in changed files:** in scope, fix them. This is work the user is doing right now.
- **Pre-existing findings:** report them separately with a count and a one-line summary, and
  offer to fix them. Don't silently expand the diff into unrelated files, because that makes the
  user's change harder to review.
  - Exception: if CI is currently red because of them, they block the build, so say so and fix
    them with the user's OK.

For each finding, look up the linter in `references/linters.md`. It gives the preferred fix for
this repo and says when suppression is ever acceptable.

## 3. Fix properly

- Make the smallest change that resolves the finding and keeps the code idiomatic for this repo.
  Match the surrounding style: early-return `switch` validation, `isValid()` on request structs,
  `fmt.Errorf("Type.Method error: %v", err)` wrapping, and multi-line signatures for long
  parameter lists.
- `gofmt` findings are mechanical: run `gofmt -w <file>`.
- `golangci-lint run --fix` can auto-fix some linters (gofmt, wsl, misspell...). It's fine to use
  on changed files. Afterwards, review the diff (`git diff`), because auto-fixes can touch
  unrelated lines.
- **Don't change exported API to satisfy a linter** (renaming identifiers, turning value params
  into pointers for gocritic `hugeParam`, and so on) without telling the user. In an SDK, that's
  a breaking change for every caller, which is a far bigger cost than a style warning.
- **Behaviour must not change.** If a "lint fix" alters logic (say, restructuring an `if` chain
  that guards validation), run the tests: `go test ./...`.

## 4. Suppression: ask first

Suggest a `//nolint` only when the finding is a genuine false positive or the fix would make the
code worse. Before adding one, show the user the finding, why you think it's a false positive,
and the exact comment you'd add. Then wait for a yes.

The repo's format names the linter and gives a reason after `//`:

```go
return tr.http.Do(req) //nolint:gosec // URL is built from validated Config, not user input
```

Never use a bare `//nolint` or `//nolint:all`. Those silence linters nobody evaluated.

The one pre-approved case is validation-error tests where the method returns a nil response:
`//nolint:bodyclose // returns nil response on validation error`. It's an established pattern
in every resource test file, so use it without asking.

## 5. The config is off-limits unless asked

Don't edit `.golangci.yml` to make findings go away (disabling a linter, raising funlen limits,
adding exclusions). If a rule seems wrong for this codebase, say so and propose the config
change. Only make it if the user asks. The config is a team decision, and CI uses it too.

## 6. Verify and report

Re-run the script, using `--ci` for the final check if the local version differs. Then run
`go test ./...` if you changed any non-trivial code. Report:

- what you fixed, grouped by linter (for example "wsl_v5 ×4 in webhooks.go: blank lines before
  returns")
- any suppression you added and its reason (all of them were approved)
- pre-existing findings you left alone, with a count
- which linter version gave the final result, local or CI-pinned

Don't commit. The `git-commit` skill handles that, and it runs the linter again as part of its
gate.
