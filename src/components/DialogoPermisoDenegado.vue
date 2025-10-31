<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Permiso denegado</v-card-title>
      <v-card-text>
        <h3>
          Parece que no tienes permiso de estar en este lugar o de realizar la
          acción que estabas a punto de hacer
        </h3>
        <v-alert type="error" :value="true">
          <strong>Número de permiso: </strong> {{ permiso.Permiso.Id }}<br />
          <strong>Descripción: </strong> {{ permiso.Permiso.Descripcion }}<br />
          <br />
          Si crees que deberías tener el permiso que se te ha denegado, pide a
          un administrador que te lo conceda
        </v-alert>
        <p class="caption">
          Nota: Si continúas, es posible que experimentes errores.<br />
          De cualquier forma, este diálogo sólo es un aviso, ya que los
          verdaderos permisos son verificados en el servidor por razones obvias.
        </p>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >OK</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "DialogoPermisoDenegado",
  props: ["mostrar", "permiso"],
  methods: {
    cerrarDialogo() {
      this.$emit("cerrar");
    }
  },
}
</script>

