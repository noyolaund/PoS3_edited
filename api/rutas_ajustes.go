package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func configurarRutasDeAjustes(enrutador *mux.Router) {
	enrutador.HandleFunc("/valor/{clave}", func(w http.ResponseWriter, r *http.Request) {
		ac := AjustesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		variables := mux.Vars(r)
		clave := variables["clave"]
		responderHttpExitoso(ac.obtenerValor(clave), w, r)
	}).Name("VerAjustes").Methods(http.MethodGet)

	enrutador.HandleFunc("/valor", func(w http.ResponseWriter, r *http.Request) {
		var claveYValor ClaveYValor
		err := json.NewDecoder(r.Body).Decode(&claveYValor)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ac := AjustesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		ac.guardarValor(claveYValor.Valor, claveYValor.Clave)
		responderHttpExitoso(true, w, r)
	}).Name("CambiarAjustes").Methods(http.MethodPut)
	enrutador.HandleFunc("/es/version/prueba", func(w http.ResponseWriter, r *http.Request) {
		responderHttpExitoso(false, w, r)
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodGet)

	enrutador.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
		responderHttpExitoso(obtenerIP(), w, r)
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodGet)

	enrutadorAjustes := enrutador.PathPrefix("/ajustes").Subrouter()
	enrutadorAjustes.HandleFunc("/otros", func(w http.ResponseWriter, r *http.Request) {
		ac := AjustesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.obtenerOtros(), w, r)
	}).Name("VerAjustes").Methods(http.MethodGet)

	enrutadorAjustes.HandleFunc("/otros", func(w http.ResponseWriter, r *http.Request) {
		var ajustes OtrosAjustes
		err := json.NewDecoder(r.Body).Decode(&ajustes)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ac := AjustesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		ac.guardarOtros(&ajustes)
		responderHttpExitoso(true, w, r)
	}).Name("CambiarAjustes").Methods(http.MethodPut)

	enrutadorAjustes.HandleFunc("/empresa", func(w http.ResponseWriter, r *http.Request) {
		ac := AjustesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.obtenerDatosEmpresa(), w, r)
	}).Name("VerAjustes").Methods(http.MethodGet)

	enrutadorAjustes.HandleFunc("/empresa", func(w http.ResponseWriter, r *http.Request) {
		var datosEmpresa DatosEmpresa
		err := json.NewDecoder(r.Body).Decode(&datosEmpresa)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ac := AjustesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		ac.guardarDatosEmpresa(&datosEmpresa)
		responderHttpExitoso(true, w, r)
	}).Name("CambiarAjustes").Methods(http.MethodPut)

	enrutador.HandleFunc("/nombre/impresora", func(w http.ResponseWriter, r *http.Request) {
		ac := AjustesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.obtenerNombreImpresora(), w, r)
	}).Name("VerAjustes").Methods(http.MethodGet)

	enrutador.HandleFunc("/nombre/impresora", func(w http.ResponseWriter, r *http.Request) {
		var nombreImpresora string
		err := json.NewDecoder(r.Body).Decode(&nombreImpresora)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ac := AjustesController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		ac.guardarNombreImpresora(nombreImpresora)
		responderHttpExitoso(true, w, r)
	}).Name("CambiarAjustes").Methods(http.MethodPut)

	enrutador.HandleFunc("/probar/impresora/{nombreImpresora}", func(w http.ResponseWriter, r *http.Request) {
		responderHttpExitoso(false, w, r)
	}).Name("CambiarAjustes").Methods(http.MethodGet)
}
