<template>
  <v-flex>
    <v-alert
      :value="!clienteSeleccionado || !clienteSeleccionado.Numero"
      type="error"
    >
      No ha seleccionado ningún cliente
    </v-alert>
    <v-alert :value="hayClienteSeleccionado" type="info">
      <strong>Cliente seleccionado</strong>
      <br />
      <span v-if="hayClienteSeleccionado">
        <strong>ID: </strong> {{ clienteSeleccionado.Numero }}
        <br />
      </span>
      <span v-if="hayClienteSeleccionado">
        <strong>Nombre: </strong> {{ clienteSeleccionado.Nombre }}
        <br />
      </span>
      <span v-if="hayClienteSeleccionado">
        <strong>Teléfono: </strong> {{ clienteSeleccionado.NumeroTelefono }}
        <br />
      </span>
    </v-alert>
    <br />
  </v-flex>
</template>
<script>
export default {
  props: ["clienteSeleccionado"],
  computed: {
    hayClienteSeleccionado() {
      if (this.clienteSeleccionado) {
        if (
          this.clienteSeleccionado.Nombre &&
          this.clienteSeleccionado.NumeroTelefono
        )
          return true;
      }
      return false;
    }
  }
};
</script>
