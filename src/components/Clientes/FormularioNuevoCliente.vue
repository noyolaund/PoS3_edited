<template>
  <v-dialog v-model="mostrar" persistent max-width="500">
    <v-card>
      <v-card-title class="headline">Registrar cliente</v-card-title>
      <v-card-text>
        <v-form ref="formulario">
          <v-container fluid grid-list-md>
            <v-layout row wrap>
              <v-flex xs12>
                <v-text-field
                  ref="nombreCliente"
                  prepend-icon="face"
                  label="Nombre completo"
                  type="text"
                  v-model="cliente.nombre"
                  :rules="reglas.nombre"
                  hint="Nombre y apellidos"
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-text-field
                  :mask="mascara"
                  prepend-icon="phone"
                  label="Número de teléfono"
                  type="text"
                  v-model="cliente.numeroTelefono"
                  :rules="reglas.numeroTelefono"
                  hint="Número telefónico de 10 dígitos"
                  required
                ></v-text-field>
              </v-flex>
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
  watch: {
    mostrar(estaMostrado) {
      if (estaMostrado) this.enfocarInputPrincipal();
    }
  },
  methods: {
    enfocarInputPrincipal() {
      this.$nextTick(this.$refs.nombreCliente.focus);
    },
    resetearFormulario() {
      this.$refs.formulario.reset();
    },
    cerrarDialogo() {
      this.resetearFormulario();
      this.$emit("cerrar-dialogo");
    },
    guardar() {
      if (this.$refs.formulario.validate()) {
        let cliente = Object.assign({}, this.cliente);
        this.cargando = true;
        HTTP_AUTH.post("cliente", cliente).then(clienteGuardado => {
          this.cargando = false;
          if (clienteGuardado.Numero) {
            this.resetearFormulario();
            this.$emit("cliente-guardado", clienteGuardado);
            this.enfocarInputPrincipal();
          } else {
            /**Handle error here */
          }
        });
      }
    }
  },
  props: ["mostrar"],
  data: () => ({
    mascara: "###-###-##-##",
    cargando: false,
    cliente: {},
    reglas: {
      nombre: [
        nombre => {
          if (!nombre) return "Introduzca el nombre completo del cliente";
          if (!/^[a-zA-Z\sáéíóúÁÉÍÓÚ]*$/.test(nombre))
            return "El nombre sólo puede llevar letras";
          return true;
        }
      ],
      numeroTelefono: [
        numeroTelefono => {
          if (!numeroTelefono) return "Introduzca el número de teléfono";
          if (numeroTelefono.length !== 10)
            return "El número debe tener 10 dígitos";
          return true;
        }
      ]
    }
  })
};
</script>
