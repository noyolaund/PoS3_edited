package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
)

const LimiteBusquedaProductosAutocompletado = 10

type ProductosController struct {
	AjustesUsuario AjustesDeUsuarioLogueado
}

func (p *ProductosController) importarExcel(archivo multipart.File, informacion *multipart.FileHeader, ajustesDeImportacion *AjustesParaImportarProductos) bool {
	inicio := time.Now()
	//Copiar el archivo ajustesDeImportacion un lugar en donde podamos examinarlo y abrirlo
	s := Sesion{
		request:        p.AjustesUsuario.httpRequest,
		responseWriter: p.AjustesUsuario.httpResponseWriter,
	}
	usuario, err := s.obtenerUsuarioLogueado()
	if err != nil {
		//TODO: return err
		return false
	}
	idNegocio := strconv.Itoa(usuario.Negocio.Id)
	ubicacion, err := p.importarArchivo(archivo, informacion, idNegocio)
	if err != nil {
		log.Printf("Error al importar archivo: %v", err.Error())
		return false
	}
	archivoExcel, err := xlsx.OpenFile(ubicacion)
	if err != nil {
		log.Printf("Error abriendo archivo excel en la ruta %s: %v", ubicacion, err)
		return false
	}

	if len(archivoExcel.Sheets) > 0 {
		hoja := archivoExcel.Sheets[0]
		comenzarDesde := 0 //La fila en donde se comienza
		if ajustesDeImportacion.TieneEncabezados {
			//Si tiene encabezados se omite la primera
			comenzarDesde = 1
		}
		cantidadDeFilas := len(hoja.Rows)
		var codigoBarras, descripcion string
		var precioCompra, precioVenta, existencia, stock float64
		var celdas []*xlsx.Cell

		db, err := p.AjustesUsuario.obtenerBaseDeDatos()

		if err != nil {
			log.Printf("Error abriendo base de datos al importar productos: %q", err)
			panic(err)
		}

		// Asegurar que las columnas idPadre y equivalencia existen
		ensureProductosColumns(db)

		defer db.Close()

		tx, err := db.Begin()
		defer tx.Commit()

		if err != nil {
			log.Printf("Error obteniendo contexto de la base de datos al importar productos:\n%q", err)
			return false
		}

		preferenciaConsulta := "replace" //Por defecto se remplazan los repetidos
		if ajustesDeImportacion.IgnorarCodigosDeBarrasRepetidos {
			preferenciaConsulta = "insert ignore" // O se ignoran si el usuario así lo quiere
		}
		sentenciaPreparada, err := tx.Prepare(preferenciaConsulta + ` INTO productos
(codigoBarras, descripcion, precioCompra, precioVenta, existencia, stock)
VALUES
(?, ?, ?, ?, ?, ?)`)
		if err != nil {
			log.Printf("Error preparando consulta al importar productos:\n%q", err)
			return false
		}

		for indiceFila := comenzarDesde; indiceFila < cantidadDeFilas; indiceFila++ {
			celdas = hoja.Rows[indiceFila].Cells
			cantidadCeldas := len(celdas)
			if ajustesDeImportacion.IndiceCodigoBarras > 0 && ajustesDeImportacion.IndiceCodigoBarras <= cantidadCeldas {
				codigoBarras = celdas[ajustesDeImportacion.IndiceCodigoBarras-1].String()
			} else {
				log.Printf("Error. El índice del código de barras debe ser mayor o igual que 0. En cambio, se intentó dar %d", ajustesDeImportacion.IndiceCodigoBarras)
				return false
			}
			if ajustesDeImportacion.IndiceDescripcion > 0 && ajustesDeImportacion.IndiceDescripcion <= cantidadCeldas {
				descripcion = celdas[ajustesDeImportacion.IndiceDescripcion-1].String()
			} else {
				log.Printf("Error. El índice de la descripción debe ser mayor o igual que 0. En cambio, se intentó dar %d", ajustesDeImportacion.IndiceDescripcion)
				return false
			}
			if ajustesDeImportacion.IndicePrecioCompra > 0 && ajustesDeImportacion.IndicePrecioCompra <= cantidadCeldas {
				precioCompra, _ = celdas[ajustesDeImportacion.IndicePrecioCompra-1].Float()
			}
			if ajustesDeImportacion.IndicePrecioVenta > 0 && ajustesDeImportacion.IndicePrecioVenta <= cantidadCeldas {
				precioVenta, _ = celdas[ajustesDeImportacion.IndicePrecioVenta-1].Float()
			}
			if ajustesDeImportacion.IndiceExistencia > 0 && ajustesDeImportacion.IndiceExistencia <= cantidadCeldas {
				existencia, _ = celdas[ajustesDeImportacion.IndiceExistencia-1].Float()
			}
			if ajustesDeImportacion.IndiceStock > 0 && ajustesDeImportacion.IndiceStock <= cantidadCeldas {
				stock, _ = celdas[ajustesDeImportacion.IndiceStock-1].Float()
			}
			_, err = sentenciaPreparada.Exec(codigoBarras, descripcion, precioCompra, precioVenta, existencia, stock)
			if err != nil {
				log.Printf("Error ejecutando sentencia:\n%q", err)
			}
		}
	}

	//Finalmente eliminar el archivo, pues no lo necesitamos
	err = os.Remove(ubicacion)
	if err != nil {
		log.Printf("Error eliminando el archivo después de importar: %v", err)
	}
	transcurrido := time.Since(inicio)
	log.Printf("El tiempo que se llevó importar los productos desde un archivo Excel: %s", transcurrido)

	return true
}
func (p *ProductosController) importarArchivo(archivo multipart.File, informacion *multipart.FileHeader, idNegocio string) (rutaGuardado string, err error) {
	defer archivo.Close()
	nombreArchivo := "temp_" + idNegocio + ".xlsx"
	directorioSubidas := path.Join(obtenerRutaActual(), NombreDirectorioParaSubidas)
	crearDirectorioSiNoExiste(directorioSubidas)
	ubicacionDelNuevoArchivo := path.Join(directorioSubidas, nombreArchivo)
	nuevoArchivo, err := os.OpenFile(ubicacionDelNuevoArchivo, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Printf("Error creando archivo para importar Excel. El error: %v", err.Error())
		return "", err
	}
	defer nuevoArchivo.Close()

	io.Copy(nuevoArchivo, archivo)
	return ubicacionDelNuevoArchivo, nil
}
func (p *ProductosController) exportar(ajustes AjustesParaExportarProductos) {
	switch ajustes.Extension {
	case "xlsx":
		p.exportarExcel(&ajustes)
		break
	case "csv":
		p.exportarCSV(&ajustes)
		break
	default:
		log.Printf("No se puede exportar porque la extensión %v es inválida", ajustes.Extension)
		break
	}
}

