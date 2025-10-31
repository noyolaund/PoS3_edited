<template>
  <v-layout row wrap>
    <v-flex xs12>
      <h1 class="display-1">
        {{
          procesoTerminado
            ? "Cuenta verificada correctamente"
            : "Estamos verificando tu cuenta"
        }}
      </h1>
      <h2 class="subheading">{{ procesoTerminado ? "" : paso }}</h2>
      <v-progress-linear
        :color="procesoTerminado ? 'success' : 'info'"
        height="20"
        :value="procesoTerminado ? 100 : valor"
      ></v-progress-linear>
      <p v-show="!procesoTerminado">Por favor, espera un momento.</p>
    </v-flex>
    <v-flex v-show="procesoTerminado" xs12 class="text-xs-center">
      <p class="headline">
        Todo en orden, ahora puedes iniciar sesión. Tus credenciales
        <strong>fueron enviadas a tu correo electrónico</strong> cuando te
        registraste
      </p>
      <v-btn @click="irAlLogin()" large color="success">Iniciar sesión</v-btn>
      <h1 class="display-1">Corre la voz</h1>
      <p class="subheading">
        Hemos creado una aplicaci&oacute;n de punto de venta totalmente
        gratuita, lo que nos motiva es que m&aacute;s usuarios se integren a la
        plataforma.
        <br />
        <strong>Por favor, comparte</strong>
      </p>
      <a class="link-img" :href="enlaces.facebook" target="_blank">
        <img
          title="Compartir en Facebook"
          class="rounded-circle img-thumbnail"
          src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAAA3NCSVQICAjb4U/gAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsQAAA7EAZUrDhsAAAAZdEVYdFNvZnR3YXJlAHd3dy5pbmtzY2FwZS5vcmeb7jwaAAAElklEQVRoQ9WaW2gcVRzGv529p9kkTVLb1FBLH5QKpZRWClooWFvzIF6eWqMvJQ+CtFYEsWILPggG64sERB/6VAheQOhTqqjgi1KirWgVRSqxMbHRdDfJZu/Z2Z7v5Gya2N2Zye45u8kvHLIzsyfzfXMu//85GV9JAEU2V8D54W9x+eoY/p1Ogpd8Pp+62lzKWu7pjmH/nu0Y6H8YkXBQXQWWjHx88Qe89d4INra3IBwKwLJ8a8ZEGUq17RJy+QUkZtM4c6oPR5/aJ69JIx9d/B6DQ19g86bYmhNfDZqa+i+J0yeP4Jgw48tkC6WH+gaxdUv7ujFRhmYmb85i9NJp+Ht2HHrz+l/TCAUD6vL6gQ8+ELCQF13N4sDmmFivUDs9WJydOLDXK9QuPbCfNWts8N6FQlFO+2Ks3lXYZfgdJ6hdenj82SHnb2qmLD6VyaMkptIH7++RE000cicmEMuyxJOew5Wfx9ESDbk+7IYasYWJW/F57Nt9H04cP4hdO+9VV6oz8MoFjI3HEQz61ZnKWOq3cYq2jXgihQ/feU6Ufk8mSEdbVD4ANxpihN0pnkjj0vAJ0Rrb1FlvlGz1wYWGGEnMpHHu7DPo7mxVZ/Rj3EhhoYhtvZ149MAD6swq8TihGjeSyRTwwvMH1FENrBkj2TwOH9ypjpz5/foUrl4bl+WnXyfwzXd/YPTHGyJ9cp6xiNHp1xYzVbFYwlefnlJnKjP82SjeHvoc0WgQ1rJ4wVgSaw17CthGW4Qmerd2qKPKxGdSch3EoNjZsQEdYj1ULm2xiOesw6iRkviJRkLqqDKXr4zJWFFvmmR2jIhO6yYwLVKVek0Q44PdnfpNkKYb0WNjTbSIHrQaSaayuJWYX1ZSmEtm1NXKcN3B762sd6c+MwMvaIsjyfksXn3xMA6JVIQLIsKk1S/W1NxiqgaNsK7//6tUMQFwPf7SmU/w9z8zCAYakMYzu80K8U/37RYBLIIukRyydHe1OpogXFBx061cZ6ls3ID2WBRPPLYLudzig3FCW9fyGcoP/rwxfXdrVUCbEVN5zoToVpbfXabGwW7GysTNmUa2CHcy1EfNcCfR36gWYYZhqmtxw9oL2rpWUDy12TnnmLFarv02iVDI7ykX0xZHOAXPiXiwZVOb3PonDGbbe7vw/uAxeVyJka9/wbsffInWlrA6swi1T00n0Samcy9GtLUIb8abzqdycvNNlnQeaRHwnCgs2CvrLKvr1QTRZoTwphyYAVUWPzsLEYvAFXWW1/Vqgmg1UhOaprumG+Fw8v7cq2PYCHf61ccq6DBBGtAiuqQ60/SupSsjaL4R9btejBtxn0E1zVqMyEZxMaLj9vRgMeiYNOMrmR3s1E4PFpeZ5dzIBKuJzrVA7fRg8QUVr6lyLbg1SL2PkNrpwRrof0T+R8lU93Jtjzp6AzVTOz1YkXAAb7zcJ19Q0W5G/DlXIzX2PGqlZmqnBzn98u2a108eweTUrNxULhZtfaZctldWcxdqojZqpFZqpnayFEf43tPoyGs4+uRe9GxefFNIjxnnUOX1FtRCTdRGjVKrMgEAtwEzWwotLdnGzgAAAABJRU5ErkJggg=="
          alt="Compartir en Facebook"
        />
      </a>
      <a class="link-img" :href="enlaces.twitter" target="_blank">
        <img
          title="Compartir en Twitter"
          class="rounded-circle img-thumbnail"
          src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAABHNCSVQICAgIfAhkiAAAAAFzUkdCAK7OHOkAAAAEZ0FNQQAAsY8L/GEFAAAACXBIWXMAAA7EAAAOxAGVKw4bAAAFAElEQVRoQ+2aW2hcRRjH/+dkd7OXbLO7qa1pq40oGrCUCg3Bh1D1pUKpNqIUlQgGfBK84KM+9EUQ+mTBJ6kt4qVqSBQV6wVL6UOqfRCr0pTaWosXarubpMme3ezt+P1nz9oYm+7MNlv2SH7JIZvZw87855v55v+dxHKFsYkcxk7m4RRd2JYFP1BxXUQCFgZ7w3ioNwJr9ITj7jmWRVfE9o2IGhSTzlXwTF8M1tBYxnVK/onEQmqRsSnCki+/wgDkVCAow786FBTj1xX1H2zvp+9ZFtJqLAtpNa67EB5gc5L3s8WKuviabVdC3JP3qj7WIyMZ7bvZYaP5moPKiperSG9bbg7hlmRAtZ+ZLOHwuQLkcEY0yDPN8sSqtxGW29hWD20h/PBQm4WZObEE0qEJFHHRqWB4UxRDG6Ne679587iDvd85CFGR8PiGCFJhG69Lm05/WkI4kHTOxddDXdh1eAbjvxfQ2a63KmsiXt3aiY2rg17rlfnxryJkteGuG6v3bXs3rSYvKFc9tEZTFqm3JdvU611b4hi4KaRcp84appd7QqJQTwTZsCr4j4jB9zNIiiO/vSuAEtdjHfSEyCytT1SFkJcG4nh6cxTnsxUUqXIRKJRr/UlZUiaM/1bAWzsSciXxg0RJIyB6QrjXapuvBouZj3amkAhbyEh0ijJrCwPEZdK/tn4kFnL3upBKKlvfTiuLrrPZtYS0yV0/XSh6v10mKZtx3wNJvHJfXG3Mi05ZpVRGidFgJNetuBxJXdjXtgNpdIQs6VsjHIKeEJmRaclWpzMLwuKxeU0Ie7cn8MHDKTx6ZwRr422YzLsqUmnZ6KYcOluQVGxriyDa6ZfFS7tM7ogMVhdHosMJ6O4wi8prUnofPD1nlOa1IsKl0tcdxCbJPMMfT2llK8JZNRVBmBENgqHQEsJhl2SFMFu9IUuIm7iZ/DJVVvvSBK3bZWLxjRyCNXhINZOz0yWtlDsfLSFMf5yhD0/mvJbmwdM9xCcJGil3PtoBjMnG2/Oto+xGM/n8zBzaPb9lgrYQzlAqYmHn6CSOiFttFgd/FiHm+UFfCOE8peQkHzmRw4uHLolFKVffWCJGJ3JiEKuTZopRPcKsO5Wv4LPHUgiY5kcN7n+nepo3UvOYRUQ+nx099cm017J07B6fVZPTaOFmJISwNrggG/6FL5dOzPHzRXx6Kq8MYqMYCyEsSU9lynjwvYxYibxWvbAYswUXz31xCSujtop4ozQkhIRl9rhNXj4yi4H9aRz7wzyT5cW/yR4VF33tz261hJRlxgvit3ix82xBXK34IQbi2f4YxodXok8csAm/TpWw/UBGGUMTl7sYdbMWRayK2ehJBDAzV0FcavVbpeztWxPE+s7qkxBT9n/vYJ9cN8hyutZI1KgrRJWrclzc2xPC8/0dXmtjHJUSdvfRWTiyL5j9GjkvFkP7KUpBxExLRGjld9wRxj097d67V+fP2TK+EtsxOpFXz7XiBlWfCYYHoqssPJ8Ocq+sllqjtyuA7riNFaFq1mE7K8Nz02VMpEvK/jMx0HY0Q0ANIyHzoShudg6UPyusWuSby4XjpQ0PSCpZyuVzNRpOvxwgZ5hOlZknJkVLTKLCM4YR4MF5vUSQhoW0GstCWo1lIa3G/0IIzw97sT97+QWOnmeazWLGz2IoIsqSgv/vVLXk/hPDMXPs1ODrfzyjixgUAzvYG8HfUJ8bBF/As8cAAAAASUVORK5CYII="
          alt="Compartir en Twitter"
        />
      </a>
      <a class="link-img" :href="enlaces.telegram" target="_blank">
        <img
          title="Compartir en Telegram"
          class="rounded-circle img-thumbnail"
          src="data:image/png;base64,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"
          alt="Compartir en Telegram"
        />
      </a>
      <a class="link-img" :href="enlaces.whatsapp" target="_blank">
        <img
          title="Compartir en WhatsApp"
          class="rounded-circle img-thumbnail"
          src="data:image/png;base64,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"
          alt="Compartir en WhatsApp"
        />
      </a>
    </v-flex>
  </v-layout>
