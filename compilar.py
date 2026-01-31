from subprocess import call
import shutil
import os
import datetime
import platform
import struct

sistema = platform.system()

print("""
    =================================================
            DETECTADO SISTEMA '{}'
    =================================================""".format(sistema))

# Detectar arquitectura
arquitectura_bits = struct.calcsize("P") * 8
arquitectura_maquina = platform.machine()
print("Arquitectura detectada: {} ({} bits)".format(arquitectura_maquina, arquitectura_bits))
def eliminar_si_existe(ruta):
    if os.path.exists(ruta):
        shutil.rmtree(ruta)

if input("""¿Deseas compilar el código JavaScript?
No deberías hacerlo si lo has compilado recientemente
(por ejemplo, si compilaste con Go la versión para 64 bits y
ahora estás compilando para 32)

Deberías hacerlo si has hecho cambios a los archivos de js, css, html o cualquier cosa del lado del cliente

[s/n] """).lower().find("s") != -1:
    #Compilar javascript
    print("Compilando el código JavaScript...")
    if sistema == "Windows":
        call(["cmd", "/c" ,"npm", "run", "build"])
    else:
        call(["npm", "run", "build"])

#Eliminar static
print("Eliminando static creado anteriormente...")
ruta_absoluta = os.path.dirname(os.path.abspath(__file__))
ruta_destino_compilacion = ruta_absoluta + "/api/static";
eliminar_si_existe(ruta_destino_compilacion)

print("Copiando nuevos archivos...")

#Copiar lo de dist a static
ruta_actual_dist = ruta_absoluta + "/dist"
shutil.copytree(ruta_actual_dist, ruta_destino_compilacion)



print("Copiando archivos a la carpeta padre...")

# Definir rutas de origen
ruta_static = ruta_destino_compilacion + "/static"
ruta_css_src = ruta_static + "/css"
ruta_js_src = ruta_static + "/js"
ruta_img_src = ruta_static + "/img"

# Definir rutas de destino
ruta_css_dest = ruta_destino_compilacion + "/css"
ruta_js_dest = ruta_destino_compilacion + "/js"
ruta_img_dest = ruta_destino_compilacion + "/img"

# Mover CSS si existe
if os.path.exists(ruta_css_src):
    if os.path.exists(ruta_css_dest):
        shutil.rmtree(ruta_css_dest)
    shutil.move(ruta_css_src, ruta_css_dest)
    print("CSS movido correctamente")
else:
    print("Advertencia: No se encontró carpeta CSS en {}".format(ruta_css_src))

# Mover fuentes en Linux si existen
if sistema == "Linux":
    ruta_fuentes = ruta_destino_compilacion + "/fonts"
    if os.path.exists(ruta_fuentes):
        ruta_css_fonts = ruta_css_dest + "/fonts"
        if not os.path.exists(ruta_css_fonts):
            os.makedirs(ruta_css_fonts)
        for item in os.listdir(ruta_fuentes):
            shutil.move(os.path.join(ruta_fuentes, item), os.path.join(ruta_css_fonts, item))
        os.rmdir(ruta_fuentes)
        print("Fuentes movidas correctamente")

# Mover IMG si existe
if os.path.exists(ruta_img_src):
    if os.path.exists(ruta_img_dest):
        shutil.rmtree(ruta_img_dest)
    shutil.move(ruta_img_src, ruta_img_dest)
    print("IMG movida correctamente")
else:
    print("Advertencia: No se encontró carpeta IMG en {}".format(ruta_img_src))

# Mover JS si existe
if os.path.exists(ruta_js_src):
    if os.path.exists(ruta_js_dest):
        shutil.rmtree(ruta_js_dest)
    shutil.move(ruta_js_src, ruta_js_dest)
    print("JS movido correctamente")
else:
    print("Advertencia: No se encontró carpeta JS en {}".format(ruta_js_src))

