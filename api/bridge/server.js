const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');
const { exec } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

const app = express();
const PORT = 8001;

// Detectar sistema operativo
const SISTEMA = os.platform();
const ARQUITECTURA = os.arch();
console.log(`[Bridge] Sistema operativo detectado: ${SISTEMA}`);
console.log(`[Bridge] Arquitectura detectada: ${ARQUITECTURA}`);

// Determinar si es Raspberry Pi
const isRaspberryPi = ARQUITECTURA === 'arm64' || ARQUITECTURA === 'arm';
if (isRaspberryPi) {
    console.log('[Bridge] ✓ Detectado como Raspberry Pi');
}

app.use(cors());
app.use(bodyParser.json());

// Endpoint para obtener información del sistema
app.get('/sistema', (req, res) => {
    res.json({
        sistema: SISTEMA,
        arquitectura: ARQUITECTURA,
        esRaspberryPi: isRaspberryPi,
        node: process.version,
        plataforma: process.platform
    });
});

// Endpoint de prueba - Imprime un ticket de ejemplo
app.get('/prueba/:nombreImpresora', (req, res) => {
    const { nombreImpresora } = req.params;
    
    console.log(`[Bridge] Solicitud de prueba para impresora: ${nombreImpresora}`);
    
    try {
        // Crear un buffer de prueba simple
        const testContent = Buffer.from([
            ...Buffer.from('TICKET DE PRUEBA\n'),
            ...Buffer.from('='.repeat(32) + '\n'),
            ...Buffer.from(`Sistema: ${SISTEMA}\n`),
            ...Buffer.from(`Arquitectura: ${ARQUITECTURA}\n`),
            ...Buffer.from(`Impresora: ${nombreImpresora}\n`),
            ...Buffer.from(`Fecha: ${new Date().toLocaleString()}\n`),
            ...Buffer.from('='.repeat(32) + '\n'),
            ...Buffer.from('✓ Prueba exitosa\n'),
            0x0C // Form feed para Raspberry Pi
        ]);
        
        const tempFilePath = path.join(os.tmpdir(), `test_print_${Date.now()}.txt`);
        fs.writeFileSync(tempFilePath, testContent);
        
        console.log(`[Bridge] Archivo de prueba creado: ${tempFilePath}`);
        
        if (SISTEMA === 'linux') {
            const lpCommand = `lp -d "${nombreImpresora}" "${tempFilePath}" 2>&1`;
            exec(lpCommand, (error, stdout, stderr) => {
                try { fs.unlinkSync(tempFilePath); } catch (e) { }
                
                if (error) {
                    console.error(`[Bridge] Error enviando prueba: ${error.message}`);
                    return res.json({
                        success: false,
                        error: error.message,
                        stderr: stderr,
                        comando: lpCommand
                    });
                }
                
                console.log(`[Bridge] ✓ Ticket de prueba enviado exitosamente`);
                res.json({
                    success: true,
                    mensaje: 'Ticket de prueba enviado a la impresora',
                    stdout: stdout
                });
            });
        } else if (SISTEMA === 'win32') {
            const copyCommand = `copy /b "${tempFilePath}" "\\\\localhost\\${nombreImpresora}"`;
            exec(copyCommand, (error, stdout, stderr) => {
                try { fs.unlinkSync(tempFilePath); } catch (e) { }
                
                if (error) {
                    const printCommand = `print /D:"${nombreImpresora}" "${tempFilePath}"`;
                    exec(printCommand, (error2) => {
                        try { fs.unlinkSync(tempFilePath); } catch (e) { }
                        
                        if (error2) {
                            return res.json({
                                success: false,
                                error: error2.message
                            });
                        }
                        
                        res.json({
                            success: true,
                            mensaje: 'Ticket de prueba enviado a la impresora'
                        });
                    });
                } else {
                    res.json({
                        success: true,
                        mensaje: 'Ticket de prueba enviado a la impresora'
                    });
                }
            });
        } else {
            try { fs.unlinkSync(tempFilePath); } catch (e) { }
            res.json({
                success: false,
                error: `Sistema ${SISTEMA} no soportado aún`
            });
        }
        
    } catch (e) {
        console.error(`[Bridge] Error en prueba: ${e.message}`);
        res.json({
            success: false,
            error: e.message
        });
    }
});

