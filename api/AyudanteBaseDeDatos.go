package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type AyudanteBaseDeDatos struct {
	nombreTabla    string
	AjustesUsuario AjustesDeUsuarioLogueado
}

func obtenerBaseDeDatosAPartirDeIdNegocio(idNegocio int) (*sql.DB, error) {
	// Debe tener la forma usuario:contraseña@host/nombreBaseDeDatos
	db, err := sql.Open(ControladorBD, CadenaConexionBDUsuarios+nombreBDAPartirDeNegocio(idNegocio)+".db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (ayudante *AyudanteBaseDeDatos) conteo() int {
	db, err := ayudante.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Fatalf("Error abriendo base de datos: %q", err)
		panic(err)
	}

	defer db.Close()
	filas, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM %s;", ayudante.nombreTabla))
	if err != nil {
		log.Printf("Error al contar datos de la tabla %s:\n%q", ayudante.nombreTabla, err)
		return 0
	}

	defer filas.Close()

	if !filas.Next() {
		return 0
	}
	var total int
	err = filas.Scan(&total)
	if err != nil {
		log.Printf("Error al contar datos de la tabla %s:\n%q", ayudante.nombreTabla, err)
		return 0
	}
	return total
}

func (ayudante *AyudanteBaseDeDatos) actualizarDonde(columna string, valor interface{}, columnasYValores map[string]interface{}) {
	var valores []interface{}
	columnasActualizadas := ""
	contador := 0
	longitud := len(columnasYValores)
	for clave, valor := range columnasYValores {
		valores = append(valores, valor)
		columnasActualizadas += fmt.Sprintf("%s = ?", clave)
		contador = contador + 1
		if contador < longitud {
			columnasActualizadas += ", "
		}
	}
	valores = append(valores, valor)
	consulta := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", ayudante.nombreTabla, columnasActualizadas, columna)

	db, err := ayudante.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Printf("Error abriendo base de datos: %q", err)
		panic(err)
	}

	defer db.Close()
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error obteniendo contexto:\n%q", err)
		return
	}

	sentenciaPreparada, err := tx.Prepare(consulta)
	if err != nil {
		log.Printf("Error preparando sentencia:\n%q", err)
		return
	}

	_, err = sentenciaPreparada.Exec(valores...)
	if err != nil {
		log.Printf("Error ejecutando sentencia:\n%q", err)
	}
	tx.Commit()
}
func (ayudante *AyudanteBaseDeDatos) eliminarDonde(columna string, valor interface{}) {
	db, err := ayudante.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Printf("Error abriendo base de datos: %q", err)
		panic(err)
	}

	defer db.Close()
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error obteniendo contexto de base de datos:\n%q", err)
		return
	}

	consulta := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", ayudante.nombreTabla, columna)

	sentenciaPreparada, err := tx.Prepare(consulta)
	if err != nil {
		log.Printf("Error preparando consulta:\n%q", err)
		return
	}

	_, err = sentenciaPreparada.Exec(valor)
	if err != nil {
		log.Printf("Error ejecutando sentencia:\n%q", err)
	}
	tx.Commit()
}

func (ayudante *AyudanteBaseDeDatos) insertar(columnasYValores map[string]interface{}) {
	var valores []interface{}
	columnas := ""
	signosDeInterrogacion := ""
	contador := 0
	longitud := len(columnasYValores)
	for clave, valor := range columnasYValores {
		valores = append(valores, valor)
		columnas += clave
		signosDeInterrogacion += "?"
		contador = contador + 1
		if contador < longitud {
			columnas += ", "
			signosDeInterrogacion += ", "
		}
	}
	consulta := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", ayudante.nombreTabla, columnas, signosDeInterrogacion)

	db, err := ayudante.AjustesUsuario.obtenerBaseDeDatos()

	if err != nil {
		log.Printf("Error abriendo base de datos: %q", err)
		panic(err)
	}

	defer db.Close()
	tx, err := db.Begin()

	if err != nil {
		log.Printf("Error obteniendo contexto de la base de datos:\n%q", err)
		return
	}

	sentenciaPreparada, err := tx.Prepare(consulta)
	if err != nil {
		log.Printf("Error preparando consulta:\n%q", err)
		return
	}

	_, err = sentenciaPreparada.Exec(valores...)
	if err != nil {
		log.Printf("Error ejecutando sentencia:\n%q", err)
	}
	tx.Commit()
}
