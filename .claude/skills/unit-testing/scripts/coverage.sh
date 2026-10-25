#!/usr/bin/env bash
# Coverage report for this repo, optionally compared against a git ref.
#
# Usage (from anywhere inside the repo):
#   bash .claude/skills/unit-testing/scripts/coverage.sh            # current tree only
#   bash .claude/skills/unit-testing/scripts/coverage.sh HEAD       # current tree vs HEAD
#   bash .claude/skills/unit-testing/scripts/coverage.sh master     # current tree vs master
#
# The baseline is measured in a temporary detached git worktree, so the user's
# working tree is never stashed, reset or otherwise touched.
# Prints total coverage, the delta vs baseline, and functions below 100% in the
# current tree. Exit code: 0 normally, 2 if coverage dropped vs the baseline,
# 1 if tests failed.

root="$(git rev-parse --show-toplevel)" || exit 1
cd "$root" || exit 1

ref="$1"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

total_of() { # $1 = profile path
	go tool cover -func="$1" | awk '/^total:/ {gsub("%", "", $NF); print $NF}'
}

echo "==> current tree"
if ! go test -coverprofile="$tmp/current.out" . >"$tmp/current.log" 2>&1; then
	cat "$tmp/current.log"
	echo "TESTS FAILED in current tree"
	exit 1
fi
current=$(total_of "$tmp/current.out")
echo "total: ${current}%"

echo
echo "==> functions below 100% (current tree)"
go tool cover -func="$tmp/current.out" | awk '$NF != "100.0%" && !/^total:/' | sed "s#^.*/dusupay-sdk-go/v2/##"

if [ -z "$ref" ]; then
	exit 0
fi

echo
echo "==> baseline: $ref (temporary worktree)"
wt="$tmp/baseline"
if ! git worktree add --detach --quiet "$wt" "$ref" 2>"$tmp/wt.err"; then
	cat "$tmp/wt.err"
	exit 1
fi
# main/ is gitignored but holds the test private key; without it some tests skip
# and the baseline would be artificially low.
[ -d main ] && cp -r main "$wt/" 2>/dev/null

(cd "$wt" && go test -coverprofile="$tmp/base.out" . >"$tmp/base.log" 2>&1)
base_status=$?
git worktree remove --force "$wt" >/dev/null 2>&1

if [ $base_status -ne 0 ]; then
	echo "baseline tests failed at $ref; cannot compare"
	exit 0
fi
base=$(total_of "$tmp/base.out")
echo "total: ${base}%"

echo
delta=$(awk -v c="$current" -v b="$base" 'BEGIN {printf "%+.1f", c - b}')
echo "COVERAGE: ${base}% -> ${current}% (${delta})"
if awk -v c="$current" -v b="$base" 'BEGIN {exit !(c < b)}'; then
	echo "COVERAGE DROPPED"
	exit 2
fi
