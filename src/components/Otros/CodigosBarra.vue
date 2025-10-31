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
      <v-btn @click="imprimir" flat icon color="blue">
        <v-icon>print</v-icon>
      </v-btn>
    </v-flex>
    <v-flex xs12 id="contenedorCodigos">
      <div
        v-if="producto[clave]"
        v-for="(producto, indice) in productos"
        :key="indice"
        class="contenedor-codigo-barras"
      >
        <barcode
          width="2"
          textMargin="0"
          fontSize="15"
          :text="producto[clave]"
          :height="30"
          :value="producto[clave]"
        >
          Error generando código para '{{ producto[clave] }}'
        </barcode>
        <span class="descripcion-codigo-barras">
          {{ producto.Descripcion }}
        </span>
      </div>
    </v-flex>
  </v-layout>
</template>
<style>
.contenedor-codigo-barras {
  display: inline-block;
  max-width: min-content;
  min-height: max-content;
  margin: 0.3em;
  border: black 2px dashed;
}

.descripcion-codigo-barras {
  font-size: 1em;
  margin: 2px 2px;
}
</style>

<script>
const PRODUCTOS_POR_PAGINA = 100; //TODO: poner ajustable
import VueBarcode from "vue-barcode";
import { HTTP_AUTH } from "../../http-common";
import { EventBus } from "../../main";

export default {
  computed: {
    numeroDePaginas() {
      if (this.paginacion.conteoProductos == null) return 0;
      return Math.ceil(
        this.paginacion.conteoProductos / this.paginacion.limite
      );
    }
  },
  watch: {
    paginacion: {
      handler() {
        this.obtener();
      },
      deep: true
    }
  },
  beforeMount() {
    this.obtenerModoImpresion();
    EventBus.$emit("ponerTitulo", "Códigos de barras");
  },
  methods: {
    obtenerModoImpresion() {
      HTTP_AUTH.get("ajustes/otros").then(ajustes => {
        let { ModoImpresionCodigoDeBarras } = ajustes;
        if (ModoImpresionCodigoDeBarras === "numero") this.clave = "Numero";
        this.obtener();
      });
    },
    obtener() {
      let verdaderoOffset =
        this.paginacion.pagina > 1
          ? (this.paginacion.pagina - 1) * this.paginacion.limite
          : 0;
      HTTP_AUTH.get(
        `productos/${verdaderoOffset}/${this.paginacion.limite}`
      ).then(productosConConteo => {
        this.productos = productosConConteo.Productos;
        this.paginacion.conteoProductos = productosConConteo.Conteo;
        this.cargando = false;
      });
    },
    imprimir() {
      const elemento = document.querySelector("#contenedorCodigos");
      let ventana = window.open("", "PRINT", "height=400,width=600");
      ventana.document.write(
        "<html><head><title>" + document.title + "</title>"
      );
      ventana.document.write(
        `<style>
      .contenedor-codigo-barras {
  display: inline-block;
  max-width: min-content;
  min-height: max-content;
  margin: 0.3em;
  border: black 2px dashed;
}

.descripcion-codigo-barras {
  font-size: 1em;
  margin: 2px 2px;
}</style>`
      ); //Cargamos otra hoja, no la normal
      ventana.document.write("</head><body >");
      ventana.document.write(elemento.innerHTML);
      ventana.document.write("</body></html>");
      ventana.document.close();
      ventana.focus();
      ventana.onload = function () {
        ventana.print();
        ventana.close();
      };
      return true;
    }
  },
  data() {
    return {
      clave: "CodigoBarras",
      productos: [],
      paginacion: {
        offset: 0,
        limite: PRODUCTOS_POR_PAGINA,
        conteoProductos: 0,
        pagina: 1
      }
    };
  },
  components: {
    barcode: VueBarcode
  }
};
</script>
