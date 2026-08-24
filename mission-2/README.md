# Mission 2 — ¿Cuál es la calificación más común?

Tres propuestas de solución al mismo problema (encontrar la calificación más
común —moda— de un grupo de alumnos, y su porcentaje), cada una implementada
dos veces: en **C++** y en **Go**, para comparar el mismo algoritmo en ambos
lenguajes.

| Propuesta | Enfoque |
|---|---|
| `propuesta-1` | Arreglo de contadores indexado por calificación (1-100) |
| `propuesta-2` | Vector dinámico + ordenamiento + conteo de corridas |
| `propuesta-3` | Tabla hash de 10 posiciones con encadenamiento externo |

Cada carpeta de propuesta tiene su propio README explicando qué se hizo. Este
README solo cubre cómo compilar y ejecutar los programas en macOS y Windows.

## Estructura

```
mission-2/
├── propuesta-1/
│   ├── cpp/
│   └── go/
├── propuesta-2/
│   ├── cpp/
│   └── go/
├── propuesta-3/
│   ├── cpp/
│   └── go/
└── calificaciones_caso_pruebas/   # Archivos .txt de prueba
```

Las tres propuestas comparten la misma interfaz: piden nombres de archivo por
teclado hasta un renglón vacío, y al final muestran el resultado acumulado de
todos los archivos procesados.

## Requisitos y verificación de versiones

### macOS

**Compilador C++** (Clang, viene con Xcode Command Line Tools):
```bash
xcode-select --install   # solo si no está instalado
clang++ --version
```

**CMake** (opcional, recomendado para compilar C++):
```bash
brew install cmake       # si falta
cmake --version
```

**Go**:
```bash
brew install go          # si falta
go version
```

### Windows

**Compilador C++**: MinGW-w64 (`g++`) o Visual Studio Build Tools (MSVC):
```powershell
g++ --version
:: o, dentro de un "Developer Command Prompt for VS":
cl
```

**CMake**:
```powershell
cmake --version
```

**Go**:
```powershell
go version
```

Si falta alguna herramienta, se puede instalar con
[Chocolatey](https://chocolatey.org/):
```powershell
choco install mingw cmake golang
```
o descargarla manualmente: [Go](https://go.dev/dl/),
[CMake](https://cmake.org/download/), [MSYS2/MinGW-w64](https://www.msys2.org/),
o Visual Studio con el workload "Desktop development with C++".

## Cómo compilar y ejecutar

Los tres proyectos de C++ generan el mismo ejecutable (`grade_mode_finder`) y
los tres de Go se ejecutan igual; solo cambia la carpeta (`propuesta-1`,
`propuesta-2` o `propuesta-3`). Reemplazar `<N>` por `1`, `2` o `3` según la
propuesta que se quiera correr.

### C++ — macOS / Linux

Con CMake:
```bash
cd propuesta-<N>/cpp
cmake -S . -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build
./build/grade_mode_finder
```

Sin CMake, directo con el compilador:
```bash
cd propuesta-<N>/cpp
clang++ -std=c++17 -Wall -Wextra -Wpedantic -O2 main.cpp -o grade_mode_finder
./grade_mode_finder
```

### C++ — Windows

Con CMake (usa el generador de Visual Studio o de MinGW, según lo que tengas
instalado):
```powershell
cd propuesta-<N>\cpp
cmake -S . -B build
cmake --build build --config Release
```
El ejecutable queda en `build\Release\grade_mode_finder.exe` (generador de
Visual Studio) o en `build\grade_mode_finder.exe` (generador MinGW Makefiles).

Sin CMake, con MinGW-w64:
```powershell
cd propuesta-<N>\cpp
g++ -std=c++17 -Wall -Wextra -Wpedantic -O2 main.cpp -o grade_mode_finder.exe
.\grade_mode_finder.exe
```

### Go — macOS / Linux

```bash
cd propuesta-<N>/go
go run .
```
o compilando el binario:
```bash
cd propuesta-<N>/go
go build -o grade_mode_finder .
./grade_mode_finder
```

### Go — Windows

```powershell
cd propuesta-<N>\go
go run .
```
o compilando el binario:
```powershell
cd propuesta-<N>\go
go build -o grade_mode_finder.exe .
.\grade_mode_finder.exe
```

## Archivos de prueba

En `calificaciones_caso_pruebas/` hay 6 archivos `.txt` con calificaciones
(de 100 a 500,000 datos) para probar los programas. `Empates/` y `Moda/`
contienen casos adicionales (empate controlado y moda clara).
