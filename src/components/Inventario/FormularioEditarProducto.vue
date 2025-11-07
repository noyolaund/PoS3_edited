<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Modificar producto</v-card-title>
      <v-card-text>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12>
                <v-text-field
                  readonly
                  label="Número"
                  type="number"
                  v-model.number="producto.Numero"
                  :rules="reglas.numero"
                  hint="Recuerde que no puede cambiar el número de producto"
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-text-field
                  label="Código de barras"
                  type="text"
                  v-model="producto.CodigoBarras"
                  @keyup.enter="enfocarDescripcion"
                  hint="Código de barras, si es que existe"
                ></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-textarea
                  label="Descripción"
                  type="text"
                  rows="3"
                  v-model.trim="producto.Descripcion"
                  :rules="reglas.descripcion"
                  ref="inputDescripcion"
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
                  v-model.number="producto.PrecioCompra"
                  :rules="reglas.precios"
                  hint="Lo que el producto cuesta"
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12 md6>
                <v-text-field
                  label="Venta"
                  type="number"
                  v-model.number="producto.PrecioVenta"
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
                  v-model.number="producto.Stock"
                  hint="Si la existencia del producto es menor al stock, le avisaremos"
                ></v-text-field>
              </v-flex>
              <v-flex xs12 md6>
                <v-checkbox
                  v-model="esHijo"
                  label="Este producto es hijo (pieza)"
                ></v-checkbox>
                <v-select
                  v-if="esHijo"
                  :items="padres"
                  item-text="Descripcion"
                  item-value="Numero"
                  label="Producto padre"
                  v-model="producto.Padre"
                  hint="Seleccione el producto contenedor (cartón)"
                ></v-select>
                <v-text-field
                  v-if="esHijo"
                  label="Equivalencia (piezas por contenedor)"
                  type="number"
                  v-model.number="producto.Equivalencia"
                  hint="Cuántas piezas contiene un cartón (p.ej. 12)"
                ></v-text-field>
                <v-text-field
                  v-if="!esHijo"
                  label="Existencia actual"
                  type="number"
                  v-model.number="producto.Existencia"
                  :rules="reglas.existencia"
                  hint="Existencia actual"
                  required
                ></v-text-field>
                <v-alert v-if="esHijo" type="info" dense>
                  La existencia se gestiona desde el producto padre. Use "Agregar inventario" para reabastecer piezas.
                </v-alert>
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
        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >Cancelar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  computed: {
    utilidad() {
      if (this.producto.PrecioCompra && this.producto.PrecioVenta) {
        return this.producto.PrecioVenta - this.producto.PrecioCompra;
      }
      return -1;
    }
  },
  methods: {
    enfocarDescripcion() {
      this.$nextTick(this.$refs.inputDescripcion.focus);
    },
    resetearFormulario() {
      this.$refs.formulario.reset();
    },
    cerrarDialogo() {
      this.resetearFormulario();
      this.$emit("cerrar-dialogo");
    },
    obtenerPadres() {
      HTTP_AUTH.get("productos/0/1000").then(resultados => {
        this.padres = resultados.Productos.filter(p => !p.Padre || p.Padre == 0);
      });
    },
    guardar() {
      if (this.$refs.formulario.validate()) {
        let producto = Object.assign({}, this.producto);
        this.cargando = true;
        HTTP_AUTH.put("producto", producto).then(resultados => {
          this.cargando = false;
          if (resultados === true) {
            this.resetearFormulario();
            this.$emit("producto-guardado");
          } else {
            /**Handle error here */
          }
        });
      }
    }
  },
  watch: {
    mostrar(esta) {
      if (esta) {
        this.obtenerPadres();
        this.esHijo = this.producto.Padre && this.producto.Padre > 0;
      }
    }
  },
  props: ["mostrar", "producto"],
  data: () => ({
    cargando: false,
    padres: [],
    esHijo: false,
    reglas: {
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
