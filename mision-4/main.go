package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/**
 * main pide el tamaño de la red, las dimensiones de cada cisterna
 * y el volumen de agua; luego imprime nivel, OVERFLOW si aplica
 * y el porcentaje de llenado de cada depósito.
 */
func main() {
	lector := bufio.NewScanner(os.Stdin)
	lector.Split(bufio.ScanLines)

	fmt.Println("=== Red de cisternas ===")
	fmt.Println("Las cisternas están conectadas por tuberías.")
	fmt.Println("El agua fluye hacia las de abajo; las que quedan por encima del nivel se quedan vacías.")
	fmt.Println()

	n := leerEntero(lector, "¿Cuántas cisternas conforman la red? ", 1)

	red := Red{
		Cisternas: make([]Cisterna, n),
	}

	for i := 0; i < n; i++ {
		fmt.Printf("\n--- Cisterna %d ---\n", i+1)
		red.Cisternas[i] = leerCisterna(lector)
	}

	imprimirResumen(red)

	fmt.Println()
	volumen := leerReal(lector, "¿Cuántos metros cúbicos de agua vas a ingresar? ", 0, false)

	resultado := red.NivelAgua(volumen)
	imprimirResultado(red, resultado)
}

/**
 * leerCisterna pide ancho, largo, alto y la altura del embase
 * respecto al piso.
 *
 * @param lector - escáner de la entrada estándar.
 * @returns cisterna con las cuatro medidas en metros.
 */
func leerCisterna(lector *bufio.Scanner) Cisterna {
	return Cisterna{
		Ancho: leerReal(lector, "  Ancho (m):  ", 0, true),
		Largo: leerReal(lector, "  Largo (m):  ", 0, true),
		Alto:  leerReal(lector, "  Alto (m):   ", 0, true),
		Base:  leerReal(lector, "  Altura de la base respecto al piso (m): ", 0, false),
	}
}

/**
 * imprimirResumen muestra dimensiones, área, volumen y el rango
 * de cotas de la red.
 *
 * @param red - conjunto de cisternas ya leído.
 */
func imprimirResumen(red Red) {
	fmt.Println()
	fmt.Println("Resumen de la red")
	fmt.Println("-----------------")
	fmt.Printf("%-4s %10s %10s %10s %10s %12s %12s\n",
		"#", "Ancho", "Largo", "Alto", "Base", "Área m²", "Volumen m³")

	for i, c := range red.Cisternas {
		fmt.Printf("%-4d %10.2f %10.2f %10.2f %10.2f %12.2f %12.2f\n",
			i+1, c.Ancho, c.Largo, c.Alto, c.Base, c.AreaBase(), c.Volumen())
	}

	fmt.Printf("\nCapacidad total: %.2f m³\n", red.VolumenTotal())
	fmt.Printf("Rango de cotas:  %.2f m  →  %.2f m (piso → techo más alto)\n",
		red.baseMinima(), red.techoMaximo())
}

/**
 * imprimirResultado muestra OVERFLOW o el nivel encontrado, el
 * orden de llenado y el porcentaje de cada cisterna.
 *
 * @param red - red sobre la que se inyectó el agua.
 * @param r - resultado de NivelAgua (nivel o desborde).
 */
func imprimirResultado(red Red, r ResultadoNivel) {
	fmt.Println()

	if r.Overflow {
		imprimirOverflow(red, r)
		return
	}

	fmt.Printf("Nivel del agua: %.2f m respecto al piso.\n", r.Nivel)
	fmt.Printf("Agua ingresada: %.2f / %.2f m³\n", r.VolumenPed, r.Capacidad)

	imprimirOrden(red, r.VolumenPed)
	imprimirPorcentajes(red, r.Nivel)
}

/**
 * imprimirOverflow avisa que el volumen supera la capacidad de
 * toda la red y deja todas las cisternas al 100%.
 *
 * @param red - red saturada.
 * @param r - resultado con Overflow en true.
 */
func imprimirOverflow(red Red, r ResultadoNivel) {
	exceso := r.VolumenPed - r.Capacidad

	fmt.Println("OVERFLOW")
	fmt.Printf("El volumen ingresado (%.2f m³) supera la capacidad de la red (%.2f m³).\n",
		r.VolumenPed, r.Capacidad)
	fmt.Printf("Excedente: %.2f m³. Toda el agua de más se desborda.\n", exceso)
	fmt.Println("Todas las cisternas quedan llenas (100%).")

	imprimirPorcentajes(red, red.techoMaximo())
}

/**
 * imprimirOrden narra cómo baja el agua: cotas bajas primero,
 * cisternas lado a lado al mismo tiempo, huecos por tuberías.
 *
 * @param red - red conectada por tuberías.
 * @param volumen - m³ inyectados.
 */
func imprimirOrden(red Red, volumen float64) {
	etapas := red.OrdenLlenado(volumen)
	if len(etapas) == 0 {
		return
	}

	fmt.Println()
	fmt.Println("Cómo se llena el agua")
	fmt.Println("---------------------")
	fmt.Println("Por las tuberías el agua baja primero. Lo que está por encima del nivel queda vacío.")
	fmt.Println()

	for i, etapa := range etapas {
		fmt.Printf("%d. ", i+1)

		if etapa.Hueco {
			fmt.Printf("Entre %.2f m y %.2f m no hay cisterna; el sobrante sigue por las tuberías hacia arriba.\n",
				etapa.Desde, etapa.Hasta)
			continue
		}

		nombres := nombrarCisternas(etapa.Indices)
		if etapa.LadoALado() {
			fmt.Printf("El agua llega a %.2f m y entra a %s (lado a lado).\n",
				etapa.Desde, nombres)
			fmt.Printf("   Se llenan juntas de %.2f m a %.2f m (%.2f m³).\n",
				etapa.Desde, etapa.Hasta, etapa.Volumen)
		} else {
			fmt.Printf("El agua baja a %.2f m y entra a %s.\n",
				etapa.Desde, nombres)
			fmt.Printf("   Llena de %.2f m a %.2f m (%.2f m³).\n",
				etapa.Desde, etapa.Hasta, etapa.Volumen)
		}

		for _, idx := range etapa.Indices {
			c := red.Cisternas[idx]
			fmt.Printf("   Cisterna %d: %.2f%%\n",
				idx+1, c.PorcentajeLlenado(etapa.Hasta))
		}
	}
}

