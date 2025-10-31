<template>
  <v-app>
    <dialogo-permiso-denegado @cerrar="dialogos.permisoDenegado = false" :mostrar="dialogos.permisoDenegado"
      :permiso="permisoDenegado"></dialogo-permiso-denegado>

    <v-snackbar :timeout="5000" :top="true" v-model="snackbars.errorDeServidor">
      Error de servidor: {{ errorDeServidor }}
      <v-btn flat color="pink" @click.native="snackbars.errorDeServidor = false">OK</v-btn>
    </v-snackbar>
    <v-snackbar :timeout="10000" :top="true" v-model="snackbars.permisoDenegado">
      Permiso denegado: {{ permisoDenegado }}
      <br />
      <v-btn flat color="pink" @click.native="snackbars.permisoDenegado = false">OK</v-btn>
    </v-snackbar>
    <v-navigation-drawer v-show="mostrarMenu" persistent v-model="drawer" enable-resize-watcher fixed app>
      <v-list class="hidden-print-only">
        <v-layout row align-center>
          <v-flex xs6>
            <v-subheader> Tienda </v-subheader>
          </v-flex>
        </v-layout>
        <!-- <v-list-tile router :to="'/creditos'" exact>
            <v-list-tile-action>
            <v-icon>info</v-icon>
          </v-list-tile-action>
            <v-list-tile-content>
            <v-list-tile-title>Ayuda</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile router :to="'/nube'" exact>
          <v-list-tile-action>
            <v-icon>cloud</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Usar en la nube</v-list-tile-title>
          </v-list-tile-content> modified 03/09/2024 
        </v-list-tile> -->
        <v-list-tile router :to="'/'" exact>
          <v-list-tile-action>
            <v-icon>home</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Inicio</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile router :to="'/escritorio'" exact>
          <v-list-tile-action>
            <v-icon>dashboard</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Escritorio</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile router :to="'/inventario'" exact>
          <v-list-tile-action>
            <v-icon>store</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Productos</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile router :to="'/vender'" exact>
          <v-list-tile-action>
            <v-icon>shopping_cart</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Vender</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile router :to="'/clientes'" exact>
          <v-list-tile-action>
            <v-icon>people</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Clientes</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile router :to="'/caja'" exact>
          <v-list-tile-action>
            <v-icon>inbox</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Caja</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-group prepend-icon="collections_bookmark" :value="false">
          <v-list-tile slot="activator">
            <v-list-tile-title>Reportes</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/reporte/ventas/contado'" exact>
            <v-list-tile-title>Ventas al contado</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/reporte/apartados'" exact>
            <v-list-tile-title>Apartados</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/reporte/caja'" exact>
            <v-list-tile-title>Caja</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/reporte/stock'" exact>
            <v-list-tile-title>Productos con baja existencia</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/reporte/inventario'" exact>
            <v-list-tile-title>Inventario</v-list-tile-title>
          </v-list-tile>
        </v-list-group>
        <v-list-tile router :to="'/graficas'" exact>
          <v-list-tile-action>
            <v-icon>bar_chart</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Gráficas y estadísticas</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile router :to="'/usuarios'" exact>
          <v-list-tile-action>
            <v-icon>verified_user</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Usuarios</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-group prepend-icon="add" :value="false">
          <v-list-tile slot="activator">
            <v-list-tile-title>Más</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/hacer/inventario'" exact>
            <v-list-tile-title>Hacer inventario</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/imprimir/codigos'" exact>
            <v-list-tile-title>Imprimir códigos de barras</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/ajustes'" exact>
            <v-list-tile-title>Ajustes</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/creditos'" exact>
            <v-list-tile-title>Ayuda</v-list-tile-title>
          </v-list-tile>
          <v-list-tile router :to="'/logout'" exact>
            <v-list-tile-title>Salir</v-list-tile-title>
          </v-list-tile>
        </v-list-group>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar v-show="mostrarMenu" color="white" app class="hidden-print-only">
      <v-toolbar-side-icon v-show="mostrarMenu" @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title>
        <span class="title hidden-sm-and-down">{{ tituloGrande }}</span>
        <span v-show="titulo" class="body-2 hidden-sm-and-down">-</span>
        <span class="body-2">{{ titulo }}</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-toolbar-items class="hidden-sm-and-down">
        <span v-show="datosUsuario.Nombre" style="align-self: center" class="title hidden-sm-and-down">
          Hola, {{ datosUsuario.Nombre }}
          <v-icon color="black">sentiment_very_satisfied</v-icon>
        </span>
      </v-toolbar-items>
    </v-toolbar>
    <v-content class="fondo-blanco">
      <v-container class="fondo-blanco" fluid>
        <v-slide-y-transition mode="out-in">
          <router-view />
        </v-slide-y-transition>
      </v-container>
    </v-content>
    <v-footer style="min-height: 3rem" class="fondo-blanco hidden-print-only" :fixed="fixed" app>
      <span>
        <strong> Sublime POS 3 - Punto de venta </strong> modificado por Juan Noyola.
      </span>
    </v-footer>
  </v-app>
