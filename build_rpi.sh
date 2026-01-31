#!/usr/bin/env bash
set -euo pipefail

# build_rpi.sh — Compila frontend y backend en modo producción para Raspberry Pi
# Agregamos la ruta de Go al PATH para que el script lo encuentre
export PATH=$PATH:/home/bicho/.local/share/go/bin
# Uso:
#   ./build_rpi.sh                   -> compila nativamente en la máquina actual (ideal si ya estás en la Pi)
#   ./build_rpi.sh --cross armv7     -> cross-compila el binario Go para ARMv7 (Raspberry Pi 3/4 32-bit)
#   ./build_rpi.sh --cross arm64     -> cross-compila para ARM64 (Raspberry Pi 3/4 64-bit)
#   ./build_rpi.sh --clean           -> limpia build previos antes de compilar
#
# Requisito: ejecutar desde la raíz del proyecto (donde está package.json). Si el frontend
# está en una subcarpeta, ajusta FRONTEND_DIR.

PROGNAME="$(basename "$0")"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"
FRONTEND_DIR="$PROJECT_ROOT"   # <- ajusta si tu frontend está en subcarpeta
API_DIR="$PROJECT_ROOT/api"
BIN_DIR="$PROJECT_ROOT/bin"
RELEASE_DIR="$PROJECT_ROOT/release"

CROSS_MODE=""
CLEAN_FIRST=false

show_help() {
  cat <<EOF
$PROGNAME — Compila frontend + backend para Raspberry Pi

Opciones:
  --cross <armv7|arm64>   Cross-compile del binario Go para la arquitectura indicada.
  --clean                 Elimina carpetas build/bin/dist antes de compilar.
  -h, --help              Muestra esta ayuda.
EOF
}

# Parse args
while [[ $# -gt 0 ]]; do
  case "$1" in
    --cross)
      CROSS_MODE="${2:-}"
      shift 2
      ;;
    --clean)
      CLEAN_FIRST=true
      shift
      ;;
    -h|--help)
      show_help
      exit 0
      ;;
    *)
      echo "Opción desconocida: $1"
      show_help
      exit 1
      ;;
  esac
done

echo "==> Iniciando build en: $PROJECT_ROOT"

# Verificar que estamos en el directorio correcto
if [ ! -f "$PROJECT_ROOT/package.json" ]; then
  echo "ERROR: No se encontró package.json en $PROJECT_ROOT"
  echo "Asegúrate de ejecutar este script desde la raíz del proyecto"
  exit 1
fi

if $CLEAN_FIRST; then
  echo "==> Limpieza previa..."
  rm -rf "$BIN_DIR" "$RELEASE_DIR" "$PROJECT_ROOT/dist" "$PROJECT_ROOT/build" node_modules
fi

# Verificar herramientas mínimas
MISSING=()
for cmd in node npm tar git; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    MISSING+=("$cmd")
  fi