func (p *ProductosController) exportarCSV(ajustes *AjustesParaExportarProductos) {
	inicio := time.Now()
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}
	// Asegurar columnas necesarias
	ensureProductosColumns(db)

	defer db.Close()
	rutaArchivo := obtenerRutaActual() + "/" + NombreArchivoExportadoCSV
	os.Remove(rutaArchivo)
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, precioCompra, precioVenta, existencia, 
stock, idPadre, equivalencia FROM productos ORDER BY idProducto DESC;`)
	if err != nil {
		log.Printf("Error al realizar la consulta para exportar productos a CSV:\n%q", err)
		return
	}

	defer filas.Close()

	archivo, err := os.Create(rutaArchivo)
	comprobarError(err)
	defer archivo.Close()

	if ajustes.IncluirEncabezado {
		archivo.WriteString("idProducto,codigoBarras,descripcion,precioCompra,precioVenta,existencia,stock\n")
	}
	var producto Producto
	for filas.Next() {

		err := filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion, &producto.PrecioCompra,
			&producto.PrecioVenta, &producto.Existencia, &producto.Stock, &producto.Padre, &producto.Equivalencia)
		for contador := 0; contador < ajustes.Copias; contador++ {
			archivo.WriteString(fmt.Sprintf("%d,%s,%s,%0.2f,%0.2f,%0.2f,%0.2f\n",
				producto.Numero, producto.CodigoBarras, producto.Descripcion, producto.PrecioCompra,
				producto.PrecioVenta, producto.Existencia, producto.Stock))
		}
		if err != nil {
			log.Printf("Error al escanear un producto para exportar:\n%q", err)
		}
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todos los productos:\n%q", err)
	}
	transcurrido := time.Since(inicio)
	log.Printf("El tiempo que se llevó exportar todos los productos a un archivo CSV: %s", transcurrido)

}
func (p *ProductosController) exportarExcel(ajustes *AjustesParaExportarProductos) {
	inicio := time.Now()
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}
	// Asegurar columnas necesarias
	ensureProductosColumns(db)

	defer db.Close()

	rutaArchivo := obtenerRutaActual() + "/" + NombreArchivoExportadoExcel
	os.Remove(rutaArchivo)

	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, precioCompra, precioVenta, existencia, 
stock, idPadre, equivalencia FROM productos ORDER BY idProducto DESC;`)
	if err != nil {
		log.Printf("Error al realizar la consulta para exportar productos a Excel:\n%q", err)
		return
	}

	defer filas.Close()

	archivo := xlsx.NewFile()
	hoja, err := archivo.AddSheet("Hoja 1")
	comprobarError(err)

	var fila *xlsx.Row
	var celda *xlsx.Cell
	if ajustes.IncluirEncabezado {
		fila = hoja.AddRow()
		for _, encabezado := range []string{"idProducto", "codigoBarras", "descripcion", "precioCompra", "precioVenta", "existencia", "stock"} {
			celda = fila.AddCell()
			celda.Value = encabezado
		}
	}
	var producto Producto
	for filas.Next() {

		err := filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion, &producto.PrecioCompra,
			&producto.PrecioVenta, &producto.Existencia, &producto.Stock, &producto.Padre, &producto.Equivalencia)
		for contador := 0; contador < ajustes.Copias; contador++ {
			fila = hoja.AddRow()
			celda = fila.AddCell()
			celda.SetInt(producto.Numero)
			celda = fila.AddCell()
			celda.SetString(producto.CodigoBarras)
			celda = fila.AddCell()
			celda.SetString(producto.Descripcion)
			celda = fila.AddCell()
			celda.SetFloat(producto.PrecioCompra)
			celda = fila.AddCell()
			celda.SetFloat(producto.PrecioVenta)
			celda = fila.AddCell()
			celda.SetFloat(producto.Existencia)
			celda = fila.AddCell()
			celda.SetFloat(producto.Stock)
		}
		if err != nil {
			log.Printf("Error al escanear un producto para exportar:\n%q", err)
		}
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todos los productos:\n%q", err)
	}
	archivo.Save(rutaArchivo)
	transcurrido := time.Since(inicio)
	log.Printf("El tiempo que se llevó exportar todos los productos a un archivo Excel: %s", transcurrido)

}

