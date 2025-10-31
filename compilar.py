from subprocess import call
import shutil
import os
import datetime
import platform

sistema = platform.system()

print("""
    =================================================
            DETECTADO SISTEMA '{}'
    =================================================""".format(sistema))
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
#Sacar los archivos a la carpeta static

ruta_css = ruta_destino_compilacion + "/static/css"
shutil.move(ruta_css, ruta_destino_compilacion)



#Mover las fuentes, sólo en linux pero creo que es por lo de postcss
if sistema == "Linux":
    ruta_fuentes = ruta_destino_compilacion + "/fonts"
    shutil.move(ruta_fuentes, ruta_destino_compilacion + "/css/fonts/")

ruta_img = ruta_destino_compilacion + "/static/img"
shutil.move(ruta_img, ruta_destino_compilacion)

ruta_js = ruta_destino_compilacion + "/static/js"
shutil.move(ruta_js, ruta_destino_compilacion)

print("Copiando el favicon...")
shutil.copy(ruta_absoluta + "/src/assets/inicio/logo.png", ruta_absoluta + "/api/static/img/logo.png")

print("Eliminando static, pues es un directorio vacío")
#Eliminar static
eliminar_si_existe(ruta_destino_compilacion + "/static")

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
if sistema == "Windows":
    call(["cmd", "/c", "go", "build", "-tags", "produccion"])
else:
    call(["go", "build", "-tags", "produccion"])




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
