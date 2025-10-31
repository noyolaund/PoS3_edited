package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func configurarRutasDeApartadosYAbonos(enrutador *mux.Router) {
	enrutador.HandleFunc("/fecha/apartado/{idApartado}", func(w http.ResponseWriter, r *http.Request) {
		var nuevaFecha string
		err := json.NewDecoder(r.Body).Decode(&nuevaFecha)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		variables := mux.Vars(r)
		idApartado, err := strconv.Atoi(variables["idApartado"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		ac.cambiarFechaDeVencimiento(nuevaFecha, idApartado)
		responderHttpExitoso(true, w, r)
	}).Name("CambiarFechaVencimientoApartado").Methods(http.MethodPut)

	enrutador.HandleFunc("/apartado", func(w http.ResponseWriter, r *http.Request) {
		var apartado Apartado
		err := json.NewDecoder(r.Body).Decode(&apartado)
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
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		apartado.Usuario.Numero = idUsuario
		ac.nuevo(&apartado)
		responderHttpExitoso(apartado, w, r)
	}).Name("RegistrarApartado").Methods(http.MethodPost)

	enrutador.HandleFunc("/productos/apartado/{idApartado}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idApartado, err := strconv.Atoi(variables["idApartado"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		productosController := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(productosController.deUnApartado(idApartado), w, r)
	}).Name("VerApartados").Methods(http.MethodGet)

	enrutador.HandleFunc("/apartado/{idApartado}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idApartado, err := strconv.Atoi(variables["idApartado"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.uno(idApartado), w, r)
	}).Name("VerApartados").Methods(http.MethodGet)

	enrutador.HandleFunc("/abono/{idAbono}/{idApartado}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idApartado, err := strconv.Atoi(variables["idApartado"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		idAbono, err := strconv.Atoi(variables["idAbono"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.detallesDeUnAbono(idAbono, idApartado), w, r)
	}).Name("VerApartados").Methods(http.MethodGet)

	enrutador.HandleFunc("/cambiar/producto/apartado/{idApartado}/{idProductoAnterior}/{idProductoNuevo}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idApartado, err := strconv.Atoi(variables["idApartado"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		idProductoAnterior, err := strconv.Atoi(variables["idProductoAnterior"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		idProductoNuevo, err := strconv.Atoi(variables["idProductoNuevo"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.cambiarProducto(idApartado, idProductoAnterior, idProductoNuevo), w, r)
	}).Name("CambiarProductoDeApartado").Methods(http.MethodGet)

	enrutador.HandleFunc("/apartados/{fechaInicio}/{fechaFin}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		fechaInicio, fechaFin := variables["fechaInicio"], variables["fechaFin"]
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.enPeriodo(fechaInicio, fechaFin), w, r)
	}).Name("VerApartados").Methods(http.MethodGet)

	enrutador.HandleFunc("/total/abonado/{fechaInicio}/{fechaFin}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		fechaInicio, fechaFin := variables["fechaInicio"], variables["fechaFin"]
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.totalAbonadoEnPeriodo(fechaInicio, fechaFin), w, r)
	}).Name("VerApartados").Methods(http.MethodGet)

	enrutador.HandleFunc("/apartados/pendientes", func(w http.ResponseWriter, r *http.Request) {
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.pendientes(), w, r)
	}).Name("VerApartados").Methods(http.MethodGet)

	enrutador.HandleFunc("/apartados/proximos/vencer", func(w http.ResponseWriter, r *http.Request) {
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.proximosAVencer(DiasParaApartadosProximosAVencer), w, r)
	}).Name("VerApartados").Methods(http.MethodGet)

	enrutador.HandleFunc("/abonos/apartado/{idApartado}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idApartado, err := strconv.Atoi(variables["idApartado"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(ac.abonosDe(idApartado), w, r)
	}).Name("VerApartados").Methods(http.MethodGet)

	enrutador.HandleFunc("/abono", func(w http.ResponseWriter, r *http.Request) {
		var abono Abono
		err := json.NewDecoder(r.Body).Decode(&abono)
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
		abono.Usuario.Numero = idUsuario
		ac := ApartadosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		idAbono := ac.abonar(&abono)
		responderHttpExitoso(idAbono, w, r)
	}).Name("RegistrarApartado").Methods(http.MethodPut)
}
