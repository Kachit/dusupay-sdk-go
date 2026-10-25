# Linter fix guide for dusupay-sdk-go

The linters enabled in `.golangci.yml`, what each one usually flags in this codebase, and the
preferred fix. "Suppress?" says when a `//nolint` is ever reasonable. Even then, ask the user
first (see SKILL.md).

In `_test.go` files, **funlen, goconst, lll, staticcheck and wsl_v5 are excluded** by config.
Everything else still applies to tests.

## Contents
- Correctness: bodyclose, errcheck, govet, staticcheck, ineffassign, unused, gosec
- Complexity: funlen, gocognit, gocyclo, nakedret
- Style: wsl_v5, nlreturn, lll, gocritic, goconst, prealloc, unconvert, misspell, dogsled
- Formatter: gofmt

---

## Correctness

### bodyclose
Every `*http.Response` returned to a caller must get its body closed.
- In SDK code, resource methods hand `rsp` back to the caller, and `unmarshalResponse` closes and
  re-wraps the body. Keep using `unmarshalResponse`. Don't read bodies by hand.
- In tests, use `defer func() { _ = resp.Body.Close() }()`. Don't use a bare
  `defer resp.Body.Close()`, because errcheck flags it.
- Suppress? Yes, but only on validation-error tests where the method returns a nil response,
  using the existing form:
  `//nolint:bodyclose // returns nil response on validation error`

### errcheck
Unchecked errors. Handle the error or explicitly discard it with `_ =` when ignoring it is
correct (closing a read body, writing to a hash). Wrap returned errors the repo way:
`fmt.Errorf("ResourceName.Method error: %v", err)`.
Suppress? No. `_ =` already documents the intent.

### govet / staticcheck / ineffassign / unused
These are almost always real bugs or dead code. Fix them. For `unused` on something exported
later in a series of changes, check whether a test or the next step will use it before deleting.
Suppress? No.

### gosec
Security findings. Common in this repo:
- `G107`/`G704`-style "HTTP request with variable URL" on `tr.http.Do(req)`. This one is already
  suppressed with a reason in `http.go`. Leave it alone.
- Weak crypto, or `math/rand` for crypto: switch to `crypto/rand`. Don't suppress.
- `G115` (integer overflow conversion) is excluded in config.
Suppress? Only for a false positive with a concrete reason, as in `http.go`. Ask first.

## Complexity

### funlen (70 lines / 50 statements, non-test code only)
Extract a well-named helper. The repo's pattern is a request struct with an `isValid()` method,
so pull long validation chains into `isValid()` rather than inlining them in the resource method.
Suppress? No.

### gocognit / gocyclo
Usually long `if` chains. Prefer an early-return `switch` (as in `Config.IsValid`, `CardData.isValid`)
or split into helpers. Don't trade readability for a lower score with clever tricks.
Suppress? No.

### nakedret
Return explicit values in functions with named results.

## Style

### wsl_v5 (whitespace linter, non-test code only)
The most frequent finding. Typical fixes:
- Put a blank line before `return` when the block has more than one statement. (nlreturn says
  the same thing.)
- Put a blank line between a variable declaration block and the following `if` or `for`, unless
  the `if` checks the variable that was just assigned. `err := f(); if err != nil` with no blank
  line is fine.
- Don't cuddle unrelated statements. Group them with blank lines.
Fix the layout. Suppress? No.

### nlreturn
Add a blank line before `return`/`break`/`continue` in blocks of more than one statement.

### lll (120 chars, non-test code only)
- Long doc comments: split the URL onto its own line, for example
  `// GetBalances get balances list.` followed by `// See https://developer.dusupay.com/...`.
- Long calls: break after `(`, one argument group per line, closing `)` on its own line, as in
  the multi-line resource method signatures already in the repo.
- Long struct tags: these rarely happen. Shorten the inline comment instead.
Suppress? No.

### gocritic
Follow its suggestion (`ifElseChain` → `switch`, `appendAssign`, `paramTypeCombine` and so on).
For `hugeParam`, consider a pointer receiver or argument, but don't change the signature of an
exported API just for gocritic without telling the user. That's a breaking change.

### goconst (non-test code only)
A string repeated 3 or more times becomes a constant. Check `common.go` first, because a typed
constant may already exist (`TransactionMethodBank`, `CallbackEventTransactionCompleted`...).

### prealloc
`make([]T, 0, len(x))` when the final length is known.

### unconvert / misspell / dogsled
Remove the redundant conversion. Fix the spelling (US English). Avoid `_, _, _ :=`; use only the
values you need.

## Formatter: gofmt
`gofmt -w <file>`. This is mechanical, so always fix it and never suppress it. `stubs/`, `main/`
and `tmp/` are excluded from formatting checks by config.