// Endpoint para obtener impresoras
app.get('/impresoras', (req, res) => {
    console.log('[Bridge] Solicitud GET /impresoras recibida');
    
    if (SISTEMA === 'win32') {
        let command = 'powershell "Get-Printer | Select-Object Name | ConvertTo-Json"';
        exec(command, (error, stdout, stderr) => {
            if (error) {
                console.error(`[Bridge] Error obteniendo impresoras: ${error.message}`);
                return res.json([]);
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
    } else if (SISTEMA === 'linux') {
        // En Linux, buscar tanto en CUPS como en dispositivos USB directos
        let printerList = [];
        
        // 1. Intentar encontrar impresoras CUPS
        exec('lpstat -p -d 2>/dev/null | grep "printer" | awk \'{print $2}\' | tr -d \'":"\' | sort | uniq', 
        (error, stdout, stderr) => {
            if (!error && stdout.trim()) {
                printerList = printerList.concat(
                    stdout.split('\n').filter(line => line.trim().length > 0)
                );
                console.log(`[Bridge] Impresoras CUPS encontradas: ${printerList.join(', ')}`);
            }
            
            // 2. Buscar dispositivos USB en /dev/usb/
            console.log('[Bridge] Buscando dispositivos USB en /dev/usb/');
            exec('ls -la /dev/usb/lp* 2>/dev/null || ls -la /dev/lp* 2>/dev/null', 
            (err, out, stderr) => {
                if (!err && out) {
                    const devices = out.split('\n')
                        .filter(line => line.includes('usb') || line.includes('/dev/lp'))
                        .map(line => {
                            const match = line.match(/(\/(dev\/usb\/lp\d+|dev\/lp\d+))/);
                            return match ? match[1] : null;
                        })
                        .filter(d => d);
                    
                    printerList = printerList.concat(devices);
                    console.log(`[Bridge] Dispositivos USB encontrados: ${devices.join(', ')}`);
                }
                
                // 3. Retornar lista final sin duplicados
                const finalList = [...new Set(printerList)].filter(p => p);
                console.log(`[Bridge] Impresoras totales: ${finalList.join(', ')}`);
                res.json(finalList);
            });
        });
    } else if (SISTEMA === 'darwin') {
        // Para macOS
        let command = 'lpstat -p | grep printer | awk \'{print $2}\' | tr -d \'"\'';
        exec(command, (error, stdout, stderr) => {
            if (error) {
                console.error(`[Bridge] Error obteniendo impresoras: ${error.message}`);
                return res.json([]);
            }
            const printerList = stdout.split('\n').filter(line => line.trim().length > 0);
            console.log(`[Bridge] Impresoras encontradas: ${printerList.join(', ')}`);
            res.json(printerList);
        });
    } else {
        res.json([]);
    }
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

        if (SISTEMA === 'win32') {
            imprimirEnWindows(tempFilePath, nombreImpresora, res);
        } else if (SISTEMA === 'linux') {
            imprimirEnLinux(tempFilePath, nombreImpresora, res);
        } else if (SISTEMA === 'darwin') {
            imprimirEnMacOS(tempFilePath, nombreImpresora, res);
        } else {
            throw new Error(`Sistema operativo no soportado: ${SISTEMA}`);
        }

    } catch (e) {
        console.error("[Bridge] Error procesando buffer:", e);
        res.status(500).json({ error: e.message });
    }
});

// Función para imprimir en Windows
function imprimirEnWindows(tempFilePath, nombreImpresora, res) {
    console.log(`[Bridge-Windows] Intentando impresión...`);
    
    // Método 1: copy /b
    const copyCommand = `copy /b "${tempFilePath}" "\\\\localhost\\${nombreImpresora}"`;
    exec(copyCommand, (error, stdout, stderr) => {
        if (error) {
            console.log(`[Bridge-Windows] copy /b falló: ${error.message}`);
            
            // Método 2: print /D
            const printCommand = `print /D:"${nombreImpresora}" "${tempFilePath}"`;
            exec(printCommand, (error2, stdout2, stderr2) => {
                try { fs.unlinkSync(tempFilePath); } catch (e) { }
                
                if (error2) {
                    console.error(`[Bridge-Windows] print /D también falló`);
                    return res.json({
                        success: false,
                        error: `Ambos métodos fallaron`
                    });
                }
                
                console.log(`[Bridge-Windows] ✓ Impresión exitosa`);
                res.json(true);
            });
        } else {
            try { fs.unlinkSync(tempFilePath); } catch (e) { }
            console.log(`[Bridge-Windows] ✓ Impresión exitosa`);
            res.json(true);
        }
    });
}

// Función para imprimir en Linux (incluyendo Raspberry Pi)
function imprimirEnLinux(tempFilePath, nombreImpresora, res) {
    console.log(`[Bridge-Linux] Intentando impresión en impresora: ${nombreImpresora}`);
    
    // Verificar si es un dispositivo USB directo
    const esDispositivoUSB = nombreImpresora.startsWith('/dev/usb/') || nombreImpresora.startsWith('/dev/lp');
    
    if (esDispositivoUSB) {
        console.log(`[Bridge-Linux] Detectado dispositivo USB: ${nombreImpresora}`);
        imprimirEnDispositivoUSB(tempFilePath, nombreImpresora, res);
    } else {
        // Usar CUPS (lp o lpr)
        console.log(`[Bridge-Linux] Usando CUPS para impresora: ${nombreImpresora}`);
        imprimirConCUPS(tempFilePath, nombreImpresora, res);
    }
}

// Función para imprimir directamente en dispositivo USB
function imprimirEnDispositivoUSB(tempFilePath, dispositivo, res) {
    console.log(`[Bridge-Linux-USB] Escribiendo en dispositivo: ${dispositivo}`);
    
    // Verificar que el dispositivo existe y es accesible
    fs.access(dispositivo, fs.constants.W_OK, (err) => {
        if (err) {
            console.error(`[Bridge-Linux-USB] No se puede acceder al dispositivo ${dispositivo}: ${err.message}`);
            try { fs.unlinkSync(tempFilePath); } catch (e) { }
            
            // Intentar con sudo
            console.log(`[Bridge-Linux-USB] Intentando con sudo...`);
            const catCommand = `sudo cat "${tempFilePath}" > "${dispositivo}"`;
            exec(catCommand, (error, stdout, stderr) => {
                try { fs.unlinkSync(tempFilePath); } catch (e) { }
                
                if (error) {
                    console.error(`[Bridge-Linux-USB] Error con sudo: ${error.message}`);
                    return res.json({
                        success: false,
                        error: `No se puede acceder al dispositivo. Verifica permisos: chmod 666 ${dispositivo}`,
                        dispositivo: dispositivo
                    });
                }
                
                console.log(`[Bridge-Linux-USB] ✓ Impresión exitosa (con sudo)`);
                res.json({
                    success: true,
                    mensaje: 'Ticket enviado al dispositivo USB (requiere permisos de sudo)'
                });
            });
            return;
        }
        
        // El dispositivo es accesible, escribir directamente
        const readStream = fs.createReadStream(tempFilePath);
        const writeStream = fs.createWriteStream(dispositivo);
        
        readStream.on('error', (err) => {
            console.error(`[Bridge-Linux-USB] Error leyendo archivo: ${err.message}`);
            res.json({
                success: false,
                error: `Error leyendo archivo: ${err.message}`
            });
        });
        
        writeStream.on('error', (err) => {
            console.error(`[Bridge-Linux-USB] Error escribiendo en dispositivo: ${err.message}`);
            res.json({
                success: false,
                error: `Error escribiendo en dispositivo: ${err.message}. Verifica permisos: chmod 666 ${dispositivo}`
            });
        });
        
        writeStream.on('finish', () => {
            try { fs.unlinkSync(tempFilePath); } catch (e) { }
            console.log(`[Bridge-Linux-USB] ✓ Impresión exitosa`);
            res.json({
                success: true,
                mensaje: 'Ticket enviado exitosamente al dispositivo USB'
            });
        });
        
        readStream.pipe(writeStream);
    });
}

// Función para imprimir usando CUPS
function imprimirConCUPS(tempFilePath, nombreImpresora, res) {
    const lpCommand = `lp -d "${nombreImpresora}" "${tempFilePath}"`;
    
    exec(lpCommand, (error, stdout, stderr) => {
        try { fs.unlinkSync(tempFilePath); } catch (e) { }
        
        if (error) {
            console.error(`[Bridge-Linux-CUPS] Error en lp: ${error.message}`);
            
            // Intentar con lpr
            const lprCommand = `lpr -P "${nombreImpresora}" "${tempFilePath}"`;
            exec(lprCommand, (error2, stdout2, stderr2) => {
                try { fs.unlinkSync(tempFilePath); } catch (e) { }
                
                if (error2) {
                    console.error(`[Bridge-Linux-CUPS] lpr también falló: ${error2.message}`);
                    return res.json({
                        success: false,
                        error: `Error enviando a impresora: ${error.message}`
                    });
                }
                
                console.log(`[Bridge-Linux-CUPS] ✓ Impresión exitosa con lpr`);
                res.json({
                    success: true,
                    mensaje: 'Ticket enviado a través de CUPS'
                });
            });
            return;
        }
        
        console.log(`[Bridge-Linux-CUPS] ✓ Impresión exitosa`);
        res.json({
            success: true,
            mensaje: 'Ticket enviado a través de CUPS'
        });
    });
}

// Función para imprimir en macOS
function imprimirEnMacOS(tempFilePath, nombreImpresora, res) {
    console.log(`[Bridge-macOS] Intentando impresión...`);
    
    const lprCommand = `lpr -P "${nombreImpresora}" "${tempFilePath}"`;
    
    exec(lprCommand, (error, stdout, stderr) => {
        try { fs.unlinkSync(tempFilePath); } catch (e) { }
        
        if (error) {
            console.error(`[Bridge-macOS] Error: ${error.message}`);
            return res.json({
                success: false,
                error: `Error enviando a impresora: ${error.message}`
            });
        }
        
        console.log(`[Bridge-macOS] ✓ Impresión exitosa`);
        res.json(true);
    });
}

app.listen(PORT, () => {
    console.log(`[Bridge] Servidor de impresión corriendo en http://localhost:${PORT}`);
    console.log(`[Bridge] Sistema: ${SISTEMA}, Arquitectura: ${ARQUITECTURA}`);
});
