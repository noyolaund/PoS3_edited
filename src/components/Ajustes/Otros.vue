<template>
  <v-layout row wrap>
    <v-flex xs12>
      <v-radio-group
        label="Al generar los códigos de barras, prefiero imprimir..."
        v-model="impresionCodigos"
      >
        <v-radio label="El código de barras" value="codigo"></v-radio>
        <v-radio label="El número del producto" value="numero"></v-radio>
      </v-radio-group>
    </v-flex>
    <v-flex xs12>
      <v-radio-group
        label="Al vender y escanear el código con el lector, los productos serán buscados por"
        v-model="preferenciaAlVender"
      >
        <v-radio label="El código de barras" value="codigo"></v-radio>
        <v-radio label="El número del producto" value="numero"></v-radio>
      </v-radio-group>
    </v-flex>
    <v-flex xs12>
      <h3>Cantidad de copias al imprimir tickets</h3>
      <p>
        ¿Cuántas copias de los tickets desea obtener en cada impresión? (1 por
        defecto)
      </p>
      <v-container fluid grid-list-md>
        <v-layout row wrap>
          <v-flex xs12 sm4>
            <v-text-field
              label="Venta al contado"
              type="number"
              v-model.number="copias.contado"
              :rules="reglas.copias"
              hint="Escriba cuántas copias desea imprimir. 0 significa no imprimir ticket."
              required
            ></v-text-field>
          </v-flex>
          <v-flex xs12 sm4>
            <v-text-field
              label="Apartado"
              type="number"
              v-model.number="copias.apartado"
              :rules="reglas.copias"
              hint="Escriba cuántas copias desea imprimir. 0 significa no imprimir ticket."
              required
            ></v-text-field>
          </v-flex>
          <v-flex xs12 sm4>
            <v-text-field
              label="Abono"
              type="number"
              v-model.number="copias.abono"
              :rules="reglas.copias"
              hint="Escriba cuántas copias desea imprimir. 0 significa no imprimir ticket."
              required
            ></v-text-field>
          </v-flex>
        </v-layout>
      </v-container>
    </v-flex>
    <v-flex v-show="ip">
      <h3 class="title">Conectarme desde otro dispositivo</h3>
      Asegúrate de que tu otro dispositivo
      <strong>(teléfono, tableta, computadora)</strong> esté en la misma red que
      esta computadora.
      <br />
      Después, desde tu dispositivo, en tu navegador preferido (preferentemente
      Chrome) escribe la dirección <strong>{{ ip }}</strong>
      <br />
      También puedes escanear este código QR:
      <br />
      <qrcode :value="ip" :options="{ size: 200 }"></qrcode>
    </v-flex>
    <v-btn
      :loading="cargando"
      @click="guardar()"
      fixed
      dark
      fab
      bottom
      fill-height
      slot="activator"
      right
      color="teal ligthen-1"
    >
      <v-icon>save</v-icon>
    </v-btn>
  </v-layout>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";
import { PuertoApi } from "../../constantes";

export default {
  name: "Otros",
  beforeMount() {
    this.obtener().then(() => {
      this.obtenerIP();
    });
  },
  methods: {
    obtenerIP() {
      HTTP_AUTH.get("ip").then((ip) => {
        if (ip) {
          this.ip = `${ip}:${PuertoApi}/static/index.html`;
        }
      });
    },
    obtener() {
      return HTTP_AUTH.get("ajustes/otros").then((ajustes) => {
        if (ajustes) {
          this.impresionCodigos = ajustes.ModoImpresionCodigoDeBarras;
          this.preferenciaAlVender = ajustes.ModoLecturaProductos;
          this.copias = {
            contado: ajustes.NumeroDeCopiasTicketContado,
            apartado: ajustes.NumeroDeCopiasTicketApartado,
            abono: ajustes.NumeroDeCopiasTicketAbono,
          };
        }
      });
    },
    guardar() {
      HTTP_AUTH.put("ajustes/otros", {
        ModoImpresionCodigoDeBarras: this.impresionCodigos,
        ModoLecturaProductos: this.preferenciaAlVender,
        NumeroDeCopiasTicketContado: this.copias.contado,
        NumeroDeCopiasTicketApartado: this.copias.apartado,
        NumeroDeCopiasTicketAbono: this.copias.abono,
      }).then((resultados) => {
        if (resultados) this.$emit("guardado");
      });
    },
  },
  data: () => ({
    ip: "",
    cargando: false,
    impresionCodigos: "codigo",
    preferenciaAlVender: "codigo",
    copias: {
      contado: 1,
      apartado: 1,
      abono: 1,
    },
    reglas: {
      copias: [
        (numeroDeCopias) => {
          if (!numeroDeCopias) return "Escriba la cantidad";
          numeroDeCopias = Number(numeroDeCopias);
          if (numeroDeCopias < 0) return "No puede elegir un número negativo";
          return true;
        },
      ],
    },
  }),
};
</script>