/**
 * imprimirPorcentajes muestra el llenado de cada cisterna, de
 * abajo hacia arriba, con su estado (vacía, parcial o llena).
 *
 * @param red - red a reportar.
 * @param nivel - cota del agua respecto al piso, en metros.
 */
func imprimirPorcentajes(red Red, nivel float64) {
	fmt.Println()
	fmt.Println("Porcentaje de llenado (de abajo hacia arriba)")
	fmt.Println("---------------------------------------------")
	fmt.Printf("%-4s %10s %14s %14s  %s\n", "#", "Base (m)", "Agua (m³)", "Llenado", "Estado")

	for _, i := range red.indicesPorBase() {
		c := red.Cisternas[i]
		pct := c.PorcentajeLlenado(nivel)
		fmt.Printf("%-4d %10.2f %14.2f %13.2f%%  %s\n",
			i+1,
			c.Base,
			c.VolumenHasta(nivel),
			pct,
			estadoLlenado(pct, c.Base, nivel),
		)
	}
}

/**
 * estadoLlenado describe si la cisterna está vacía (el agua se
 * fue abajo), parcial o llena.
 *
 * @param porcentaje - llenado de 0 a 100.
 * @param base - altura del embase respecto al piso, en metros.
 * @param nivel - cota actual del agua, en metros.
 * @returns etiqueta de estado para la tabla.
 */
func estadoLlenado(porcentaje float64, base float64, nivel float64) string {
	if porcentaje >= 100-1e-6 {
		return "llena"
	}
	if porcentaje <= 1e-6 {
		if base > nivel {
			return "vacía (el agua fluyó hacia abajo)"
		}
		return "vacía"
	}
	return "parcial"
}

/**
 * nombrarCisternas arma una frase con los números de cisterna
 * ("la cisterna 1" / "las cisternas 1 y 2").
 *
 * @param indices - índices 0-based de las cisternas.
 * @returns texto listo para imprimir.
 */
func nombrarCisternas(indices []int) string {
	if len(indices) == 0 {
		return "ninguna cisterna"
	}

	nums := make([]string, len(indices))
	for i, idx := range indices {
		nums[i] = strconv.Itoa(idx + 1)
	}

	if len(nums) == 1 {
		return "la cisterna " + nums[0]
	}
	if len(nums) == 2 {
		return "las cisternas " + nums[0] + " y " + nums[1]
	}

	return "las cisternas " + strings.Join(nums[:len(nums)-1], ", ") + " y " + nums[len(nums)-1]
}

/**
 * leerEntero pide un entero hasta que la entrada sea válida.
 *
 * @param lector - escáner de la entrada estándar.
 * @param prompt - texto a mostrar antes de leer.
 * @param minimo - valor mínimo aceptado (inclusive).
 * @returns entero ≥ minimo.
 */
func leerEntero(lector *bufio.Scanner, prompt string, minimo int) int {
	for {
		fmt.Print(prompt)
		texto, ok := leerLinea(lector)
		if !ok {
			os.Exit(1)
		}

		valor, err := strconv.Atoi(texto)
		if err != nil || valor < minimo {
			fmt.Printf("Ingrese un entero mayor o igual a %d.\n", minimo)
			continue
		}

		return valor
	}
}

/**
 * leerReal pide un flotante hasta que la entrada sea válida.
 * Si estricto, el valor debe ser mayor que minimo; si no, se
 * acepta mayor o igual.
 *
 * @param lector - escáner de la entrada estándar.
 * @param prompt - texto a mostrar antes de leer.
 * @param minimo - umbral numérico.
 * @param estricto - true exige > minimo; false acepta ≥ minimo.
 * @returns número leído.
 */
func leerReal(lector *bufio.Scanner, prompt string, minimo float64, estricto bool) float64 {
	for {
		fmt.Print(prompt)
		texto, ok := leerLinea(lector)
		if !ok {
			os.Exit(1)
		}

		valor, err := strconv.ParseFloat(texto, 64)
		if err != nil {
			fmt.Println("Ingrese un número válido.")
			continue
		}

		if estricto && valor <= minimo {
			fmt.Printf("El valor debe ser mayor que %.2f.\n", minimo)
			continue
		}
		if !estricto && valor < minimo {
			fmt.Printf("El valor no puede ser menor que %.2f.\n", minimo)
			continue
		}

		return valor
	}
}

/**
 * leerLinea lee un renglón de la entrada y quita espacios de borde.
 *
 * @param lector - escáner de la entrada estándar.
 * @returns texto recortado y true si hubo línea; "" y false al EOF o error.
 */
func leerLinea(lector *bufio.Scanner) (string, bool) {
	if !lector.Scan() {
		if err := lector.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error de lectura:", err)
		}
		return "", false
	}
	return strings.TrimSpace(lector.Text()), true
}