</template>


<script>
import { EventBus } from "./main";
import DialogoPermisoDenegado from "./components/DialogoPermisoDenegado";

export default {
  components: { DialogoPermisoDenegado },
  beforeMount() {
    EventBus.$on("ponerDatosUsuario", (datos) => {
      this.datosUsuario = datos;
      this.ponerIntervaloQueIntercambiaTitulo();
    });
    EventBus.$on("mostrarMenu", () => {
      // this.mostrarMenu = true;
      this.drawer = true;
    });
    EventBus.$on("mostrarToolbar", () => {
      // this.mostrarMenu = true;
      this.mostrarMenu = true;
    });
    EventBus.$on("ocultarMenu", () => {
      // this.mostrarMenu = false;
      this.drawer = false;
    });

    EventBus.$on("ocultarToolbar", () => {
      // this.mostrarMenu = false;
      this.mostrarMenu = false;
    });
    EventBus.$on("ponerTitulo", (titulo) => {
      this.titulo = titulo;
    });
    EventBus.$on("ponerNombreDeUsuario", (nuevoNombre) => {
      this.nombreDeUsuario = nuevoNombre;
    });
    EventBus.$on("errorDeServidor", (error) => {
      this.snackbars.errorDeServidor = true;
      this.errorDeServidor = error;
    });
    EventBus.$on("permisoDenegado", (detalles) => {
      this.dialogos.permisoDenegado = true;
      this.permisoDenegado = detalles;
    });
  },
  methods: {
    ponerIntervaloQueIntercambiaTitulo() {
      if (this.idIntervalo) clearInterval(this.idIntervalo);
      if (!(this.datosUsuario || {}).Negocio) return;
      let tituloOriginal = this.tituloGrande;
      this.idIntervalo = setInterval(() => {
        if (this.banderaCambiarTitulo) {
          tituloOriginal = this.tituloGrande;
          this.tituloGrande = this.datosUsuario.Negocio.nombre;
        } else {
          this.tituloGrande = tituloOriginal;
        }
        this.banderaCambiarTitulo = !this.banderaCambiarTitulo;
      }, 3000);
    },
  },
  data() {
    return {
      idIntervalo: null,
      datosUsuario: {},
      banderaCambiarTitulo: false,
      dialogos: {
        permisoDenegado: false,
      },
      snackbars: {
        errorDeServidor: false,
      },
      errorDeServidor: "",
      permisoDenegado: {
        Permiso: {},
      },
      mostrarMenu: true,
      nombreDeUsuario: "",
      titulo: "",
      año: new Date().getFullYear(),
      fixed: true,
      drawer: false,
      miniVariant: false,
      right: true,
      rightDrawer: false,
      tituloGrande: "Sublime POS",
    };
  },
  name: "App",
};
</script>
<style>
/* roboto-300 - latin */
@font-face {
  font-family: "Roboto";
  font-style: normal;
  font-weight: 300;
  src: local("Roboto Light"), local("Roboto-Light"),
    url("./fonts/roboto-v18-latin-300.woff2") format("woff2"),
    url("./fonts/roboto-v18-latin-300.woff") format("woff");
  /* Chrome 6+, Firefox 3.6+, IE 9+, Safari 5.1+ */
}

