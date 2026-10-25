#!/usr/bin/env bash
# Pre-commit quality gate for this repo. Mirrors CI (.github/workflows/tests.yml):
# gofmt, go vet, golangci-lint, go test. Runs every check even if an earlier one
# fails, so a single run reports everything that needs fixing.
# Usage: bash .claude/skills/git-commit/scripts/precheck.sh   (from repo root)
# Exit code: 0 if all checks pass, 1 otherwise.

cd "$(git rev-parse --show-toplevel)" || exit 1

failed=()

step() {
	echo "==> $1"
}

step "gofmt"
unformatted=$(gofmt -l . 2>&1)
if [ -n "$unformatted" ]; then
	echo "Unformatted files (fix with: gofmt -w <file>):"
	echo "$unformatted"
	failed+=("gofmt")
else
	echo "ok"
fi

step "go vet"
if go vet ./... 2>&1; then echo "ok"; else failed+=("go vet"); fi

step "golangci-lint"
if command -v golangci-lint >/dev/null 2>&1; then
	if golangci-lint run ./... 2>&1; then echo "ok"; else failed+=("golangci-lint"); fi
else
	echo "SKIPPED: golangci-lint not installed (CI will still run it)"
fi

step "go test"
if go test ./... 2>&1; then echo "ok"; else failed+=("go test"); fi

echo
if [ ${#failed[@]} -eq 0 ]; then
	echo "PRECHECK PASSED"
	exit 0
fi

echo "PRECHECK FAILED: ${failed[*]}"
exit 1
