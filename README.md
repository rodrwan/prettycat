# prettycat

`prettycat` es un CLI en Go tipo `cat`, pero con salida bonita para terminal:

- Detecta Markdown automáticamente y lo renderiza con estilo.
- Resalta sintaxis en archivos de código por extensión.
- Tiene fallback a texto plano para archivos no reconocidos.
- Incluye pager interactivo con navegación por teclas.

## Características

- Entrada por archivo(s): `prettycat archivo.txt`
- Entrada por `stdin`: `cat archivo.md | prettycat`
- Múltiples archivos con encabezados visuales por sección
- Política Unix de errores: continúa en fallos parciales y retorna `exit code 1` si hubo errores
- Soporte `--no-color`, `--help`, `--version`

## Tipos de archivos soportados

- Markdown (render semántico): `.md`, `.markdown`, `.mdown`, `.mkd`
- Código con resaltado por keywords: `.go`, `.js`, `.ts`, `.py`, `.rb`, `.java`
- Cualquier otro formato: fallback a texto plano

## Paleta de colores actual (ANSI 256)

`prettycat` usa colores ANSI de 256 tonos. Referencia actual:

- `212` (rosa): bullets Markdown, encabezados de archivo, prompt de búsqueda
- `159` (cian claro): títulos Markdown `#`
- `117` (azul claro): títulos Markdown `##`
- `111` (turquesa): títulos Markdown `###`
- `179` (ámbar): bloques de código Markdown e inline code
- `81` (azul brillante): keywords en archivos de código (`go/js/ts/py/rb/java`)
- `216` (durazno): strings en código
- `244` (gris): status del pager, quotes Markdown y comentarios
- `240` (gris oscuro): separadores visuales entre archivos
- `250` (gris claro): fallback para código genérico

## Estructura del repositorio

- `cmd/prettycat/main.go`: entrypoint del CLI
- `internal/app`: orquestación principal
- `internal/input`: carga de archivos y stdin
- `internal/render`: render Markdown/código/plano
- `internal/pager`: pager interactivo
- `internal/style`: helpers de estilo ANSI
- `testdata/`: archivos de ejemplo
- `Makefile`: comandos de desarrollo

## Instalación

### Requisitos

- Go 1.22+

### Build local

```bash
make build
```

Genera el binario en `bin/prettycat`.

### Instalación en GOPATH/bin

```bash
make install
```

## Uso

```bash
prettycat [flags] [file ...]
```

### Flags

- `--no-color`: desactiva colores ANSI
- `--version`: muestra versión
- `--help`: ayuda

### Ejemplos

```bash
# Markdown bonito
prettycat README.md

# Código con highlight
prettycat cmd/prettycat/main.go

# Múltiples archivos
prettycat README.md Makefile

# Desde stdin
cat testdata/sample.md | prettycat

# Sin color
prettycat --no-color testdata/sample.go
```

## Controles del pager interactivo

Cuando la salida va a una TTY, se activa el pager:

- `j` / `k` o `↑` / `↓`: mover línea
- `f` / `b` / `space`: avanzar o retroceder página
- `g` / `G`: inicio / fin
- `/`: buscar (Enter confirma, Esc cancela)
- `n` / `N`: siguiente/anterior match
- `q`: salir

## Desarrollo

Comandos útiles:

```bash
make help       # lista de targets
make test       # ejecutar tests
make testv      # tests en modo verbose
make fmt        # formatear Go
make fmt-check  # validar formato
make tidy       # limpiar/sincronizar dependencias
make check      # fmt-check + test + build
make clean      # limpiar artefactos
```

## Testing

```bash
make test
```

El proyecto usa `go test` estándar con pruebas unitarias en `internal/app` y `internal/render`.

## Estado actual

Proyecto funcional para uso local y desarrollo iterativo del render/pager.
