<template>
  <v-layout row wrap>
    <seleccionador-fechas
      @cambio="comprobarFechasYRefrescarSiEsNecesario"
      ref="seleccionadorFechas"
    ></seleccionador-fechas>
    <v-flex xs12>
      <!-- <Publicidad></Publicidad> -->
    </v-flex>
    <v-flex xs12 sm6>
      <h1>
        <span class="display-1">{{ totales.utilidad | currency }}</span>
        <span class="title">Utilidad</span>
      </h1>
    </v-flex>
    <v-flex xs12 sm6>
      <h1>
        <span class="display-1">{{ totales.venta | currency }}</span>
        <span class="title">Vendido</span>
      </h1>
    </v-flex>
    <v-flex xs12>
      <v-data-table
        :headers="encabezadosVentas"
        :items="ventas"
        hide-actions
        item-key="IdVenta"
      >
        <template slot="items" slot-scope="props">
          <tr
            class="cursor-manita"
            :class="{ expandido: props.expanded }"
            @click="props.expanded = !props.expanded"
          >
            <td>{{ props.item.IdVenta }}</td>
            <td>{{ props.item.Monto | currency }}</td>
            <td>
              <strong>{{ props.item.Utilidad | currency }}</strong>
            </td>
            <td>{{ props.item.Fecha | fechaExpresiva }}</td>
            <td>{{ props.item.Cliente.Nombre }}</td>
            <td>{{ props.item.Usuario.Nombre }}</td>
            <td class="justify-center layout px-0">
              <v-btn
                title="Reimprimir ticket"
                icon
                class="mx-0"
                @click="imprimirTicket(props.item.IdVenta)"
              >
                <v-icon :color="props.expanded ? 'white' : 'blue'"
                  >print</v-icon
                >
              </v-btn>
            </td>
            <td class="justify-center">
              <v-btn
                title="Anular venta"
                icon
                class="mx-0"
                @click="anularVenta(props.item.IdVenta)"
              >
                <v-icon :color="props.expanded ? 'white' : 'red'"
                  >delete</v-icon
                >
              </v-btn>
            </td>
          </tr>
        </template>
        <template slot="expand" slot-scope="props">
          <v-card flat class="productos-vendidos">
            <v-data-table
              :headers="encabezadosProductos"
              :items="props.item.productos"
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
                </tr>
              </template>
            </v-data-table>
          </v-card>
        </template>
      </v-data-table>
    </v-flex>
  </v-layout>
</template>
<style>
.productos-vendidos {
  padding-left: 2em;
}

tr.expandido {
  background-color: #2196f3;
  color: white;
}

tr.expandido:hover {
  background-color: #2196f3 !important;
  color: white !important;
}
</style>

<script>
import { HTTP_AUTH } from "../../../http-common";
import SeleccionadorFechas from "../../Reportes/SeleccionadorFechas";
import { EventBus } from "../../../main";
import Publicidad from "../../Publicidad";
import { FUNCIONES } from '../../../funciones';

export default {
  components: {
    Publicidad,
    SeleccionadorFechas,
  },
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Ventas al contado");
  },
  data() {
    return {
      totales: {
        utilidad: 0,
        venta: 0,
      },
      ventas: [],
      encabezadosVentas: [
        {
          text: "#",
          value: "IdVenta",
        },
        {
          text: "Monto",
          value: "Monto",
        },
        {
          text: "Utilidad",
          value: "Utilidad",
        },
        {
          text: "Fecha",
          value: "Fecha",
        },
        {
          text: "Cliente",
          value: "IdCliente",
        },
        {
          text: "Usuario",
          value: "IdUsuario",
        },
        {
          text: "Reimprimir ticket",
          value: "",
          sortable: false,
        },
        {
          text: "Anular venta",
          value: "",
          sortable: false,
        },
      ],
      encabezadosProductos: [
        {
          text: "#",
          value: "Numero",
        },
        {
          text: "Código de Barras",
          value: "CodigoBarras",
        },
        {
          text: "Descripción",
          value: "Descripcion",
        },
        {
          text: "Precio de venta original",
          value: "PrecioVentaOriginal",
        },
        {
          text: "Precio de venta",
          value: "PrecioVenta",
        },
        {
          text: "Precio de compra",
          value: "PrecioCompra",
        },
        {
          text: "Utilidad",
          value: "Utilidad",
          sortable: false,
        },
        {
          text: "Cantidad",
          value: "Cantidad",
        },
        {
          text: "Total",
          value: "Total",
        },
      ],
    };
  },
  methods: {
    async anularVenta(idVenta) {
      if (
        !confirm(
          "¿Realmente quiere anular la venta? los productos volverán al inventario y la venta se va a eliminar"
        )
      ) {
        return;
      }
      await HTTP_AUTH.delete("venta/contado/" + idVenta);
      this.$refs.seleccionadorFechas.onFechasCambiadas();
    },
    async imprimirTicket(idVenta) {
      await FUNCIONES.imprimirTicketVentaContado(idVenta);
    },
    comprobarFechasYRefrescarSiEsNecesario({ inicio, fin }) {
      if (inicio && fin) {
        this.consultarVentas(inicio, fin);
      }
    },
    consultarVentas(fechaInicio, fechaFin) {
      HTTP_AUTH.get(`ventas/contado/${fechaInicio}/${fechaFin}`).then(
        (ventasRaw) => {
          this.procesarVentas(ventasRaw);
        }
      );
    },
    procesarVentas(ventasRaw) {
      /**
       * Se encarga de agrupar los productos dentro de una venta
       */
      let ventas = [];
      this.totales.utilidad = 0;
      this.totales.venta = 0;
      ventasRaw.forEach((ventaRaw) => {
        let ventaExistente = ventas.find(
          (venta) => venta.IdVenta === ventaRaw.IdVenta
        );
        if (ventaExistente) {
          ventaExistente.productos.push(ventaRaw.Producto);
          let utilidad =
            (ventaRaw.Producto.PrecioVenta - ventaRaw.Producto.PrecioCompra) *
            ventaRaw.Producto.Cantidad;
          ventaExistente.Utilidad += utilidad;

          this.totales.utilidad += utilidad;
          this.totales.venta +=
            ventaRaw.Producto.PrecioVenta * ventaRaw.Producto.Cantidad;
        } else {
          let nuevaVenta = Object.assign({}, ventaRaw);
          let { Producto } = nuevaVenta;
          nuevaVenta.Utilidad =
            (Producto.PrecioVenta - Producto.PrecioCompra) * Producto.Cantidad;
          this.totales.utilidad += nuevaVenta.Utilidad;
          this.totales.venta += Producto.PrecioVenta * Producto.Cantidad;
          delete nuevaVenta.Producto;
          nuevaVenta.productos = [Producto];
          ventas.push(nuevaVenta);
        }
      });
      this.ventas = ventas;
    },
  },
};
</script>
