/*
	Este archivo se encarga de guardar las funciones
	que no pertenecen a ningún lugar, pero que necesitamos
	en muchas ocasiones
*/

package main

import (
	"bytes"
	"math"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func hanModificadoArchivos() bool {
	// Si no estamos en producción, deshabilitamos el chequeo
	if !PRODUCCION {
		return false
	}
	for ubicacion, hashOriginal := range ubicacionesConHash {
		hash, _ := obtenerSha256DeArchivo(ubicacion)
		if hash != hashOriginal {
			log.Printf("El archivo '%s' ha sido modificado o no existe. Suma de verificación esperada: '%s', suma de verificación obtenida: '%s' El programa no puede iniciar", ubicacion, hashOriginal, hash)
			return true
		}
	}
	return false
}

func obtenerSha256DeArchivo(ubicacion string) (string, error) {
	f, err := os.Open(ubicacion)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func responderHttpExitoso(valor interface{}, w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(valor)
}
func responderHttpConError(err error, w http.ResponseWriter, r *http.Request) {
	log.Printf("Error al servir respuesta para %s: %v", r.RemoteAddr, err)
	json.NewEncoder(w).Encode(fmt.Sprintf("Error en el servidor: %v", err))

}
func comprobarError(e error) {
	if e != nil {
		panic(e)
	}
}
func obtenerRutaActual() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error obteniendo ruta actual: %v", err)
	}
	return dir
}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func obtenerIP() net.IP {
	if runtime.GOOS == "windows" {

		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		localAddr := conn.LocalAddr().(*net.UDPAddr)

		return localAddr.IP
	}
	return nil
}

func abrirAppEnNavegador() {
	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/C", "start", "http://localhost:"+PuertoServidor+"/static/index.html").Output()
	}
}
func obtenerTiempoAPartirDeCadena(fecha string) time.Time {
	justoAhora, _ := time.Parse("2006-01-02T15:04:05.000Z", fecha+".000Z")
	return justoAhora
}

/*
  Formatea una fecha en formato corto. Devuelve algo como:
  2 de enero del 2018
*/
func formatearFechaCortaParaUsuario(t *time.Time) string {

	return fmt.Sprintf("%02d de %s del %d",
		t.Day(), []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}[t.Month()-1], t.Year())
}
func formatearFechaParaUsuario(t *time.Time) string {
	hora := t.Hour()
	sufijo := "a.m."
	if hora > 12 {
		hora = hora - 12
		sufijo = "p.m."
	}
	return fmt.Sprintf("%02d de %s del %d %02d:%02d:%02d %s",
		t.Day(), []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}[t.Month()-1], t.Year(),
		hora, t.Minute(), t.Second(), sufijo)
}

func formatearFecha(t *time.Time) string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
func obtenerFechaActualFormateada() string {
	/*
		Regresa la fecha y hora actual en un formato que es
		entendible (más o menos) para el ser humano, el cual
		también es compatible para el ordenamiento lexicográficamente
		y de igual manera puede filtrarse en el lado del cliente

		Devuelve algo como 2018-05-29T10:34:22
	*/
	t := time.Now()
	return formatearFecha(&t)
}

func crearDirectorioSiNoExiste(directorio string) {
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.Mkdir(directorio, 0755)
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("El error es: %v", err)
	}
}

func obtenerDiaActual() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d",
		t.Year(), t.Month(), t.Day())
}
func fechaDeHoyAMedianoche() string {
	/*
		El día de hoy pero a medianoche. Algo como:
		2018-05-31T00:00:00
	*/
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT00:00:00",
		t.Year(), t.Month(), t.Day())
}

func fechaDeHoySinHora() string {
	return time.Now().Format("2006-01-02")
}

func obtenerListaDeImpresorasCompartidasWindows() []string {
	if runtime.GOOS == "windows" {

		salida, err := exec.Command("cmd", "/C", "wmic", "printer", "get", "name").Output()
		if err != nil {
			log.Fatal(err)
		}
		salidaCadena := string(salida)
		listaDeImpresoras := strings.Split(salidaCadena, "\r\r\n")
		var listaDeImpresorasLimpias []string
		for _, impresora := range listaDeImpresoras {
			nombreLimpio := strings.TrimRight(impresora, " ")
			if len(nombreLimpio) > 0 && nombreLimpio != "Name" {
				listaDeImpresorasLimpias = append(listaDeImpresorasLimpias, nombreLimpio)
			}
		}
		return listaDeImpresorasLimpias
	}
	return []string{}
}

func prepararPassPlana(passPlana string) string {
	passHasheada := sha512.Sum512_256([]byte(passPlana))
	hashSinNull := bytes.Trim(passHasheada[:], "\x00")
	return string(hashSinNull)
}

/***
* @see https://parzibyte.me/blog/2018/05/31/hasheando-comprobando-contrasenas-golang/
 */
func hashearPassword(passPlana string) string {
	if len(passPlana) <= 0 {
		log.Printf("Se ha hasheado una contraseña vacía")
	}
	//Al momento de escribir este código, bcrypt generaba un hash de 60 caracteres
	hash, err := bcrypt.GenerateFromPassword([]byte(prepararPassPlana(passPlana)), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

/***
* @see https://parzibyte.me/blog/2018/05/31/hasheando-comprobando-contrasenas-golang/
 */
func coincidenHashYPass(pass, hash string) bool {
	e := bcrypt.CompareHashAndPassword([]byte(hash), []byte(prepararPassPlana(pass)))
	return e == nil
}

func obtenerUUID() string {
	if runtime.GOOS == "windows" {

		salida, err := exec.Command("cmd", "/C", "wmic", "csproduct", "get", "UUID").Output()
		if err != nil {
			log.Fatal(err)
		}
		salidaCadena := string(salida)
		lineas := strings.Split(salidaCadena, "\r\r\n")
		uuid := strings.TrimRight(lineas[1], " ")
		return uuid
	}
	return "WE ARE ON LINUX"
}
func obtenerSerialDisco() string {
	if runtime.GOOS == "windows" {

		salida, err := exec.Command("cmd", "/C", "wmic", "DISKDRIVE", "get", "SerialNumber").Output()
		if err != nil {
			log.Fatal(err)
		}
		salidaCadena := string(salida)
		lineas := strings.Split(salidaCadena, "\r\r\n")
		serial := strings.TrimRight(lineas[1], " ")
		return serial
	}
	return "WE ARE ON LINUX"
}

func nombreBDAPartirDeNegocio(idNegocio int) string {
	return fmt.Sprintf("%s%d", PrefijoBDNegocios, idNegocio)
}

func generarTokenSeguro() (string, error) {
	var clave [32]byte
	_, err := io.ReadFull(rand.Reader, clave[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(clave[:]), nil
}

/*
  Regresa algo como 2019-22-01T18:28:00
*/
func timestampParaMySQL() string {
	//https://golang.org/src/time/format.go
	return time.Now().Format("2006-01-02T15:04:05")
}

// RoundToTwoDecimals redondea un float64 a 2 decimales.
// Se usa para almacenar valores parciales de existencia con 2 decimales
func RoundToTwoDecimals(v float64) float64 {
	return math.Round(v*100) / 100
}
