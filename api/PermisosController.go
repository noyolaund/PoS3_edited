package main

import (
	"log"
)

type PermisosController struct {
	AjustesUsuario AjustesDeUsuarioLogueado
}

func (p *PermisosController) obtenerPorClave(clave string) Permiso {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		log.Panicf("Error abriendo la base de datos: %v\n", err)
		panic(err)
	}

	var permiso Permiso

	defer db.Close()
	filas, err := db.Query(`SELECT idPermiso, clave, descripcion FROM permisos WHERE clave = ? LIMIT 1`, clave)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener un permiso por clave:\n%q", err)
		return permiso
	}

	defer filas.Close()

	if filas.Next() {
		err := filas.Scan(&permiso.Id, &permiso.Clave, &permiso.Descripcion)
		if err != nil {
			log.Printf("Error al leer permiso por clave:\n%q", err)
		}
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al consultar permiso por clave:\n%q", err)
	}
	return permiso
}
func (p *PermisosController) deUnUsuario(idUsuario int) *[]int {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		log.Panicf("Error abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT idPermiso FROM permisos_usuarios WHERE idUsuario = ? ORDER BY idPermiso ASC;`, idUsuario)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener los idsPermisos de un usuario:\n%q", err)
		return nil
	}

	defer filas.Close()

	idsPermisos := []int{}
	for filas.Next() {
		var idPermiso int
		err := filas.Scan(&idPermiso)
		if err != nil {
			log.Printf("Error al leer permiso de un usuario:\n%q", err)
		}
		idsPermisos = append(idsPermisos, idPermiso)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todos los permisos de un usuario:\n%q", err)
	}
	return &idsPermisos
}
func (p *PermisosController) asignarA(idUsuario int, idsPermisos []int) bool {
	if idUsuario == 1 {
		//Previene la eliminación del usuario administrador
		return false
	}
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Printf("Error abriendo base de datos: %q", err)
		panic(err)
	}

	defer db.Close()
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error obteniendo contexto de base de datos:\n%q", err)
		return false
	}

	//Eliminar permisos...

	consulta := "DELETE FROM permisos_usuarios WHERE idUsuario = ?"

	sentenciaPreparada, err := tx.Prepare(consulta)
	if err != nil {
		log.Printf("Error preparando consulta para eliminar permisos:\n%q", err)
		return false
	}

	_, err = sentenciaPreparada.Exec(idUsuario)
	if err != nil {
		log.Printf("Error ejecutando sentencia para eliminar permisos:\n%q", err)
		tx.Rollback()
	}

	//Insertar nuevos

	consultaParaInsertar := "INSERT INTO permisos_usuarios(idUsuario, idPermiso) VALUES (?, ?)"
	sentenciaPreparadaParaInsertar, err := tx.Prepare(consultaParaInsertar)
	if err != nil {
		log.Printf("Error preparando consulta para insertar nuevos permisos:\n%q", err)
		tx.Rollback()
		return false
	}

	for _, idPermiso := range idsPermisos {
		sentenciaPreparadaParaInsertar.Exec(idUsuario, idPermiso)
	}

	//Y si las cosas salen bien, hacer un commit para "guardar" los cambios
	//que le hicimos a la base de datos
	tx.Commit()
	return true
}
func (p *PermisosController) todos() *[]Permiso {
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()
	if err != nil {
		log.Panicf("Error abriendo la base de datos: %v\n", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(`SELECT idPermiso, clave, descripcion FROM permisos;`)
	if err != nil {
		log.Printf("Error al realizar la consulta para obtener todos los permisos:\n%q", err)
		return nil
	}

	defer filas.Close()

	permisos := []Permiso{}
	for filas.Next() {
		var permiso Permiso
		err := filas.Scan(&permiso.Id, &permiso.Clave, &permiso.Descripcion)
		if err != nil {
			log.Printf("Error al obtener todos los permisos:\n%q", err)
		}
		permisos = append(permisos, permiso)
	}

	if err = filas.Err(); err != nil {
		log.Printf("Error en las filas al escanear todos los permisos:\n%q", err)
	}
	return &permisos
}
func (p *PermisosController) usuarioPuede(idUsuario int, clavePermiso string) bool {
	consulta := `
SELECT permisos_usuarios.idPermiso, permisos_usuarios.idUsuario, permisos.clave 
FROM permisos_usuarios
INNER JOIN permisos
ON permisos.idPermiso = permisos_usuarios.idPermiso
WHERE permisos.clave = ? AND permisos_usuarios.idUsuario = ?
LIMIT 1;
`
	db, err := p.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Panicf("Error abriendo la base de datos: %v\n", err)
		panic(err)
	}
	defer db.Close()
	filas, err := db.Query(consulta, clavePermiso, idUsuario)
	if err != nil {
		log.Printf("Error al consultar permiso:\n%q", err)
		return false
	}

	defer filas.Close()
	return filas.Next()
}

func (p *PermisosController) asignarTodosAlPrimerUsuario() {
	ayudanteUsuarios := AyudanteBaseDeDatos{
		nombreTabla:    "usuarios",
		AjustesUsuario: p.AjustesUsuario,
	}
	cantidadDeUsuarios := ayudanteUsuarios.conteo()
	if cantidadDeUsuarios == 1 {
		/*
		   Copiar todos los id de la tabla permisos.
		   @see https://parzibyte.me/blog/2017/06/01/clonar-estructura-tabla-copiar-datos-en-mysql/
		*/
		consultaPermisos := "INSERT INTO permisos_usuarios(idUsuario, idPermiso) SELECT 1, idPermiso FROM permisos;"

		db, err := p.AjustesUsuario.obtenerBaseDeDatos()

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

		sentenciaPreparadaVenta, err := tx.Prepare(consultaPermisos)
		if err != nil {
			log.Printf("Error preparando consulta:\n%q", err)
			return
		}

		_, err = sentenciaPreparadaVenta.Exec()

		if err != nil {
			log.Printf("Error insertando permisos del primer usuario:\n%q", err)
		}
		tx.Commit()
	} else {
		log.Printf("No se registraron los permisos porque el conteo de usuarios es %d\n", cantidadDeUsuarios)
	}
}
