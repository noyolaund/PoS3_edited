package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func configurarRutasDeVentasAlContado(enrutador *mux.Router) {
	enrutador.HandleFunc("/venta/contado/{idVenta}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idVenta, err := strconv.Atoi(variables["idVenta"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		vc := VentasContadoController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		vc.anular(idVenta)
		responderHttpExitoso(true, w, r)
	}).Name("AnularVenta").Methods(http.MethodDelete)
	enrutador.HandleFunc("/venta/contado/{idVenta}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idVenta, err := strconv.Atoi(variables["idVenta"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		vc := VentasContadoController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(vc.una(idVenta), w, r)
	}).Name("VerVentasContado").Methods(http.MethodGet)

	enrutador.HandleFunc("/venta/contado", func(w http.ResponseWriter, r *http.Request) {
		var ventaContado VentaContado
		err := json.NewDecoder(r.Body).Decode(&ventaContado)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		s := Sesion{
			request:        r,
			responseWriter: w,
		}
		idUsuario, err := s.obtenerIdUsuario()
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		vc := VentasContadoController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		ventaContado.Usuario.Numero = idUsuario
		vc.nueva(&ventaContado)
		responderHttpExitoso(&ventaContado, w, r)
	}).Name("RegistrarVentaContado").Methods(http.MethodPost)

	enrutador.HandleFunc("/ventas/contado/{fechaInicio}/{fechaFin}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		fechaInicio, fechaFin := variables["fechaInicio"], variables["fechaFin"]
		vc := VentasContadoController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(vc.enPeriodo(fechaInicio, fechaFin), w, r)
	}).Name("VerVentasContado").Methods(http.MethodGet)

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
