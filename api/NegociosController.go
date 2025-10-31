package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
)

type NegociosController struct {
	AjustesUsuario AjustesDeUsuarioLogueado
}

func (n *NegociosController) obtenerNegocioParaEliminarPorToken(token string) (Negocio, error) {
	var negocio Negocio
	bd, err := obtenerBDNegocios()
	if err != nil {
		log.Printf("Error abriendo base de datos: %v", err)
	}
	defer bd.Close()
	fila := bd.QueryRow(`select negocios.id, negocios.nombre, negocios.correo
from notificaciones_eliminacion_negocios
       inner join negocios on notificaciones_eliminacion_negocios.id_negocio = negocios.id
where notificaciones_eliminacion_negocios.token = ?
limit 1;`, token)
	err = fila.Scan(&negocio.Id, &negocio.Nombre, &negocio.Correo)
	if err != nil {
		return negocio, err
	}
	return negocio, nil
}
func (n *NegociosController) eliminarNegocioPorToken(token string) bool {
	negocio, err := n.obtenerNegocioParaEliminarPorToken(token)
	if err != nil {
		log.Printf("Error obteniendo negocio para eliminar por token: %v", err)
		return false
	}
	if len(negocio.Correo) > 0 {
		n.eliminarNegocio(&negocio)
		n.reiniciarAutoincrementoNegocios()
		return true
	}
	return false
}

// No se debería usar si no se usa MySQLo si el sistema es para uso local
func (n *NegociosController) obtenerNegociosAntiguosQueNoHanIngresadoNunca() []Negocio {
	negocios := []Negocio{}
	bd, err := obtenerBDNegocios()
	if err != nil {
		log.Printf("Error obteniendo base de datos: %v", err)
		return negocios
	}
	defer bd.Close()
	filas, err := bd.Query(`select negocios.id, negocios.nombre, negocios.correo
from negocios
       left join accesos_negocios on negocios.id = accesos_negocios.id_negocio
where accesos_negocios.id_negocio is null
  and negocios.fecha_registro < date_sub(?, interval ? day)
group by negocios.id, negocios.nombre, negocios.correo, accesos_negocios.id_negocio;`, fechaDeHoySinHora(), DiasParaMarcarCuentaComoInactiva)
	if err != nil {
		log.Printf("Error consultando negocios antiguos que no han ingresado nunca: %v", err)
		return negocios
	}
	for filas.Next() {
		var negocio Negocio
		err := filas.Scan(&negocio.Id, &negocio.Nombre, &negocio.Correo)
		if err != nil {
			log.Printf("Error escaneando negocio antiguos que no han ingresado nunca: %v", err)
			return negocios
		}
		negocios = append(negocios, negocio)
	}
	return negocios
}

// No se debería usar si no se usa MySQL o si el sistema es para uso local
func (n *NegociosController) obtenerNegociosQueNoHanUsadoElSistemaRecientemente() []Negocio {
	negocios := []Negocio{}
	bd, err := obtenerBDNegocios()
	if err != nil {
		log.Printf("Error obteniendo base de datos: %v", err)
		return negocios
	}
	defer bd.Close()
	filas, err := bd.Query(`select negocios.id, negocios.nombre, negocios.correo
from accesos_negocios
       inner join negocios on accesos_negocios.id_negocio = negocios.id
where accesos_negocios.exitoso
  and negocios.verificado
group by negocios.id
having max(accesos_negocios.momento) < date_sub(?, interval ? day);`, fechaDeHoySinHora(), DiasParaMarcarCuentaComoInactiva)
	if err != nil {
		log.Printf("Error consultando negocios que no han accedido recientemente: %v", err)
		return negocios
	}
	for filas.Next() {
		var negocio Negocio
		err := filas.Scan(&negocio.Id, &negocio.Nombre, &negocio.Correo)
		if err != nil {
			log.Printf("Error escaneando negocio al obtener negocios que no han usado el sistema: %v", err)
			return negocios
		}
		negocios = append(negocios, negocio)
	}
	return negocios
}

func (n *NegociosController) numeroDeNotificacionesEnviadasANegocio(negocio *Negocio) (uint8, error) {
	bd, err := obtenerBDNegocios()
	if err != nil {
		return 0, err
	}
	defer bd.Close()
	var conteo uint8
	fila := bd.QueryRow("select count(id) from notificaciones_eliminacion_negocios where id_negocio = ?;", negocio.Id)
	err = fila.Scan(&conteo)
	if err != nil {
		return 0, err
	} else {
		return conteo, nil
	}

}

