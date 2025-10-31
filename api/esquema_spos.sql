CREATE TABLE IF NOT EXISTS productos (
  idProducto   BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  codigoBarras VARCHAR(255) UNIQUE,
  descripcion  VARCHAR(255)    NOT NULL,
  precioCompra DECIMAL(9, 2)   NOT NULL,
  precioVenta  DECIMAL(9, 2)   NOT NULL,
  existencia   DECIMAL(9, 2)   NOT NULL,
  stock        DECIMAL(9, 2)   NOT NULL
);

CREATE TABLE IF NOT EXISTS clientes (
  idCliente      BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  nombreCompleto VARCHAR(255)    NOT NULL,
  numeroTelefono VARCHAR(20)     NOT NULL
);

CREATE TABLE IF NOT EXISTS usuarios (
  idUsuario  BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  nombre     VARCHAR(255)    NOT NULL UNIQUE,
  contraseña VARCHAR(255)    NOT NULL
);

CREATE TABLE IF NOT EXISTS permisos (
  idPermiso   BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  clave       VARCHAR(255)    NOT NULL,
  descripcion VARCHAR(255)    NOT NULL
);

CREATE TABLE IF NOT EXISTS permisos_usuarios (
  idUsuario BIGINT UNSIGNED NOT NULL,
  idPermiso BIGINT UNSIGNED NOT NULL,
  foreign key (idUsuario) references usuarios (idUsuario)
    on delete cascade
    on update cascade,
  foreign key (idPermiso) references permisos (idPermiso)
    on delete cascade
    on update cascade
);

CREATE TABLE IF NOT EXISTS ventas_contado (
  idVenta   BIGINT UNSIGNED NOT NULL PRIMARY KEY    AUTO_INCREMENT,
  monto     DECIMAL(9, 2)   NOT NULL,
  pago      DECIMAL(9, 2)   NOT NULL                DEFAULT 0,
  fecha     VARCHAR(255)    NOT NULL,
  idCliente BIGINT UNSIGNED NOT NULL,
  idUsuario BIGINT UNSIGNED NOT NULL,
  foreign key (idCliente) references clientes (idCliente)
    on delete cascade
    on update cascade,
  foreign key (idUsuario) references usuarios (idUsuario)
    on delete cascade
    on update cascade
);

CREATE TABLE IF NOT EXISTS productos_vendidos (
  idProducto          BIGINT UNSIGNED NOT NULL,
  codigoBarras        VARCHAR(255),
  idVenta             BIGINT UNSIGNED NOT NULL,
  descripcion         VARCHAR(255)    NOT NULL,
  precioCompra        DECIMAL(9, 2)   NOT NULL,
  precioVenta         DECIMAL(9, 2)   NOT NULL,
  precioVentaOriginal DECIMAL(9, 2)   NOT NULL, /**/
  cantidadVendida     DECIMAL(9, 2)   NOT NULL,
  foreign key (idVenta) references ventas_contado (idVenta)
    on delete cascade
    on update cascade
);

CREATE TABLE IF NOT EXISTS apartados (
  idApartado       BIGINT UNSIGNED NOT NULL PRIMARY KEY    AUTO_INCREMENT,
  monto            DECIMAL(9, 2)   NOT NULL,
  pago             DECIMAL(9, 2)   NOT NULL                DEFAULT 0,
  abonado          DECIMAL(9, 2)   NOT NULL,
  anticipo         DECIMAL(9, 2)   NOT NULL,
  fecha            VARCHAR(255)    NOT NULL,
  fechaVencimiento VARCHAR(255)    NOT NULL,
  idCliente        BIGINT UNSIGNED NOT NULL,
  idUsuario        BIGINT UNSIGNED NOT NULL,
  foreign key (idCliente) references clientes (idCliente)
    on delete cascade
    on update cascade,
  foreign key (idUsuario) references usuarios (idUsuario)
    on delete cascade
    on update cascade
);

CREATE TABLE IF NOT EXISTS productos_apartados (
  idApartado          BIGINT UNSIGNED NOT NULL,
  idProducto          BIGINT UNSIGNED NOT NULL,
  codigoBarras        VARCHAR(255),
  descripcion         VARCHAR(255)    NOT NULL,
  precioVenta         DECIMAL(9, 2)   NOT NULL,
  precioVentaOriginal DECIMAL(9, 2)   NOT NULL,
  precioCompra        DECIMAL(9, 2)   NOT NULL,
  cantidadVendida     DECIMAL(9, 2)   NOT NULL,
  foreign key (idApartado) references apartados (idApartado)
    on delete cascade
    on update cascade,
  foreign key (idProducto) references productos (idProducto)
    on delete restrict
    on update cascade
);

CREATE TABLE IF NOT EXISTS abonos (
  idAbono    BIGINT UNSIGNED NOT NULL PRIMARY KEY    AUTO_INCREMENT,
  monto      DECIMAL(9, 2)   NOT NULL,
  pago       DECIMAL(9, 2)   NOT NULL                DEFAULT 0,
  fecha      VARCHAR(255)    NOT NULL,
  idApartado BIGINT UNSIGNED NOT NULL,
  idUsuario  BIGINT UNSIGNED NOT NULL,
  foreign key (idApartado) references apartados (idApartado)
    on delete cascade
    on update cascade,
  foreign key (idusuario) references usuarios (idUsuario)
    on delete cascade
    on update cascade
);

CREATE TABLE IF NOT EXISTS ingresos (
  monto       DECIMAL(9, 2)   NOT NULL,
  descripcion VARCHAR(255)    NOT NULL,
  fecha       VARCHAR(255)    NOT NULL,
  idUsuario   BIGINT UNSIGNED NOT NULL,
  foreign key (idusuario) references usuarios (idUsuario)
    on delete cascade
    on update cascade
);

CREATE TABLE IF NOT EXISTS egresos (
  monto       DECIMAL(9, 2)   NOT NULL,
  descripcion VARCHAR(255)    NOT NULL,
  fecha       VARCHAR(255)    NOT NULL,
  idUsuario   BIGINT UNSIGNED NOT NULL,
  foreign key (idusuario) references usuarios (idUsuario)
    on delete cascade
    on update cascade
);

CREATE TABLE IF NOT EXISTS empresa (
  nombre          VARCHAR(255),
  direccion       VARCHAR(255),
  telefono        VARCHAR(255),
  mensajePersonal VARCHAR(255)
);
CREATE TABLE IF NOT EXISTS comun (
  clave VARCHAR(255),
  valor VARCHAR(255)
);

CREATE INDEX idVenta_indice
  ON productos_vendidos (idVenta);
CREATE INDEX clave_permiso
  ON permisos (clave);
