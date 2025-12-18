
export default class ConectorJavascript {
    constructor(ruta, serial) {
        this.ruta = ruta || "http://localhost:8001";
        this.serial = serial;
        this.operaciones = [];
    }

    static URL_PLUGIN_POR_DEFECTO = "http://localhost:8001";

    static async obtenerImpresoras(ruta) {
        if (ruta) ConectorJavascript.URL_PLUGIN_POR_DEFECTO = ruta;
        console.log(`[ConectorJS] Obteniendo impresoras de: ${ConectorJavascript.URL_PLUGIN_POR_DEFECTO}/impresoras`);
        try {
            const response = await fetch(ConectorJavascript.URL_PLUGIN_POR_DEFECTO + "/impresoras");
            console.log(`[ConectorJS] Respuesta status: ${response.status}`);
            if (!response.ok) throw new Error("Error conectando con el bridge");
            const printers = await response.json();
            console.log(`[ConectorJS] Impresoras recibidas:`, printers);
            return printers;
        } catch (e) {
            console.error("[ConectorJS] Error obteniendo impresoras:", e);
            return [];
        }
    }

    Iniciar() {
        this.operaciones = [];
        this.agregarComando([0x1B, 0x40]);
        return this;
    }

    EscribirTexto(texto) {
        const encoder = new TextEncoder();
        const data = encoder.encode(texto);
        this.agregarDatos(Array.from(data));
        return this;
    }

    Feed(lineas) {
        this.agregarComando([0x1B, 0x64, lineas]);
        return this;
    }

    Corte(lineas) {
        this.agregarComando([0x1D, 0x56, 66, lineas]);
        return this;
    }

    CorteParcial() {
        this.agregarComando([0x1D, 0x56, 1]);
        return this;
    }

    EstablecerEnfatizado(enfatizado) {
        this.agregarComando([0x1B, 0x45, enfatizado ? 1 : 0]);
        return this;
    }

    EstablecerAlineacion(alineacion) {
        this.agregarComando([0x1B, 0x61, alineacion]);
        return this;
    }

    DeshabilitarElModoDeCaracteresChinos() {
        this.agregarComando([0x1C, 0x2E]);
        return this;
    }

    Pulso(pin, t1, t2) {
        let m = 0;
        if (pin === 48) m = 0;
        else if (pin === 49) m = 1;
        this.agregarComando([0x1B, 0x70, m, t1, t2]);
        return this;
    }

    ImprimirImagenEnBase64(imagenBase64, tamaño, anchoMaximo) {
        console.warn("Impresión de imagen no soportada en bridge JS simple por ahora.");
        return this;
    }

    agregarComando(bytes) {
        this.operaciones.push(bytes);
    }

    agregarDatos(bytes) {
        this.operaciones.push(bytes);
    }

    async imprimirEn(nombreImpresora) {
        const payload = {
            nombreImpresora,
            operaciones: this.operaciones
        };

        console.log(`[ConectorJS] Imprimiendo en: ${nombreImpresora}`);
        console.log(`[ConectorJS] Operaciones a enviar: ${this.operaciones.length}`);
        console.log(`[ConectorJS] URL: ${this.ruta}/imprimir`);

        try {
            const response = await fetch(this.ruta + "/imprimir", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(payload)
            });
            console.log(`[ConectorJS] Respuesta status: ${response.status}`);
            const result = await response.json();
            console.log(`[ConectorJS] Resultado:`, result);
            return result;
        } catch (e) {
            console.error("[ConectorJS] Error al imprimir:", e);
            return false;
        }
    }

    static TAMAÑO_IMAGEN_NORMAL = 0;
    static ALINEACION_IZQUIERDA = 0;
    static ALINEACION_CENTRO = 1;
    static ALINEACION_DERECHA = 2;
}
