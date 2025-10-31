<template>
  <v-layout row wrap>
    <v-flex xs12>
      <seleccionador-fechas
        @cambio="comprobarFechasYRefrescarSiEsNecesario"
      ></seleccionador-fechas>
    </v-flex>
    <v-flex xs12 sm3>
      <v-radio-group label="Tipo de reporte" v-model="tipoReporte">
        <v-radio label="General" value="general"></v-radio>
        <v-radio :loading="true" label="Por usuario" value="usuario"></v-radio>
      </v-radio-group>
    </v-flex>
    <v-flex xs12 sm9>
      <v-slide-y-transition>
        <v-select
          @change="onUsuarioCambiado"
          v-show="tipoReporte === 'usuario'"
          :loading="cargandoListaDeUsuarios"
          :items="usuarios"
          v-model="usuarioSeleccionado"
          label="Seleccione un usuario"
          item-value="Numero"
          return-object
        >
          <template slot="selection" slot-scope="data">
            <v-subheader
              >{{ data.item.Nombre }} - #{{ data.item.Numero }}</v-subheader
            >
          </template>
          <template slot="item" slot-scope="data">
            <span>{{ data.item.Nombre }} - #{{ data.item.Numero }}</span>
          </template>
        </v-select>
      </v-slide-y-transition>
    </v-flex>
    <v-flex xs12>
      <h1>
        <span class="headline">{{ reporte.VentasContado | currency }}</span>
        <span class="subheading">Ventas al contado</span>
      </h1>
    </v-flex>
    <v-flex xs12>
      <h1>
        <span class="headline">&plus;{{ reporte.Anticipos | currency }}</span>
        <span class="subheading">Anticipo de apartados</span>
      </h1>
    </v-flex>
    <v-flex xs12>
      <h1>
        <span class="headline">&plus;{{ reporte.Abonos | currency }}</span>
        <span class="subheading">Abonos</span>
      </h1>
    </v-flex>
    <v-flex xs12>
      <h1>
        <span class="headline">&plus;{{ reporte.Ingresos | currency }}</span>
        <span class="subheading">Ingresos</span>
      </h1>
    </v-flex>
    <v-flex xs12>
      <h1>
        <span class="headline">&minus;{{ reporte.Egresos | currency }}</span>
        <span class="subheading">Egresos</span>
      </h1>
    </v-flex>
    <v-flex xs12 sm4>
      <v-divider></v-divider>
    </v-flex>
    <v-flex xs12>
      <h1>
        <span class="headline">{{ total | currency }}</span>
        <span class="subheading">Total en caja</span>
      </h1>
    </v-flex>
    <v-flex xs12>
      <v-btn :loading="cargando" @click="imprimir" small color="success"
        >Imprimir</v-btn
      >
    </v-flex>
  </v-layout>
</template>
<script>
import { HTTP_AUTH } from "../../../http-common";
import SeleccionadorFechas from "../../Reportes/SeleccionadorFechas";
import { EventBus } from "../../../main";
import { FUNCIONES } from '../../../funciones';

export default {
  watch: {
    tipoReporte() {
      this.$nextTick(this.refrescarReporteDependiendoDelTipoSeleccionado);
    },
  },
  computed: {
    total() {
      if (
        null === this.reporte.VentasContado ||
        undefined === this.reporte.VentasContado
      )
        return 0;
      return (
        this.reporte.VentasContado +
        this.reporte.Anticipos +
        this.reporte.Abonos +
        this.reporte.Ingresos -
        this.reporte.Egresos
      );
    }
  },
  data: () => ({
    cargandoListaDeUsuarios: false,
    usuarios: [],
    usuarioSeleccionado: {},
    tipoReporte: "general",
    cargando: false,
    ultimaFechaInicio: null,
    ultimaFechaFin: null,
    reporte: {
      VentasContado: null, Anticipos: null, Abonos: null, Ingresos: null, Egresos: null
    }
  }),
  components: { SeleccionadorFechas },
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Reporte de caja");
    this.obtenerListaDeUsuariosSiEstaVacia();
  },
  methods: {
    onUsuarioCambiado(usuario) {
      this.refrescarReporteConUsuario({
        inicio: this.ultimaFechaInicio, fin: this.ultimaFechaFin,
        idUsuario: usuario.Numero
      });
    },
    obtenerListaDeUsuariosSiEstaVacia() {
      if (this.usuarios.length <= 0) {
        this.obtenerListaDeUsuarios();
      }
    },
    obtenerListaDeUsuarios() {
      this.cargandoListaDeUsuarios = true;
      HTTP_AUTH.get("usuarios").then(usuarios => {
        this.cargandoListaDeUsuarios = false;
        this.usuarios = usuarios;
      });
    },
    async imprimir() {
      if (this.tipoReporte === "usuario" && this.usuarioSeleccionado.Numero) {
        await FUNCIONES.imprimirReporteCaja(this.ultimaFechaInicio, this.ultimaFechaFin, this.usuarioSeleccionado.Numero);
      } else {
        await FUNCIONES.imprimirReporteCaja(this.ultimaFechaInicio, this.ultimaFechaFin, null);
      }
    },
    comprobarFechasYRefrescarSiEsNecesario(fechas) {
      this.refrescarReporteDependiendoDelTipoSeleccionado(fechas);
    },
    refrescarReporteDependiendoDelTipoSeleccionado(fechas) {
      if (fechas) {
        let { inicio, fin } = fechas;
        if (inicio && fin) {
          this.ultimaFechaInicio = inicio;
          this.ultimaFechaFin = fin;
        }
      }
      if (this.tipoReporte === "usuario" && this.usuarioSeleccionado.Numero) {
        this.refrescarReporteConUsuario({
          inicio: this.ultimaFechaInicio, fin: this.ultimaFechaFin,
          idUsuario: this.usuarioSeleccionado.Numero
        });
      } else {
        this.refrescarReporte({ inicio: this.ultimaFechaInicio, fin: this.ultimaFechaFin });
      }
    },
    refrescarReporte({ inicio, fin }) {
      this.cargando = true;
      this.ultimaFechaInicio = inicio;
      this.ultimaFechaFin = fin;
      HTTP_AUTH.get(`reporte/caja/${inicio}/${fin}`).then(reporte => {
        this.reporte = reporte;
        this.cargando = false;
      });
    },
    refrescarReporteConUsuario({ inicio, fin, idUsuario }) {
      this.cargando = true;
      this.ultimaFechaInicio = inicio;
      this.ultimaFechaFin = fin;
      HTTP_AUTH.get(`reporte/caja/${inicio}/${fin}/${idUsuario}`).then(reporte => {
        this.reporte = reporte;
        this.cargando = false;
      });
    },
  }
};
</script>
