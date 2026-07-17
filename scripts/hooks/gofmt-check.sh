#!/usr/bin/env bash
# Checks that staged Go files are gofmt-formatted. Invoked by lefthook
# pre-commit (see lefthook.yml). Kept as a standalone script rather than an
# inline `run:` block because lefthook on Windows passes multi-line scripts
# through sh.exe -c "<script>", and embedded double quotes there terminate
# the outer string early.
set -euo pipefail

unformatted="$(gofmt -l "$@")"
if [ -n "$unformatted" ]; then
  echo "gofmt found unformatted files:"
  echo "$unformatted"
  echo "run: gofmt -w <file>"
  exit 1
fi
