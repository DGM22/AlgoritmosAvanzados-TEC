# Propuesta de solución #1 — Calificación más común

Implementación en **C++** y **Go** de la Propuesta #1 del caso "¿Cuál es la
calificación más común entre un grupo de alumnos?": un **arreglo de
contadores** de tamaño fijo, donde el índice del arreglo es directamente el
valor de la calificación (1-100). Al leer cada dato del archivo se
incrementa el contador correspondiente; al final se recorre el arreglo para
encontrar el máximo (la moda) y se calcula su porcentaje.

Se implementó dos veces (mismo diseño, mismo algoritmo) para comparar cómo
se expresa el mismo paradigma orientado a objetos en un lenguaje con clases
clásicas (C++) y en un lenguaje sin clases pero con structs + interfaces
(Go).

## Estructura del proyecto

```
propuesta-1/
├── cpp/
│   ├── include/            # Headers (.hpp) — declaraciones de las clases
│   ├── src/                 # Implementaciones (.cpp) + main
│   └── CMakeLists.txt
└── go/
    ├── grades/               # Paquete con la lógica de dominio
    └── main.go               # Punto de entrada
```

## Diseño orientado a objetos

Ambas versiones dividen el problema en las mismas cuatro responsabilidades
(principio de responsabilidad única — SRP):

| Responsabilidad | C++ (clase) | Go (struct + métodos) |
|---|---|---|
| Guardar y consultar los contadores, calcular la moda | `GradeFrequencyTable` | `grades.FrequencyTable` |
| Leer un archivo y alimentar la tabla de frecuencias | `GradeFileReader` | `grades.FileReader` |
| Dar formato al resultado para el usuario | `GradeAnalyzer` | `grades.Analyzer` |
| Orquestar el flujo (pedir archivos, acumular resultado global) | `GradeModeApplication` (en `main.cpp`) | `gradeModeApp` (en `main.go`) |

Nota sobre POO en Go: Go no tiene clases ni herencia. El equivalente
idiomático es un `struct` con *métodos con receptor* (`func (t *FrequencyTable) Increment(...)`),
que da el mismo encapsulamiento y una API igual de clara; el polimorfismo,
cuando se necesita, se logra con `interface` en vez de con jerarquías de
herencia. Por eso a la versión Go se le llama "orientada a objetos" en un
sentido de composición, no de herencia clásica.

## Algoritmo (Propuesta #1)

1. Crear un arreglo `contadores[1..100]` inicializado en cero.
2. Por cada calificación `g` leída del archivo: `contadores[g]++`.
3. Al terminar de leer, recorrer `contadores` de **100 hacia 1** buscando el
   valor máximo. Se usa comparación estricta (`>`, no `>=`), así el primer
   máximo encontrado en ese recorrido —que corresponde a la calificación
   más alta— gana cualquier empate, como exige el enunciado.
4. Porcentaje = `(contador_moda / total_de_datos) * 100`.

## Análisis de complejidad

Sea **n** = número de calificaciones en el archivo (hasta 500,000) y
**k** = 100 (tamaño fijo de la escala de calificaciones).

| Operación | Tiempo | Espacio extra |
|---|---|---|
| Leer el archivo e incrementar contadores | O(n) | O(1) (se lee línea por línea, nunca se carga el archivo completo en memoria) |
| Encontrar la moda (recorrer el arreglo) | O(k) = O(1)* | O(1) |
| Fusionar resultados de varios archivos (resultado global) | O(k) por archivo | O(1) |
| **Total del programa** | **O(n)** | **O(1)** |

\* Es la ventaja clave de esta propuesta: como la escala de calificaciones
está acotada (1-100), buscar la moda cuesta siempre lo mismo sin importar
si el archivo tiene 100 o 500,000 datos. El costo real del programa está
dominado por la **lectura** del archivo (O(n)), no por el cálculo de la
moda. La memoria usada por los contadores también es constante: 101 enteros,
sin importar n.

Comparado con otras propuestas típicas para este problema (p. ej. ordenar
las calificaciones primero, O(n log n)), la Propuesta #1 es óptima en
tiempo (no se puede hacer mejor que O(n), hay que leer todo el archivo al
menos una vez) y en espacio, precisamente porque aprovecha que el dominio
de valores (1-100) es pequeño y conocido de antemano.

## Manejo de casos especiales

- **Archivo no encontrado**: se informa el error y se le vuelve a pedir un
  archivo al usuario (no se aborta el programa).
- **Líneas vacías**: se ignoran silenciosamente.
- **Líneas no numéricas o fuera de rango [1, 100]**: se cuentan como
  "líneas ignoradas" y se reporta cuántas hubo, sin detener el procesamiento.
- **Archivo sin ninguna calificación válida**: se informa y no se acumula
  al resultado global.
- **Empate**: se resuelve a favor de la calificación mayor (ver algoritmo).
- **Fin de la captura de archivos**: renglón vacío (Enter sin nombre de archivo).

## Cómo compilar y ejecutar

### C++

Con CMake:

```bash
cd propuesta-1/cpp
cmake -S . -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build
./build/grade_mode_finder
```

Sin CMake (directo con el compilador):

```bash
cd propuesta-1/cpp
clang++ -std=c++17 -Wall -Wextra -Wpedantic -O2 -Iinclude src/*.cpp -o grade_mode_finder
./grade_mode_finder
```

### Go

```bash
cd propuesta-1/go
go build -o grade_mode_finder .
./grade_mode_finder
```

O directamente, sin generar el binario:

```bash
cd propuesta-1/go
go run .
```

