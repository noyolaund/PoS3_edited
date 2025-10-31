<template>
  <v-layout row wrap>
    <cambiar-fecha-vencimiento
      @cerrar="dialogos.cambiarFechaVencimiento = false"
      @cambiada="onFechaVencimientoCambiada"
      :mostrar="dialogos.cambiarFechaVencimiento"
      :fechaVencimiento="apartadoEditado.fechaVencimiento"
      :idApartado="apartadoEditado.id"
    >
    </cambiar-fecha-vencimiento>
    <seleccionador-fechas
      v-show="tipoReporte === 'todos'"
      @cambio="comprobarFechasYRefrescarSiEsNecesario"
    >
    </seleccionador-fechas>
    <v-flex xs12>
      <v-layout row wrap>
        <v-flex xs12 sm6>
          <v-text-field
            v-model="numeroApartadoParaAbonar"
            ref="nombreCliente"
            :loading="buscandoApartadoPorNumero"
            @keyup.enter="mostrarAbonosDesdeInput()"
            prepend-icon="monetization_on"
            label="Abonar a un apartado"
            type="number"
            hint="Escriba el número y presione Enter"
            required
          ></v-text-field>
        </v-flex>
      </v-layout>
    </v-flex>
    <v-flex xs12 sm6 v-show="tipoReporte === 'todos'">
      <h1>
        <span class="display-1">{{ totales.anticipo | currency }}</span>
        <span class="title">Anticipo</span>
      </h1>
    </v-flex>
    <v-flex xs12 sm6 v-show="tipoReporte === 'todos'">
      <h1>
        <span class="display-1">{{ totales.abonado | currency }}</span>
        <span class="title">Abonado</span>
      </h1>
    </v-flex>
    <v-flex v-show="tipoReporte === 'todos'" xs12>
      <v-divider></v-divider>
    </v-flex>
    <productos-apartados
      @cerrar="dialogos.productos = false"
      :mostrar="dialogos.productos"
      :apartado="apartadoParaMostrarDetalles"
    >
    </productos-apartados>
    <abonos
      @abonado="onApartadoAbonado"
      @liquidar="snackbars.entregarProductos = true"
      @cantidadSuperior="snackbars.cantidadSuperior = true"
      @cantidadNegativa="snackbars.cantidadNegativa = true"
      :apartado="apartadoSeleccionado"
      :mostrar="dialogos.abonos"
      @cerrar="dialogos.abonos = false"
    >
    </abonos>
    <v-flex xs12>
      <v-radio-group label="Mostrar..." v-model="tipoReporte" row>
        <v-radio label="Todos" value="todos"></v-radio>
        <v-radio label="Pendientes" value="pendientes"></v-radio>
        <v-radio
          label="Próximos a vencer (1 semana)"
          value="proximos"
        ></v-radio>
      </v-radio-group>
    </v-flex>
    <v-flex xs12>
      <v-data-table
        :headers="encabezadosApartados"
        :items="apartados"
        hide-actions
        item-key="IdApartado"
      >
        <template slot="items" slot-scope="props">
          <tr>
            <td>{{ props.item.IdApartado }}</td>
            <td>{{ props.item.Monto | currency }}</td>
            <td>{{ props.item.Anticipo | currency }}</td>
            <td>{{ props.item.Abonado | currency }}</td>
            <td>
              {{
                (props.item.Monto - props.item.Anticipo - props.item.Abonado)
                  | currency
              }}
            </td>
            <td>{{ props.item.Fecha | fechaSinHora }}</td>
            <td>{{ props.item.FechaVencimiento | fechaSinHora }}</td>
            <td>{{ props.item.Cliente.Nombre }}</td>
            <td>{{ props.item.Usuario.Nombre }}</td>
            <td class="justify-center layout px-0">
              <v-btn
                title="Cambiar fecha de vencimiento"
                icon
                class="mx-0"
                @click="cambiarFechaVencimiento(props.item)"
              >
                <v-icon color="error">update</v-icon>
              </v-btn>
              <v-btn
                title="Detalles"
                icon
                class="mx-0"
                @click="mostrarDetalles(props.item)"
              >
                <v-icon color="blue">info</v-icon>
              </v-btn>
              <v-btn
                title="Abonar"
                icon
                class="mx-0"
                @click="mostrarAbonos(props.item)"
              >
                <v-icon color="green">monetization_on</v-icon>
              </v-btn>
              <v-btn
                title="Imprimir ticket"
                icon
                class="mx-0"
                @click="imprimir(props.item)"
              >
                <v-icon color="orange">print</v-icon>
              </v-btn>
            </td>
          </tr>
        </template>
      </v-data-table>
    </v-flex>
    <v-snackbar
      :timeout="2000"
      :top="true"
      v-model="snackbars.cantidadNegativa"
    >
      La cantidad abonada debe ser mayor que cero
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.cantidadNegativa = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :top="true"
      v-model="snackbars.cantidadSuperior"
    >
      Introduce una cantidad menor o igual al restante
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.cantidadSuperior = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar :timeout="0" :top="true" v-model="snackbars.entregarProductos">
      Apartado liquidado. Ahora puede entregar los productos
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.entregarProductos = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="3000"
      :bottom="true"
      v-model="snackbars.fechaVencimientoCambiada"
    >
      Fecha de vencimiento cambiada correctamente
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.fechaVencimientoCambiada = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="3000"
      :bottom="true"
      v-model="snackbars.numeroApartadoInexistenteOErroneo"
    >
      No existe un apartado con ese número
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.numeroApartadoInexistenteOErroneo = false"
        >OK</v-btn
      >
    </v-snackbar>
  </v-layout>
