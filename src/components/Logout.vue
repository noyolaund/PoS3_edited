<template></template>

<script>
import { HTTP } from "../http-common";
import { EventBus } from "../main";

export default {
  name: "Logout",
  beforeMount() {
    HTTP.get("logout").then(resultados => {
      EventBus.$emit("ponerNombreDeUsuario", "");
      EventBus.$emit("ocultarMenu");
      EventBus.$emit("ocultarToolbar");
      EventBus.$emit("ponerDatosUsuario", {});
      this.$router.push({ name: "Login" });
    });
  },
}
</script>
