package main

import (
	"database/sql"
	"net/http"
)

type ClaveYValor struct {
	Clave, Valor string
}

/*
  Aquí residen todas las clases o tipos
  que son utilizados en el sistema

  Se pusieron todos aquí para evitar tener un
  archivo por cada clase
*/

type AjustesParaExportarProductos struct {
	Extension         string //Puede ser csv, xlsx o txt
	Copias            int    // Cuántas copias del producto?
	IncluirEncabezado bool   // En el archivo exportado incluir encabezados?
}

type AjustesParaImportarProductos struct {
	TieneEncabezados, IgnorarCodigosDeBarrasRepetidos                                                           bool
	IndiceCodigoBarras, IndiceDescripcion, IndicePrecioCompra, IndicePrecioVenta, IndiceExistencia, IndiceStock int
}

type Abono struct {
	IdApartado, IdAbono int
	Monto, Pago         float64
	Fecha               string
	Usuario             Usuario
}

type Ajustes struct {
	MensajePersonal, Nombre, Direccion, Telefono string
}

type Apartado struct {
	Numero                         int
	Total, Abonado, Anticipo, Pago float64
	FechaVencimiento, Fecha        string
	Productos                      []ProductoApartado
	Cliente                        Cliente
	Usuario                        Usuario
}

type Usuario struct {
	Numero                          int
	Nombre, Password, CorreoNegocio string
	Negocio                         Negocio
}

type ApartadoConAbonos struct {
	Apartado Apartado
	Abonos   []Abono
}

type Cliente struct {
	Numero                 int
	Nombre, NumeroTelefono string
}

type DatosEmpresa struct {
	MensajePersonal, Nombre, Direccion, Telefono string
}

type DetalleApartado struct {
	IdApartado               int
	Monto, Abonado, Anticipo float64
	Fecha, FechaVencimiento  string
	Producto                 ProductoApartado
	Cliente                  Cliente
	Usuario                  Usuario
}

type Egreso struct {
	Monto              float64
	Descripcion, Fecha string
	Usuario            Usuario
}

type EtiquetaYValor struct {
	Etiqueta string
	Valor    float64
}

type HistorialCliente struct {
	Apartados       *[]Apartado
	VentasAlContado *[]VentaContado
}

type Ingreso struct {
	Monto              float64
	Descripcion, Fecha string
	Usuario            Usuario
}

type OtrosAjustes struct {
	ModoImpresionCodigoDeBarras, ModoLecturaProductos                                    string
	NumeroDeCopiasTicketContado, NumeroDeCopiasTicketApartado, NumeroDeCopiasTicketAbono int
}

type Producto struct {
	Numero                                       int
	Descripcion, CodigoBarras                    string
	PrecioCompra, PrecioVenta, Existencia, Stock float64
	// Relación padre-hijo
	Padre        int // idProducto del producto padre. 0 si no tiene padre
	Equivalencia int // cuántas unidades hijo equivalen a 1 unidad padre (p.ej. 12)
}

type ProductoApartado struct {
	Numero                                                   int
	Descripcion, CodigoBarras                                string
	PrecioVenta, PrecioVentaOriginal, PrecioCompra, Cantidad float64
}

type ProductosConConteo struct {
	Conteo    int
	Productos []Producto
}

type ProductoVendido struct {
	Numero                                                   int
	Descripcion, CodigoBarras                                string
	PrecioVenta, PrecioVentaOriginal, PrecioCompra, Cantidad float64
}

type ProductoVendidoParaGrafica struct {
	IdProducto                int
	Descripcion, CodigoBarras string
	VecesVendido              float64
}

type ReporteCaja struct {
	VentasContado, Anticipos, Abonos, Ingresos, Egresos float64
}

type ReporteInventario struct {
	PrecioCompra, PrecioVenta, CantidadProductos float64
}

type VentaContado struct {
	Numero      int
	Total, Pago float64
	Fecha       string
	Productos   []ProductoVendido
	Cliente     Cliente
	Usuario     Usuario
}

type DetalleVentaContado struct {
	IdVenta  int
	Monto    float64
	Fecha    string
	Producto ProductoVendido
	Cliente  Cliente
	Usuario  Usuario
}

type Permiso struct {
	Id                 int
	Clave, Descripcion string
}

type ErrorDePermiso struct {
	Clave, Mensaje string
	Numero         int
	Permiso        Permiso
}

type Negocio struct {
	Id         int    `json:"id"`
	Nombre     string `json:"nombre"`
	Correo     string `json:"correo"`
	Token      string `json:"token"`
	Pass       string `json:"pass"`
	Verificado bool   `json:"verificado"`
}

type AjustesDeUsuarioLogueado struct {
	IdNegocio          int
	httpResponseWriter http.ResponseWriter
	httpRequest        *http.Request
}

func (a *AjustesDeUsuarioLogueado) obtenerIdNegocio() (int, error) {
	sc := Sesion{
		request:        a.httpRequest,
		responseWriter: a.httpResponseWriter,
	}
	usuario, err := sc.obtenerUsuarioLogueado()
	if err != nil {
		return 0, err
	}
	return usuario.Negocio.Id, nil
}

func (a *AjustesDeUsuarioLogueado) obtenerBaseDeDatos() (*sql.DB, error) {
	sc := Sesion{
		request:        a.httpRequest,
		responseWriter: a.httpResponseWriter,
	}
	usuario, err := sc.obtenerUsuarioLogueado()
	if err != nil {
		return nil, err
	}
	db, err := obtenerBaseDeDatosAPartirDeIdNegocio(usuario.Negocio.Id)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type DetallesDeNegocio struct {
	Negocio, Token, Correo, Clave string
}
