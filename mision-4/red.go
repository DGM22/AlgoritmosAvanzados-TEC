package main

import "sort"

const (
	epsilon            = 1e-9
	iteracionesBinaria = 80
)

/**
 * evento marca un cambio de área transversal en una cota.
 * En la base de una cisterna el área aumenta; en el techo disminuye.
 */
type evento struct {
	z  float64
	dA float64
}

/**
 * ResultadoNivel es la respuesta al inyectar un volumen en la red.
 * Overflow es true si el agua no cabe en el conjunto de cisternas.
 */
type ResultadoNivel struct {
	Overflow   bool
	Nivel      float64
	Capacidad  float64
	VolumenPed float64
}

/**
 * Red es el conjunto de cisternas unidas por tuberías de volumen
 * despreciable. El agua fluye hacia las de abajo; las que quedan
 * por encima del nivel se quedan vacías. Si el volumen supera la
 * capacidad total de la red, hay OVERFLOW.
 */
type Red struct {
	Cisternas []Cisterna
}

/**
 * VolumenTotal suma la capacidad de todas las cisternas de la red.
 *
 * @returns capacidad agregada, en m³.
 */
func (r Red) VolumenTotal() float64 {
	total := 0.0
	for _, c := range r.Cisternas {
		total += c.Volumen()
	}
	return total
}

/**
 * VolumenHasta calcula el agua almacenada en toda la red si el
 * nivel libre está en una cota dada.
 *
 * @param nivel - cota absoluta del agua, en metros sobre el piso.
 * @returns volumen total a esa cota, en m³.
 */
func (r Red) VolumenHasta(nivel float64) float64 {
	total := 0.0
	for _, c := range r.Cisternas {
		total += c.VolumenHasta(nivel)
	}
	return total
}

/**
 * NivelAgua calcula la cota que alcanza un volumen dado.
 *
 * Si el volumen supera la capacidad de la red, hay OVERFLOW.
 * Si cabe, busca la altura con dividir y vencerás: el rango es
 * [base más baja, techo más alto] y en cada paso se parte a la
 * mitad. En la cota media se calculan los m³; si faltan, se sigue
 * en la mitad de arriba; si sobran, en la de abajo.
 *
 * Complejidad: O(n · k), n cisternas y k ≈ 80 cortes del rango.
 *
 * @param volumen - agua a inyectar, en m³.
 * @returns resultado con Overflow, Nivel, Capacidad y VolumenPed.
 */
func (r Red) NivelAgua(volumen float64) ResultadoNivel {
	capacidad := r.VolumenTotal()
	resultado := ResultadoNivel{
		Capacidad:  capacidad,
		VolumenPed: volumen,
	}

	if len(r.Cisternas) == 0 {
		resultado.Overflow = volumen > epsilon
		return resultado
	}

	if volumen > capacidad+epsilon {
		resultado.Overflow = true
		return resultado
	}

	baja, alta := r.rangoAlturas()

	if volumen <= epsilon {
		resultado.Nivel = baja
		return resultado
	}

	resultado.Nivel = r.buscarNivel(baja, alta, volumen, iteracionesBinaria)
	return resultado
}

/**
 * rangoAlturas obtiene el intervalo de búsqueda de la red.
 *
 * @returns cota mínima (base más baja) y cota máxima (mayor Base + Alto).
 */
func (r Red) rangoAlturas() (float64, float64) {
	return r.baseMinima(), r.techoMaximo()
}

/**
 * buscarNivel aplica dividir y vencerás sobre [baja, alta].
 * En la mitad evalúa VolumenHasta; según falte o sobre agua,
 * reduce el rango a la mitad superior o a la inferior.
 *
 * @param baja - extremo inferior del rango, en metros.
 * @param alta - extremo superior del rango, en metros.
 * @param volumen - m³ que se quieren almacenar.
 * @param restantes - cortes que quedan (caso base: 0 o rango ya mínimo).
 * @returns cota estimada del agua, en metros.
 */
func (r Red) buscarNivel(baja float64, alta float64, volumen float64, restantes int) float64 {
	if restantes <= 0 || alta-baja <= epsilon {
		return (baja + alta) / 2
	}

	mitad := (baja + alta) / 2
	cubicos := r.VolumenHasta(mitad)

	if cubicos < volumen {
		return r.buscarNivel(mitad, alta, volumen, restantes-1)
	}

	return r.buscarNivel(baja, mitad, volumen, restantes-1)
}

/**
 * EtapaLlenado es un tramo de cota donde el agua se comporta igual:
 * o llena un conjunto fijo de cisternas, o cruza un hueco por tuberías.
 *
 * Si hay más de un índice, esas cisternas están lado a lado en esa
 * cota y se llenan al mismo tiempo (se reparte el agua según su área).
 */
type EtapaLlenado struct {
	Desde   float64
	Hasta   float64
	Indices []int
	Volumen float64
	Hueco   bool
}

/**
 * LadoALado indica si en esta cota hay dos o más cisternas
 * compartiendo el agua al mismo tiempo.
 *
 * @returns true si no es hueco y hay más de un índice.
 */
func (e EtapaLlenado) LadoALado() bool {
	return !e.Hueco && len(e.Indices) > 1
}

