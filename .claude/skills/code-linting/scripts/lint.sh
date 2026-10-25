#!/usr/bin/env bash
# Run golangci-lint for this repo and split findings into "changed files" vs "pre-existing".
#
# Usage (from anywhere inside the repo):
#   bash .claude/skills/code-linting/scripts/lint.sh        # local golangci-lint binary
#   bash .claude/skills/code-linting/scripts/lint.sh --ci   # exact version pinned in CI (via go run)
#   bash .claude/skills/code-linting/scripts/lint.sh --ci --fix   # extra args go to golangci-lint
#
# "Changed files" = modified/staged vs HEAD plus untracked .go files.
# Exit code: golangci-lint's own (0 = clean, 1 = issues found, other = run error).

root="$(git rev-parse --show-toplevel)" || exit 1
cd "$root" || exit 1

use_ci=0
args=()
for a in "$@"; do
	if [ "$a" = "--ci" ]; then use_ci=1; else args+=("$a"); fi
done

ci_version=$(grep -oE 'version: *v[0-9]+\.[0-9]+\.[0-9]+' .github/workflows/tests.yml | head -1 | grep -oE 'v[0-9.]+')

if [ $use_ci -eq 1 ]; then
	if [ -z "$ci_version" ]; then
		echo "Could not read the golangci-lint version from .github/workflows/tests.yml"
		exit 3
	fi
	echo "==> golangci-lint $ci_version (CI-pinned, via go run; first run downloads it)"
	lint=(go run "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$ci_version")
else
	if ! command -v golangci-lint >/dev/null 2>&1; then
		echo "golangci-lint not installed; re-run with --ci to use the CI-pinned version via go run"
		exit 3
	fi
	local_version=$(golangci-lint version 2>/dev/null | grep -oE 'version [v]?[0-9.]+' | grep -oE '[0-9.]+$')
	echo "==> golangci-lint $local_version (local)"
	if [ -n "$ci_version" ] && [ "v$local_version" != "$ci_version" ]; then
		echo "WARNING: CI pins $ci_version; findings may differ. Use --ci to match CI exactly."
	fi
	lint=(golangci-lint)
fi

out=$(mktemp)
trap 'rm -f "$out"' EXIT
"${lint[@]}" run ./... "${args[@]}" >"$out" 2>&1
status=$?

changed=$( (git diff --name-only HEAD; git ls-files --others --exclude-standard -- '*.go') 2>/dev/null | sort -u)

# A finding starts with "path/file.go:line:col:"; the following source/caret lines belong to it.
awk -v changed="$changed" '
	BEGIN { n = split(changed, c, "\n"); for (i = 1; i <= n; i++) if (c[i] != "") is_changed[c[i]] = 1 }
	/^[^ \t]+\.go:[0-9]+:[0-9]+: / {
		f = $0; sub(/:.*/, "", f); gsub(/\\/, "/", f)
		bucket = (f in is_changed) ? "new" : "old"
	}
	/^[0-9]+ issues?[.:]?$/ || /^\* / { summary = summary $0 "\n"; next }
	bucket == "new" { new_out = new_out $0 "\n"; next }
	bucket == "old" { old_out = old_out $0 "\n"; next }
	{ other = other $0 "\n" }
	END {
		if (other != "") printf "%s", other
		print "\n==> findings in changed files"
		printf "%s", (new_out == "" ? "none\n" : new_out)
		print "\n==> pre-existing findings (files not changed vs HEAD)"
		printf "%s", (old_out == "" ? "none\n" : old_out)
		print "\n==> summary"
		printf "%s", summary
	}
' "$out"

exit $status
