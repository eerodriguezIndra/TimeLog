# TimeLog

[![Build & Release](https://github.com/eerodriguezIndra/TimeLog/actions/workflows/release.yml/badge.svg)](https://github.com/eerodriguezIndra/TimeLog/actions/workflows/release.yml)
[![Release](https://img.shields.io/github/v/release/eerodriguezIndra/TimeLog?include_prereleases&label=nightly)](https://github.com/eerodriguezIndra/TimeLog/releases/tag/nightly)
[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)](https://go.dev/)

Aplicación de escritorio multiplataforma para **registro horario por interrupción**: cada cierto tiempo (configurable) aparece un panel preguntando *¿qué estás haciendo, para qué cliente y de qué tipo?*, y guarda la respuesta en un CSV.

Pensada para consultores y técnicos que rebotan entre tareas y necesitan reconstruir el día sin tener que recordar todo al final.

---

## Descarga rápida

Última build automática del commit más reciente en `main` — sin instalación, portable:

| Plataforma | Archivo | Descarga directa |
|---|---|---|
| **Windows x64** | `TimeLog-windows-amd64.exe` | [⬇ Descargar](https://github.com/eerodriguezIndra/TimeLog/releases/download/nightly/TimeLog-windows-amd64.exe) |
| **macOS Apple Silicon** (M1/M2/M3/M4) | `TimeLog-darwin-arm64.zip` | [⬇ Descargar](https://github.com/eerodriguezIndra/TimeLog/releases/download/nightly/TimeLog-darwin-arm64.zip) |
| **macOS Intel** | `TimeLog-darwin-amd64.zip` | [⬇ Descargar](https://github.com/eerodriguezIndra/TimeLog/releases/download/nightly/TimeLog-darwin-amd64.zip) |

¿No sabes qué Mac tienes? **Menú Apple → Acerca de este Mac** — si dice "Chip Apple M…" es Apple Silicon, si dice "Intel" es Intel.

Ver todas las builds y checksums en la [página del release nightly](https://github.com/eerodriguezIndra/TimeLog/releases/tag/nightly).

---

## Características

- **Sin botón de cierre real** — la X de la ventana solo oculta; la app sigue viva en el system tray hasta que tú la cierres desde el menú.
- **Always-on-top en Windows** — el icono flotante queda por encima de cualquier ventana, anclado en la esquina superior derecha.
- **System tray + widget flotante** — menubar en macOS, system tray en Windows/Linux, más una mini-ventana arrastrable con el icono.
- **Prompt periódico con auto-grow** — si no respondes, la ventana va creciendo hasta ocupar ~1/3 de la pantalla. No la podrás ignorar.
- **Notificación nativa del SO** al disparar el prompt (toast en Windows, banner en macOS).
- **CSV append por día** — todas las entradas se concatenan al mismo archivo, columnas: `fecha, hora, cliente, actividad, descripcion`.
- **Configurable y persistente**:
  - Intervalo entre prompts (minutos)
  - Ruta del archivo CSV (con file picker)
  - Autoarranque al iniciar sesión
  - Tipos de actividad y clientes recientes
- **Autostart cross-platform** — LaunchAgent en macOS, `.desktop` en Linux (`~/.config/autostart/`), entrada en Registry `HKCU\…\Run` en Windows.
- **Tema oscuro moderno** con acento azul.
- **Binario portable**: un solo `.exe` autocontenido en Windows (sin instalador, sin DLLs externas).

---

## Instalación

### Windows (portable)

Descarga **[TimeLog-windows-amd64.exe](https://github.com/eerodriguezIndra/TimeLog/releases/download/nightly/TimeLog-windows-amd64.exe)** y ejecútalo directamente — no necesita instalación. La primera vez abrirá la ventana de configuración para que escojas el intervalo, dónde guardar el CSV y si quieres que arranque al iniciar sesión.

> **SmartScreen**: la primera vez Windows puede mostrar "Windows protegió tu PC" (el `.exe` no está firmado). Clic en **Más información → Ejecutar de todos modos**.

### macOS

Descarga el `.zip` según el chip de tu Mac:

- **Apple Silicon (M1/M2/M3/M4)** → [TimeLog-darwin-arm64.zip](https://github.com/eerodriguezIndra/TimeLog/releases/download/nightly/TimeLog-darwin-arm64.zip)
- **Intel** → [TimeLog-darwin-amd64.zip](https://github.com/eerodriguezIndra/TimeLog/releases/download/nightly/TimeLog-darwin-amd64.zip)

Descomprime (doble clic) y obtienes `TimeLog.app`. Muévelo a `/Applications` (opcional pero recomendado).

**La primera vez, Gatekeeper la bloqueará** porque la app no está firmada/notarizada. Tienes dos opciones:

**Opción A — desde Finder (recomendada):**
1. Clic derecho (o Ctrl + clic) sobre `TimeLog.app` → **Abrir**.
2. Aparece un diálogo "¿Estás seguro?" → **Abrir**.
3. Próximas veces se abre normalmente con doble clic.

**Opción B — desde Terminal:**
```bash
xattr -dr com.apple.quarantine /Applications/TimeLog.app
open /Applications/TimeLog.app
```

El icono quedará en la **menubar** (arriba a la derecha). Clic en él para "+ Nueva tarea", Configuración o Salir.

### Linux

Compila desde código. Necesitarás dependencias gráficas:

```bash
sudo apt-get install -y libgl1-mesa-dev xorg-dev libx11-dev libxcursor-dev \
                        libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev
```

---

## Uso

1. **Primer arranque**: se abre la ventana de configuración. Define cada cuántos minutos quieres el prompt y dónde guardar el `timelog.csv`.
2. **El icono queda en el system tray** (menubar Mac / area de notificación Windows / tray Linux). Click derecho para:
   - **Registrar ahora** — fuerza el prompt manualmente.
   - **Configuración…** — cambia intervalo, ruta CSV, autoarranque.
   - **Salir** — cierra la app de verdad.
3. **Widget flotante** — pequeña ventana siempre visible con el icono. En Windows queda anclada en la esquina superior derecha por encima de todo.
4. **Cuando aparezca el prompt**:
   - Escribe el cliente (los anteriores aparecen como sugerencia)
   - Escoge el tipo de actividad
   - Describe brevemente qué estás haciendo
   - **Guardar** lo añade al CSV con timestamp
   - **Recordar luego** lo oculta sin guardar (volverá a aparecer en el próximo ciclo)
5. **Si ignoras el prompt**, la ventana va creciendo cada 30 segundos hasta ocupar ~1/3 de la pantalla.

### Formato del CSV

```csv
fecha,hora,cliente,actividad,descripcion
2026-05-16,09:15:32,Acme Corp,Reunión,Daily standup del equipo de backend
2026-05-16,10:18:04,Acme Corp,Documentación,Actualizando ADR del módulo de pagos
```

Tipos de actividad predefinidos:
- Reunión
- Plan de trabajo
- Documentación
- Ejecución de operación
- Comité de cambio
- Respuesta de correo

---

## Configuración

Archivo: `~/Library/Application Support/TimeLog/config.json` (macOS) · `%AppData%\TimeLog\config.json` (Windows) · `~/.config/TimeLog/config.json` (Linux).

```json
{
  "interval_minutes": 60,
  "csv_path": "/Users/tu/Documents/timelog.csv",
  "autostart": false,
  "clients": ["Acme Corp", "Beta Ltd"],
  "activities": [
    "Reunión",
    "Plan de trabajo",
    "Documentación",
    "Ejecución de operación",
    "Comité de cambio",
    "Respuesta de correo"
  ]
}
```

Puedes editarlo a mano o desde la ventana de Configuración. Cambiar el intervalo desde la UI lo aplica en caliente, sin reiniciar la app.

---

## Compilación desde código

Requiere **Go 1.22+** y, en Linux/macOS, las dependencias gráficas listadas arriba.

```bash
git clone https://github.com/eerodriguezIndra/TimeLog.git
cd TimeLog
go mod download

# Build local del SO actual
go build -o timelog .

# Windows portable (con icono embebido, sin consola)
go install fyne.io/tools/cmd/fyne@latest
fyne package --os windows --icon icono.png --name TimeLog --app-id com.edwin.timelog

# macOS (.app)
fyne package --os darwin  --icon icono.png --name TimeLog --app-id com.edwin.timelog
```

### Estructura

```
TimeLog/
├── main.go                        # Entry point: wiring de config + scheduler + UI + tray
├── icono.png  / icono.svg         # Icono de la aplicación (.exe / tray)
├── avatar.png / avatar.svg        # Logo mostrado en ventanas
├── internal/
│   ├── config/        # JSON persistente, snapshot lock-free
│   ├── storage/       # Append a CSV con encabezado automático
│   ├── scheduler/     # Timer reconfigurable en caliente
│   ├── autostart/     # LaunchAgent / .desktop / Registry Run
│   └── ui/
│       ├── theme.go              # Tema oscuro moderno
│       ├── tray.go               # System tray (menubar/tray)
│       ├── floating.go           # Widget flotante con icono
│       ├── prompt.go             # Ventana de captura con auto-grow
│       ├── settings.go           # Ventana de configuración
│       ├── topmost_windows.go    # user32.SetWindowPos(HWND_TOPMOST)
│       ├── topmost_other.go      # no-op
│       ├── position_windows.go   # GetSystemMetrics → esquina sup-der
│       └── position_other.go     # no-op
└── .github/workflows/release.yml  # CI: build Windows + macOS Intel/ARM → Release
```

---

## Distribución automática

Cada `git push` a `main` dispara [`release.yml`](.github/workflows/release.yml), que en paralelo:

1. **Windows** (`windows-latest`): `fyne package` produce `TimeLog.exe` con icono embebido, subsistema GUI y sin consola.
2. **macOS Intel** (`macos-13`): `fyne package -os darwin` produce `TimeLog.app`, comprimido con `ditto` a `.zip`.
3. **macOS Apple Silicon** (`macos-latest`): igual al anterior pero `arm64`.

Todos los artefactos (más sus `.sha256`) se publican en la misma **pre-release rolling** `nightly` en GitHub Releases.

Si haces `git tag v1.2.3 && git push --tags`, en lugar de actualizar `nightly` crea una release etiquetada `TimeLog v1.2.3` con notas autogeneradas.

---

## Limitaciones conocidas

- **macOS / Linux no tienen always-on-top a nivel SO** desde código Go portable. En macOS el menubar suple esa función; en Linux depende del gestor de ventanas. Si necesitas always-on-top estricto fuera de Windows, abre un issue.
- **El cálculo de "1/3 de pantalla"** asume monitor 1920×1080. Si tu pantalla es muy distinta, el crecimiento puede quedar grande o pequeño.
- **No hay firma de código** — Windows mostrará SmartScreen la primera vez (clic en "Más información" → "Ejecutar de todos modos"). macOS bloquearía un binario sin notarizar.

---

## Licencia

MIT. Ver [LICENSE](LICENSE) si lo añades.
