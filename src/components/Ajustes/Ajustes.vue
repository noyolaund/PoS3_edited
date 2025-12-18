<template>
  <v-layout row wrap>
    <v-flex xs12>
      <!-- <Publicidad></Publicidad> -->
    </v-flex>
    <v-flex xs12>
      <v-tabs icons-and-text centered dark color="teal">
        <v-tabs-slider color="yellow"></v-tabs-slider>
        <v-tab href="#tab-1">
          Datos negocio
          <v-icon>subject</v-icon>
        </v-tab>
        <v-tab href="#tab-2">
          Impresora
          <v-icon>local_printshop</v-icon>
        </v-tab>
        <v-tab href="#tab-3">
          Otros
          <v-icon>build</v-icon>
        </v-tab>
        <v-tab href="#tab-4">
          Ticket
          <v-icon>receipt</v-icon>
        </v-tab>
        <v-tab-item value="tab-1">
          <DatosNegocio @guardado="guardado = true"></DatosNegocio>
        </v-tab-item>
        <v-tab-item value="tab-2">
          <Impresora></Impresora>
        </v-tab-item>
        <v-tab-item value="tab-3">
          <Otros @guardado="guardado = true"></Otros>
        </v-tab-item>
        <v-tab-item value="tab-4">
          <PersonalizacionTicket @guardado="guardado = true"></PersonalizacionTicket>
        </v-tab-item>
      </v-tabs>
      <v-snackbar
        :timeout="2000"
        :right="true"
        :bottom="true"
        v-model="guardado"
      >
        Ajustes guardados
        <v-btn flat color="pink" @click.native="guardado = false">OK</v-btn>
      </v-snackbar>
    </v-flex>
  </v-layout>
</template>

<script>
import DatosNegocio from './DatosNegocio'
import Impresora from './Impresora'
import Otros from './Otros'
import PersonalizacionTicket from './PersonalizacionTicket'
import { EventBus } from "../../main";
import Publicidad from "../Publicidad";

export default {
  components: { Publicidad, DatosNegocio, Impresora, Otros, PersonalizacionTicket },
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Ajustes");
  },
  data: () => {
    return {
      guardado: false,
    }
  }
}
</script>


