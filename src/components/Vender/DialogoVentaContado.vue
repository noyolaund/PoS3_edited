<template>
  <v-dialog v-model="mostrar" persistent max-width="700">
    <v-card>
      <v-card-title class="headline">Terminar venta al contado</v-card-title>
      <v-card-text>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12>
                <v-text-field
                  @keyup.enter="guardar()"
                  ref="pagoCliente"
                  prepend-icon="attach_money"
                  label="Cantidad recibida"
                  type="number"
                  v-model.number="pagoDelCliente"
                  :rules="reglas.pago"
                  hint="¿Con cuánto paga el cliente?"
                  required
                ></v-text-field>
              </v-flex>
            </v-layout>
          </v-container>
          <v-flex xs12>
            <span class="body-2">Tipo de cliente</span>
            <v-radio-group v-model="tipoCliente" row>
              <v-radio label="Mostrador" value="mostrador"></v-radio>
              <v-radio
                label="Buscar o crear nuevo"
                value="existenteONuevo"
              ></v-radio>
            </v-radio-group>
          </v-flex>
          <v-flex v-show="tipoCliente === 'existenteONuevo'">
            <detalles-cliente-seleccionado
              :clienteSeleccionado="clienteSeleccionado"
            ></detalles-cliente-seleccionado>
            <autocompletado-clientes
              ref="autocompletado"
              @cliente-cancelado="onClienteCancelado"
              @cliente-seleccionado="onClienteSeleccionado"
              @agregar-cliente="agregarNuevoCliente"
            ></autocompletado-clientes>
            <br />
          </v-flex>
        </v-form>
        <v-flex xs12>
          <p class="title text-xs-right" color="blue">
            Total: {{ datosVenta.total | currency }}
          </p>
          <p
            v-show="pagoDelCliente > 0"
            class="title text-xs-right"
            color="blue"
          >
            Pago: {{ pagoDelCliente | currency }}
          </p>
          <p v-show="cambio >= 0" class="title text-xs-right" color="blue">
            Cambio: {{ cambio | currency }}
          </p>
        </v-flex>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          :loading="cargando"
          color="green darken-1"
          flat="flat"
          @click.native="guardar()"
          >Terminar venta</v-btn
        >
        <v-btn
          :loading="cargando"
          color="indigo darken-1"
          flat="flat"
          @click.native="guardarConTicket()"
          >Terminar venta con Ticket</v-btn
        >
        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import { HTTP_AUTH } from "../../http-common";
import DetallesClienteSeleccionado from "../Vender/DetallesClienteSeleccionado";
import AutocompletadoClientes from "../Vender/AutocompletadoClientes";
import { FUNCIONES } from '../../funciones';
import { obtenerModoImpresion } from '../../modoImpresionHelper';
import ConectorJavascript from '../../ConectorJavascript';
import ConectorPluginV3 from '../../ConectorPluginV3';

