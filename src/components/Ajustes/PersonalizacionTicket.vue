<template>
  <v-layout row wrap>
    <v-flex xs12>
      <v-form ref="formulario">
        <v-container fluid grid-list-md>
          <v-layout row wrap>
            <v-flex xs12>
              <v-text-field
                prepend-icon="text_fields"
                label="Texto personalizado del ticket"
                type="text"
                v-model="textoTicket"
                hint="Este texto aparecerá en lugar de 'Deposito Beer Broos'"
                persistent-hint
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
  name: "PersonalizacionTicket",
  beforeMount() {
    this.obtener();
  },
  methods: {
    guardar() {
      this.cargando = true;
      HTTP_AUTH.put("valor", {
        Clave: "TEXTO_TICKET_CUSTOM",
        Valor: this.textoTicket
      })
        .then(respuesta => {
          this.cargando = false;
          if (respuesta) this.$emit("guardado");
        })
    },
    obtener() {
      HTTP_AUTH.get("valor/TEXTO_TICKET_CUSTOM").then(valor => {
        this.textoTicket = valor;
      });
    },
  },
  data: () => ({
    cargando: false,
    textoTicket: ""
  })
}
</script>
