package main

import (
	"log"
)

type VentasContadoController struct {
	AjustesUsuario AjustesDeUsuarioLogueado
}

func (v *VentasContadoController) anular(idVenta int) {
	pc := ProductosController{
		AjustesUsuario: AjustesDeUsuarioLogueado{
			httpResponseWriter: v.AjustesUsuario.httpResponseWriter,
			httpRequest:        v.AjustesUsuario.httpRequest,
		},
	}
	productos := pc.deUnaVenta(idVenta)
	db, err := v.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error obteniendo contexto:\n%q", err)
		return
	}
	consultaActualizarProducto := "UPDATE productos SET existencia = existencia + ? WHERE idProducto = ?"
	sentenciaPreparadaActualizarProducto, err := tx.Prepare(consultaActualizarProducto)
	if err != nil {
		log.Printf("Error preparando consulta para actualizar existencia de producto:\n%q", err)
		return
	}
	for _, producto := range productos {
		_, err = sentenciaPreparadaActualizarProducto.Exec(producto.Cantidad, producto.Numero)
		if err != nil {
			log.Printf("Error actualizando existencia de producto:\n%q", err)
		}
	}
	consultaEliminarVenta := "DELETE FROM ventas_contado WHERE idVenta = ?"
	consultaPreparadaEliminarVenta, err := tx.Prepare(consultaEliminarVenta)
	if err != nil {
		log.Printf("Error preparando consulta para eliminar venta:\n%q", err)
	}
	_, err = consultaPreparadaEliminarVenta.Exec(idVenta)
	if err != nil {
		log.Printf("Error ejecutando sentencia de eliminar venta:\n%q", err)
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Error haciendo commit:\n%q", err)
	}
}

func (v *VentasContadoController) deUnCliente(idCliente int) *[]VentaContado {
	db, err := v.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT idVenta, monto, fecha
						FROM ventas_contado
          WHERE idCliente = ?
						`, idCliente)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener ventas al contado de un cliente:\n%q", err)
		return nil
	}

	defer filas.Close()

	ventas := []VentaContado{}
	for filas.Next() {
		var venta VentaContado
		err := filas.Scan(&venta.Numero, &venta.Total, &venta.Fecha)
		if err != nil {
			log.Printf("Error al escanear detalles de venta al contado :\n%q", err)
		}
		ventas = append(ventas, venta)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear ventas al contado de un cliente:\n%q", err)
	}
	return &ventas
}
func (v *VentasContadoController) una(id int) *VentaContado {
	db, err := v.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT ventas_contado.idVenta, ventas_contado.monto, ventas_contado.pago, ventas_contado.fecha, clientes.nombreCompleto, usuarios.nombre 
FROM ventas_contado 
INNER JOIN usuarios ON usuarios.idUsuario = ventas_contado.idUsuario 
INNER JOIN clientes ON clientes.idCliente = ventas_contado.idCliente 
WHERE ventas_contado.idVenta = ?;`, id)
	if err != nil {
		log.Printf("Error en consulta al consultar una venta:\n%q", err)
		return nil
	}

	defer filas.Close()
	var venta VentaContado
	if filas.Next() {
		err := filas.Scan(&venta.Numero, &venta.Total, &venta.Pago, &venta.Fecha, &venta.Cliente.Nombre, &venta.Usuario.Nombre)
		if err != nil {
			log.Printf("Error en consulta al consultar una venta:\n%q", err)
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: v.AjustesUsuario.httpResponseWriter,
				httpRequest:        v.AjustesUsuario.httpRequest,
			},
		}
		venta.Productos = pc.deUnaVenta(id)
	}
	return &venta
}
func (v *VentasContadoController) nueva(vc *VentaContado) int64 {

	consultaVenta := "INSERT INTO ventas_contado (monto, pago, fecha, idCliente, idUsuario) VALUES (?, ?, ?, ?, ?)"

	db, err := v.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error obteniendo contexto:\n%q", err)
		return 0
	}

	sentenciaPreparadaVenta, err := tx.Prepare(consultaVenta)
	if err != nil {
		log.Printf("Error preparando consulta:\n%q", err)
		return 0
	}

	resultado, err := sentenciaPreparadaVenta.Exec(
		vc.Total,
		vc.Pago,
		obtenerFechaActualFormateada(),
		vc.Cliente.Numero,
		vc.Usuario.Numero)

	if err != nil {
		log.Printf("Error insertando venta:\n%q", err)
	}
	tx.Commit()

	tx, err = db.Begin()

	consultaProducto := "INSERT INTO productos_vendidos (idProducto, codigoBarras, idVenta, descripcion, precioCompra, precioVenta, precioVentaOriginal, cantidadVendida) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	consultaActualizarProducto := "UPDATE productos SET existencia = existencia - ? WHERE idProducto = ?"
	numeroVenta, _ := v.ultimoNumero()
	sentenciaPreparadaProducto, err := tx.Prepare(consultaProducto)
	if err != nil {
		log.Printf("Error preparando consulta para insertar producto vendido:\n%q", err)
		return 0
	}
	sentenciaPreparadaActualizarProducto, err := tx.Prepare(consultaActualizarProducto)
	if err != nil {
		log.Printf("Error preparando consulta para actualizar existencia de producto:\n%q", err)
		return 0
	}

	for _, producto := range vc.Productos {
		_, err = sentenciaPreparadaProducto.Exec(producto.Numero, producto.CodigoBarras, numeroVenta, producto.Descripcion,
			producto.PrecioCompra, producto.PrecioVenta, producto.PrecioVentaOriginal, producto.Cantidad)
		if err != nil {
			log.Printf("Error insertando producto vendido:\n%q", err)
		}
		_, err = sentenciaPreparadaActualizarProducto.Exec(producto.Cantidad, producto.Numero)
		if err != nil {
			log.Printf("Error actualizando existencia de producto:\n%q", err)
		}
	}
	tx.Commit()
	idVenta, _ := resultado.LastInsertId()
	vc.Numero = int(idVenta)
	return idVenta
}

