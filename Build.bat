@echo off
setlocal

REM Default to "build" if no argument passed
set MODE=%1
if "%MODE%"=="" set MODE=build

if "%MODE%"=="build" (
    go build -o .\bin\release\main.exe main.go
    goto :eof
)

if "%MODE%"=="run" (
    go build -o .\bin\release\main.exe main.go
    if errorlevel 1 goto :eof
    shift
    .\bin\release\main.exe %*
    goto :eof
)

if "%MODE%"=="debug" (
    go build -gcflags="all=-N -l" -o .\bin\debug\main.exe main.go
    goto :eof
)

echo Usage: %0 {build^|run^|debug}
exit /b 1