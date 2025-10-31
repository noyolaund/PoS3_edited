<template>
  <v-container fluid grid-list-lg>
    <v-layout row wrap>
      <v-flex xs12>
        <v-card>
          <v-card-title><h1 class="title">Ventas por año</h1></v-card-title>
          <v-card-text>
            <v-combobox
              v-model="aniosSeleccionados"
              :loading="cargandoAnios"
              :items="aniosParaGraficas"
              label="Seleccione uno o más años para comparar"
              chips
              multiple
            >
              <template slot="selection" slot-scope="data">
                <v-chip
                  :selected="data.selected"
                  :disabled="data.disabled"
                  :key="data.item"
                  class="chip--select-multi"
                  @input="data.parent.selectItem(data.item)"
                >
                  {{ data.item }}
                </v-chip>
              </template>
            </v-combobox>
            <ventas-anio
              v-show="aniosSeleccionados.length > 0"
              ref="graficaVentasAnio"
            ></ventas-anio>
            <v-alert
              v-show="aniosSeleccionados.length <= 0"
              value="true"
              type="info"
            >
              Seleccione uno o más años para mostrar la gráfica
            </v-alert>
          </v-card-text>
        </v-card>
      </v-flex>
      <!-- <v-flex xs12>
        <v-card>
          <v-card-text>
            <v-flex xs12>
               <Publicidad></Publicidad> 
            </v-flex>
          </v-card-text>
        </v-card>
      </v-flex> -->
      <v-flex xs12>
        <v-card>
          <v-card-title>
            <h1 class="title">Ventas por mes</h1>
          </v-card-title>
          <v-card-text>
            <v-layout row wrap>
              <v-flex xs12 sm6>
                <v-select
                  v-model="anioSeleccionadoParaGraficaPorMes"
                  :loading="cargandoAnios"
                  :items="aniosParaGraficas"
                  label="Seleccione un año"
                >
                </v-select>
              </v-flex>
              <v-flex xs12 sm6>
                <v-select
                  :loading="cargandoMeses"
                  v-show="anioSeleccionadoParaGraficaPorMes"
                  v-model="mesSeleccionadoParaGraficaPorMes"
                  :items="mesesParaGraficas"
                  return-object
                  item-text="etiqueta"
                  item-value="mes"
                  label="Ahora un mes"
                >
                </v-select>
              </v-flex>
            </v-layout>
            <ventas-mes
              ref="graficaVentasMes"
              v-show="mesSeleccionadoParaGraficaPorMes.mes !== undefined"
            ></ventas-mes>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12>
        <v-card>
          <v-card-title>
            <h1 class="title">Productos más vendidos</h1>
          </v-card-title>
          <v-card-text>
            <seleccionador-fechas
              @cambio="consultarProductosMasVendidos"
            ></seleccionador-fechas>
            <productos-mas-vendidos
              :productos="productosMasVendidos"
            ></productos-mas-vendidos>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12>
        <v-card>
          <v-card-title>
            <h1 class="title">Productos menos vendidos</h1>
          </v-card-title>
          <v-card-text>
            <seleccionador-fechas
              @cambio="consultarProductosMenosVendidos"
            ></seleccionador-fechas>
            <productos-menos-vendidos
              :productos="productosMenosVendidos"
            ></productos-menos-vendidos>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12>
        <v-card>
          <v-card-title>
            <h1 class="title">Productos nunca vendidos</h1>
          </v-card-title>
          <v-card-text>
            <p>Los productos que nunca han sido vendidos al contado</p>
            <productos-nunca-vendidos
              :productos="productosNuncaVendidos"
            ></productos-nunca-vendidos>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import VentasMes from './VentasMes'
import SeleccionadorFechas from '../Reportes/SeleccionadorFechas'
import ProductosMasVendidos from '../Reportes/ProductosMasVendidos'
import ProductosMenosVendidos from '../Reportes/ProductosMenosVendidos'
import ProductosNuncaVendidos from '../Reportes/ProductosNuncaVendidos'
import VentasAnio from './VentasAnio'
import { HTTP_AUTH } from "../../http-common";
import { EventBus } from "../../main";
import { MESES } from "../../constantes";
import Publicidad from "../Publicidad";

export default {
  data: () => ({
    cargandoMeses: false,
    cargandoAnios: false,
    aniosParaGraficas: [],
    aniosSeleccionados: [],
    mesesParaGraficas: [],
    anioSeleccionadoParaGraficaPorMes: null,
    mesSeleccionadoParaGraficaPorMes: {},
    productosMasVendidos: [],
    productosMenosVendidos: [],
    productosNuncaVendidos: [],
  }),
  watch: {
    aniosSeleccionados(nuevo) {
      this.$refs.graficaVentasAnio.setAnios(Array.from(nuevo));
    },
    anioSeleccionadoParaGraficaPorMes() {
      this.obtenerMesesParaGraficas();
    },
    mesSeleccionadoParaGraficaPorMes: {
      handler() {
        this.$refs.graficaVentasMes.setAnioYMes(this.anioSeleccionadoParaGraficaPorMes, this.mesSeleccionadoParaGraficaPorMes.mes);
      },
      deep: true
    },
  },
  methods: {
    obtenerProductosNuncaVendidos() {
      HTTP_AUTH.get("productos/nunca/vendidos/al/contado").then(productos => {
        this.productosNuncaVendidos = productos;
      });
    },
    obtenerAniosParaGraficas() {
      this.cargandoAnios = true;
      return HTTP_AUTH.get("anios/graficas/ventas/contado").then(anios => {
        this.cargandoAnios = false;
        this.aniosParaGraficas = anios;
      });
    },
    obtenerMesesParaGraficas() {
      if (this.anioSeleccionadoParaGraficaPorMes) {
        this.cargandoMeses = true;
        HTTP_AUTH.get(`meses/graficas/ventas/contado/anio/${this.anioSeleccionadoParaGraficaPorMes}`).then(meses => {
          meses.sort((a, b) => a - b);
          meses = meses.map(mes => ({
            mes: mes,
            etiqueta: MESES[parseInt(mes) - 1]
          }));
          this.cargandoMeses = false;
          this.mesesParaGraficas = meses;
        });
      }
    },
    consultarProductosMasVendidos({ inicio, fin }) {
      HTTP_AUTH.get(`productos/mas/vendidos/${inicio}/${fin}`).then(
        productosMasVendidos => {
          this.productosMasVendidos = productosMasVendidos;
        }
      );
    },
    consultarProductosMenosVendidos({ inicio, fin }) {
      HTTP_AUTH.get(`productos/menos/vendidos/${inicio}/${fin}`).then(
        productosMenosVendidos => {
          this.productosMenosVendidos = productosMenosVendidos;
        }
      );
    },
  },
  beforeMount() {
    this.obtenerAniosParaGraficas().then(() => {
      this.obtenerProductosNuncaVendidos();
    });
    EventBus.$emit("ponerTitulo", "Gráficas y estadísticas");
  },
  components: {
    Publicidad,
    VentasAnio,
    VentasMes,
    SeleccionadorFechas,
    ProductosMasVendidos,
    ProductosNuncaVendidos,
    ProductosMenosVendidos
  }
}
</script>

