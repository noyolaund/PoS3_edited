<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Error al registrar</v-card-title>
      <v-card-text>
        <p>
          Ha ocurrido un error al registrar tu negocio. Intenta registrarte de
          nuevo más tarde. Si continúas teniendo problemas, asegúrate de haber
          instalado el programa correctamente. Enlace de ayuda:
          <a
            target="_blank"
            href="https://parzibyte.me/blog/2021/06/19/ayuda-soporte-sublime-pos-3/"
            >https://parzibyte.me/blog/2021/06/19/ayuda-soporte-sublime-pos-3/</a
          >
        </p>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="success" @click.native="cerrarDialogo()">Cerrar</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { CorreoSoporteYContacto } from "../../constantes";

export default {
  name: "DialogoNegocioNoRegistrado",
  props: ["mostrar"],
  methods: {
    cerrarDialogo() {
      this.$emit("cerrar-dialogo");
    },
    data: () => ({
      correo: CorreoSoporteYContacto,
    })
  }
};
</script>

<style scoped>
</style>
