<template>
  <v-container>
    <v-autocomplete
      ref="input"
      :loading="cargando"
      :items="productos"
      :search-input.sync="busquedaAutocompletado"
      v-model="productoSeleccionado"
      label="Escriba parte de la descripción del producto"
      return-object
      item-text="Descripcion"
      item-value="Numero"
      no-data-text="No existe ningún producto con esa descripción"
      solo
      @change="onProductoSeleccionado"
    >
    </v-autocomplete>
  </v-container>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  props: ["mostrar"],
  mounted() {
    this.productos = [];
  },
  watch: {
    mostrar(mostrado) {
      if (mostrado) {
        this.$nextTick(this.$refs.input.focus);
      }
    },
    busquedaAutocompletado(busqueda) {
      if (busqueda) this.buscar(busqueda);
    }
  },
  data: () => ({
    cargando: false,
    productos: [],
    busquedaAutocompletado: "",
    productoSeleccionado: {}
  }),
  methods: {
    onProductoSeleccionado(producto) {
      if (producto.Numero) {
        setTimeout(() => {
          let p = Object.assign({}, producto);
          this.productoSeleccionado = {};
          this.busquedaAutocompletado = "";
          this.$emit("producto-seleccionado", p);
          this.productos = [];
        }, 100);//TODO: arreglar este maldito hack
      }
    },
    buscar(busqueda) {
      this.cargando = true;
      HTTP_AUTH.get(`buscar/productos/autocompletado/${busqueda}`).then(productosQueCoinciden => {
        this.productos = productosQueCoinciden;
        this.cargando = false;
      });
    }
  }
}
</script>

