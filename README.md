# URL Shortener

Its a project I'm doing for fun this is to shorten the URLs

## how to run

### Windows

1. Build `Build.bat`
2. Debug `BuildAndDebug.bat`
3. Build and Run `BuildAndRun.bat`

### Linux

`build.sh` contains all the commands can be used as per need

```sh
#!/bin/bash

# -gcflags=all="-N -l" preserves debugging symbols (disables optimizations and inlining)
# -o flag specifies the output path

# Build
# go build -o ./bin/release/main main.go

# Build and Run
# go build -o ./bin/release/main main.go && ./bin/release/main "$@"

# Build and Debug
go build -gcflags="all=-N -l" -o ./bin/debug/main main.go
```