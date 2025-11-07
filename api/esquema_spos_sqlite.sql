CREATE TABLE IF NOT EXISTS productos (
  idProducto INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  codigoBarras TEXT NOT NULL,
  descripcion TEXT NOT NULL,
  precioCompra REAL NOT NULL,
  precioVenta REAL NOT NULL,
  existencia REAL NOT NULL,
  stock REAL NOT NULL,
  -- Para relación padre-hijo
  idPadre INTEGER NOT NULL DEFAULT 0,
  equivalencia INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS clientes (
  idCliente INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  nombreCompleto TEXT NOT NULL,
  numeroTelefono TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS usuarios (
  idUsuario INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  nombre TEXT NOT NULL UNIQUE,
  contraseña TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS permisos (
  idPermiso INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  clave TEXT NOT NULL,
  descripcion TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS permisos_usuarios (
  idUsuario INTEGER NOT NULL,
  idPermiso INTEGER NOT NULL,
  FOREIGN KEY (idUsuario) REFERENCES usuarios (idUsuario) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (idPermiso) REFERENCES permisos (idPermiso) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS ventas_contado (
  idVenta INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  monto REAL NOT NULL,
  pago REAL NOT NULL DEFAULT 0,
  fecha TEXT NOT NULL,
  idCliente INTEGER NOT NULL,
  idUsuario INTEGER NOT NULL,
  FOREIGN KEY (idCliente) REFERENCES clientes (idCliente) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (idUsuario) REFERENCES usuarios (idUsuario) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS productos_vendidos (
  idProducto INTEGER NOT NULL,
  codigoBarras TEXT,
  idVenta INTEGER NOT NULL,
  descripcion TEXT NOT NULL,
  precioCompra REAL NOT NULL,
  precioVenta REAL NOT NULL,
  precioVentaOriginal REAL NOT NULL,
  /**/
  cantidadVendida REAL NOT NULL,
  FOREIGN KEY (idVenta) REFERENCES ventas_contado (idVenta) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS apartados (
  idApartado INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  monto REAL NOT NULL,
  pago REAL NOT NULL DEFAULT 0,
  abonado REAL NOT NULL,
  anticipo REAL NOT NULL,
  fecha TEXT NOT NULL,
  fechaVencimiento TEXT NOT NULL,
  idCliente INTEGER NOT NULL,
  idUsuario INTEGER NOT NULL,
  FOREIGN KEY (idCliente) REFERENCES clientes (idCliente) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (idUsuario) REFERENCES usuarios (idUsuario) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS productos_apartados (
  idApartado INTEGER NOT NULL,
  idProducto INTEGER NOT NULL,
  codigoBarras TEXT,
  descripcion TEXT NOT NULL,
  precioVenta REAL NOT NULL,
  precioVentaOriginal REAL NOT NULL,
  precioCompra REAL NOT NULL,
  cantidadVendida REAL NOT NULL,
  FOREIGN KEY (idApartado) REFERENCES apartados (idApartado) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (idProducto) REFERENCES productos (idProducto) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS abonos (
  idAbono INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  monto REAL NOT NULL,
  pago REAL NOT NULL DEFAULT 0,
  fecha TEXT NOT NULL,
  idApartado INTEGER NOT NULL,
  idUsuario INTEGER NOT NULL,
  FOREIGN KEY (idApartado) REFERENCES apartados (idApartado) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (idusuario) REFERENCES usuarios (idUsuario) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS ingresos (
  monto REAL NOT NULL,
  descripcion TEXT NOT NULL,
  fecha TEXT NOT NULL,
  idUsuario INTEGER NOT NULL,
  FOREIGN KEY (idusuario) REFERENCES usuarios (idUsuario) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS egresos (
  monto REAL NOT NULL,
  descripcion TEXT NOT NULL,
  fecha TEXT NOT NULL,
  idUsuario INTEGER NOT NULL,
  FOREIGN KEY (idusuario) REFERENCES usuarios (idUsuario) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS empresa (
  nombre TEXT,
  direccion TEXT,
  telefono TEXT,
  mensajePersonal TEXT
);

CREATE TABLE IF NOT EXISTS comun (clave TEXT, valor TEXT);

CREATE INDEX idVenta_indice ON productos_vendidos (idVenta);

CREATE INDEX clave_permiso ON permisos (clave);