// Los negocios a los que se les ha mandado notificación y no han hecho caso
func (n *NegociosController) obtenerNegociosQueDebenEliminarse() []Negocio {
	//var negocios []Negocio
	negocios := []Negocio{}
	//clientes := []Cliente{}
	bd, err := obtenerBDNegocios()
	if err != nil {
		log.Printf("Error obteniendo base de datos al obtener negocios que deben eliminarse: %v", err)
		return negocios
	}
	defer bd.Close()
	filas, err := bd.Query(`select negocios.id, negocios.nombre, negocios.correo
from negocios
       inner join notificaciones_eliminacion_negocios on negocios.id = notificaciones_eliminacion_negocios.id_negocio
group by negocios.id
having count(notificaciones_eliminacion_negocios.id) >= ?;`, NumeroAvisosAntesDeEliminarCuenta)
	if err != nil {
		log.Printf("Error consultando negocios que deben eliminarse: %v", err)
		return negocios
	}
	for filas.Next() {
		var negocio Negocio
		err := filas.Scan(&negocio.Id, &negocio.Nombre, &negocio.Correo)
		if err != nil {
			log.Printf("Error escaneando negocio al obtener negocios que deben eliminarse: %v", err)
			return negocios
		}
		negocios = append(negocios, negocio)
	}
	return negocios
}

func (n *NegociosController) eliminarNegociosQueNoHanUsadoElSistema() {
	log.Println("Eliminando negocios que no han usado el sistema")
	negocios := n.obtenerNegociosQueDebenEliminarse()
	for _, negocio := range negocios {
		log.Printf("Eliminando negocio con correo %s y nombre %s", negocio.Correo, negocio.Nombre)
		err := n.eliminarNegocio(&negocio)
		if err != nil {
			log.Printf("Error eliminando negocio: %v", err)
		} else {
			log.Println("OK. Negocio eliminado")
		}
	}
	n.reiniciarAutoincrementoNegocios()
}

func (n *NegociosController) reiniciarAutoincrementoNegocios() {
	bd, err := obtenerBDNegocios()
	if err != nil {
		log.Printf("Error obteniendo BD para reiniciar autoincremento: %v", err)
		return
	}
	defer bd.Close()
	_, err = bd.Exec("ALTER TABLE negocios AUTO_INCREMENT = 1;")
	if err != nil {
		log.Printf("Error reiniciando contador: %v", err)
	}
}

func (n *NegociosController) eliminarNegocio(negocio *Negocio) error {
	bd, err := obtenerBDNegocios()
	if err != nil {
		log.Printf("Error obteniendo DB negocios para eliminar un negocio: %v", err)
		return err
	}
	defer bd.Close()
	_, err = bd.Exec("delete from negocios where id = ?", negocio.Id)
	if err != nil {
		log.Printf("Error eliminando de BD negocios: %v", err)
		return err
	}
	_, err = bd.Exec("delete from notificaciones_eliminacion_negocios where id_negocio = ?", negocio.Id)
	if err != nil {
		log.Printf("Error eliminando de notificaciones: %v", err)
		return err
	}
	_, err = bd.Exec("delete from accesos_negocios where id_negocio = ?", negocio.Id)
	if err != nil {
		log.Printf("Error eliminando de accesos: %v", err)
		return err
	}
	_, err = bd.Exec(fmt.Sprintf("drop database if exists %s", nombreBDAPartirDeNegocio(negocio.Id)))
	if err != nil {
		log.Printf("Error eliminando base de datos de negocio: %v", err)
		return err
	}
	return nil

}