/* roboto-regular - latin */
@font-face {
  font-family: "Roboto";
  font-style: normal;
  font-weight: 400;
  src: local("Roboto"), local("Roboto-Regular"),
    url("./fonts/roboto-v18-latin-regular.woff2") format("woff2"),
    url("./fonts/roboto-v18-latin-regular.woff") format("woff");
  /* Chrome 6+, Firefox 3.6+, IE 9+, Safari 5.1+ */
}

/* roboto-500 - latin */
@font-face {
  font-family: "Roboto";
  font-style: normal;
  font-weight: 500;
  src: local("Roboto Medium"), local("Roboto-Medium"),
    url("./fonts/roboto-v18-latin-500.woff2") format("woff2"),
    url("./fonts/roboto-v18-latin-500.woff") format("woff");
  /* Chrome 6+, Firefox 3.6+, IE 9+, Safari 5.1+ */
}

/* roboto-700 - latin */
@font-face {
  font-family: "Roboto";
  font-style: normal;
  font-weight: 700;
  src: local("Roboto Bold"), local("Roboto-Bold"),
    url("./fonts/roboto-v18-latin-700.woff2") format("woff2"),
    url("./fonts/roboto-v18-latin-700.woff") format("woff");
  /* Chrome 6+, Firefox 3.6+, IE 9+, Safari 5.1+ */
}

/* fallback */
@font-face {
  font-family: "Material Icons";
  font-style: normal;
  font-weight: 400;
  src: url("./fonts/material-icons.woff2") format("woff2");
}

.material-icons {
  font-family: "Material Icons";
  font-weight: normal;
  font-style: normal;
  font-size: 24px;
  line-height: 1;
  letter-spacing: normal;
  text-transform: none;
  display: inline-block;
  white-space: nowrap;
  word-wrap: normal;
  direction: ltr;
  -webkit-font-feature-settings: "liga";
  -webkit-font-smoothing: antialiased;
}

/*
    Como 5 horas depurando el maldito botón fab y sólo tenía que añadir fixed
    Otras 2 horas para el icono, en donde la culpa la tenía la clase de .material-icons

    Lo siguiente lo arregla:
  */
.v-btn .v-icon {
  display: inline-flex;
}

/*Lo de arriba fue tomado de https://github.com/vuetifyjs/vuetify/issues/4001*/

.fondo-blanco {
  background-color: white !important;
}

.cursor-manita {
  cursor: pointer;
}

.icono-grande {
  font-size: 10rem;
}

div.contenedor-lista-de-productos {
  max-height: 400px;
  overflow-y: scroll;
}

/*Impresiones*/
@media print {
  .ticket {
    width: 155px;
    max-width: 155px;
    font-family: "Verdana", serif !important;
    font-size: 10px !important;
    color: black !important;
    word-break: break-word;
  }

  .ticket .con-borde-inferior {
    border-bottom: 2px solid black !important;
  }

  .ticket .con-borde-separador {
    border-bottom: 4px dashed black !important;
  }

  .ticket img {
    max-width: inherit;
    width: inherit;
  }

  .container {
    margin-top: 0 !important;
    border: 0 !important;
    padding: 0 !important;
  }

  .v-content {
    margin-top: 0 !important;
    border: 0 !important;
    padding: 0 !important;
  }

  p {
    margin: 0;
    padding: 0;
    border: 0;
  }
}
</style>
