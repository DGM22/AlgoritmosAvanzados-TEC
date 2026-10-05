INF = float("inf")
import time


# /**
#  * cargar_archivo lee ciudades y vuelos desde un .txt.
#  * Primero vienen los nombres de ciudad, luego un "#",
#  * y después líneas "ciudad1 ciudad2 costo".
#  *
#  * @param nombre_archivo - ruta del archivo de vuelos.
#  * @returns lista de ciudades y lista de tuplas (origen, destino, costo).
#  */
def cargar_archivo(nombre_archivo):
    ciudades = []
    vuelos = []

    leyendo_ciudades = True

    with open(nombre_archivo, "r", encoding="utf-8") as archivo:

        for linea in archivo:

            linea = linea.strip()

            if linea == "":
                continue

            if leyendo_ciudades:

                if linea == "#":
                    leyendo_ciudades = False
                    continue

                ciudades.append(linea)

            else:
                partes = linea.split()

                if len(partes) != 3:
                    print("Línea inválida:", linea)
                    continue

                ciudad1 = partes[0]
                ciudad2 = partes[1]
                costo = int(partes[2])

                vuelos.append(
                    (ciudad1, ciudad2, costo)
                )

    return ciudades, vuelos


# /**
#  * crear_matrices arma las matrices de costo y de siguiente vuelo.
#  * Los vuelos se tratan como bidireccionales. Si hay dos vuelos
#  * entre el mismo par, se queda el de menor costo.
#  *
#  * @param ciudades - nombres de ciudad en el orden del archivo.
#  * @param vuelos - tuplas (ciudad1, ciudad2, costo).
#  * @returns dist, siguiente e indice (nombre en minúsculas -> posición).
#  */
def crear_matrices(ciudades, vuelos):
    n = len(ciudades)

    indice = {}

    for i in range(n):
        indice[ciudades[i].lower()] = i

    dist = []
    siguiente = []

    for i in range(n):

        fila_dist = []
        fila_siguiente = []

        for j in range(n):

            if i == j:
                fila_dist.append(0)
                fila_siguiente.append(i)

            else:
                fila_dist.append(INF)
                fila_siguiente.append(None)

        dist.append(fila_dist)
        siguiente.append(fila_siguiente)

    for ciudad1, ciudad2, costo in vuelos:

        ciudad1_lower = ciudad1.lower()
        ciudad2_lower = ciudad2.lower()

        if ciudad1_lower not in indice:
            print("Ciudad no encontrada:", ciudad1)
            continue

        if ciudad2_lower not in indice:
            print("Ciudad no encontrada:", ciudad2)
            continue

        i = indice[ciudad1_lower]
        j = indice[ciudad2_lower]

        if costo < dist[i][j]:

            dist[i][j] = costo
            dist[j][i] = costo

            siguiente[i][j] = j
            siguiente[j][i] = i

    return dist, siguiente, indice


# /**
#  * floyd_warshall calcula el costo mínimo entre todos los pares.
#  * Para cada posible escala k, prueba si i -> k -> j es más barato
#  * que el camino actual i -> j y actualiza `siguiente` para reconstruir.
#  *
#  * Complejidad: Theta(n³).
#  *
#  * @param dist - matriz de costos; se modifica in-place.
#  * @param siguiente - matriz del siguiente vértice; se modifica in-place.
#  */
def floyd_warshall(dist, siguiente):
    n = len(dist)

    for k in range(n):

        for i in range(n):

            for j in range(n):

                if dist[i][k] == INF or dist[k][j] == INF:
                    continue

                nuevo_costo = (
                    dist[i][k]
                    + dist[k][j]
                )

                if nuevo_costo < dist[i][j]:

                    dist[i][j] = nuevo_costo
                    siguiente[i][j] = siguiente[i][k]


