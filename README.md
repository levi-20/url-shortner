# URL Shortener

Its a project I'm doing for fun this is to shorten the URLs

## Setup

1. Clone `git clone https://github.com/levi-20/url-shortner.git`
2. `go mod tidy`
3. Migration
    Command `migrate create -ext sql -dir migration/db/migrations -seq create_url_tables`.

    Generated file names

    1. `00001_create_url_tables.up.sql`
    2. `00001_create_url_tables.down.sql`
4. Running
    - Windows
      1. Build `Build.bat`
      2. Debug `BuildAndDebug.bat`
      3. Build and Run `BuildAndRun.bat`

    - Linux

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