done
if [ ${#MISSING[@]} -ne 0 ]; then
  echo "ERROR: faltan comandos: ${MISSING[*]}"
  echo "Instala: sudo apt update && sudo apt install -y build-essential git curl tar"
  exit 1
fi

# 1) Instalar dependencias Node (IMPORTANTE: incluir devDependencies para el proceso de build)
echo "==> Instalando dependencias Node (incluyendo devDependencies para el build)..."
cd "$FRONTEND_DIR"

# Si NODE_ENV está exportado, lo quitamos para que npm instale devDependencies
if [ "${NODE_ENV:-}" = "production" ]; then
  unset NODE_ENV
fi

if [ -f package-lock.json ]; then
  npm ci --prefer-offline --no-audit
else
  npm install --prefer-offline --no-audit
fi

# 2) Build frontend en modo producción
echo "==> Ejecutando build frontend (NODE_ENV=production)..."
echo "    Directorio de trabajo: $(pwd)"
echo "    Verificando que build/build.js existe..."
if [ ! -f "build/build.js" ]; then
  echo "ERROR: No se encuentra build/build.js en $(pwd)/build/build.js"
  ls -la build/ || echo "No existe la carpeta build"
  exit 1
fi

NODE_ENV=production npm run build || {
  echo "ERROR: 'npm run build' falló"
  exit 1
}

# Determinar carpeta de salida del frontend
if [ -d "$PROJECT_ROOT/dist" ]; then
  FRONT_DIST="$PROJECT_ROOT/dist"
elif [ -d "$PROJECT_ROOT/build" ]; then
  FRONT_DIST="$PROJECT_ROOT/build"
else
  FRONT_DIST=""
  echo "ADVERTENCIA: no se encontró 'dist/' ni 'build/'. Ajusta el proyecto si usa otra carpeta."
fi
echo "Frontend output: ${FRONT_DIST:-(no-detectada)}"

# 3) Compilar backend Go (si existe)
if [ -d "$API_DIR" ]; then
  echo "==> Compilando backend Go..."
  mkdir -p "$BIN_DIR"
  pushd "$API_DIR" >/dev/null

  # Cross-compile si se pidió
  if [ -n "$CROSS_MODE" ]; then
    case "$CROSS_MODE" in
      armv7)
        export GOOS=linux
        export GOARCH=arm
        export GOARM=7
        export CC=arm-linux-gnueabihf-gcc
        export CGO_ENABLED=1
        echo "Cross-compiling: GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM"
        ;;
      arm64)
        export GOOS=linux
        export GOARCH=arm64
        export CC=aarch64-linux-gnu-gcc
        export CGO_ENABLED=1
        echo "Cross-compiling: GOOS=$GOOS GOARCH=$GOARCH"
        ;;
      *)
        echo "Arquitectura de cross desconocida: $CROSS_MODE"
        popd >/dev/null
        exit 1
        ;;
    esac
  else
    # Compilación nativa - detectar arquitectura actual
    CURRENT_ARCH=$(uname -m)
    unset GOOS GOARCH GOARM
    
    # Si estamos en ARM y pedimos arm64, usar compilación nativa
    if [ "$CURRENT_ARCH" = "aarch64" ] || [ "$CURRENT_ARCH" = "arm64" ]; then
      echo "Compilación nativa en ARM64"
      export GOOS=linux
      export GOARCH=arm64
      export CGO_ENABLED=1
    elif [ "$CURRENT_ARCH" = "armv7l" ] || [ "$CURRENT_ARCH" = "armv6l" ]; then
      echo "Compilación nativa en ARM 32-bit"
      export GOOS=linux
      export GOARCH=arm
      export GOARM=7
      export CGO_ENABLED=1
    else
      echo "Compilación nativa en $(uname -m)"
    fi
  fi

  OUT_BIN="$BIN_DIR/pos_api"
  LDFLAGS="-s -w"
  
  # Debug: mostrar qué archivos existen
  echo "Archivos Go en la carpeta actual:"
  ls -1 *.go | grep -E "(Produccion|Desarrollo)" || echo "  (No encontrados Produccion.go o Desarrollo.go)"
  echo ""
  
  # Compila con build tag produccion - especificar explícitamente el paquete actual
  echo "Ejecutando: go build -tags produccion -v -trimpath -ldflags \"$LDFLAGS\" -o \"$OUT_BIN\" ."
  go build -tags produccion -v -trimpath -ldflags "$LDFLAGS" -o "$OUT_BIN" . || {
    echo "ERROR: go build falló"
    popd >/dev/null
    exit 1
  }

  popd >/dev/null
  echo "Backend: binario creado en $OUT_BIN"
else
  echo "No existe '$API_DIR' — omitiendo compilación Go."
fi

# 4) Empaquetar artifacts
echo "==> Empaquetando artifacts..."
mkdir -p "$RELEASE_DIR"
DATE=$(date +%Y%m%d_%H%M%S)
TARFILE="$RELEASE_DIR/pos3_rpi_${DATE}.tar.gz"
FILES_TO_TAR=()

[ -d "$BIN_DIR" ] && FILES_TO_TAR+=("bin")
[ -n "$FRONT_DIST" ] && [ -d "$FRONT_DIST" ] && FILES_TO_TAR+=("$(basename "$FRONT_DIST")")
[ -f "$PROJECT_ROOT/.env" ] && FILES_TO_TAR+=(".env")

if [ ${#FILES_TO_TAR[@]} -eq 0 ]; then
  echo "Nada para empaquetar. Comprueba que el build se ejecutó correctamente."
else
  pushd "$PROJECT_ROOT" >/dev/null
  tar -czf "$TARFILE" "${FILES_TO_TAR[@]}"
  popd >/dev/null
  echo "Paquete creado: $TARFILE"
fi

echo "==> Build completado."
[ -d "$BIN_DIR" ] && echo "  - Binarios: $BIN_DIR"
[ -n "$FRONT_DIST" ] && [ -d "$FRONT_DIST" ] && echo "  - Frontend: $FRONT_DIST"
echo "  - Paquetes: $RELEASE_DIR"
exit 0
