#!/bin/bash

set -euo pipefail

main_file="$1"
pr_file="$2"
threshold=20

echo "Comparing benchmark results..."

benchstat "$main_file" "$pr_file" > result.txt
cat result.txt

# Fail if regression exceeds threshold
fail=0
awk -v threshold=$threshold '
/±/ {
  split($1, name, "-")
  name = name[1]
  old = $3
  new = $6
  delta = $7

  if (match(delta, /\+([0-9.]+)%/, m)) {
    change = m[1]
    if (change > threshold) {
      printf "❌ %s regressed by %s%%\n", name, change
      exit 1
    }
  }
}
' result.txt || fail=1

if [ "$fail" -eq 1 ]; then
  echo "🚨 Performance regression > $threshold% detected!"
  exit 1
else
  echo "✅ No significant regressions detected."
fi