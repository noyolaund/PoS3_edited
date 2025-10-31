package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func configurarRutasDeGraficas(enrutador *mux.Router) {
	enrutador.HandleFunc("/total/vendido/por/mes/{anio}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		anio := variables["anio"]
		graficasController := DatosGraficasController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(graficasController.totalVentasPorAnio(anio), w, r)
	}).Name("VerGraficas").Methods(http.MethodGet)

	enrutador.HandleFunc("/total/vendido/por/dia/{anio}/{mes}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		anio, mes := variables["anio"], variables["mes"]
		graficasController := DatosGraficasController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(graficasController.totalVentasPorDiaEnMesYAnio(mes, anio), w, r)
	}).Name("VerGraficas").Methods(http.MethodGet)

	enrutador.HandleFunc("/anios/graficas/ventas/contado", func(w http.ResponseWriter, r *http.Request) {
		d := DatosGraficasController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(d.aniosEnLosQueHayRegistrosDeVentasAlContado(), w, r)
	}).Name("VerGraficas").Methods(http.MethodGet)

	enrutador.HandleFunc("/meses/graficas/ventas/contado/anio/{anio}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		anio, err := strconv.Atoi(variables["anio"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}

		d := DatosGraficasController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(d.mesesEnLosQueHayRegistrosDeVentasAlContadoDependiendoDeAnio(anio), w, r)
	}).Name("VerGraficas").Methods(http.MethodGet)

	enrutador.HandleFunc("/productos/mas/vendidos/{fechaInicio}/{fechaFin}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		fechaInicio, fechaFin := variables["fechaInicio"], variables["fechaFin"]
		graficasController := DatosGraficasController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(graficasController.productosMasVendidos(fechaInicio, fechaFin), w, r)
	}).Name("VerGraficas").Methods(http.MethodGet)

	enrutador.HandleFunc("/productos/nunca/vendidos/al/contado", func(w http.ResponseWriter, r *http.Request) {
		graficasController := DatosGraficasController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(graficasController.productosNuncaVendidosAlContado(), w, r)
	}).Name("VerGraficas").Methods(http.MethodGet)

	enrutador.HandleFunc("/productos/menos/vendidos/{fechaInicio}/{fechaFin}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		fechaInicio, fechaFin := variables["fechaInicio"], variables["fechaFin"]
		graficasController := DatosGraficasController{

			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(graficasController.productosMenosVendidos(fechaInicio, fechaFin), w, r)
	}).Name("VerGraficas").Methods(http.MethodGet)
}
