package main

import "testing"

/**
 * casiIgual falla el test si |got − want| supera la tolerancia.
 *
 * @param t - contexto del test.
 * @param got - valor obtenido.
 * @param want - valor esperado.
 * @param tol - diferencia máxima aceptada.
 * @param msg - etiqueta del aserto en el mensaje de error.
 */
func casiIgual(t *testing.T, got, want, tol float64, msg string) {
	t.Helper()
	if abs(got-want) > tol {
		t.Fatalf("%s: got=%.6f want=%.6f", msg, got, want)
	}
}

func TestCisternaVolumenHasta(t *testing.T) {
	c := Cisterna{Ancho: 2, Largo: 3, Alto: 4, Base: 1}

	casiIgual(t, c.Volumen(), 24, 1e-9, "capacidad")
	casiIgual(t, c.VolumenHasta(0.5), 0, 1e-9, "debajo de la base")
	casiIgual(t, c.VolumenHasta(1), 0, 1e-9, "justo en la base")
	casiIgual(t, c.VolumenHasta(3), 12, 1e-9, "media cisterna")
	casiIgual(t, c.VolumenHasta(5), 24, 1e-9, "llena")
	casiIgual(t, c.VolumenHasta(10), 24, 1e-9, "arriba del techo")
}

func TestNivelUnaCisterna(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 2, Largo: 3, Alto: 4, Base: 0},
	}}

	r := red.NivelAgua(12)
	if r.Overflow {
		t.Fatal("no debería desbordar")
	}
	casiIgual(t, r.Nivel, 2, 1e-6, "media cisterna")
}

func TestNivelDosCisternasLadoALado(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 2, Largo: 2, Alto: 5, Base: 0},
		{Ancho: 3, Largo: 1, Alto: 5, Base: 0},
	}}

	r := red.NivelAgua(14)
	if r.Overflow {
		t.Fatal("no debería desbordar")
	}
	casiIgual(t, r.Nivel, 2, 1e-6, "área conjunta 7 m²")
}

func TestNivelConHueco(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 10, Largo: 10, Alto: 10, Base: 0},
		{Ancho: 10, Largo: 10, Alto: 10, Base: 20},
	}}

	exacta := red.NivelAgua(1000)
	if exacta.Overflow {
		t.Fatal("1000 m³ caben justo en la inferior")
	}
	casiIgual(t, exacta.Nivel, 10, 1e-6, "techo de la inferior")

	mediaAlta := red.NivelAgua(1500)
	if mediaAlta.Overflow {
		t.Fatal("1500 m³ caben")
	}
	casiIgual(t, mediaAlta.Nivel, 25, 1e-6, "mitad de la superior")
}

func TestOverflow(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 1, Largo: 1, Alto: 1, Base: 0},
	}}

	r := red.NivelAgua(1.01)
	if !r.Overflow {
		t.Fatalf("debería desbordar, nivel=%f", r.Nivel)
	}
}

func TestOverflowSuperaCapacidadDeLaRed(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 10, Largo: 10, Alto: 10, Base: 0},
		{Ancho: 10, Largo: 10, Alto: 10, Base: 20},
	}}

	capacidad := red.VolumenTotal()
	casiIgual(t, capacidad, 2000, 1e-9, "capacidad de la red")

	cabe := red.NivelAgua(capacidad)
	if cabe.Overflow {
		t.Fatal("el volumen exacto de la red no es overflow")
	}

	demas := red.NivelAgua(capacidad + 0.01)
	if !demas.Overflow {
		t.Fatal("si se pasa la capacidad de todas las cisternas debe haber OVERFLOW")
	}
}

func TestCapacidadExacta(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 2, Largo: 2, Alto: 3, Base: 5},
	}}

	r := red.NivelAgua(12)
	if r.Overflow {
		t.Fatal("volumen exacto no desborda")
	}
	casiIgual(t, r.Nivel, 8, 1e-6, "techo = base + alto")
}

func TestVolumenCero(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 1, Largo: 1, Alto: 2, Base: 4},
	}}

	r := red.NivelAgua(0)
	if r.Overflow {
		t.Fatal("volumen 0 no desborda")
	}
	casiIgual(t, r.Nivel, 4, 1e-9, "nivel en la base más baja")
}

