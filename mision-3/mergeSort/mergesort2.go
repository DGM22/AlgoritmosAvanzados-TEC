/**
 * Merge Sort de 2 particiones (divide el arreglo a la mitad).
 *
 * Generado con inteligencia artificial OpenAI.
 *
 * Complejidad de tiempo: Theta(n log n) en todos los casos.
 * Complejidad de espacio: Theta(n) por el arreglo temporal del merge.
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
 * Usa un arreglo temporal y luego copia el resultado de vuelta a a.
 *
 * Merge Sort no intercambia en sitio; intercambios permanece en 0.
 *
 * @param a - arreglo sobre el que se escribe el resultado.
 * @param low - inicio inclusivo del primer segmento.
 * @param mid - fin exclusivo del primero / inicio del segundo.
 * @param high - fin exclusivo del segundo segmento.
 */
func merge(
	a []int,
	low int,
	mid int,
	high int,
) {

	// Buffer auxiliar del tamaño exacto del rango a fusionar.
	temp := make(
		[]int,
		high-low,
	)

	i := low
	j := mid
	k := 0

	// Consume el menor de las dos cabezas mientras ambas tengan datos.
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

	// Copia lo que quedó en la mitad izquierda.
	for i < mid {

		temp[k] = a[i]
		movimientos++
		i++
		k++
	}

	// Copia lo que quedó en la mitad derecha.
	for j < high {

		temp[k] = a[j]
		movimientos++
		j++
		k++
	}

	// Vuelve a escribir el rango ordenado sobre el arreglo original.
	for x := 0; x < len(temp); x++ {

		a[low+x] = temp[x]
		movimientos++
	}
}

/**
 * mergeSort ordena a[low:high) partiendo a la mitad y fusionando.
 *
 * @param a - arreglo a ordenar en sitio (con buffer auxiliar en merge).
 * @param low - inicio inclusivo.
 * @param high - fin exclusivo.
 */
func mergeSort(
	a []int,
	low int,
	high int,
) {

	// Un segmento de 0 o 1 elementos ya está ordenado.
	if high-low <= 1 {
		return
	}

	mid := low + (high-low)/2

	mergeSort(a, low, mid)
	mergeSort(a, mid, high)
	merge(a, low, mid, high)
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

	imprimirEncabezado("MERGE SORT - 2 PARTICIONES")

	for _, ruta := range archivos {

		original, err := cargar(ruta)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Ordenar una copia evita mutar los datos leídos del archivo.
		trabajo := copia(original)

		comparaciones = 0
		movimientos = 0
		intercambios = 0

		inicio := time.Now()

		mergeSort(
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
