CREATE TABLE IF NOT EXISTS negocios (
  id             INTEGER PRIMARY KEY AUTO_INCREMENT,
  nombre         VARCHAR(255)        NOT NULL,
  correo         VARCHAR(255) UNIQUE NOT NULL,
  pass           VARCHAR(255)        NOT NULL,
  token          VARCHAR(64)         NOT NULL,
  verificado     BOOL                NOT NULL,
  fecha_registro timestamp           not null
);

CREATE TABLE IF NOT EXISTS accesos_negocios (
  id         bigint unsigned PRIMARY KEY AUTO_INCREMENT,
  id_negocio INTEGER not null,
  momento    timestamp,
  exitoso    bool,
  foreign key (id_negocio) references negocios (id)
    on delete cascade
);

create table if not exists notificaciones_eliminacion_negocios (
  id         bigint unsigned primary key auto_increment,
  id_negocio integer not null,
  token      varchar(64),
  foreign key (id_negocio) references negocios (id)
    on delete cascade
);
