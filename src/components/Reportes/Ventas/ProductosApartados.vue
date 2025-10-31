<template>
  <v-container>
    <dialogo-cambiar-producto
      @cerrar="dialogos.buscarProductoParaIntercambiar = false"
      @producto-cambiado="onProductoCambiado"
      :producto="productoParaIntercambiar"
      :idApartado="idApartadoParaCambiar"
      :mostrar="dialogos.buscarProductoParaIntercambiar"
    ></dialogo-cambiar-producto>
    <v-dialog v-model="mostrar && !dialogos.buscarProductoParaIntercambiar">
      <v-card>
        <v-card-title class="headline">Productos apartados</v-card-title>
        <v-card-text>
          <v-data-table
            :headers="encabezadosProductos"
            :items="apartado.productos"
            hide-actions
            item-key="Numero"
          >
            <template slot="items" slot-scope="props">
              <tr>
                <td>{{ props.item.Numero }}</td>
                <td>{{ props.item.CodigoBarras }}</td>
                <td>{{ props.item.Descripcion }}</td>
                <td>{{ props.item.PrecioVentaOriginal | currency }}</td>
                <td>{{ props.item.PrecioVenta | currency }}</td>
                <td>{{ props.item.PrecioCompra | currency }}</td>
                <td>
                  <strong>{{
                    (props.item.PrecioVenta - props.item.PrecioCompra)
                      | currency
                  }}</strong>
                </td>
                <td>{{ props.item.Cantidad }}</td>
                <td>
                  {{
                    (props.item.PrecioVenta * props.item.Cantidad) | currency
                  }}
                </td>
                <td class="justify-center layout px-0">
                  <v-btn
                    title="Cambiar producto"
                    icon
                    class="mx-0"
                    @click="cambiarProducto(props.item)"
                  >
                    <v-icon color="error">find_replace</v-icon>
                  </v-btn>
                </td>
              </tr>
            </template>
          </v-data-table>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="green darken-1"
            flat="flat"
            @click.native="ocultarDialogo()"
            >Cerrar</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
<script>
import DialogoCambiarProducto from '../../Otros/DialogoCambiarProducto'
import { HTTP_AUTH } from "../../../http-common";

export default {
  props: ["apartado", "mostrar"],
  components: { DialogoCambiarProducto },
  methods: {
    onProductoCambiado() {
      this.dialogos.buscarProductoParaIntercambiar = false;
      this.refrescarProductos();
    },
    refrescarProductos() {
      HTTP_AUTH.get(`productos/apartado/${this.apartado.IdApartado}`)
        .then(productos => {
          this.apartado.Monto = productos
            .reduce(
              (acumulador, siguiente) =>
              ({
                PrecioVenta: acumulador.PrecioVenta + siguiente.PrecioVenta
              }),
              { PrecioVenta: 0 }).PrecioVenta;
          this.apartado.Utilidad = productos
            .reduce(
              (acumulador, siguiente) =>
              ({
                PrecioVenta: (acumulador.PrecioVenta - acumulador.PrecioCompra) + (siguiente.PrecioVenta - siguiente.PrecioCompra)
              }),
              { PrecioVenta: 0, PrecioCompra: 0 }).PrecioVenta;
          this.apartado.productos = productos;
        })
    },
    ocultarDialogo() {
      this.$emit("cerrar");
    },
    cambiarProducto(producto) {
      this.productoParaIntercambiar = producto;
      this.idApartadoParaCambiar = this.apartado.IdApartado;
      this.dialogos.buscarProductoParaIntercambiar = true;
    },
  },
  data: () => ({
    productoParaIntercambiar: {},
    idApartadoParaCambiar: 0,
    dialogos: {
      buscarProductoParaIntercambiar: false,
    },
    encabezadosProductos: [
      {
        text: "#",
        value: "Numero"
      },
      {
        text: "Código de Barras",
        value: "CodigoBarras"
      },
      {
        text: "Descripción",
        value: "Descripcion"
      },
      {
        text: "Precio de venta original",
        value: "PrecioVentaOriginal"
      },
      {
        text: "Precio de venta",
        value: "PrecioVenta"
      },
      {
        text: "Precio de compra",
        value: "PrecioCompra"
      },
      {
        text: "Utilidad",
        value: "Utilidad",
        sortable: false
      },
      {
        text: "Cantidad",
        value: "Cantidad"
      },
      {
        text: "Total",
        value: "Total"
      },
      {
        text: "Opciones",
        sortable: false,
      }
    ]
  })
};
</script>
