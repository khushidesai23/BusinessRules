@echo off
REM Quick test setup and runner for Windows
REM Usage: run_tests.bat [test-path]

setlocal enabledelayedexpansion

echo.
echo Business Rules Engine - Test Runner (Windows)
echo ========================================

REM Check if Docker is available
where docker >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Docker is not installed. Please install Docker to run tests.
    exit /b 1
)

REM Parse arguments
set "TEST_PATH=%1"
if "!TEST_PATH!"=="" set "TEST_PATH=."

echo.
echo Starting test database...

REM Check if container exists and remove it
docker ps -a | find "business-rules-test-db" >nul 2>&1
if !ERRORLEVEL! EQU 0 (
    echo    (Removing existing container...)
    docker stop business-rules-test-db >nul 2>&1
    docker rm business-rules-test-db >nul 2>&1
)

REM Start containers
docker-compose -f docker-compose.test.yml up -d
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to start the test database with Docker Compose.
    exit /b 1
)

REM Wait for database to be ready
echo    (Waiting for database to be ready...)
set "DB_READY="
for /l %%i in (1,1,30) do (
    docker exec business-rules-test-db pg_isready -U postgres -d postgres >nul 2>&1
    if !ERRORLEVEL! EQU 0 (
        set "DB_READY=1"
        goto :db_ready
    )
    timeout /t 1 /nobreak >nul
)

:db_ready
if not defined DB_READY (
    echo ERROR: Test database did not become ready in time.
    docker-compose -f docker-compose.test.yml down
    exit /b 1
)

REM Ensure the postgres password matches the DSN used by tests
set "DB_CONFIGURED="
for /l %%i in (1,1,30) do (
    docker exec business-rules-test-db psql -U postgres -d postgres -c "ALTER USER postgres PASSWORD 'postgres';" >nul 2>&1
    if !ERRORLEVEL! EQU 0 (
        set "DB_CONFIGURED=1"
        goto :db_configured
    )
    timeout /t 1 /nobreak >nul
)

:db_configured
if not defined DB_CONFIGURED (
    echo ERROR: Failed to configure the test database credentials.
    docker-compose -f docker-compose.test.yml down
    exit /b 1
)

REM Set environment variable for tests
set "POSTGRES_TEST_DSN=host=127.0.0.1 user=postgres password=postgres dbname=postgres port=55432 sslmode=disable"

echo Database ready
echo.

REM Run tests
echo Running tests...
echo.

if "!TEST_PATH!"=="." (
    echo Running: go test ./... -v
    go test ./... -v
) else (
    echo Running: go test !TEST_PATH! -v
    go test !TEST_PATH! -v
)

set "TEST_RESULT=%ERRORLEVEL%"

echo.
echo ========================================

REM Cleanup
echo Cleaning up...
docker-compose -f docker-compose.test.yml down

if %TEST_RESULT% EQU 0 (
    echo.
    echo All tests passed!
    exit /b 0
) else (
    echo.
    echo Some tests failed
    exit /b 1
)
