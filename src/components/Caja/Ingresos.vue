<template>
  <v-layout row wrap>
    <v-flex xs12>
      <h1>
        <span class="display-1">{{ total | currency }}</span>
        <span class="title">Ingresos</span>
      </h1>
    </v-flex>
    <v-flex xs12>
      <v-data-table
        :headers="encabezados"
        :items="ingresos"
        hide-actions
        item-key="props.item.index"
      >
        <template slot="items" slot-scope="props">
          <tr>
            <td>{{ props.item.Monto | currency }}</td>
            <td>{{ props.item.Descripcion }}</td>
            <td>{{ props.item.Fecha | fechaExpresiva }}</td>
            <td>{{ props.item.Usuario.Nombre }}</td>
          </tr>
        </template>
      </v-data-table>
    </v-flex>
  </v-layout>
</template>
<script>
import { HTTP_AUTH } from "../../http-common.js";

export default {
  methods: {
    obtener(inicio, fin) {
      /**
       * No arregla ninguna fecha, así que se debe pasar la fecha ya arreglada
       */
      this.total = 0;
      HTTP_AUTH.get(`ingresos/${inicio}/${fin}`).then(ingresos => {
        this.total = ingresos.reduce(
          (acumulador, siguiente) => ({
            Monto: acumulador.Monto + siguiente.Monto
          }),
          {
            Monto: 0
          }
        ).Monto;
        this.ingresos = ingresos;
      });
    }
  },
  data: () => ({
    ingresos: [],
    total: 0,
    encabezados: [
      {
        text: "Monto",
        value: "Monto"
      },
      {
        text: "Descripcion",
        value: "Descripcion"
      },
      {
        text: "Fecha",
        value: "Fecha"
      },
      {
        text: "Usuario",
        value: "Usuario.Nombre"
      }
    ]
  })
};
</script>
