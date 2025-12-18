package main

import (
	"log"
)

type AjustesController struct {
	AjustesUsuario AjustesDeUsuarioLogueado
}

func (a *AjustesController) guardarValor(valor string, clave string) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "comun",
		AjustesUsuario: a.AjustesUsuario,
	}
	if a.existeClave(clave) {
		ayudante.actualizarDonde("clave", clave, map[string]interface{}{
			"valor": valor,
		})
	} else {
		ayudante.insertar(map[string]interface{}{
			"clave": clave,
			"valor": valor,
		})
	}
}

func (a *AjustesController) existeClave(clave string) bool {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		return false
	}
	defer db.Close()
	var conteo int
	err = db.QueryRow("SELECT COUNT(*) FROM comun WHERE clave = ?", clave).Scan(&conteo)
	return err == nil && conteo > 0
}

func (a *AjustesController) obtenerValor(clave string) string {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo base de datos: %q", err)
		panic(err)
	}

	var valor string
	defer db.Close()
	fila := db.QueryRow(`SELECT valor FROM comun WHERE clave = ?;`, clave)
	err = fila.Scan(&valor)
	if err != nil {
		log.Printf("Error al escanear modo impresión:\n%q", err)
	}
	return valor
}

func (a *AjustesController) obtenerOtros() *OtrosAjustes {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo base de datos: %q", err)
		panic(err)
	}

	var ajustes OtrosAjustes
	defer db.Close()
	filas, err := db.Query(`SELECT 
(SELECT valor FROM comun WHERE clave = "NUMERO_COPIAS_TICKET_CONTADO") AS copiasContado, 
(SELECT valor FROM comun WHERE clave = "NUMERO_COPIAS_TICKET_APARTADO") AS copiasApartado, 
(SELECT valor FROM comun WHERE clave = "NUMERO_COPIAS_TICKET_ABONO") AS copiasAbono, 
(SELECT valor FROM comun WHERE clave = "MODO_IMPRESION_CODIGOS") AS modoImpresion, 
(SELECT valor FROM comun WHERE clave = "MODO_LECTURA_CODIGOS") AS modoLectura;`)
	if err != nil {
		log.Printf("Error al consultar otros ajustes:\n%q", err)
	}

	defer filas.Close()
	if filas.Next() {
		err = filas.Scan(&ajustes.NumeroDeCopiasTicketContado, &ajustes.NumeroDeCopiasTicketApartado,
			&ajustes.NumeroDeCopiasTicketAbono, &ajustes.ModoImpresionCodigoDeBarras, &ajustes.ModoLecturaProductos)
		if err != nil {
			log.Printf("Error al escanear otros ajustes:\n%q", err)
		}
	}
	return &ajustes
}
func (a *AjustesController) guardarOtros(ajustes *OtrosAjustes) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "comun",
		AjustesUsuario: a.AjustesUsuario,
	}
	ayudante.actualizarDonde("clave", "MODO_IMPRESION_CODIGOS", map[string]interface{}{
		"valor": ajustes.ModoImpresionCodigoDeBarras,
	})
	ayudante.actualizarDonde("clave", "MODO_LECTURA_CODIGOS", map[string]interface{}{
		"valor": ajustes.ModoLecturaProductos,
	})
	ayudante.actualizarDonde("clave", "NUMERO_COPIAS_TICKET_CONTADO", map[string]interface{}{
		"valor": ajustes.NumeroDeCopiasTicketContado,
	})
	ayudante.actualizarDonde("clave", "NUMERO_COPIAS_TICKET_APARTADO", map[string]interface{}{
		"valor": ajustes.NumeroDeCopiasTicketApartado,
	})
	ayudante.actualizarDonde("clave", "NUMERO_COPIAS_TICKET_ABONO", map[string]interface{}{
		"valor": ajustes.NumeroDeCopiasTicketAbono,
	})
}
func (a *AjustesController) guardarNombreImpresora(nombre string) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "comun",
		AjustesUsuario: a.AjustesUsuario,
	}
	ayudante.actualizarDonde("clave", "NOMBRE_IMPRESORA", map[string]interface{}{
		"valor": nombre,
	})
}
func (a *AjustesController) obtenerNombreImpresora() string {
	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo base de datos: %q", err)
		panic(err)
	}

	var nombreImpresora string
	defer db.Close()
	filas, err := db.Query(`SELECT valor FROM comun WHERE clave = 'NOMBRE_IMPRESORA';`)
	if err != nil {
		log.Printf("Error al consultar nombre de la impresora:\n%q", err)
	}

	defer filas.Close()
	if filas.Next() {
		err = filas.Scan(&nombreImpresora)
		if err != nil {
			log.Printf("Error al escanear nombre de la impresora:\n%q", err)
		}
	}
	return nombreImpresora
}

func (a *AjustesController) obtenerDatosEmpresa() *DatosEmpresa {

	db, err := a.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error fatal abriendo base de datos: %q", err)
		panic(err)
	}

	var datosEmpresa DatosEmpresa
	defer db.Close()
	filas, err := db.Query(`SELECT nombre, direccion, telefono, mensajePersonal FROM empresa;`)
	if err != nil {
		log.Printf("Error al consultar los datos de la empresa:\n%q", err)
	}

	defer filas.Close()
	if filas.Next() {
		err = filas.Scan(&datosEmpresa.Nombre, &datosEmpresa.Direccion, &datosEmpresa.Telefono, &datosEmpresa.MensajePersonal)
		if err != nil {
			log.Printf("Error al escanear  los datos de la empresa:\n%q", err)
		}
	}
	return &datosEmpresa
}

func (a *AjustesController) guardarDatosEmpresa(empresa *DatosEmpresa) {
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "empresa",
		AjustesUsuario: a.AjustesUsuario,
	}
	ayudante.actualizarDonde("1", 1, map[string]interface{}{
		"nombre":          empresa.Nombre,
		"telefono":        empresa.Telefono,
		"mensajePersonal": empresa.MensajePersonal,
		"direccion":       empresa.Direccion,
	})
}
