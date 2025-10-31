package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var rutasPermitidasSinSesion = []string{RutaGeneralNoNecesitaComprobacion}

func configurarRutas(enrutador *mux.Router) {
	if PRODUCCION {
		// Con esto servimos el contenido de la carpeta dist
		enrutador.PathPrefix(RutaDirectorioContenidoEstatico).Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(NombreDirectorioContenidoStaticoParaCliente))))
	}
	// Configuración de todas las rutas
	configurarRutasApagarServidor(enrutador)
	configurarRutasDeLogin(enrutador)
	configurarRutasDeNegocios(enrutador)
	configurarRutasUtiles(enrutador)
	configurarRutasMisc(enrutador)
	configurarRutasAdmin(enrutador)

	// CORS solo cuando estamos en desarrollo, pues JS se sirve en otro puerto
	if !PRODUCCION {
		// Necesitamos el método OPTIONS porque JS manda una petición antes de que se haga un DELETE
		enrutador.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", DominioPermitidoCORS)
		}).Methods("OPTIONS")

		enrutador.Use(middlewareCors)
	}
}
func rutaPermitidaSinIniciarSesion(ruta string) bool {
	// Si necesitan algo de static, no importa que no estén logueados pues es el simple archivo
	if strings.HasPrefix(ruta, RutaDirectorioContenidoEstatico) {
		return true
	}
	for _, r := range rutasPermitidasSinSesion {
		if ruta == r {
			return true
		}
	}
	return false
}

func middlewareSesion(siguienteManejador http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s := Sesion{
				request:        r,
				responseWriter: w,
			}
			usuarioEstaLogueado, err := s.estaLogueado()
			if err != nil {
				responderHttpConError(err, w, r)
				return
			}
			if !usuarioEstaLogueado {
				responderHttpConError(errors.New("usuario no ha iniciado sesión"), w, r)
				return
			}
			nombreDeRuta := mux.CurrentRoute(r).GetName()
			if nombreDeRuta == RutaGeneralNoNecesitaComprobacion {
				siguienteManejador.ServeHTTP(w, r)
				return
			}
			idUsuario, err := s.obtenerIdUsuario()
			if err != nil {
				responderHttpConError(err, w, r)
				return
			}

			pc := PermisosController{
				AjustesUsuario: AjustesDeUsuarioLogueado{
					httpResponseWriter: w,
					httpRequest:        r,
				},
			}
			if !pc.usuarioPuede(idUsuario, nombreDeRuta) {
				errorDePermiso := ErrorDePermiso{
					Clave:   nombreDeRuta,
					Mensaje: "Permiso denegado",
					Numero:  21,
					Permiso: pc.obtenerPorClave(nombreDeRuta),
				}
				responderHttpExitoso(errorDePermiso, w, r)
				return
			}
			// Si llegamos hasta aquí el usuario tiene permiso
			siguienteManejador.ServeHTTP(w, r)
		})
}

// Y también necesitamos un middleware para CORS, para todas aquellas que no sean OPTIONS
func middlewareCors(siguienteManejador http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", DominioPermitidoCORS)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			siguienteManejador.ServeHTTP(w, req)
		})
}
