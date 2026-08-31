/**
 * Quick Sort de 3 particiones (dos pivotes: primero y último).
 *
 * Generado con inteligencia artificial OpenAI.
 *
 * Tras partir: menores | pivote1 | intermedios | pivote2 | mayores.
 * No es el 3-way de Dijkstra (iguales a un pivote); es Dual-Pivot.
 *
 * Complejidad de tiempo: Theta(n log n) promedio, Theta(n^2) peor caso
 * (por ejemplo invertido: los pivotes son max y min).
 * Complejidad de espacio: O(1) extra + pila O(log n) promedio, O(n) peor.
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
 * Si i == j no cuenta. Un intercambio = 3 movimientos.
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
 * partition3 clasifica a[low:high] con dos pivotes (primero y último).
 * Garantiza pivot1 <= pivot2 intercambiándolos al inicio si hace falta.
 *
 * @param a - arreglo.
 * @param low - inicio inclusivo.
 * @param high - fin inclusivo.
 * @returns posiciones finales de pivot1 y pivot2.
 */
func partition3(
	a []int,
	low int,
	high int,
) (int, int) {

	comparaciones++

	if a[low] > a[high] {
		swap(a, low, high)
	}

	pivot1 := a[low]
	pivot2 := a[high]

	left := low + 1
	right := high - 1
	i := left

	for i <= right {

		comparaciones++

		if a[i] < pivot1 {

			swap(a, i, left)
			left++
			i++
			continue
		}

		comparaciones++

		if a[i] > pivot2 {

			swap(a, i, right)
			right--

			// El valor que llegó de la derecha aún no está clasificado.
			continue
		}

		// pivot1 <= a[i] <= pivot2: se queda en la región central.
		i++
	}

	left--
	right++

	swap(a, low, left)
	swap(a, high, right)

	return left, right
}

/**
 * quickSort3 ordena a[low:high] con tres llamadas recursivas (índices inclusivos).
 *
 * @param a - arreglo a ordenar en sitio.
 * @param low - inicio inclusivo.
 * @param high - fin inclusivo.
 */
func quickSort3(
	a []int,
	low int,
	high int,
) {

	if low >= high {
		return
	}

	pivot1, pivot2 := partition3(a, low, high)

	quickSort3(a, low, pivot1-1)
	quickSort3(a, pivot1+1, pivot2-1)
	quickSort3(a, pivot2+1, high)
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

	imprimirEncabezado("QUICK SORT - 3 PARTICIONES")

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

		quickSort3(
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
