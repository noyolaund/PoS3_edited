@echo off
setlocal enabledelayedexpansion
chcp 65001 > nul
color 0A

echo ====================================
echo   Sistema POS - Inicializador
echo ====================================
echo.
echo Iniciando servicios...
echo.

:: Variables para tracking de procesos
set BACKEND_PID=
set BRIDGE_PID=
set BACKEND_OK=0
set BRIDGE_OK=0

:: ========================================
:: 1. Iniciar Backend (API Go)
:: ========================================
echo [1/2] Iniciando Backend API...
cd /d "%~dp0"

if not exist "SPOS3.exe" (
    echo [ERROR] No se encontró SPOS3.exe en la raíz
    goto :ERROR_EXIT
)

:: Iniciar backend en nueva ventana
start "POS Backend API" cmd /c "SPOS3.exe"
timeout /t 3 /nobreak > nul

:: Verificar que el backend esté respondiendo
echo Verificando Backend API en localhost:2106...
for /L %%i in (1,1,10) do (
    curl -s http://localhost:2106/auth/salud > nul 2>&1
    if !errorlevel! equ 0 (
        set BACKEND_OK=1
        echo [OK] Backend API iniciado correctamente
        goto :BACKEND_SUCCESS
    )
    if %%i lss 10 (
        echo Esperando Backend API... intento %%i/10
        timeout /t 2 /nobreak > nul
    )
)

:BACKEND_SUCCESS
if !BACKEND_OK! equ 0 (
    echo [ERROR] Backend API no responde después de 10 intentos
    echo Verifica que el puerto 2106 esté disponible
    goto :ERROR_EXIT
)

echo.

:: ========================================
:: 2. Iniciar Printer Bridge (Node.js)
:: ========================================
echo [2/2] Iniciando Printer Bridge...
cd /d "%~dp0bridge"

if not exist "bridge.exe" (
    echo [ERROR] No se encontró bridge.exe en la carpeta bridge
    goto :ERROR_EXIT
)

:: Iniciar bridge en nueva ventana
start "POS Printer Bridge" cmd /c "bridge.exe"
timeout /t 3 /nobreak > nul

:: Verificar que el bridge esté respondiendo
echo Verificando Printer Bridge en localhost:8001...
for /L %%i in (1,1,10) do (
    curl -s http://localhost:8001/health > nul 2>&1
    if !errorlevel! equ 0 (
        set BRIDGE_OK=1
        echo [OK] Printer Bridge iniciado correctamente
        goto :BRIDGE_SUCCESS
    )
    if %%i lss 10 (
        echo Esperando Printer Bridge... intento %%i/10
        timeout /t 2 /nobreak > nul
    )
)

:BRIDGE_SUCCESS
if !BRIDGE_OK! equ 0 (
    echo [ERROR] Printer Bridge no responde después de 10 intentos
    echo Verifica que el puerto 3001 esté disponible
    goto :ERROR_EXIT
)

echo.

:: ========================================
:: 3. Verificación Final
:: ========================================
if !BACKEND_OK! equ 1 if !BRIDGE_OK! equ 1 (
    color 0A
    echo ====================================
    echo   ✓ TODOS LOS SERVICIOS INICIADOS
    echo ====================================
    echo.
    echo Backend API:      http://localhost:2106
    echo Printer Bridge:   http://localhost:8001
    echo.
    echo Ambos servicios están corriendo en ventanas separadas.
    echo Puedes cerrar esta ventana.
    echo.
    timeout /t 3 /nobreak > nul
    exit /b 0
)

goto :ERROR_EXIT

:: ========================================
:: Manejo de Errores
:: ========================================
:ERROR_EXIT
color 0C
echo.
echo ====================================
echo   ✗ ERROR AL INICIAR SERVICIOS
echo ====================================
echo.
echo Estado de servicios:
if !BACKEND_OK! equ 1 (
    echo   Backend API:      [OK] ✓
) else (
    echo   Backend API:      [ERROR] ✗
)
if !BRIDGE_OK! equ 1 (
    echo   Printer Bridge:   [OK] ✓
) else (
    echo   Printer Bridge:   [ERROR] ✗
)
echo.
echo Posibles causas:
echo   - Puertos 2106 o 8001 ya están en uso
echo   - Dependencias no instaladas (Go, Node.js)
echo   - Archivos de configuración faltantes
echo.
echo Presiona cualquier tecla para cerrar...
pause > nul
exit /b 1
