@echo off

rem Bypass the STUPID "Terminate Batch Job" prompt
if "%~1"=="-FIXED_CTRL_C" (
   REM Remove the -FIXED_CTRL_C parameter
   SHIFT
) ELSE (
   REM Run the batch with <NUL and -FIXED_CTRL_C
   CALL <NUL %0 -FIXED_CTRL_C %*
   GOTO :EOF
)

REM -o flag is for specifying the output directory
go build -o .\bin\release\main.exe main.go && .\bin\release\main.exe