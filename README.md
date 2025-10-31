# tiendas

> Sistema de ventas personalizado

# Versión de node y npm
Probado con:
```
C:\Users\parzibyte>node -v
v14.17.0

C:\Users\parzibyte>npm -v
9.6.7

C:\Windows\System32>go version
go version go1.17 windows/amd64

C:\Windows\System32>gcc -v
Using built-in specs.
COLLECT_GCC=gcc
COLLECT_LTO_WRAPPER=C:/msys64/mingw64/bin/../lib/gcc/x86_64-w64-mingw32/13.2.0/lto-wrapper.exe
Target: x86_64-w64-mingw32
Configured with: ../gcc-13.2.0/configure --prefix=/mingw64 --with-local-prefix=/mingw64/local --build=x86_64-w64-mingw32 --host=x86_64-w64-mingw32 --target=x86_64-w64-mingw32 --with-native-system-header-dir=/mingw64/include --libexecdir=/mingw64/lib --enable-bootstrap --enable-checking=release --with-arch=nocona --with-tune=generic --enable-languages=c,lto,c++,fortran,ada,objc,obj-c++,jit --enable-shared --enable-static --enable-libatomic --enable-threads=posix --enable-graphite --enable-fully-dynamic-string --enable-libstdcxx-filesystem-ts --enable-libstdcxx-time --disable-libstdcxx-pch --enable-lto --enable-libgomp --disable-libssp --disable-multilib --disable-rpath --disable-win32-registry --disable-nls --disable-werror --disable-symvers --with-libiconv --with-system-zlib --with-gmp=/mingw64 --with-mpfr=/mingw64 --with-mpc=/mingw64 --with-isl=/mingw64 --with-pkgversion='Rev3, Built by MSYS2 project' --with-bugurl=https://github.com/msys2/MINGW-packages/issues --with-gnu-as --with-gnu-ld --disable-libstdcxx-debug --with-boot-ldflags=-static-libstdc++ --with-stage1-ldflags=-static-libstdc++
Thread model: posix
Supported LTO compression algorithms: zlib zstd
gcc version 13.2.0 (Rev3, Built by MSYS2 project)
```

# Verificación de archivos
En la última actualización se hace una verificación de los archivos del lado del cliente para comprobar que no se han modificado. Si se quiere desactivar, hay que ir a las líneas 61 a 63 de main.go y removerlas.

Por otro lado, si esto se quiere usar hay que compilar el lado del cliente, calcular el sha512 de cada archivo, ir a checksum_archivos.go:29, y ajustar nombres de archivo junto con el checksum

# Preparando entorno de desarrollo
Se necesita contar con NPM y Go, además de GCC si se usa SQLite3 pues el mismo debe ser compilado.
Por otro lado, si quieres que el makefile funcione, debes contar con `make`


- NPM: https://parzibyte.me/blog/2018/09/27/instalar-npm-node-js-windows/
- Go: https://parzibyte.me/blog/2017/12/07/instalar-configurar-go-golang-en-windows-10/
- GCC: https://parzibyte.me/blog/2018/09/27/instalar-gcc-64-bits-en-windows-con-mingw/

Una vez que tengas el proyecto, abre una terminal en el mismo. Instala las dependencias de Node con:

`npm install`

Después de eso ya puedes iniciar el servidor de desarrollo con:
`npm run start`

Ahora vamos al lado del servidor así que abre otra terminal y navega a este directorio pero entra a la carpeta **api**.


Si tienes **make** ejecuta:

`make`

Y él se encargará de todo.

O en caso de que no tengas **make** instala las dependencias de Go con:
`go mod tidy`

Luego formatea el código (no es necesario, pero se ve bonito):
`gofmt -w .`

Finalmente compila con:
`go build -o tiendas.exe -tags desarrollo`


Y luego ejecuta `tiendas.exe`

# Compilando para producción

## Lado del servidor
Recuerda que debes estar dentro de **api**. Compila para producción con: `make prod`

O si no cuentas con make, compila así:
`go build -o tiendas_prod_64.exe -tags produccion -ldflags "-H windowsgui"`

Ahora copia los siguientes archivos a una carpeta limpia:
- tiendas_prod_64.exe
- esquema_negocios_sqlite.sql
- esquema_spos_sqlite.sql
- init.sql

Y dentro de esa misma carpeta crea otra carpeta vacía llamada **dist**

Hasta el momento el árbol debe verse así:

