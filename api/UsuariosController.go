package main

import (
	"log"
)

type UsuariosController struct {
	AjustesUsuario AjustesDeUsuarioLogueado
}

func (u *UsuariosController) uno(idUsuario int) *Usuario {
	db, err := u.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		log.Panicf("Error abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT idUsuario, nombre FROM usuarios WHERE idUsuario = ? LIMIT 1;`, idUsuario)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener un usuario por id:\n%q", err)
		return nil
	}

	defer filas.Close()

	var usuario Usuario
	if filas.Next() {
		err := filas.Scan(&usuario.Numero, &usuario.Nombre)
		if err != nil {
			log.Printf("Error al obtener un usuario por id:\n%q", err)
		}
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear un usuario por id:\n%q", err)
	}
	return &usuario
}

func (u *UsuariosController) todos() *[]Usuario {
	db, err := u.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		log.Panicf("Error abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT idUsuario, nombre FROM usuarios;`)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener todos los usuarios:\n%q", err)
		return nil
	}

	defer filas.Close()

	usuarios := []Usuario{}
	for filas.Next() {
		var usuario Usuario
		err := filas.Scan(&usuario.Numero, &usuario.Nombre)
		if err != nil {
			log.Printf("Error al obtener todos los usuarios:\n%q", err)
		}
		usuarios = append(usuarios, usuario)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todos los usuarios:\n%q", err)
	}
	return &usuarios
}

func (u *UsuariosController) eliminar(usuario *Usuario) {
	if usuario.Numero == 1 {
		//Previene la eliminación del usuario administrador
		return
	}
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "usuarios",
		AjustesUsuario: u.AjustesUsuario,
	}
	ayudante.eliminarDonde("idUsuario", usuario.Numero)
}

func (u *UsuariosController) conteo() int {
	db, err := u.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Panicf("Error abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query("SELECT COUNT(idUsuario) FROM usuarios;")
	if err != nil {
		log.Printf("Error al conteo usuarios:\n%q", err)
		return 0
	}

	defer filas.Close()

	if !filas.Next() {
		return 0
	}
	var total int
	err = filas.Scan(&total)
	if err != nil {
		log.Printf("Error al conteo usuarios:\n%q", err)
		return 0
	}
	return total
}

func (u *UsuariosController) nuevo(usuario *Usuario) bool {
	if u.existe(usuario) {
		return false
	}
	ayudante := AyudanteBaseDeDatos{
		nombreTabla:    "usuarios",
		AjustesUsuario: u.AjustesUsuario,
	}
	ayudante.insertar(map[string]interface{}{
		"nombre":     usuario.Nombre,
		"contraseña": hashearPassword(usuario.Password),
	})
	return true
}

func (u *UsuariosController) existe(usuario *Usuario) bool {
	db, err := u.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Panicf("Error abriendo la base de datos: %v\n", err)
		panic(err)
	}
	defer db.Close()
	filas, err := db.Query("SELECT nombre FROM usuarios WHERE nombre = ? LIMIT 1;", usuario.Nombre)
	if err != nil {
		log.Printf("Error al consultar usuario para ver si existe:\n%q", err)
		return false
	}

	defer filas.Close()
	return filas.Next()
}

func (u *UsuariosController) login(usuario *Usuario) uint8 {
	// Para evitar o ralentizar ataques de fuerza bruta
	//TODO: poner captcha
	//time.Sleep(time.Second * 3)
	nc := NegociosController{}
	negocio, err := nc.obtenerUnoPorCorreo(usuario.Negocio.Correo)
	if err != nil {
		log.Printf("Error obteniendo negocio: %v", err)
		return RespuestaLoginIncorrecto
	}
	if !negocio.Verificado {
		log.Printf("Negocio llamado %s con correo %s no verificado", negocio.Nombre, negocio.Correo)
		return RespuestaLoginNegocioNoVerificado
	}
	//Esta vez obtenemos directamente la BD, después se obtendrá de la sesión
	db, err := obtenerBaseDeDatosAPartirDeIdNegocio(negocio.Id)

	if err != nil {
		log.Panicf("Error abriendo la base de datos: %v\n", err)
		panic(err)
	}
	defer db.Close()
	filas, err := db.Query("SELECT idUsuario, contraseña FROM usuarios WHERE nombre = ? LIMIT 1;", usuario.Nombre)
	if err != nil {
		log.Printf("Error al consultar usuario para login:\n%q", err)
		nc.registrarAcceso(negocio.Id, false)
		return RespuestaLoginError
	}

	defer filas.Close()

	if !filas.Next() {
		log.Printf("Al consultar un usuario por nombre, no hay filas")
		nc.registrarAcceso(negocio.Id, false)
		return RespuestaLoginIncorrecto
	}

	var hash string
	err = filas.Scan(&usuario.Numero, &hash)
	if err != nil {
		log.Printf("Error al consultar usuario para login:\n%q", err)
		nc.registrarAcceso(negocio.Id, false)
		return RespuestaLoginError
	}
	// Se asignan las cosas para luego propagarlas
	usuario.Negocio = *negocio

	// Registrar login exitoso o no exitoso
	exitoso := coincidenHashYPass(usuario.Password, hash)
	nc.registrarAcceso(negocio.Id, exitoso)
	if exitoso {
		return RespuestaLoginCorrecto
	}
	return RespuestaLoginIncorrecto
}
