package main

import (
	"net/http"
)

type Sesion struct {
	request        *http.Request
	responseWriter http.ResponseWriter
}

func (s *Sesion) propagarDatosDeUsuario(usuario Usuario) error {
	sesion, err := almacenamientoDeSesion.Get(s.request, NombreSesion)
	if err != nil {
		return err
	}
	sesion.Values["autenticado"] = true
	sesion.Values["idUsuario"] = usuario.Numero
	sesion.Values["idNegocio"] = usuario.Negocio.Id
	sesion.Values["correoNegocio"] = usuario.Negocio.Correo
	sesion.Values["nombreNegocio"] = usuario.Negocio.Nombre
	sesion.Values["nombreUsuario"] = usuario.Nombre
	return sesion.Save(s.request, s.responseWriter)
}

func (s *Sesion) obtenerIdUsuario() (int, error) {
	sesion, err := almacenamientoDeSesion.Get(s.request, NombreSesion)
	if err != nil {
		return -1, err
	}
	idUsuario := sesion.Values["idUsuario"]
	if idUsuario == nil {
		return -1, nil
	}
	return idUsuario.(int), nil
}

func (s *Sesion) obtenerUsuarioLogueado() (Usuario, error) {
	var usuario Usuario
	sesion, err := almacenamientoDeSesion.Get(s.request, NombreSesion)
	if err != nil {
		return usuario, err
	}
	usuario.Numero = sesion.Values["idUsuario"].(int)
	usuario.Nombre = sesion.Values["nombreUsuario"].(string)
	usuario.Negocio.Id = sesion.Values["idNegocio"].(int)
	usuario.Negocio.Correo = sesion.Values["correoNegocio"].(string)
	usuario.Negocio.Nombre = sesion.Values["nombreNegocio"].(string)
	return usuario, nil
}

func (s *Sesion) estaLogueado() (bool, error) {
	sesion, err := almacenamientoDeSesion.Get(s.request, NombreSesion)
	if err != nil {
		return false, nil
	}
	autenticado := sesion.Values["autenticado"]
	if autenticado != nil {
		return autenticado.(bool), nil
	}
	return false, nil
}

func (s *Sesion) cerrarSesion() error {
	sesion, err := almacenamientoDeSesion.Get(s.request, NombreSesion)
	if err != nil {
		return err
	}
	sesion.Options.MaxAge = -1
	err = sesion.Save(s.request, s.responseWriter)
	return err
}
