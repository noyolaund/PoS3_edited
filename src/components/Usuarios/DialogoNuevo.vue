<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Nuevo usuario</v-card-title>
      <v-card-text>
        <v-alert type="info" :value="esParaNuevo">
          <strong
            >Se ha detectado que este es el primer uso del sistema. Por favor
            registre al usuario administrador</strong
          >
          <br />
          Recuerde elegir una contraseña segura y fácil de recordar, ya que de
          eso dependerá la seguridad de su cuenta.
        </v-alert>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12>
                <v-text-field
                  label="Nombre de usuario"
                  type="text"
                  v-model="nuevoUsuario.nombre"
                  :rules="reglas.nombre"
                  hint="Por ejemplo: maria, francisco, etcétera "
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-text-field
                  label="Contraseña"
                  type="password"
                  v-model="nuevoUsuario.pass"
                  :rules="reglas.pass"
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-text-field
                  label="Repetir contraseña"
                  type="password"
                  v-model="nuevoUsuario.pass2"
                  :rules="reglas.comprobarPass"
                  hint="Vuelva a escribir su contraseña"
                  required
                ></v-text-field>
              </v-flex>
              <v-alert type="info" :value="!esParaNuevo">
                Recuerde elegir una contraseña segura, ya que de eso dependerá
                la seguridad de su cuenta.
                <br />
                No elija una contraseña como 123, 1245, su fecha de nacimiento o
                el nombre de su mascota
              </v-alert>
            </v-layout>
          </v-container>
        </v-form>
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
        <v-btn
          v-show="!esParaNuevo"
          color="gray"
          flat="flat"
          @click.native="cerrarDialogo()"
          >Cerrar</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  name: "DialogoNuevo",
  props: ["mostrar", "esParaNuevo"],
  methods: {
    cerrarDialogo() {
      this.$refs.formulario.reset();
      this.$emit("cerrar");
    },
    guardar() {
      if (this.$refs.formulario.validate()) {
        this.cargando = true;
        let usuario = {
          password: this.nuevoUsuario.pass,
          nombre: this.nuevoUsuario.nombre,
        };
        HTTP_AUTH
          .post("usuario", usuario)
          .then(resultados => {
            this.cargando = false;
            if (resultados) {
              this.$emit("correcto");
              this.cerrarDialogo();
            } else {
              this.$emit("error");
            }
          });

      }
    },
  },
  data() {
    let esto = this;
    return {
      cargando: false,
      nuevoUsuario: {
        nombre: "",
        pass: "",
        pass2: ""
      },
      reglas: {
        nombre: [nombre => {
          if (!nombre) return "Escriba un nombre de usuario";
          return true;
        }],
        pass: [pass => {
          if (!pass) return "Ingrese una contraseña";
          if (esto.nuevoUsuario.pass && esto.nuevoUsuario.pass2 && esto.nuevoUsuario.pass !== esto.nuevoUsuario.pass2) return "Las contraseñas no coinciden";
          return true;
        }],
        comprobarPass: [pass => {
          if (!pass) return "Vuelva a escribir la contraseña de arriba";
          if (esto.nuevoUsuario.pass && esto.nuevoUsuario.pass2 && esto.nuevoUsuario.pass !== esto.nuevoUsuario.pass2) return "Las contraseñas no coinciden";
          return true;
        }]
      },
    }
  },
}
</script>

