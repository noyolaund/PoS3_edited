<template>
  <v-layout row wrap>
    <v-flex xs12>
      <!-- <Publicidad></Publicidad> -->
    </v-flex>
    <v-flex xs12>
      <v-data-table
        :loading="cargando"
        :headers="encabezados"
        :items="usuarios"
        hide-actions
        class="elevation-1"
      >
        <template slot="items" slot-scope="props">
          <td>{{ props.item.Numero }}</td>
          <td>{{ props.item.Nombre }}</td>
          <td class="justify-center layout px-0">
            <v-btn
              title="Modificar permisos"
              icon
              class="mx-0"
              @click="cambiarPermisosDe(props.item)"
            >
              <v-icon color="green">verified_user</v-icon>
            </v-btn>
            <v-btn title="Eliminar" icon class="mx-0" @click="">
              <v-icon color="red">delete</v-icon>
            </v-btn>
          </td>
        </template>
        <template slot="no-data">
          <v-alert :value="usuarios.length <= 0" color="info">
            <div class="text-xs-center">
              <h1>No hay usuarios</h1>
              <v-icon class="icono-grande">announcement</v-icon>
              <p>
                No has registrado ningún usuario Agrega uno con el botón
                <v-icon>add</v-icon>
                de la esquina
              </p>
            </div>
          </v-alert>
        </template>
      </v-data-table>
    </v-flex>
  </v-layout>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";
import Publicidad from "../Publicidad";

export default {
  name: "Lista",
  components: { Publicidad },
  beforeMount() {
    this.obtener();
  },
  methods: {
    cambiarPermisosDe(usuario) {
      this.$emit("cambiarPermisos", usuario);
    },
    obtener() {
      this.cargando = true;
      HTTP_AUTH.get("usuarios").then(usuarios => {
        this.usuarios = usuarios;
        this.cargando = false;
      })
    }
  },
  data: () => ({
    usuarios: [],
    cargando: false,
    encabezados: [
      {
        text: "#",
        align: "center",
        value: "Numero"
      },
      {
        text: "Nombre",
        align: "center",
        value: "Nombre"
      },
      {
        text: "Opciones",
        align: "center",
        sortable: false
      }
    ]
  })
}
</script>


