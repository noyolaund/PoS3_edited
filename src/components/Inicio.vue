<template>
  <v-container grid-list-lg fluid>
    <v-snackbar :timeout="5000" :top="true" v-model="snackbarStock">
      Advertencia: Hay {{ cantidadEnStock }} productos con existencia menor al
      stock
      <v-btn
        flat
        color="pink"
        @click.native="irAProductosEnStockYOcultarSnackbar()"
        >ACCIÓN</v-btn
      >
    </v-snackbar>
    <v-layout row wrap fill-height>
      <v-flex xs12>
        <h1 class="display-3 text-xs-center">Bienvenido. Elija una opción</h1>
        <p>
          Recuerde que siempre puede explorar el menú completo haciendo click en
          la parte superior izquierda
        </p>
      </v-flex>
      <v-flex xs12>
        <!-- <publicidad></publicidad> -->
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('Vender')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/bolsa-de-la-compra.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Realizar una venta</div>
              <span class="grey--text"
                >Hacer una venta al contado o un apartado</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('Inventario')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/que-hacer.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Ver inventario</div>
              <span class="grey--text"
                >Registrar, eliminar o actualizar detalles productos</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('Clientes')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/red.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Mis clientes</div>
              <span class="grey--text"
                >Registrar, eliminar o actualizar detalles clientes</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('ReporteApartados')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/dinero.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Abonos</div>
              <span class="grey--text"
                >Registrar o ver abonos de un apartado</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('Caja')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/monedas.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Caja</div>
              <span class="grey--text"
                >Registrar una salida o entrada de efectivo</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('Escritorio')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/salpicadero.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Escritorio</div>
              <span class="grey--text"
                >Resumen de todos los movimientos de su tienda</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('Graficas')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/estadistica.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Gráficas</div>
              <span class="grey--text"
                >Compare ventas, vea los productos más vendidos y mucho
                más</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card class="text-xs-center" style="min-height: 250px">
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/tareas.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Reportes</div>
              <span class="grey--text"
                >Ver el <a @click="irA('ReporteCaja')">reporte de caja</a>, de
                <a @click="irA('ReporteVentasContado')">ventas al contado</a> o
                de <a @click="irA('ReporteApartados')">apartados</a></span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('Usuarios')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/escudo.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Usuarios</div>
              <span class="grey--text"
                >Administrar usuarios (o cajeros) así como sus permisos</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('ImprimirCodigosDeBarras')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/codigo-de-barras.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Códigos de barras</div>
              <span class="grey--text"
                >Imprimir códigos de barras de productos</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex xs12 sm6 md3 lg2>
        <v-card
          class="cursor-manita text-xs-center"
          style="min-height: 250px"
          @click.native="irA('Ajustes')"
        >
          <v-img
            class="text-xs-center"
            :src="require('@/assets/inicio/ajustes.png')"
            style="max-height: 100px; width: 100px; margin: auto"
            height="200px"
          >
          </v-img>
          <v-card-title primary-title>
            <div>
              <div class="headline">Ajustes</div>
              <span class="grey--text"
                >Cambiar los datos de su tienda, configurar la impresora,
                etcétera</span
              >
            </div>
          </v-card-title>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { EventBus } from "../main";
import { HTTP_AUTH } from "../http-common";
import Publicidad from "./Publicidad";

export default {
  components: { Publicidad },
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Inicio");
    this.comprobarProductosEnStock();
  },
  methods: {
    irA(donde) {
      this.$router.push({ name: donde });
    },
    irAProductosEnStockYOcultarSnackbar() {
      this.snackbarStock = false;
      this.$router.push({ name: "ReporteStock" });
    },
    comprobarProductosEnStock() {
      HTTP_AUTH
        .get("usuario/logueado")
        .then(datosUsuario => {
          EventBus.$emit("ponerDatosUsuario", datosUsuario)
        })
        .finally(() => {
          HTTP_AUTH
            .get("conteo/productos/stock")
            .then(cantidad => {
              if (!isNaN(cantidad) && cantidad > 0) {
                this.snackbarStock = true;
                this.cantidadEnStock = cantidad;
              }

            })
        });
    }
  },
  data: () => ({
    snackbarStock: false,
    cantidadEnStock: 0,
  })
};

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.card:hover {
  background-color: #e0d7dd;
}

div.card__media__background:hover {
  background-color: yellow;
  -moz-transform: scale(1.1);
  -webkit-transform: scale(1.1);
  transform: scale(1.1);
}

h1,
h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