# /**
#  * reconstruir_ruta arma el camino óptimo de origen a destino
#  * usando la matriz `siguiente` de Floyd-Warshall.
#  *
#  * @param origen - índice de la ciudad de origen.
#  * @param destino - índice de la ciudad de destino.
#  * @param siguiente - matriz del siguiente vértice en la ruta óptima.
#  * @returns lista de índices de la ruta; vacía si no hay camino.
#  */
def reconstruir_ruta(origen, destino, siguiente):
    if siguiente[origen][destino] is None:
        return []

    ruta = [origen]

    actual = origen

    while actual != destino:

        actual = siguiente[actual][destino]

        ruta.append(actual)

    return ruta


# /**
#  * ruta_a_texto convierte índices de ciudad en una cadena "A -> B -> C".
#  *
#  * @param ruta - índices de las ciudades en orden.
#  * @param ciudades - nombres indexados igual que las matrices.
#  * @returns ruta formateada para imprimir.
#  */
def ruta_a_texto(ruta, ciudades):
    nombres = []

    for indice in ruta:
        nombres.append(ciudades[indice])

    return " -> ".join(nombres)


# /**
#  * pedir_ciudad lee un nombre de ciudad del usuario y lo busca
#  * en el índice (sin distinguir mayúsculas).
#  *
#  * @param mensaje - texto del input.
#  * @param indice - mapa nombre en minúsculas -> posición en la matriz.
#  * @returns índice de la ciudad, o None si no existe.
#  */
def pedir_ciudad(mensaje, indice):
    ciudad = input(mensaje).strip().lower()

    if ciudad not in indice:

        print("La ciudad no existe.")

        return None

    return indice[ciudad]


# /**
#  * opcion_1 muestra la ruta más económica entre dos ciudades,
#  * su costo y el número de escalas.
#  *
#  * @param ciudades - nombres de ciudad.
#  * @param dist - matriz de costos mínimo ya calculada.
#  * @param siguiente - matriz para reconstruir caminos.
#  * @param indice - mapa nombre -> posición.
#  */
def opcion_1(ciudades, dist, siguiente, indice):
    origen = pedir_ciudad(
        "Ciudad origen: ",
        indice
    )

    if origen is None:
        return

    destino = pedir_ciudad(
        "Ciudad destino: ",
        indice
    )

    tiempo_inicio = time.time()
    if destino is None:
        return

    if dist[origen][destino] == INF:

        print(
            "No existe una ruta entre esas ciudades."
        )

        return

    ruta = reconstruir_ruta(
        origen,
        destino,
        siguiente
    )

    escalas = max(0, len(ruta) - 2)

    print()
    print(
        "Costo mínimo:",
        f"${dist[origen][destino]}"
    )

    print(
        "Ruta:",
        ruta_a_texto(ruta, ciudades)
    )

    print(
        "Número de escalas:",
        escalas
    )

    tiempo_fin = time.time()
    tiempo_ejecucion = tiempo_fin - tiempo_inicio
    print(f"Tiempo de ejecución: {tiempo_ejecucion:.6f} segundos")


# /**
#  * opcion_2 lista todos los destinos alcanzables desde un origen,
#  * con costo y ruta óptima de cada uno.
#  *
#  * @param ciudades - nombres de ciudad.
#  * @param dist - matriz de costos mínimo ya calculada.
#  * @param siguiente - matriz para reconstruir caminos.
#  * @param indice - mapa nombre -> posición.
#  */
def opcion_2(ciudades, dist, siguiente, indice):
    origen = pedir_ciudad(
        "Ciudad origen: ",
        indice
    )

    if origen is None:
        return

    print()
    print(
        f"Destinos alcanzables desde {ciudades[origen]}:"
    )

    encontrado = False
    tiempo_inicio = time.time()
    for destino in range(len(ciudades)):

        if origen == destino:
            continue

        if dist[origen][destino] == INF:
            continue

        encontrado = True

        ruta = reconstruir_ruta(
            origen,
            destino,
            siguiente
        )

        print()
        print(
            f"Destino: {ciudades[destino]}"
        )

        print(
            f"Costo: ${dist[origen][destino]}"
        )

        print(
            "Ruta:",
            ruta_a_texto(ruta, ciudades)
        )

    if not encontrado:

        print(
            "No hay ciudades alcanzables desde este origen."
        )

    tiempo_fin = time.time()
    tiempo_ejecucion = tiempo_fin - tiempo_inicio
    print(f"Tiempo de ejecución: {tiempo_ejecucion:.6f} segundos")


