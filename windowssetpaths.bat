@echo off

set "CURRENT_DIR=%CD%"

echo %PATH% | findstr /C:"%CURRENT_DIR%" >nul
if %errorlevel% neq 0 (
    echo export PATH="%PATH%;%CURRENT_DIR%" >> %USERPROFILE%\bashrc
    echo Successfully added PATH
) else (
    echo PATH Already added
)

echo %LATER_PROJECT_DIR% | findstr /C:"%CURRENT_DIR%" >nul
if %errorlevel% neq 0 (
    echo export LATER_PROJECT_DIR="%CURRENT_DIR%" >> %USERPROFILE%\bashrc
    echo Successfully added LATER_PROJECT_DIR
) else (
    echo LATER_PROJECT_DIR Already added
)

bash
