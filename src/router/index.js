import Vue from 'vue'
import Router from 'vue-router'
import { HTTP } from "../http-common";

Vue.use(Router);

let router = new Router({
  routes: [
    {
      path: '/',
      name: 'Inicio',
      component: require("@/components/Inicio").default
    },
    {
      path: "/inventario",
      name: "Inventario",
      component: require("@/components/Inventario/Productos").default
    },
    {
      path: "/vender",
      name: "Vender",
      component: require("@/components/Vender/Vender").default
    },
    {
      path: "/clientes",
      name: "Clientes",
      component: require("@/components/Clientes/Clientes").default
    },
    {
      path: "/reporte/ventas/contado",
      name: "ReporteVentasContado",
      component: require("@/components/Reportes/Ventas/Contado").default
    },
    {
      path: "/reporte/apartados",
      name: "ReporteApartados",
      component: require("@/components/Reportes/Ventas/Apartado").default
    },
    {
      path: "/reporte/caja",
      name: "ReporteCaja",
      component: require("@/components/Reportes/Caja/Caja").default
    },
    {
      path: "/reporte/inventario",
      name: "ReporteInventario",
      component: require("@/components/Reportes/Inventario").default
    },
    {
      path: "/imprimir/codigos",
      name: "ImprimirCodigosDeBarras",
      component: require("@/components/Otros/CodigosBarra").default
    },
    {
      path: "/hacer/inventario",
      name: "HacerInventario",
      component: require("@/components/Otros/HacerInventario").default
    },
    {
      path: "/caja",
      name: "Caja",
      component: require("@/components/Caja/Caja").default
    },
    {
      path: "/reporte/stock",
      name: "ReporteStock",
      component: require("@/components/Reportes/Stock").default
    },
    {
      path: "/creditos",
      name: "AcercaDe",
      component: require("@/components/AcercaDe").default
    },
    {
      path: "/escritorio",
      name: "Escritorio",
      component: require("@/components/Escritorio").default
    },
    {
      path: "/graficas",
      name: "Graficas",
      component: require("@/components/Graficas/Graficas").default
    },
    {
      path: "/ajustes",
      name: "Ajustes",
      component: require("@/components/Ajustes/Ajustes").default
    },
    {
      path: "/login",
      name: "Login",
      component: require("@/components/Login").default
    },
    {
      path: '/usuarios',
      name: 'Usuarios',
      component: require("@/components/Usuarios/Usuarios").default
    },
    {
      path: '/logout',
      name: 'Logout',
      component: require("@/components/Logout").default
    },
    {
      path: "/registro",
      name: "Registro",
      component: require("@/components/Registro/Registro").default
    },
    {
      path: "/verificar/:token",
      name: "VerificarNegocio",
      component: require("@/components/VerificarNegocio").default
    },
    {
      path: "/eliminar/:token",
      name: "EliminarNegocio",
      component: require("@/components/EliminarNegocio").default
    },

    // Le tickets
    {
      path: "/ticket/venta/contado/:idVenta",
      name: "TicketDeVentaContado",
      component: require("@/components/Tickets/TicketVentaContado").default
    },
    {
      path: "/ticket/apartado/:idApartado",
      name: "TicketDeApartado",
      component: require("@/components/Tickets/TicketApartado").default
    },
    {
      path: "/ticket/abono/:idAbono/apartado/:idApartado",
      name: "TicketDeAbono",
      component: require("@/components/Tickets/TicketAbono").default
    },
    {
      path: "/ticket/caja/",
      name: "TicketDeCaja",
      component: require("@/components/Tickets/Caja").default
    },
 {
      path: "/nube",
      name: "UsarEnLaNube",
      component: require("@/components/Nube").default
    },
  ]
});
router.beforeEach((destino, origen, siguiente) => {
  if (["Login", "Logout", "Registro", "VerificarNegocio", "EliminarNegocio"].indexOf(destino.name) !== -1) siguiente();
  else {
    HTTP.get("estoy/logueado").then(estaLogueado => {
      if (estaLogueado) {
        siguiente();
      }
      else {
        siguiente({ name: 'Login' });
      }
    });
  }
});
export default router
