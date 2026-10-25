---
name: git-commit
description: Commit changes in the dusupay-sdk-go repo the project's way. Runs the quality gate (gofmt, go vet, golangci-lint, go test), proposes a logical split when the working tree mixes unrelated changes, and writes Conventional Commits messages with SDK-specific scopes and BREAKING CHANGE footers. Use this whenever the user asks to commit, make a commit, "commit this", save or checkpoint work in git, write a commit message, or says things like "ok ship it to git", "commit the changes", or "split this into commits", even if they don't say "skill". Commits only; it never pushes.
---

# git-commit

Turn the current working-tree changes into clean, well-described commits for this Go SDK
(`github.com/kachit/dusupay-sdk-go/v2`, a wrapper around the DusuPay payments API).

The goal is a history that someone can read later to answer "what changed in the public API
and why", because SDK users upgrade by reading commits and release notes. That is why the
message format, the scopes and the breaking-change marking matter more here than in an app repo.

## Workflow

### 1. Look before touching anything

Run these together:

```bash
git status --short
git diff --stat
git diff --cached --stat
git branch --show-current
git log --oneline -5
```

- **Branch guard:** if the current branch is `master` (the main branch), stop and offer to create
  a branch first. Suggest a name from the change, e.g. `feat/nested-callbacks`. Committing straight
  to `master` bypasses PR review and CI on the PR. Any other branch, `v2` included, is fine.
- **Nothing to commit:** say so and stop. Don't create empty commits.
- **Already staged changes:** the user may have staged something on purpose. Treat the staged set
  as their intended commit unless they say otherwise, and mention any unstaged leftovers.

Read the actual diffs (`git diff`, `git diff --cached`) for the files involved. Writing the message
from the filenames alone produces vague messages. You need the content to explain *why*.

### 2. Guard against things that must not be committed

Before staging, scan the change list for:

- **Private keys and secrets:** `*.pem` files other than public keys in `stubs/rsa/`, `private*`,
  `.env`, or real-looking `public-key`/`secret-key`/signing-key values in code or fixtures. The
  repo deliberately gitignores `main/` (it holds `private.pem` for generating test signatures).
  If any such file shows up, stop and ask. Leaked merchant credentials can move real money.
- **Local or generated junk:** `coverage.out`, `coverage.html`, `*.test`, binaries,
  `.claude/settings.local.json`, workspace or scratch directories.

Stage files by explicit path (`git add path/a path/b`). Never use `git add -A` or `git add .`.
Explicit paths make it impossible to sweep in a stray file you didn't review.

### 3. Run the quality gate

```bash
bash .claude/skills/git-commit/scripts/precheck.sh
```

It mirrors CI (`.github/workflows/tests.yml`) by running gofmt, go vet, golangci-lint and go test.
It runs every step and reports all failures at once.

- **Passed:** continue.
- **Failed:** show the user the failing output and stop. Don't commit "to fix later", and don't
  skip checks with `--no-verify`. A red commit breaks the CI badge and anyone bisecting.
  - If the only failure is gofmt, the fix is mechanical. Offer to run `gofmt -w` on the listed
    files and re-run the gate.
  - Anything else (a failing test, a lint finding) needs a real decision, so ask.
- **Docs-only changes** (`*.md` and nothing else): the gate can be skipped. Say that you skipped it.
- **golangci-lint not installed:** the script reports SKIPPED. Mention that CI will still lint.

### 4. Decide: one commit or several

One commit should be one logical change: something you could revert or cherry-pick on its own.
If the diff mixes unrelated work (say, a new feature in `cards.go` plus an unrelated README typo
fix plus a CI tweak), propose a split before committing anything. Show it as a short plan:

```
Proposed commits:
1. feat(cards): add CardEncryptor for direct card payments
   cards.go, cards_test.go, env.go
2. docs: fix typo in README installation section
   README.md
3. ci: bump golangci-lint to v2.8.0
   .github/workflows/tests.yml
Commit as planned?
```

Wait for confirmation, then commit in that order. Keep each commit green: a test file belongs in
the same commit as the code it tests, and fixture changes under `stubs/` go with the tests that
use them. If a split would leave an intermediate commit that doesn't build, merge those groups.

