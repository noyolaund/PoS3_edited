<template>
  <v-dialog v-model="mostrar" persistent max-width="700">
    <v-card>
      <v-card-title class="headline"
        >Cambiando fecha de vencimiento</v-card-title
      >
      <v-card-text>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12>
                <v-menu
                  ref="menu"
                  :close-on-content-click="false"
                  v-model="mostrarDialogoFecha"
                  :nudge-right="40"
                  :return-value.sync="nuevaFechaVencimiento"
                  lazy
                  transition="scale-transition"
                  offset-y
                  full-width
                  min-width="290px"
                >
                  <v-text-field
                    slot="activator"
                    v-model="nuevaFechaVencimiento"
                    label="Fecha de vencimiento"
                    prepend-icon="event"
                    readonly
                  ></v-text-field>
                  <v-date-picker
                    color="green lighten-1"
                    locale="es-419"
                    v-model="nuevaFechaVencimiento"
                    :min="hoy"
                  >
                    <v-spacer></v-spacer>
                    <v-btn
                      flat
                      color="primary"
                      @click="mostrarDialogoFecha = false"
                      >Cerrar</v-btn
                    >
                    <v-btn
                      flat
                      color="primary"
                      @click="$refs.menu.save(nuevaFechaVencimiento)"
                      >OK</v-btn
                    >
                  </v-date-picker>
                </v-menu>
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
          >Guardar
        </v-btn>
        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { FUNCIONES } from "../../../../funciones";
import { HTTP_AUTH } from "../../../../http-common";

export default {
  name: "CambiarFechaVencimiento",
  props: ["mostrar", "idApartado", "fechaVencimiento"],
  data: () => ({
    nuevaFechaVencimiento: null,
    mostrarDialogoFecha: true,
    hoy: FUNCIONES.hoyComoCadena(),
    cargando: false,
  }),
  watch: {
    mostrar(mostrado) {
      if (mostrado) {
        this.$nextTick(() => {
          this.nuevaFechaVencimiento = this.fechaVencimiento;
        });
      }
    }
  },
  methods: {
    cerrarDialogo() {
      this.$emit("cerrar");
    },
    guardar() {
      if (this.idApartado && this.nuevaFechaVencimiento) {
        this.cargando = true;
        HTTP_AUTH.put(`fecha/apartado/${this.idApartado}`, FUNCIONES.agregarHoraCeroAFecha(this.nuevaFechaVencimiento)).then(respuestaAlCambiarFecha => {
          this.cargando = false;
          this.$emit("cambiada");
        });
      }

    }
  },
}
</script>

