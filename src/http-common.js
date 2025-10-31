import { EventBus } from "./main";
import { PuertoApi } from "./constantes";
const RUTA_SERVIDOR_CON_AUTH = `${window.document.location.protocol}//${window.document.location.hostname}:${PuertoApi}/auth/`;
const RUTA_SERVIDOR_SIN_AUTH = `${window.document.location.protocol}//${window.document.location.hostname}:${PuertoApi}/`;

const checarSiPermisoFueDenegado = datos => {
  if (null !== datos && undefined !== datos && "" !== datos) {
    if (datos.Clave && datos.Numero && datos.Numero === 21) {
      EventBus.$emit("permisoDenegado", datos);
    }
  }
  return datos;
};

const manejarError = error => {
  EventBus.$emit("errorDeServidor", error);
  console.error(error);
};
export const HTTP = {
  get(ruta) {
    return fetch(RUTA_SERVIDOR_SIN_AUTH + ruta, {
      credentials: 'include'
    })
      .then(respuesta => respuesta.json())
      .catch(error => {
        manejarError(error);
      })
  },
  post(ruta, objeto) {
    return fetch(RUTA_SERVIDOR_SIN_AUTH + ruta, {
      method: "POST",
      credentials: 'include',
      body: JSON.stringify(objeto)
    })
      .then(respuesta => respuesta.json())
      .catch(error => {
        manejarError(error);
      })
  },
  put(ruta, objeto) {
    return fetch(RUTA_SERVIDOR_SIN_AUTH + ruta, {
      method: "PUT",
      credentials: 'include',
      body: JSON.stringify(objeto),
    })
      .then(respuesta => respuesta.json())
      .catch(error => {
        manejarError(error);
      })
  },
};

export const HTTP_AUTH = {
  postArchivo(ruta, formData) {
    return fetch(RUTA_SERVIDOR_CON_AUTH + ruta, {
      method: "POST",
      credentials: 'include',
      body: formData
    })
      .then(respuesta => {
        return respuesta.json().then(respuestaDecodificada => {
          return checarSiPermisoFueDenegado(respuestaDecodificada);
        });
      })
      .catch(error => {
        manejarError(error);
      })
  },
  post(ruta, objeto) {
    return fetch(RUTA_SERVIDOR_CON_AUTH + ruta, {
      method: "POST",
      credentials: 'include',
      body: JSON.stringify(objeto)
    })
      .then(respuesta => {
        return respuesta.json().then(respuestaDecodificada => {
          return checarSiPermisoFueDenegado(respuestaDecodificada);
        });
      })
      .catch(error => {
        manejarError(error);
      })
  },
  put(ruta, objeto) {
    return fetch(RUTA_SERVIDOR_CON_AUTH + ruta, {
      method: "PUT",
      credentials: 'include',
      body: JSON.stringify(objeto),
    })
      .then(respuesta => {
        return respuesta.json().then(respuestaDecodificada => {
          return checarSiPermisoFueDenegado(respuestaDecodificada);
        });
      })
      .catch(error => {
        manejarError(error);
      })
  },
  get(ruta) {
    return fetch(RUTA_SERVIDOR_CON_AUTH + ruta, {
      credentials: 'include'
    })
      .then(respuesta => {
        return respuesta.json().then(respuestaDecodificada => {
          return checarSiPermisoFueDenegado(respuestaDecodificada);
        });
      })
      .catch(error => {
        manejarError(error);
      })
  },
  delete(ruta) {
    return fetch(RUTA_SERVIDOR_CON_AUTH + ruta, {
      method: "DELETE",
      credentials: 'include',
    })
      .then(respuesta => {
        return respuesta.json().then(respuestaDecodificada => {
          return checarSiPermisoFueDenegado(respuestaDecodificada);
        });
      })
      .catch(error => {
        manejarError(error);
      })
  }
};