func (p *ProductosController) reporteInventario() *ReporteInventario {
	reporte := ReporteInventario{}
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}
	defer db.Close()
	filas, err := db.Query(`SELECT SUM(precioCompra * existencia) AS precioCompra, 
SUM(precioVenta * existencia) AS precioVenta, 
COUNT(*) AS cantidadProductos FROM productos;`)
	if err != nil {
		log.Printf("Error al consultar reporte de inventario:\n%v", err)
		return &reporte
	}

	defer filas.Close()

	if !filas.Next() {
		log.Printf("No hay filas que escanear")
		return &reporte
	}
	err = filas.Scan(&reporte.PrecioCompra, &reporte.PrecioVenta, &reporte.CantidadProductos)
	if err != nil {
		log.Printf("Error al escanear reporte de inventario:\n%v", err)
		return &reporte
	}
	return &reporte
}

func (p *ProductosController) alAzar(cuantos int) *[]Producto {

	cantidadDeProductosExistentes, err := p.conteo()
	if err != nil {
		log.Fatalf("Error contando productos existentes: %v", err)
	}
	//Arreglar, ya que Intn no acepta un 0
	if cantidadDeProductosExistentes <= 0 {
		cantidadDeProductosExistentes = 1
	}
	fuente := rand.NewSource(time.Now().Unix())
	random := rand.New(fuente)
	offset := random.Intn(cantidadDeProductosExistentes)

	productos := p.todos(offset, cuantos)
	return &productos.Productos
}

