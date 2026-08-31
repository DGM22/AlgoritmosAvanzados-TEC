# Misión 3 — Merge Sort y Quick Sort (2 y 3 particiones)

Comparación de cuatro variantes de ordenamiento en Go, sobre los **mismos**
archivos de `datos/`:

| Programa | Algoritmo |
|---|---|
| `mergeSort/mergesort2.go` | Merge Sort clásico (divide a la **mitad**) |
| `mergeSort/mergesort3.go` | Merge Sort en **tercios** (luego dos merges binarios) |
| `quickSort/quicksort2.go` | Quick Sort con **1 pivote** (último elemento, Lomuto) |
| `quickSort/quicksort3.go` | Quick Sort con **2 pivotes** (primero y último) |

Cada corrida reporta comparaciones, movimientos, intercambios y tiempo.
Merge no intercambia en sitio: la columna de intercambios sale en **0**.

Los códigos de los algoritmos se documentaron con comentarios estilo JSDoc
(Airbnb) y encabezado de generación con inteligencia artificial OpenAI.

## Estructura

```
mision-3/
├── datos/                 # .txt de prueba (ya generados; no hace falta volver a crearlos)
├── mergeSort/
│   ├── mergesort2.go
│   ├── mergesort3.go
│   └── io.go             # lectura de archivos (hacer go run junto con el .go del algoritmo)
├── quickSort/
│   ├── quicksort2.go
│   ├── quicksort3.go
│   └── io.go
├── correr_todo.sh        # los 4 algoritmos, todos los tamaños
└── correr_50k.sh         # los 4 algoritmos, solo n = 50 000
```

## Requisitos

Go instalado (`go version`). Desde la carpeta `mision-3`.

## Cómo correr

Los `.txt` **ya existen** en `datos/`. Los scripts **no** los regeneran.

### Todos los conjuntos

Desde `mision-3`:

```bash
chmod +x correr_todo.sh   # solo la primera vez
./correr_todo.sh
```

Corre merge2, merge3, quick2 y quick3 sobre **cada** archivo de `datos/`
(100, 1 000, 8 000, 10 000, 50 000, y los distintos órdenes).

### Solo 50 000 datos

```bash
chmod +x correr_50k.sh    # solo la primera vez
./correr_50k.sh
```

Solo usa:

- `n050000_aleatorio.txt`
- `n050000_casi_ordenado.txt`
- `n050000_duplicados.txt`

No hay invertido de 50 000: Quick Sort con pivote último/primero-último
queda `O(n²)` y esa corrida sería demasiado pesada.

### Un algoritmo a mano

Hay que pasar también `io.go` (si no, falla el `go run`):

```bash
go run ./mergeSort/mergesort2.go ./mergeSort/io.go
go run ./mergeSort/mergesort3.go ./mergeSort/io.go
go run ./quickSort/quicksort2.go ./quickSort/io.go
go run ./quickSort/quicksort3.go ./quickSort/io.go
```

Solo un tamaño, con `-n`:

```bash
go run ./mergeSort/mergesort2.go ./mergeSort/io.go -n 1000
go run ./quickSort/quicksort2.go ./quickSort/io.go -n 50000
```

`-n 0` (o sin flag) = todos los `.txt`.

## Qué hay en `datos/`

Mismos números para los cuatro programas. Órdenes:

| Sufijo | Qué es |
|---|---|
| `aleatorio` | permutación revuelta de 1..n |
| `invertido` | n, n-1, …, 1 (peor caso típico de Quick) |
| `casi_ordenado` | 1..n con pocos intercambios |
| `duplicados` | valores 0–9 (muchos repetidos) |

Tamaños: 100, 1 000, 10 000 y 50 000, salvo invertido (hasta 10 000, más el
archivo extra `n008000_invertido.txt`).

## Cómo leer la tabla

| Columna | Significado |
|---|---|
| `n` | cantidad de enteros |
| `comparaciones` | veces que se compararon dos datos |
| `movimientos` | copias de un valor (en Quick, cada swap cuenta 3) |
| `intercambios` | `swap` reales; en Merge siempre 0 |
| `tiempo` | reloj; útil sobre todo en 10k y 50k |
| `ok` | si el arreglo quedó ordenado |

## Complejidad (estas implementaciones)

| Algoritmo | Tiempo promedio | Peor caso | Espacio extra |
|---|---|---|---|
| Merge 2 | `Θ(n log n)` | `Θ(n log n)` | `Θ(n)` |
| Merge 3 | `Θ(n log n)` | `Θ(n log n)` | `Θ(n)` |
| Quick 2 | `Θ(n log n)` | `Θ(n²)` | pila `O(log n)` / `O(n)` |
| Quick 3 | `Θ(n log n)` | `Θ(n²)` | pila `O(log n)` / `O(n)` |
