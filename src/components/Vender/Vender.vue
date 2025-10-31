<template>
  <v-layout>
    <v-flex xs12>
      <dialogo-busqueda-producto
        @producto-seleccionado="onProductoSeleccionadoDesdeBusqueda"
        @cerrar-dialogo="dialogos.buscar = false"
        :mostrar="dialogos.buscar"
      ></dialogo-busqueda-producto>

      <!-- Autocompletado inline para buscar productos mientras se vende -->
      <autocompletado-productos
        ref="autocompletado"
        class="inline-autocomplete"
        @producto-seleccionado="onProductoSeleccionadoDesdeBusqueda"
      ></autocompletado-productos>

      <lista-de-productos
        @buscar-producto="mostrarDialogoParaBuscarProducto()"
        @producto-no-existente="onProductoNoExistente"
        ref="listaDeProductos"
      ></lista-de-productos>

      <dialogo-venta-contado
        @error-pago-incompleto="onErrorPagoIncompleto"
        @no-hay-cliente="onErrorNoCliente"
        @venta-realizada="onVentaContadoRealizada"
        @cerrar-dialogo="dialogos.contado = false"
        ref="dialogoVentaContado"
        @agregar-cliente="onAgregarCliente"
        :mostrar="dialogos.contado && !dialogos.nuevoCliente"
        :datosVenta="datosVenta"
      ></dialogo-venta-contado>

      <dialogo-apartado
        @error-pago-excedido="onErrorPagoExcedido"
        @no-hay-cliente="onErrorNoCliente"
        @no-hay-fecha="onErrorNoFecha"
        @apartado-realizado="onApartadoRealizado"
        @cerrar-dialogo="dialogos.apartado = false"
        ref="dialogoApartado"
        @agregar-cliente="onAgregarCliente"
        :mostrar="dialogos.apartado && !dialogos.nuevoCliente"
        :datosVenta="datosVenta"
      ></dialogo-apartado>

      <dialogo-confirmacion-vaciar-lista
        :mostrar="dialogos.confirmarEliminacion"
        @cerrar-dialogo="dialogos.confirmarEliminacion = false"
        @confirmado="onCancelarVenta"
      >
      </dialogo-confirmacion-vaciar-lista>

      <formulario-nuevo-cliente
        @cerrar-dialogo="dialogos.nuevoCliente = false"
        @cliente-guardado="onClienteGuardado"
        :mostrar="dialogos.nuevoCliente"
      ></formulario-nuevo-cliente>
      <speed-dial
        @venta-contado="onVentaContado()"
        @apartado="onApartado()"
        @cancelar-venta="mostrarDialogoConfirmacionSiEsNecesario()"
      >
      </speed-dial>
    </v-flex>

    <v-snackbar
      :timeout="5000"
      :top="true"
      :right="true"
      v-model="snackbars.ventaCorrecta"
    >
      Venta realizada correctamente
      <v-btn flat color="pink" @click.native="snackbars.ventaCorrecta = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="5000"
      :top="true"
      :right="true"
      v-model="snackbars.apartadoCorrecto"
    >
      Apartado realizado correctamente
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.apartadoCorrecto = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :top="true"
      :right="true"
      v-model="snackbars.productoNoExistente"
    >
      El producto no existe
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.productoNoExistente = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :top="true"
      :right="true"
      v-model="snackbars.pagoCliente"
    >
      El pago del cliente debe ser mayor o igual que el total
      <v-btn flat color="pink" @click.native="snackbars.pagoCliente = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :top="true"
      :right="true"
      v-model="snackbars.seleccionarCliente"
    >
      Elige un cliente para continuar
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.seleccionarCliente = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :top="true"
      :right="true"
      v-model="snackbars.seleccionarFecha"
    >
      Elige la fecha en la que el apartado vence
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.seleccionarFecha = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :top="true"
      :right="true"
      v-model="snackbars.agregarProductos"
    >
      Para vender necesitas agregar productos
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.agregarProductos = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :top="true"
      :right="true"
      v-model="snackbars.pagoExcedido"
    >
      El anticipo debe ser menor al pago total
      <v-btn flat color="pink" @click.native="snackbars.pagoExcedido = false"
        >OK</v-btn
      >
    </v-snackbar>
  </v-layout>
</template>


<script>
import ListaDeProductos from "../Vender/ListaDeProductos";
import FormularioNuevoCliente from "../Clientes/FormularioNuevoCliente";
import DialogoVentaContado from "../Vender/DialogoVentaContado";
import DialogoApartado from "../Vender/DialogoApartado";
import SpeedDial from "../Vender/SpeedDial";
import DialogoBusquedaProducto from '../Vender/DialogoBusquedaProducto';
import AutocompletadoProductos from '../Vender/AutocompletadoProductos';
import DialogoConfirmacionVaciarLista from '../Vender/DialogoConfirmacionVaciarLista';
import { EventBus } from "../../main";


