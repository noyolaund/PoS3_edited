package main

import (
	"fmt"

	"log"
	"time"
)

type ApartadosController struct {
	AjustesUsuario AjustesDeUsuarioLogueado
}

func (a *ApartadosController) detallesDeUnAbono(idAbono, idApartado int) *Abono {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}
	defer db.Close()
	filas, err := db.Query(`SELECT abonos.idAbono,
abonos.monto,
abonos.fecha,
abonos.idApartado,
abonos.pago,
usuarios.nombre AS usuario
FROM abonos
INNER JOIN usuarios
ON abonos.idUsuario = usuarios.idUsuario
AND abonos.idAbono = ? AND abonos.idApartado = ? LIMIT 1;`, idAbono, idApartado)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener detalles de un abono de un apartado:\n%q", err)
		return nil
	}

	defer filas.Close()

	var abono Abono
	if filas.Next() {

		err = filas.Scan(&abono.IdAbono, &abono.Monto, &abono.Fecha, &abono.IdApartado, &abono.Pago, &abono.Usuario.Nombre)
	} else {
		log.Printf("No había filas al escanear un abono. idAbono: %v. IdApartado: %v", idAbono, idApartado)
	}
	if err != nil {
		log.Printf("Error al escanear detalles de abono:\n%q", err)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear detalles de un abono:\n%q", err)
	}
	return &abono
}

func (a *ApartadosController) deUnCliente(idCliente int) *[]Apartado {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT apartados.idApartado,
								apartados.monto,
								apartados.abonado,
								apartados.anticipo,
								apartados.fecha,
								apartados.fechaVencimiento
						FROM apartados
          WHERE apartados.idCliente = ?
						`, idCliente)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener apartados de un cliente:\n%q", err)
		return nil
	}

	defer filas.Close()

	apartados := []Apartado{}
	for filas.Next() {
		var apartado Apartado
		err := filas.Scan(&apartado.Numero, &apartado.Total, &apartado.Abonado,
			&apartado.Anticipo, &apartado.Fecha, &apartado.FechaVencimiento)
		if err != nil {
			log.Printf("Error al escanear detalles de apartado:\n%q", err)
		}
		apartados = append(apartados, apartado)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear apartados de un cliente:\n%q", err)
	}
	return &apartados
}
func (a *ApartadosController) cambiarProducto(idApartado, idProductoAnterior, idProductoNuevo int) bool {
	pc := ProductosController{
		AjustesUsuario: AjustesDeUsuarioLogueado{
			httpResponseWriter: a.AjustesUsuario.httpResponseWriter,
			httpRequest:        a.AjustesUsuario.httpRequest,
		},
	}
	productoAnterior := pc.unoApartado(idApartado, idProductoAnterior)
	productoNuevo := pc.porRowid(idProductoNuevo)
	if productoAnterior == nil || productoNuevo == nil {
		return false
	}
	apartado := a.uno(idApartado)
	if productoNuevo.PrecioVenta >= productoAnterior.PrecioVenta {
		// Actualizar el monto del apartado
		ayudante := AyudanteBaseDeDatos{
			nombreTabla:    "apartados",
			AjustesUsuario: a.AjustesUsuario,
		}
		nuevoMonto := apartado.Total - productoAnterior.PrecioVenta + productoNuevo.PrecioVenta
		ayudante.actualizarDonde("idApartado", idApartado, map[string]interface{}{
			"monto": nuevoMonto,
		})

		//Eliminar el producto que será cambiado, e insertar el nuevo
		consultaParaEliminar := `DELETE FROM productos_apartados WHERE idApartado = ? AND idProducto = ?;`
		consultaParaInsertar := `INSERT INTO productos_apartados 
(idApartado, idProducto, codigoBarras, descripcion, precioVenta, precioVentaOriginal, precioCompra, cantidadVendida) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		db, err := a.AjustesUsuario.obtenerBaseDeDatos()
		if err != nil {
			panic(err)
		}

		defer db.Close()
		tx, err := db.Begin()

		if err != nil {
			log.Printf("Error obteniendo contexto:\n%q", err)
			return false
		}

		sentenciaParaEliminar, err := tx.Prepare(consultaParaEliminar)
		if err != nil {
			log.Printf("Error preparando consulta para eliminar:\n%q", err)
			return false
		}

		sentenciaParaInsertar, err := tx.Prepare(consultaParaInsertar)
		if err != nil {
			log.Printf("Error preparando sentencia para insertar:\n%q", err)
			return false
		}

		_, err = sentenciaParaEliminar.Exec(idApartado, idProductoAnterior)
		if err != nil {
			log.Printf("Error ejecutando sentencia para eliminar:\n%q", err)
		}
		_, err = sentenciaParaInsertar.Exec(idApartado, productoNuevo.Numero, productoNuevo.CodigoBarras,
			productoNuevo.Descripcion, productoNuevo.PrecioVenta, productoNuevo.PrecioVenta,
			productoNuevo.PrecioCompra, 1)

		if err != nil {
			log.Printf("Error ejecutando sentencia para insertar:\n%q", err)
		}
		tx.Commit()
		return true
	} else {
		log.Println("No se puede cambiar el producto, porque el precio es menor al anterior")
	}
	return false
}
func (a *ApartadosController) cambiarFechaDeVencimiento(nuevaFecha string, idApartado int) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "apartados",
		AjustesUsuario: a.AjustesUsuario,
	}
	ayudante.actualizarDonde("idApartado", idApartado, map[string]interface{}{
		"fechaVencimiento": nuevaFecha,
	})
}

