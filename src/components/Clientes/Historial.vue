<template>
  <v-dialog v-model="mostrar" persistent max-width="800">
    <v-card>
      <v-card-title class="headline">Historial de cliente</v-card-title>
      <v-card-text>
        <v-tabs v-model="tabActiva" color="cyan" dark slider-color="yellow">
          <v-tab key="0" ripple> Resumen </v-tab>
          <v-tab-item key="0">
            <v-card flat>
              <v-card-text>
                <v-layout row wrap>
                  <v-flex xs12>
                    <h3 class="headline">
                      <strong>{{ resumen.cantidadContado }}</strong>
                      {{
                        resumen.cantidadContado <= 0
                          ? "ventas"
                          : resumen.cantidadContado === 1
                          ? "venta"
                          : "ventas"
                      }}
                      al contado, con valor de
                      <strong>{{ resumen.contado | currency }}</strong>
                    </h3>
                    <h3 class="headline">
                      <strong>{{ resumen.cantidadApartados }}</strong>
                      {{
                        resumen.cantidadApartados <= 0 ||
                        resumen.cantidadApartados > 1
                          ? "apartados"
                          : "apartado"
                      }}, con valor de
                      <strong>{{ resumen.apartadoYAbonado | currency }}</strong
                      >. Adeuda
                      <strong>{{
                        resumen.cuentasPendientes | currency
                      }}</strong>
                    </h3>
                    <h3 class="headline">
                      Total comprado:
                      <strong>{{
                        (resumen.contado + resumen.apartadoYAbonado) | currency
                      }}</strong>
                    </h3>
                  </v-flex>
                </v-layout>
              </v-card-text>
            </v-card>
          </v-tab-item>
          <v-tab key="1" ripple> Ventas al contado </v-tab>
          <v-tab-item key="1">
            <v-card flat>
              <v-card-text>
                <v-data-table
                  :loading="cargando"
                  :headers="tabla.encabezados.ventasAlContado"
                  :items="historial.VentasAlContado"
                  hide-actions
                  class="elevation-1"
                >
                  <template slot="items" slot-scope="props">
                    <td>{{ props.item.Numero }}</td>
                    <td>{{ props.item.Total | currency }}</td>
                    <td>{{ props.item.Fecha | fechaExpresiva }}</td>
                  </template>
                  <template slot="no-data">
                    <v-alert :value="true" color="info">
                      <div class="text-xs-center">
                        <h1>Sin registros</h1>
                        <v-icon class="icono-grande">announcement</v-icon>
                        <p>Este cliente no ha realizado ventas al contado</p>
                      </div>
                    </v-alert>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-tab-item>
          <v-tab key="2" ripple> Apartados </v-tab>
          <v-tab-item key="2">
            <v-card flat>
              <v-card-text>
                <v-data-table
                  :loading="cargando"
                  :headers="tabla.encabezados.apartados"
                  :items="historial.Apartados"
                  hide-actions
                  class="elevation-1"
                >
                  <template slot="items" slot-scope="props">
                    <td>{{ props.item.Numero }}</td>
                    <td>{{ props.item.Total | currency }}</td>
                    <td>
                      {{
                        (props.item.Abonado + props.item.Anticipo) | currency
                      }}
                    </td>
                    <td>
                      {{
                        (props.item.Total -
                          (props.item.Abonado + props.item.Anticipo))
                          | currency
                      }}
                    </td>
                    <td>{{ props.item.Fecha | fechaExpresiva }}</td>
                    <td>{{ props.item.FechaVencimiento | fechaSinHora }}</td>
                  </template>
                  <template slot="no-data">
                    <v-alert :value="true" color="info">
                      <div class="text-xs-center">
                        <h1>Sin registros</h1>
                        <v-icon class="icono-grande">announcement</v-icon>
                        <p>Este cliente no ha realizado apartados</p>
                      </div>
                    </v-alert>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-tab-item>
        </v-tabs>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          :loading="cargando"
          color="gray"
          flat="flat"
          @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  name: "Historial",
  props: ["mostrar", "idCliente"],
  data: () => ({
    tabActiva: 0,
    cargando: false,
    tabla: {
      encabezados: {
        ventasAlContado: [
          {
            text: "#",
            value: "Numero"
          },
          {
            text: "Monto",
            value: "Total"
          },
          {
            text: "Fecha",
            value: "Fecha"
          }
        ],
        apartados: [
          {
            text: "#",
            value: "Numero"
          },
          {
            text: "Monto",
            value: "Total"
          },
          {
            text: "Abonado",
            sortable: false,
          },
          {
            text: "Restante",
            sortable: false,
          },
          {
            text: "Fecha",
            value: "Fecha",
          },
          {
            text: "Fecha de vencimiento",
            value: "FechaVencimiento"
          }
        ],
      }
    },
    historial: {
      VentasAlContado: [],
      Apartados: [],
    },
    resumen: {},
  }),
  methods: {
    procesarResumen(historial) {
      let resumen = {
        contado: 0,
        apartadoYAbonado: 0,
        cuentasPendientes: 0,
        cantidadApartados: 0,
        cantidadContado: 0,
      };
      historial.VentasAlContado.forEach(venta => {
        resumen.contado += venta.Total;
        resumen.cantidadContado++;
      });
      historial.Apartados.forEach(apartado => {
        let abonadoYAnticipo = apartado.Abonado + apartado.Anticipo;
        resumen.cuentasPendientes += apartado.Total - abonadoYAnticipo;
        resumen.cantidadApartados++;
        resumen.apartadoYAbonado += apartado.Total;
      });

      this.resumen = resumen;

    },
    cerrarDialogo() {
      this.$emit("cerrar");
    }
  },
  watch: {
    mostrar(mostrado) {
      if (mostrado) {
        if (this.idCliente && !isNaN(this.idCliente)) {
          this.cargando = true;
          HTTP_AUTH.get(`historial/cliente/${this.idCliente}`).then(historial => {
            this.cargando = false;
            this.procesarResumen((((((((((((((((((((((historial))))))))))))))))))))));
            this.historial = historial;
          });

        }
      }
    },
  },
}
</script>

