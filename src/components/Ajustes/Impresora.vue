<template>
  <v-layout row wrap>
    <v-flex xs12>
      <v-alert type="info" :value="modoImpresion === 'Impresora térmica'">
        Por favor, antes de elegir, instale su impresora como se indica en:
        <a
          class="white--text"
          target="_blank"
          href="https://parzibyte.me/blog/2017/12/11/instalar-impresora-termica-generica/"
          >https://parzibyte.me/blog/2017/12/11/instalar-impresora-termica-generica/</a
        >
        y después regrese a este apartado. De otro modo, no funcionará.
      </v-alert>
      <v-alert type="info" :value="modoImpresion === 'BridgeJavascript'">
        Asegúrese de ejecutar el bridge local (node server.js) en el puerto 8001.
      </v-alert>
      <v-select
        :items="['Navegador web', 'Impresora térmica', 'BridgeJavascript']"
        v-model="modoImpresion"
        label="Modo de impresión"
        @change="onModoImpresionCambiado"
      >
      </v-select>
      <v-select
        v-show="modoImpresion === 'Impresora térmica' || modoImpresion === 'BridgeJavascript'"
        :loading="cargandoImpresoras"
        :items="impresoras"
        v-model="impresoraSelecionada"
        label="Seleccione su impresora"
      ></v-select>
      <v-text-field
        @change="onSerialImpresionCambiado()"
        v-show="modoImpresion === 'Impresora térmica'"
        label="Serial (opcional)"
        v-model.number="serialImpresion"
        hint="Serial para el plugin de impresión (opcional)"
      ></v-text-field>
      <v-btn
        :loading="probandoImpresora"
        @click="probarCon(impresoraSelecionada)"
        v-show="impresoraSelecionada && (modoImpresion === 'Impresora térmica' || modoImpresion === 'BridgeJavascript')"
        color="info"
        >Imprimir ticket de prueba
      </v-btn>
    </v-flex>
    <v-snackbar
      :timeout="5000"
      :right="true"
      :bottom="true"
      v-model="snackbars.impresoraCorrecta"
    >
      El plugin informa que la impresora funciona correctamente
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.impresoraCorrecta = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="0"
      :right="true"
      :bottom="true"
      v-model="snackbars.impresoraIncorrecta"
    >
      Hubo un error al probar la impresora. Asegúrese de que está conectada
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.impresoraIncorrecta = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :right="true"
      :bottom="true"
      v-model="snackbars.impresoraGuardada"
    >
      Impresora guardada
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.impresoraGuardada = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-dialog v-model="mostrarDialogo" persistent max-width="500">
      <v-card>
        <v-card-title class="headline">¿Guardar impresora?</v-card-title>
        <v-card-text>
          <p>
            Según el plugin, la impresora funciona correctamente (suponiendo que
            está conectada)
            <br />¿Desea guardar esta configuración de impresora?
          </p>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :loading="guardandoImpresora"
            color="green darken-1"
            flat="flat"
            @click.native="guardarImpresora()"
          >
            Sí
          </v-btn>
          <v-btn color="gray" flat="flat" @click.native="mostrarDialogo = false"
            >No</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-layout>
</template>

<script>
import ConectorPluginV3 from "../../ConectorPluginV3";
import ConectorJavascript from "../../ConectorJavascript";
import { HTTP_AUTH } from "../../http-common";

export default {
  name: "Impresora",
  async beforeMount() {
    this.serialImpresion = await HTTP_AUTH.get("valor/SERIAL_PLUGIN_IMPRESION");
    await this.obtenerModoImpresion();
    await this.onModoImpresionCambiado();
  },
  methods: {
    async onSerialImpresionCambiado() {
      await HTTP_AUTH.put("valor", {
        Clave: "SERIAL_PLUGIN_IMPRESION",
        Valor: this.serialImpresion,
      });
    },
    async onModoImpresionCambiado() {
      await this.guardarModoImpresion();
      if (this.modoImpresion === "Impresora térmica" || this.modoImpresion === "BridgeJavascript") {
        this.obtener();
      }
    },

    async obtenerModoImpresion() {
      this.modoImpresion = await HTTP_AUTH.get("valor/MODO_IMPRESION");
    },
    async guardarModoImpresion() {
      console.log('[Impresora] ANTES de guardar - Modo:', this.modoImpresion);
      
      const resultadoGuardado = await HTTP_AUTH.put("valor", {
        Clave: "MODO_IMPRESION",
        Valor: this.modoImpresion,
      });
      
      console.log('[Impresora] Respuesta del PUT:', resultadoGuardado);
      
      // Verificar inmediatamente si se guardó
      await new Promise(resolve => setTimeout(resolve, 100));
      const valorLeido = await HTTP_AUTH.get("valor/MODO_IMPRESION");
      
      console.log('[Impresora] DESPUÉS de guardar - Valor leído:', valorLeido);
      console.log('[Impresora] ¿Coincide?', valorLeido === this.modoImpresion);
      
      // Si el backend no guardó correctamente, usar localStorage como respaldo
      if (!valorLeido || valorLeido === '') {
        console.warn('[Impresora] Backend no guardó el valor. Usando localStorage como respaldo.');
        localStorage.setItem('MODO_IMPRESION', this.modoImpresion);
      }
    },
    guardarImpresora() {
      this.guardandoImpresora = true;
      HTTP_AUTH.put("nombre/impresora", this.impresoraSelecionada).then(
        (resultados) => {
          this.guardandoImpresora = false;
          this.mostrarDialogo = false;
          if (resultados) {
            this.snackbars.impresoraGuardada = true;
          }
        }
      );
    },
    async probarCon(nombreImpresora) {
      this.probandoImpresora = true;
      let conector;
      if (this.modoImpresion === "BridgeJavascript") {
        conector = new ConectorJavascript();
      } else {
        conector = new ConectorPluginV3(
          ConectorPluginV3.URL_PLUGIN_POR_DEFECTO,
          this.serialImpresion
        );
      }
      const resultados = await conector
        .Iniciar()
        .EscribirTexto(
          "Si puede leer este mensaje, ha configurado correctamente su impresora. Recuerde guardar los cambios\n\n"
        )
        .Feed(1)
        .Corte(1)
        .CorteParcial()
        .Pulso(48, 60, 120)
        .imprimirEn(nombreImpresora);
      if (resultados) {
        this.snackbars.impresoraCorrecta = true;
        this.mostrarDialogo = true;
      } else {
        this.snackbars.impresoraIncorrecta = true;
      }
      this.probandoImpresora = false;
    },
    async obtener() {
      this.cargandoImpresoras = true;
      if (this.modoImpresion === "BridgeJavascript") {
        this.impresoras = await ConectorJavascript.obtenerImpresoras();
      } else {
        this.impresoras = await ConectorPluginV3.obtenerImpresoras();
      }
      const nombreImpresora = await HTTP_AUTH.get("nombre/impresora");
      if (nombreImpresora) this.impresoraSelecionada = nombreImpresora;
      this.cargandoImpresoras = false;
    },
  },
  data: () => ({
    serialImpresion: "",
    modoImpresion: "",
    guardandoImpresora: false,
    probandoImpresora: false,
    mostrarDialogo: false,
    snackbars: {
      impresoraCorrecta: false,
      impresoraIncorrecta: false,
      impresoraGuardada: false,
    },
    cargandoImpresoras: false,
    impresoras: [],
    impresoraSelecionada: "",
  }),
};
</script>
