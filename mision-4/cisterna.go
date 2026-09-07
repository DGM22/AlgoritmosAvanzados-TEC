package main

/**
 * Cisterna es un depósito rectangular dentro de la red.
 *
 * Ancho, Largo y Alto son las dimensiones internas.
 * Base es la altura del piso de la cisterna respecto al suelo.
 * El techo queda en Base + Alto. Todas las medidas están en metros.
 */
type Cisterna struct {
	Ancho float64
	Largo float64
	Alto  float64
	Base  float64
}

/**
 * AreaBase calcula el área de la sección horizontal (m²).
 * En un intervalo de altura donde la cisterna está presente,
 * el volumen crece linealmente con esta área.
 *
 * @returns área ancho × largo, en m².
 */
func (c Cisterna) AreaBase() float64 {
	return c.Ancho * c.Largo
}

/**
 * Volumen calcula la capacidad total de la cisterna.
 *
 * @returns capacidad en m³ (área de la base × alto).
 */
func (c Cisterna) Volumen() float64 {
	return c.AreaBase() * c.Alto
}

/**
 * Techo calcula la cota absoluta del borde superior respecto al piso.
 * Equivale a Base + Alto: distancia del suelo al embase más la
 * altura interna del depósito.
 *
 * @returns altura del borde superior, en metros.
 */
func (c Cisterna) Techo() float64 {
	return c.Base + c.Alto
}

/**
 * VolumenHasta calcula cuánta agua cabe en esta cisterna si el
 * nivel libre de la red está en una cota dada.
 *
 * @param nivel - cota del agua respecto al piso de la red, en metros.
 * @returns volumen almacenado en esta cisterna, en m³ (0 si el
 *          nivel no alcanza la base; capacidad plena si la rebasa).
 */
func (c Cisterna) VolumenHasta(nivel float64) float64 {
	if nivel <= c.Base {
		return 0
	}

	columna := nivel - c.Base
	if columna > c.Alto {
		columna = c.Alto
	}

	return columna * c.AreaBase()
}

/**
 * ColumnaAgua calcula la altura de agua dentro de la cisterna.
 *
 * @param nivel - cota del agua respecto al piso de la red, en metros.
 * @returns columna de agua en metros (0 si el nivel no la alcanza,
 *          Alto si ya rebosó esta unidad).
 */
func (c Cisterna) ColumnaAgua(nivel float64) float64 {
	if nivel <= c.Base {
		return 0
	}

	columna := nivel - c.Base
	if columna > c.Alto {
		return c.Alto
	}

	return columna
}

/**
 * PorcentajeLlenado calcula el llenado de esta cisterna según el
 * nivel absoluto del agua en la red.
 *
 * @param nivel - cota del agua respecto al piso de la red, en metros.
 * @returns porcentaje de 0 a 100.
 */
func (c Cisterna) PorcentajeLlenado(nivel float64) float64 {
	if c.Alto <= 0 {
		return 0
	}

	return 100.0 * c.ColumnaAgua(nivel) / c.Alto
}
