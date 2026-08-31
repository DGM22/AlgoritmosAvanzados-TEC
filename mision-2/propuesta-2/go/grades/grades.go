// Propuesta #2: vector dinámico + ordenamiento + contador único de corridas contiguas.
package grades

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

// GradeVector: vector dinámico de calificaciones leídas.
type GradeVector struct {
	values []int
}

func NewGradeVector() *GradeVector {
	return &GradeVector{}
}

/**
 * Add agrega una calificación al vector.
 * @param grade - calificación a agregar; debe estar en [MinGrade, MaxGrade].
 * @returns error si grade está fuera de rango.
 */
func (v *GradeVector) Add(grade int) error {
	if grade < MinGrade || grade > MaxGrade {
		return &OutOfRangeError{Grade: grade}
	}
	v.values = append(v.values, grade)
	return nil
}

// Total devuelve la cantidad de calificaciones registradas.
func (v *GradeVector) Total() int64 {
	return int64(len(v.values))
}

// IsEmpty indica si el vector está vacío.
func (v *GradeVector) IsEmpty() bool {
	return len(v.values) == 0
}

// Sort ordena el vector de forma ascendente (sort.Ints).
func (v *GradeVector) Sort() {
	sort.Ints(v.values)
}


/**
 * Merge agrega los valores de otro vector a este.
 * @param other - vector cuyos valores se agregarán a este.
 */
func (v *GradeVector) Merge(other *GradeVector) {
	v.values = append(v.values, other.values...)
}

// Mode: calificación más común, su frecuencia y porcentaje.
type Mode struct {
	Grade      int
	Count      int64
	Percentage float64
}

// MostCommonGrade cuenta corridas contiguas en el vector ya ordenado (llamar Sort() antes); empata a favor del valor mayor.
func (v *GradeVector) MostCommonGrade() Mode {
	var result Mode

	total := len(v.values)
	if total == 0 {
		return result
	}

	runGrade := v.values[0]
	runCount := 1
	result.Grade = runGrade
	result.Count = 1

	for i := 1; i < total; i++ {
		if v.values[i] == runGrade {
			runCount++
		} else {
			runGrade = v.values[i]
			runCount = 1
		}

		if int64(runCount) >= result.Count {
			result.Count = int64(runCount)
			result.Grade = runGrade
		}
	}

	result.Percentage = (float64(result.Count) / float64(total)) * 100
	return result
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
 * ReadInto lee un archivo línea por línea y agrega cada calificación
 * válida al vector.
 * @param filePath - ruta del archivo a leer.
 * @param vector - vector al que se agregará cada calificación válida.
 */
func (r *FileReader) ReadInto(filePath string, vector *GradeVector) (ReadResult, error) {
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

		if addErr := vector.Add(grade); addErr != nil {
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
 * @param elapsed - tiempo de leer + ordenar + contar (sin imprimir).
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
 * @param modeElapsed - tiempo de volver a ordenar + calcular la moda global.
 * @param totalElapsed - suma de los tiempos de procesamiento por archivo.
 */
func (a *Analyzer) DescribeGlobal(totalCount int64, mode Mode, modeElapsed, totalElapsed time.Duration) string {
	return fmt.Sprintf(
		"RESULTADO ACUMULADO DE TODOS LOS ARCHIVOS\n-----------------------------------------\nCantidad total de datos: %s\nCalificación más común global: %d\nFrecuencia global: %s\nPorcentaje global: %s%%\nTiempo del ordenamiento + cálculo de la moda global: %s\nTiempo total acumulado (todos los archivos): %s",
		withThousandsSeparator(totalCount),
		mode.Grade,
		withThousandsSeparator(mode.Count),
		formatPercentage(mode.Percentage),
		modeElapsed.Round(time.Microsecond),
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
