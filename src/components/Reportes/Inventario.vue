<template>
  <v-container fluid grid-list-lg>
    <v-layout row wrap>
      <v-flex xs12 sm6 md3>
        <v-card height="100%" color="blue lighten-2" dark>
          <v-card-title class="sin-padding-inferior">
            <h1 class="headline">{{ reporte.CantidadProductos }}</h1>
          </v-card-title>
          <v-card-text class="sin-padding-superior">
            <p
              class="subheading"
              title="Cantidad de productos distintos que hay en el inventario"
            >
              Productos en inventario
            </p>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3>
        <v-card height="100%" color="indigo lighten-2" dark>
          <v-card-title class="sin-padding-inferior">
            <h1 class="headline">{{ reporte.PrecioVenta | currency }}</h1>
          </v-card-title>
          <v-card-text class="sin-padding-superior">
            <p
              class="subheading"
              title="Sumatoria del precio de venta multiplicado por la existencia de cada producto"
            >
              Precio del inventario
            </p>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3>
        <v-card height="100%" color="amber lighten-1" dark>
          <v-card-title class="sin-padding-inferior">
            <h1 class="headline">{{ reporte.PrecioCompra | currency }}</h1>
          </v-card-title>
          <v-card-text class="sin-padding-superior">
            <p
              class="subheading"
              title="Sumatoria del precio de compra multiplicado por la existencia de cada producto"
            >
              Costo del inventario
            </p>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3>
        <v-card height="100%" color="green darken-1" dark>
          <v-card-title class="sin-padding-inferior">
            <h1 class="headline">
              {{ (reporte.PrecioVenta - reporte.PrecioCompra) | currency }}
            </h1>
          </v-card-title>
          <v-card-text class="sin-padding-superior">
            <p
              class="subheading"
              title="Precio del inventario - Costo del inventario"
            >
              Utilidad total
            </p>
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
    </v-layout>
  </v-container>
</template>
<style>
.sin-padding-inferior {
  padding-bottom: 0.1em !important;
}

.sin-padding-superior {
  padding-top: 0.1em !important;
}
</style>
<script>
import { EventBus } from "../../main";
import { HTTP_AUTH } from "../../http-common";
import ProductosMasVendidos from './ProductosMasVendidos'
import ProductosMenosVendidos from './ProductosMenosVendidos'
import ProductosNuncaVendidos from './ProductosNuncaVendidos'
import SeleccionadorFechas from './SeleccionadorFechas'

export default {
  name: "Inventario",
  components: { ProductosMasVendidos, ProductosMenosVendidos, ProductosNuncaVendidos, SeleccionadorFechas },
  data: () => ({
    reporte: {},
    productosMasVendidos: [],
    productosMenosVendidos: [],
    productosNuncaVendidos: [],
  }),
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Reporte de inventario");
    this.obtenerReporte();
  },
  methods: {
    obtenerProductosNuncaVendidos() {
      HTTP_AUTH.get("productos/nunca/vendidos/al/contado").then(productos => {
        this.productosNuncaVendidos = productos;
      });
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
    obtenerReporte() {
      HTTP_AUTH.get("reporte/inventario")
        .then(reporte => {
          this.reporte = reporte;
        })
    }
  },
}
</script>

