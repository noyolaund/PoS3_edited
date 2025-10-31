package main

import (
	"fmt"

	"log"
)

type DatosGraficasController struct {
	AjustesUsuario AjustesDeUsuarioLogueado
}

func (d *DatosGraficasController) mesesEnLosQueHayRegistrosDeVentasAlContadoDependiendoDeAnio(anio int) []string {
	db, err := d.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT DISTINCT STRFTIME("%m", fecha) AS mes FROM ventas_contado WHERE fecha BETWEEN ? AND ?;`, anio, anio+1)
	if err != nil {
		log.Printf("Error al consultar meses en los que hay datos para ventas al contado, dependiendo de un año:\n%q", err)
	}

	defer filas.Close()

	var meses []string

	for filas.Next() {

		var mes string
		err = filas.Scan(&mes)
		if err != nil {
			log.Printf("Error al escanear meses en los que hay datos para ventas al contado, dependiendo de un año:\n%q", err)
		}
		meses = append(meses, mes)
	}
	return meses
}
func (d *DatosGraficasController) aniosEnLosQueHayRegistrosDeVentasAlContado() []string {
	db, err := d.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT DISTINCT STRFTIME("%Y", fecha) AS anio FROM ventas_contado ORDER BY anio DESC;`)
	if err != nil {
		log.Printf("Error al consultar años en los que hay registros de ventas al contado:\n%q", err)
	}

	defer filas.Close()

	anios := []string{}

	for filas.Next() {

		var anio string
		err = filas.Scan(&anio)
		if err != nil {
			log.Printf("Error al escanear años en los que hay registros de ventas al contado:\n%q", err)
		}
		anios = append(anios, anio)
	}
	return anios
}

func (d *DatosGraficasController) productosMasVendidos(fechaInicio, fechaFin string) []ProductoVendidoParaGrafica {
	db, err := d.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT SUM(productos_vendidos.cantidadVendida) AS vecesVendido,
								productos_vendidos.idProducto,
								productos_vendidos.codigoBarras,
								productos_vendidos.descripcion
							FROM productos_vendidos
							INNER JOIN ventas_contado
							ON productos_vendidos.idVenta = ventas_contado.idVenta
							AND ventas_contado.fecha
							BETWEEN ?
							AND ?
							GROUP BY  productos_vendidos.idProducto,codigoBarras, descripcion
							ORDER BY  vecesVendido DESC LIMIT ?;`, fechaInicio, fechaFin, LimiteProductosMasVendidos)
	if err != nil {
		log.Printf("Error al consultar productos más vendidos:\n%q", err)
	}

	defer filas.Close()

	productos := []ProductoVendidoParaGrafica{}
	for filas.Next() {

		var producto ProductoVendidoParaGrafica
		err = filas.Scan(&producto.VecesVendido, &producto.IdProducto, &producto.CodigoBarras, &producto.Descripcion)
		if err != nil {
			log.Printf("Error al escanear producto más vendido:\n%q", err)
		}
		productos = append(productos, producto)
	}
	return productos
}

func (d *DatosGraficasController) productosMenosVendidos(fechaInicio, fechaFin string) []ProductoVendidoParaGrafica {
	db, err := d.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT SUM(productos_vendidos.cantidadVendida) AS vecesVendido,
								productos_vendidos.idProducto,
								productos_vendidos.codigoBarras,
								productos_vendidos.descripcion
							FROM productos_vendidos
							INNER JOIN ventas_contado
							ON productos_vendidos.idVenta = ventas_contado.idVenta
							AND ventas_contado.fecha
							BETWEEN ?
							AND ?
							GROUP BY  productos_vendidos.idProducto
							ORDER BY  vecesVendido ASC LIMIT ?;`, fechaInicio, fechaFin, LimiteProductosMenosVendidos)
	if err != nil {
		log.Printf("Error al consultar productos más vendidos:\n%q", err)
	}

	defer filas.Close()

	productos := []ProductoVendidoParaGrafica{}
	for filas.Next() {

		var producto ProductoVendidoParaGrafica
		err = filas.Scan(&producto.VecesVendido, &producto.IdProducto, &producto.CodigoBarras, &producto.Descripcion)
		if err != nil {
			log.Printf("Error al escanear producto más vendido:\n%q", err)
		}
		productos = append(productos, producto)
	}
	return productos
}

func (d *DatosGraficasController) productosNuncaVendidosAlContado() []Producto {
	//TODO: tal vez poner cuáles no se venden en determinada fecha?
	db, err := d.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT productos.idProducto, productos.codigoBarras, productos.descripcion, 
productos.precioCompra, productos.precioVenta, productos.existencia, productos.stock 
FROM productos 
LEFT OUTER JOIN productos_vendidos 
ON productos.idProducto = productos_vendidos.idProducto 
WHERE productos_vendidos.idProducto IS NULL LIMIT ?;`, LimiteProductosNuncaVendidos)
	if err != nil {
		log.Printf("Error al consultar productos nunca vendidos:\n%q", err)
	}

	defer filas.Close()

	productos := []Producto{}
	for filas.Next() {

		var producto Producto
		err = filas.Scan(&producto.Numero, &producto.CodigoBarras, &producto.Descripcion, &producto.PrecioCompra, &producto.PrecioVenta, &producto.Existencia, &producto.Stock)
		if err != nil {
			log.Printf("Error al escanear  productos nunca vendidos:\n%q", err)
		}
		productos = append(productos, producto)
	}
	return productos
}

func (d *DatosGraficasController) totalVentasPorAnio(anio string) []EtiquetaYValor {

	/*
		Año como cadena, por ejemplo 2018
	*/

	db, err := d.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT STRFTIME("%m", fecha) AS mes, COALESCE(SUM(monto), 0) AS total FROM ventas_contado WHERE STRFTIME("%Y", fecha) = ? GROUP BY mes;`, anio)
	if err != nil {
		log.Printf("Error al consultar total vendido por mes en un año:\n%q", err)
	}

	defer filas.Close()

	meses := []EtiquetaYValor{}
	for filas.Next() {

		var ev EtiquetaYValor
		err = filas.Scan(&ev.Etiqueta, &ev.Valor)
		if err != nil {
			log.Printf("Error al escanear total vendido por mes en un año:\n%q", err)
		}
		meses = append(meses, ev)
	}
	return meses
}
func (d *DatosGraficasController) totalVentasPorDiaEnMesYAnio(mes, anio string) []EtiquetaYValor {

	/*
		Año como cadena, por ejemplo 2018
		Mes como cadena, en donde enero es 1 y diciembre 12. Se pasa con el cero si es que < 10
		Por ejemplo enero es 01 y diciembre 12
	*/
	mesYAnio := fmt.Sprintf("%s-%s", anio, mes)

	db, err := d.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT STRFTIME("%d", fecha) AS dia, COALESCE(SUM(monto), 0) AS total FROM ventas_contado WHERE STRFTIME("%Y-%m", fecha) = ? GROUP BY dia;`, mesYAnio)
	if err != nil {
		log.Printf("Error al consultar totalVendidoPorDia:\n%q", err)
	}

	defer filas.Close()

	dias := []EtiquetaYValor{}
	for filas.Next() {

		var ev EtiquetaYValor
		err = filas.Scan(&ev.Etiqueta, &ev.Valor)
		if err != nil {
			log.Printf("Error al escanear totalVendidoPorDia:\n%q", err)
		}
		dias = append(dias, ev)
	}
	return dias
}
