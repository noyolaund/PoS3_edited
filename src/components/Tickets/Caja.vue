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
      <v-img class="text-xs-center" :src="require('@/assets/inicio/logo.png')">
      </v-img>
      <div class="text-xs-center">
        <p>
          <strong>Estado de caja</strong>
          <br />
          <strong>Tipo: </strong>
          {{
            usuario.Numero
              ? `De #${usuario.Numero} ${usuario.Nombre} `
              : "General"
          }}
          <br />
          <strong>Desde: </strong> {{ $route.query.fechaInicio }}
          <br />
          <strong>Hasta: </strong> {{ $route.query.fechaFin }}
          <br />
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
          Impreso por: <strong>{{ usuarioLogueado.Nombre }}</strong>
        </p>
        <p>
          {{ hoy | fechaExpresiva }}
        </p>
        <div class="text-xs-center">
          <v-flex xs9 offset-xs1 class="text-xs-right con-borde-separador">
            <br />
          </v-flex>
        </div>
      </div>
      <v-layout wrap row>
        <v-flex xs12 class="text-xs-left">Ventas al contado</v-flex>
        <v-flex xs12 class="text-xs-right con-borde-inferior">{{
          estado.VentasContado | currency
        }}</v-flex>
        <v-flex xs12 class="text-xs-left">Anticipo de apartados</v-flex>
        <v-flex xs12 class="text-xs-right con-borde-inferior">{{
          estado.Anticipos | currency
        }}</v-flex>
        <v-flex xs12 class="text-xs-left">Abonos</v-flex>
        <v-flex xs12 class="text-xs-right con-borde-inferior">{{
          estado.Abonos | currency
        }}</v-flex>
        <v-flex xs12 class="text-xs-left">Ingresos</v-flex>
        <v-flex xs12 class="text-xs-right con-borde-inferior">{{
          estado.Ingresos | currency
        }}</v-flex>
        <v-flex xs12 class="text-xs-left">Egresos</v-flex>
        <v-flex xs12 class="text-xs-right con-borde-inferior">{{
          estado.Egresos | currency
        }}</v-flex>
      </v-layout>
      <div class="text-xs-center">
        <v-flex xs9 offset-xs1 class="text-xs-right con-borde-separador">
          <br />
        </v-flex>
      </div>
      <div class="text-xs-right">
        <p>
          <strong>Total</strong>
          {{
            (estado.VentasContado +
              estado.Anticipos +
              estado.Abonos +
              estado.Ingresos -
              estado.Egresos)
              | currency
          }}
        </p>
      </div>
      <div class="text-xs-center">
        <pie></pie>
      </div>
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
import { HTTP, HTTP_AUTH } from "../../http-common";
import { EventBus } from "../../main";
import Pie from "./Pie";
import Encabezado from "./Encabezado";
import { TimeoutOcultarMenuTickets } from "../../constantes";

export default {
  name: "TicketDeCaja",
  components: { Encabezado, Pie },
  data: () => ({
    cargando: false,
    ajustes: {},
    estado: {},
    usuario: {},
    hoy: null,
    usuarioLogueado: {},
  }),
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Impresión de ticket");
    this.obtenerDetallesParaEstadoDeCaja(
      this.$route.query.fechaInicio,
      this.$route.query.fechaFin,
      this.$route.query.usuario
    );
  },
  beforeRouteUpdate(ruta) {
    this.obtenerDetallesParaEstadoDeCaja(
      ruta.query.fechaInicio,
      ruta.query.fechaFin,
      ruta.query.usuario
    );
  },
  methods: {
    obtenerDetallesParaEstadoDeCaja(fechaInicio, fechaFin, idUsuario) {
      let ruta = `reporte/caja/${fechaInicio}/${fechaFin}`;
      if (idUsuario) {
        ruta += `/${idUsuario}`;
      }
      this.cargando = true;
      HTTP_AUTH.get("ajustes/empresa")
        .then((ajustes) => {
          this.ajustes = ajustes;
        })
        .then(() =>
          HTTP_AUTH.get(ruta).then((estado) => (this.estado = estado))
        )
        .then(() =>
          HTTP.get("fechaYHora").then((fechaYHora) => (this.hoy = fechaYHora))
        )
        .then(() =>
          HTTP_AUTH.get("usuario/logueado").then(
            (usuario) => (this.usuarioLogueado = usuario)
          )
        )
        .then(() => {
          if (idUsuario) {
            HTTP_AUTH.get(`usuario/caja/${idUsuario}`)
              .then((usuario) => (this.usuario = usuario))
              .finally(() => (this.cargando = false));
          } else {
            this.cargando = false;
          }
        });
    },
    imprimir() {
      if (this.cargando) return;
      EventBus.$emit("ocultarMenu");
      // Esperar a que el menú esté oculto
      setTimeout(() => {
        let tituloOriginal = document.title;
        document.title = `Estado de caja`;
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
