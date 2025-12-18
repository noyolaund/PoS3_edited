@echo off
echo ========================================
echo Compilando PoS3 con servidor de impresion
echo ========================================

REM Crear directorio dist si no existe
if not exist "dist" mkdir dist
if not exist "dist\bridge" mkdir dist\bridge

echo.
echo [1/3] Instalando dependencias del bridge...
cd bridge
call npm install
if errorlevel 1 (
    echo Error instalando dependencias
    pause
    exit /b 1
)

echo.
echo [2/3] Compilando servidor de impresion...
call npm run build
if errorlevel 1 (
    echo Error compilando bridge
    pause
    exit /b 1
)

cd ..

echo.
echo [3/3] Compilando backend Go...
cd api
go build -o ..\dist\pos3.exe
if errorlevel 1 (
    echo Error compilando backend
    pause
    exit /b 1
)

cd ..

echo.
echo ========================================
echo Compilacion completada!
echo ========================================
echo.
echo Archivos generados en: dist\
echo   - pos3.exe (backend principal)
echo   - bridge\printer-bridge.exe (servidor de impresion)
echo.
echo Para distribuir, copia toda la carpeta 'dist' con su contenido.
echo.
pause
