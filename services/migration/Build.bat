@echo off

REM -gcflags=all=-N -l is used to preserve debugging symbols in the build binary
REM -o flag is for specifying the output directory
go build -o .\bin\release\main.exe main.go