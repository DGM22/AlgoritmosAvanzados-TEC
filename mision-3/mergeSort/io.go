package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func raizMision() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for i := 0; i < 10; i++ {
		if esRaiz(dir) {
			return dir, nil
		}

		anidada := filepath.Join(dir, "mision-3")
		if esRaiz(anidada) {
			return anidada, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("no se encontró mision-3")
}

func esRaiz(dir string) bool {
	_, errMerge := os.Stat(filepath.Join(dir, "mergeSort"))
	_, errQuick := os.Stat(filepath.Join(dir, "quickSort"))
	return errMerge == nil && errQuick == nil
}

func dirDatos() (string, error) {
	raiz, err := raizMision()
	if err != nil {
		return "", err
	}
	return filepath.Join(raiz, "datos"), nil
}

func archivosTamano(n int) ([]string, error) {
	dir, err := dirDatos()
	if err != nil {
		return nil, err
	}

	entradas, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer %q: %w", dir, err)
	}

	var rutas []string
	prefijo := ""
	if n > 0 {
		prefijo = fmt.Sprintf("n%06d_", n)
	}

	for _, e := range entradas {
		if e.IsDir() || strings.ToLower(filepath.Ext(e.Name())) != ".txt" {
			continue
		}
		if prefijo != "" && !strings.HasPrefix(e.Name(), prefijo) {
			continue
		}
		rutas = append(rutas, filepath.Join(dir, e.Name()))
	}

	sort.Strings(rutas)

	if len(rutas) == 0 {
		if n > 0 {
			return nil, fmt.Errorf("no hay .txt de n=%d en %s", n, dir)
		}
		return nil, fmt.Errorf("no hay .txt en %s", dir)
	}

	return rutas, nil
}

func cargar(ruta string) ([]int, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var valores []int
	escaner := bufio.NewScanner(f)
	escaner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	linea := 0
	for escaner.Scan() {
		linea++
		texto := strings.TrimSpace(escaner.Text())
		if texto == "" {
			continue
		}

		v, err := strconv.Atoi(texto)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: no es un entero: %q", ruta, linea, texto)
		}
		valores = append(valores, v)
	}

	return valores, escaner.Err()
}

func copia(origen []int) []int {
	destino := make([]int, len(origen))
	copy(destino, origen)
	return destino
}

func imprimirEncabezado(titulo string) {
	fmt.Println()
	fmt.Println("==========", titulo, "==========")
	fmt.Printf(
		"%-28s %8s %14s %14s %14s %14s %6s\n",
		"archivo",
		"n",
		"comparaciones",
		"movimientos",
		"intercambios",
		"tiempo",
		"ok",
	)
}

func imprimirFila(
	archivo string,
	n int,
	comparaciones int64,
	movimientos int64,
	intercambios int64,
	tiempo time.Duration,
	ok bool,
) {
	okTxt := "false"
	if ok {
		okTxt = "true"
	}

	fmt.Printf(
		"%-28s %8d %14d %14d %14d %14s %6s\n",
		archivo,
		n,
		comparaciones,
		movimientos,
		intercambios,
		tiempo,
		okTxt,
	)
}
