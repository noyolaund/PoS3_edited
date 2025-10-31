<template>
  <v-layout row wrap>
    <v-flex xs12>
      <h1 class="display-1">Lamentamos que te vayas</h1>

      <p>
        Un último paso: confirma que realmente deseas eliminar tu cuenta
        (<strong>esta acción no se puede deshacer</strong>)
      </p>
      <v-btn
        :disabled="deshabilitarBoton"
        :loading="cargando"
        @click="confirmarEliminacion()"
        large
        color="error"
      >
        Eliminar mi cuenta
      </v-btn>
      <v-alert :value="alerta.mostrar" :type="alerta.tipo">
        {{ alerta.mensaje }}
      </v-alert>
    </v-flex>
  </v-layout>
</template>

<script>
import { HTTP } from "../http-common";

export default {
  name: "EliminarNegocio",
  data: () => ({
    cargando: false,
    procesoTerminado: false,
    deshabilitarBoton: false,
    alerta: {
      mostrar: false,
      tipo: "success",
      mensaje: ""
    }
  }),
  methods: {
    confirmarEliminacion() {
      let token = this.$route.params.token;
      this.cargando = true;
      HTTP.get("logout")
        .then(() =>
          HTTP
            .get(`negocio/eliminar/${token}`)
            .then(respuesta => {
              this.alerta.mostrar = true;
              if (respuesta === true) {
                this.alerta.mensaje = "Cuenta eliminada correctamente. En cualquier momento del futuro puedes volver a registrarte";
                this.alerta.tipo = "success";
              } else {
                this.alerta.mensaje = "Error eliminando cuenta. ¿Tal vez el token es incorrecto o el negocio no existe?";
                this.alerta.tipo = "error";
              }
            })
            .finally(() => {
              this.cargando = false;
              this.deshabilitarBoton = true;
            }));
    }
  }
}
</script>

<style scoped>
</style>
