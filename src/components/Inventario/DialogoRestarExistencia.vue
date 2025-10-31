<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Restar inventario</v-card-title>
      <v-card-text>
        <p>
          <strong>Producto: </strong> {{ producto.Descripcion }} <br />
          <strong>ID: </strong> {{ producto.Numero }}<br />
          <strong>Código de barras: </strong> {{ producto.CodigoBarras }}<br />
          <strong>Existencia actual: </strong> {{ producto.Existencia }} <br />
        </p>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12>
                <v-text-field
                  @keydown.prevent.native.enter="guardar()"
                  ref="cantidad"
                  label="Cantidad a restar"
                  type="number"
                  :rules="reglas.cantidad"
                  v-model.number="cantidadQueSeResta"
                  hint="¿Cuántos productos de este tipo resta?"
                  required
                ></v-text-field>
              </v-flex>
            </v-layout>
          </v-container>
        </v-form>
        <p v-show="nuevaExistencia > -1">
          Una vez guardada, la existencia del producto será de:<br />
          <span class="title">{{ nuevaExistencia }}</span> <br />
        </p>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          :loading="cargando"
          color="green darken-1"
          flat="flat"
          @click.native="guardar()"
          >Guardar</v-btn
        >
        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  name: "DialogoRestarExistencia",
  props: ["producto", "mostrar"],
  watch: {
    mostrar(mostrado) {
      if (mostrado) {
        this.$nextTick(this.$refs.cantidad.focus);
        this.$refs.formulario.reset();
        this.cantidadQueSeResta = null;
      }
    }
  },
  methods: {
    guardar() {
      if (this.cargando) return;
      if (!this.$refs.formulario.validate()) return;
      let producto = Object.assign({}, this.producto);
      producto.Existencia -= this.cantidadQueSeResta;
      this.cargando = true;
      HTTP_AUTH.put("producto", producto).then(respuesta => {
        if (respuesta) {
          this.$emit("restado");
          this.cantidadQueSeResta = null;
          this.cargando = false;
        }
      });
    },
    cerrarDialogo() {
      this.$refs.formulario.reset();
      this.$emit("cerrar");
    }
  },
  data: function () {
    let esto = this;
    return {
      cargando: false,
      cantidadQueSeResta: null,
      reglas: {
        cantidad: [
          cantidad => {
            if (!cantidad) return "Ingrese la cantidad";
            cantidad = parseFloat(cantidad);
            if (cantidad <= 0) return "Cantidad inválida";
            if (cantidad > esto.producto.Existencia) return "No puede quitar más de lo existente";
            return true;
          }
        ]
      }
    }
  },
  computed: {
    nuevaExistencia() {
      if (!this.cantidadQueSeResta || !this.producto.Existencia) return -1;
      return this.producto.Existencia - this.cantidadQueSeResta;

    }
  }
}
</script>

