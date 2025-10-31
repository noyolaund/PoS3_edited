<template>
  <v-dialog v-model="mostrar" persistent max-width="700">
    <v-card>
      <v-card-title class="headline">Buscar producto</v-card-title>
      <v-card-text>
        <autocompletado-productos
          @producto-seleccionado="productoSeleccionado"
        ></autocompletado-productos>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import AutocompletadoProductos from './AutocompletadoProductos';

export default {
  props: ["mostrar"],
  components: { AutocompletadoProductos },
  methods: {
    productoSeleccionado(producto) {
      this.$emit("producto-seleccionado", producto);
    },
    cerrarDialogo() {
      this.$emit("cerrar-dialogo");
    }
  }
}
</script>
