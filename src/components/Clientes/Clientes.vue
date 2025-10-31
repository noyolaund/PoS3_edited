<template>
  <v-layout>
    <lista-de-clientes
      @eliminar="onEliminarCliente"
      @editar="onEditarCliente"
      @mostrar-historial="onMostrarHistorial"
      ref="listaDeClientes"
    ></lista-de-clientes>

    <formulario-nuevo-cliente
      :mostrar="dialogos.nuevo"
      @cliente-guardado="onClienteGuardado()"
      @cerrar-dialogo="dialogos.nuevo = false"
    ></formulario-nuevo-cliente>

    <formulario-editar-cliente
      @cliente-guardado="onCambiosClienteGuardados()"
      :mostrar="dialogos.editar"
      :cliente="clienteEditado"
      @cerrar-dialogo="dialogos.editar = false"
    ></formulario-editar-cliente>

    <dialogo-confirmar-eliminacion
      :cliente="clienteParaEliminar"
      @cliente-eliminado="onClienteEliminado"
      @cerrar-dialogo="dialogos.eliminar = false"
      :mostrar="dialogos.eliminar"
    ></dialogo-confirmar-eliminacion>

    <historial
      :mostrar="dialogos.historial"
      :idCliente="idClienteHistorial"
      @cerrar="dialogos.historial = false"
    ></historial>
    <v-tooltip top open-delay="0">
      <v-btn
        @click="dialogos.nuevo = true"
        fixed
        dark
        fab
        bottom
        fill-height
        slot="activator"
        right
        color="pink"
      >
        <v-icon>add</v-icon>
      </v-btn>
      <span>Registrar cliente</span>
    </v-tooltip>
    <v-snackbar
      :timeout="2000"
      :bottom="true"
      v-model="snackbars.clienteGuardado"
    >
      Cliente guardado correctamente
      <v-btn flat color="pink" @click.native="snackbars.clienteGuardado = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :bottom="true"
      v-model="snackbars.clienteEditado"
    >
      Cambios guardados
      <v-btn flat color="pink" @click.native="snackbars.clienteEditado = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :bottom="true"
      v-model="snackbars.clienteEliminado"
    >
      Cliente eliminado
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.clienteEliminado = false"
        >OK</v-btn
      >
    </v-snackbar>
  </v-layout>
</template>
<script>
import ListaDeClientes from "../Clientes/ListaDeClientes";
import FormularioNuevoCliente from "../Clientes/FormularioNuevoCliente";
import FormularioEditarCliente from "../Clientes/FormularioEditarCliente";
import DialogoConfirmarEliminacion from "../Clientes/DialogoConfirmarEliminacion";
import Historial from "./Historial"
import { EventBus } from "../../main";

export default {
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Clientes");
  },
  methods: {
    onMostrarHistorial(idCliente) {
      this.idClienteHistorial = idCliente;
      this.dialogos.historial = true;
    },
    onEditarCliente(cliente) {
      this.clienteEditado = Object.assign({}, cliente);
      this.dialogos.editar = true;
    },
    onEliminarCliente(cliente) {
      this.clienteParaEliminar = Object.assign({}, cliente);
      this.dialogos.eliminar = true;
    },
    onClienteEliminado() {
      this.snackbars.clienteEliminado = true;
      this.$refs.listaDeClientes.obtener();
    },
    onClienteGuardado() {
      this.snackbars.clienteGuardado = true;
      this.$refs.listaDeClientes.obtener();
    },
    onCambiosClienteGuardados() {
      this.snackbars.clienteEditado = true;
      this.$refs.listaDeClientes.obtener();
    }
  },
  data: () => ({
    idClienteHistorial: null,
    clienteEditado: {},
    clienteParaEliminar: {},
    dialogos: {
      nuevo: false,
      editar: false,
      eliminar: false,
      historial: false,
    },
    snackbars: {
      clienteGuardado: false,
      clienteEditado: false,
      clienteEliminado: false
    }
  }),
  components: {
    ListaDeClientes,
    FormularioNuevoCliente,
    FormularioEditarCliente,
    DialogoConfirmarEliminacion,
    Historial,
  }
};
</script>
