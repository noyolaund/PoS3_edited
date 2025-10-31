<template>
  <v-flex xs12>
    <!-- <Publicidad></Publicidad> -->
    <v-flex xs12>
      <v-text-field
        label="Buscar un producto por su descripción"
        v-model="busqueda"
        prepend-icon="search"
        solo
        clearable
      ></v-text-field>
    </v-flex>
    <br />
    <div class="text-xs-center pt-2">
      <v-pagination
        v-model="paginacion.pagina"
        :length="numeroDePaginas"
      ></v-pagination>
    </div>
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
        <td>{{ props.item.PrecioCompra | currency }}</td>
        <td>{{ props.item.PrecioVenta | currency }}</td>
        <td>
          {{ (props.item.PrecioVenta - props.item.PrecioCompra) | currency }}
        </td>
        <td>{{ props.item.Existencia }}</td>
        <td>{{ props.item.Stock }}</td>
        <td class="justify-center layout px-0">
          <v-btn
            title="Crear una copia"
            icon
            class="mx-0"
            @click="duplicar(props.item)"
          >
            <v-icon color="green darken-3">file_copy</v-icon>
          </v-btn>
          <v-btn
            title="Restar existencia"
            icon
            class="mx-0"
            @click="restarExistencia(props.item)"
          >
            <v-icon color="orange darken-3">remove_circle</v-icon>
          </v-btn>
          <v-btn
            title="Aumentar existencia"
            icon
            class="mx-0"
            @click="aumentarExistencia(props.item)"
          >
            <v-icon color="blue darken-3">add_circle</v-icon>
          </v-btn>
          <v-btn
            title="Editar/modificar"
            icon
            class="mx-0"
            @click="editar(props.item)"
          >
            <v-icon color="amber darken-4">edit</v-icon>
          </v-btn>
          <v-btn
            title="Eliminar"
            icon
            class="mx-0"
            @click="eliminar(props.item)"
          >
            <v-icon color="red">delete</v-icon>
          </v-btn>
        </td>
      </template>
      <template slot="no-data">
        <v-alert v-show="!busqueda" :value="true" color="info">
          <div class="text-xs-center">
            <h1>No hay productos</h1>
            <v-icon class="icono-grande">announcement</v-icon>
            <p>
              No has registrado ningún producto. Agrega uno con el botón
              <v-icon>add</v-icon>
              de la esquina
            </p>
          </div>
        </v-alert>
        <v-alert v-show="busqueda" :value="true" color="info" dark>
          <div class="text-xs-center">
            <h1>Sin resultados</h1>
            <v-icon class="icono-grande">highlight_off</v-icon>
            <p>
              No hay productos que coincidan con "
              <strong>{{ busqueda }}</strong
              >"
            </p>
          </div>
        </v-alert>
      </template>
    </v-data-table>
    <exportar-importar @importado="obtener"></exportar-importar>
  </v-flex>
</template>
<script>
import { HTTP_AUTH } from "../../http-common";
import ExportarImportar from "./ExportarImportar"
import Publicidad from "../Publicidad";

const PRODUCTOS_POR_PAGINA = 7;
export default {
  components: { Publicidad, ExportarImportar },
  computed: {
    numeroDePaginas() {
      if (this.paginacion.conteoProductos == null) return 0;
      return Math.ceil(
        this.paginacion.conteoProductos / this.paginacion.limite
      );
    }
  },
  data: () => ({
    deberiaReiniciarPaginacionAlBuscar: false,
    cargando: false,
    paginacion: {
      offset: 0,
      limite: PRODUCTOS_POR_PAGINA,
      conteoProductos: 0,
      pagina: 1
    },
    busqueda: "",
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
        text: "P. compra",
        value: "PrecioCompra"
      },
      {
        text: "P. venta",
        value: "PrecioVenta"
      },
      {
        text: "Utilidad",
        value: "Utilidad"
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
        text: "Opciones",
        sortable: false
      }
    ]
  }),
  watch: {
    busqueda() {
      this.obtener();
    },
    paginacion: {
      handler() {
        this.obtener();
      },
      deep: true
    }
  },
  beforeMount() {
    this.obtener();
  },
  methods: {
    obtener() {
      //TODO: poner debounce en búsqueda
      this.cargando = true;
      let verdaderoOffset =
        this.paginacion.pagina > 1
          ? (this.paginacion.pagina - 1) * this.paginacion.limite
          : 0;
      if (this.busqueda) {
        if (this.deberiaReiniciarPaginacionAlBuscar) {
          this.paginacion.pagina = 1;
          verdaderoOffset = 0;
          this.deberiaReiniciarPaginacionAlBuscar = false;
        }
        HTTP_AUTH.get(
          `buscar/productos/${verdaderoOffset}/${this.paginacion.limite
          }/${encodeURIComponent(this.busqueda)}`
        ).then(resultadosDeBusqueda => {
          this.productos = resultadosDeBusqueda.Productos;
          this.paginacion.conteoProductos = resultadosDeBusqueda.Conteo;
          this.cargando = false;
        });
      } else {
        this.deberiaReiniciarPaginacionAlBuscar = true;
        HTTP_AUTH.get(
          `productos/${verdaderoOffset}/${this.paginacion.limite}`
        ).then(productosConConteo => {
          this.productos = productosConConteo.Productos;
          this.paginacion.conteoProductos = productosConConteo.Conteo;
          this.cargando = false;
        });
      }
    },
    editar(producto) {
      this.$emit("editar-producto", producto);
    },
    eliminar(producto) {
      this.$emit("eliminar-producto", producto);
    },
    aumentarExistencia(producto) {
      this.$emit("aumentar-existencia", producto);
    },
    restarExistencia(producto) {
      this.$emit("restar-existencia", producto);
    },
    duplicar(producto) {
      this.$emit("duplicar-producto", producto);
    },
  }
};
</script>