/*
  Apartados que vencerán en los próximos días a partir
  de hoy.
  Se incluyen aquellos que vencen hoy, y también los que vencen
  dentro de "dias" incluyentes

  Por ejemplo, si hoy es 20 de junio y se quieren ver los próximos a vencer
  en 7 días, incluirá los que vencen el 20, 21, ... hasta el 27

*/
func (a *ApartadosController) proximosAVencer(dias int) *[]DetalleApartado {
	hoy := time.Now()
	//Se agrega el +1 a días para que en SQLite se pueda trabajar con BETWEEN al formatear fechas
	otroDia := hoy.AddDate(0, 0, dias+1)
	fechaInicio := fmt.Sprintf("%d-%02d-%02d",
		hoy.Year(), hoy.Month(), hoy.Day())
	fechaFin := fmt.Sprintf("%d-%02d-%02d",
		otroDia.Year(), otroDia.Month(), otroDia.Day())
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT apartados.idApartado, apartados.monto, apartados.abonado, apartados.anticipo, 
apartados.fecha, apartados.fechaVencimiento, apartados.idCliente, apartados.idUsuario,
productos_apartados.idProducto,
productos_apartados.codigoBarras,
productos_apartados.descripcion,
productos_apartados.precioCompra,
productos_apartados.precioVenta,
productos_apartados.precioVentaOriginal,
productos_apartados.cantidadVendida,
clientes.nombreCompleto AS cliente,
usuarios.nombre AS usuario
FROM apartados
INNER JOIN productos_apartados
ON apartados.idApartado = productos_apartados.idApartado
INNER JOIN clientes
ON apartados.idCliente = clientes.idCliente
INNER JOIN usuarios
ON apartados.idUsuario = usuarios.idUsuario 
AND
ROUND(apartados.monto - apartados.abonado - apartados.anticipo) > 0 
AND apartados.fechaVencimiento BETWEEN ? AND ?
ORDER BY  apartados.idApartado DESC;`, fechaInicio, fechaFin)

	if err != nil {
		log.Printf("Error al realizar la consulta para obtener apartados próximos a vencer:\n%q", err)
		return nil
	}

	defer filas.Close()

	apartados := []DetalleApartado{}
	for filas.Next() {
		var da DetalleApartado
		err := filas.Scan(&da.IdApartado, &da.Monto, &da.Abonado, &da.Anticipo, &da.Fecha, &da.FechaVencimiento,
			&da.Cliente.Numero, &da.Usuario.Numero, &da.Producto.Numero, &da.Producto.CodigoBarras,
			&da.Producto.Descripcion, &da.Producto.PrecioCompra, &da.Producto.PrecioVenta, &da.Producto.PrecioVentaOriginal,
			&da.Producto.Cantidad, &da.Cliente.Nombre, &da.Usuario.Nombre)
		if err != nil {
			log.Printf("Error al escanear detalles de apartado:\n%q", err)
		}
		apartados = append(apartados, da)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear apartados próximos a vencer:\n%q", err)
	}
	return &apartados
}

func (a *ApartadosController) uno(idApartado int) *Apartado {

	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT apartados.idApartado, apartados.monto, apartados.pago, apartados.abonado, apartados.anticipo, 
apartados.fecha, apartados.fechaVencimiento,clientes.nombreCompleto, usuarios.nombre 
FROM apartados 
INNER JOIN usuarios ON usuarios.idUsuario = apartados.idUsuario 
INNER JOIN clientes ON clientes.idCliente = apartados.idCliente 
WHERE apartados.idApartado = ?;`, idApartado)
	if err != nil {
		log.Printf("Error en consulta al consultar un apartado:\n%q", err)
		return nil
	}

	defer filas.Close()
	var apartado Apartado
	if filas.Next() {
		err := filas.Scan(&apartado.Numero, &apartado.Total, &apartado.Pago, &apartado.Abonado, &apartado.Anticipo, &apartado.Fecha,
			&apartado.FechaVencimiento, &apartado.Cliente.Nombre, &apartado.Usuario.Nombre)
		if err != nil {
			log.Printf("Error en consulta al consultar un apartado:\n%q", err)
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: a.AjustesUsuario.httpResponseWriter,
				httpRequest:        a.AjustesUsuario.httpRequest,
			},
		}
		apartado.Productos = pc.deUnApartado(idApartado)
	}
	return &apartado
}
func (a *ApartadosController) liquidar(idApartado int) {
	/*
		Resta la existencia de los productos apartados. Se supone que los productos se restan
		cuando se liquida el apartado
	*/
	consulta := "UPDATE productos SET existencia = existencia - ? WHERE idProducto = ?;"

	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error obteniendo contexto de base de datos:\n%q", err)
		return
	}

	sentenciaPreparadaRestarProductos, err := tx.Prepare(consulta)
	if err != nil {
		log.Printf("Error preparando consulta:\n%q", err)
		return
	}

	pc := ProductosController{
		AjustesUsuario: AjustesDeUsuarioLogueado{
			httpResponseWriter: a.AjustesUsuario.httpResponseWriter,
			httpRequest:        a.AjustesUsuario.httpRequest,
		},
	}

	productosApartados := pc.deUnApartado(idApartado)
	for _, productoApartado := range productosApartados {
		_, err = sentenciaPreparadaRestarProductos.Exec(productoApartado.Cantidad, productoApartado.Numero)
		if err != nil {
			log.Printf("Error ejecutando sentencia:\n%q", err)
		} else {
			log.Println("Dato insertado correctamente")
		}
	}

	err = tx.Commit()
	log.Printf("Al hacer commit del apartado tenemos error: %v", err)
}

