<template>
  <v-flex xs12>
    <v-card>
      <v-toolbar color="light-blue" dark>
        <v-toolbar-title v-show="!estaBuscando">Clientes</v-toolbar-title>
        <v-toolbar-title v-show="estaBuscando"></v-toolbar-title>
        <v-spacer></v-spacer>
        <v-flex xs12 v-show="estaBuscando">
          <v-text-field
            v-model="busqueda"
            ref="inputBusqueda"
            hide-details
            label="Escriba parte del nombre"
            solo-inverted
            single-line
          ></v-text-field>
        </v-flex>
        <v-btn @click="prepararParaBuscar()" v-show="!estaBuscando" icon>
          <v-icon>search</v-icon>
        </v-btn>
        <v-btn @click="cancelarBusqueda()" v-show="estaBuscando" icon>
          <v-icon>close</v-icon>
        </v-btn>
      </v-toolbar>
      <v-flex xs12>
        <!-- <Publicidad></Publicidad> -->
      </v-flex>
      <v-list two-line subheader>
        <v-list-tile v-for="(cliente, indice) in clientes" :key="indice" avatar>
          <v-list-tile-content>
            <v-list-tile-title>{{ cliente.Nombre }}</v-list-tile-title>
            <v-list-tile-sub-title>{{
              cliente.NumeroTelefono
            }}</v-list-tile-sub-title>
          </v-list-tile-content>
          <v-list-tile-action>
            <v-btn
              target="_blank"
              :href="'tel:' + cliente.NumeroTelefono"
              icon
              ripple
            >
              <v-icon color="blue">phone</v-icon>
            </v-btn>
          </v-list-tile-action>
          <v-list-tile-action>
            <v-btn @click="mostrarHistorial(cliente.Numero)" icon ripple>
              <v-icon color="blue">info</v-icon>
            </v-btn>
          </v-list-tile-action>
          <v-list-tile-action>
            <v-btn @click="editar(cliente)" icon ripple>
              <v-icon color="orange">edit</v-icon>
            </v-btn>
          </v-list-tile-action>
          <v-list-tile-action>
            <v-btn icon ripple @click="eliminar(cliente)">
              <v-icon color="red">delete</v-icon>
            </v-btn>
          </v-list-tile-action>
        </v-list-tile>
      </v-list>
    </v-card>
  </v-flex>
</template>
<script>
import { HTTP_AUTH } from "../../http-common";
import Publicidad from "../Publicidad";

export default {
  components: { Publicidad },
  beforeMount() {
    this.obtener();
  },
  watch: {
    busqueda() {
      this.obtener();
    }
  },
  methods: {
    mostrarHistorial(numeroCliente) {
      this.$emit("mostrar-historial", numeroCliente);
    },
    editar(cliente) {
      this.$emit("editar", cliente);
    },
    eliminar(cliente) {
      this.$emit("eliminar", cliente);
    },
    prepararParaBuscar() {
      this.estaBuscando = true;
      this.$nextTick(this.$refs.inputBusqueda.focus);
    },
    cancelarBusqueda() {
      this.busqueda = "";
      this.estaBuscando = false;
    },
    obtener() {
      if (this.busqueda) {
        HTTP_AUTH.get(
          `buscar/clientes/${encodeURIComponent(this.busqueda)}`
        ).then(clientes => {
          this.clientes = clientes;
        });
      } else {
        HTTP_AUTH.get("clientes").then(clientes => {
          this.clientes = clientes;
        });
      }
    }
  },
  data: () => ({
    busqueda: "",
    estaBuscando: false,
    clientes: []
  })
};
</script>
