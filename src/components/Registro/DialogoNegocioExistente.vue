<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Correo electrónico existente</v-card-title>
      <v-card-text>
        <p>
          Parece que ya hay un negocio registrado con ese correo electrónico.
          Intenta de nuevo con otra dirección. Si continúas teniendo problemas,
          asegúrate de haber instalado el programa correctamente. Enlace de
          ayuda:
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
  name: "DialogoNegocioExistente",
  props: ["mostrar"],
  methods: {
    cerrarDialogo() {
      this.$emit("cerrar-dialogo");
    },
    data: () => ({
      correo: CorreoSoporteYContacto,
    })
  }
}
</script>

<style scoped>
</style>
