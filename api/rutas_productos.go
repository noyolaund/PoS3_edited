package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func configurarRutasDeProductos(enrutador *mux.Router) {
	enrutador.HandleFunc("/reporte/inventario", func(w http.ResponseWriter, r *http.Request) {
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.reporteInventario(), w, r)
	}).Name("VerReporteDeInventario").Methods(http.MethodGet)

	enrutador.HandleFunc("/productos/aleatorios/{cantidad}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		cantidad, err := strconv.Atoi(variables["cantidad"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.alAzar(cantidad), w, r)
	}).Name("VerProductos").Methods(http.MethodGet)

	enrutador.HandleFunc("/buscar/productos/autocompletado/{busqueda}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.buscarParaAutocompletado(variables["busqueda"]), w, r)
	}).Name("AutocompletarProductos").Methods(http.MethodGet)

	enrutador.HandleFunc("/productos/{offset}/{limite}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		offset, err := strconv.Atoi(variables["offset"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		limite, err := strconv.Atoi(variables["limite"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.todos(offset, limite), w, r)
	}).Name("VerProductos").Methods(http.MethodGet)

	enrutador.HandleFunc("/conteo/productos/stock", func(w http.ResponseWriter, r *http.Request) {
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.conteoParaStock(), w, r)
	}).Name("VerProductos").Methods(http.MethodGet)

	enrutador.HandleFunc("/productos/stock/{offset}/{limite}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		offset, err := strconv.Atoi(variables["offset"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		limite, err := strconv.Atoi(variables["limite"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.enStock(offset, limite), w, r)
	}).Name("VerProductos").Methods(http.MethodGet)

	enrutador.HandleFunc("/producto/{idProducto}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idProducto, err := strconv.Atoi(variables["idProducto"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.porRowid(idProducto), w, r)
	}).Name("VerProductoPorCodigoONumero").Methods(http.MethodGet)

	enrutador.HandleFunc("/producto/codigo/barras/{codigo}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.porCodigoDeBarras(variables["codigo"]), w, r)
	}).Name("VerProductoPorCodigoONumero").Methods(http.MethodGet)

	enrutador.HandleFunc("/producto/codigoBarras/{codigo}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		responderHttpExitoso(pc.porCodigoDeBarras(variables["codigo"]), w, r)
	}).Name("VerProductoPorCodigoONumero").Methods(http.MethodGet)

	enrutador.HandleFunc("/siguiente/numero/producto", func(w http.ResponseWriter, r *http.Request) {
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		siguienteNumero, err := pc.siguienteNumero()
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		responderHttpExitoso(siguienteNumero, w, r)
	}).Name("VerProductos").Methods(http.MethodGet)

	enrutador.HandleFunc("/buscar/productos/{offset}/{limite}/{descripcion}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		offset, err := strconv.Atoi(variables["offset"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		limite, err := strconv.Atoi(variables["limite"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}

		responderHttpExitoso(pc.buscar(offset, limite, variables["descripcion"]), w, r)
	}).Name("VerProductos").Methods(http.MethodGet)

	enrutador.HandleFunc("/producto", func(w http.ResponseWriter, r *http.Request) {
		var producto Producto
		err := json.NewDecoder(r.Body).Decode(&producto)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		pc.nuevo(&producto)
		responderHttpExitoso(true, w, r)
	}).Name("RegistrarProducto").Methods(http.MethodPost)

	enrutador.HandleFunc("/producto", func(w http.ResponseWriter, r *http.Request) {
		var producto Producto
		err := json.NewDecoder(r.Body).Decode(&producto)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		pc.actualizar(&producto)
		responderHttpExitoso(true, w, r)
	}).Name("ActualizarProducto").Methods(http.MethodPut)

	enrutador.HandleFunc("/producto/{idProducto}", func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		idProducto, err := strconv.Atoi(variables["idProducto"])
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		producto := Producto{
			Numero: idProducto,
		}
		pc.eliminar(&producto)
		responderHttpExitoso(true, w, r)
	}).Name("EliminarProducto").Methods(http.MethodDelete)

	enrutador.HandleFunc("/importar/excel", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(TamanioMaximoArchivoImportacionExcel)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		encabezadosArchivo := r.MultipartForm.File["archivo"]
		if encabezadosArchivo == nil || len(encabezadosArchivo) < 0 {
			responderHttpConError(errors.New("no hay archivo"), w, r)
			return
		}
		encabezadoArchivo := encabezadosArchivo[0]
		archivo, err := encabezadoArchivo.Open()
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		indiceCodigoBarras, err := strconv.Atoi(r.FormValue("IndiceCodigoBarras"))
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		indiceDescripcion, err := strconv.Atoi(r.FormValue("IndiceDescripcion"))
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		indicePrecioCompra, err := strconv.Atoi(r.FormValue("IndicePrecioCompra"))
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		indicePrecioVenta, err := strconv.Atoi(r.FormValue("IndicePrecioVenta"))
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		indiceExistencia, err := strconv.Atoi(r.FormValue("IndiceExistencia"))
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		indiceStock, err := strconv.Atoi(r.FormValue("IndiceStock"))
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		ajustes := AjustesParaImportarProductos{
			TieneEncabezados:                r.FormValue("TieneEncabezados") == "true",
			IgnorarCodigosDeBarrasRepetidos: r.FormValue("IgnorarCodigosDeBarrasRepetidos") == "true",
			IndiceCodigoBarras:              indiceCodigoBarras,
			IndiceDescripcion:               indiceDescripcion,
			IndicePrecioCompra:              indicePrecioCompra,
			IndicePrecioVenta:               indicePrecioVenta,
			IndiceExistencia:                indiceExistencia,
			IndiceStock:                     indiceStock,
		}
		responderHttpExitoso(pc.importarExcel(archivo, encabezadoArchivo, &ajustes), w, r)
	}).Name("RegistrarProducto").Methods(http.MethodPost)

	enrutador.HandleFunc("/exportar", func(w http.ResponseWriter, r *http.Request) {
		var ajustes AjustesParaExportarProductos
		err := json.NewDecoder(r.Body).Decode(&ajustes)
		if err != nil {
			responderHttpConError(err, w, r)
			return
		}
		pc := ProductosController{
			AjustesUsuario: AjustesDeUsuarioLogueado{
				httpResponseWriter: w,
				httpRequest:        r,
			},
		}
		pc.exportar(ajustes)
		responderHttpExitoso(true, w, r)
	}).Name("VerProductos").Methods(http.MethodPut)
}
