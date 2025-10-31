package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func configurarRutasDeClientes(enrutador *mux.Router) {
	enrutador.HandleFunc("/clientes", func(w http.ResponseWriter, r *http.Request) {
		cc := ClientesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(cc.todos(), w, r)
	}).Name("VerClientes").Methods(http.MethodGet)

	enrutador.HandleFunc("/cliente/{idCliente}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idCliente, err := strconv.Atoi(variables["idCliente"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		cc := ClientesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(cc.porRowid(idCliente), w, r)
	}).Name("VerClientes").Methods(http.MethodGet)

	enrutador.HandleFunc("/historial/cliente/{idCliente}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idCliente, err := strconv.Atoi(variables["idCliente"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		cc := ClientesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(cc.historial(idCliente), w, r)
	}).Name("VerClientes").Methods(http.MethodGet)

	enrutador.HandleFunc("/buscar/clientes/{nombre}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		cc := ClientesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(cc.buscar(variables["nombre"]), w, r)
	}).Name("VerClientes").Methods(http.MethodGet)

	enrutador.HandleFunc("/autocompletado/clientes/{nombre}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		cc := ClientesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(cc.buscarParaAutocompletado(variables["nombre"]), w, r)
	}).Name("AutocompletarClientes").Methods(http.MethodGet)

	enrutador.HandleFunc("/cliente", func(w http.ResponseWriter, r *http.Request) {
		var cliente Cliente
		err := json.NewDecoder(r.Body).Decode(&cliente)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		cc := ClientesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		cc.nuevo(&cliente)
		responderHttpExitoso(cliente, w, r)
	}).Name("RegistrarCliente").Methods(http.MethodPost)

	enrutador.HandleFunc("/cliente", func(w http.ResponseWriter, r *http.Request) {
		var cliente Cliente
		err := json.NewDecoder(r.Body).Decode(&cliente)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		cc := ClientesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		cc.actualizar(&cliente)
		responderHttpExitoso(cliente, w, r)
	}).Name("ActualizarCliente").Methods(http.MethodPut)

	enrutador.HandleFunc("/cliente/{idCliente}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idCliente, err := strconv.Atoi(variables["idCliente"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		cc := ClientesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		clienteParaEliminar := Cliente{
			Numero: idCliente,
		}
		cc.eliminar(&clienteParaEliminar)
		responderHttpExitoso(true, w, r)
	}).Name("EliminarCliente").Methods(http.MethodDelete)
}
