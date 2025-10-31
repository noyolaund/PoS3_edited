<template>
  <v-layout row wrap>
    <v-flex xs12>
      <Lista ref="lista" @cambiarPermisos="onCambiarPermisos"></Lista>
    </v-flex>
    <DialogoNuevo
      @correcto="onRegistroCorrecto"
      @error="snackbars.registroIncorrecto = true"
      :mostrar="dialogos.nuevo"
      @cerrar="dialogos.nuevo = false"
    ></DialogoNuevo>
    <Permisos
      :mostrar="dialogos.permisos"
      :usuario="usuarioSeleccionado"
      @cerrar="dialogos.permisos = false"
      @asignados-correctamente="onPermisosAsignadosCorrectamente"
      @error-asignando="onErrorAlAsignarPermisos"
    ></Permisos>
    <v-snackbar
      :timeout="2000"
      :top="true"
      :right="true"
      v-model="snackbars.registroCorrecto"
    >
      Usuario registrado. No olvide asignarle permisos
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.registroCorrecto = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="5000"
      :top="true"
      :right="true"
      v-model="snackbars.registroIncorrecto"
    >
      Error al registrar. ¿Tal vez el nombre de usuario está ocupado?
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.registroIncorrecto = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="5000"
      :top="true"
      :right="true"
      v-model="snackbars.permisosAsignados"
    >
      Permisos asignados correctamente
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.permisosAsignados = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="0"
      :top="true"
      :right="true"
      v-model="snackbars.errorAlAsignarPermisos"
    >
      Error al asignar. ¿tal vez intentó cambiar los permisos del administrador?
      (recuerde que estos no se pueden cambiar)
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.errorAlAsignarPermisos = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-btn
      @click="dialogos.nuevo = true"
      fixed
      dark
      fab
      bottom
      fill-height
      slot="activator"
      open-delay="0"
      right
      color="indigo"
    >
      <v-icon>add</v-icon>
    </v-btn>
  </v-layout>
</template>

<script>
import DialogoNuevo from './DialogoNuevo'
import Lista from './Lista'
import Permisos from './Permisos'
import { EventBus } from "../../main";

export default {
  methods: {
    onErrorAlAsignarPermisos() {
      this.snackbars.errorAlAsignarPermisos = true;
    },
    onPermisosAsignadosCorrectamente() {
      this.snackbars.permisosAsignados = true;
    },
    onRegistroCorrecto() {
      this.$refs.lista.obtener();
      this.snackbars.registroCorrecto = true;
    },
    onCambiarPermisos(usuario) {
      this.usuarioSeleccionado = Object.assign({}, usuario);
      this.dialogos.permisos = true;
    }
  },
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Usuarios")
  },
  data: () => ({
    snackbars: {
      errorAlAsignarPermisos: false,
      registroCorrecto: false,
      registroIncorrecto: false,
      permisosAsignados: false,
    },
    dialogos: {
      nuevo: false,
      permisos: false,
    },
    usuarioSeleccionado: {}
  }),
  name: "Usuarios",
  components: { Lista, DialogoNuevo, Permisos }
}
</script>

<style>
</style>