func (v *VentasContadoController) ultimoNumero() (ultimoNumero int, err error) {
	db, err := v.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT idVenta FROM ventas_contado ORDER BY idVenta DESC LIMIT 1;")
	if err != nil {
		log.Printf("Error al consultar el último número de venta_contado:\n%q", err)
		return -1, err
	}

	defer filas.Close()

	if !filas.Next() {
		return 1, nil
	}
	var ultimoRowid int
	err = filas.Scan(&ultimoRowid)
	if err != nil {
		log.Printf("Error al escanear idVenta de venta_contado:\n%q", err)
		return -1, err
	}
	return ultimoRowid, nil
}

func (v *VentasContadoController) enPeriodo(fechaInicio, fechaFin string) []DetalleVentaContado {
	db, err := v.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	/*
		Devuelve todas las ventas al contado, pero el id de venta se puede
		repetir en caso de que sea más de un producto el vendido.

		Por lo tanto, se tiene que implementar un algoritmo para acomodar
		estos datos

		Pudo hacerse un group_concat pero no queremos que se rompa
		en caso de que sean muchos productos
	*/
	filas, err := db.Query(`SELECT ventas_contado.idVenta,
									ventas_contado.monto,
									ventas_contado.fecha,
									ventas_contado.idCliente,
									ventas_contado.idUsuario,
									productos_vendidos.idProducto,
									productos_vendidos.codigoBarras,
									productos_vendidos.descripcion,
									productos_vendidos.precioCompra,
									productos_vendidos.precioVenta,
									productos_vendidos.precioVentaOriginal,
									productos_vendidos.cantidadVendida,
									clientes.nombreCompleto AS cliente,
									usuarios.nombre AS usuario
							FROM ventas_contado
							INNER JOIN productos_vendidos
							ON ventas_contado.idVenta = productos_vendidos.idVenta
							INNER JOIN clientes
							ON ventas_contado.idCliente = clientes.idCliente
							INNER JOIN usuarios
							ON ventas_contado.idUsuario = usuarios.idUsuario
							WHERE fecha > ?
							AND fecha < ?
							ORDER BY  ventas_contado.idVenta DESC;`, fechaInicio, fechaFin)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener todos los clientes:\n%q", err)
		return nil
	}

	defer filas.Close()

	ventas := []DetalleVentaContado{}
	for filas.Next() {
		var v DetalleVentaContado
		err := filas.Scan(&v.IdVenta, &v.Monto, &v.Fecha, &v.Cliente.Numero, &v.Usuario.Numero, &v.Producto.Numero, &v.Producto.CodigoBarras, &v.Producto.Descripcion, &v.Producto.PrecioCompra, &v.Producto.PrecioVenta, &v.Producto.PrecioVentaOriginal, &v.Producto.Cantidad, &v.Cliente.Nombre, &v.Usuario.Nombre)
		if err != nil {
			log.Printf("Error al escanear detalles de venta:\n%q", err)
		}
		ventas = append(ventas, v)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todas las ventas:\n%q", err)
	}
	return ventas
}
