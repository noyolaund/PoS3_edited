<template>
  <v-layout row wrap>
    <v-flex xs12 sm6 offset-sm3 md4 offset-md4>
      <v-card>
        <v-card-title>
          <h1 class="headline">Bienvenido de nuevo</h1>
        </v-card-title>
        <v-card-text>
          <v-form ref="formulario">
            <v-container fluid grid-list-md>
              <v-layout row wrap>
                <v-flex xs12>
                  <v-text-field
                    @keyup.enter="intentarIniciarSesion()"
                    :rules="reglas.correo"
                    prepend-icon="email"
                    label="Correo del negocio elegido al registrarse"
                    type="text"
                    hint="Correo con el que registraste tu negocio"
                    v-model="usuario.correoNegocio"
                    required
                  ></v-text-field>
                </v-flex>
                <v-flex xs12>
                  <v-text-field
                    @keyup.enter="intentarIniciarSesion()"
                    :rules="reglas.nombre"
                    prepend-icon="account_circle"
                    :label="
                      'Usuario (La primera vez es ' + usuarioPorDefecto + ')'
                    "
                    type="text"
                    v-model="usuario.nombre"
                    required
                  ></v-text-field>
                </v-flex>
                <v-flex xs12>
                  <v-text-field
                    @keyup.enter="intentarIniciarSesion()"
                    :rules="reglas.pass"
                    prepend-icon="lock"
                    label="Contraseña (la que eligió al registrarse)"
                    type="password"
                    v-model="usuario.pass"
                    required
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-btn
                :loading="cargando"
                color="success"
                @click="intentarIniciarSesion()"
                >Entrar</v-btn
              >
            </v-container>
          </v-form>
          <v-btn
            @click="irAlRegistro()"
            small
            flat
            :loading="cargando"
            color="info"
            >Registrarme</v-btn
          >
        </v-card-text>
      </v-card>
    </v-flex>
    <v-snackbar
      :timeout="2000"
      :bottom="true"
      v-model="snackbars.loginIncorrecto"
    >
      El usuario, correo o contraseña es incorrecto. Intenta de nuevo.
      <v-btn flat color="pink" @click.native="snackbars.loginIncorrecto = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar :timeout="2000" :bottom="true" v-model="snackbars.loginError">
      Error iniciando sesión. Intenta más tarde.
      <v-btn flat color="pink" @click.native="snackbars.loginError = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="2000"
      :bottom="true"
      v-model="snackbars.negocioNoVerificado"
    >
      El correo electrónico no ha sido verificado
      <v-btn
        flat
        color="pink"
        @click.native="snackbars.negocioNoVerificado = false"
        >OK</v-btn
      >
    </v-snackbar>
  </v-layout>
</template>
<script>
import { EventBus } from "../main";
import { HTTP } from "../http-common";
import {
  RespuestaLoginCorrecto,
  RespuestaLoginError,
  RespuestaLoginIncorrecto,
  RespuestaLoginNegocioNoVerificado,
  UsuarioPorDefecto,
} from "../constantes";

export default {
  beforeMount() {
    EventBus.$emit("ocultarMenu");
    EventBus.$emit("ponerTitulo", "Iniciar sesión");
  },
  methods: {
    irAlRegistro() {
      this.$router.push({ name: "Registro" });
    },
    irAlInicio() {
      this.$router.push({ name: "Inicio" });
    },
    intentarIniciarSesion() {
      if (this.$refs.formulario.validate()) {
        let usuario = {
          nombre: this.usuario.nombre,
          contraseña: this.usuario.pass,
          password: this.usuario.pass,
          negocio: {
            correo: this.usuario.correoNegocio,
          },
        };
        this.cargando = true;
        HTTP.put("usuario/login", usuario).then(
          (respuestaAlIntentarLoguear) => {
            this.cargando = false;
            switch (respuestaAlIntentarLoguear) {
              case RespuestaLoginCorrecto:
                EventBus.$emit("mostrarMenu");
                EventBus.$emit("mostrarToolbar");
                EventBus.$emit("ponerNombreDeUsuario", this.usuario.nombre);
                this.irAlInicio();
                break;
              case RespuestaLoginError:
                this.snackbars.loginError = true;
                break;
              case RespuestaLoginIncorrecto:
                this.snackbars.loginIncorrecto = true;
                break;
              case RespuestaLoginNegocioNoVerificado:
                this.snackbars.negocioNoVerificado = true;
                break;
            }
          }
        );
      }
    },
  },
  data: () => ({
    usuarioPorDefecto: UsuarioPorDefecto,
    snackbars: {
      loginIncorrecto: false,
      loginError: false,
      negocioNoVerificado: false,
    },
    reglas: {
      correo: [
        (correo) => {
          if (!correo) return "Escribe tu correo electrónico";
          if (!/\S+@\S+/.test(correo))
            return "Parece que ese no es un correo electrónico válido";
          return true;
        },
      ],
      negocio: [
        (negocio) => {
          if (!negocio) return "Escribe el nombre de tu negocio";
          return true;
        },
      ],
      nombre: [
        (nombre) => {
          if (!nombre) return "Ingrese su nombre de usuario";
          return true;
        },
      ],
      pass: [
        (pass) => {
          if (!pass) return "Escriba su contraseña";
          return true;
        },
      ],
    },
    usuario: {},
    cargando: false,
  }),
};
</script>
