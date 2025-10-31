/*
  Cliente HTTP en Go con net/http
  Ejemplo de petición HTTP Get en Golang
  @author parzibyte
*/
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func detener(url string) {

	clienteHttp := &http.Client{
		Timeout: time.Second * 4, // Debe ser menor al "sleep" del desinstalador
	}
	// Si quieres agregar parámetros a la URL simplemente haz una
	// concatenación :)
	peticion, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// Maneja el error de acuerdo a tu situación
		log.Printf("Error creando petición: %v", err)
		return

	}
	// Podemos agregar encabezados
	respuesta, err := clienteHttp.Do(peticion)
	if err != nil {
		// Maneja el error de acuerdo a tu situación
		log.Printf("Error haciendo petición: %v", err)
		return
	}
	// No olvides cerrar el cuerpo al terminar
	defer respuesta.Body.Close()

	_, err = ioutil.ReadAll(respuesta.Body)
	if err != nil {
		log.Printf("Error leyendo respuesta: %v", err)
		return
	}
	println("Plugin detenido")

}

func main() {
	detener("http://localhost:2106/apagar")
	detener("http://localhost:8000/apagar")
}
