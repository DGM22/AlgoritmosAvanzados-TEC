/**
 * Quick Sort de 2 particiones (un pivote: el último elemento).
 *
 * Generado con inteligencia artificial OpenAI.
 *
 * Esquema de Lomuto: menores-o-iguales | pivote | mayores.
 *
 * Complejidad de tiempo: Theta(n log n) promedio, Theta(n^2) peor caso
 * (arreglo invertido o casi ordenado con este pivote).
 * Complejidad de espacio: O(1) extra + O(log n) de pila en promedio,
 * O(n) de pila en el peor caso.
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
 * swap intercambia a[i] y a[j].
 * Si i == j no cuenta: no hubo intercambio real.
 * Un intercambio equivale a 3 movimientos (temp, a[i], a[j]).
 *
 * @param a - arreglo.
 * @param i - primera posición.
 * @param j - segunda posición.
 */
func swap(a []int, i int, j int) {

	if i == j {
		return
	}

	temp := a[i]
	a[i] = a[j]
	a[j] = temp

	intercambios++
	movimientos += 3
}

/**
 * partition reordena a[low:high] usando a[high] como pivote.
 *
 * @param a - arreglo.
 * @param low - inicio inclusivo.
 * @param high - fin inclusivo (posición del pivote).
 * @returns índice final del pivote.
 */
func partition(a []int, low int, high int) int {

	pivot := a[high]

	// i marca el final de la región de elementos <= pivote.
	i := low - 1

	for j := low; j < high; j++ {

		comparaciones++

		if a[j] <= pivot {

			i++
			swap(a, i, j)
		}
	}

	swap(a, i+1, high)

	return i + 1
}

/**
 * quickSort ordena a[low:high] de forma recursiva (índices inclusivos).
 *
 * @param a - arreglo a ordenar en sitio.
 * @param low - inicio inclusivo.
 * @param high - fin inclusivo.
 */
func quickSort(a []int, low int, high int) {

	if low >= high {
		return
	}

	pivotPosition := partition(a, low, high)

	quickSort(a, low, pivotPosition-1)
	quickSort(a, pivotPosition+1, high)
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

	imprimirEncabezado("QUICK SORT - 2 PARTICIONES")

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

		quickSort(
			trabajo,
			0,
			len(trabajo)-1,
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
