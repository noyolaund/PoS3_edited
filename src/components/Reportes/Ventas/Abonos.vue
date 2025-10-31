<template>
  <v-dialog persistent v-model="mostrar">
    <v-card>
      <v-card-title class="headline"
        >Abonos del apartado #{{ apartado.IdApartado }}</v-card-title
      >
      <v-card-text>
        <v-layout row wrap>
          <v-flex xs12>
            <v-alert :value="estaVencido && restante > 0" type="error">
              Este apartado está vencido y no se pueden realizar más abonos
            </v-alert>
            <v-alert :value="restante <= 0" type="success">
              Este apartado ha sido liquidado
            </v-alert>
            <v-text-field
              ref="inputPago"
              v-show="restante > 0 && !estaVencido"
              @keyup.enter.prevent="abonar()"
              v-model.number="pago"
              prepend-icon="monetization_on"
              label="Pago del cliente"
              type="text"
              hint="Cantidad que da el cliente"
              required
            ></v-text-field>
            <v-text-field
              ref="inputAbono"
              v-show="restante > 0 && !estaVencido"
              @keyup.enter.prevent="abonar()"
              v-model.number="cantidadAbonada"
              prepend-icon="monetization_on"
              label="Cantidad a abonar"
              type="text"
              hint="Del pago del cliente, ¿Cuánto se abona?"
              required
            ></v-text-field>
            <v-flex xs12>
              <p class="title">Pago: {{ pago | currency }}</p>
              <p class="title">Abono: {{ cantidadAbonada | currency }}</p>
              <p class="title" v-show="cambio >= 0">
                Cambio: {{ cambio | currency }}
              </p>
              <p class="title" v-show="cambio >= 0">
                Restante: {{ (restante - cantidadAbonada) | currency }}
              </p>
            </v-flex>
          </v-flex>
          <v-flex xs12>
            <v-btn
              title="Haga click aquí para realizar el abono"
              :loading="cargandoAbonando"
              v-show="cantidadAbonada > 0"
              @click="abonar()"
              small
              color="success"
            >
              Terminar abono
            </v-btn>
          </v-flex>
          <v-flex xs12 sm6>
            <h1>
              <span class="display-1">{{ totalAbonado | currency }}</span>
              <span class="title">Abonado</span>
            </h1>
          </v-flex>
          <v-flex xs12 sm6>
            <h1>
              <span class="display-1">{{ restante | currency }}</span>
              <span class="title">Restante</span>
            </h1>
          </v-flex>
          <v-flex xs12>
            <v-data-table
              :headers="encabezados"
              :items="abonos"
              hide-actions
              item-key="props.item.index"
            >
              <template slot="items" slot-scope="props">
                <tr>
                  <td>{{ props.item.Monto | currency }}</td>
                  <td>{{ props.item.Fecha | fechaExpresiva }}</td>
                  <td>{{ props.item.Usuario.Nombre }}</td>
                  <td class="justify-center layout px-0">
                    <v-btn
                      title="Imprimir"
                      icon
                      class="mx-0"
                      @click="imprimir(props.item.IdAbono, apartado.IdApartado)"
                    >
                      <v-icon color="orange">print</v-icon>
                    </v-btn>
                  </td>
                </tr>
              </template>
            </v-data-table>
          </v-flex>
        </v-layout>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="green darken-1"
          flat="flat"
          @click.native="ocultarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import { FUNCIONES } from '../../../funciones';
import { HTTP_AUTH } from "../../../http-common";

export default {
  props: ["mostrar", "apartado"],
  created() {
    this.obtenerAbonos();
  },
  computed: {
    estaVencido() {
      if (!this.apartado.FechaVencimiento) return true;
      let fechaVencimiento = new Date(this.apartado.FechaVencimiento).getTime(),
        hoy = new Date().getTime();
      return fechaVencimiento < hoy;
    },
    totalAbonado() {
      return this.abonos.reduce(
        (acumulador, siguiente) => ({
          Monto: acumulador.Monto + siguiente.Monto
        }),
        {
          Monto: 0
        }
      ).Monto;
    },
    restante() {
      return this.apartado.Monto - this.apartado.Anticipo - this.totalAbonado;
    },
    cambio() {
      if (this.pago >= 0 && this.cantidadAbonada >= 0) {
        return this.pago - this.cantidadAbonada;
      }
      return 0;
    }
  },
  watch: {
    mostrar(mostrado) {
      if (mostrado) {
        this.obtenerAbonos();
        this.$nextTick(this.$refs.inputPago.focus);
      }
    }
  },
  data: () => ({
    cargandoAbonando: false,
    cantidadAbonada: 0,
    pago: 0,//Lo que el cliente paga, para calcular el cambio
    abonos: [],
    encabezados: [
      {
        text: "Monto",
        value: "Monto"
      },
      {
        text: "Fecha",
        value: "Fecha"
      },
      {
        text: "Usuario",
        value: "Usuario",
        sortable: false
      },
      {
        text: "Opciones",
        value: "Opciones",
        sortable: false
      }
    ]
  }),
  methods: {
    ocultarDialogo() {
      this.$emit("cerrar");
    },
    obtenerAbonos() {
      if (this.apartado.IdApartado) {
        HTTP_AUTH.get(
          `abonos/apartado/${encodeURIComponent(this.apartado.IdApartado)}`
        ).then(abonos => {
          this.abonos = abonos;
        });
      }
    },
    async imprimir(idAbono, idApartado) {
      await FUNCIONES.imprimirTicketAbono(idAbono, idApartado);
    },
    async abonar() {
      if (this.cantidadAbonada <= 0) return this.$emit("cantidadNegativa");
      if (this.cantidadAbonada > this.restante)
        return this.$emit("cantidadSuperior");
      this.cargandoAbonando = true;
      const resultados = await HTTP_AUTH.put("abono", {
        IdApartado: this.apartado.IdApartado,
        Monto: this.cantidadAbonada,
        Pago: this.pago,
      });
      console.log({ resultados });
      this.cargandoAbonando = false;
      if (this.cantidadAbonada === this.restante) this.$emit("liquidar");
      this.$emit("abonado");
      this.cantidadAbonada = 0;
      this.obtenerAbonos();
      await FUNCIONES.imprimirTicketAbono(resultados, this.apartado.IdApartado)

    }
  }
};
</script>
