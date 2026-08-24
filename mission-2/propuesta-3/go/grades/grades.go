// Propuesta #3: tabla hash (10 posiciones) con encadenamiento externo.
package grades

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Rango de calificaciones permitido. tableSize = posiciones de la tabla.
const (
	MinGrade  = 1
	MaxGrade  = 100
	tableSize = 10
)

// OutOfRangeError: calificación fuera de [MinGrade, MaxGrade].
type OutOfRangeError struct {
	Grade int
}

func (e *OutOfRangeError) Error() string {
	return fmt.Sprintf("la calificación %d está fuera del rango permitido [%d, %d]",
		e.Grade, MinGrade, MaxGrade)
}

// hashNode: nodo de la cadena de cada posición (calificación única + contador).
type hashNode struct {
	grade int
	count int64
	next  *hashNode
}

// HashTable: arreglo de tableSize posiciones, cada una con su cadena de nodos.
type HashTable struct {
	buckets [tableSize]*hashNode
	total   int64
}

func NewHashTable() *HashTable {
	return &HashTable{}
}

/**
 * hashIndex calcula la posición de la tabla para una calificación.
 * @param grade - calificación a mapear (usa su último dígito).
 */
func hashIndex(grade int) int {
	return grade % tableSize
}

/**
 * Insert agrega una aparición de una calificación. Si ya existe en su
 * cadena, solo incrementa el contador; si no, crea un nodo nuevo.
 * @param grade - calificación a insertar; debe estar en [MinGrade, MaxGrade].
 * @returns error si grade está fuera de rango.
 */
func (t *HashTable) Insert(grade int) error {
	if grade < MinGrade || grade > MaxGrade {
		return &OutOfRangeError{Grade: grade}
	}

	idx := hashIndex(grade)
	for node := t.buckets[idx]; node != nil; node = node.next {
		if node.grade == grade {
			node.count++
			t.total++
			return nil
		}
	}

	t.buckets[idx] = &hashNode{grade: grade, count: 1, next: t.buckets[idx]}
	t.total++
	return nil
}

// Total devuelve la cantidad de calificaciones registradas.
func (t *HashTable) Total() int64 {
	return t.total
}

// IsEmpty indica si no hay calificaciones registradas.
func (t *HashTable) IsEmpty() bool {
	return t.total == 0
}

// Mode: calificación más común, su frecuencia y porcentaje.
type Mode struct {
	Grade      int
	Count      int64
	Percentage float64
}

// MostCommonGrade recorre toda la tabla; en empate gana la calificación mayor.
func (t *HashTable) MostCommonGrade() Mode {
	var result Mode
	if t.total == 0 {
		return result
	}

	for _, head := range t.buckets {
		for node := head; node != nil; node = node.next {
			if node.count > result.Count ||
				(node.count == result.Count && node.grade > result.Grade) {
				result.Count = node.count
				result.Grade = node.grade
			}
		}
	}

	result.Percentage = (float64(result.Count) / float64(t.total)) * 100
	return result
}

/**
 * MergeFrom acumula los pares (calificación, contador) de otra tabla.
 * @param other - tabla cuyos pares se acumularán en esta.
 */
func (t *HashTable) MergeFrom(other *HashTable) {
	for _, head := range other.buckets {
		for node := head; node != nil; node = node.next {
			t.addCount(node.grade, node.count)
		}
	}
}

/**
 * addCount suma apariciones de una calificación (usado por Insert y MergeFrom).
 * @param grade - calificación a la que se le suma el contador.
 * @param count - cantidad a sumar.
 */
func (t *HashTable) addCount(grade int, count int64) {
	idx := hashIndex(grade)
	for node := t.buckets[idx]; node != nil; node = node.next {
		if node.grade == grade {
			node.count += count
			t.total += count
			return
		}
	}
	t.buckets[idx] = &hashNode{grade: grade, count: count, next: t.buckets[idx]}
	t.total += count
}

// ReadResult: líneas procesadas y omitidas al leer un archivo.
type ReadResult struct {
	LinesProcessed int64
	LinesSkipped   int64
}

