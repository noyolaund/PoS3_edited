<template>
  <v-flex xs12>
    <a @click="mostrarDialogo = true">Exportar o importar productos</a>
    <v-snackbar :timeout="5000" :bottom="true" v-model="snackbars.exportacion">
      Exportado con éxito
      <v-btn flat color="pink" @click.native="snackbars.exportacion = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar
      :timeout="5000"
      :bottom="true"
      v-model="snackbars.noHayArchivos"
    >
      No ha seleccionado ningún archivo
      <v-btn flat color="pink" @click.native="snackbars.noHayArchivos = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-snackbar :timeout="5000" :bottom="true" v-model="snackbars.importacion">
      Importado con éxito
      <v-btn flat color="pink" @click.native="snackbars.importacion = false"
        >OK</v-btn
      >
    </v-snackbar>
    <v-dialog v-model="mostrarDialogo" persistent fullscreen>
      <v-card>
        <v-toolbar dark color="cyan">
          <v-btn icon dark @click.native="mostrarDialogo = false">
            <v-icon>close</v-icon>
          </v-btn>
          <v-toolbar-title>Exportar o importar productos</v-toolbar-title>
        </v-toolbar>
        <v-card-text>
          <v-tabs centered color="cyan" dark icons-and-text>
            <v-tabs-slider color="yellow"></v-tabs-slider>
            <v-tab href="#tab-1"
              >Exportar
              <v-icon>cloud_download</v-icon>
            </v-tab>
            <v-tab href="#tab-2"
              >Importar
              <v-icon>cloud_upload</v-icon>
            </v-tab>
            <v-tab-item value="tab-1">
              <v-card flat>
                <v-card-text>
                  <v-radio-group v-model="configuracionParaExportar.Extension">
                    <p>Formato de exportación</p>
                    <v-radio
                      label="CSV o valores separados por coma"
                      value="csv"
                    ></v-radio>
                    <v-radio
                      label="Formato de Microsoft Excel"
                      value="xlsx"
                    ></v-radio>
                  </v-radio-group>
                  <v-checkbox
                    label="Incluir encabezado"
                    v-model="configuracionParaExportar.IncluirEncabezado"
                  ></v-checkbox>
                  <v-text-field
                    v-model.number="configuracionParaExportar.Copias"
                    label="Copias"
                    type="number"
                    hint="¿Cuántos copias de cada producto desea exportar? por defecto es 1"
                  ></v-text-field>
                  <v-alert outline type="info" v-model="mostrarDialogo">
                    <p>
                      El archivo será exportado en el directorio en donde se
                      está ejecutando esta aplicación
                      <br />
                      Tenga en cuenta que el tiempo para generar el archivo
                      dependerá de la velocidad de su equipo y de la cantidad de
                      productos
                      <br />
                      <strong>Advertencia: </strong> cualquier archivo exportado
                      anteriormente será eliminado
                    </p>
                  </v-alert>
                  <v-btn
                    :loading="cargando"
                    small
                    color="success"
                    @click="exportar"
                    >Comenzar a exportar</v-btn
                  >
                </v-card-text>
              </v-card>
            </v-tab-item>
            <v-tab-item value="tab-2">
              <v-card flat>
                <v-card-text>
                  <p>Seleccione el archivo de Excel</p>
                  <v-container fluid grid-list-md>
                    <v-layout row wrap>
                      <v-flex xs12>
                        <input
                          accept=".xlsx"
                          ref="archivoImportar"
                          name="archivoImportar"
                          id="archivoImportar"
                          type="file"
                        />
                      </v-flex>
                    </v-layout>
                  </v-container>
                  <v-checkbox
                    label="¿El archivo tiene encabezado?"
                    v-model="configuracionParaImportar.TieneEncabezados"
                    hint="Si selecciona esta opción, la primera fila será omitida"
                  ></v-checkbox>
                  <v-select
                    :items="opcionesParaImportarRepetidos"
                    label="En caso de que se encuentren códigos de barras repetidos..."
                    hint="Si se encuentran códigos de barras repetidos, ¿qué desea hacer?"
                    outline
                    v-model="configuracionParaImportar.EnRepetidos"
                  ></v-select>
                  <h2>Configuración de columnas</h2>
                  <v-select
                    :items="indicesColumnas"
                    item-text="descripcion"
                    item-value="indice"
                    label="La columna que tiene el código de barras es la número..."
                    outline
                    v-model="configuracionParaImportar.IndiceCodigoBarras"
                  ></v-select>
                  <v-select
                    :items="indicesColumnas"
                    item-text="descripcion"
                    item-value="indice"
                    label="La columna que tiene la descripción del producto es la número..."
                    outline
                    v-model="configuracionParaImportar.IndiceDescripcion"
                  ></v-select>
                  <v-select
                    :items="indicesColumnas"
                    item-text="descripcion"
                    item-value="indice"
                    label="La columna que tiene el precio de compra es la número..."
                    outline
                    v-model="configuracionParaImportar.IndicePrecioCompra"
                  ></v-select>
                  <v-select
                    :items="indicesColumnas"
                    item-text="descripcion"
                    item-value="indice"
                    label="La columna que tiene el precio de venta es la número..."
                    outline
                    v-model="configuracionParaImportar.IndicePrecioVenta"
                  ></v-select>
                  <v-select
                    :items="indicesColumnas"
                    item-text="descripcion"
                    item-value="indice"
                    label="La columna que tiene la existencia es la número..."
                    outline
                    v-model="configuracionParaImportar.IndiceExistencia"
                  ></v-select>
                  <v-select
                    :items="indicesColumnas"
                    item-text="descripcion"
                    item-value="indice"
                    label="La columna que tiene la cantidad mínima del producto es la número..."
                    outline
                    v-model="configuracionParaImportar.IndiceStock"
                  ></v-select>
                  <v-alert outline type="info" v-model="mostrarDialogo">
                    <p>
                      Tenga en cuenta que el tiempo para importar los productos
                      dependerá de la velocidad de su equipo y de la cantidad de
                      los mismos
                      <br />
                      <strong>Advertencia: </strong> revise bien las opciones
                      que seleccionó antes de importar, para evitar resultados
                      inesperados
                    </p>
                  </v-alert>
                  <v-btn
                    :loading="cargando"
                    small
                    color="success"
                    @click="importar"
                    >Comenzar a importar</v-btn
                  >
                </v-card-text>
              </v-card>
            </v-tab-item>
          </v-tabs>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="gray" flat="flat" @click.native="mostrarDialogo = false"
            >Cerrar</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-flex>
