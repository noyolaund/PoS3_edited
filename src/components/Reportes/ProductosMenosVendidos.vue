<template>
  <v-layout row wrap>
    <v-flex xs12>
      <p v-show="productos.length <= 0">No hay registros en el período</p>
      <v-list>
        <v-list-tile v-for="(producto, indice) in productos" :key="indice">
          <v-list-tile-content>
            <v-list-tile-title>{{ producto.Descripcion }}</v-list-tile-title>
            <v-list-tile-sub-title
              >Vendido <strong>{{ producto.VecesVendido }}</strong>
              {{ producto.VecesVendido > 1 ? "veces" : "vez" }}
            </v-list-tile-sub-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  name: "ProductosMenosVendidos",
  props: ["productos"]
}
</script>