# /**
#  * opcion_3 muestra las rutas óptimas (i < j, sin duplicar ida y
#  * vuelta) que tienen exactamente el número de escalas pedido.
#  *
#  * @param ciudades - nombres de ciudad.
#  * @param dist - matriz de costos mínimo ya calculada.
#  * @param siguiente - matriz para reconstruir caminos.
#  */
def opcion_3(ciudades, dist, siguiente):
    try:

        numero_escalas = int(
            input("Número de escalas: ")
        )
    

    except ValueError:

        print(
            "Debes introducir un número entero."
        )

        return

    tiempo_inicio = time.time()

    if numero_escalas < 0:

        print(
            "El número de escalas no puede ser negativo."
        )

        return

    print()
    print(
        f"Rutas óptimas con {numero_escalas} escalas:"
    )

    encontrado = False

    for i in range(len(ciudades)):

        for j in range(i + 1, len(ciudades)):

            if dist[i][j] == INF:
                continue

            ruta = reconstruir_ruta(
                i,
                j,
                siguiente
            )

            escalas = len(ruta) - 2

            if escalas == numero_escalas:

                encontrado = True

                print()
                print(
                    ruta_a_texto(
                        ruta,
                        ciudades
                    )
                )

                print(
                    f"Costo: ${dist[i][j]}"
                )

    if not encontrado:

        print(
            "No existen rutas óptimas "
            "con ese número de escalas."
        )
    
    tiempo_fin = time.time()
    tiempo_ejecucion = tiempo_fin - tiempo_inicio
    print(f"Tiempo de ejecución: {tiempo_ejecucion:.6f} segundos")
   


# /**
#  * main carga el archivo de vuelos, ejecuta Floyd-Warshall una
#  * sola vez y abre el menú de consultas.
#  */
def main():
    eleccion_archivo = input(
        """
        Selecciona el archivo de vuelos:
        1. vuelos.txt
        2. vuelos_100.txt
        """
    )

    if eleccion_archivo == "1":
        nombre_archivo = "vuelos.txt"
    elif eleccion_archivo == "2":
        nombre_archivo = "vuelos_100.txt"
    else:
        print("Opción inválida.")
        return

    try:

        ciudades, vuelos = cargar_archivo(
            nombre_archivo
        )

    except FileNotFoundError:

        print(
            "No se encontró el archivo."
        )

        return

    dist, siguiente, indice = crear_matrices(
        ciudades,
        vuelos
    )

    floyd_warshall(
        dist,
        siguiente
    )

    print()
    print(
        f"Se cargaron {len(ciudades)} ciudades."
    )

    while True:

        print()
        print("============================")
        print("      SISTEMA AEROLÍNEA")
        print("============================")
        print()
        print("1. Ruta más económica")
        print("2. Destinos desde una ciudad")
        print("3. Rutas por número de escalas")
        print("0. Salir")

        opcion = input(
            "\nSelecciona una opción: "
        )

        if opcion == "1":

            opcion_1(
                ciudades,
                dist,
                siguiente,
                indice
            )

        elif opcion == "2":

            opcion_2(
                ciudades,
                dist,
                siguiente,
                indice
            )

        elif opcion == "3":

            opcion_3(
                ciudades,
                dist,
                siguiente
            )


        elif opcion == "0":

            print(
                "Programa terminado."
            )

            break

        else:

            print(
                "Opción inválida."
            )


if __name__ == "__main__":
    main()
