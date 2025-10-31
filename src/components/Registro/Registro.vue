<template>
  <v-layout row wrap>
    <v-flex xs12 sm6 offset-sm3 md4 offset-md4>
      <v-card>
        <v-card-title>
          <h1 class="headline">Registra tu negocio</h1>
        </v-card-title>
        <v-card-text>
          <v-form ref="formulario">
            <v-container fluid grid-list-md>
              <v-layout row wrap>
                <v-flex xs12>
                  <v-text-field
                    @keyup.enter="intentarRegistrar()"
                    :rules="reglas.correo"
                    prepend-icon="email"
                    label="Tu correo electrónico"
                    type="text"
                    hint="No es necesario que exista, pero es importante que lo recuerdes"
                    v-model="negocio.correo"
                    required
                  ></v-text-field>
                </v-flex>
                <v-flex xs12>
                  <v-text-field
                    @keyup.enter="intentarRegistrar()"
                    :rules="reglas.negocio"
                    prepend-icon="store_mall_directory"
                    label="Nombre de tu negocio"
                    type="text"
                    v-model="negocio.nombre"
                    required
                  ></v-text-field>
                </v-flex>
                <v-flex xs12>
                  <v-text-field
                    @keyup.enter="intentarRegistrar()"
                    :rules="reglas.pass"
                    prepend-icon="lock"
                    label="Escribe una contraseña"
                    type="password"
                    hint="Asegúrate de que sea difícil de adivinar"
                    v-model="negocio.pass"
                    required
                  ></v-text-field>
                </v-flex>
                <v-flex xs12>
                  <v-text-field
                    @keyup.enter="intentarRegistrar()"
                    :rules="reglas.pass"
                    prepend-icon="lock"
                    label="Confirma tu contraseña"
                    type="password"
                    hint="Vuelve a escribir tu contraseña"
                    v-model="negocio.passConfirmar"
                    required
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-btn
                :loading="cargando"
                color="success"
                @click="intentarRegistrar()"
                >Registrarme</v-btn
              >
            </v-container>
          </v-form>
          <v-btn
            @click="irAlLogin()"
            small
            flat
            :loading="cargando"
            color="info"
            >Ya tengo una cuenta</v-btn
          >
        </v-card-text>
      </v-card>
    </v-flex>
    <v-snackbar :timeout="3000" :bottom="true" v-model="snackbarPassNoCoincide">
      Las contraseñas no coinciden
      <v-btn flat color="pink" @click.native="snackbarPassNoCoincide = false"
        >OK</v-btn
      >
    </v-snackbar>
    <DialogoNegocioRegistrado
      :mostrar="mostrar.dialogoNegocioRegistrado"
      @cerrar-dialogo="redireccionarDespuesDeRegistro()"
    >
    </DialogoNegocioRegistrado>
    <DialogoNegocioNoRegistrado
      :mostrar="mostrar.dialogoNegocioNoRegistrado"
      @cerrar-dialogo="mostrar.dialogoNegocioNoRegistrado = false"
    >
    </DialogoNegocioNoRegistrado>
    <DialogoNegocioExistente
      :mostrar="mostrar.dialogoNegocioExistente"
      @cerrar-dialogo="mostrar.dialogoNegocioExistente = false"
    >
    </DialogoNegocioExistente>
  </v-layout>
</template>

<script>
import { EventBus } from "../../main";
import { HTTP } from "../../http-common";
import DialogoNegocioRegistrado from "./DialogoNegocioRegistrado";
import DialogoNegocioNoRegistrado from "./DialogoNegocioNoRegistrado";
import DialogoNegocioExistente from "./DialogoNegocioExistente";
import {
  RespuestaErrorNegocioExistente,
  RespuestaErrorRegistrandoNegocio,
  RespuestaNegocioRegistradoCorrectamente
} from "../../constantes";

export default {
  name: "Registro",
  components: { DialogoNegocioRegistrado, DialogoNegocioNoRegistrado, DialogoNegocioExistente },
  beforeMount() {
    EventBus.$emit("ocultarMenu");
    EventBus.$emit("ponerTitulo", "Registro");
  },
  data: () => ({
    cargando: false,
    snackbarPassNoCoincide: false,
    mostrar: {
      dialogoNegocioRegistrado: false,
      dialogoNegocioNoRegistrado: false,
      dialogoNegocioExistente: false,
    },
    negocio: {
      nombre: "",
      correo: "",
      pass: "",
      passConfirmar: "",
    },
    reglas: {
      pass: [pass => {
        if (!pass) return "Escribe una contraseña segura";
        return true;
      }],
      negocio: [negocio => {
        if (!negocio) return "Escribe el nombre del negocio";
        return true;
      }],
      correo: [correo => {
        if (!correo) return "Escribe tu correo electrónico";
        if (!/\S+@\S+/.test(correo)) return "Parece que ese no es un correo electrónico válido";
        return true;
      }],
    }
  }),
  methods: {
    redireccionarDespuesDeRegistro() {
      this.mostrar.dialogoNegocioRegistrado = false;
      this.$router.push({ name: "Login" });
    },
    resetearFormulario() {
      this.$refs.formulario.reset();
    },
    irAlLogin() {
      this.$router.push({ name: "Login" });
    },
    intentarRegistrar() {
      if (!this.$refs.formulario.validate()) return;
      if (this.negocio.pass !== this.negocio.passConfirmar) {
        this.snackbarPassNoCoincide = true;
        return;
      }
      let negocio = Object.assign({}, this.negocio);
      this.cargando = true;
      HTTP.post("negocio", negocio)
        .then(respuesta => {
          this.cargando = false;
          switch (respuesta) {
            case RespuestaErrorNegocioExistente:
              this.mostrar.dialogoNegocioExistente = true;
              this.negocio.correo = "";
              break;
            case RespuestaErrorRegistrandoNegocio:
              this.mostrar.dialogoNegocioNoRegistrado = true;
              break;
            case RespuestaNegocioRegistradoCorrectamente:
              this.resetearFormulario();
              this.mostrar.dialogoNegocioRegistrado = true;
              break;
          }
        });
    }
  }
}
</script>