func TestPorcentajeLlenadoAbajoHaciaArriba(t *testing.T) {
	baja := Cisterna{Ancho: 10, Largo: 10, Alto: 10, Base: 0}
	alta := Cisterna{Ancho: 10, Largo: 10, Alto: 10, Base: 20}
	red := Red{Cisternas: []Cisterna{baja, alta}}

	r := red.NivelAgua(1500)
	if r.Overflow {
		t.Fatal("1500 m³ caben")
	}

	casiIgual(t, baja.PorcentajeLlenado(r.Nivel), 100, 1e-6, "la de abajo se llena primero")
	casiIgual(t, alta.PorcentajeLlenado(r.Nivel), 50, 1e-6, "la de arriba queda a la mitad")
}

func TestPorcentajeOverflowEsCien(t *testing.T) {
	c := Cisterna{Ancho: 1, Largo: 1, Alto: 2, Base: 3}
	casiIgual(t, c.PorcentajeLlenado(c.Techo()), 100, 1e-9, "techo = 100%")
	casiIgual(t, c.PorcentajeLlenado(c.Base), 0, 1e-9, "base = 0%")
}

func TestOrdenLadoALado(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 2, Largo: 2, Alto: 5, Base: 0},
		{Ancho: 3, Largo: 1, Alto: 5, Base: 0},
	}}

	etapas := red.OrdenLlenado(14)
	if len(etapas) != 1 {
		t.Fatalf("una sola etapa, got=%d", len(etapas))
	}
	if !etapas[0].LadoALado() {
		t.Fatal("deben llenarse lado a lado")
	}
	if len(etapas[0].Indices) != 2 {
		t.Fatalf("indices=%v", etapas[0].Indices)
	}
}

func TestOrdenAbajoLuegoArriba(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 10, Largo: 10, Alto: 10, Base: 0},
		{Ancho: 10, Largo: 10, Alto: 10, Base: 20},
	}}

	etapas := red.OrdenLlenado(1500)
	if len(etapas) != 3 {
		t.Fatalf("llenar abajo, hueco y arriba; got=%d etapas", len(etapas))
	}
	if etapas[0].Hueco || len(etapas[0].Indices) != 1 || etapas[0].Indices[0] != 0 {
		t.Fatalf("primera etapa debe ser la cisterna baja: %+v", etapas[0])
	}
	if !etapas[1].Hueco {
		t.Fatal("entre 10 y 20 debe haber hueco de tuberías")
	}
	if etapas[2].Hueco || len(etapas[2].Indices) != 1 || etapas[2].Indices[0] != 1 {
		t.Fatalf("tercera etapa debe ser la cisterna alta: %+v", etapas[2])
	}
}

func TestOrdenSolapadasParcialmente(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 1, Largo: 1, Alto: 10, Base: 0},
		{Ancho: 1, Largo: 1, Alto: 10, Base: 5},
	}}

	etapas := red.OrdenLlenado(17)
	if len(etapas) != 3 {
		t.Fatalf("sola, lado a lado, sola; got=%d", len(etapas))
	}
	if etapas[0].LadoALado() || etapas[0].Indices[0] != 0 {
		t.Fatalf("0–5 solo cisterna 1: %+v", etapas[0])
	}
	if !etapas[1].LadoALado() {
		t.Fatalf("5–10 lado a lado: %+v", etapas[1])
	}
	if etapas[2].LadoALado() || etapas[2].Indices[0] != 1 {
		t.Fatalf("10–15 solo cisterna 2: %+v", etapas[2])
	}
}

func TestRangoYBusquedaBinaria(t *testing.T) {
	red := Red{Cisternas: []Cisterna{
		{Ancho: 10, Largo: 10, Alto: 10, Base: 0},
		{Ancho: 10, Largo: 10, Alto: 10, Base: 20},
	}}

	baja, alta := red.rangoAlturas()
	casiIgual(t, baja, 0, 1e-9, "altura mínima")
	casiIgual(t, alta, 30, 1e-9, "altura máxima")

	r := red.NivelAgua(1500)
	if r.Overflow {
		t.Fatal("1500 m³ caben")
	}
	casiIgual(t, r.Nivel, 25, 1e-6, "mitad por búsqueda binaria")
	casiIgual(t, red.VolumenHasta(r.Nivel), 1500, 1e-4, "los m³ de esa altura coinciden")
}
