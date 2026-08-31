#!/bin/sh
set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

echo "==> Usando los .txt de n=50000 (aleatorio, casi_ordenado, duplicados)"
echo

echo "==> Merge Sort 2"
go run ./mergeSort/mergesort2.go ./mergeSort/io.go -n 50000

echo
echo "==> Merge Sort 3"
go run ./mergeSort/mergesort3.go ./mergeSort/io.go -n 50000

echo
echo "==> Quick Sort 2"
go run ./quickSort/quicksort2.go ./quickSort/io.go -n 50000

echo
echo "==> Quick Sort 3"
go run ./quickSort/quicksort3.go ./quickSort/io.go -n 50000
