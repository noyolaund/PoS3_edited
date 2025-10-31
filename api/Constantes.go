package main

const (
	/*Algunas constantes usadas para respuesta en JSON. Verificar constantes.js*/
	RespuestaErrorNegocioExistente          = 0
	RespuestaErrorRegistrandoNegocio        = 1
	RespuestaNegocioRegistradoCorrectamente = 2
	RespuestaLoginNegocioNoVerificado       = 3
	RespuestaLoginCorrecto                  = 4
	RespuestaLoginError                     = 5
	RespuestaLoginIncorrecto                = 6
	// Constantes del sistema
	ImpresionNativa                             = false
	TamanioMaximoArchivoImportacionExcel        = 5 << 20 // 5MB
	NombreArchivoExportadoCSV                   = "ProductosExportados_SPOS3.csv"
	NombreArchivoExportadoExcel                 = "ProductosExportados_SPOS3.xlsx"
	LimiteProductosMasVendidos                  = 15
	DiasParaApartadosProximosAVencer            = 7
	LimiteProductosMenosVendidos                = 15
	LimiteProductosNuncaVendidos                = 15
	LimiteAutoCompletadoClientes                = 10
	MensajePieDeTicket                          = "bit.ly/sublime-pos"
	MIME                                        = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	NumeroAvisosAntesDeEliminarCuenta           = 3
	DiasParaMarcarCuentaComoInactiva            = 7
	AsuntoReactivarCuenta                       = "Reactiva tu cuenta de SPOS 3"
	MensajeCompartir                            = "Conoce Sublime POS 3. Un punto de venta completo y gratuito"
	UrlAcortadaCompartirSistema                 = "http://bit.ly/sublime-pos"
	RutaDirectorioContenidoEstatico             = "/static/"
	NombreSesion                                = "SublimePosSesion"
	NombreBaseDeDatosSesiones                   = "./sesiones.sqlite"
	ClaveCookie                                 = "905ff3f39809b9f9e204204596e234589a99a8fab5e691c08474584a5f3b7523"
	EdadDeSesionEnSegundos                      = 86400  // Una semana en segundos
	DominioPermitidoCORS                        = "http://localhost:8080"
	PuertoServidor                              = ":2106" // Debe coincidir con constantes.js
	RutaGeneralNoNecesitaComprobacion           = "RUTA_NO_AUTH"
	ControladorBD                               = "sqlite3"
	NombreDirectorioParaSubidas                 = "subidas"
	CadenaConexionBDUsuarios                    = "negocio_parzibyte.me_"
	CadenaConexionBDNegocios                    = "negocios_spos3_by_parzibyte.me.db"
	PrefijoBDNegocios                           = "spos"
	NombreArchivoEsquemaSQLInit                 = "init.sql"
	NombreArchivoEsquemaSQLNegocios             = "esquema_negocios_sqlite.sql"
	NombreArchivoEsquemaSQLSistemas             = "esquema_spos_sqlite.sql"
	NombrePrimerUsuarioAdmin                    = "parzibyte"
	ClaveAPICorreo                              = ""
	PrefijoRutasAdmin                           = "/auth"
	GmailServidor                               = ""
	GmailPuerto                                 = ""
	GmailCorreo                                 = ""
	GmailPass                                   = ""
	UrlBaseApp                                  = ""
	UrlBasePaginaWeb                            = "https://parzibyte.me/apps/sublime-pos-3/"
	CorreoSoporteYContacto                      = "parzibyte@gmail.com"
	NombreDirectorioContenidoStaticoParaCliente = "./dist" // El nombre de la carpeta que tiene HTML, CSS y JS para ser servido
)