func (a *ApartadosController) abonosDe(idApartado int) []Abono {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT abonos.idAbono, abonos.monto, abonos.fecha, usuarios.nombre FROM abonos INNER JOIN usuarios ON usuarios.idUsuario = abonos.idUsuario WHERE abonos.idApartado = ?;", idApartado)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener todos los clientes:\n%q", err)
		return nil
	}

	defer filas.Close()

	abonos := []Abono{}
	for filas.Next() {
		var abono Abono
		err := filas.Scan(&abono.IdAbono, &abono.Monto, &abono.Fecha, &abono.Usuario.Nombre)
		if err != nil {
			log.Printf("Error al escanear detalles de venta:\n%q", err)
		}
		abonos = append(abonos, abono)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todas las abonos:\n%q", err)
	}
	return abonos
}

func (a *ApartadosController) abonar(abono *Abono) int64 {
	if !a.puedeAbonar(abono.IdApartado) {
		return 0
	}
	abono.Fecha = obtenerFechaActualFormateada()
	consultaAbono := "INSERT INTO abonos(monto, pago, fecha, idApartado, idUsuario) VALUES(?, ?, ?, ?, ?)"
	consultaApartado := "UPDATE apartados SET abonado = abonado + ? WHERE idApartado = ?"
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error obteniendo contexto de base de datos:\n%q", err)
		return 0
	}

	sentenciaAbono, err := tx.Prepare(consultaAbono)
	if err != nil {
		log.Printf("Error preparando consulta para abonar:\n%q", err)
		return 0
	}

	sentenciaApartado, err := tx.Prepare(consultaApartado)
	if err != nil {
		log.Printf("Error preparando consulta para apartado:\n%q", err)
		return 0
	}

	var ultimoId int64

	resultado, err := sentenciaAbono.Exec(abono.Monto, abono.Pago, abono.Fecha, abono.IdApartado, abono.Usuario.Numero)
	if err != nil {
		log.Printf("Error ejecutando sentencia de abono:\n%q", err)
	}
	ultimoId, err = resultado.LastInsertId()
	if err != nil {
		log.Printf("Error obteniendo último ID:\n%q", err)
	}
	_, err = sentenciaApartado.Exec(abono.Monto, abono.IdApartado)
	if err != nil {
		log.Printf("Error ejecutando sentencia de apartado:\n%q", err)
	}
	tx.Commit()

	//Aquí
	filas, err := db.Query("SELECT abonado + anticipo >= monto FROM apartados WHERE idApartado = ? LIMIT 1;", abono.IdApartado)
	if err != nil {
		log.Printf("Error al consultar si el apartado debería liquidarse:\n%q", err)
		return 0
	}

	if !filas.Next() {
		return 0
	}
	var estaLiquidado int
	err = filas.Scan(&estaLiquidado)
	if err != nil {
		log.Printf("Error al escanear si está liquidado el apartado:\n%q", err)
		return 0
	}
	filas.Close()
	db.Close()
	if estaLiquidado == 1 {
		a.liquidar(abono.IdApartado)
	}
	return ultimoId

}
func (a *ApartadosController) nuevo(apartado *Apartado) *Apartado {

	fecha := obtenerFechaActualFormateada()

	consultaApartado := `INSERT INTO apartados (monto, pago, anticipo, abonado, fecha, fechaVencimiento, idCliente, idUsuario) 
VALUES (?, ?, ?, 0, ?, ?, ?, ?)`

	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error obteniendo contexto de base de datos:\n%q", err)
		return nil
	}

	sentenciaPreparadaApartado, err := tx.Prepare(consultaApartado)
	if err != nil {
		log.Printf("Error preparando consulta de apartado:\n%q", err)
		return nil
	}

	resultado, err := sentenciaPreparadaApartado.Exec(apartado.Total, apartado.Pago, apartado.Anticipo, fecha,
		apartado.FechaVencimiento, apartado.Cliente.Numero, apartado.Usuario.Numero)
	if err != nil {
		log.Printf("Error ejecutando sentencia de apartado:\n%q", err)
	}
	tx.Commit()

	tx, err = db.Begin()

	consultaProducto := `INSERT INTO productos_apartados 
(idApartado, idProducto, codigoBarras, descripcion, precioVenta, precioVentaOriginal, precioCompra, cantidadVendida) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	numeroApartado, _ := a.ultimoNumero()
	sentenciaPreparadaProducto, err := tx.Prepare(consultaProducto)
	if err != nil {
		log.Printf("Error preparando consulta de producto:\n%q", err)
		return nil
	}

	for _, producto := range apartado.Productos {
		_, err = sentenciaPreparadaProducto.Exec(numeroApartado, producto.Numero, producto.CodigoBarras,
			producto.Descripcion, producto.PrecioVenta, producto.PrecioVentaOriginal, producto.PrecioCompra,
			producto.Cantidad)
		if err != nil {
			log.Printf("Error al hacer nuevo apartado:\n%q", err)
		}
	}
	tx.Commit()
	id, _ := resultado.LastInsertId()
	apartado.Numero = int(id)
	return apartado
}

func (a *ApartadosController) puedeAbonar(idApartado int) bool {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT fechaVencimiento >= ? FROM apartados WHERE idApartado = ?;",
		fechaDeHoyAMedianoche(), idApartado)
	if err != nil {
		log.Printf("Error al consultar si un apartado se puede abonar:\n%q", err)
		return false
	}

	defer filas.Close()

	if !filas.Next() {
		return false
	}
	var puede int
	err = filas.Scan(&puede)
	if err != nil {
		log.Printf("Error al escanear si un apartado se puede abonar:\n%q", err)
		return false
	}
	return puede == 1
}

func (a *ApartadosController) ultimoNumero() (ultimoNumero int, err error) {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT idApartado FROM apartados ORDER BY fecha DESC;")
	if err != nil {
		log.Printf("Error al consultar el último número de apartado:\n%q", err)
		return -1, err
	}

	defer filas.Close()

	if !filas.Next() {
		return 1, nil
	}
	var ultimoRowid int
	err = filas.Scan(&ultimoRowid)
	if err != nil {
		log.Printf("Error al escanear idApartado de apartado:\n%q", err)
		return -1, err
	}
	return ultimoRowid, nil
}

func (a *ApartadosController) enPeriodo(fechaInicio, fechaFin string) []DetalleApartado {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT apartados.idApartado,
								apartados.monto,
								apartados.abonado,
								apartados.anticipo,
								apartados.fecha,
								apartados.fechaVencimiento,
								apartados.idCliente,
								apartados.idUsuario,
								productos_apartados.idProducto,
								productos_apartados.codigoBarras,
								productos_apartados.descripcion,
								productos_apartados.precioCompra,
								productos_apartados.precioVenta,
								productos_apartados.precioVentaOriginal,
								productos_apartados.cantidadVendida,
								clientes.nombreCompleto AS cliente,
								usuarios.nombre AS usuario
						FROM apartados
						INNER JOIN productos_apartados
						ON apartados.idApartado = productos_apartados.idApartado
						INNER JOIN clientes
						ON apartados.idCliente = clientes.idCliente
						INNER JOIN usuarios
						ON apartados.idUsuario = usuarios.idUsuario
						WHERE fecha > ?
							AND fecha < ?
						ORDER BY  apartados.idApartado DESC;`, fechaInicio, fechaFin)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener todos los clientes:\n%q", err)
		return nil
	}

	defer filas.Close()

	apartados := []DetalleApartado{}
	for filas.Next() {
		var da DetalleApartado
		err := filas.Scan(&da.IdApartado, &da.Monto, &da.Abonado, &da.Anticipo, &da.Fecha, &da.FechaVencimiento,
			&da.Cliente.Numero, &da.Usuario.Numero, &da.Producto.Numero, &da.Producto.CodigoBarras,
			&da.Producto.Descripcion, &da.Producto.PrecioCompra, &da.Producto.PrecioVenta, &da.Producto.PrecioVentaOriginal,
			&da.Producto.Cantidad, &da.Cliente.Nombre, &da.Usuario.Nombre)
		if err != nil {
			log.Printf("Error al escanear detalles de apartado:\n%q", err)
		}
		apartados = append(apartados, da)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todas las apartados:\n%q", err)
	}
	return apartados
}

