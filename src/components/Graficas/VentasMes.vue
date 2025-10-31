<script>
import { Line } from "vue-chartjs";
import { HTTP_AUTH } from "../../http-common";
import { FUNCIONES } from "../../funciones";
import { MESES } from "../../constantes";

export default {
  extends: Line,
  data: () => ({
    anio: null,
    mes: null,
  }),
  methods: {
    setAnioYMes(anio, mes) {
      this.anio = anio;
      this.mes = mes;
      this.refrescarConMesYAnio(anio, mes);
    },
    refrescarConMesYAnio(anio, mes) {
      HTTP_AUTH.get(`total/vendido/por/dia/${anio}/${mes}`).then(
        etiquetasYValores => {
          let etiquetas = [];
          let data = [];
          etiquetasYValores.forEach(etiquetaYValor => {
            etiquetas.push(etiquetaYValor.Etiqueta);
            data.push(etiquetaYValor.Valor);
          });
          let datasets = [
            {
              label: "Cantidad vendida",
              backgroundColor: FUNCIONES.colorHexadecimalAleatorio(),
              data
            }
          ];
          this.$nextTick(() => {
            if (this.$data._chart) {
              //Reiniciar gráfica. Gracias a https://parzibyte.me/blog/2018/05/03/reiniciar-limpiar-grafica-chart-js/
              this.$data._chart.clear();
              this.$data._chart.destroy();
            }
            this.renderChart(
              { labels: etiquetas, datasets },
              {
                tooltips: {
                  callbacks: {
                    label: (t, d) => {
                      let xLabel = d.datasets[t.datasetIndex].label;
                      let yLabel = t.yLabel;
                      return xLabel + ": " + this.$options.filters.currency(yLabel);
                    },
                    title: (a) => `${a[0].xLabel} de ${MESES[parseInt(this.mes) - 1]}`
                  }
                },
                title: {
                  text: `Ventas de ${MESES[parseInt(this.mes) - 1]}`,
                  display: true
                },

                responsive: true,
                maintainAspectRatio: false,
                scales: {
                  yAxes: [
                    {
                      display: true,
                      ticks: {
                        beginAtZero: true
                      }
                    }
                  ]
                }
              }
            );
          });

        }
      );
    }
  }
};
</script>
