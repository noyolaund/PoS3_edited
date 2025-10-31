<template>
  <v-layout>
    <dialogo-aumentar-existencia
      :mostrar="dialogos.aumentarExistencia"
      :producto="productoEditado"
      @cerrar="dialogos.aumentarExistencia = false"
      @aumentado="onCantidadAumentada"
    ></dialogo-aumentar-existencia>

    <dialogo-restar-existencia
      :mostrar="dialogos.restarExistencia"
      :producto="productoEditado"
      @cerrar="dialogos.restarExistencia = false"
      @restado="onCantidadRestada"
    ></dialogo-restar-existencia>

    <formulario-nuevo-producto
      @producto-guardado="onProductoGuardado()"
      @cerrar-dialogo="dialogos.nuevo = false"
      :productoParaDuplicar="productoEditado"
      :mostrar="dialogos.nuevo"
    ></formulario-nuevo-producto>

    <formulario-editar-producto
      @producto-guardado="onCambiosDeProductoGuardados()"
      @cerrar-dialogo="dialogos.editar = false"
      :mostrar="dialogos.editar"
      :producto="productoEditado"
    ></formulario-editar-producto>

    <lista-de-productos
      @editar-producto="editarProducto"
      @eliminar-producto="eliminarProducto"
      @aumentar-existencia="aumentarExistencia"
      @restar-existencia="restarExistencia"
      @duplicar-producto="duplicarProducto"
      ref="productos"
    ></lista-de-productos>

    <dialogo-confirmar-eliminacion
      :mostrar="dialogos.confirmarEliminacion"
      :producto="productoParaEliminar"
      @cerrar-dialogo="dialogos.confirmarEliminacion = false"
      @producto-eliminado="onProductoEliminado()"
    ></dialogo-confirmar-eliminacion>

    <v-snackbar
      :timeout="2000"
      :bottom="true"
      v-model="snackbars.existenciaAumentada"
    >
      Existencia aumentada
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.existenciaAumentada = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :bottom="true"
      v-model="snackbars.existenciaRestada"
    >
      Existencia restada
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.existenciaRestada = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :bottom="true"
      v-model="snackbars.productoGuardado"
    >
      Producto guardado correctamente
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.productoGuardado = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :bottom="true"
      v-model="snackbars.cambiosProductoGuardados"
    >
      Producto guardado correctamente
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.productoGuardado = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-tooltip top>
      <v-btn
        @click="mostrarDialogoParaAgregar()"
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
      <span>Agregar producto</span>
    </v-tooltip>
  </v-layout>
</template>
<script>
import FormularioNuevoProducto from "../Inventario/FormularioNuevoProducto";
import FormularioEditarProducto from "../Inventario/FormularioEditarProducto";
import ListaDeProductos from "../Inventario/ListaDeProductos";
import DialogoConfirmarEliminacion from "../Inventario/DialogoConfirmarEliminacion";
import { EventBus } from "../../main";
import DialogoAumentarExistencia from "./DialogoAumentarExistencia"
import DialogoRestarExistencia from "./DialogoRestarExistencia"

export default {
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Inventario");
  },
  components: {
    FormularioNuevoProducto,
    ListaDeProductos,
    FormularioEditarProducto,
    DialogoConfirmarEliminacion,
    DialogoAumentarExistencia,
    DialogoRestarExistencia
  },
  data: () => ({
    dialogos: {
      nuevo: false,
      editar: false,
      aumentarExistencia: false,
      restarExistencia: false,
      confirmarEliminacion: false
    },
    snackbars: {
      productoGuardado: false,
      cambiosProductoGuardados: false,
      existenciaAumentada: false,
      existenciaRestada: false,
    },
    productoEditado: {},
    productoParaEliminar: {}
  }),
  methods: {
    onCantidadAumentada() {
      this.dialogos.aumentarExistencia = false;
      this.snackbars.existenciaAumentada = true;
      this.$refs.productos.obtener();
    },
    onCantidadRestada() {
      this.dialogos.restarExistencia = false;
      this.snackbars.existenciaRestada = true;
      this.$refs.productos.obtener();
    },
    onProductoGuardado() {
      this.snackbars.productoGuardado = true;
      this.$refs.productos.obtener();
    },
    onCambiosDeProductoGuardados() {
      this.snackbars.cambiosProductoGuardados = true;
      this.dialogos.editar = false;
      this.$refs.productos.obtener();
    },
    onProductoEliminado() {
      this.dialogos.confirmarEliminacion = false;
      this.$refs.productos.obtener();
    },
    mostrarDialogoParaAgregar() {
      this.dialogos.nuevo = true;
    },
    editarProducto(producto) {
      this.dialogos.editar = true;
      this.productoEditado = Object.assign({}, producto);
    },
    eliminarProducto(producto) {
      this.productoParaEliminar = Object.assign({}, producto);
      this.dialogos.confirmarEliminacion = true;
    },
    aumentarExistencia(producto) {
      this.productoEditado = Object.assign({}, producto);
      this.dialogos.aumentarExistencia = true;
    },
    restarExistencia(producto) {
      this.productoEditado = Object.assign({}, producto);
      this.dialogos.restarExistencia = true;
    },
    duplicarProducto(producto) {
      this.productoEditado = Object.assign({}, producto);
      this.mostrarDialogoParaAgregar();
    }
  }
};
</script>
