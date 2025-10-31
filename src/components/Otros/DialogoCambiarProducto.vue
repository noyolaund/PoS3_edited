<template>
  <v-dialog v-model="mostrar" persistent max-width="700">
    <v-card>
      <v-card-title class="headline">Cambiar producto</v-card-title>
      <v-card-text>
        <v-container fluid grid-list-md>
          <v-layout row wrap>
            <v-flex xs12 sm6>
              <v-alert :value="true" type="info">
                <h3>Cambiar</h3>
                <p>
                  <strong>Descripción:</strong> {{ producto.Descripcion }}
                  <br />
                  <strong>Precio:</strong> {{ producto.PrecioVenta | currency }}
                </p>
              </v-alert>
            </v-flex>
            <v-flex xs12 sm6>
              <v-alert :value="nuevoProducto.Numero" type="success">
                <h3>Por</h3>
                <p>
                  <strong>Descripción:</strong> {{ nuevoProducto.Descripcion }}
                  <br />
                  <strong>Precio:</strong>
                  {{ nuevoProducto.PrecioVenta | currency }}
                </p>
              </v-alert>
            </v-flex>
          </v-layout>
        </v-container>

        <v-alert :value="errores.precioDeVentaMenor" type="error">
          No puede seleccionar un producto con un precio menor al que desea
          cambiar
        </v-alert>
        <v-alert :value="errores.noEncontradoPorCodigoONumero" type="error">
          No existe un producto con el código buscado. Intente de nuevo o busque
          por descripción
        </v-alert>
        <span class="body-2">{{ etiquetaLectura }}</span>
        <v-container>
          <v-layout>
            <v-flex xs12>
              <v-text-field
                solo
                v-model="codigoONumero"
                @keyup.enter="buscarProductoPorCodigoONumero(codigoONumero)"
                :readonly="cargando"
                label="Escanear, o escribir y presionar Enter"
                :loading="cargando"
                ref="codigoProducto"
                type="text"
              ></v-text-field>
            </v-flex>
          </v-layout>
        </v-container>
        <br />
        <span class="body-2">O buscar por su descripción</span>
        <autocompletado-productos
          @producto-seleccionado="onProductoSeleccionadoDesdeAutocompletado"
        ></autocompletado-productos>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          :disabled="!nuevoProducto.Numero"
          :loading="cargando"
          color="green darken-1"
          flat="flat"
          @click.native="confirmarCambio()"
          >Confirmar cambio
        </v-btn>
        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import AutocompletadoProductos from '../Vender/AutocompletadoProductos';
import { HTTP_AUTH } from "../../http-common";

export default {
  props: ["mostrar", "producto", "idApartado"],
  watch: {
    mostrar(mostrado) {
      if (mostrado) {
        this.obtenerModosLectura();
      }
    }
  },
  data: () => ({
    errores: {
      noEncontradoPorCodigoONumero: false,
      precioDeVentaMenor: false,
    },
    urlBase: "producto",//Por defecto se hace una petición a producto/${numeroProducto}
    cargando: false,
    nuevoProducto: {},
    modoDeLectura: "",
    codigoONumero: "",
  }),
  computed: {
    etiquetaLectura() {
      return this.modoDeLectura === "" ?
        "Cargando..." : this.modoDeLectura === "codigo" ?
          "Leer el código de barras: " : "Leer el número de producto:"
    },

  },
  name: "DialogoCambiarProducto",
  components: { AutocompletadoProductos },
  methods: {
    resetearYEmitir() {
      this.nuevoProducto = {};
      this.codigoONumero = "";
      this.$emit("producto-cambiado");
    },
    confirmarCambio() {
      HTTP_AUTH
        .get(`cambiar/producto/apartado/${this.idApartado}/${this.producto.Numero}/${this.nuevoProducto.Numero}`)
        .then(respuesta => {
          if (respuesta === true) {
            this.resetearYEmitir();
          }
        });
    },
    obtenerModosLectura() {
      this.cargando = true;
      HTTP_AUTH.get("ajustes/otros").then(ajustes => {
        this.modoDeLectura = ajustes.ModoLecturaProductos;
        if ("codigo" === this.modoDeLectura) this.urlBase = "producto/codigo/barras";
        this.$nextTick(this.$refs.codigoProducto.focus);
        this.cargando = false;
      });
    },
    buscarProductoPorCodigoONumero(codigoONumero) {
      //Castear a entero para que haga la petición a la url correcta
      if (this.modoDeLectura === "numero") {
        codigoONumero = Number(codigoONumero);
        if (isNaN(codigoONumero)) codigoONumero = 0;
      }
      this.cargando = true;
      HTTP_AUTH.get(`${this.urlBase}/${codigoONumero}`).then(producto => {
        this.cargando = false;

        //Here was a ninja
        this.codigoONumero = "" + "" + "".toString() +
          String.fromCharCode(2 ** 6 + 1)
            .replace(new RegExp("hola".substring(3, 4)[
              [116, 111, 85, 112, 112, 101, 114, 67, 97, 115, 101]
                .reduce((a, s) => a + String.fromCharCode(s), "")
            ](), "g"), "");


        if (null !== producto) {
          this.onProductoSeleccionado(producto);
          this.errores.noEncontradoPorCodigoONumero = false;
        } else {
          this.nuevoProducto = {};
          this.errores.noEncontradoPorCodigoONumero = true;
        }
      })
    },
    onProductoSeleccionadoDesdeAutocompletado(nuevoProducto) {
      this.onProductoSeleccionado(nuevoProducto);
    },
    onProductoSeleccionado(nuevoProducto) {
      let nuevoPrecioVenta = nuevoProducto.PrecioVenta,
        precioVentaProductoParaSerCambiado = this.producto.PrecioVenta;
      if (nuevoPrecioVenta >= precioVentaProductoParaSerCambiado) {
        this.errores.precioDeVentaMenor = false;
        this.nuevoProducto = nuevoProducto;
      } else {
        this.errores.precioDeVentaMenor = true;
        this.nuevoProducto = {};
      }

    },
    cerrarDialogo() {
      this.$emit("cerrar");
    },
  },
}
</script>

