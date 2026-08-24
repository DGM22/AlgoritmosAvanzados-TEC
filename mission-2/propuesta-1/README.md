# Propuesta de solución #1 — Calificación más común

Implementación en C++ y Go de la Propuesta #1: un arreglo de contadores de
tamaño fijo (101 posiciones, una por cada calificación de 1 a 100). Al leer
el archivo, cada calificación incrementa directamente su contador
(`contadores[calificacion]++`); al terminar, se recorre el arreglo de 100 a
1 buscando el contador más alto, que es la moda.

Se implementó dos veces con el mismo diseño: en C++ con clases
(`FrequencyTable`, `FileReader`, `Analyzer`) y en Go con structs + métodos
con receptor, que es el equivalente idiomático a una clase en un lenguaje
sin herencia.

## Algoritmo

1. Crear un arreglo `contadores[1..100]` en cero.
2. Por cada calificación leída: `contadores[calificacion]++`.
3. Recorrer el arreglo de 100 a 1 buscando el valor más alto. Al usar
   comparación estricta (`>`), el primer máximo encontrado en ese recorrido
   es siempre la calificación mayor, resolviendo empates a su favor.
4. Porcentaje = (frecuencia de la moda / total de datos) × 100.

## Complejidad

Leer el archivo y llenar los contadores es O(n), con n = cantidad de
calificaciones. Encontrar la moda es O(100) = O(1), porque el arreglo
siempre tiene el mismo tamaño sin importar cuántos datos haya. El programa
completo es O(n) en tiempo y O(1) en espacio extra: solo se guardan 101
contadores, nunca los datos crudos.

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
cd propuesta-1/cpp && cmake -S . -B build && cmake --build build && ./build/grade_mode_finder
cd propuesta-1/go && go run .
```

## Validación

Se probó contra los 6 archivos de `calificaciones_caso_pruebas/`,
comparando contra `resultados_esperados.txt`. Ambas versiones (C++ y Go)
reproducen exactamente los valores esperados, incluyendo el caso de empate
(90 vs. 95, gana 95 por ser la mayor) y el resultado global acumulado de
611,124 datos.
