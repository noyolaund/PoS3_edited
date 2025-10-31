package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

func enviarCorreoParaVerificarCuenta(destinatario []string, asunto string, datos interface{}) (bool, error) {
	t, err := template.ParseFiles("plantillas_correos/verificar_cuenta.html")
	if err != nil {
		return false, err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, datos); err != nil {
		return false, err
	}
	contenidoPlantilla := buffer.String()
	cuerpo := "To: " + destinatario[0] + "\r\nSubject: " + asunto + "\r\n" + MIME + "\r\n" + contenidoPlantilla
	direccionServidor := fmt.Sprintf("%s:%s", GmailServidor, GmailPuerto)
	err = smtp.SendMail(direccionServidor, smtp.PlainAuth("", GmailCorreo, GmailPass, GmailServidor), GmailCorreo, destinatario, []byte(cuerpo))
	if err != nil {
		return false, err
	}
	return true, nil
}
