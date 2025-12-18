const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');
const { exec } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

const app = express();
const PORT = 8001;

app.use(cors());
app.use(bodyParser.json());

// Endpoint para obtener impresoras
app.get('/impresoras', (req, res) => {
    console.log('[Bridge] Solicitud GET /impresoras recibida');
    const command = 'powershell "Get-Printer | Select-Object Name | ConvertTo-Json"';

    exec(command, (error, stdout, stderr) => {
        if (error) {
            console.error(`[Bridge] Error obteniendo impresoras: ${error.message}`);
            return res.status(500).json({ error: error.message });
        }
        try {
            const printers = JSON.parse(stdout);
            const printerList = Array.isArray(printers) ? printers.map(p => p.Name) : [printers.Name];
            console.log(`[Bridge] Impresoras encontradas: ${printerList.join(', ')}`);
            res.json(printerList);
        } catch (e) {
            console.error("[Bridge] Error parseando impresoras:", e);
            res.json([]);
        }
    });
});

// Endpoint para imprimir
app.post('/imprimir', (req, res) => {
    const { nombreImpresora, operaciones } = req.body;

    console.log(`[Bridge] Solicitud POST /imprimir recibida`);
    console.log(`[Bridge] Impresora: ${nombreImpresora}`);
    console.log(`[Bridge] Número de operaciones: ${operaciones ? operaciones.length : 0}`);

    if (!nombreImpresora || !operaciones) {
        console.error("[Bridge] Faltan datos en la solicitud");
        return res.status(400).json({ error: "Faltan datos" });
    }

    let bufferParts = [];

    try {
        for (const op of operaciones) {
            if (Array.isArray(op) || Buffer.isBuffer(op)) {
                bufferParts.push(Buffer.from(op));
            } else if (typeof op === 'object') {
                const values = Object.values(op);
                bufferParts.push(Buffer.from(values));
            }
        }

        const finalBuffer = Buffer.concat(bufferParts);
        console.log(`[Bridge] Buffer total: ${finalBuffer.length} bytes`);

        const tempFilePath = path.join(os.tmpdir(), `print_${Date.now()}.bin`);
        fs.writeFileSync(tempFilePath, finalBuffer);
        console.log(`[Bridge] Archivo temporal creado: ${tempFilePath}`);

        // Método 1: Intentar con copy /b (requiere impresora compartida)
        console.log(`[Bridge] Intentando método 1: copy /b...`);
        const copyCommand = `copy /b "${tempFilePath}" "\\\\localhost\\${nombreImpresora}"`;

        exec(copyCommand, (error, stdout, stderr) => {
            if (error) {
                console.log(`[Bridge] copy /b falló: ${error.message}`);
                console.log(`[Bridge] Intentando método 2: print /D...`);

                // Método 2: Usar el comando print de Windows
                const printCommand = `print /D:"${nombreImpresora}" "${tempFilePath}"`;

                exec(printCommand, (error2, stdout2, stderr2) => {
                    // Limpiar archivo temporal
                    try { fs.unlinkSync(tempFilePath); } catch (e) { }

                    if (error2) {
                        console.error(`[Bridge] print /D también falló: ${error2.message}`);
                        console.error(`[Bridge] stderr: ${stderr2}`);
                        return res.json({
                            success: false,
                            error: `Ambos métodos fallaron. copy: ${error.message.substring(0, 100)}, print: ${error2.message.substring(0, 100)}`
                        });
                    }

                    console.log(`[Bridge] ✓ Impresión exitosa con print /D`);
                    console.log(`[Bridge] stdout: ${stdout2}`);
                    res.json(true);
                });
            } else {
                // Limpiar archivo temporal
                try { fs.unlinkSync(tempFilePath); } catch (e) { }

                console.log(`[Bridge] ✓ Impresión exitosa con copy /b`);
                console.log(`[Bridge] stdout: ${stdout}`);
                res.json(true);
            }
        });

    } catch (e) {
        console.error("[Bridge] Error procesando buffer:", e);
        res.status(500).json({ error: e.message });
    }
});

app.listen(PORT, () => {
    console.log(`Bridge de impresión corriendo en http://localhost:${PORT}`);
});
