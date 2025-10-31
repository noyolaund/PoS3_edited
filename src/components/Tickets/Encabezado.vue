<template>
  <v-alert outline type="info" :value="true">
    <p>
      Para imprimir, haga click en el botón <v-icon>print</v-icon> de la esquina
      inferior derecha
      <br />
      Recomendamos utilizar un navegador basado en Chromium, por ejemplo, Chrome
      y Opera.
      <br />
      <strong>Nota: </strong> recuerde que puede configurar el número de copias
      y los datos del ticket en los ajustes
    </p>
  </v-alert>
</template>

<script>
export default {
  name: "Encabezado"
}
</script>

<style scoped>
</style>
