<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Cambiar precio</v-card-title>
      <v-card-text>
        <p>
          <strong>Precio original: </strong
          >{{ producto.PrecioVentaOriginal | currency }}
        </p>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12>
                <v-text-field
                  @keydown.prevent.native.enter="guardarNuevoPrecio"
                  :rules="reglas.nuevoPrecio"
                  ref="nuevoPrecio"
                  prepend-icon="attach_money"
                  label="Nuevo precio de venta"
                  type="number"
                  hint="Escriba el nuevo precio de venta"
                  v-model.number="nuevoPrecio"
                  required
                ></v-text-field>
              </v-flex>
            </v-layout>
          </v-container>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >Cancelar</v-btn
        >
        <v-btn color="green darken-1" flat="flat" @click="guardarNuevoPrecio()"
          >Ok</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "DialogoCambiarPrecioVenta",
  watch: {
    mostrar(estaMostrado) {
      if (estaMostrado) this.$nextTick(() => {
        this.$refs.nuevoPrecio.focus();
        this.$refs.formulario.reset();
      });
    }
  },
  data: () => ({
    nuevoPrecio: null,
    reglas: {
      nuevoPrecio: [
        nuevoPrecio => {
          if (!isNaN(nuevoPrecio) && nuevoPrecio <= 0) return "Escriba un número mayor a 0";
          if (!nuevoPrecio) return "Escriba un precio";
          return true;
        }
      ]
    },
  }),
  methods: {
    cerrarDialogo() {
      // this.$refs.formulario.reset();
      this.$emit("cerrar");
    },
    guardarNuevoPrecio() {
      if (!this.$refs.formulario.validate()) return;
      this.producto.PrecioVenta = this.nuevoPrecio;
      this.producto.Total = this.producto.PrecioVenta * this.producto.Cantidad;
      // this.nuevoPrecio = null;
      this.$refs.formulario.reset();
      this.cerrarDialogo();
    }
  },
  props: ["mostrar", "producto"]
}
</script>


