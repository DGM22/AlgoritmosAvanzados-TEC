# Misión 4 — Red de cisternas

Programa en Go que simula cisternas **conectadas por tuberías**. El agua
fluye hacia las de **abajo**; las que quedan por encima del nivel se
quedan **vacías**. Si el volumen supera la capacidad de la red, imprime
**OVERFLOW**.

El nivel se busca con **dividir y vencerás**: rango `[base más baja,
Base+Alto más alto]`, se parte a la mitad y en esa cota se calculan los
m³ hasta acotar la altura.

## Cómo correr

Go instalado. Desde `mision-4`:

```bash
go run .
```

Pide, en metros: número de cisternas; para cada una ancho, largo, alto
y altura de la base respecto al piso; al final los m³ de agua. Devuelve
el nivel (o OVERFLOW) y el **porcentaje de llenado** de cada cisterna.

Tests:

```bash
go test .
```
