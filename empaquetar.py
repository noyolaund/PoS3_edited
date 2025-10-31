import argparse
import datetime
import os
import shutil
from subprocess import check_output
import hashlib

NOMBRE_DIRECTORIO_DIST = "dist"
NOMBRE_DIRECTORIO_API = "api"
NOMBRE_DE_EJECUTABLE_SIN_EXTENSION = "sublime_pos_3_by_parzibyte"
NOMBRE_ARCHIVO_CHECKSUM = "checksum_archivos.go"

fecha_y_hora = datetime.datetime.now().strftime("%Y-%m-%d_%H-%M-%S")

directorio_inicial = os.getcwd()

parser = argparse.ArgumentParser()
parser.add_argument(
    "arquitectura", help="Arquitectura para la que se compila. Puede ser 32 o 64 bits, pero si compilas para 32 bits recuerda modificar la PATH")
argumentos = parser.parse_args()
arquitecturas = ["32", "64"]
if not argumentos.arquitectura in arquitecturas:
    print(
        f"La arquitectura debe ser una de: {arquitecturas} pero {argumentos.arquitectura} fue proporcionado")
    exit()

arquitectura = argumentos.arquitectura
print("""
Si vas a compilar para 32 bits, probablemente quieras:
SET PATH=C:\Go32\go\\bin;C:\MinGW\\bin;%PATH% && SET GOROOT=C:\Go32\go\
""")


def hash_file(filename):

    # make a hash object
    h = hashlib.sha512()

    # open file for reading in binary mode
    with open(filename, 'rb') as file:

        # loop till the end of the file
        chunk = 0
        while chunk != b'':
            # read only 1024 bytes at a time
            chunk = file.read(1024)
            h.update(chunk)

    # return the hex representation of digest
    return h.hexdigest()


def obtener_codigo_go_para_checksum(archivos_checksum):
    codigo = """
package main

/*
    Este código es autogenerado por el script de Python. No se recomienda tocarlo,
    pero si quieres puedes hacerlo. Solo es un mapa en donde la clave es la ubicación relativa
    del archivo, y el valor es el SHA512 de ese archivo

    Generado el: """+fecha_y_hora + """
*/

var ubicacionesConHash = map[string]string{
"""

    for archivo in archivos_checksum:
        codigo += f"\t\"./{NOMBRE_DIRECTORIO_DIST}/{archivo['ubicacion']}\" : \"{archivo['hash']}\",\n"

    codigo += """
}"""
    return codigo


def obtener_lista_archivos_para_checksum(directorio_dist):
    archivos = []
    for directorio in os.listdir(directorio_dist):
        ruta_completa = os.path.join(directorio_dist, directorio)
        if os.path.isfile(ruta_completa):
            archivos.append({
                "hash": hash_file(ruta_completa),
                "ubicacion": directorio,
            })
        if os.path.isdir(ruta_completa):
            for subdirectorio in os.listdir(ruta_completa):
                if os.path.isfile(os.path.join(ruta_completa, subdirectorio)):
                    archivos.append({
                        "hash": hash_file(os.path.join(ruta_completa, subdirectorio)),
                        "ubicacion": directorio+"/" + subdirectorio,
                    })
    return archivos


nombre_carpeta_release = os.path.join(
    directorio_inicial, "sublime_pos_3_{}_{}".format(arquitectura, fecha_y_hora))

print(f"Creando carpeta {nombre_carpeta_release}...")
if os.path.exists(nombre_carpeta_release):
    shutil.rmtree(nombre_carpeta_release)
os.mkdir(nombre_carpeta_release)

# Compilar cliente
comando_compilar_cliente = "npm run build"
print(f"Compilando cliente con {comando_compilar_cliente}...")
check_output(comando_compilar_cliente, shell=True)


directorio_dist = os.path.join(directorio_inicial, NOMBRE_DIRECTORIO_DIST)
directorio_dist_dentro_de_release = os.path.join(
    nombre_carpeta_release, NOMBRE_DIRECTORIO_DIST)
print(
    f"Copiando archivos de cliente {directorio_dist} a carpeta release {nombre_carpeta_release}...")
shutil.copytree(directorio_dist, directorio_dist_dentro_de_release)
directorio_fonts = os.path.join(directorio_dist_dentro_de_release, "fonts")
directorio_css = os.path.join(directorio_dist_dentro_de_release, "css")
print("Moviendo fonts a css")
shutil.copytree(directorio_fonts, os.path.join(directorio_css, "fonts"))
shutil.rmtree(directorio_fonts)

print("Obteniendo hash de archivos dist")
archivos_checksum = obtener_lista_archivos_para_checksum(
    directorio_dist_dentro_de_release)

ruta_api = os.path.join(directorio_inicial, NOMBRE_DIRECTORIO_API)
ruta_archivo_checksum = os.path.join(ruta_api, NOMBRE_ARCHIVO_CHECKSUM)
print(f"Escribiendo código Golang en {ruta_archivo_checksum}")
with open(ruta_archivo_checksum, "w+", encoding="utf-8") as archivo:
    archivo.write(obtener_codigo_go_para_checksum(archivos_checksum))

print(f"Cambiando directorio a {ruta_api}")
os.chdir(ruta_api)
nombre_ejecutable = f"{NOMBRE_DE_EJECUTABLE_SIN_EXTENSION}_{arquitectura}.exe"
comando = f"go build -tags produccion -ldflags \"-H windowsgui\" -o {nombre_ejecutable}"
print(f"Compilando API con {comando}...")
check_output(comando, shell=True)
print(f"Cambiando directorio a {directorio_inicial}")
os.chdir(directorio_inicial)
print(f"Copiando archivos de {ruta_api} a {nombre_carpeta_release}")
archivos_api = [nombre_ejecutable, "init.sql",
                "esquema_negocios_sqlite.sql", "esquema_spos_sqlite.sql"]
for archivo in archivos_api:
    ruta_completa = os.path.join(ruta_api, archivo)
    print(f"Copiando {ruta_completa}...")
    shutil.copyfile(ruta_completa, os.path.join(
        nombre_carpeta_release, archivo))

print(f"Creando {nombre_carpeta_release}.zip...")
shutil.make_archive(nombre_carpeta_release, "zip", nombre_carpeta_release)
print(f"Removiendo {nombre_carpeta_release}")
shutil.rmtree(nombre_carpeta_release)
print(f"""
  ==========================================
             ___  _ __ _
            | . || / /| |
            | | ||  \ |_/
            `___'|_\_\<_>

Ya puedes distribuir {nombre_carpeta_release}.zip
  ==========================================
  """)
