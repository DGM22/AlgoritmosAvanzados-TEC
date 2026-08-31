#!/bin/sh
set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

echo "==> Usando los .txt ya existentes en datos/"
echo
echo "==> Merge Sort 2"
go run ./mergeSort/mergesort2.go ./mergeSort/io.go

echo
echo "==> Merge Sort 3"
go run ./mergeSort/mergesort3.go ./mergeSort/io.go

echo
echo "==> Quick Sort 2"
go run ./quickSort/quicksort2.go ./quickSort/io.go

echo
echo "==> Quick Sort 3"
go run ./quickSort/quicksort3.go ./quickSort/io.go
