// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import VueCurrencyFilter from 'vue-currency-filter'
import VueQrcode from '@xkeshi/vue-qrcode'

export const EventBus = new Vue();
Vue.component('qrcode', VueQrcode);
Vue.filter('fechaExpresiva', function (fecha) {
  if (!fecha) return "";
  fecha = new Date(fecha);
  return ["domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"][fecha.getDay()] + ", " + fecha.getDate().toString() + " de " + ["enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"][fecha.getMonth()] + " de " + fecha.getFullYear().toString() + " " + (fecha.getHours() > 9 ? fecha.getHours() : "0" + fecha.getHours()).toString() + ":" + (fecha.getMinutes() > 9 ? fecha.getMinutes() : "0" + fecha.getMinutes()).toString();
});
Vue.filter('fechaSinHora', function (fecha) {
  if (!fecha) return "";
  fecha = new Date(fecha);
  return fecha.getDate().toString() + " de " + ["enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"][fecha.getMonth()] + " de " + fecha.getFullYear().toString();
});

Vue.use(VueCurrencyFilter, {
  symbol: '$',
  thousandsSeparator: ',',
  fractionCount: 2,
  fractionSeparator: '.',
  symbolPosition: 'front',
  symbolSpacing: true
})

Vue.use(Vuetify)

Vue.config.productionTip = false
Vue.$router = router;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: {
    App
  },
  template: '<App/>'
})
