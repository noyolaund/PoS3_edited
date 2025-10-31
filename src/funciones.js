import { DIA_EN_MILISEGUNDOS } from "./constantes";
import Vue from "vue";
import { HTTP, HTTP_AUTH } from "./http-common";
import ConectorPluginV3 from "./ConectorPluginV3.js";
export const FUNCIONES = {
  async imprimirReporteCaja(fechaInicio, fechaFin, idUsuario) {
    const impresionEscpos = (await HTTP_AUTH.get("valor/MODO_IMPRESION") === "Impresora térmica");
    if (!impresionEscpos) {
      Vue.$router.push({
        name: "TicketDeCaja",
        query: {
          fechaInicio,
          fechaFin,
          usuario: idUsuario
        }
      });
      return;
    }

    const serial = await HTTP_AUTH.get("valor/SERIAL_PLUGIN_IMPRESION");
    const filtros = Vue.options.filters;
    const ajustesEmpresa = await HTTP_AUTH.get("ajustes/empresa");
    let ruta = `reporte/caja/${fechaInicio}/${fechaFin}`;
    if (idUsuario) {
      ruta += `/${idUsuario}`;
    }
    const estado = await HTTP_AUTH.get(ruta);
    const fechaYHora = await HTTP.get("fechaYHora");
    const usuarioLogueado = await HTTP_AUTH.get("usuario/logueado");
    let usuario = {};
    if (idUsuario) {
      usuario = await HTTP_AUTH.get(`usuario/caja/${idUsuario}`)
    }
    const logotipo = require("@/assets/inicio/logo.png");
    const conector = new ConectorPluginV3(ConectorPluginV3.URL_PLUGIN_POR_DEFECTO, serial)
      .Iniciar()
      .ImprimirImagenEnBase64(
        logotipo,
        ConectorPluginV3.TAMAÑO_IMAGEN_NORMAL,
        320
      )
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Estado de caja\nTipo: ")
      .EstablecerEnfatizado(false)
    if (usuario.Numero) {
      conector.EscribirTexto(`De #${usuario.Numero} ${usuario.Nombre}`)
    } else {
      conector.EscribirTexto("General")
    }
    conector
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Desde: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(fechaInicio)
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Hasta: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(fechaFin)
      .Feed(1)
    const claves = ["Nombre", "Direccion", "Telefono"];
    for (const clave of claves) {
      if (ajustesEmpresa[clave]) {
        conector.EscribirTexto(ajustesEmpresa[clave]);
        conector.Feed(1);
      }
    }

    conector
      .EstablecerEnfatizado(false)
      .EscribirTexto("Impreso por: ")
      .EstablecerEnfatizado(true)
      .EscribirTexto(usuarioLogueado.Nombre)
      .Feed(1)
      .EscribirTexto(filtros.fechaExpresiva(fechaYHora))
      .Feed(1)
      .EscribirTexto("--------------------------\n")
      .EstablecerEnfatizado(false)
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_IZQUIERDA)
      .EscribirTexto("Ventas al contado\n")
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
      .EscribirTexto(filtros.currency(estado.VentasContado))
      .Feed(1)
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_IZQUIERDA)
      .EscribirTexto("Anticipo de apartados\n")
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
      .EscribirTexto(filtros.currency(estado.Anticipos))
      .Feed(1)
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_IZQUIERDA)
      .EscribirTexto("Abonos\n")
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
      .EscribirTexto(filtros.currency(estado.Abonos))
      .Feed(1)
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_IZQUIERDA)
      .EscribirTexto("Ingresos\n")
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
      .EscribirTexto(filtros.currency(estado.Ingresos))
      .Feed(1)
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_IZQUIERDA)
      .EscribirTexto("Egresos\n")
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
      .EscribirTexto(filtros.currency(estado.Egresos))
      .Feed(1)
      .EscribirTexto("--------------------------\n")
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Total: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(estado.VentasContado +
        estado.Anticipos +
        estado.Abonos +
        estado.Ingresos -
        estado.Egresos))
      .Feed(1)
    conector.Pulso(48, 60, 120).CorteParcial().Corte(1);
    return await conector.imprimirEn(await HTTP_AUTH.get("nombre/impresora"));
  },
  async imprimirTicketAbono(idAbono, idApartado) {
    const impresionEscpos = (await HTTP_AUTH.get("valor/MODO_IMPRESION") === "Impresora térmica");
    if (!impresionEscpos) {
      Vue.$router.push({
        name: "TicketDeAbono",
        params: { idApartado, idAbono },
      });
      return;
    }
    const serial = await HTTP_AUTH.get("valor/SERIAL_PLUGIN_IMPRESION");
    const filtros = Vue.options.filters;
    const ajustesEmpresa = await HTTP_AUTH.get("ajustes/empresa");
    const apartado = await HTTP_AUTH.get("apartado/" + idApartado);
    const abono = await HTTP_AUTH.get(`abono/${idAbono}/${idApartado}`);
    const logotipo = require("@/assets/inicio/logo.png");
    const conector = new ConectorPluginV3(ConectorPluginV3.URL_PLUGIN_POR_DEFECTO, serial)
      .Iniciar()
      .ImprimirImagenEnBase64(
        logotipo,
        ConectorPluginV3.TAMAÑO_IMAGEN_NORMAL,
        320
      )
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Ticket de abono #" + idAbono + "\n")
      .Feed(1);
    const claves = ["Nombre", "Direccion", "Telefono"];
    for (const clave of claves) {
      if (ajustesEmpresa[clave]) {
        conector.EscribirTexto(ajustesEmpresa[clave]);
        conector.Feed(1);
      }
    }
    conector
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.fechaExpresiva(apartado.Fecha))
      .Feed(1)
      .DeshabilitarElModoDeCaracteresChinos()
      .EscribirTexto("Lo atendió: ")
      .EstablecerEnfatizado(true)
      .EscribirTexto(abono.Usuario.Nombre)
      .EstablecerEnfatizado(false)
      .Feed(1)
      .EscribirTexto("Cliente: ")
      .EstablecerEnfatizado(true)
      .EscribirTexto(apartado.Cliente.Nombre)
      .EstablecerEnfatizado(false)
      .Feed(1)
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EscribirTexto("--------------------------\n");

    for (const producto of apartado.Productos) {
      conector
        .EstablecerAlineacion(ConectorPluginV3.ALINEACION_IZQUIERDA)
        .EscribirTexto(producto.Descripcion)
        .Feed(1)
        .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
        .EscribirTexto(
          `${producto.Cantidad} x ${filtros.currency(
            producto.PrecioVenta
          )} = ${filtros.currency(
            producto.Cantidad * producto.PrecioVenta
          )}`
        )
        .Feed(1);
    }
    conector
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EscribirTexto("--------------------------\n")
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Total: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(apartado.Total))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Restante anterior: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(apartado.Total - apartado.Abonado - apartado.Anticipo + abono.Monto))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Cantidad abonada: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(abono.Monto))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Su pago: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(abono.Pago))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Cambio: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(
        filtros.currency(apartado.Pago - apartado.Anticipo)
      )
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Restante actual: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(apartado.Total - apartado.Abonado - apartado.Anticipo))
      .Feed(1)
    if (ajustesEmpresa.MensajePersonal) {
      conector
        .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
        .EstablecerEnfatizado(true)
        .EscribirTexto(ajustesEmpresa.MensajePersonal)
        .Feed(1);
    }
    conector.Pulso(48, 60, 120).CorteParcial().Corte(1);
    return await conector.imprimirEn(await HTTP_AUTH.get("nombre/impresora"));
  },

  async imprimirTicketApartado(idApartado) {
    const impresionEscpos = (await HTTP_AUTH.get("valor/MODO_IMPRESION") === "Impresora térmica");
    if (!impresionEscpos) {
      Vue.$router.push({
        name: "TicketDeApartado",
        params: { idApartado },
      });
      return;
    }
    const serial = await HTTP_AUTH.get("valor/SERIAL_PLUGIN_IMPRESION");
    const filtros = Vue.options.filters;
    const ajustesEmpresa = await HTTP_AUTH.get("ajustes/empresa");
    const apartado = await HTTP_AUTH.get("apartado/" + idApartado);
    const logotipo = require("@/assets/inicio/logo.png");
    const conector = new ConectorPluginV3(ConectorPluginV3.URL_PLUGIN_POR_DEFECTO, serial)
      .Iniciar()
      .ImprimirImagenEnBase64(
        logotipo,
        ConectorPluginV3.TAMAÑO_IMAGEN_NORMAL,
        320
      )
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Ticket de apartado #" + idApartado + "\n")
      .Feed(1);
    const claves = ["Nombre", "Direccion", "Telefono"];
    for (const clave of claves) {
      if (ajustesEmpresa[clave]) {
        conector.EscribirTexto(ajustesEmpresa[clave]);
        conector.Feed(1);
      }
    }
    conector
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.fechaExpresiva(apartado.Fecha))
      .Feed(1)
      .DeshabilitarElModoDeCaracteresChinos()
      .EscribirTexto("Lo atendió: ")
      .EstablecerEnfatizado(true)
      .EscribirTexto(apartado.Usuario.Nombre)
      .EstablecerEnfatizado(false)
      .Feed(1)
      .EscribirTexto("Cliente: ")
      .EstablecerEnfatizado(true)
      .EscribirTexto(apartado.Cliente.Nombre)
      .EstablecerEnfatizado(false)
      .Feed(1)
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EscribirTexto("--------------------------\n");
    for (const producto of apartado.Productos) {
      conector
        .EstablecerAlineacion(ConectorPluginV3.ALINEACION_IZQUIERDA)
        .EscribirTexto(producto.Descripcion)
        .Feed(1)
        .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
        .EscribirTexto(
          `${producto.Cantidad} x ${filtros.currency(
            producto.PrecioVenta
          )} = ${filtros.currency(
            producto.Cantidad * producto.PrecioVenta
          )}`
        )
        .Feed(1);
    }
    conector
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EscribirTexto("--------------------------\n")
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Total: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(apartado.Total))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Su pago: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(apartado.Pago))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Anticipo: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(apartado.Anticipo))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Cambio: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(
        filtros.currency(apartado.Pago - apartado.Anticipo)
      )
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Restante: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(apartado.Total - apartado.Anticipo - apartado.Abonado))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Nota: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(`no nos hacemos responsables de los articulos después de la fecha de vencimiento: ${filtros.fechaSinHora(apartado.FechaVencimiento)}`)
      .Feed(1)
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EscribirTexto("Firma del cliente\n\n\n")
      .EscribirTexto("_________________________")
      .Feed(1)

    if (ajustesEmpresa.MensajePersonal) {
      conector
        .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
        .EstablecerEnfatizado(true)
        .EscribirTexto(ajustesEmpresa.MensajePersonal)
        .Feed(1);
    }
    conector.Pulso(48, 60, 120).CorteParcial().Corte(1);
    return await conector.imprimirEn(await HTTP_AUTH.get("nombre/impresora"));
  },
  async imprimirTicketVentaContado(idVenta) {
    const impresionEscpos = (await HTTP_AUTH.get("valor/MODO_IMPRESION") === "Impresora térmica");
    if (!impresionEscpos) {
      Vue.$router.push({
        name: "TicketDeVentaContado",
        params: { idVenta },
      });
      return;
    }
    const serial = await HTTP_AUTH.get("valor/SERIAL_PLUGIN_IMPRESION");
    const filtros = Vue.options.filters;
    const ajustesEmpresa = await HTTP_AUTH.get("ajustes/empresa");
    const venta = await HTTP_AUTH.get("venta/contado/" + idVenta);
    const logotipo = require("@/assets/inicio/logo.png");
    const conector = new ConectorPluginV3(ConectorPluginV3.URL_PLUGIN_POR_DEFECTO, serial)
      .Iniciar()
      .ImprimirImagenEnBase64(
        logotipo,
        ConectorPluginV3.TAMAÑO_IMAGEN_NORMAL,
        320
      )
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Ticket de venta #" + idVenta + "\n")
      .Feed(1);
    const claves = ["Nombre", "Direccion", "Telefono"];
    for (const clave of claves) {
      if (ajustesEmpresa[clave]) {
        conector.EscribirTexto(ajustesEmpresa[clave]);
        conector.Feed(1);
      }
    }
    conector
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.fechaExpresiva(venta.Fecha))
      .Feed(1)
      .DeshabilitarElModoDeCaracteresChinos()
      .EscribirTexto("Lo atendió: ")
      .EstablecerEnfatizado(true)
      .EscribirTexto(venta.Usuario.Nombre)
      .EstablecerEnfatizado(false)
      .Feed(1)
      .EscribirTexto("Cliente: ")
      .EstablecerEnfatizado(true)
      .EscribirTexto(venta.Cliente.Nombre)
      .EstablecerEnfatizado(false)
      .Feed(1)
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EscribirTexto("--------------------------\n");
    for (const producto of venta.Productos) {
      conector
        .EstablecerAlineacion(ConectorPluginV3.ALINEACION_IZQUIERDA)
        .EscribirTexto(producto.Descripcion)
        .Feed(1)
        .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
        .EscribirTexto(
          `${producto.Cantidad} x ${filtros.currency(
            producto.PrecioVenta
          )
          } = ${filtros.currency(
            producto.Cantidad * producto.PrecioVenta
          )
          }`
        )
        .Feed(1);
    }
    conector
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
      .EscribirTexto("--------------------------\n")
      .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Total: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(venta.Total))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Su pago: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(filtros.currency(venta.Pago))
      .Feed(1)
      .EstablecerEnfatizado(true)
      .EscribirTexto("Cambio: ")
      .EstablecerEnfatizado(false)
      .EscribirTexto(
        filtros.currency(venta.Pago - venta.Total)
      )
      .Feed(1);
    if (ajustesEmpresa.MensajePersonal) {
      conector
        .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
        .EstablecerEnfatizado(true)
        .EscribirTexto(ajustesEmpresa.MensajePersonal)
        .Feed(1);
    }
    conector.Pulso(48, 60, 120).CorteParcial().Corte(1);
    return await conector.imprimirEn(await HTTP_AUTH.get("nombre/impresora"));
  },
  /**
   * Devuelve la fecha que se le pasa, más la hora a las 12 de la noche
   *
   * Por ejemplo, si se le pasa algo como 2018-06-20 devuelve 2018-06-20T00:00:00
   *
   * */
  agregarHoraCeroAFecha(fecha) {
    return `${fecha}T00:00:00`
  },
  /**
   * Devuelve la fecha de hoy en el formato:
   * YYYY-MM-DD
   *
   * Por ejemplo:
   *
   * 2018-06-20
   * */
  hoyComoCadena() {
    let fecha = new Date(Date.now() + 864e5),
      anio = fecha.getFullYear(),
      mes = this.agregarCerosALaIzquierdaSiEsNecesario(fecha.getMonth() + 1),
      dia = this.agregarCerosALaIzquierdaSiEsNecesario(fecha.getDate());
    return `${anio}-${mes}-${dia}`;
  },
  /**
   * Lo único que hace es convertir la fecha a algo como
   * 2018-05-24
   *
   * Por cierto, getMonth() devuelve el mes pero tomando en cuenta que enero es 0 y
   * diciembre es 11
   */
  formatearFecha(fecha) {
    if (!fecha instanceof Date)
      throw new TypeError("La fecha debe ser un objeto de tipo Date");

    let mes = fecha.getMonth() + 1, diaDelMes = fecha.getDate();
    return `${fecha.getFullYear()}-${this.agregarCerosALaIzquierdaSiEsNecesario(mes)}-${this.agregarCerosALaIzquierdaSiEsNecesario(diaDelMes)}`;
  },
  componerFechaParaFin(fecha) {
    fecha = new Date(fecha);
    fecha.setTime(fecha.getTime() + fecha.getTimezoneOffset() * 60 * 1000); //Necesario si viene de cadena. Ver https://parzibyte.me/blog/2018/03/05/ajustando-fechas-javascript/
    fecha.setTime(fecha.getTime() + DIA_EN_MILISEGUNDOS);
    return this.formatearFecha(fecha);
  },
  componerFechaParaInicio(fecha) {
    return fecha;
  },
  esteAnioComoCadena() {
    return new Date().getFullYear().toString();
  },
  esteMesComoCadena() {
    let mes = new Date().getMonth() + 1;

    return this.agregarCerosALaIzquierdaSiEsNecesario(mes);
  },
  agregarCerosALaIzquierdaSiEsNecesario(valor) {
    valor = valor.toString();
    return valor.length < 2 ? `0${valor}` : valor.toString();
  },
  colorHexadecimalAleatorio() {
    let colores = ["#f2476a", "#fb654e", "#eb2d3a", "#add8e6", "#90ee90", "#ffcb7e", "#ff9464", "#d5e389", "#d0b191", "#c18390", "#F44336", "#9C27B0", "#673AB7", "#3F51B5", "##2196F3", "#4CAF50", "#8BC34A", "#CDDC39", "##FFEB3B", "#FF9800", "#FF5722"];
    return colores[Math.floor(Math.random() * colores.length)];
  },
  /*
  * @deprecated
  * */
  _colorHexadecimalAleatorio() {
    return "#000000".replace(/0/g, () => (~~(Math.random() * 16)).toString(16));
  }
};
