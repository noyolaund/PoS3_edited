CREATE TABLE IF NOT EXISTS negocios (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  nombre         TEXT        NOT NULL,
  correo         TEXT UNIQUE NOT NULL,
  pass           TEXT        NOT NULL,
  token          TEXT         NOT NULL,
  verificado     INTEGER                NOT NULL,
  fecha_registro TEXT           not null
);

CREATE TABLE IF NOT EXISTS accesos_negocios (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  id_negocio INTEGER not null,
  momento    TEXT,
  exitoso    INTEGER,
  foreign key (id_negocio) references negocios (id)
    on delete cascade
);

create table if not exists notificaciones_eliminacion_negocios (
  id         INTEGER primary key AUTOINCREMENT,
  id_negocio integer not null,
  token      TEXT,
  foreign key (id_negocio) references negocios (id)
    on delete cascade
);

