<script>
import { MESES } from "../../constantes";
import { Bar } from "vue-chartjs";
import { HTTP_AUTH } from "../../http-common";
import { FUNCIONES } from "../../funciones";

const NUMERO_DE_MESES = 12;
export default {
  extends: Bar,
  data: () => ({
    anio: null,
    aniosParaGraficar: [],
  }),
  methods: {
    agregarAnio(anio) {
      this.aniosParaGraficar.push(anio);
      this.refrescar();
    },
    removerAnio(indice) {
      this.aniosParaGraficar.splice(indice, 1);
      this.refrescar();
    },
    setAnios(anios) {
      this.aniosParaGraficar = anios;
      this.refrescar();
    },
    refrescar() {
      let contador = 0, etiquetas = [], datasets = [];
      const traerMasAnios = indice => {
        if (indice >= 0) {
          let datos = [];
          let anio = this.aniosParaGraficar[indice];
          HTTP_AUTH.get(`total/vendido/por/mes/${anio}`).then(etiquetasYValores => {
            contador++;
            let verdaderosValoresYEtiquetas = [];
            for (let x = 1; x < NUMERO_DE_MESES + 1; x++) {
              let mesComoString = FUNCIONES.agregarCerosALaIzquierdaSiEsNecesario(x);
              let indice = etiquetasYValores.findIndex(etiquetaYValor => etiquetaYValor.Etiqueta === mesComoString);
              if (indice !== -1) {
                verdaderosValoresYEtiquetas.push(etiquetasYValores[indice]);
              } else {
                verdaderosValoresYEtiquetas.push({ Etiqueta: mesComoString, Valor: 0 });
              }
            }

            verdaderosValoresYEtiquetas.forEach(etiquetaYValor => {

              if (etiquetas.length < NUMERO_DE_MESES) {
                etiquetas.push(MESES[parseInt(etiquetaYValor.Etiqueta) - 1]);
              }
              datos.push(etiquetaYValor.Valor);
            });

            datasets.push({
              label: anio,
              backgroundColor: FUNCIONES.colorHexadecimalAleatorio(),
              data: Array.from(datos)
            });
            traerMasAnios(indice - 1);
          });
        } else {
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
                      let yLabel = t.yLabel;
                      return "Vendido al contado: " + this.$options.filters.currency(yLabel);
                    },
                    title: a => {
                      return `${a[0].xLabel} del ${this.aniosParaGraficar[a[0].datasetIndex]}`
                    }
                  }
                },
                title: {
                  text: `Ventas por año`,
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
          })
        }
      };
      traerMasAnios(this.aniosParaGraficar.length - 1);
    }
  }
};
</script>
