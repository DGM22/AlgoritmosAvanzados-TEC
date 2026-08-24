// Propuesta #2: vector ordenado + contador único.
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
        std::cout << "=== Calificación más común (Propuesta #2: vector ordenado + contador único) ===\n";
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
     * processFile lee, ordena, procesa y muestra el resultado de un archivo.
     * @param fileName - ruta o nombre del archivo a procesar.
     */
    void processFile(const std::string& fileName) {
        const auto start = std::chrono::steady_clock::now();

        grades::GradeVector fileVector;
        const auto result = reader_.readInto(fileName, fileVector);

        if (!result.success) {
            std::cout << "Error: " << result.errorMessage << "\n";
            return;
        }

        if (fileVector.isEmpty()) {
            std::cout << "No se encontraron calificaciones válidas en '" << fileName << "'.\n";
            return;
        }

        fileVector.sort();
        const auto mode = fileVector.mostCommonGrade();
        const auto elapsedNanos = std::chrono::duration_cast<std::chrono::nanoseconds>(
                                       std::chrono::steady_clock::now() - start)
                                       .count();

        if (result.linesSkipped > 0) {
            std::cout << "Aviso: se ignoraron " << result.linesSkipped
                      << " línea(s) inválida(s) en '" << fileName << "'.\n";
        }
        std::cout << analyzer_.describeFile(fileName, fileVector.totalCount(), mode, elapsedNanos) << "\n";

        globalVector_.merge(fileVector);
        totalElapsedNanos_ += elapsedNanos;
        ++filesProcessed_;
    }

    // printGlobalResult ordena el vector global (necesario) y muestra el resultado acumulado.
    void printGlobalResult() {
        std::cout << "\n";
        if (filesProcessed_ == 0 || globalVector_.isEmpty()) {
            std::cout << "No se procesó ningún archivo válido. Fin del programa.\n";
            return;
        }

        const auto modeStart = std::chrono::steady_clock::now();
        globalVector_.sort();
        const auto globalMode = globalVector_.mostCommonGrade();
        const auto modeElapsedNanos = std::chrono::duration_cast<std::chrono::nanoseconds>(
                                           std::chrono::steady_clock::now() - modeStart)
                                           .count();

        std::cout << analyzer_.describeGlobal(globalVector_.totalCount(), globalMode, modeElapsedNanos,
                                               totalElapsedNanos_)
                  << "\n";
    }

    grades::FileReader reader_;
    grades::Analyzer analyzer_;
    grades::GradeVector globalVector_;
    int filesProcessed_ = 0;
    std::int64_t totalElapsedNanos_ = 0;
};

} // namespace

int main() {
    GradeModeApplication app;
    app.run();
    return 0;
}
