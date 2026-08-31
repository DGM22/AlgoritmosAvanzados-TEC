/**
 * Merge Sort de 3 particiones (divide el arreglo en tercios).
 *
 * Generado con inteligencia artificial OpenAI.
 *
 * El combine no es un merge ternario: junta tercio1+tercio2 y luego
 * ese resultado con tercio3, usando el mismo merge binario.
 *
 * Complejidad de tiempo: Theta(n log n).
 * Complejidad de espacio: Theta(n) por el arreglo temporal.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

/**
 * Contadores globales de la corrida actual.
 * Se reinician antes de ordenar cada archivo.
 */
var comparaciones int64
var movimientos int64
var intercambios int64

/**
 * merge une dos segmentos ya ordenados: [low, mid) y [mid, high).
 * Si alguno está vacío, no hace trabajo.
 *
 * @param a - arreglo sobre el que se escribe el resultado.
 * @param low - inicio inclusivo del primer segmento.
 * @param mid - frontera entre los dos segmentos.
 * @param high - fin exclusivo del segundo segmento.
 */
func merge(
	a []int,
	low int,
	mid int,
	high int,
) {

	if low >= mid || mid >= high {
		return
	}

	temp := make(
		[]int,
		high-low,
	)

	i := low
	j := mid
	k := 0

	// Elige el menor de las dos cabezas hasta agotar un lado.
	for i < mid && j < high {

		comparaciones++

		if a[i] <= a[j] {

			temp[k] = a[i]
			i++

		} else {

			temp[k] = a[j]
			j++
		}

		movimientos++
		k++
	}

	for i < mid {

		temp[k] = a[i]
		movimientos++
		i++
		k++
	}

	for j < high {

		temp[k] = a[j]
		movimientos++
		j++
		k++
	}

	for x := 0; x < len(temp); x++ {

		a[low+x] = temp[x]
		movimientos++
	}
}

/**
 * mergeSort3 ordena a[low:high) dividiendo en tres tercios y fusionando
 * en dos pasos binarios.
 *
 * Rangos:
 *   [low, corte1) | [corte1, corte2) | [corte2, high)
 *
 * @param a - arreglo a ordenar.
 * @param low - inicio inclusivo.
 * @param high - fin exclusivo.
 */
func mergeSort3(
	a []int,
	low int,
	high int,
) {

	n := high - low

	if n <= 1 {
		return
	}

	corte1 := low + n/3
	corte2 := low + (2*n)/3

	mergeSort3(a, low, corte1)
	mergeSort3(a, corte1, corte2)
	mergeSort3(a, corte2, high)

	// Primero quedan ordenados los dos tercios de la izquierda.
	merge(a, low, corte1, corte2)

	// Luego se incorpora el tercio derecho al bloque ya fusionado.
	merge(a, low, corte2, high)
}

/**
 * main recorre los .txt de datos/, ordena una copia de cada uno y reporta métricas.
 * Usa -n para filtrar por tamaño (0 = todos los archivos).
 */
func main() {

	soloN := flag.Int("n", 0, "solo archivos de este tamaño (0 = todos)")
	flag.Parse()

	archivos, err := archivosTamano(*soloN)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	imprimirEncabezado("MERGE SORT - 3 PARTICIONES")

	for _, ruta := range archivos {

		original, err := cargar(ruta)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		trabajo := copia(original)

		comparaciones = 0
		movimientos = 0
		intercambios = 0

		inicio := time.Now()

		mergeSort3(
			trabajo,
			0,
			len(trabajo),
		)

		tiempo := time.Since(inicio)

		imprimirFila(
			filepath.Base(ruta),
			len(trabajo),
			comparaciones,
			movimientos,
			intercambios,
			tiempo,
			sort.IntsAreSorted(trabajo),
		)
	}
}
