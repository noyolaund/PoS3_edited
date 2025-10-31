package main

import (
	"log"
)

type CajaController struct {
	AjustesUsuario AjustesDeUsuarioLogueado
}

func (c *CajaController) paraReportePorUsuario(fechaInicio, fechaFin string, idUsuario int) *ReporteCaja {
	consulta := `SELECT 
					(SELECT COALESCE(SUM(monto),
						0) AS total
					FROM ventas_contado
					WHERE fecha > ?
							AND fecha < ? AND idUsuario = ?) AS total_ventas_contado, 
					(SELECT COALESCE(SUM(anticipo),
						0) AS total
					FROM apartados
					WHERE fecha > ?
							AND fecha < ? AND idUsuario = ?) AS total_anticipo, 
					(SELECT COALESCE(SUM(monto),
						0) AS total
					FROM abonos
					WHERE fecha > ?
							AND fecha < ? AND idUsuario = ?) AS total_abonado, 
					(SELECT COALESCE(SUM(monto),
						0) AS total
					FROM ingresos
					WHERE fecha > ?
							AND fecha < ? AND idUsuario = ?) AS total_ingresos, 
					(SELECT COALESCE(SUM(monto),
						0) AS total
					FROM egresos
					WHERE fecha > ?
							AND fecha < ? AND idUsuario = ?) AS total_egresos;`

	db, err := c.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo base de datos: %q", err)
		panic(err)
	}

	var reporte ReporteCaja
	defer db.Close()
	filas, err := db.Query(consulta, fechaInicio, fechaFin, idUsuario, fechaInicio, fechaFin, idUsuario, fechaInicio, fechaFin, idUsuario, fechaInicio, fechaFin, idUsuario, fechaInicio, fechaFin, idUsuario)
	if err != nil {
		log.Printf("Error al consultar el reporte para caja en un período y por usuario:\n%q", err)
	}

	defer filas.Close()

	filas.Next()

	err = filas.Scan(&reporte.VentasContado, &reporte.Anticipos, &reporte.Abonos, &reporte.Ingresos, &reporte.Egresos)
	if err != nil {
		log.Printf("Error al escanear el reporte para caja en un período y por usuario:\n%q", err)
	}
	return &reporte
}
func (c *CajaController) paraReporte(fechaInicio, fechaFin string) *ReporteCaja {
	consulta := `SELECT 
					(SELECT COALESCE(SUM(monto),
						0) AS total
					FROM ventas_contado
					WHERE fecha > ?
							AND fecha < ?) AS total_ventas_contado, 
					(SELECT COALESCE(SUM(anticipo),
						0) AS total
					FROM apartados
					WHERE fecha > ?
							AND fecha < ?) AS total_anticipo, 
					(SELECT COALESCE(SUM(monto),
						0) AS total
					FROM abonos
					WHERE fecha > ?
							AND fecha < ?) AS total_abonado, 
					(SELECT COALESCE(SUM(monto),
						0) AS total
					FROM ingresos
					WHERE fecha > ?
							AND fecha < ?) AS total_ingresos, 
					(SELECT COALESCE(SUM(monto),
						0) AS total
					FROM egresos
					WHERE fecha > ?
							AND fecha < ?) AS total_egresos;`

	db, err := c.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo base de datos: %q", err)
		panic(err)
	}

	var reporte ReporteCaja
	defer db.Close()
	filas, err := db.Query(consulta, fechaInicio, fechaFin, fechaInicio, fechaFin, fechaInicio, fechaFin, fechaInicio, fechaFin, fechaInicio, fechaFin)
	if err != nil {
		log.Printf("Error al consultar el reporte para caja en un período:\n%q", err)
	}

	defer filas.Close()

	filas.Next()

	err = filas.Scan(&reporte.VentasContado, &reporte.Anticipos, &reporte.Abonos, &reporte.Ingresos, &reporte.Egresos)
	if err != nil {
		log.Printf("Error al escanear el reporte para caja:\n%q", err)
	}
	return &reporte
}
func (c *CajaController) nuevoIngreso(ingreso *Ingreso) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "ingresos",
		AjustesUsuario: c.AjustesUsuario,
	}
	ayudante.insertar(map[string]interface{}{
		"monto":       ingreso.Monto,
		"descripcion": ingreso.Descripcion,
		"fecha":       obtenerFechaActualFormateada(),
		"idUsuario":   ingreso.Usuario.Numero,
	})
}
func (c *CajaController) nuevoEgreso(egreso *Egreso) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "egresos",
		AjustesUsuario: c.AjustesUsuario,
	}
	ayudante.insertar(map[string]interface{}{
		"monto":       egreso.Monto,
		"descripcion": egreso.Descripcion,
		"fecha":       obtenerFechaActualFormateada(),
		"idUsuario":   egreso.Usuario.Numero,
	})
}

func (c *CajaController) ingresosEnPeriodo(fechaInicio, fechaFin string) []Ingreso {
	db, err := c.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo base de datos: %q", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT ingresos.monto,
								ingresos.descripcion,
								ingresos.fecha,
								usuarios.nombre
							FROM ingresos
							INNER JOIN usuarios
							ON usuarios.idUsuario = ingresos.idUsuario
							AND ingresos.fecha > ?
							AND ingresos.fecha < ?;`, fechaInicio, fechaFin)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener los ingresos en un período:\n%q", err)
		return nil
	}

	defer filas.Close()

	ingresos := []Ingreso{}
	for filas.Next() {
		var ingreso Ingreso
		err := filas.Scan(&ingreso.Monto, &ingreso.Descripcion, &ingreso.Fecha, &ingreso.Usuario.Nombre)
		if err != nil {
			log.Printf("Error al escanear un ingreso:\n%q", err)
		}
		ingresos = append(ingresos, ingreso)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todos los ingresos:\n%q", err)
	}
	return ingresos
}

func (c *CajaController) egresosEnPeriodo(fechaInicio, fechaFin string) []Egreso {
	db, err := c.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo base de datos: %q", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT egresos.monto,
								egresos.descripcion,
								egresos.fecha,
								usuarios.nombre
							FROM egresos
							INNER JOIN usuarios
							ON usuarios.idUsuario = egresos.idUsuario
							AND egresos.fecha > ?
							AND egresos.fecha < ?;`, fechaInicio, fechaFin)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener los egresos en un período:\n%q", err)
		return nil
	}

	defer filas.Close()

	egresos := []Egreso{}
	for filas.Next() {
		var egreso Egreso
		err := filas.Scan(&egreso.Monto, &egreso.Descripcion, &egreso.Fecha, &egreso.Usuario.Nombre)
		if err != nil {
			log.Printf("Error al escanear un egreso:\n%q", err)
		}
		egresos = append(egresos, egreso)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todos los egresos:\n%q", err)
	}
	return egresos
}