func (p *ProductosController) deUnApartado(idApartado int) []ProductoApartado {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, precioCompra, precioVenta, 
precioVentaOriginal, cantidadVendida FROM productos_apartados WHERE idApartado = ?;`, idApartado)
	if err != nil {
		log.Printf("Error consultando productos de un apartado\n%q", err)
		return nil
	}

	defer filas.Close()

	productos := []ProductoApartado{}
	for filas.Next() {
		var producto ProductoApartado
		err := filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion, &producto.PrecioCompra,
			&producto.PrecioVenta, &producto.PrecioVentaOriginal, &producto.Cantidad)
		if err != nil {
			log.Printf("Error consultando productos de un apartado\n%q", err)
		}
		productos = append(productos, producto)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas:\n%q", err)
	}
	return productos
}
func (p *ProductosController) deUnaVenta(idVenta int) []ProductoVendido {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, precioCompra, precioVenta, precioVentaOriginal, 
cantidadVendida FROM productos_vendidos WHERE idVenta = ?;`, idVenta)
	if err != nil {
		log.Printf("Error consultando productos de una venta\n%q", err)
		return nil
	}

	defer filas.Close()

	productos := []ProductoVendido{}
	for filas.Next() {
		var producto ProductoVendido
		err := filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion, &producto.PrecioCompra,
			&producto.PrecioVenta, &producto.PrecioVentaOriginal, &producto.Cantidad)
		if err != nil {
			log.Printf("Error consultando productos de una venta\n%q", err)
		}
		productos = append(productos, producto)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas:\n%q", err)
	}
	return productos
}
func (p *ProductosController) buscarParaAutocompletado(busqueda string) []Producto {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}
	ensureProductosColumns(db)

	defer db.Close()
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, precioCompra, precioVenta, existencia, 
stock, idPadre, equivalencia FROM productos WHERE descripcion LIKE ? ORDER BY idProducto DESC LIMIT ?;`,
		"%"+busqueda+"%", LimiteBusquedaProductosAutocompletado)
	if err != nil {
		log.Printf("Error en consulta al realizar una búsqueda en productos para autocompletado:\n%q", err)
		return nil
	}

	defer filas.Close()

	productos := []Producto{}
	for filas.Next() {
		var producto Producto
		err := filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion, &producto.PrecioCompra,
			&producto.PrecioVenta, &producto.Existencia, &producto.Stock, &producto.Padre, &producto.Equivalencia)
		if err != nil {
			log.Printf("Error al escanear resultados de búsqueda en productos para autocompletado:\n%q", err)
		}
		productos = append(productos, producto)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas:\n%q", err)
	}
	return productos
}
func (p *ProductosController) nuevo(producto *Producto) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "productos",
		AjustesUsuario: p.AjustesUsuario,
	}
	existencia := producto.Existencia
	if producto.Padre > 0 {
		// Los hijos no almacenan existencia propia; la maneja el padre
		existencia = 0
	}
	ayudante.insertar(map[string]interface{}{
		"descripcion":  producto.Descripcion,
		"codigoBarras": producto.CodigoBarras,
		"precioCompra": producto.PrecioCompra,
		"precioVenta":  producto.PrecioVenta,
		"existencia":   existencia,
		"stock":        producto.Stock,
		"idPadre":      producto.Padre,
		"equivalencia": producto.Equivalencia,
	})
}
func (p *ProductosController) actualizar(producto *Producto) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "productos",
		AjustesUsuario: p.AjustesUsuario,
	}
	campos := map[string]interface{}{
		"descripcion":  producto.Descripcion,
		"codigoBarras": producto.CodigoBarras,
		"precioCompra": producto.PrecioCompra,
		"precioVenta":  producto.PrecioVenta,
		"stock":        producto.Stock,
		"idPadre":      producto.Padre,
		"equivalencia": producto.Equivalencia,
	}
	// No permitir actualizar existencia de un producto hijo directamente
	if producto.Padre == 0 {
		campos["existencia"] = producto.Existencia
	}
	ayudante.actualizarDonde("idProducto", producto.Numero, campos)
}

// reabastecer incrementa existencia. Si se pasa el id de un producto hijo,
// convierte la cantidad en unidades padre y actualiza la existencia del padre.
func (p *ProductosController) reabastecer(idProducto int, cantidad float64) error {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Obtener relación padre/equivalencia
	var idPadre int
	var equivalencia int
	fila := tx.QueryRow("SELECT idPadre, equivalencia FROM productos WHERE idProducto = ?;", idProducto)
	if err := fila.Scan(&idPadre, &equivalencia); err != nil {
		// Si falla, aplicar al propio producto
		_, err = tx.Exec("UPDATE productos SET existencia = existencia + ? WHERE idProducto = ?;", cantidad, idProducto)
		if err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	}

	if idPadre > 0 && equivalencia > 0 {
		delta := RoundToTwoDecimals(cantidad / float64(equivalencia))
		_, err = tx.Exec("UPDATE productos SET existencia = existencia + ? WHERE idProducto = ?;", delta, idPadre)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err = tx.Exec("UPDATE productos SET existencia = existencia + ? WHERE idProducto = ?;", cantidad, idProducto)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
func (p *ProductosController) eliminar(producto *Producto) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "productos",
		AjustesUsuario: p.AjustesUsuario,
	}
	ayudante.eliminarDonde("idProducto", producto.Numero)
}
func (p *ProductosController) todos(offset, limite int) ProductosConConteo {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()
	pc := ProductosConConteo{}
	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}
	ensureProductosColumns(db)

	defer db.Close()
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, precioCompra, precioVenta, existencia, 
stock, idPadre, equivalencia FROM productos ORDER BY idProducto DESC LIMIT ? OFFSET ?;`, limite, offset)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener todos los productos:\n%q", err)
		return pc
	}

	defer filas.Close()

	productos := []Producto{}
	for filas.Next() {
		var producto Producto
		err := filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion, &producto.PrecioCompra,
			&producto.PrecioVenta, &producto.Existencia, &producto.Stock, &producto.Padre, &producto.Equivalencia)
		if err != nil {
			log.Printf("Error al escanear todos los productos:\n%q", err)
		}
		productos = append(productos, producto)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todos los productos:\n%q", err)
	}
	conteo, _ := p.conteo()
	pc.Conteo = conteo
	pc.Productos = productos
	return pc
}
func (p *ProductosController) enStock(offset, limite int) ProductosConConteo {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()
	pc := ProductosConConteo{}
	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, existencia, stock, idPadre, equivalencia FROM productos 
WHERE existencia < stock ORDER BY idProducto DESC LIMIT ? OFFSET ?;`, limite, offset)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener productos en stock:\n%q", err)
		return pc
	}

	defer filas.Close()

	productos := []Producto{}
	for filas.Next() {
		var producto Producto
		err := filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion,
			&producto.Existencia, &producto.Stock, &producto.Padre, &producto.Equivalencia)
		if err != nil {
			log.Printf("Error al escanear productos en stock:\n%q", err)
		}
		productos = append(productos, producto)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todos los productos:\n%q", err)
	}
	pc.Conteo = p.conteoParaStock()
	pc.Productos = productos
	return pc
}
func (p *ProductosController) buscar(offset, limite int, descripcion string) ProductosConConteo {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()
	pc := ProductosConConteo{}
	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, precioCompra, 
precioVenta, existencia, stock, idPadre, equivalencia FROM productos WHERE descripcion LIKE ? ORDER BY idProducto DESC LIMIT ? OFFSET ?;`,
		"%"+descripcion+"%", limite, offset)
	if err != nil {
		log.Printf("Error en consulta al realizar una búsqueda en productos:\n%q", err)
		return pc
	}

	defer filas.Close()

	productos := []Producto{}
	for filas.Next() {
		var producto Producto
		err := filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion,
			&producto.PrecioCompra, &producto.PrecioVenta, &producto.Existencia, &producto.Stock, &producto.Padre, &producto.Equivalencia)
		if err != nil {
			log.Printf("Error al escanear resultados de búsqueda en productos:\n%q", err)
		}
		productos = append(productos, producto)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas:\n%q", err)
	}
	pc.Productos = productos
	conteo, _ := p.conteoBusqueda(descripcion)
	pc.Conteo = conteo
	return pc
}

