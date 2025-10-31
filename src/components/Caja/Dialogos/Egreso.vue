<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Registrar egreso</v-card-title>
      <v-card-text>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12>
                <v-text-field
                  @keyup.enter="guardar()"
                  ref="monto"
                  prepend-icon="monetization_on"
                  label="Monto"
                  type="number"
                  v-model.number="egreso.monto"
                  :rules="reglas.monto"
                  hint="La cantidad que sale"
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-text-field
                  @keyup.enter="guardar()"
                  prepend-icon="note_add"
                  label="Descripción"
                  type="text"
                  v-model="egreso.descripcion"
                  :rules="reglas.descripcion"
                  hint="La descripción del egreso o salida"
                  required
                ></v-text-field>
              </v-flex>
            </v-layout>
          </v-container>
        </v-form>
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
import { HTTP_AUTH } from "../../../http-common";

export default {
  props: ["mostrar"],
  watch: {
    mostrar(estaMostrado) {
      if (estaMostrado) this.enfocarInputPrincipal();
    }
  },
  methods: {
    enfocarInputPrincipal() {
      this.$nextTick(this.$refs.monto.focus);
    },
    cerrarDialogo() {
      this.$refs.formulario.reset();
      this.$emit("cerrar");
    },
    guardar() {
      if (this.$refs.formulario.validate() && !this.cargando) {
        this.cargando = true;
        HTTP_AUTH.post("egreso", Object.assign({}, this.egreso)).then(
          resultados => {
            this.$emit("guardado");
            this.cargando = false;
            this.$refs.formulario.reset();
          }
        );
      }
    }
  },
  data: () => ({
    cargando: false,
    egreso: {},
    reglas: {
      monto: [
        monto => {
          if (monto <= 0) return "Cantidad inválida";
          if (!monto) return "Introduce la cantidad";
          return true;
        }
      ],
      descripcion: [
        descripcion => {
          if (!descripcion)
            return "¿Cuál es la razón o descripción del egreso?";
          return true;
        }
      ]
    }
  })
};
</script>
