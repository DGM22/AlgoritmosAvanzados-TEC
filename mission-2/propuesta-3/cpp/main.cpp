// Propuesta #3: tabla hash con encadenamiento externo.
#include <chrono>
#include <iostream>
#include <string>

#include "grades.hpp"

namespace {

// Orquesta el flujo de lectura, procesamiento y resultado global.
class GradeModeApplication {
public:
    void run() {
        printHeader();

        std::string fileName;
        while (promptForFileName(fileName)) {
            processFile(fileName);
        }

        printGlobalResult();
    }

private:
    void printHeader() const {
        std::cout << "=== Calificación más común (Propuesta #3: tabla hash con encadenamiento externo) ===\n";
    }

    /**
     * promptForFileName pide un nombre de archivo por teclado.
     * @param outFileName - variable de salida con el nombre ingresado.
     * @returns false si el usuario ya no quiere dar más archivos.
     */
    bool promptForFileName(std::string& outFileName) const {
        std::cout << "\nIngrese el nombre del archivo (Enter para terminar): ";
        std::getline(std::cin, outFileName);
        return !outFileName.empty();
    }

    /**
     * processFile lee, procesa y muestra el resultado de un archivo.
     * @param fileName - ruta o nombre del archivo a procesar.
     */
    void processFile(const std::string& fileName) {
        const auto start = std::chrono::steady_clock::now();

        grades::HashTable fileTable;
        const auto result = reader_.readInto(fileName, fileTable);

        if (!result.success) {
            std::cout << "Error: " << result.errorMessage << "\n";
            return;
        }

        if (fileTable.isEmpty()) {
            std::cout << "No se encontraron calificaciones válidas en '" << fileName << "'.\n";
            return;
        }

        const auto mode = fileTable.mostCommonGrade();
        const auto elapsedNanos = std::chrono::duration_cast<std::chrono::nanoseconds>(
                                       std::chrono::steady_clock::now() - start)
                                       .count();

        if (result.linesSkipped > 0) {
            std::cout << "Aviso: se ignoraron " << result.linesSkipped
                      << " línea(s) inválida(s) en '" << fileName << "'.\n";
        }
        std::cout << analyzer_.describeFile(fileName, fileTable.totalCount(), mode, elapsedNanos) << "\n";

        globalTable_.mergeFrom(fileTable);
        totalElapsedNanos_ += elapsedNanos;
        ++filesProcessed_;
    }

    // printGlobalResult muestra el resultado acumulado de todos los archivos.
    void printGlobalResult() const {
        std::cout << "\n";
        if (filesProcessed_ == 0 || globalTable_.isEmpty()) {
            std::cout << "No se procesó ningún archivo válido. Fin del programa.\n";
            return;
        }

        const auto modeStart = std::chrono::steady_clock::now();
        const auto globalMode = globalTable_.mostCommonGrade();
        const auto modeElapsedNanos = std::chrono::duration_cast<std::chrono::nanoseconds>(
                                           std::chrono::steady_clock::now() - modeStart)
                                           .count();

        std::cout << analyzer_.describeGlobal(globalTable_.totalCount(), globalMode, modeElapsedNanos,
                                               totalElapsedNanos_)
                  << "\n";
    }

    grades::FileReader reader_;
    grades::Analyzer analyzer_;
    grades::HashTable globalTable_;
    int filesProcessed_ = 0;
    std::int64_t totalElapsedNanos_ = 0;
};

} // namespace

int main() {
    GradeModeApplication app;
    app.run();
    return 0;
}