func (n *NegociosController) enviarCorreoANegocioQueNoHaUsadoElSistema(negocio *Negocio, tokenEliminar string, numeroDeAviso uint8) (bool, error) {
	t, err := template.ParseFiles("plantillas_correos/recordatorio_uso_sistema.html")
	if err != nil {
		return false, err
	}
	buffer := new(bytes.Buffer)
	datos := map[string]string{
		"urlCompartir":           UrlAcortadaCompartirSistema,
		"urlEliminar":            fmt.Sprintf("%s/#/eliminar/%s", UrlBaseApp, tokenEliminar),
		"urlLoginSistema":        fmt.Sprintf("%s/#/login", UrlBaseApp),
		"negocio":                negocio.Nombre,
		"correoSoporteYContacto": CorreoSoporteYContacto,
		"urlBasePaginaWeb":       UrlBasePaginaWeb,
		"diasRestantes":          fmt.Sprintf("%d", DiasParaMarcarCuentaComoInactiva),
		"mensajeCompartir":       MensajeCompartir,
		"numeroDeAviso":          fmt.Sprintf("%d", numeroDeAviso),
		"numeroTotalDeAvisos":    fmt.Sprintf("%d", NumeroAvisosAntesDeEliminarCuenta),
	}
	destinatarios := []string{negocio.Correo}
	if err = t.Execute(buffer, datos); err != nil {
		return false, err
	}
	contenidoPlantilla := buffer.String()
	cuerpo := "To: " + destinatarios[0] + "\r\nSubject: " + AsuntoReactivarCuenta + "\r\n" + MIME + "\r\n" + contenidoPlantilla
	direccionServidor := fmt.Sprintf("%s:%s", GmailServidor, GmailPuerto)
	err = smtp.SendMail(direccionServidor, smtp.PlainAuth("", GmailCorreo, GmailPass, GmailServidor), GmailCorreo, destinatarios, []byte(cuerpo))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (n *NegociosController) enviarCorreoANegociosQueNoHanUsadoElSistema() {
	log.Println("Comenzando a enviar correos a negocios que no han usado el sistema")
	negocios := n.obtenerNegociosQueNoHanUsadoElSistemaRecientemente()
	for _, negocio := range negocios {
		log.Printf("Enviando correo al negocio %s <%s>", negocio.Nombre, negocio.Correo)
		token, err := n.obtenerTokenNoExistenteParaNotificarEliminacion()
		if err != nil {
			log.Printf("Error obteniendo token: %v", err)
			continue
		}
		n.registrarNotificacionDeEliminacion(&negocio, token)
		numeroDeAvisos, err := n.numeroDeNotificacionesEnviadasANegocio(&negocio)
		if err != nil {
			log.Printf("Error consultando número de avisos:%v", err)
			continue
		}
		log.Printf("Aviso %d de %d", numeroDeAvisos, NumeroAvisosAntesDeEliminarCuenta)

		_, err = n.enviarCorreoANegocioQueNoHaUsadoElSistema(&negocio, token, numeroDeAvisos)
		if err != nil {
			log.Printf("Error al enviar correo de notificación de eliminación: %v", err)
		}
	}
	log.Println("Terminado envío de correos a negocios que no han usado el sistema")
	log.Println("Comenzando a enviar correos a negocios antiguos que nunca han ingresado")
	negocios = n.obtenerNegociosAntiguosQueNoHanIngresadoNunca()
	for _, negocio := range negocios {
		log.Printf("Enviando correo al negocio %s <%s>", negocio.Nombre, negocio.Correo)
		token, err := n.obtenerTokenNoExistenteParaNotificarEliminacion()
		if err != nil {
			log.Printf("Error obteniendo token: %v", err)
			continue
		}
		n.registrarNotificacionDeEliminacion(&negocio, token)
		numeroDeAvisos, err := n.numeroDeNotificacionesEnviadasANegocio(&negocio)
		if err != nil {
			log.Printf("Error consultando número de avisos:%v", err)
			continue
		}
		log.Printf("Aviso %d de %d", numeroDeAvisos, NumeroAvisosAntesDeEliminarCuenta)

		_, err = n.enviarCorreoANegocioQueNoHaUsadoElSistema(&negocio, token, numeroDeAvisos)
		if err != nil {
			log.Printf("Error al enviar correo de notificación de eliminación: %v", err)
		}
	}
	log.Println("Terminado envío de correos a negocios antiguos que nunca han ingresado")
}

func (n *NegociosController) registrarNotificacionDeEliminacion(negocio *Negocio, token string) {
	bd, err := obtenerBDNegocios()
	if err != nil {
		log.Printf("Error obteniendo base de datos de negocios para guardar notificación de eliminación: %v", err)
		return
	}
	defer bd.Close()
	bd.Exec("insert into notificaciones_eliminacion_negocios(id_negocio, token) values (?,?)", negocio.Id, token)
}

func (n *NegociosController) eliminarNotificacionesDeEliminacionDeNegocio(idNegocio int) error {
	bd, err := obtenerBDNegocios()
	if err != nil {
		return err
	}
	defer bd.Close()
	_, err = bd.Exec("delete from notificaciones_eliminacion_negocios where id_negocio = ?", idNegocio)
	return err
}

func (n *NegociosController) registrarAcceso(idNegocio int, exitoso bool) error {

	db, err := obtenerBDNegocios()
	if err != nil {
		return err
	}
	defer db.Close()

	if err != nil {
		log.Printf("Error abriendo base de datos: %q", err)
		return err
	}

	_, err = db.Exec("INSERT INTO accesos_negocios (id_negocio,momento, exitoso) VALUES (?, ?, ?)", idNegocio, timestampParaMySQL(), exitoso)
	if exitoso {
		err = n.eliminarNotificacionesDeEliminacionDeNegocio(idNegocio)
		if err != nil {
			log.Printf("Error eliminando notificaciones de negocio: %v", err)
		}
	}
	return err
}

func (n *NegociosController) prepararDBMaestra() (bool, error) {
	db, err := obtenerBDNegocios()
	if err != nil {
		return false, err
	}
	defer db.Close()

	// Leer lo que hay en el archivo sql
	bytesLeidos, err := ioutil.ReadFile(NombreArchivoEsquemaSQLNegocios)
	if err != nil {
		log.Printf("Error leyendo archivo %s: %v", NombreArchivoEsquemaSQLNegocios, err)
	}

	tablas := strings.Split(string(bytesLeidos), ";")
	for _, tabla := range tablas {
		tabla = strings.Trim(tabla, "\n\r")
		if len(tabla) > 0 {
			_, err = db.Exec(tabla)
			if err != nil {
				log.Printf("Error preparando una tabla para BD negocios: %v", err)
				return false, nil
			}
		}
	}

	return true, nil
}

/*
  Una vez que el correo electrónico es verificado, se crea el usuario
  administrador con la contraseña hasheada del correo
*/
func (n *NegociosController) prepararPrimerUso(negocio *Negocio) {
	/*
		Solo es necesario descomentar en caso de que NO se esté usando SQLite3
	*/
	// dbMaestra, err := obtenerBDNegocios()
	// if err != nil {
	// 	log.Printf("Error obteniendo DB negocios en el primer uso: %v", err)
	// 	return
	// }
	// _, err = dbMaestra.Exec(fmt.Sprintf("create database if not exists %s", nombreBDAPartirDeNegocio(negocio.Id)))
	// if err != nil {
	// 	log.Printf("Error creando base de datos para nuevo negocio: %v", err)
	// 	return
	// }
	// dbMaestra.Close()
	/*

	 */
	db, err := obtenerBaseDeDatosAPartirDeIdNegocio(negocio.Id)
	if err != nil {
		log.Printf("Error obteniendo bd para negocio: %v", err)
	}
	defer db.Close()
	contextoBDUsuario, err := db.Begin()
	if err != nil {
		log.Printf("Error obteniendo contexto: %v", err)
	}
	// Crear las tablas...
	bytesLeidos, err := ioutil.ReadFile(NombreArchivoEsquemaSQLSistemas)
	if err != nil {
		log.Printf("Error leyendo archivo %s: %v", NombreArchivoEsquemaSQLSistemas, err)
	}
	tablas := strings.Split(string(bytesLeidos), ";")
	for _, tabla := range tablas {
		tabla = strings.Trim(tabla, "\n\r")
		if len(tabla) > 0 {
			_, err = contextoBDUsuario.Exec(tabla)
			if err != nil {
				log.Printf("Error creando una tabla para el negocio con nombre %s... %q\n. La tabla es: '%q'", negocio.Nombre, err, tabla)
				contextoBDUsuario.Rollback()
				return
			}
		}
	}

	// Más tarde, el usuario administrador
	sentencia, err := contextoBDUsuario.Prepare("INSERT INTO usuarios(nombre, contraseña) VALUES (?,?)")
	if err != nil {
		log.Printf("Error preparando: %v", err)
		return
	}
	_, err = sentencia.Exec(NombrePrimerUsuarioAdmin, negocio.Pass)
	if err != nil {
		log.Printf("Error creando admin para negocio con nombre %s y correo %s: %v", negocio.Nombre, negocio.Correo, err)
	}

	// Luego los permisos
	bytesLeidos, err = ioutil.ReadFile(NombreArchivoEsquemaSQLInit)
	if err != nil {
		log.Printf("Error leyendo archivo %s: %v", NombreArchivoEsquemaSQLInit, err)
	}

	permisos := strings.Split(string(bytesLeidos), ";")
	for _, permiso := range permisos {
		permiso = strings.Trim(permiso, "\n\r")
		if len(permiso) > 0 {
			_, err = contextoBDUsuario.Exec(permiso)
			if err != nil {
				log.Printf("Error creando una permiso para el negocio con nombre %s y correo %s... %q\n. El permiso es: '%q'", negocio.Nombre, negocio.Correo, err, permiso)
				contextoBDUsuario.Rollback()
				return
			}
		}
	}

	if err != nil {
		log.Printf("Error iniciando un valor: %v", err)
	}
	// y después, Darle los permisos...
	consultaPermisos := "INSERT INTO permisos_usuarios(idUsuario, idPermiso) SELECT 1, idPermiso FROM permisos;"

	if err != nil {
		log.Printf("Error al verificar negocio obteniendo contexto:\n%q", err)
		return
	}

	sentenciaPreparada, err := contextoBDUsuario.Prepare(consultaPermisos)
	if err != nil {
		log.Printf("Error preparando consulta:\n%q", err)
		return
	}

	_, err = sentenciaPreparada.Exec()
	if err != nil {
		log.Printf("Error copiando permisos: %v", err)
	}

	if err != nil {
		log.Printf("Error insertando permisos del primer usuario:\n%q", err)
		return
	}
	contextoBDUsuario.Commit()
}

func (n *NegociosController) pedirEnvioDeCorreo(detalles DetallesDeNegocio) (bool, error) {
	destinatarios := []string{detalles.Correo}
	return enviarCorreoParaVerificarCuenta(destinatarios, "Activa tu cuenta de Sublime POS 3", map[string]string{
		"urlCompartir":           UrlAcortadaCompartirSistema,
		"urlRedirigir":           fmt.Sprintf("%s/#/verificar/%s", UrlBaseApp, detalles.Token),
		"urlBasePaginaWeb":       UrlBasePaginaWeb,
		"correoSoporteYContacto": CorreoSoporteYContacto,
		"negocio":                detalles.Negocio,
		"correo":                 detalles.Correo,
		"mensajeCompartir":       MensajeCompartir,
	})

}

func (n *NegociosController) obtenerUnoPorToken(token string) (*Negocio, error) {
	db, err := obtenerBDNegocios()

	if err != nil {
		log.Printf("Error obteniendo BD negocios: %v", err)
		return nil, err
	}
	defer db.Close()
	filas, err := db.Query("SELECT id, nombre, correo, pass, verificado, token FROM negocios WHERe token = ? LIMIT 1;", token)
	if err != nil {
		log.Printf("Error al consultar negocio por token:\n%q", err)
		return nil, err
	}

	defer filas.Close()

	if !filas.Next() {
		return nil, errors.New("No existe un negocio con ese token")
	}
	var negocio Negocio
	err = filas.Scan(&negocio.Id, &negocio.Nombre, &negocio.Correo, &negocio.Pass, &negocio.Verificado, &negocio.Token)
	if err != nil {
		log.Printf("Error al escanear negocio por token:\n%q", err)
		return nil, err
	}
	return &negocio, nil
}

func (n *NegociosController) obtenerUnoPorCorreo(correo string) (*Negocio, error) {
	db, err := obtenerBDNegocios()

	if err != nil {
		log.Printf("Error obteniendo BD negocios: %v", err)
		return nil, err
	}
	defer db.Close()
	filas, err := db.Query("SELECT id, nombre, correo, pass, verificado, token FROM negocios WHERE correo = ? LIMIT 1;", correo)
	if err != nil {
		log.Printf("Error al consultar negocio por correo:\n%q", err)
		return nil, err
	}

	defer filas.Close()

	if !filas.Next() {
		return nil, errors.New("No existe un negocio con ese correo electrónico")
	}
	var negocio Negocio
	err = filas.Scan(&negocio.Id, &negocio.Nombre, &negocio.Correo, &negocio.Pass, &negocio.Verificado, &negocio.Token)
	if err != nil {
		log.Printf("Error al escanear negocio por correo:\n%q", err)
		return nil, err
	}
	return &negocio, nil
}

func (n *NegociosController) marcarComoVerificado(id int) (bool, error) {
	_, err := n.prepararDBMaestra()
	if err != nil {
		log.Printf("Error preparando BD de negocios: %v", err)
		return false, err
	}
	db, err := obtenerBDNegocios()
	if err != nil {
		return false, err
	}
	defer db.Close()
	consulta := "UPDATE negocios SET verificado = TRUE WHERE id = ?"

	sentenciaPreparada, err := db.Prepare(consulta)
	if err != nil {
		log.Printf("Error preparando sentencia para verificar negocio: %v", err)
		return false, err
	}
	_, err = sentenciaPreparada.Exec(id)
	if err != nil {
		log.Printf("Error ejecutando sentencia:\n%v", err)
		return false, nil
	}
	return true, nil
}

func (n *NegociosController) verificarPorToken(token string) bool {
	negocioNoVerificado, err := n.obtenerUnoPorToken(token)
	if err != nil {
		log.Printf("Error obteniendo negocio por token al verificar: %v", err)
		return false
	}
	n.prepararPrimerUso(negocioNoVerificado)
	n.marcarComoVerificado(negocioNoVerificado.Id)
	return true

}

func (n *NegociosController) existeTokenEnBaseDeDatosParaVerificaciones(token string) (bool, error) {
	db, err := obtenerBDNegocios()
	if err != nil {
		return true, err
	}
	defer db.Close()
	var numero uint8
	err = db.QueryRow("select 1 from negocios where token = ? LIMIT 1", token).Scan(&numero)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, err
}
func (n *NegociosController) existeTokenEnBaseDeDatosParaEliminaciones(token string) (bool, error) {
	db, err := obtenerBDNegocios()
	if err != nil {
		return true, err
	}
	defer db.Close()
	var numero uint8
	err = db.QueryRow("select 1 from notificaciones_eliminacion_negocios where token = ? LIMIT 1", token).Scan(&numero)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, err
}

// Obtiene un token aleatorio pero que no existe en la base de datos
func (n *NegociosController) obtenerTokenNoExistenteParaVerificarCuenta() (string, error) {
	for {
		log.Println("Dentro del ciclo para obtener un token no existente")
		token, err := generarTokenSeguro()
		if err != nil {
			return "", err
		}
		existe, err := n.existeTokenEnBaseDeDatosParaVerificaciones(token)
		if err != nil {
			return "", err
		}
		if err == nil && !existe {
			return token, nil
		}
	}
}

// Obtiene un token aleatorio pero que no existe en la base de datos de eliminaciones
func (n *NegociosController) obtenerTokenNoExistenteParaNotificarEliminacion() (string, error) {
	for {
		log.Println("Dentro del ciclo para obtener un token no existente para eliminar")
		token, err := generarTokenSeguro()
		if err != nil {
			return "", err
		}
		existe, err := n.existeTokenEnBaseDeDatosParaEliminaciones(token)
		if err != nil {
			return "", err
		}
		if err == nil && !existe {
			return token, nil
		}
	}
}

func (n *NegociosController) nuevo(negocio *Negocio) uint8 {
	_, err := n.prepararDBMaestra()
	if err != nil {
		log.Printf("Error preparando BD de negocios: %v", err)
		return RespuestaErrorRegistrandoNegocio
	}

	posibleNegocio, _ := n.obtenerUnoPorCorreo(negocio.Correo)
	if posibleNegocio != nil {
		log.Printf("Al intentar registrar un negocio llamado %s con correo %s ya existía", negocio.Nombre, negocio.Correo)
		return RespuestaErrorNegocioExistente
	}
	db, err := obtenerBDNegocios()
	if err != nil {
		return RespuestaErrorRegistrandoNegocio
	}
	defer db.Close()
	consulta := "INSERT INTO negocios (nombre, correo, pass, token, verificado, fecha_registro) VALUES (?, ?, ?, ?, ?, ?)"

	if err != nil {
		log.Printf("Error abriendo base de datos: %q", err)
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error obteniendo contexto de la base de datos:\n%q", err)
		return RespuestaErrorRegistrandoNegocio
	}

	sentenciaPreparada, err := tx.Prepare(consulta)

	if err != nil {
		log.Printf("Error preparando consulta:\n%q", err)
		tx.Rollback()
		return RespuestaErrorRegistrandoNegocio
	}

	token, err := n.obtenerTokenNoExistenteParaVerificarCuenta()
	if err != nil {
		log.Printf("Error obteniendo token no existente: %v", err)
		return RespuestaErrorRegistrandoNegocio
	}
	_, err = sentenciaPreparada.Exec(negocio.Nombre, negocio.Correo, hashearPassword(negocio.Pass), token, 0, timestampParaMySQL())
	if err != nil {
		log.Printf("Error ejecutando sentencia:\n%q", err)
		tx.Rollback()
		return RespuestaErrorRegistrandoNegocio
	}
	tx.Commit()
	n.verificarPorToken(token)

	return RespuestaNegocioRegistradoCorrectamente
}