func (a *ApartadosController) pendientes() []DetalleApartado {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT apartados.idApartado,
								apartados.monto,
								apartados.abonado,
								apartados.anticipo,
								apartados.fecha,
								apartados.fechaVencimiento,
								apartados.idCliente,
								apartados.idUsuario,
								productos_apartados.idProducto,
								productos_apartados.codigoBarras,
								productos_apartados.descripcion,
								productos_apartados.precioCompra,
								productos_apartados.precioVenta,
								productos_apartados.precioVentaOriginal,
								productos_apartados.cantidadVendida,
								clientes.nombreCompleto AS cliente,
								usuarios.nombre AS usuario
						FROM apartados
						INNER JOIN productos_apartados
						ON apartados.idApartado = productos_apartados.idApartado
						INNER JOIN clientes
						ON apartados.idCliente = clientes.idCliente
						INNER JOIN usuarios
						ON apartados.idUsuario = usuarios.idUsuario
							AND ROUND(apartados.monto - apartados.abonado - apartados.anticipo) > 0
							AND apartados.fechaVencimiento >= ?
						ORDER BY  apartados.idApartado DESC;`, obtenerDiaActual())
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener todos los clientes:\n%q", err)
		return nil
	}

	defer filas.Close()

	apartados := []DetalleApartado{}
	for filas.Next() {
		var da DetalleApartado
		err := filas.Scan(&da.IdApartado, &da.Monto, &da.Abonado, &da.Anticipo, &da.Fecha, &da.FechaVencimiento,
			&da.Cliente.Numero, &da.Usuario.Numero, &da.Producto.Numero, &da.Producto.CodigoBarras,
			&da.Producto.Descripcion, &da.Producto.PrecioCompra, &da.Producto.PrecioVenta, &da.Producto.PrecioVentaOriginal,
			&da.Producto.Cantidad, &da.Cliente.Nombre, &da.Usuario.Nombre)
		if err != nil {
			log.Printf("Error al escanear detalles de apartado:\n%q", err)
		}
		apartados = append(apartados, da)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todas las apartados:\n%q", err)
	}
	return apartados
}

func (a *ApartadosController) totalAbonadoEnPeriodo(fechaInicio, fechaFin string) float64 {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT COALESCE(SUM(monto), 0) AS total FROM abonos WHERE fecha > ? AND fecha < ?;",
		fechaInicio, fechaFin)
	if err != nil {
		log.Printf("Error al consultar la suma de abonos en un período:\n%q", err)
		return 0
	}

	defer filas.Close()

	if !filas.Next() {
		return 0
	}
	var total float64
	err = filas.Scan(&total)
	if err != nil {
		log.Printf("Error al escanear la suma de abonos:\n%q", err)
		return 0
	}
	return total
}
