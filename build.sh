#!/bin/bash

# -gcflags=all="-N -l" preserves debugging symbols (disables optimizations and inlining)
# -o flag specifies the output path

# Build
# go build -o ./bin/release/main main.go

# Build and Run
# go build -o ./bin/release/main main.go && ./bin/release/main "$@"

# Build and Debug
go build -gcflags="all=-N -l" -o ./bin/debug/main main.go