Don't over-split. A feature with its tests, fixtures, README example and CLAUDE.md note is *one*
change, because all of those files describe the same thing. Only propose a split when the parts
have different reasons to exist. When everything is one change, commit without asking.

### 5. Write the message (Conventional Commits)

```
<type>(<scope>)<!>: <subject>

<body>

<footers>
```

**Type:** pick by what the change does for SDK users.

| type | use for |
|---|---|
| `feat` | new public API: types, fields, methods, constants |
| `fix` | wrong behaviour corrected (bad JSON tag, a value that could never be sent, a broken parse) |
| `refactor` | internal restructuring, no behaviour change |
| `test` | tests or `stubs/` fixtures only |
| `docs` | README, CLAUDE.md, comments only |
| `ci` | `.github/workflows` |
| `build` | `go.mod`/`go.sum`, dependency bumps |
| `chore` | anything else: tooling, `.gitignore`, `.claude/` |

**Scope:** the SDK area touched. Use these names so history stays greppable:
`collections`, `payouts`, `merchants`, `providers`, `banks`, `webhooks`, `signature`, `cards`,
`http`, `config`, `client`, `common`. Leave the scope out if the change spans several areas evenly.
Don't invent a vague scope like `sdk` or `misc`.

**Subject:** imperative mood ("add", not "added"), lowercase, no trailing period, 72 characters or
fewer for the whole first line. Say what changed for the user of the SDK, not which file was edited.

**Body:** explain *why*, wrapped at about 72 columns. For this repo the "why" is usually a DusuPay
API behaviour or a doc page. Cite it (for example `See https://developer.dusupay.com/...`) when it
motivated the change. Leave out the body only for trivial commits.

**Breaking changes:** this is a library, so flag any change that can break a caller's code:
renamed, removed or retyped exported identifiers, changed function signatures, changed JSON
shape, or changed validation that now rejects previously accepted input. Add `!` after the
type/scope, plus a `BREAKING CHANGE:` footer that tells users how to migrate. Check
`git diff` for changed `func (` signatures and exported struct fields. This is the case most
likely to be missed.

**Attribution footer:** if the session's system instructions specify a commit trailer (such as
`Co-Authored-By: ...`), end the message with it exactly as given. Don't invent one.

### 6. Commit

Pass the message through a heredoc so newlines and quotes survive. In Git Bash:

```bash
git commit -F - <<'EOF'
feat(webhooks)!: parse nested event/payload callbacks

DusuPay now sends callbacks as {"event": ..., "payload": {...}}.
The flat structs left every field empty, so signature validation
failed on real callbacks.
See https://developer.dusupay.com/funds-collection/getting-started.md

BREAKING CHANGE: CollectionCallbackPayload no longer carries Event.
Unmarshal into CollectionCallback and pass it to the validators.
EOF
```

If a pre-commit hook rejects the commit, fix the cause and make a **new** commit. Don't use
`--amend` (the failed commit never happened, so amending would rewrite the previous, unrelated
commit). Don't use `--no-verify`.

### 7. Report and stop

Show `git log --oneline -<n>` for the commits you just made and `git status --short`, so the
user sees anything left uncommitted. Then stop. **Don't push.** Pushing publishes the branch and
triggers CI, and this skill only commits. If the user wants to push, they will ask; mention
`git push` as the next step at most.

## Examples

**Single feature with tests and docs → one commit**
Changed: `payouts.go`, `payouts_test.go`, `README.md` (new `SwiftCode` field, its test, an example).
→ `feat(payouts): add swift and branch fields to extra_params`

**Bug fix that changes a type → breaking**
Changed: `collections.go` (`AllowFinalStatusChange bool` → `*bool`), plus its test.
→ `fix(collections)!: allow sending allow_final_status_change=false`
Body explains that `omitempty` on a bool dropped `false`, so the API default (`true`) always won.
Footer: `BREAKING CHANGE: AllowFinalStatusChange is now *bool; pass a pointer.`

**Mixed tree → propose a split**
Changed: `banks.go` (new filter field), `.github/workflows/tests.yml` (Go version bump),
`README.md` (typo unrelated to banks).
→ Three commits (`feat(banks)`, `ci`, `docs`), proposed first and committed after a "yes".