export default {
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Vender");
  },
  mounted() {
    // Atajo de teclado: F8 para abrir directamente el diálogo de venta al contado
    // Se añade el listener en mounted y se remueve en beforeDestroy
    window.addEventListener("keydown", this.handleKeydown);
  },
  beforeDestroy() {
    window.removeEventListener("keydown", this.handleKeydown);
  },
  components: {
    ListaDeProductos,
    SpeedDial,
    DialogoVentaContado,
    FormularioNuevoCliente,
    DialogoApartado,
    DialogoBusquedaProducto,
    AutocompletadoProductos,
    DialogoConfirmacionVaciarLista
  },
  data: () => ({
    fab: false,
    dialogos: {
      contado: false,
      apartado: false,
      nuevoCliente: false,
      buscar: false,
      confirmarEliminacion: false,
    },
    snackbars: {
      agregarProductos: false,
      ventaCorrecta: false,
      productoNoExistente: false,
      pagoCliente: false,
      seleccionarCliente: false,
      apartadoCorrecto: false,
      seleccionarFecha: false,
      pagoExcedido: false
    },
    datosVenta: {}
  }),
  methods: {
    handleKeydown(event) {
      // F8 -> abrir diálogo de venta al contado
      if (event.key === "F8" || event.keyCode === 119) {
        event.preventDefault();
        this.onVentaContado();
        return;
      }

      // F9 -> enfocar el autocompletado inline
      if (event.key === "F9" || event.keyCode === 120) {
        event.preventDefault();
        this.focusAutocompletado();
        return;
      }

      // / -> enfocar el campo inicial de código/entrada en ListaDeProductos
      if (event.key === "/") {
        event.preventDefault();
        this.focusCodigoInicial();
        return;
      }
    },

    focusAutocompletado() {
      // Accede al input interno del componente AutocompletadoProductos
      this.$nextTick(() => {
        try {
          const ac = this.$refs.autocompletado;
          if (ac && ac.$refs && ac.$refs.input && typeof ac.$refs.input.focus === 'function') {
            ac.$refs.input.focus();
          }
        } catch (e) {
          // no-op
        }
      });
    },

    focusCodigoInicial() {
      this.$nextTick(() => {
        try {
          const lista = this.$refs.listaDeProductos;
          if (lista && lista.$refs && lista.$refs.codigoInput && typeof lista.$refs.codigoInput.focus === 'function') {
            lista.$refs.codigoInput.focus();
          }
        } catch (e) {
          // no-op
        }
      });
    },
    mostrarDialogoConfirmacionSiEsNecesario() {
      if (this.$refs.listaDeProductos.totalVenta > 0) this.dialogos.confirmarEliminacion = true;
    },
    onProductoSeleccionadoDesdeBusqueda(producto) {
      this.dialogos.buscar = false;
      this.$refs.listaDeProductos.agregarOModificarExistenciaDeProductoEnLista(producto);
    },
    mostrarDialogoParaBuscarProducto() {
      this.dialogos.buscar = true;
    },
    onErrorPagoIncompleto() {
      this.snackbars.pagoCliente = true;
    },
    onErrorPagoExcedido() {
      this.snackbars.pagoExcedido = true;
    },
    onErrorNoCliente() {
      this.snackbars.seleccionarCliente = true;
    },
    onErrorNoFecha() {
      this.snackbars.seleccionarFecha = true;
    },
    onProductoNoExistente() {
      this.snackbars.productoNoExistente = true;
    },
    onVentaContadoRealizada() {
      this.dialogos.contado = false;
      this.onCancelarVenta();
      this.snackbars.ventaCorrecta = true;
    },
    onApartadoRealizado() {
      this.dialogos.apartado = false;
      this.onCancelarVenta();
      this.snackbars.apartadoCorrecto = true;
    },
    onClienteGuardado(clienteGuardado) {
      this.dialogos.nuevoCliente = false;
      this.$refs.dialogoVentaContado.setCliente(
        Object.assign({}, clienteGuardado)
      );
      this.$refs.dialogoApartado.setCliente(Object.assign({}, clienteGuardado));
    },
    onAgregarCliente() {
      this.dialogos.nuevoCliente = true;
    },
    onVentaContado() {
      if (this.$refs.listaDeProductos.totalVenta) {
        this.datosVenta.total = this.$refs.listaDeProductos.totalVenta;
        this.datosVenta.lista = Array.from(
          this.$refs.listaDeProductos.listaDeProductos
        );
        this.dialogos.contado = true;
      } else {
        this.snackbars.agregarProductos = true;
      }
    },
    onApartado() {
      if (this.$refs.listaDeProductos.totalVenta) {
        this.datosVenta.total = this.$refs.listaDeProductos.totalVenta;
        this.datosVenta.lista = Array.from(
          this.$refs.listaDeProductos.listaDeProductos
        );
        this.dialogos.apartado = true;
      } else {
        this.snackbars.agregarProductos = true;
      }
    },
    onCancelarVenta() {
      this.dialogos.confirmarEliminacion = false;
      this.$refs.listaDeProductos.cancelarVenta();
    }
  }
};
</script>

<style scoped>
/* Reducir el espacio superior alrededor del autocompletado inline */
.inline-autocomplete {
  margin-top: -20px;
  margin-bottom: -20px;
}

/* Ajuste adicional para el layout si hace falta */
v-layout {
  padding-top: 4px;
}
</style>
