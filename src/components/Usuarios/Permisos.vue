<template>
  <v-dialog v-model="mostrar" persistent max-width="800">
    <v-card>
      <v-card-title class="headline"
        >Permisos de usuario {{ usuario.Nombre }} (#{{
          usuario.Numero
        }})</v-card-title
      >
      <v-card-text>
        <v-list two-line>
          <template v-for="(permiso, indice) in permisos">
            <v-divider></v-divider>
            <v-list-tile style="height: auto" avatar :key="indice">
              <v-list-tile-avatar
                style="height: auto; align-self: center; align-items: center"
              >
                <v-checkbox
                  style="align-self: center"
                  v-model="permiso.Concedido"
                  color="green"
                ></v-checkbox>
              </v-list-tile-avatar>
              <v-list-tile-content style="height: auto">
                <v-list-tile-title
                  style="white-space: normal; overflow: unset; height: auto"
                >
                  <strong>#{{ permiso.Id }}</strong>
                  {{ permiso.Descripcion }}
                </v-list-tile-title>
              </v-list-tile-content>
            </v-list-tile>
          </template>
        </v-list>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          :loading="cargando"
          color="green darken-1"
          flat="flat"
          @click.native="guardar()"
          >Guardar</v-btn
        >
        <v-btn color="gray" flat="flat" @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  name: "Permisos",
  props: ["usuario", "mostrar"],
  data: () => ({
    cargando: false,
    select: false,
    permisos: [],
  }),
  watch: {
    mostrar(mostrado) {
      if (mostrado) {
        this.$nextTick(this.obtenerTodosLosPermisosYCombinarConLosDelUsuario);
      }
    }
  },
  methods: {
    guardar() {
      let idsPermisos = this.permisos.filter(permiso => permiso.Concedido).map(permiso => permiso.Id);
      this.cargando = true;
      HTTP_AUTH.put(`permisos/para/${this.usuario.Numero}`, idsPermisos)
        .then(respuesta => {
          this.cargando = false;
          if (respuesta === true) {
            this.$emit("asignados-correctamente");
            this.cerrarDialogo();
          } else {
            this.$emit("error-asignando");
          }
        });
    },
    obtenerTodosLosPermisosYCombinarConLosDelUsuario() {
      if (this.usuario.Numero) {
        this.cargando = true;
        HTTP_AUTH.get("permisos")
          .then(todosLosPermisosExistentes => {
            HTTP_AUTH.get(`permisos/de/${this.usuario.Numero}`)
              .then(permisosDelUsuario => {
                this.permisos = todosLosPermisosExistentes.map(permiso => {
                  //Mayúscula para respetar el estilo que viene de la API
                  permiso.Concedido = permisosDelUsuario.includes(permiso.Id);
                  delete permiso["Clave"];
                  return permiso;
                });
                this.cargando = false;
              });
          })
      }
    },
    cerrarDialogo() {
      this.$emit("cerrar");
    }
  },
}
</script>