</template>
<script>
import { HTTP_AUTH } from "../../http-common";

export default {
  name: "ExportarImportar",
  methods: {
    exportar() {
      this.cargando = true;
      HTTP_AUTH.put("exportar", this.configuracionParaExportar)
        .then(respuesta => {
          this.cargando = false;
          if (respuesta) {
            this.snackbars.exportacion = true;
          }
        });
    },
    importar() {
      let archivos = this.$refs.archivoImportar.files;
      if (archivos.length > 0) {
        this.cargando = true;
        let archivo = archivos[0];
        let formData = new FormData();
        formData.append("archivo", archivo);
        formData.append("TieneEncabezados", this.configuracionParaImportar.TieneEncabezados.toString());
        formData.append("IgnorarCodigosDeBarrasRepetidos", (this.configuracionParaImportar.EnRepetidos === "Ignorar").toString());
        formData.append("IndiceCodigoBarras", this.configuracionParaImportar.IndiceCodigoBarras.toString());
        formData.append("IndiceDescripcion", this.configuracionParaImportar.IndiceDescripcion.toString());
        formData.append("IndicePrecioCompra", this.configuracionParaImportar.IndicePrecioCompra.toString());
        formData.append("IndicePrecioVenta", this.configuracionParaImportar.IndicePrecioVenta.toString());
        formData.append("IndiceExistencia", this.configuracionParaImportar.IndiceExistencia.toString());
        formData.append("IndiceStock", this.configuracionParaImportar.IndiceStock.toString());
        HTTP_AUTH.postArchivo("importar/excel", formData).then(respuesta => {
          this.cargando = false;
          if (respuesta) {
            this.snackbars.importacion = true;
            this.$emit("importado");
          }
        });
      } else {
        this.snackbars.noHayArchivos = true;
      }
    }
  },
  data: () => ({
    snackbars: {
      exportacion: false,
      noHayArchivos: false,
      importacion: false,
    },
    opcionesParaImportarRepetidos: ["Ignorar", "Remplazar"],
    indicesColumnas: [
      {
        indice: 0,
        descripcion: "El archivo no contiene esta columna"
      },
      {
        indice: 1,
        descripcion: 1
      },
      {
        indice: 2,
        descripcion: 2
      },
      {
        indice: 3,
        descripcion: 3
      },
      {
        indice: 4,
        descripcion: 4
      },
      {
        indice: 5,
        descripcion: 5
      },
      {
        indice: 6,
        descripcion: 6
      },
      {
        indice: 7,
        descripcion: 7
      },
      {
        indice: 8,
        descripcion: 8
      },
      {
        indice: 9,
        descripcion: 9
      },
      {
        indice: 10,
        descripcion: 10
      },
    ],
    configuracionParaExportar: {
      Extension: "csv",
      Copias: 1,
      IncluirEncabezado: true,
    },
    configuracionParaImportar: {
      TieneEncabezados: true,
      IndiceCodigoBarras: 1,
      IndiceDescripcion: 2,
      IndicePrecioCompra: 3,
      IndicePrecioVenta: 4,
      IndiceExistencia: 5,
      IndiceStock: 6,
      EnRepetidos: "Ignorar",
    },
    mostrarDialogo: false,
    cargando: false,
  })
}
</script>
