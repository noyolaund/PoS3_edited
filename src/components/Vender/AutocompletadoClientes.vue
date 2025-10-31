<template>
  <v-autocomplete
    :loading="cargando"
    :items="clientes"
    :search-input.sync="busquedaAutocompletado"
    v-model="cliente"
    label="Nombre del cliente"
    return-object
    item-text="Nombre"
    item-value="Nombre"
    clearable
    solo
    prepend-icon="add"
    @click:prepend="agregarNuevoCliente"
    required
  ></v-autocomplete>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  data: () => ({
    clientes: [],
    cliente: {},
    cargando: false,
    busquedaAutocompletado: ""
  }),
  methods: {
    limpiar() {
      this.busqueda = "";
      this.cliente = {};
    },
    agregarNuevoCliente() {
      this.$emit("agregar-cliente");
    },
    buscarClientes(busqueda) {
      this.cargando = true;
      HTTP_AUTH.get(
        `autocompletado/clientes/${encodeURIComponent(busqueda)}`
      ).then(clientes => {
        this.clientes = clientes;
        this.cargando = false;
      });
    }
  },
  watch: {
    busquedaAutocompletado(busqueda) {
      if (busqueda) this.buscarClientes(busqueda);
    },
    cliente(cliente) {
      if (null === cliente) {
        this.$emit("cliente-cancelado");
      } else if (cliente.Numero) {
        this.$emit("cliente-seleccionado", cliente);
      }
    }
  }
};
</script>
