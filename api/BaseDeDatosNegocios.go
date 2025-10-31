package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func obtenerBDNegocios() (*sql.DB, error) {
	db, err := sql.Open(ControladorBD, CadenaConexionBDNegocios)
	if err != nil {
		log.Printf("Error abriendo base de datos maestra: %v", err)
		return nil, err
	}
	return db, nil
}
