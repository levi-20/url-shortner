#!/bin/bash
set -e

MODE="${1:-build}"  # default to "build" if no arg passed

case "$MODE" in
  build)
    go build -o ./bin/release/main main.go
    ;;
  run)
    go build -o ./bin/release/main main.go
    shift  # drop the "run" arg
    ./bin/release/main "$@"
    ;;
  debug)
    go build -gcflags="all=-N -l" -o ./bin/debug/main main.go
    ;;
  *)
    echo "Usage: $0 {build|run|debug}"
    exit 1
    ;;
esac