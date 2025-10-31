package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/michaeljs1990/sqlitestore"
)

var almacenamientoDeSesion *sqlitestore.SqliteStore

func init() {
	/*
		La definimos, así al asignar usamos =
		En lugar de usar :=
		Esto es por el scope; si usamos := entonces Go va a pensar que la variable
		almacenamientoDeSesion está dentro de esta función y la declarará como nueva,
		ignorando la global que está allá arriba
		Por eso es que debemos hacer que Go sepa que nos referimos a la global

	*/
	var err error
	//TODO: la clave cookie debe ser leída de env
	almacenamientoDeSesion, err = sqlitestore.NewSqliteStore(NombreBaseDeDatosSesiones, "sesiones", "/", 3600, []byte(ClaveCookie))
	if err != nil {
		log.Fatal(err)
	}
	almacenamientoDeSesion.Options = &sessions.Options{
		MaxAge:   EdadDeSesionEnSegundos,
		HttpOnly: true,
		Path:     "/",
	}

}
func main() {
	/*
	   Poner un log para, valga la redundancia, loguear
	   los movimientos ahí en lugar de la terminal. De esta
	   manera podremos ver los errores en un archivo, con fecha
	   hora y también fichero en donde se produjeron. Por ejemplo:
	   2018/07/06 13:25:03 AyudanteBaseDeDatos.go:312: Dato insertado correctamente
	*/
	t := time.Now()
	ficheroLog, err := os.OpenFile(
		fmt.Sprintf("%d-%02d.log", t.Year(), t.Month()),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ficheroLog.Close()
	log.SetOutput(ficheroLog)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	enrutador := mux.NewRouter()
	configurarRutas(enrutador)
	/*
	if hanModificadoArchivos() {
		return
	}
	*/
	// Preparar y encender servidor

	servidor := &http.Server{
		Handler: enrutador,
		Addr:    PuertoServidor,
		// Timeouts para evitar que el servidor se quede "colgado" por siempre
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Servidor iniciado en http://localhost%s", PuertoServidor)

	log.Fatal(servidor.ListenAndServe())

	// Ahora agendar el envío de correos
	//c := cron.New()
	//
	//c.AddFunc("0 35 9 * * 4", func() {
	//	nc := NegociosController{}
	//	// Primero se eliminan los que ya llevan las notificaciones acumuladas
	//	nc.eliminarNegociosQueNoHanUsadoElSistema()
	//	// Para que no se les tenga que enviar correo más tarde :)
	//	nc.enviarCorreoANegociosQueNoHanUsadoElSistema()
	//})
	//c.Start()
	//defer c.Stop()

	/*
	   Descomentar lo de abajo cuando se necesiten listar las rutas en el log
	*/

	//httpRequest := app.GetRoutes()
	//var rutasExistentes []string
	//for _, ruta := range httpRequest {
	//deberiaAgregar := true
	//for _, rutaExistente := range rutasExistentes {
	//  if rutaExistente == ruta.Name {
	//    deberiaAgregar = false
	//    break
	//  }
	//}
	//if deberiaAgregar {
	//  rutasExistentes = append(rutasExistentes, ruta.Name)
	//}
	//}
	//
	//for _, b := range rutasExistentes {
	//log.Printf("%s\n", b)
	//}
	// Y ejecutar la app!
}
