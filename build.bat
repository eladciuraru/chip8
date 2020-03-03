@echo off

:: Variables
set "ParentDir=%~dp0"
set "BuildDir=%ParentDir%\build"
set "SourceDir=%ParentDir%\cmd"

if not exist "%BuildDir%" (
    echo [*] Creating %BuildDir%
    mkdir "%BuildDir%"
)

pushd %BuildDir%

echo [*] Building disasm
go build "%SourceDir%\disasm"

popd
