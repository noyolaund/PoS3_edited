<template>
  <v-layout row wrap>
    <v-flex xs12>
      <v-form ref="formulario">
        <v-container fluid grid-list-md>
          <v-layout row wrap>
            <v-flex xs12 sm6>
              <v-text-field
                prepend-icon="info"
                label="Nombre"
                type="text"
                v-model="datos.Nombre"
                hint="Nombre de la tienda"
                required
              ></v-text-field>
            </v-flex>
            <v-flex xs12 sm6>
              <v-text-field
                prepend-icon="location_on"
                label="Dirección"
                type="text"
                v-model="datos.Direccion"
                hint="La dirección que saldrá en el ticket"
                required
              ></v-text-field>
            </v-flex>
            <v-flex xs12 sm6>
              <v-text-field
                prepend-icon="local_phone"
                label="Teléfono"
                type="text"
                v-model="datos.Telefono"
                hint="Número de teléfono para dudas o aclaraciones"
                required
              ></v-text-field>
            </v-flex>
            <v-flex xs12 sm6>
              <v-text-field
                prepend-icon="message"
                label="Mensaje personal"
                type="text"
                v-model="datos.MensajePersonal"
                hint="Texto al final del ticket, por ejemplo 'Gracias por su compra'"
                required
              ></v-text-field>
            </v-flex>
          </v-layout>
          <v-btn
            :loading="cargando"
            @click="guardar()"
            fixed
            dark
            fab
            bottom
            fill-height
            slot="activator"
            right
            color="teal ligthen-1"
          >
            <v-icon>save</v-icon>
          </v-btn>
        </v-container>
      </v-form>
    </v-flex>
  </v-layout>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  beforeMount() {
    this.obtener();
  },
  methods: {
    guardar() {
      this.cargando = true;
      HTTP_AUTH.put("ajustes/empresa", Object.assign({}, this.datos))
        .then(respuesta => {
          this.cargando = false;
          if (respuesta) this.$emit("guardado");
        })
    },
    obtener() {
      HTTP_AUTH.get("ajustes/empresa").then(datosEmpresa => {
        this.datos = Object.assign({}, datosEmpresa);
      });
    },
  },
  data: () => ({
    cargando: false,
    datos: {
      Direccion: "",
      Nombre: "",
      MensajePersonal: "",
      Telefono: ""
    }
  })
}
</script>


