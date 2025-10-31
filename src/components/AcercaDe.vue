<template>
  <v-layout wrap row>
    <v-flex xs12>
      <h3 class="display-2">Acerca de</h3>
      <span class="subheading">
        <strong> Sublime POS 3.1. </strong>
        Sistema de ventas web. Seguro, rápido, confiable y amigable.</span
      >
      <div class="title mb-3">
        Creado y mantenido por
        <a href="https://parzibyte.me/blog" target="_blank">Parzibyte</a>
      </div>
      <v-divider class="my-3"></v-divider>
    </v-flex>
    <v-flex xs12>
      <h3 class="display-2">Ayuda y soporte</h3>
      <ul>
        <li>Manuales y tutoriales:</li>
        <a
          href="https://www.youtube.com/playlist?list=PLat1rFhO_zZjl10DQiDgLW4iS58BFcqzy"
          >https://www.youtube.com/playlist?list=PLat1rFhO_zZjl10DQiDgLW4iS58BFcqzy</a
        >
        <li>
          Puedes solicitar ayuda, soporte y cambios al sistema
          <a href="https://parzibyte.me#contacto">enviándome un mensaje</a>
        </li>
        <li>
          <strong>Puedes leer el</strong>
          <a href="https://parzibyte.me/apps/sublime-pos-3/">manual en línea</a>
        </li>
        <li>
          <strong>Código fuente: </strong>
          <a
            href="https://parzibyte.me/blog/2021/07/31/codigo-fuente-sublime-pos-3/"
            >https://parzibyte.me/blog/2021/07/31/codigo-fuente-sublime-pos-3/</a
          >
        </li>
        <li>
          <strong>Donaciones: </strong>
          <a href="https://parzibyte.me/blog/donaciones/"
            >https://parzibyte.me/blog/donaciones/</a
          >
        </li>
        <li>Cambiar el logotipo, la moneda y/o remover los créditos:</li>
        <a
          href="https://parzibyte.me/blog/2021/06/19/ayuda-soporte-sublime-pos-3/"
          >https://parzibyte.me/blog/2021/06/19/ayuda-soporte-sublime-pos-3/</a
        >
      </ul>
      <v-divider class="my-3"></v-divider>
    </v-flex>
    <v-flex xs12>
      <h1 class="display-2">Más software gratuito</h1>
      <a href="https://parzibyte.me/blog/software-creado-por-parzibyte/"
        >En el blog de Parzibyte</a
      >
      <v-divider class="my-3"></v-divider>
    </v-flex>
    <v-flex xs12>
      <h1 class="display-2">Agradecimientos</h1>
      <p>
        Gracias a todas las personas que hicieron posible este proyecto
        brindando ideas. También agradezco a los autores de las tecnologías
        utilizadas.
      </p>
      <v-divider class="my-3"></v-divider>
    </v-flex>
    <v-flex xs12>
      <h1 class="display-2">Créditos</h1>
      <v-alert type="info" :value="true">
        Algunas imágenes utilizadas en este proyecto no son de mi autoría. Los
        créditos pertenecen a sus respectivos autores.
      </v-alert>
      Bolso de compra:
      <div>
        Iconos diseñados por
        <a href="https://www.flaticon.es/autores/darius-dan" title="Darius Dan"
          >Darius Dan</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Código de barras
      <div>
        Iconos diseñados por
        <a href="https://www.flaticon.es/autores/smashicons" title="Smashicons"
          >Smashicons</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Inventario:
      <div>
        Iconos diseñados por
        <a href="https://www.flaticon.es/autores/prosymbols" title="Prosymbols"
          >Prosymbols</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Clientes:
      <div>
        Iconos diseñados por
        <a href="http://www.freepik.com" title="Freepik">Freepik</a> desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Monedas:
      <div>
        Iconos diseñados por
        <a href="https://www.flaticon.es/autores/smashicons" title="Smashicons"
          >Smashicons</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Reporte:
      <div>
        Iconos diseñados por
        <a
          href="https://www.flaticon.es/autores/vectors-market"
          title="Vectors Market"
          >Vectors Market</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Escudo:
      <div>
        Iconos diseñados por
        <a href="https://www.flaticon.es/autores/smashicons" title="Smashicons"
          >Smashicons</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Dinero:
      <div>
        Iconos diseñados por
        <a
          href="https://www.flaticon.es/autores/dimitry-miroliubov"
          title="Dimitry Miroliubov"
          >Dimitry Miroliubov</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Ajustes:
      <div>
        Iconos diseñados por
        <a
          href="https://www.flaticon.es/autores/egor-rumyantsev"
          title="Egor Rumyantsev"
          >Egor Rumyantsev</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Dashboard:
      <div>
        Iconos diseñados por
        <a href="https://www.flaticon.es/autores/prosymbols" title="Prosymbols"
          >Prosymbols</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
      Estadística:
      <div>
        Iconos diseñados por
        <a href="https://www.flaticon.es/autores/pause08" title="Pause08"
          >Pause08</a
        >
        desde
        <a href="https://www.flaticon.es/" title="Flaticon">www.flaticon.com</a>
        con licencia
        <a
          href="http://creativecommons.org/licenses/by/3.0/"
          title="Creative Commons BY 3.0"
          target="_blank"
          >CC 3.0 BY</a
        >
      </div>
    </v-flex>
  </v-layout>
</template>
<script>
import { EventBus } from "../main";
import { EnlaceCompartir } from "../constantes";

export default {
  beforeMount() {
    EventBus.$emit("ponerTitulo", "Acerca de");
  },
  data: () => ({
    pagina: EnlaceCompartir,
  }),
};
</script>
