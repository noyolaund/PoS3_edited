<template>
  <v-speed-dial
    v-model="fab"
    fixed
    transition="slide-x"
    bottom
    right
    direction="top"
  >
    <v-btn slot="activator" color="blue darken-2" dark fab hover v-model="fab">
      <v-icon>check</v-icon>
      <v-icon>close</v-icon>
    </v-btn>
    <v-tooltip disabled :value="true" left>
      <v-btn
        @click="ventaContado()"
        fab
        dark
        small
        color="green"
        slot="activator"
      >
        <v-icon>attach_money</v-icon>
      </v-btn>
      <span>Al contado</span>
    </v-tooltip>
    <v-tooltip disabled :value="true" left>
      <v-btn @click="apartado()" fab dark small color="indigo" slot="activator">
        <v-icon>credit_card</v-icon>
      </v-btn>
      <span>Apartado</span>
    </v-tooltip>
    <v-tooltip disabled :value="true" left>
      <v-btn
        @click="cancelarVenta()"
        fab
        dark
        small
        color="red"
        slot="activator"
      >
        <v-icon>delete</v-icon>
      </v-btn>
      <span>Cancelar venta</span>
    </v-tooltip>
  </v-speed-dial>
</template>
<script>
export default {
  data: () => ({
    fab: false
  }),
  methods: {
    ventaContado() {
      this.$emit("venta-contado");
    },
    apartado() {
      this.$emit("apartado");
    },
    cancelarVenta() {
      this.$emit("cancelar-venta");
    }
  }
};
</script>
