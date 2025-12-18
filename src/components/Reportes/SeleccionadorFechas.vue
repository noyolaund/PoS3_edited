<template>
  <v-container fluid fluid-xs>
    <v-layout row wrap>
        <v-flex xs12 sm3>
          <v-dialog
            ref="menuFechaInicio"
            v-model="menuFechaInicio"
            :return-value.sync="fechaInicio"
            persistent
            lazy
            full-width
            width="290px"
          >
            <v-text-field
              slot="activator"
              v-model="fechaInicio"
              label="Ver desde"
              prepend-icon="event"
              readonly
              @click.native="menuFechaInicio = true"
            ></v-text-field>
            <v-date-picker locale="es-419" v-model="fechaInicio" scrollable>
              <v-spacer></v-spacer>
              <v-btn flat color="primary" @click="menuFechaInicio = false"
                >Cancelar</v-btn
              >
              <v-btn flat color="primary" @click="guardarFechaInicio()"
                >OK</v-btn
              >
            </v-date-picker>
          </v-dialog>
        </v-flex>
        <v-flex xs12 sm3>
          <v-dialog
            ref="menuFechaFin"
            v-model="menuFechaFin"
            :return-value.sync="fechaFin"
            persistent
            lazy
            full-width
            width="290px"
          >
            <v-text-field
              slot="activator"
              v-model="fechaFin"
              label="Hasta"
              prepend-icon="event"
              readonly
              @click.native="menuFechaFin = true"
            ></v-text-field>
            <v-date-picker locale="es-419" v-model="fechaFin" scrollable>
              <v-spacer></v-spacer>
              <v-btn flat color="primary" @click="menuFechaFin = false"
                >Cancelar</v-btn
              >
              <v-btn flat color="primary" @click="guardarFechaFin()">OK</v-btn>
            </v-date-picker>
          </v-dialog>
        </v-flex>
        <v-flex xs12 sm6>
          <v-btn
            @click="hoy()"
            :flat="seleccionado === 'hoy'"
            small
            :color="seleccionado === 'hoy' ? 'success' : 'primary'"
            >Hoy
            <v-icon v-show="seleccionado === 'hoy'">check</v-icon>
          </v-btn>
          <v-btn
            @click="semana()"
            :flat="seleccionado === 'semana'"
            small
            :color="seleccionado === 'semana' ? 'success' : 'primary'"
            >Esta semana
            <v-icon v-show="seleccionado === 'semana'">check</v-icon>
          </v-btn>
          <v-btn
            @click="mes()"
            :flat="seleccionado === 'mes'"
            small
            :color="seleccionado === 'mes' ? 'success' : 'primary'"
            >Este mes
            <v-icon v-show="seleccionado === 'mes'">check</v-icon>
          </v-btn>
        </v-flex>
      </v-layout>
  </v-container>
</template>
<script>
const DIA_EN_MILISEGUNDOS = 864e5; // 1000 * 60 * 60 * 24
import { FUNCIONES } from "../../funciones"

export default {
  beforeMount() {
    let hoy = new Date();
    this.fechaInicio = FUNCIONES.formatearFecha(hoy);
    this.fechaFin = FUNCIONES.formatearFecha(hoy);
    this.hoy();
  },
  methods: {
    hoy() {
      /**
       * Pone las fechas en los calendarios y consulta las fechas de hoy
       */
      let hoy = new Date();
      this.fechaInicio = FUNCIONES.formatearFecha(hoy);
      this.fechaFin = FUNCIONES.formatearFecha(hoy);
      this.onFechasCambiadas();
      this.seleccionado = "hoy";
    },
    semana() {
      /**
       * Pone las fechas en los calendarios y consulta las fechas de la semana
       */
      //Dependiendo de la fecha de hoy, establece la fecha al lunes de esta semana
      let lunes = new Date();
      lunes.setDate(
        lunes.getDate() - lunes.getDay() + (lunes.getDay() === 0 ? 0 : 1)
      );
      //Ahora el próximo domingo...
      let domingo = new Date();
      domingo.setTime(lunes.getTime() + DIA_EN_MILISEGUNDOS * 6);
      this.fechaInicio = FUNCIONES.formatearFecha(lunes);
      this.fechaFin = FUNCIONES.formatearFecha(domingo);
      this.onFechasCambiadas();
      this.seleccionado = "semana";
    },
    mes() {
      /**
       * Pone las fechas en los calendarios y consulta las fechas del mes
       */
      //Dependiendo de la fecha de hoy, establece la fecha al 1
      let primerDia = new Date();
      primerDia.setDate(1);
      //Ahora calcula la fecha dentro de un mes pero le resta un día. Robado de https://stackoverflow.com/questions/222309/calculate-last-day-of-month-in-javascript
      let ultimoDia = new Date(
        primerDia.getFullYear(),
        primerDia.getMonth() + 1,
        0
      );
      this.fechaInicio = FUNCIONES.formatearFecha(primerDia);
      this.fechaFin = FUNCIONES.formatearFecha(ultimoDia);
      this.onFechasCambiadas();
      this.seleccionado = "mes";
    },
    guardarFechaInicio() {
      this.seleccionado = "";
      this.$refs.menuFechaInicio.save(this.fechaInicio);
      this.onFechasCambiadas();
    },
    guardarFechaFin() {
      this.seleccionado = "";
      this.$refs.menuFechaFin.save(this.fechaFin);
      this.onFechasCambiadas();
    },
    onFechasCambiadas() {
      this.$emit("cambio", {
        inicio: FUNCIONES.componerFechaParaInicio(this.fechaInicio),
        fin: FUNCIONES.componerFechaParaFin(this.fechaFin)
      });
    }
  },
  data: () => ({
    mostrar: false,
    seleccionado: '',
    fechaInicio: null,
    fechaFin: null,
    menuFechaInicio: false,
    menuFechaFin: false
  })
};
</script>