func (p *ProductosController) siguienteNumero() (siguienteNumero int, err error) {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT idProducto FROM productos ORDER BY idProducto DESC LIMIT 1;")
	if err != nil {
		log.Printf("Error al consultar el siguiente número de producto:\n%q", err)
		return -1, err
	}

	defer filas.Close()

	if !filas.Next() {
		return 1, nil
	}
	var ultimoRowid int
	err = filas.Scan(&ultimoRowid)
	if err != nil {
		log.Printf("Error al escanear idProducto de productos:\n%q", err)
		return -1, err
	}
	return ultimoRowid + 1, nil
}

func (p *ProductosController) conteo() (conteo int, err error) {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT COUNT(idProducto) FROM productos;")
	if err != nil {
		log.Printf("Error al conteo productos:\n%q", err)
		return 0, err
	}

	defer filas.Close()

	if !filas.Next() {
		return 0, nil
	}
	var total int
	err = filas.Scan(&total)
	if err != nil {
		log.Printf("Error al conteo productos:\n%q", err)
		return 0, err
	}
	return total, nil
}
func (p *ProductosController) conteoParaStock() int {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT COUNT(*) AS conteo FROM productos WHERE existencia < stock;")
	if err != nil {
		log.Printf("Error al conteo productos con existencia menor que el stock:\n%q", err)
		return 0
	}

	defer filas.Close()

	if !filas.Next() {
		return 0
	}
	var total int
	err = filas.Scan(&total)
	if err != nil {
		log.Printf("Error al conteo productos con existencia menor que el stock:\n%q", err)
		return 0
	}
	return total
}
func (p *ProductosController) conteoBusqueda(descripcion string) (conteo int, err error) {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT COUNT(idProducto) FROM productos WHERE descripcion LIKE ?;",
		"%"+descripcion+"%")
	if err != nil {
		log.Printf("Error al conteo productos de búsqueda:\n%q", err)
		return 0, err
	}

	defer filas.Close()

	if !filas.Next() {
		return 0, nil
	}
	var total int
	err = filas.Scan(&total)
	if err != nil {
		log.Printf("Error al conteo productos de búsqueda:\n%q", err)
		return 0, err
	}
	return total, nil
}
func (p *ProductosController) unoApartado(idApartado, idProducto int) *ProductoApartado {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}
	defer db.Close()
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, precioCompra, precioVenta, 
precioVentaOriginal,cantidadVendida FROM productos_apartados WHERE idProducto = ? AND idApartado = ? LIMIT 1;`,
		idProducto, idApartado)
	if err != nil {
		log.Printf("Error al consultar producto de un apartado por idProducto:\n%q", err)
		return nil
	}

	defer filas.Close()

	if !filas.Next() {
		return nil
	}
	var producto ProductoApartado
	err = filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion, &producto.PrecioCompra,
		&producto.PrecioVenta, &producto.PrecioVentaOriginal, &producto.Cantidad)
	if err != nil {
		log.Printf("Error al escanear producto de un apartado por idProducto:\n%q", err)
		return nil
	}
	return &producto
}
func (p *ProductosController) porRowid(idProducto int) *Producto {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}
	// Asegurar que las columnas idPadre y equivalencia existen
	ensureProductosColumns(db)

	defer db.Close()
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, 
precioCompra, precioVenta, existencia, stock  FROM productos WHERE idProducto = ? LIMIT 1;`, idProducto)
	if err != nil {
		log.Printf("Error al consultar producto por idProducto:\n%q", err)
		return nil
	}

	defer filas.Close()

	if !filas.Next() {
		return nil
	}
	var producto Producto
	err = filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion,
		&producto.PrecioCompra, &producto.PrecioVenta, &producto.Existencia, &producto.Stock)
	if err != nil {
		log.Printf("Error al escanear producto por idProducto:\n%q", err)
		return nil
	}
	// Si el producto es hijo (tiene idPadre), calcular existencia y stock equivalentes
	filaPadre := db.QueryRow(`SELECT idPadre, equivalencia FROM productos WHERE idProducto = ? LIMIT 1;`, idProducto)
	var idPadre int
	var equivalencia int
	if err := filaPadre.Scan(&idPadre, &equivalencia); err == nil {
		if idPadre > 0 && equivalencia > 0 {
			fila := db.QueryRow(`SELECT existencia, stock FROM productos WHERE idProducto = ? LIMIT 1;`, idPadre)
			var existenciaPadre, stockPadre float64
			if err := fila.Scan(&existenciaPadre, &stockPadre); err == nil {
				producto.Existencia = RoundToTwoDecimals(existenciaPadre * float64(equivalencia))
				producto.Stock = RoundToTwoDecimals(stockPadre * float64(equivalencia))
			}
		}
	}
	return &producto
}
func (p *ProductosController) porCodigoDeBarras(codigoDeBarras string) *Producto {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}
	// Asegurar que las columnas idPadre y equivalencia existen
	ensureProductosColumns(db)

	defer db.Close()
	filas, err := db.Query(`SELECT idProducto, codigoBarras, descripcion, precioCompra, 
precioVenta, existencia, stock  FROM productos WHERE codigoBarras = ? LIMIT 1;`, codigoDeBarras)
	if err != nil {
		log.Printf("Error al consultar producto por código de barras:\n%q", err)
		return nil
	}

	defer filas.Close()

	if !filas.Next() {
		return nil
	}
	var producto Producto
	err = filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion,
		&producto.PrecioCompra, &producto.PrecioVenta, &producto.Existencia, &producto.Stock)
	if err != nil {
		log.Printf("Error al escanear producto por código de barras:\n%q", err)
		return nil
	}
	// Si el producto es hijo (tiene idPadre), calcular existencia y stock equivalentes
	filaPadre := db.QueryRow(`SELECT idPadre, equivalencia FROM productos WHERE idProducto = ? LIMIT 1;`, producto.Numero)
	var idPadre int
	var equivalencia int
	if err := filaPadre.Scan(&idPadre, &equivalencia); err == nil {
		if idPadre > 0 && equivalencia > 0 {
			fila := db.QueryRow(`SELECT existencia, stock FROM productos WHERE idProducto = ? LIMIT 1;`, idPadre)
			var existenciaPadre, stockPadre float64
			if err := fila.Scan(&existenciaPadre, &stockPadre); err == nil {
				producto.Existencia = RoundToTwoDecimals(existenciaPadre * float64(equivalencia))
				producto.Stock = RoundToTwoDecimals(stockPadre * float64(equivalencia))
			}
		}
	}
	return &producto
}

// ensureProductosColumns agrega las columnas idPadre y equivalencia a la tabla productos
// si no existen. Esto evita errores en instalaciones antiguas que no tienen el nuevo esquema.
func ensureProductosColumns(db *sql.DB) {
	filas, err := db.Query("PRAGMA table_info(productos);")
	if err != nil {
		return
	}
	defer filas.Close()

	hasIdPadre := false
	hasEquivalencia := false
	for filas.Next() {
		var cid int
		var name string
		var ctype string
		var notnull int
		var dflt sql.NullString
		var pk int
		if err := filas.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			continue
		}
		if name == "idPadre" {
			hasIdPadre = true
		}
		if name == "equivalencia" {
			hasEquivalencia = true
		}
	}

	if !hasIdPadre {
		_, _ = db.Exec("ALTER TABLE productos ADD COLUMN idPadre INTEGER NOT NULL DEFAULT 0;")
	}
	if !hasEquivalencia {
		_, _ = db.Exec("ALTER TABLE productos ADD COLUMN equivalencia INTEGER NOT NULL DEFAULT 1;")
	}
}
