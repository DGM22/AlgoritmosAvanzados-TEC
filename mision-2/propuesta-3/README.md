# Propuesta de solución #3 — Calificación más común

Implementación en C++ y Go de la Propuesta #3: una tabla hash de 10
posiciones, donde la posición de una calificación es el residuo de
dividirla entre 10 (`calificacion % 10`). Cada posición guarda una lista
encadenada de las calificaciones distintas que caen ahí junto con su
contador; al leer una calificación se busca en su lista: si ya existe se
incrementa su contador, si no se agrega un nodo nuevo.

Mismo diseño en ambos lenguajes: `HashTable` (C++, con nodos gestionados
por `unique_ptr`) y su equivalente en Go (con punteros y garbage
collector), además de un `FileReader` y un `Analyzer` compartidos con las
otras propuestas.

## Algoritmo

1. Calcular el índice de una calificación como `calificacion % 10`.
2. Insertar: buscar la calificación en la lista de esa posición; si existe,
   sumar 1 a su contador; si no, crear un nodo nuevo con contador 1.
3. Al terminar de leer el archivo, recorrer las 10 listas completas
   comparando contadores; en caso de empate (mismo contador) se elige la
   calificación mayor.
4. Porcentaje = (frecuencia de la moda / total de datos) × 100.

## Complejidad

Leer e insertar es O(n) en tiempo: cada posición de la tabla tiene en
promedio pocas calificaciones distintas, así que cada inserción es
prácticamente O(1). Encontrar la moda recorre las 10 listas, con a lo más
100 nodos en total, así que es O(k) con k ≤ 100. El programa completo es
O(n) en tiempo. En espacio solo se guardan pares (calificación, contador)
—máximo 100 nodos—, nunca los datos crudos, por lo que el espacio extra es
O(1), igual que la Propuesta #1.

## Casos especiales

- Archivo no encontrado: se avisa y se vuelve a pedir un archivo, sin
  abortar el programa.
- Líneas vacías: se ignoran.
- Líneas no numéricas o fuera de [1, 100]: se cuentan como "líneas
  ignoradas" y se reporta cuántas hubo.
- Archivo sin calificaciones válidas: se informa y no se acumula al
  resultado global.
- Renglón vacío: termina la captura de archivos.

## Cómo ejecutar

Ver el [README de la raíz](../README.md) para instrucciones completas en
macOS y Windows. En resumen:

```bash
cd propuesta-3/cpp && cmake -S . -B build && cmake --build build && ./build/grade_mode_finder
cd propuesta-3/go && go run .
```

## Comparación con las otras propuestas

Logra el mismo O(n) en tiempo y espacio constante que la Propuesta #1, pero
con más código y una estructura más compleja (listas encadenadas) para
resolver colisiones. Como el dominio de calificaciones es pequeño y
conocido (1-100), la tabla hash no aporta una ventaja real sobre el arreglo
directo de la Propuesta #1; sirve más como ejercicio para practicar el
manejo de colisiones con encadenamiento.
