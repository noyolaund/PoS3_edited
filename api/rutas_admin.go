package main

import (
	"github.com/gorilla/mux"
)

func configurarRutasAdmin(enrutador *mux.Router) {
	subrouter := enrutador.PathPrefix("/auth").Subrouter()
	configurarRutasMisc(subrouter)
	configurarRutasPermisos(subrouter)
	configurarRutasDeUsuarios(subrouter)
	configurarRutasDeAjustes(subrouter)
	configurarRutasDeGraficas(subrouter)
	configurarRutasDeCaja(subrouter)
	configurarRutasDeVentasAlContado(subrouter)
	configurarRutasDeApartadosYAbonos(subrouter)
	configurarRutasDeClientes(subrouter)
	configurarRutasDeProductos(subrouter)
	// La sesión
	subrouter.Use(middlewareSesion)
}
