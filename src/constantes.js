export const MESES = "Enero Febrero Marzo Abril Mayo Junio Julio Agosto Septiembre Octubre Noviembre Diciembre".split(" ");
export const DIA_EN_MILISEGUNDOS = 864e5; // 1000 * 60 * 60 * 24
/*Algunas constantes usadas para respuesta en JSON. Verificar Constantes.go*/
export const RespuestaErrorNegocioExistente = 0;
export const RespuestaErrorRegistrandoNegocio = 1;
export const RespuestaNegocioRegistradoCorrectamente = 2;
export const RespuestaLoginNegocioNoVerificado = 3;
export const RespuestaLoginCorrecto = 4;
export const RespuestaLoginError = 5;
export const RespuestaLoginIncorrecto = 6;
export const MensajeCompartir = "Conoce Sublime POS 3. Un punto de venta completo y gratuito";
export const EnlaceCompartir = encodeURIComponent("http://bit.ly/sublime-pos");
export const EnlaceCompartirParaHumanos = "bit.ly/sublime-pos";
export const CorreoSoporteYContacto = "parzibyte@gmail.com";
export const TimeoutOcultarMenuTickets = 200;
export const PuertoApi = "2106"; // Debe coincidir con Constantes.go
export const UsuarioPorDefecto = "parzibyte"; // Debe coincidir con Constantes.go