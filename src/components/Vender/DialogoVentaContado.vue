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
          await FUNCIONES.imprimirTicketVentaContado(resultados.Numero);
        }
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
