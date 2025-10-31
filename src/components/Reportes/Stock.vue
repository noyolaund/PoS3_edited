<template>
  <v-layout row wrap>
    <v-flex xs12>
      <div class="text-xs-center pt-2">
        <v-pagination
          v-model="paginacion.pagina"
          :length="numeroDePaginas"
        ></v-pagination>
      </div>
    </v-flex>
    <v-flex xs12>
      <v-data-table
        :loading="cargando"
        :headers="encabezados"
        :items="productos"
        hide-actions
        class="elevation-1"
      >
        <template slot="items" slot-scope="props">
          <td>{{ props.item.Numero }}</td>
          <td>{{ props.item.CodigoBarras }}</td>
          <td>{{ props.item.Descripcion }}</td>
          <td>{{ props.item.Existencia }}</td>
          <td>{{ props.item.Stock }}</td>
          <td>
            <strong>{{ props.item.Stock - props.item.Existencia }}</strong>
          </td>
        </template>
        <template slot="no-data">
          <v-alert :value="true" color="success">
            <div class="text-xs-center">
              <h1>No hay productos en stock</h1>
              <v-icon class="icono-grande">tag_faces</v-icon>
              <p>La existencia de los productos está en los niveles óptimos</p>
            </div>
          </v-alert>
        </template>
      </v-data-table>
    </v-flex>
  </v-layout>
</template>
<script>
import { HTTP_AUTH } from "../../http-common";
import { EventBus } from "../../main";

const PRODUCTOS_POR_PAGINA = 20;

export default {
  computed: {
    numeroDePaginas() {
      if (this.paginacion.conteoProductos == null) return 0;
      return Math.ceil(
        this.paginacion.conteoProductos / this.paginacion.limite
      );
    }
  },
  beforeMount() {
    this.obtener();
    EventBus.$emit("ponerTitulo", "Productos en stock");
  },
  watch: {
    paginacion: {
      handler() {
        this.obtener();
      },
      deep: true
    }
  },
  methods: {
    obtener() {
      this.cargando = true;
      let verdaderoOffset =
        this.paginacion.pagina > 1
          ? (this.paginacion.pagina - 1) * this.paginacion.limite
          : 0;
      HTTP_AUTH.get(
        `productos/stock/${verdaderoOffset}/${this.paginacion.limite}`
      ).then(productosConConteo => {
        this.productos = productosConConteo.Productos;
        this.paginacion.conteoProductos = productosConConteo.Conteo;
        this.cargando = false;
      });
    }
  },
  data: () => ({
    cargando: false,
    paginacion: {
      offset: 0,
      limite: PRODUCTOS_POR_PAGINA,
      conteoProductos: 0,
      pagina: 1
    },
    productos: [],
    encabezados: [
      {
        text: "#",
        align: "left",
        value: "Numero"
      },
      {
        text: "Código de barras",
        align: "left",
        value: "CodigoBarras"
      },
      {
        text: "Descripción",
        value: "Descripcion"
      },
      {
        text: "Existencia",
        value: "Existencia"
      },
      {
        text: "Stock",
        value: "Stock"
      },
      {
        text: "Diferencia",
        value: "Diferencia",
        sortable: false
      }
    ]
  })
};
</script>
