<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Registrar producto</v-card-title>
      <v-card-text>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12>
                <v-text-field
                  ref="inputCodigoBarras"
                  label="Código de barras (no debería repetirse)"
                  type="text"
                  @keyup.enter="enfocarDescripcion"
                  v-model="producto.codigoBarras"
                  :rules="reglas.codigoBarras"
                  hint="Código de barras. Si no existe, invente uno"
                ></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-textarea
                  label="Descripción"
                  type="text"
                  rows="3"
                  ref="inputDescripcion"
                  v-model.trim="producto.descripcion"
                  :rules="reglas.descripcion"
                  hint="Color, tamaño, talla, lo que describe al producto"
                  required
                ></v-textarea>
              </v-flex>
              <v-flex xs12>
                <span class="subheading">Precios</span>
              </v-flex>
              <v-flex xs12 md6>
                <v-text-field
                  label="Compra"
                  type="number"
                  v-model.number="producto.precioCompra"
                  :rules="reglas.precios"
                  hint="Lo que el producto cuesta"
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12 md6>
                <v-text-field
                  label="Venta"
                  type="number"
                  v-model.number="producto.precioVenta"
                  :rules="reglas.precios"
                  hint="Precio en el que el producto se vende"
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12>
                <span v-show="utilidad >= 0" class="subheading">
                  Utilidad: <strong>{{ utilidad | currency }}</strong>
                </span>
              </v-flex>
              <v-flex xs12 md6>
                <v-text-field
                  label="Stock"
                  type="number"
                  v-model.number="producto.stock"
                  hint="Si la existencia del producto es menor al stock, le avisaremos"
                ></v-text-field>
              </v-flex>
              <v-flex xs12 md6>
                <v-text-field
                  label="Existencia actual"
                  type="number"
                  v-model.number="producto.existencia"
                  :rules="reglas.existencia"
                  hint="Existencia actual"
                  required
                ></v-text-field>
              </v-flex>
            </v-layout>
          </v-container>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          :loading="cargando"
          color="green darken-1"
          flat="flat"
          @click.native="guardar()"
          >Guardar</v-btn
        >
        <v-btn color="black" flat="flat" @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  watch: {
    mostrar(estaMostrado) {
      if (estaMostrado) {
        this.enfocarCodigoDeBarras();
        this.obtenerSiguienteNumero();
        if (this.productoParaDuplicar.Numero) {

          // No queremos el número, ese es calculado automáticamente
          //Tampoco el código de barras
          delete this.productoParaDuplicar.Numero;
          this.productoParaDuplicar.CodigoBarras = "";

          /*
          *   Como la API devuelve las claves con
          *   mayúscula al inicio, por ejemplo:
          *     producto.Numero
          *   nosotros debemos convertir la clave a minúscula,
          *   sólo la primera letra, así:
          *     producto.numero
          *   Esto es porque aquí usamos objeto.algunaPropiedad
          *   mientras que la API devuelve objeto.AlgunaPropiedad
          * */
          let claves = Object.keys(this.productoParaDuplicar);
          claves.forEach(clave => {
            let claveConMinusculaAlInicio = clave[0].toLowerCase() + clave.substr(1, clave.length - 1);
            this.producto[claveConMinusculaAlInicio] = this.productoParaDuplicar[clave];
          });
        }
      }
    }
  },
  computed: {
    utilidad() {
      if (this.producto.precioCompra && this.producto.precioVenta) {
        return this.producto.precioVenta - this.producto.precioCompra;
      }
      return -1;
    }
  },
  methods: {
    enfocarCodigoDeBarras() {
      this.$nextTick(this.$refs.inputCodigoBarras.focus);
    },
    enfocarDescripcion() {
      this.$nextTick(this.$refs.inputDescripcion.focus);
    },
    obtenerSiguienteNumero() {
      HTTP_AUTH.get("siguiente/numero/producto").then(siguienteNumero => {
        this.producto.numero = siguienteNumero;
      });
    },
    resetearFormulario() {
      this.$refs.formulario.reset();
    },
    cerrarDialogo() {
      this.resetearFormulario();
      this.$emit("cerrar-dialogo");
    },
    guardar() {
      if (this.$refs.formulario.validate()) {
        let producto = Object.assign({}, this.producto);
        this.cargando = true;
        HTTP_AUTH.post("producto", producto).then(resultados => {
          this.cargando = false;
          if (resultados === true) {
            this.resetearFormulario();
            this.enfocarCodigoDeBarras();
            this.obtenerSiguienteNumero();
            this.$emit("producto-guardado");
          } else {
            /**Handle error here */
          }
        });
      }
    }
  },
  props: ["mostrar", "productoParaDuplicar"],
  data: () => ({
    cargando: false,
    producto: { existencia: 0, numero: 1 },
    reglas: {
      codigoBarras: [
        codigoBarras => {
          if (!codigoBarras) return "Introduzca el código de barras o invente uno que no exista";
          return true;
        }
      ],
      numero: [
        numero => {
          if (!numero) return "Introduzca un código o número";
          if (numero < 0) return "El número debe ser positivo";
          return true;
        }
      ],
      descripcion: [
        descripcion => {
          if (!descripcion) return "Introduzca la descripción del producto";
          return true;
        }
      ],
      precios: [
        precio => {
          if (!precio) return "Introduzca un precio válido";
          if (precio <= 0) return "El precio no puede ser negativo ni cero";
          return true;
        }
      ],
      existencia: [
        existencia => {
          if (
            "undefined" === typeof existencia ||
            null === existencia ||
            ("string" === typeof existencia && existencia.length <= 0)
          )
            return "Introduzca la existencia o 0";
          if (existencia < 0) return "La existencia no puede ser negativa";
          return true;
        }
      ]
    }
  })
};
</script>