</template>

<script>
import { EventBus } from "../main";
import { HTTP } from "../http-common";
import { EnlaceCompartir, MensajeCompartir } from "../constantes";

export default {
  name: "VerificarNegocio",
  mounted() {
    EventBus.$emit("ocultarMenu");
    EventBus.$emit("ponerTitulo", "Iniciar sesión");
    this.idIntervalo = setInterval(() => this.siguientePaso(), 1000);
    this.verificarToken(this.$route.params.token);
  },
  data: () => ({
    enlaces: {
      facebook: `https://www.facebook.com/sharer.php?u=${EnlaceCompartir}&t=${MensajeCompartir}`,
      twitter: `https://www.twitter.com/intent/tweet?hashtags=PuntoDeVenta,POS,App,SublimePOS&text=${MensajeCompartir}&url=${EnlaceCompartir}`,
      whatsapp: `https://api.whatsapp.com/send?text=${MensajeCompartir}${encodeURIComponent(": ")}${EnlaceCompartir}`,
      telegram: `https://telegram.me/share/url?url=${EnlaceCompartir}&text=${MensajeCompartir}`,
    },
    idIntervalo: null,
    valor: 0,
    procesoTerminado: false,
    numeroPaso: -1,
    pasos: ["Creando base de datos", "Preparando primer uso", "Registrando primer usuario", "Dando permisos al usuario", "A punto de terminar..."],
  }),
  computed: {
    paso() {
      if (this.numeroPaso < 0) return "";
      if (this.numeroPaso < this.pasos.length) {
        return this.pasos[this.numeroPaso];
      }
      return "";
    }
  },
  methods: {
    irAlLogin() {
      this.$router.push({ name: "Login" });
    },
    verificarToken(token) {
      HTTP.get(`negocio/verificar/${token}`)
        .then(respuesta => {
          console.log("Al verificar con el servidor:", respuesta);
        }).finally(() => {
          this.procesoTerminado = true;
        });
    },
    siguientePaso() {
      if (this.numeroPaso < this.pasos.length - 1) {
        this.numeroPaso++;
        this.valor += (100 / (this.pasos.length + 1));
      } else {
        clearInterval(this.idIntervalo);
      }
    }
  },
}
</script>

<style scoped>
</style>
