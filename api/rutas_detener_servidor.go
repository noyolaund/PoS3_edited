package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

func configurarRutasApagarServidor(enrutador *mux.Router) {
	enrutador.HandleFunc("/apagar", func(w http.ResponseWriter, r *http.Request) {
		pid := os.Getpid()
		cmd := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid))
		err := cmd.Run()
		if err != nil {
			responderHttpConError(err, w, r)
		} else {
			responderHttpExitoso(true, w, r)
		}
	}).Methods(http.MethodGet).Name(RutaGeneralNoNecesitaComprobacion)
}
