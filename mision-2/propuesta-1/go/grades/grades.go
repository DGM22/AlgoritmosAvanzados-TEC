// Propuesta #1: arreglo de contadores indexado por calificación.
package grades

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Rango de calificaciones permitido.
const (
	MinGrade = 1
	MaxGrade = 100
)

// OutOfRangeError: calificación fuera de [MinGrade, MaxGrade].
type OutOfRangeError struct {
	Grade int
}

func (e *OutOfRangeError) Error() string {
	return fmt.Sprintf("la calificación %d está fuera del rango permitido [%d, %d]",
		e.Grade, MinGrade, MaxGrade)
}

// FrequencyTable: contador por calificación.
type FrequencyTable struct {
	counters [MaxGrade + 1]int64
	total    int64
}

func NewFrequencyTable() *FrequencyTable {
	return &FrequencyTable{}
}

/**
 * Increment suma 1 al contador de una calificación.
 * @param grade - calificación a incrementar; debe estar en [MinGrade, MaxGrade].
 * @returns error si grade está fuera de rango.
 */
func (t *FrequencyTable) Increment(grade int) error {
	if grade < MinGrade || grade > MaxGrade {
		return &OutOfRangeError{Grade: grade}
	}
	t.counters[grade]++
	t.total++
	return nil
}

/**
 * CountFor devuelve cuántas veces se registró una calificación.
 * @param grade - calificación a consultar.
 * @returns 0 si grade está fuera de rango.
 */
func (t *FrequencyTable) CountFor(grade int) int64 {
	if grade < MinGrade || grade > MaxGrade {
		return 0
	}
	return t.counters[grade]
}

// Total devuelve la cantidad de calificaciones registradas.
func (t *FrequencyTable) Total() int64 {
	return t.total
}

// IsEmpty indica si no hay calificaciones registradas.
func (t *FrequencyTable) IsEmpty() bool {
	return t.total == 0
}

// Mode: calificación más común, su frecuencia y porcentaje.
type Mode struct {
	Grade      int
	Count      int64
	Percentage float64
}

// MostCommonGrade recorre de 100 a 1; en empate gana la calificación mayor.
func (t *FrequencyTable) MostCommonGrade() Mode {
	var result Mode
	if t.total == 0 {
		return result
	}

	for grade := MaxGrade; grade >= MinGrade; grade-- {
		count := t.counters[grade]
		if count > result.Count {
			result.Grade = grade
			result.Count = count
		}
	}

	result.Percentage = (float64(result.Count) / float64(t.total)) * 100
	return result
}

/**
 * MergeFrom acumula los contadores de otra tabla en esta.
 * @param other - tabla cuyos contadores se sumarán a esta.
 */
func (t *FrequencyTable) MergeFrom(other *FrequencyTable) {
	for grade := MinGrade; grade <= MaxGrade; grade++ {
		t.counters[grade] += other.counters[grade]
	}
	t.total += other.total
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
 * ReadInto lee un archivo línea por línea e incrementa la tabla por
 * cada calificación válida.
 * @param filePath - ruta del archivo a leer.
 * @param table - tabla que se incrementará con cada calificación válida.
 */
func (r *FileReader) ReadInto(filePath string, table *FrequencyTable) (ReadResult, error) {
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

		if incErr := table.Increment(grade); incErr != nil {
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
 * @param elapsed - tiempo de leer + contar + calcular la moda (sin imprimir).
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
 * @param modeElapsed - tiempo de calcular la moda global.
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