print("Copiando el favicon...")
shutil.copy(ruta_absoluta + "/src/assets/inicio/logo.png", ruta_absoluta + "/api/static/img/logo.png")

print("Eliminando static, pues es un directorio vacío")
#Eliminar static
ruta_static_vacia = ruta_destino_compilacion + "/static"
if os.path.exists(ruta_static_vacia):
    shutil.rmtree(ruta_static_vacia)
    print("Carpeta static eliminada")

# Eliminar archivo index.html si existe (no lo necesitamos)
ruta_index = ruta_destino_compilacion + "/index.html"
if os.path.exists(ruta_index):
    os.remove(ruta_index)
    print("archivo index.html eliminado")

#Eliminar archivos js

ruta_js = ruta_destino_compilacion + "/js"
ruta_css = ruta_destino_compilacion + "/css"
print("Limpiando carpeta js...")
for raiz, directorio, archivos in os.walk(ruta_js):
    for archivo in archivos:
        if archivo.endswith(".js.map"):
            os.remove(ruta_js + "/" + archivo)
            print("Se ha eliminado {} porque era un archivo MAP".format(archivo))

print("Limpiando carpeta css...")
for raiz, directorio, archivos in os.walk(ruta_css):
    for archivo in archivos:
        if archivo.endswith(".css.map"):
            os.remove(ruta_css + "/" + archivo)
            print("Se ha eliminado {} porque era un archivo MAP".format(archivo))
print("Moviéndonos a la carpeta api")
os.chdir("api")
print("Compilando el código en un archivo exe (este es el proceso más tardado)")

# Configurar variables de entorno para go build
env = os.environ.copy()

if sistema == "Windows":
    print("Compilando para Windows x64...")
    call(["cmd", "/c", "go", "build", "-tags", "produccion", "-o", "api.exe"])
else:
    # En Linux, detectar arquitectura y configurar Go
    env["CGO_ENABLED"] = "1"
    env["GOOS"] = "linux"
    
    if "arm" in arquitectura_maquina.lower() or "aarch64" in arquitectura_maquina.lower():
        # Compilación nativa en ARM
        if arquitectura_bits == 64 or "aarch64" in arquitectura_maquina.lower():
            env["GOARCH"] = "arm64"
            print("Compilando para ARM 64-bit (Raspberry Pi 64-bit)...")
        else:
            env["GOARCH"] = "arm"
            env["GOARM"] = "7"
            print("Compilando para ARM 32-bit (Raspberry Pi 32-bit)...")
    elif "x86_64" in arquitectura_maquina.lower() or "amd64" in arquitectura_maquina.lower():
        env["GOARCH"] = "amd64"
        print("Compilando para AMD64 64-bit...")
    else:
        print("Compilando con arquitectura nativa detectada: {}".format(arquitectura_maquina))
    
    print("Variables de entorno configuradas:")
    print("  GOOS: {}".format(env.get("GOOS", "no-configurado")))
    print("  GOARCH: {}".format(env.get("GOARCH", "no-configurado")))
    if "GOARM" in env:
        print("  GOARM: {}".format(env["GOARM"]))
    print("  CGO_ENABLED: {}".format(env.get("CGO_ENABLED", "no-configurado")))
    
    # Ejecutar go build con build tag produccion
    print("Ejecutando: go build -tags produccion -o api")
    call(["go", "build", "-tags", "produccion", "-o", "api"], env=env)




print("Creando copia")
extension = ".exe"
if sistema == "Linux":
    extension = ""
if os.path.isfile("api"+extension):
    nombre = ("SPOS3_{}"+extension).format(datetime.datetime.now().strftime("%Y-%m-%d_%H-%M-%S"))
    shutil.copy("api"+extension, nombre)

print("""
    =================================================
    |                                               |
    |         Compilado correctamente :)            |
    |                                               |
    =================================================
Ahora toma el archivo {} y la carpeta static con todo su contenido,
es lo único que necesitas para ejecutar el programa
    """.format(nombre))
