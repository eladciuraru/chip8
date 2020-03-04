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

:: Check first parameter
if "%ToolName%"=="" (
    for /d %%n in ("%SourceDir%\*") do (
        echo [*] Building %%~nn
        go build %%n
    )
) else (
    echo [*] Building %ToolName%
    go build "%SourceDir%\%ToolName%"
)

popd
