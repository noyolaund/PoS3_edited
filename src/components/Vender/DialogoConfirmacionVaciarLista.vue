<template>
  <v-dialog v-model="mostrar" persistent max-width="300">
    <v-card>
      <v-card-title class="headline">¿Cancelar venta actual?</v-card-title>
      <v-card-text>
        ¿Realmente desea vaciar la lista de productos?
        <br />
        <strong>Esta acción no se puede revertir</strong>.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
        <v-btn
          color="green darken-1"
          flat="flat"
          @click="confirmarEliminacion()"
          >Sí, vaciar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  props: ["mostrar"],
  methods: {
    cerrarDialogo() {
      this.$emit("cerrar-dialogo");
    },
    confirmarEliminacion() {
      this.$emit("confirmado");
    }
  }
}
</script>
