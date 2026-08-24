# Propuesta de solución #2 — Calificación más común

Implementación en C++ y Go de la Propuesta #2: se guardan todas las
calificaciones leídas en un vector dinámico, se ordenan de forma ascendente
y luego se recorre el vector una sola vez contando "corridas" de valores
iguales para encontrar cuál se repite más veces (la moda).

Mismo diseño en ambos lenguajes: una estructura que guarda y ordena las
calificaciones (`GradeVector` en C++, su equivalente en Go), un `FileReader`
para leer el archivo y un `Analyzer` que da formato al resultado.

## Algoritmo

1. Leer cada calificación válida del archivo y agregarla a un vector.
2. Ordenar el vector de forma ascendente (`std::sort` en C++, `sort.Ints` en
   Go).
3. Recorrer el vector ordenado contando cuántos elementos iguales aparecen
   seguidos (una "corrida"). Si una corrida es mayor o igual a la más
   grande vista hasta el momento, se actualiza la moda; como se recorre de
   menor a mayor usando `>=`, un empate siempre lo termina ganando la
   calificación mayor.
4. Porcentaje = (frecuencia de la moda / total de datos) × 100.

## Complejidad

Leer el archivo es O(n). Ordenar es O(n log n) y domina el costo total.
Contar las corridas sobre el vector ya ordenado es O(n). El programa
completo es O(n log n) en tiempo y O(n) en espacio, porque aquí sí se
guarda cada calificación en memoria (a diferencia de la Propuesta #1, que
solo guarda contadores).

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
cd propuesta-2/cpp && cmake -S . -B build && cmake --build build && ./build/grade_mode_finder
cd propuesta-2/go && go run .
```

## Comparación con la Propuesta #1

Con calificaciones acotadas a [1, 100], el arreglo de contadores de la
Propuesta #1 es más eficiente en tiempo (O(n) vs. O(n log n)) y en espacio
extra (O(1) vs. O(n)). Ordenar sería la opción natural si el rango de
valores no fuera pequeño y conocido de antemano.
