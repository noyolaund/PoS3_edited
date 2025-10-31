<template>
  <v-dialog v-model="mostrar" persistent max-width="700">
    <v-card>
      <v-card-title class="headline">Hacer apartado</v-card-title>
      <v-card-text>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12 sm4>
                <v-text-field
                  @keydown.native.enter="guardar()"
                  ref="pago"
                  prepend-icon="attach_money"
                  label="Pago del cliente"
                  type="number"
                  v-model.number="pago"
                  :rules="reglas.pago"
                  hint="¿Con cuánto paga el cliente (para calcular el cambio)?"
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12 sm4>
                <v-text-field
                  @keydown.native.enter="guardar()"
                  ref="anticipo"
                  prepend-icon="attach_money"
                  label="Anticipo"
                  type="number"
                  v-model.number="anticipo"
                  :rules="reglas.anticipo"
                  hint="¿Cuál es el anticipo del cliente?"
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12 sm4>
                <v-menu
                  ref="menu"
                  :close-on-content-click="false"
                  v-model="mostrarDialogoFecha"
                  :nudge-right="40"
                  :return-value.sync="fechaVencimiento"
                  lazy
                  transition="scale-transition"
                  offset-y
                  full-width
                  min-width="290px"
                >
                  <v-text-field
                    slot="activator"
                    v-model="fechaVencimiento"
                    label="Fecha de vencimiento"
                    prepend-icon="event"
                    readonly
                  ></v-text-field>
                  <v-date-picker
                    color="green lighten-1"
                    :events="fechasRecomendadas"
                    locale="es-419"
                    v-model="fechaVencimiento"
                    :min="hoy"
                  >
                    <v-spacer></v-spacer>
                    <v-btn
                      flat
                      color="primary"
                      @click="mostrarDialogoFecha = false"
                      >Cerrar</v-btn
                    >
                    <v-btn
                      flat
                      color="primary"
                      @click="$refs.menu.save(fechaVencimiento)"
                      >OK</v-btn
                    >
                  </v-date-picker>
                </v-menu>
              </v-flex>
            </v-layout>
          </v-container>
        </v-form>
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
        <v-flex xs12>
          <p class="title text-xs-right" color="blue">
            Total: {{ datosVenta.total | currency }}
          </p>
          <p
            v-show="pago > 0 && cambio >= 0"
            class="title text-xs-right"
            color="blue"
          >
            Pago: {{ pago | currency }}
          </p>
          <p
            v-show="pago > 0 && cambio >= 0"
            class="title text-xs-right"
            color="blue"
          >
            Anticipo: {{ anticipo | currency }}
          </p>
          <p v-show="cambio >= 0" class="title text-xs-right" color="blue">
            Cambio: {{ cambio | currency }}
          </p>
          <p v-show="cambio >= 0" class="title text-xs-right" color="blue">
            Restante: {{ (datosVenta.total - anticipo) | currency }}
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
          >Terminar apartado
        </v-btn>
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
import { FUNCIONES } from "../../funciones";

export default {
  mounted() {
    let fecha = new Date(),
      diasRecomendados = [14, 29, 59]; // A 14 días, 29 o 59
    //Nota: 864e5 es igual a 86400000, lo que equivale a un día expresado en milisegundos
    this.fechasRecomendadas = diasRecomendados.map(numeroDeDias =>
      new Date(fecha.getTime() + 864e5 * numeroDeDias)
        .toISOString()
        .substr(0, 10)
    );
  },
  components: {
    DetallesClienteSeleccionado,
    AutocompletadoClientes
  },
  computed: {
    cambio() {
      if (this.datosVenta && this.pago > 0 && this.anticipo)
        return this.pago - this.anticipo;
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
      this.$nextTick(this.$refs.pago.focus);
    },
    resetearFormulario() {
      this.$refs.formulario.reset();
    },
    cerrarDialogo() {
      this.resetearFormulario();
      this.$emit("cerrar-dialogo");
    },
    async guardar() {
      if (!this.$refs.formulario.validate()) {
        return;
      }
      if (this.anticipo >= this.datosVenta.total) return this.$emit("error-pago-excedido");
      if (null === this.fechaVencimiento) {
        this.$emit("no-hay-fecha");
        this.mostrarDialogoFecha = true;
        return;
      }
      if (
        null === this.clienteSeleccionado ||
        !this.clienteSeleccionado.Nombre
      ) {
        return this.$emit("no-hay-cliente");
      }
      let apartado = {
        Total: this.datosVenta.total,
        Productos: this.datosVenta.lista,
        Cliente: this.clienteSeleccionado,
        FechaVencimiento: FUNCIONES.agregarHoraCeroAFecha(this.fechaVencimiento),
        Anticipo: this.anticipo,
        Pago: this.pago,
      };
      this.cargando = true;
      const resultados = await HTTP_AUTH.post("apartado", apartado);
      this.cargando = false;
      if (resultados) {
        this.$emit("apartado-realizado");
        this.prepararNuevaVenta();
        FUNCIONES.imprimirTicketApartado(resultados.Numero);
      }
    }
  },
  props: ["mostrar", "datosVenta"],
  data: () => ({
    cargando: false,
    mostrarDialogoFecha: null,
    fechasRecomendadas: null,
    hoy: FUNCIONES.hoyComoCadena(),
    fechaVencimiento: null,
    clienteSeleccionado: {},
    anticipo: null,
    pago: null,
    reglas: {
      pago: [
        pago => {
          if (pago < 0) return "Introduzca un valor mayor que 0";
          return true;
        }
      ],
      anticipo: [
        anticipo => {
          if (anticipo < 0) return "Introduzca un valor mayor que 0";
          return true;
        }
      ]
    }
  })
};
</script>
