// Propuesta #1: arreglo de contadores.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gradeapp/grades"
)

// nonDataFiles: archivos .txt que no son datos de calificaciones.
var nonDataFiles = map[string]bool{
	"leeme.txt":                true,
	"resultados_esperados.txt": true,
}

// gradeModeApp orquesta el flujo de lectura, procesamiento y resultado global.
type gradeModeApp struct {
	reader       *grades.FileReader
	analyzer     *grades.Analyzer
	globalTable  *grades.FrequencyTable
	filesOK      int
	totalElapsed time.Duration
	stdinScanner *bufio.Scanner
}

func newGradeModeApp() *gradeModeApp {
	return &gradeModeApp{
		reader:       grades.NewFileReader(),
		analyzer:     grades.NewAnalyzer(),
		globalTable:  grades.NewFrequencyTable(),
		stdinScanner: bufio.NewScanner(os.Stdin),
	}
}

// run: modo interactivo, pide archivos hasta un renglón vacío.
func (app *gradeModeApp) run() {
	fmt.Println("=== Calificación más común (Propuesta #1: arreglo de contadores) ===")

	for {
		fileName, hasMore := app.promptForFileName()
		if !hasMore {
			break
		}
		app.processFile(fileName)
	}

	app.printGlobalResult()
}

/**
 * runAutoDir procesa automáticamente todos los .txt de un directorio,
 * en orden alfabético, sin pedir nada por teclado.
 * @param dir - ruta del directorio a procesar.
 */
func (app *gradeModeApp) runAutoDir(dir string) error {
	fmt.Printf("=== Calificación más común (Propuesta #1: arreglo de contadores) ===\n")
	fmt.Printf("Modo automático: procesando todos los archivos de %q\n", dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("no se pudo leer el directorio %q: %w", dir, err)
	}

	var fileNames []string
	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".txt" {
			continue
		}
		if nonDataFiles[strings.ToLower(entry.Name())] {
			continue
		}
		fileNames = append(fileNames, entry.Name())
	}
	sort.Strings(fileNames)

	if len(fileNames) == 0 {
		return fmt.Errorf("no se encontraron archivos .txt de calificaciones en %q", dir)
	}

	for _, name := range fileNames {
		app.processFile(filepath.Join(dir, name))
	}

	app.printGlobalResult()
	return nil
}

// promptForFileName devuelve false si el usuario ya no quiere dar más archivos.
func (app *gradeModeApp) promptForFileName() (string, bool) {
	fmt.Print("\nIngrese el nombre del archivo (Enter para terminar): ")
	if !app.stdinScanner.Scan() {
		return "", false
	}
	fileName := strings.TrimSpace(app.stdinScanner.Text())
	return fileName, fileName != ""
}

/**
 * processFile lee, procesa y muestra el resultado de un archivo.
 * @param fileName - ruta o nombre del archivo a procesar.
 */
func (app *gradeModeApp) processFile(fileName string) {
	start := time.Now()

	fileTable := grades.NewFrequencyTable()
	result, err := app.reader.ReadInto(fileName, fileTable)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if fileTable.IsEmpty() {
		fmt.Printf("No se encontraron calificaciones válidas en %q.\n", fileName)
		return
	}

	mode := fileTable.MostCommonGrade()
	elapsed := time.Since(start)

	if result.LinesSkipped > 0 {
		fmt.Printf("Aviso: se ignoraron %d línea(s) inválida(s) en %q.\n", result.LinesSkipped, fileName)
	}
	fmt.Println(app.analyzer.DescribeFile(fileName, fileTable.Total(), mode, elapsed))

	app.globalTable.MergeFrom(fileTable)
	app.totalElapsed += elapsed
	app.filesOK++
}

// printGlobalResult muestra el resultado acumulado de todos los archivos.
func (app *gradeModeApp) printGlobalResult() {
	fmt.Println()
	if app.filesOK == 0 || app.globalTable.IsEmpty() {
		fmt.Println("No se procesó ningún archivo válido. Fin del programa.")
		return
	}

	modeStart := time.Now()
	globalMode := app.globalTable.MostCommonGrade()
	modeElapsed := time.Since(modeStart)

	fmt.Println(app.analyzer.DescribeGlobal(app.globalTable.Total(), globalMode, modeElapsed, app.totalElapsed))
}

func main() {
	dir := flag.String("dir", "", "Si se especifica, procesa automáticamente todos los .txt de este directorio en vez de preguntar por archivos uno por uno.")
	flag.Parse()

	app := newGradeModeApp()

	if *dir != "" {
		if err := app.runAutoDir(*dir); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		return
	}

	app.run()
}