/**
 * OrdenLlenado recorre las cotas de abajo hacia arriba y describe cómo
 * baja el agua: primero las más bajas; las que coinciden en altura se
 * llenan juntas; si hay un hueco, el sobrante sigue por las tuberías.
 *
 * @param volumen - agua a inyectar, en m³.
 * @returns etapas de llenado en orden ascendente de cota.
 */
func (r Red) OrdenLlenado(volumen float64) []EtapaLlenado {
	if len(r.Cisternas) == 0 || volumen <= epsilon {
		return nil
	}

	tope := volumen
	capacidad := r.VolumenTotal()
	if tope > capacidad {
		tope = capacidad
	}

	eventos := compactarEventos(r.eventosAltura())
	areaActiva := 0.0
	acumulado := 0.0
	etapas := make([]EtapaLlenado, 0)

	for i := 0; i+1 < len(eventos); i++ {
		areaActiva += eventos[i].dA
		z0 := eventos[i].z
		z1 := eventos[i+1].z
		dz := z1 - z0

		if dz <= epsilon {
			continue
		}

		if areaActiva <= epsilon {
			if acumulado > epsilon && acumulado < tope-epsilon {
				etapas = append(etapas, EtapaLlenado{
					Desde: z0,
					Hasta: z1,
					Hueco: true,
				})
			}
			continue
		}

		if acumulado >= tope-epsilon {
			break
		}

		losas := areaActiva * dz
		usado := losas
		hasta := z1

		if acumulado+losas >= tope-epsilon {
			falta := tope - acumulado
			if falta < 0 {
				falta = 0
			}
			usado = falta
			hasta = z0 + falta/areaActiva
		}

		etapas = append(etapas, EtapaLlenado{
			Desde:   z0,
			Hasta:   hasta,
			Indices: r.cisternasEnCota((z0 + hasta) / 2),
			Volumen: usado,
		})

		acumulado += usado
		if acumulado >= tope-epsilon {
			break
		}
	}

	return etapas
}

/**
 * cisternasEnCota busca las cisternas que ocupan una altura dada:
 * están lado a lado en esa cota.
 *
 * @param z - cota absoluta, en metros.
 * @returns índices 0-based de las cisternas que contienen esa cota.
 */
func (r Red) cisternasEnCota(z float64) []int {
	indices := make([]int, 0)
	for i, c := range r.Cisternas {
		if c.Base <= z+epsilon && z < c.Techo()+epsilon {
			indices = append(indices, i)
		}
	}
	return indices
}

/**
 * indicesPorBase ordena las cisternas de abajo hacia arriba según
 * la distancia del piso al embase.
 *
 * @returns índices 0-based ordenados por Base ascendente.
 */
func (r Red) indicesPorBase() []int {
	orden := make([]int, len(r.Cisternas))
	for i := range orden {
		orden[i] = i
	}

	sort.SliceStable(orden, func(i, j int) bool {
		return r.Cisternas[orden[i]].Base < r.Cisternas[orden[j]].Base
	})

	return orden
}

/**
 * eventosAltura genera un +área en cada base y un −área en cada techo.
 *
 * @returns 2n eventos, sin ordenar ni fusionar.
 */
func (r Red) eventosAltura() []evento {
	eventos := make([]evento, 0, 2*len(r.Cisternas))
	for _, c := range r.Cisternas {
		area := c.AreaBase()
		eventos = append(eventos, evento{z: c.Base, dA: area})
		eventos = append(eventos, evento{z: c.Techo(), dA: -area})
	}
	return eventos
}

/**
 * compactarEventos ordena los eventos por cota y fusiona los que
 * caen en la misma altura.
 *
 * @param eventos - cambios de área en bases y techos.
 * @returns eventos ordenados, con un solo delta por cota.
 */
func compactarEventos(eventos []evento) []evento {
	sort.Slice(eventos, func(i, j int) bool {
		return eventos[i].z < eventos[j].z
	})

	compactos := make([]evento, 0, len(eventos))
	for _, e := range eventos {
		if len(compactos) > 0 && abs(compactos[len(compactos)-1].z-e.z) <= epsilon {
			compactos[len(compactos)-1].dA += e.dA
			continue
		}
		compactos = append(compactos, e)
	}

	return compactos
}

/**
 * baseMinima busca la Base más baja de la red: la distancia mínima
 * del piso al embase de alguna cisterna.
 *
 * @returns cota mínima, en metros. Requiere al menos una cisterna.
 */
func (r Red) baseMinima() float64 {
	minBase := r.Cisternas[0].Base
	for _, c := range r.Cisternas[1:] {
		if c.Base < minBase {
			minBase = c.Base
		}
	}
	return minBase
}

/**
 * techoMaximo busca el mayor Base + Alto de la red: hasta dónde
 * puede llegar el agua si se llena toda la capacidad.
 *
 * @returns cota máxima, en metros. Requiere al menos una cisterna.
 */
func (r Red) techoMaximo() float64 {
	maxTecho := r.Cisternas[0].Techo()
	for _, c := range r.Cisternas[1:] {
		if t := c.Techo(); t > maxTecho {
			maxTecho = t
		}
	}
	return maxTecho
}

/**
 * abs calcula el valor absoluto de un número.
 *
 * @param x - valor a convertir.
 * @returns x si es ≥ 0; −x si es negativo.
 */
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
