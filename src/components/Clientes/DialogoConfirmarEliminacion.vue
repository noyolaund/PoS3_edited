<template>
  <v-dialog v-model="mostrar" max-width="400">
    <v-card>
      <v-card-title class="headline">Confirmar eliminación</v-card-title>
      <v-card-text>
        ¿Realmente desea eliminar el cliente con el nombre
        <strong>{{ cliente.Nombre }}</strong> y el número
        <strong>{{ cliente.NumeroTelefono }}</strong
        >? <br />
        Esta acción no se puede revertir.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          :loading="cargando"
          color="green darken-1"
          flat="flat"
          @click="confirmarEliminacion()"
          >Sí, eliminar
        </v-btn>
        <v-btn color="black" flat="flat" @click.native="ocultarDialogo()"
          >Cancelar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  props: ["mostrar", "cliente"],
  data: () => ({
    cargando: false,
  }),
  methods: {
    confirmarEliminacion() {
      this.cargando = true;
      HTTP_AUTH.delete(`cliente/${this.cliente.Numero}`).then(respuesta => {
        this.cargando = false;
        this.$emit("cliente-eliminado");
        this.ocultarDialogo();
      });
    },
    ocultarDialogo() {
      this.$emit("cerrar-dialogo");
    }
  }
};
</script>
