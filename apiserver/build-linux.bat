@echo off
setlocal enabledelayedexpansion

set "VERSION=v1.0.0"
if not "%~1"=="" set "VERSION=%~1"

if "%VERSION%"=="" (
    echo ERROR: Image version cannot be empty.
    echo Usage: %~nx0 [version]
    echo Example: %~nx0 v1.2.3
    exit /b 1
)

echo Using version: %VERSION%

set "GOOS=linux"
set "GOARCH=amd64"
set "CGO_ENABLED=0"

if not exist "out" mkdir "out"

echo Building Linux binary...
go build -o out\apiserver main.go
if errorlevel 1 (
    echo ERROR: Go build failed.
    exit /b 1
)
echo OK: Linux binary built at out\apiserver

echo Removing old images...
for /f "tokens=*" %%i in ('docker images --format "{{.Repository}}:{{.Tag}}" ^| findstr /i "openai/apiserver:" 2^>nul') do (
    docker rmi "%%i" >nul 2>&1
)

set "IMAGE_NAME=openai/apiserver:%VERSION%"
echo Building Docker image: %IMAGE_NAME%
docker build -t %IMAGE_NAME% -f deploy\Dockerfile .
if errorlevel 1 (
    echo ERROR: Docker build failed.
    exit /b 1
)
echo OK: Docker built: %IMAGE_NAME%

endlocal