</template>
<script>
import { HTTP_AUTH } from "../../../http-common";
import SeleccionadorFechas from "../../Reportes/SeleccionadorFechas";
import ProductosApartados from "./ProductosApartados";
import Abonos from "./Abonos";
import { EventBus } from "../../../main";
import CambiarFechaVencimiento from './Apartados/CambiarFechaVencimiento'
import { FUNCIONES } from '../../../funciones';

export default {
  components: {
    SeleccionadorFechas,
    ProductosApartados,
    Abonos,
    CambiarFechaVencimiento
  },
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Reporte de apartados");
  },
  watch: {
    tipoReporte() {
      this.refrescarApartadosConElTipoSeleccionado();
    }
  },
  methods: {
    mostrarAbonosDesdeInput() {

      let verdaderoNumero = Number(this.numeroApartadoParaAbonar);
      if (!isNaN(verdaderoNumero) && verdaderoNumero > 0) {
        this.buscandoApartadoPorNumero = true;
        HTTP_AUTH.get(`apartado/${verdaderoNumero}`).then(apartado => {
          this.numeroApartadoParaAbonar = "";
          this.buscandoApartadoPorNumero = false;
          if (null !== apartado && apartado.Total > 0) {
            /**
             * Renombrar propiedades, ya que <b>Apartado</b> y <b>DetalleApartado</b>
             * son diferentes tipos. Y <b>mostrarAbonos</b> espera un objeto
             * de tipo <b>DetalleApartado</b>
             *
             * Para más información, consultar <i>api/Apartado.go</i> y <i>api/DetalleApartado.go</i>
             * */
            apartado.IdApartado = apartado.Numero;
            delete apartado.Numero;
            apartado.Monto = apartado.Total;
            delete apartado.Total;
            this.mostrarAbonos(apartado);
          } else {
            this.snackbars.numeroApartadoInexistenteOErroneo = true;
          }
        });
      } else {
        this.snackbars.numeroApartadoInexistenteOErroneo = true;
      }
    },
    onFechaVencimientoCambiada() {
      this.snackbars.fechaVencimientoCambiada = true;
      this.dialogos.cambiarFechaVencimiento = false;
      this.refrescarApartadosConElTipoSeleccionado();
    },
    refrescarApartadosConElTipoSeleccionado() {
      switch (this.tipoReporte) {
        case "todos":
          this.consultarApartados(this.ultimaFechaInicio, this.ultimaFechaFin);
          break;
        case "pendientes":
          this.consultarApartadosPendientes();
          break;
        case "proximos":
          this.consultarApartadosProximos();
          break;
      }
    },
    cambiarFechaVencimiento({ IdApartado, FechaVencimiento }) {
      this.dialogos.cambiarFechaVencimiento = true;
      this.apartadoEditado.id = IdApartado;
      this.apartadoEditado.fechaVencimiento = FechaVencimiento.split("T")[0];//Remover T00:00:00
    },
    async imprimir(apartado) {
      await FUNCIONES.imprimirTicketApartado(apartado.IdApartado);
    },
    onApartadoAbonado() {
      this.refrescarApartadosConElTipoSeleccionado();
    },
    mostrarDetalles(apartado) {
      this.dialogos.productos = true;
      this.apartadoParaMostrarDetalles = apartado;
    },
    mostrarAbonos(apartado) {

      this.apartadoSeleccionado = apartado;
      this.dialogos.abonos = true;
    },
    comprobarFechasYRefrescarSiEsNecesario({ inicio, fin }) {
      if (inicio && fin) {
        this.consultarApartados(inicio, fin);
      }
    },
    consultarApartados(fechaInicio, fechaFin) {
      this.ultimaFechaInicio = fechaInicio;
      this.ultimaFechaFin = fechaFin;
      //Simplemente cambiamos para que el watch se encargue de llamar a esta función de nuevo
      //así evitamos una doble petición
      if (this.mostrarPendientes) return (this.mostrarPendientes = false);
      HTTP_AUTH.get(`apartados/${fechaInicio}/${fechaFin}`)
        .then(apartados => {
          this.procesarApartados(apartados);
        })
        .then(() => {
          return HTTP_AUTH.get(`total/abonado/${fechaInicio}/${fechaFin}`);
        })
        .then(totalAbonadoEnPeriodo => {
          this.totales.abonado = totalAbonadoEnPeriodo;
        });
    },
    consultarApartadosPendientes() {
      HTTP_AUTH.get("apartados/pendientes").then(apartados => {
        this.procesarApartados(apartados);
      });
    },
    consultarApartadosProximos() {
      HTTP_AUTH.get("apartados/proximos/vencer").then(apartados => {
        this.procesarApartados(apartados);
      });
    },
    procesarApartados(apartadosRaw) {
      /**
       * Se encarga de agrupar los productos dentro de un apartado
       */
      let apartados = [];
      this.totales.anticipo = 0;
      this.totales.abonado = 0;
      apartadosRaw.forEach(apartadoRaw => {
        let apartadoExistente = apartados.find(
          apartado => apartado.IdApartado === apartadoRaw.IdApartado
        );
        if (apartadoExistente) {
          apartadoExistente.productos.push(apartadoRaw.Producto);
          let utilidad =
            (apartadoRaw.Producto.PrecioVenta -
              apartadoRaw.Producto.PrecioCompra) *
            apartadoRaw.Producto.Cantidad;
          apartadoExistente.Utilidad += utilidad;
        } else {
          let nuevoApartado = Object.assign({}, apartadoRaw);
          let { Producto } = nuevoApartado;
          nuevoApartado.Utilidad =
            (Producto.PrecioVenta - Producto.PrecioCompra) * Producto.Cantidad;
          this.totales.anticipo += nuevoApartado.Anticipo;
          delete nuevoApartado.Producto;
          nuevoApartado.productos = [Producto];
          apartados.push(nuevoApartado);
        }
      });
      this.apartados = apartados;
    }
  },
  data: () => ({
    buscandoApartadoPorNumero: false,
    numeroApartadoParaAbonar: "",
    tipoReporte: "todos",
    apartadoParaMostrarDetalles: {},
    apartadoEditado: {
      id: null,
      fechaVencimiento: null
    },

    snackbars: {
      cantidadNegativa: false,
      cantidadSuperior: false,
      entregarProductos: false,
      fechaVencimientoCambiada: false,
      numeroApartadoInexistenteOErroneo: false,
    },
    ultimaFechaInicio: null,
    ultimaFechaFin: null,
    apartadoSeleccionado: {},
    dialogos: {
      productos: false,
      abonos: false,
      cambiarFechaVencimiento: false,
    },
    apartados: [],
    totales: {
      anticipo: 0,
      abonado: 0
    },
    encabezadosApartados: [
      {
        text: "#",
        value: "IdVenta"
      },
      {
        text: "Monto",
        value: "Monto"
      },
      {
        text: "Anticipo",
        value: "Anticipo"
      },
      {
        text: "Abonado",
        value: "Abonado"
      },
      {
        text: "Restante",
        value: "Restante",
        sortable: false
      },
      {
        text: "Fecha",
        value: "Fecha"
      },
      {
        text: "Vence el",
        value: "FechaVencimiento"
      },
      {
        text: "Cliente",
        value: "IdCliente"
      },
      {
        text: "Usuario",
        value: "IdUsuario"
      },
      {
        text: "Opciones",
        sortable: false
      }
    ]
  })
};
</script>
