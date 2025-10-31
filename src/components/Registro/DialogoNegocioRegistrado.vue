<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Registro correcto</v-card-title>
      <v-card-text>
        <p>
          Registro correcto, será redireccionado
          <br />
          Para iniciar sesión coloque los siguientes datos:
        </p>
        <strong>Correo: </strong> El que usted eligió
        <br />
        <strong>Usuario: </strong> <code>{{ usuarioPorDefecto }}</code>
        <br />
        <strong>Contraseña: </strong> la que usted eligió
        <br />
        <strong>Más adelante podrá crear más usuarios</strong>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="success" @click.native="cerrarDialogo()">OK</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { CorreoSoporteYContacto, UsuarioPorDefecto } from "../../constantes";

export default {
  name: "DialogoNegocioRegistrado",
  props: ["mostrar"],
  methods: {
    cerrarDialogo() {
      this.$emit("cerrar-dialogo");
    },
  },
  data: () => ({
    correo: CorreoSoporteYContacto,
    usuarioPorDefecto: UsuarioPorDefecto,
  }),
};
</script>

<style scoped>
</style>