// FileReader lee calificaciones (una por línea) desde un archivo.
type FileReader struct{}

func NewFileReader() *FileReader {
	return &FileReader{}
}

/**
 * ReadInto lee un archivo línea por línea e inserta cada calificación
 * válida en la tabla.
 * @param filePath - ruta del archivo a leer.
 * @param table - tabla en la que se insertará cada calificación válida.
 */
func (r *FileReader) ReadInto(filePath string, table *HashTable) (ReadResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ReadResult{}, fmt.Errorf("no se pudo abrir el archivo %q: %w", filePath, err)
	}
	defer file.Close()

	var result ReadResult

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		grade, parseErr := strconv.Atoi(line)
		if parseErr != nil {
			result.LinesSkipped++
			continue
		}

		if insErr := table.Insert(grade); insErr != nil {
			result.LinesSkipped++
			continue
		}

		result.LinesProcessed++
	}

	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("error leyendo el archivo %q: %w", filePath, err)
	}

	return result, nil
}

// Analyzer da formato a los resultados para el usuario.
type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

/**
 * DescribeFile formatea el resultado de un archivo procesado.
 * @param fileName - nombre del archivo procesado.
 * @param totalCount - cantidad total de calificaciones válidas leídas.
 * @param mode - calificación más común calculada.
 * @param elapsed - tiempo de leer + insertar + calcular la moda (sin imprimir).
 */
func (a *Analyzer) DescribeFile(fileName string, totalCount int64, mode Mode, elapsed time.Duration) string {
	return fmt.Sprintf(
		"Archivo: %s\nCantidad de datos: %s\nCalificación más común: %d\nFrecuencia: %s\nPorcentaje: %s%%\nTiempo de procesamiento: %s",
		fileName,
		withThousandsSeparator(totalCount),
		mode.Grade,
		withThousandsSeparator(mode.Count),
		formatPercentage(mode.Percentage),
		elapsed.Round(time.Microsecond),
	)
}

/**
 * DescribeGlobal formatea el resultado acumulado de todos los archivos.
 * @param totalCount - cantidad total de calificaciones válidas de todos los archivos.
 * @param mode - calificación más común global calculada.
 * @param modeElapsed - tiempo de recorrer la tabla global para hallar la moda.
 * @param totalElapsed - suma de los tiempos de procesamiento por archivo.
 */
func (a *Analyzer) DescribeGlobal(totalCount int64, mode Mode, modeElapsed, totalElapsed time.Duration) string {
	return fmt.Sprintf(
		"RESULTADO ACUMULADO DE TODOS LOS ARCHIVOS\n-----------------------------------------\nCantidad total de datos: %s\nCalificación más común global: %d\nFrecuencia global: %s\nPorcentaje global: %s%%\nTiempo del cálculo de la moda global: %s\nTiempo total acumulado (todos los archivos): %s",
		withThousandsSeparator(totalCount),
		mode.Grade,
		withThousandsSeparator(mode.Count),
		formatPercentage(mode.Percentage),
		modeElapsed,
		totalElapsed.Round(time.Microsecond),
	)
}

/**
 * withThousandsSeparator agrega separador de miles a un entero.
 * @param value - valor a formatear (p. ej. 611124 -> "611,124").
 */
func withThousandsSeparator(value int64) string {
	negative := value < 0
	if negative {
		value = -value
	}

	digits := strconv.FormatInt(value, 10)

	var groups []string
	for len(digits) > 3 {
		cut := len(digits) - 3
		groups = append([]string{digits[cut:]}, groups...)
		digits = digits[:cut]
	}
	groups = append([]string{digits}, groups...)

	result := strings.Join(groups, ",")
	if negative {
		result = "-" + result
	}
	return result
}

/**
 * formatPercentage da formato a un porcentaje con decimales fijos.
 * @param value - valor a formatear (p. ej. 18.0 -> "18.0000").
 */
func formatPercentage(value float64) string {
	return strconv.FormatFloat(value, 'f', 4, 64)
}
