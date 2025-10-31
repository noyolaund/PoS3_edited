<template>
  <v-layout row wrap>
    <v-layout row wrap v-show="productosParaVerificar.length <= 0">
      <v-flex xs12>
        <v-alert :value="true" type="info">
          Seleccione una cantidad de productos que desea verificar o contar;
          serán seleccionados aleatoriamente
        </v-alert>
      </v-flex>
      <v-flex xs12 sm6>
        <v-text-field
          label="Productos  para verificar. P.ej.: 50"
          type="number"
          ref="cantidadProductos"
          hint="¿Cuántos productos va a verificar?"
          v-model.number="cantidadDeProductosParaVerificar"
          @keyup.enter="obtenerProductosParaVerificar()"
          required
        ></v-text-field>
        <v-btn
          :loading="cargando"
          dark
          small
          color="green darken-1"
          @click.native="obtenerProductosParaVerificar()"
        >
          Comenzar
        </v-btn>
      </v-flex>
    </v-layout>
    <v-layout row wrap v-show="productosParaVerificar.length > 0">
      <v-flex xs12>
        <v-btn
          :loading="cargando"
          dark
          small
          color="green darken-1"
          @click.native="volverASeleccionarProductos()"
        >
          <v-icon>arrow_back</v-icon>
          Volver a seleccionar productos
        </v-btn>
        <v-btn
          :loading="cargando"
          dark
          small
          color="info"
          @click.native="obtenerProductosParaVerificar()"
        >
          <v-icon>refresh</v-icon>
          Refrescar
        </v-btn>
      </v-flex>
      <v-flex xs12 sm6>
        <v-text-field
          label="Escanear código de barras o número de producto"
          type="number"
          ref="inputProducto"
          hint="Utilice el lector de códigos, o escriba el número y presione Enter"
          v-model="codigoEscaneado"
          @keyup.enter="aumentarEncontradoSiExiste(codigoEscaneado)"
          required
        ></v-text-field>
      </v-flex>
      <v-flex xs12 v-show="log.length > 0">
        <v-btn
          dark
          small
          color="error"
          @click.native="mostrarLog = !mostrarLog"
        >
          <v-icon>{{ mostrarLog ? "expand_less" : "expand_more" }}</v-icon>
          {{ mostrarLog ? "Ocultar" : "Mostrar" }} registro de errores
        </v-btn>
        <v-card
          v-show="mostrarLog"
          style="max-height: 200px; overflow-y: scroll"
        >
          <v-card-text>
            <ul class="lista-log">
              <li v-for="(mensaje, i) in log" :key="i">{{ mensaje }}</li>
            </ul>
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs12>
        <v-data-table
          :loading="cargando"
          :headers="encabezados"
          :items="productosParaVerificar"
          hide-actions
          class="elevation-1"
        >
          <template slot="items" slot-scope="props">
            <td>{{ props.item.Numero }}</td>
            <td>{{ props.item.Descripcion }}</td>
            <td>{{ props.item.CodigoBarras }}</td>
            <td>{{ props.item.Existencia }}</td>
            <td>
              <strong>{{ props.item.Encontrados }}</strong>
            </td>
            <td>
              <strong>{{
                props.item.Existencia - props.item.Encontrados
              }}</strong>
            </td>
          </template>
        </v-data-table>
      </v-flex>
    </v-layout>
  </v-layout>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";
import { EventBus } from "../../main";


export default {
  name: "HacerInventario",
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Hacer inventario");
    HTTP_AUTH.get("ajustes/otros").then(ajustes => {
      this.modoDeLectura = ajustes.ModoLecturaProductos;
      this.enfocarInputParaCantidadDeProductos();
    });
  },
  data: () => ({
    mostrarLog: false,
    codigoEscaneado: "",
    encabezados: [
      {
        text: "No.",
        value: "Numero"
      },
      {
        text: "Descripción",
        value: "Descripcion"
      },
      {
        text: "Código de barras",
        value: "CodigoBarras"
      },
      {
        text: "Existencia",
        value: "Existencia"
      },
      {
        text: 'Encontrados',
        value: 'Encontrados'
      },
      {
        text: 'Faltantes',
        sortable: 1 === 0,
      }
    ],
    cantidadDeProductosParaVerificar: null,
    productosParaVerificar: [],
    cargando: false,
    log: [],
    modoDeLectura: "",
  }),
  methods: {
    volverASeleccionarProductos() {
      this.productosParaVerificar = [];
      this.log = [];
      this.enfocarInputParaCantidadDeProductos();
    },
    enfocarInputParaBuscarProducto() {
      this.$nextTick(this.$refs.inputProducto.focus);
    },
    enfocarInputParaCantidadDeProductos() {
      this.$nextTick(this.$refs.cantidadProductos.focus);
    },
    aumentarEncontradoSiExiste(codigo) {
      if (this.modoDeLectura !== "") {
        let indice = -1;
        switch (this.modoDeLectura) {
          case "codigo":
            indice = this.productosParaVerificar.findIndex(producto => producto.CodigoBarras === codigo);
            break;
          case "numero":
            indice = this.productosParaVerificar.findIndex(producto => producto.Numero === codigo);
        }
        if (indice !== -1) {
          let producto = this.productosParaVerificar[indice];
          if (producto.Encontrados < producto.Existencia) {
            producto.Encontrados++;
          } else {
            this.log.push(`Se encontró el producto ${producto.Descripcion} con el código ${codigo}, pero ya estaba completo`);
          }
        } else {
          this.log.push(`Se intentó buscar el producto con el código ${codigo}, pero no existe`);
        }

        this.codigoEscaneado = "";
      }
    },
    obtenerProductosParaVerificar() {
      if (this.cantidadDeProductosParaVerificar > 0) {
        this.log = [];
        this.cargando = true;
        HTTP_AUTH.get(`productos/aleatorios/${this.cantidadDeProductosParaVerificar}`)
          .then(productos => {
            this.cargando = false;
            this.productosParaVerificar = productos.map(producto => {
              return {
                Numero: producto.Numero,
                Descripcion: producto.Descripcion,
                CodigoBarras: producto.CodigoBarras,
                Existencia: producto.Existencia,
                Encontrados: 0
              }
            });
            this.enfocarInputParaBuscarProducto();
          })
      }
    },
  },
}
</script>

<style scoped>
ul.lista-log {
  list-style-type: none;
  padding: 0;
  margin: 0;
}
</style>