### Uso — modo interactivo (idéntico en ambas versiones)

```
=== Calificación más común (Propuesta #1: arreglo de contadores) ===

Ingrese el nombre del archivo (Enter para terminar): ../../calificaciones_caso_pruebas/grupo_pequeno_100.txt
Archivo: ../../calificaciones_caso_pruebas/grupo_pequeno_100.txt
Cantidad de datos: 100
Calificación más común: 85
Frecuencia: 18
Porcentaje: 18.0000%
Tiempo de procesamiento: 38µs

Ingrese el nombre del archivo (Enter para terminar):
RESULTADO ACUMULADO DE TODOS LOS ARCHIVOS
-----------------------------------------
Cantidad total de datos: 100
Calificación más común global: 85
Frecuencia global: 18
Porcentaje global: 18.0000%
Tiempo del cálculo de la moda global: 209ns
Tiempo total acumulado (todos los archivos): 38µs
```

### Uso — modo automático `-dir` (solo Go, pensado para pruebas/benchmark)

En vez de escribir cada nombre de archivo a mano, la versión Go acepta una
bandera `-dir` que procesa automáticamente, en orden alfabético, todos los
`.txt` de esa carpeta (ignorando `LEEME.txt` y `resultados_esperados.txt`,
que no son datos de calificaciones). Es solo un atajo de conveniencia para
probar rápido durante el desarrollo — el programa entregable sigue siendo el
modo interactivo de arriba, que es el que pide el enunciado.

```bash
cd propuesta-1/go
go run . -dir ../../calificaciones_caso_pruebas
```

### Medición de tiempo (Go)

La versión en Go mide internamente (con `time.Now()` / `time.Since()`, sin
depender del `time` del shell) cuánto tarda **solo el algoritmo** —leer el
archivo, contar y calcular la moda—, sin contar lo que tarda imprimir en
pantalla:

- **Por archivo**: línea `Tiempo de procesamiento: ...` — domina la lectura
  del archivo, es la parte O(n).
- **Moda global**: línea `Tiempo del cálculo de la moda global: ...` — solo
  recorrer 100 posiciones del arreglo acumulado; suele medirse en
  nanosegundos, evidencia visible de que esta operación es O(k) constante,
  sin importar cuántos archivos o datos se hayan procesado.
- **Total acumulado**: línea `Tiempo total acumulado (todos los archivos): ...`
  — suma de los tiempos de cada archivo individual.

## Validación contra los casos de prueba

Se ejecutaron ambas versiones contra los 6 archivos en
`calificaciones_caso_pruebas/` y se compararon contra
`resultados_esperados.txt`. **Los dos programas (C++ y Go) reproducen
exactamente** todos los valores esperados, incluyendo el caso de empate
(90 vs. 95, frecuencia 8 cada uno → se elige 95) y el resultado global
acumulado de los 611,124 datos:

| Archivo | Moda esperada | Frecuencia | Porcentaje | C++ | Go |
|---|---|---|---|---|---|
| grupo_pequeno_100.txt | 85 | 18 | 18.0000% | ✅ | ✅ |
| grupo_1000.txt | 85 | 140 | 14.0000% | ✅ | ✅ |
| grupo_10000.txt | 85 | 1,200 | 12.0000% | ✅ | ✅ |
| grupo_100000.txt | 85 | 10,000 | 10.0000% | ✅ | ✅ |
| grupo_500000.txt | 85 | 40,000 | 8.0000% | ✅ | ✅ |
| empate_controlado.txt | 95 | 8 | 33.3333% | ✅ | ✅ |
| **Global (611,124 datos)** | **85** | **51,358** | **8.4039%** | ✅ | ✅ |

## Comparación C++ vs. Go

| Aspecto | C++ | Go |
|---|---|---|
| Paradigma OOP | Clases con encapsulamiento explícito (`private`/`public`), sin herencia en este diseño | Structs + métodos con receptor; sin herencia, composición + `interface` para polimorfismo |
| Manejo de errores | Excepciones para invariantes internas (`std::out_of_range`) + `struct` de resultado (`ReadResult`) para errores esperables (evita el costo y la ambigüedad de excepciones para control de flujo normal) | Errores como valores de retorno (`error`), sin excepciones; es el estilo idiomático de Go en todo el lenguaje |
| Gestión de memoria | Manual/RAII: `std::array` en la pila, sin punteros crudos ni `new`/`delete` | Automática con *garbage collector*; los `*FrequencyTable` se liberan solos cuando ya no se usan |
| Rendimiento de E/S | `std::ifstream` + `std::getline`, buffer manejado por la STL | `bufio.Scanner`, buffer configurable explícitamente |
| Formato de números | Se implementó manualmente (separador de miles, 4 decimales) para no depender del locale del sistema | Igual: implementado manualmente con `strconv`, mismo motivo |
| Curva de aprendizaje | Más verboso (headers separados, gestión explícita de tipos) | Más compacto; el manejo explícito de `error` en cada llamada es distinto a lo que se usa en C++/Java pero se vuelve natural rápido |

## Siguientes pasos sugeridos

- Agregar pruebas unitarias (Google Test en C++, `testing` estándar en Go)
  que verifiquen `FrequencyTable` de forma aislada (sin necesidad de
  archivos), incluyendo el caso de empate y el caso de archivo vacío.
- Cuando se aborden las Propuestas #2 y #3, agregar una carpeta hermana
  `propuesta-2/` y `propuesta-3/` con la misma estructura, para poder medir
  tiempos de ejecución reales (`time ./grade_mode_finder`) y comparar contra
  esta primera propuesta.
