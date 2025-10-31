<template>
  <v-layout class="ticket" wrap row>
    <v-flex xs12 class="hidden-print-only">
      <encabezado></encabezado>
      <v-btn small color="success" @click="volver()">
        <v-icon>arrow_back</v-icon>
        Volver
      </v-btn>
    </v-flex>
    <v-flex xs12>
      <br />
      <!-- <v-img class="text-xs-center" :src="require('@/assets/inicio/logo.png')">
      </v-img> -->
      <div class="text-xs-center">
        <p>
          <strong>Ticket de venta #{{ venta.Numero }}</strong>
        </p>
        <p v-if="ajustes.Nombre">
          {{ ajustes.Nombre }}
        </p>
        <p v-if="ajustes.Direccion">
          {{ ajustes.Direccion }}
        </p>
        <p v-if="ajustes.Telefono">Teléfono: {{ ajustes.Telefono }}</p>
        <br />
        <p>
          {{ venta.Fecha | fechaExpresiva }}
        </p>
        <p>
          Lo atendió: <strong>{{ venta.Usuario.Nombre }}</strong>
        </p>
        <p>
          Cliente: <strong>{{ venta.Cliente.Nombre }}</strong>
        </p>
        <div class="text-xs-center">
          <v-flex xs9 offset-xs1 class="text-xs-right con-borde-separador">
            <br />
          </v-flex>
        </div>
      </div>
      <v-layout wrap row>
        <template v-for="producto in venta.Productos">
          <v-flex xs12 class="text-xs-left">{{ producto.Descripcion }}</v-flex>
          <v-flex xs12 class="text-xs-right con-borde-inferior"
            >{{ producto.Cantidad }} x {{ producto.PrecioVenta | currency }}
            =
            {{ (producto.Cantidad * producto.PrecioVenta) | currency }}
          </v-flex>
        </template>
      </v-layout>
      <div class="text-xs-center">
        <v-flex xs9 offset-xs1 class="text-xs-right con-borde-separador">
          <br />
        </v-flex>
      </div>
      <div class="text-xs-right">
        <p><strong>Total</strong> {{ venta.Total | currency }}</p>
        <p><strong>Su pago</strong> {{ venta.Pago | currency }}</p>
        <p>
          <strong>Cambio</strong> {{ (venta.Pago - venta.Total) | currency }}
        </p>
      </div>
      <!-- <div class="text-xs-center">
        <p v-if="ajustes.MensajePersonal">
          <strong>{{ ajustes.MensajePersonal }}</strong>
        </p>
        <Pie></Pie>
      </div> -->
    </v-flex>
    <v-btn
      :loading="cargando"
      class="hidden-print-only"
      @click="imprimir()"
      fixed
      dark
      fab
      bottom
      fill-height
      slot="activator"
      right
      color="green"
    >
      <v-icon>print</v-icon>
    </v-btn>
  </v-layout>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";
import { EventBus } from "../../main";
import Encabezado from "./Encabezado";
import Pie from "./Pie";
export const TimeoutOcultarMenuTickets = 200;

export default {
  name: "TicketVentaContado",
  components: { Pie, Encabezado },
  beforeRouteUpdate(detallesRuta) {
    this.obtenerDetallesDeVenta(detallesRuta.params.idVenta);
  },
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Impresión de ticket");
    this.obtenerDetallesDeVenta(this.$route.params.idVenta);
  },
  data: () => ({
    cargando: false,
    venta: {
      Usuario: {},
      Cliente: {},
    },
    ajustes: {},
  }),
  methods: {
    obtenerDetallesDeVenta(idVenta) {
      if (!idVenta) return this.$router.go(-1);
      this.cargando = true;
      HTTP_AUTH.get("ajustes/empresa")
        .then((ajustes) => {
          this.ajustes = ajustes;
        })
        .then(() => {
          HTTP_AUTH.get(`venta/contado/${idVenta}`)
            .then((venta) => {
              this.venta = venta;
            })
            .finally(() => (this.cargando = false));
        });
    },
    imprimir() {
      if (this.cargando) return;
      EventBus.$emit("ocultarMenu");
      setTimeout(() => {
        let tituloOriginal = document.title;
        document.title = `Venta al contado #${this.venta.Numero}`;
        window.print();
        document.title = tituloOriginal;
        EventBus.$emit("mostrarMenu");
      }, TimeoutOcultarMenuTickets);
    },
    volver() {
      this.$router.go(-1);
    },
  },
};
</script>
