#!/usr/bin/env bash
# Compare freshly fetched docs with the committed snapshot, or update the snapshot.
#
#   bash .claude/skills/docs-sync/scripts/snapshot.sh compare <fetched_dir>   # what changed since last sync
#   bash .claude/skills/docs-sync/scripts/snapshot.sh update  <fetched_dir>   # record a completed sync
#
# The snapshot lives in .claude/skills/docs-sync/snapshot/ and is committed with the sync,
# so the next run can show exactly which pages and spec fields changed.

set -u
root="$(git rev-parse --show-toplevel)" || exit 1
snap="$root/.claude/skills/docs-sync/snapshot"
cmd="${1:-}"
src="${2:-}"

if [ -z "$cmd" ] || [ -z "$src" ] || [ ! -d "$src" ]; then
	echo "usage: snapshot.sh compare|update <fetched_dir>" >&2
	exit 2
fi

files_of() { # $1 = dir; relative paths of tracked sources
	(cd "$1" && find . -type f \( -name '*.md' -o -name 'llms.txt' -o -name 'openapi.json' -o -name 'postman.json' \) | sed 's#^\./##' | sort)
}

case "$cmd" in
compare)
	if [ ! -d "$snap" ]; then
		echo "No snapshot yet (first sync): treat every page as new and do a full audit."
		exit 0
	fi
	[ -f "$snap/SYNCED_AT" ] && echo "Last sync: $(cat "$snap/SYNCED_AT")"
	added=$(comm -13 <(files_of "$snap") <(files_of "$src"))
	removed=$(comm -23 <(files_of "$snap") <(files_of "$src"))
	changed=""
	for f in $(comm -12 <(files_of "$snap") <(files_of "$src")); do
		cmp -s "$snap/$f" "$src/$f" || changed="$changed $f"
	done
	echo
	echo "==> added:";   [ -n "$added" ] && echo "$added" | sed 's/^/  + /' || echo "  none"
	echo "==> removed:"; [ -n "$removed" ] && echo "$removed" | sed 's/^/  - /' || echo "  none"
	echo "==> changed:"; [ -n "$changed" ] && printf '  ~ %s\n' $changed || echo "  none"
	for f in $changed; do
		case "$f" in
		*.json)
			# JSON is often minified; compare pretty-printed key/value lines instead of raw bytes.
			echo
			echo "--- diff $f (normalized)"
			diff <(tr ',{}' '\n\n\n' <"$snap/$f" | sort -u) <(tr ',{}' '\n\n\n' <"$src/$f" | sort -u) | grep '^[<>]' | head -80
			;;
		*)
			echo
			echo "--- diff $f"
			# Pages keep whole tables on one line, so a line diff is unreadable. Porcelain
			# word diff emits only the changed fragments: "-old words" / "+new words".
			git diff --no-index --no-color --word-diff=porcelain -U0 "$snap/$f" "$src/$f" 2>/dev/null |
				awk '/^(---|\+\+\+) / {next} /^[-+]/ {print "  " $0}' | cut -c1-300 | head -60
			;;
		esac
	done
	;;
update)
	rm -rf "$snap"
	mkdir -p "$snap"
	(cd "$src" && files_of . | while read -r f; do mkdir -p "$snap/$(dirname "$f")"; cp "$f" "$snap/$f"; done)
	echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) (SDK commit $(git -C "$root" rev-parse --short HEAD))" >"$snap/SYNCED_AT"
	echo "snapshot updated: $(files_of "$snap" | wc -l) files"
	;;
*)
	echo "unknown command: $cmd" >&2
	exit 2
	;;
esac
