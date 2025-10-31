<template>
  <v-container fluid grid-list-lg>
    <v-layout row wrap>
      <v-flex xs12 sm6 md3>
        <v-card color="blue lighten-2" dark>
          <v-card-title>
            <h1 class="headline">
              {{ reporteCaja.VentasContado | currency }}
              <v-icon color="blue lighten-5" class="icono-tarjeta"
                >attach_money</v-icon
              >
            </h1>
          </v-card-title>
          <v-card-text>
            <p class="title">Vendido al contado</p>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3>
        <v-card color="indigo lighten-2" dark>
          <v-card-title>
            <h1 class="headline">
              {{ (reporteCaja.Anticipos + reporteCaja.Abonos) | currency }}
              <v-icon color="indigo lighten-5" class="icono-tarjeta"
                >credit_card</v-icon
              >
            </h1>
          </v-card-title>
          <v-card-text>
            <p class="title">Anticipos y abonos</p>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3>
        <v-card color="amber lighten-1" dark>
          <v-card-title>
            <h1 class="headline">
              {{ reporteCaja.Ingresos | currency }}
              <v-icon color="amber lighten-5" class="icono-tarjeta"
                >trending_up</v-icon
              >
            </h1>
          </v-card-title>
          <v-card-text>
            <p class="title">Entradas de efectivo</p>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3>
        <v-card color="red darken-1" dark>
          <v-card-title>
            <h1 class="headline">
              {{ reporteCaja.Egresos | currency }}
              <v-icon color="red lighten-5" class="icono-tarjeta"
                >trending_down</v-icon
              >
            </h1>
          </v-card-title>
          <v-card-text>
            <p class="title">Salidas de efectivo</p>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12>
        <!-- <Publicidad></Publicidad> -->
      </v-flex>
      <v-flex xs12 sm12 md6>
        <v-card>
          <v-card-text>
            <ventas-mes
              ref="graficaVentasMes"
              style="max-height: 400px"
            ></ventas-mes>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm12 md6>
        <v-card>
          <v-card-text>
            <ventas-anio
              ref="graficaVentasAnio"
              style="max-height: 400px"
            ></ventas-anio>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm4>
        <v-card>
          <v-card-title>
            <h1 class="headline">Productos más vendidos</h1>
          </v-card-title>
          <v-card-text>
            <productos-mas-vendidos
              :productos="productosMasVendidos"
            ></productos-mas-vendidos>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm4>
        <v-card>
          <v-card-title>
            <h1 class="headline">Clientes frecuentes</h1>
          </v-card-title>
          <v-card-text> Muy pronto... </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm4>
        <v-card>
          <v-card-title>
            <h1 class="headline">Top cajeros</h1>
          </v-card-title>
          <v-card-text> Muy pronto... </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>
<style>
.icono-tarjeta {
  position: absolute;
  right: 5px;
  top: 5px;
  font-size: 70px;
  z-index: 1;
}
</style>

<script>
import VentasMes from "./Graficas/VentasMes";
import VentasAnio from "./Graficas/VentasAnio";
import { FUNCIONES } from "../funciones";
import { HTTP_AUTH } from "../http-common";
import { EventBus } from "../main";
import ProductosMasVendidos from './Reportes/ProductosMasVendidos'
import Publicidad from "./Publicidad";

export default {
  beforeMount() {
    this.consultarDatosParaEscritorio();
    EventBus.$emit("ponerTitulo", "Escritorio");
    this.$nextTick(() => {
      this.$refs.graficaVentasAnio.agregarAnio(FUNCIONES.esteAnioComoCadena());
      this.$refs.graficaVentasMes.setAnioYMes(FUNCIONES.esteAnioComoCadena(), FUNCIONES.esteMesComoCadena());
    });
  },
  data: () => ({
    reporteCaja: {},
    productosMasVendidos: []
  }),
  methods: {
    consultarDatosParaEscritorio() {
      let hoy = new Date();
      let inicio = FUNCIONES.componerFechaParaInicio(FUNCIONES.formatearFecha(hoy)),
        fin = FUNCIONES.componerFechaParaFin(FUNCIONES.formatearFecha(hoy));
      HTTP_AUTH.get(`reporte/caja/${inicio}/${fin}`)
        .then(reporteCaja => {
          this.reporteCaja = reporteCaja;
        })
        .then(() => {
          return HTTP_AUTH.get(`productos/mas/vendidos/${inicio}/${fin}`).then(
            productosMasVendidos => {
              this.productosMasVendidos = productosMasVendidos;
            }
          );
        });
    }
  },
  components: {
    Publicidad,
    VentasMes,
    VentasAnio,
    ProductosMasVendidos
  }
};
</script>
