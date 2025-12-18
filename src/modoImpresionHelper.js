import { HTTP_AUTH } from "./http-common";

// Helper para leer modo de impresión con fallback a localStorage
export async function obtenerModoImpresion() {
    let modo = await HTTP_AUTH.get("valor/MODO_IMPRESION");
    // Si el backend devuelve vacío, usar localStorage como respaldo
    if (!modo || modo === '') {
        modo = localStorage.getItem('MODO_IMPRESION') || '';
        console.log('[Helper] Backend devolvió vacío, usando localStorage:', modo);
    }
    return modo;
}
