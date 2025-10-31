package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func configurarRutasDeLogin(enrutador *mux.Router) {

	enrutador.HandleFunc("/estoy/logueado", func(w http.ResponseWriter, r *http.Request) {
		s := Sesion{
			request:        r,
			responseWriter: w,
		}
		logueado, err := s.estaLogueado()
		if err != nil {
			responderHttpConError(err, w, r)
		} else {
			responderHttpExitoso(logueado, w, r)
		}
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodGet)

	enrutador.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		s := Sesion{
			request:        r,
			responseWriter: w,
		}
		err := s.cerrarSesion()
		if err == nil {
			responderHttpExitoso(true, w, r)
		} else {
			responderHttpConError(err, w, r)
		}
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodGet)

	enrutador.HandleFunc("/es/primer/uso", func(w http.ResponseWriter, r *http.Request) {
		uc := UsuariosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(uc.conteo() == 0, w, r)
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodGet)

	enrutador.HandleFunc("/usuario/login", func(w http.ResponseWriter, r *http.Request) {
		var usuario Usuario
		err := json.NewDecoder(r.Body).Decode(&usuario)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		uc := UsuariosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		logueado := uc.login(&usuario)
		if logueado == RespuestaLoginCorrecto {
			s := Sesion{
				request:        r,
				responseWriter: w,
			}
			err := s.propagarDatosDeUsuario(usuario)
			if err != nil {
				responderHttpConError(err, w, r)
			} else {
				responderHttpExitoso(logueado, w, r)
			}
		} else {
			responderHttpExitoso(logueado, w, r)
		}
	}).Name(RutaGeneralNoNecesitaComprobacion).Methods(http.MethodPut)
}
