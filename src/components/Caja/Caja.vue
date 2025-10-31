<template>
  <v-layout row wrap>
    <dialogo-ingreso
      @guardado="onIngresoGuardado"
      @cerrar="dialogos.ingreso = false"
      :mostrar="dialogos.ingreso"
    ></dialogo-ingreso>
    <dialogo-egreso
      @guardado="onEgresoGuardado"
      @cerrar="dialogos.egreso = false"
      :mostrar="dialogos.egreso"
    ></dialogo-egreso>
    <v-flex xs12>
      <v-btn @click="dialogos.ingreso = true" small color="success"
        >Registrar ingreso</v-btn
      >
      <v-btn @click="dialogos.egreso = true" small color="orange"
        >Registrar egreso</v-btn
      >
    </v-flex>
    <seleccionador-de-fechas
      @cambio="comprobarFechasYRefrescarSiEsNecesario"
    ></seleccionador-de-fechas>
    <v-flex xs12>
      <!-- <Publicidad></Publicidad> -->
    </v-flex>
    <v-flex xs12 sm6>
      <ingresos ref="ingresos"></ingresos>
    </v-flex>
    <v-flex xs12 sm6>
      <egresos ref="egresos"></egresos>
    </v-flex>
    <v-snackbar
      :timeout="5000"
      :top="true"
      :right="true"
      v-model="snackbars.ingresoGuardado"
    >
      Ingreso guardado
      <v-btn flat color="pink" @click.native="snackbars.ingresoGuardado = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="5000"
      :top="true"
      :right="true"
      v-model="snackbars.egresoGuardado"
    >
      Egreso guardado
      <v-btn flat color="pink" @click.native="snackbars.egresoGuardado = false"
        >OK</v-btn
      >
    </v-snackbar>
  </v-layout>
</template>
<script>
import SeleccionadorDeFechas from "../Reportes/SeleccionadorFechas";
import DialogoIngreso from "./Dialogos/Ingreso";
import DialogoEgreso from "./Dialogos/Egreso";
import Ingresos from "./Ingresos";
import Egresos from "./Egresos";
import { EventBus } from "../../main";
import Publicidad from "../Publicidad";

export default {
  components: {
    Publicidad,
    SeleccionadorDeFechas,
    DialogoIngreso,
    DialogoEgreso,
    Ingresos,
    Egresos
  },
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Caja");
  },
  methods: {
    comprobarFechasYRefrescarSiEsNecesario({ inicio, fin }) {
      if (inicio && fin) {
        this.ultimaFechaInicio = inicio;
        this.ultimaFechaFin = fin;
        this.consultarIngresosYEgresos(inicio, fin);
      }
    },
    consultarIngresosYEgresos(inicio, fin) {
      this.$nextTick(() => {
        this.$refs.ingresos.obtener(inicio, fin);
        this.$refs.egresos.obtener(inicio, fin);
      });
    },
    onIngresoGuardado() {
      this.consultarIngresosYEgresos(
        this.ultimaFechaInicio,
        this.ultimaFechaFin
      );
      this.dialogos.ingreso = false;
      this.snackbars.ingresoGuardado = true;
    },
    onEgresoGuardado() {
      this.consultarIngresosYEgresos(
        this.ultimaFechaInicio,
        this.ultimaFechaFin
      );
      this.dialogos.egreso = false;
      this.snackbars.egresoGuardado = true;
    }
  },
  data: () => ({
    ultimaFechaInicio: null,
    ultimaFechaFin: null,
    dialogos: {
      ingreso: false,
      egreso: false
    },
    snackbars: {
      ingresoGuardado: false,
      egresoGuardado: false
    }
  })
};
</script>
