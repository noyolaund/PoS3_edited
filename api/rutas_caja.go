package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func configurarRutasDeCaja(enrutador *mux.Router) {
	enrutador.HandleFunc("/reporte/caja/{fechaInicio}/{fechaFin}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		fechaInicio, fechaFin := variables["fechaInicio"], variables["fechaFin"]
		cc := CajaController{AjustesUsuario: AjustesDeUsuarioLogueado{
			httpResponseWriter: w,
			httpRequest:        r,
		}}
		responderHttpExitoso(cc.paraReporte(fechaInicio, fechaFin), w, r)
	}).Name("VerReporteCaja").Methods(http.MethodGet)

	enrutador.HandleFunc("/reporte/caja/{fechaInicio}/{fechaFin}/{idUsuario}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		fechaInicio, fechaFin := variables["fechaInicio"], variables["fechaFin"]
		idUsuario, err := strconv.Atoi(variables["idUsuario"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		cc := CajaController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			}}
		responderHttpExitoso(cc.paraReportePorUsuario(fechaInicio, fechaFin, idUsuario), w, r)
	}).Name("VerReporteCaja").Methods(http.MethodGet)

	enrutador.HandleFunc("/ingreso", func(w http.ResponseWriter, r *http.Request) {
		var ingreso Ingreso
		err := json.NewDecoder(r.Body).Decode(&ingreso)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		s := Sesion{
			request:        r,
			responseWriter: w,
		}
		idUsuarioLogueado, err := s.obtenerIdUsuario()
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		cc := CajaController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			}}

		ingreso.Usuario.Numero = idUsuarioLogueado
		cc.nuevoIngreso(&ingreso)
		responderHttpExitoso(true, w, r)
	}).Name("RegistrarIngreso").Methods(http.MethodPost)

	enrutador.HandleFunc("/egreso", func(w http.ResponseWriter, r *http.Request) {
		var egreso Egreso
		err := json.NewDecoder(r.Body).Decode(&egreso)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		s := Sesion{
			request:        r,
			responseWriter: w,
		}
		idUsuarioLogueado, err := s.obtenerIdUsuario()
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		cc := CajaController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			}}

		egreso.Usuario.Numero = idUsuarioLogueado
		cc.nuevoEgreso(&egreso)
		responderHttpExitoso(true, w, r)
	}).Name("RegistrarEgreso").Methods(http.MethodPost)

	enrutador.HandleFunc("/ingresos/{fechaInicio}/{fechaFin}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		fechaInicio, fechaFin := variables["fechaInicio"], variables["fechaFin"]
		cc := CajaController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(cc.ingresosEnPeriodo(fechaInicio, fechaFin), w, r)
	}).Name("VerReporteCaja").Methods(http.MethodGet)

	enrutador.HandleFunc("/egresos/{fechaInicio}/{fechaFin}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		fechaInicio, fechaFin := variables["fechaInicio"], variables["fechaFin"]
		cc := CajaController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(cc.egresosEnPeriodo(fechaInicio, fechaFin), w, r)
	}).Name("VerReporteCaja").Methods(http.MethodGet)
}