```
λ tree /F
Listado de rutas de carpetas
El número de serie del volumen es 
C:.
│   esquema_negocios_sqlite.sql
│   esquema_spos_sqlite.sql
│   init.sql
│   tiendas_prod_64.exe
│
└───dist
```

## Lado del cliente
Recuerda que NO debes estar dentro de **api**.

1. Compila los archivos con: `npm run build`
2. Dentro de `dist` (carpeta recién creada por el script, no te confundas con la que creaste anteriormente) mover la carpeta `fonts` a `css` de manera que su ruta ahora sea `css/fonts`
3. Entra a las carpetas **js** y **css**, busca los archivos que terminen en **map** y elimínalos
4. Copia todo lo de dist (TODO LO DE dist, no la carpeta en sí, sino su contenido) a la carpeta **dist** que creaste anteriormente

## Verificando
Ahora el directorio se debería ver así:

```
λ tree /F
Listado de rutas de carpetas
El número de serie del volumen es 
C:.
│   esquema_negocios_sqlite.sql
│   esquema_spos_sqlite.sql
│   init.sql
│   tiendas_prod_64.exe
│
└───dist
    │   index.html
    │
    ├───css
    │   │   app.9afe257e7d77a264e0675dacea965a46.css
    │   │
    │   └───fonts
    │           material-icons.ac188f9.woff2
    │           roboto-v18-latin-300.55536c8.woff2
    │           roboto-v18-latin-300.a1471d1.woff
    │           roboto-v18-latin-500.2854671.woff2
    │           roboto-v18-latin-500.de8b743.woff
    │           roboto-v18-latin-700.037d830.woff2
    │           roboto-v18-latin-700.cf6613d.woff
    │           roboto-v18-latin-regular.5d4aeb4.woff2
    │           roboto-v18-latin-regular.bafb105.woff
    │
    ├───img
    │       ajustes.06556a3.png
    │       bolsa-de-la-compra.8aefc13.png
    │       dinero.280fff4.png
    │       escudo.2db8426.png
    │       estadistica.d676772.png
    │       monedas.cc6526e.png
    │       que-hacer.c1d47a0.png
    │       red.524a30a.png
    │       salpicadero.f66d1ab.png
    │       tareas.ecd565c.png
    │
    └───js
            app.6a416a79713a4b7ab346.js
            manifest.5ebdb86dbeef19313db4.js
            vendor.01fdaafe859388bfc024.js
```

Ahora solo resta distribuir la app. Cuando lo hagas, solo hay que ejecutar:
`tiendas_prod_64.exe`

Y luego visitar http://localhost:2106/static/

        



# Usar MySQL en lugar de SQLite3
Este sistema es compatible con MySQL, pero por ahora solo se usa pensando en SQLite3. Si quieres cambiar a MySQL aquí están los apuntes

## Funciones de fecha
Si se usa SQLite3 se debe usar la función STRFTIME en lugar de DATE_FORMAT usada en MySQL.
Recuerda que hay diferencias, por ejemplo, para extraer el mes con SQLite3 es:


`strftime("%m", fecha)`

Mientras que con MySQL es:

`date_format(fecha, "%m")`

Por ahora se ha encontrado el uso de la función en:

- DatosGraficasController.go

## Archivos
En lugar de usar los archivos que terminan en `sqlite.sql` deberían usarse solo los de `.sql`, haciendo los cambios también en `Constantes.go`

## CREATE DATABASE
Con SQLite3 no es necesario el `create database...` pero con MySQL sí. Si vas a usar MySQL, además de los cambios de arriba, debes cambiar el if en `NegociosController.go:380`

## Permisos
Para que todo vaya como un encanto, el usuario necesita permisos de:
* CREATE
* REFERENCES
* INDEX
* SELECT
* UPDATE
* DELETE
* INSERT
* ALTER (esto es para que pueda reiniciar el autoincremento)
* DROP (para eliminar las de los negocios)

Lo que se logra con:
```mysql
grant alter, drop, create, references, index, select, update, delete, insert 
on *.* to 'usuario'@'localhost';

```

## Build Setup

``` bash
# install dependencies
npm install

# serve with hot reload at localhost:8080
npm run dev

# build for production with minification
npm run build

# build for production and view the bundle analyzer report
npm run build --report
```

For a detailed explanation on how things work, check out the [guide](http://vuejs-templates.github.io/webpack/) and [docs for vue-loader](http://vuejs.github.io/vue-loader).