export default {
  components: { DetallesClienteSeleccionado, AutocompletadoClientes },
  computed: {
    cambio() {
      if (this.datosVenta && this.pagoDelCliente > 0)
        return this.pagoDelCliente - this.datosVenta.total;
      return -1;
    }
  },
  watch: {
    mostrar(estaMostrado) {
      if (estaMostrado) this.enfocarInputPrincipal();
    }
  },
  methods: {
    prepararNuevaVenta() {
      this.setCliente({});
      this.$refs.formulario.reset();
      this.$refs.autocompletado.limpiar();
      this.tipoCliente = "mostrador";
    },
    setCliente(cliente) {
      this.clienteSeleccionado = cliente;
    },
    onClienteSeleccionado(cliente) {
      this.setCliente(Object.assign({}, cliente));
    },
    onClienteCancelado() {
      this.setCliente({});
    },
    agregarNuevoCliente() {
      this.$emit("agregar-cliente");
    },
    enfocarInputPrincipal() {
      this.$nextTick(this.$refs.pagoCliente.focus);
    },
    resetearFormulario() {
      this.$refs.formulario.reset();
    },
    cerrarDialogo() {
      this.resetearFormulario();
      this.$emit("cerrar-dialogo");
    },
    async guardar() {
      if (this.$refs.formulario.validate()) {
        if (this.cambio < 0) return this.$emit("error-pago-incompleto");
        let cliente = Object.assign({}, this.clienteSeleccionado);
        if (this.tipoCliente === "existenteONuevo") {
          if (null === cliente || !cliente.Nombre) {
            return this.$emit("no-hay-cliente");
          }
        } else {
          cliente.Numero = 1; //Para que tome el mostrador, que se supone será el cliente con el id 1
        }
        let venta = {
          Total: this.datosVenta.total,
          Productos: this.datosVenta.lista,
          Cliente: cliente,
          Pago: this.pagoDelCliente,
        };
        this.cargando = true;
        const resultados = await HTTP_AUTH.post("venta/contado", venta);
        this.cargando = false;
        if (resultados) {
          this.$emit("venta-realizada");
          this.prepararNuevaVenta();
        }
      }
    }
    ,
    async guardarConTicket() {
      console.log('[DialogoVentaContado] guardarConTicket iniciado');
      if (this.$refs.formulario.validate()) {
        console.log('[DialogoVentaContado] Formulario válido');
        if (this.cambio < 0) {
          console.log('[DialogoVentaContado] Pago incompleto');
          return this.$emit("error-pago-incompleto");
        }
        let cliente = Object.assign({}, this.clienteSeleccionado);
        if (this.tipoCliente === "existenteONuevo") {
          if (null === cliente || !cliente.Nombre) {
            console.log('[DialogoVentaContado] No hay cliente seleccionado');
            return this.$emit("no-hay-cliente");
          }
        } else {
          cliente.Numero = 1;
        }
        let venta = {
          Total: this.datosVenta.total,
          Productos: this.datosVenta.lista,
          Cliente: cliente,
          Pago: this.pagoDelCliente,
        };
        console.log('[DialogoVentaContado] Guardando venta:', venta);
        this.cargando = true;
        const resultados = await HTTP_AUTH.post("venta/contado", venta);
        this.cargando = false;
        console.log('[DialogoVentaContado] Venta guardada, resultado:', resultados);
        if (resultados) {
          console.log('[DialogoVentaContado] Llamando a imprimirTicketVentaContado con ID:', resultados.Numero);
          
          // Verificar modo de impresión con localStorage fallback
          const modoImpresion = await obtenerModoImpresion();
          console.log('[DialogoVentaContado] Modo de impresión obtenido:', modoImpresion);
          
          if (modoImpresion === "Impresora térmica" || modoImpresion === "BridgeJavascript") {
            console.log('[DialogoVentaContado] Modo es térmico/bridge, imprimiendo directamente');
            
            // Obtener datos para imprimir
            const venta = await HTTP_AUTH.get("venta/contado/" + resultados.Numero);
            const ajustesEmpresa = await HTTP_AUTH.get("ajustes/empresa");
            const nombreImpresora = await HTTP_AUTH.get("nombre/impresora");
            const logotipo = require("@/assets/inicio/logo.png");
            
            // Crear conector según el modo
            let conector;
            if (modoImpresion === "BridgeJavascript") {
              console.log('[DialogoVentaContado] Usando ConectorJavascript');
              conector = new ConectorJavascript();
            } else {
              console.log('[DialogoVentaContado] Usando ConectorPluginV3');
              const serial = await HTTP_AUTH.get("valor/SERIAL_PLUGIN_IMPRESION");
              conector = new ConectorPluginV3(ConectorPluginV3.URL_PLUGIN_POR_DEFECTO, serial);
            }
            
            // Construir ticket
            conector
              .Iniciar()
              .DeshabilitarElModoDeCaracteresChinos()
              .ImprimirImagenEnBase64(logotipo, ConectorPluginV3.TAMAÑO_IMAGEN_NORMAL, 320)
              .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
              .EstablecerEnfatizado(true)
              .EscribirTexto("Ticket de venta #" + resultados.Numero + "\n")
              .EstablecerEnfatizado(false);
            
            // Datos de la empresa
            const claves = ["Nombre", "Direccion", "Telefono"];
            for (const clave of claves) {
              if (ajustesEmpresa[clave]) {
                conector.EscribirTexto(ajustesEmpresa[clave] + "\n");
              }
            }
            
            const textoCustom = await HTTP_AUTH.get("valor/TEXTO_TICKET_CUSTOM");
            const nombreNegocio = textoCustom ? textoCustom : "Deposito Beer Broos";

            conector
              .EscribirTexto(new Date().toLocaleString() + "\n")
              .EscribirTexto("Lo atendio: ")
              .EstablecerEnfatizado(true)
              .EscribirTexto(venta.Usuario.Nombre + "\n")
              .EstablecerEnfatizado(false)
              .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
              .EstablecerEnfatizado(true)
              .EscribirTexto("*** " + nombreNegocio + " ***\n")
              .EstablecerEnfatizado(false)
              .EscribirTexto("--------------------------\n");
            
            // Productos
            for (const producto of venta.Productos) {
              conector
                .EstablecerAlineacion(ConectorPluginV3.ALINEACION_IZQUIERDA)
                .EscribirTexto(producto.Descripcion + "\n")
                .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
                .EscribirTexto(`${producto.Cantidad} x $${producto.PrecioVenta} = $${producto.Cantidad * producto.PrecioVenta}\n`);
            }
            
            conector
              .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
              .EscribirTexto("--------------------------\n")
              .EstablecerAlineacion(ConectorPluginV3.ALINEACION_DERECHA)
              .EstablecerEnfatizado(true)
              .EscribirTexto("Total: ")
              .EstablecerEnfatizado(false)
              .EscribirTexto("$" + venta.Total + "\n")
              .EstablecerEnfatizado(true)
              .EscribirTexto("Su pago: ")
              .EstablecerEnfatizado(false)
              .EscribirTexto("$" + venta.Pago + "\n")
              .EstablecerEnfatizado(true)
              .EscribirTexto("Cambio: ")
              .EstablecerEnfatizado(false)
              .EscribirTexto("$" + (venta.Pago - venta.Total) + "\n");
            
            if (ajustesEmpresa.MensajePersonal) {
              conector
                .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
                .EstablecerEnfatizado(true)
                .EscribirTexto(ajustesEmpresa.MensajePersonal + "\n");
            }
            
            conector
              .EstablecerAlineacion(ConectorPluginV3.ALINEACION_CENTRO)
              .EscribirTexto("\n")
              .EstablecerEnfatizado(true)
              .EscribirTexto("GRACIAS POR SU COMPRA\n")
              .EstablecerEnfatizado(false);
            
            conector.Pulso(48, 60, 120).CorteParcial().Corte(1);
            
            console.log('[DialogoVentaContado] Enviando a imprimir en:', nombreImpresora);
            const resultado = await conector.imprimirEn(nombreImpresora);
            console.log('[DialogoVentaContado] Resultado de impresión:', resultado);
          } else {
            console.log('[DialogoVentaContado] Modo no es térmico/bridge, saltando impresión');
          }
          
          console.log('[DialogoVentaContado] Impresión completada');
          this.$emit("venta-realizada-con-ticket");
          this.prepararNuevaVenta();
        } else {
          console.error('[DialogoVentaContado] No se recibió resultado de la venta');
        }
      } else {
        console.log('[DialogoVentaContado] Formulario inválido');
      }
    }
  },
  props: ["mostrar", "datosVenta"],
  data: () => ({
    cargando: false,
    clienteSeleccionado: {},
    pagoDelCliente: null,
    tipoCliente: "mostrador",
    reglas: {
      pago: [
        pago => {
          if (pago <= 0) return "Introduzca un valor mayor que 0";
          return true;
        }
      ]
    }
  })
};
</script>
