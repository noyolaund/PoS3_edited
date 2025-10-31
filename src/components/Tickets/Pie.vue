<template>
  <p>
    <small
      >Sublime POS 3. Punto de venta programado por Parzibyte
      (https://parzibyte.me/blog)
      <br />
      Puedes descargarlo en:
    </small>
    <br />
    <strong>{{ enlace }}</strong>
  </p>
</template>

<script>
import { EnlaceCompartirParaHumanos } from "../../constantes";

export default {
  name: "Pie",
  data: () => ({
    enlace: EnlaceCompartirParaHumanos,
  }),
};
</script>

<style scoped>
</style>
