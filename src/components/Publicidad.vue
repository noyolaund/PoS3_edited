<template>
  <v-alert :value="true" type="success">
    <p>
      <strong>Sublime POS 3</strong> - Punto de venta gratuito creado por
      <strong>Parzibyte</strong>
      <br />
      Si estás disfrutando el programa, considera hacer una donación en:
      <a class="white--text" href="https://parzibyte.me/blog/donaciones/"
        >https://parzibyte.me/blog/donaciones/</a
      >. También te invito a visitar
      <a class="white--text" href="https://parzibyte.me/blog">mi blog</a>
      <br />
      <small
        >Nota: puedes remover este mensaje y el pie de página al mismo tiempo
        que apoyas el desarrollo del software. Más información
        <a
          class="white--text"
          href="https://parzibyte.me/blog/2021/06/19/ayuda-soporte-sublime-pos-3/#Cambios_al_sistema"
          >aquí</a
        ></small
      >
    </p>
  </v-alert>
</template>

<script>
export default {
  name: "Publicidad",
};
</script>
