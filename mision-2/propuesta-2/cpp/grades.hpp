#pragma once

// Propuesta #2: vector dinámico + std::sort + contador único (header-only; main.cpp orquesta).
#include <algorithm>
#include <cstdint>
#include <fstream>
#include <iomanip>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

namespace grades {

constexpr int kMinGrade = 1;
constexpr int kMaxGrade = 100;

namespace number_format {

/**
 * withThousandsSeparator agrega separador de miles a un entero.
 * @param value - valor a formatear (p. ej. 611124 -> "611,124").
 */
inline std::string withThousandsSeparator(std::int64_t value) {
    const bool isNegative = value < 0;
    std::string digits = std::to_string(isNegative ? -value : value);

    std::string grouped;
    grouped.reserve(digits.size() + digits.size() / 3);

    int digitsSinceSeparator = 0;
    for (auto it = digits.rbegin(); it != digits.rend(); ++it) {
        if (digitsSinceSeparator == 3) {
            grouped.push_back(',');
            digitsSinceSeparator = 0;
        }
        grouped.push_back(*it);
        ++digitsSinceSeparator;
    }

    std::reverse(grouped.begin(), grouped.end());
    return isNegative ? ("-" + grouped) : grouped;
}

/**
 * percentage da formato a un porcentaje con decimales fijos.
 * @param value - valor a formatear (p. ej. 18.0 -> "18.0000").
 * @param decimals - cantidad de decimales fijos a mostrar.
 */
inline std::string percentage(double value, int decimals = 4) {
    std::ostringstream oss;
    oss << std::fixed << std::setprecision(decimals) << value;
    return oss.str();
}

} // namespace number_format

namespace time_format {

/**
 * roundToMicroseconds redondea una duración al microsegundo más cercano.
 * @param nanoseconds - duración en nanosegundos a redondear.
 */
inline std::int64_t roundToMicroseconds(std::int64_t nanoseconds) {
    constexpr std::int64_t kMicro = 1000;
    const std::int64_t remainder = nanoseconds % kMicro;
    if (remainder * 2 >= kMicro) {
        return nanoseconds - remainder + kMicro;
    }
    return nanoseconds - remainder;
}

/**
 * formatDuration elige la unidad más legible (ns / µs / ms / s).
 * @param nanoseconds - duración en nanosegundos a formatear.
 */
inline std::string formatDuration(std::int64_t nanoseconds) {
    if (nanoseconds < 1000) {
        return std::to_string(nanoseconds) + "ns";
    }

    double value = 0.0;
    std::string unit;
    if (nanoseconds < 1000000) {
        value = static_cast<double>(nanoseconds) / 1000.0;

        // "µs" en UTF-8.
        unit = "\xC2\xB5s";
    } else if (nanoseconds < 1000000000) {
        value = static_cast<double>(nanoseconds) / 1000000.0;
        unit = "ms";
    } else {
        value = static_cast<double>(nanoseconds) / 1000000000.0;
        unit = "s";
    }

    std::ostringstream oss;
    oss << value << unit;
    return oss.str();
}

} // namespace time_format

namespace detail {

/**
 * trim quita espacios en blanco al inicio y al final de un texto.
 * @param text - texto a limpiar.
 */
inline std::string trim(const std::string& text) {
    const auto first = text.find_first_not_of(" \t\r\n");
    if (first == std::string::npos) {
        return "";
    }
    const auto last = text.find_last_not_of(" \t\r\n");
    return text.substr(first, last - first + 1);
}

/**
 * tryParseInt intenta convertir un texto a entero completo (sin sobrantes).
 * @param text - texto a convertir.
 * @param outValue - variable de salida con el valor convertido.
 * @returns true si la conversión fue exitosa.
 */
inline bool tryParseInt(const std::string& text, int& outValue) {
    if (text.empty()) {
        return false;
    }
    try {
        std::size_t consumed = 0;
        const int value = std::stoi(text, &consumed);
        if (consumed != text.size()) {
            return false;
        }
        outValue = value;
        return true;
    } catch (const std::exception&) {
        return false;
    }
}

} // namespace detail

// Mode: calificación más común, su frecuencia y porcentaje.
struct Mode {
    int grade = 0;
    std::int64_t count = 0;
    double percentage = 0.0;
};

// GradeVector: vector dinámico de calificaciones leídas.
class GradeVector {
public:
    /**
     * add agrega una calificación al vector.
     * @param grade - calificación a agregar; debe estar en [kMinGrade, kMaxGrade].
     * @throws std::out_of_range si grade está fuera de rango.
     */
    void add(int grade) {
        if (grade < kMinGrade || grade > kMaxGrade) {
            throw std::out_of_range(
                "La calificación " + std::to_string(grade) +
                " está fuera del rango permitido [" + std::to_string(kMinGrade) +
                ", " + std::to_string(kMaxGrade) + "].");
        }
        values_.push_back(grade);
    }

    std::int64_t totalCount() const noexcept { return static_cast<std::int64_t>(values_.size()); }
    bool isEmpty() const noexcept { return values_.empty(); }

    // sort ordena el vector de forma ascendente con std::sort.
    void sort() { std::sort(values_.begin(), values_.end()); }

    /**
     * merge agrega los valores de otro vector a este.
     * @param other - vector cuyos valores se agregarán a este.
     */
    void merge(const GradeVector& other) {
        values_.insert(values_.end(), other.values_.begin(), other.values_.end());
    }

    // mostCommonGrade cuenta corridas contiguas en el vector ya ordenado (llamar sort() antes); empata a favor del valor mayor.
    Mode mostCommonGrade() const {
        Mode result;

        const std::size_t total = values_.size();
        if (total == 0) {
            return result;
        }

        int runGrade = values_[0];
        std::int64_t runCount = 1;
        result.grade = runGrade;
        result.count = 1;

        for (std::size_t i = 1; i < total; ++i) {
            if (values_[i] == runGrade) {
                ++runCount;
            } else {
                runGrade = values_[i];
                runCount = 1;
            }

            if (runCount >= result.count) {
                result.count = runCount;
                result.grade = runGrade;
            }
        }

        result.percentage = (static_cast<double>(result.count) / static_cast<double>(total)) * 100.0;
        return result;
    }

private:
    std::vector<int> values_;
};

// FileReader lee calificaciones (una por línea) desde un archivo.
class FileReader {
public:
    struct ReadResult {
        bool success = false;
        std::string errorMessage;
        std::int64_t linesProcessed = 0;
        std::int64_t linesSkipped = 0;
    };

    /**
     * readInto lee un archivo línea por línea y agrega cada
     * calificación válida al vector.
     * @param filePath - ruta del archivo a leer.
     * @param vector - vector al que se agregará cada calificación válida.
     */
    ReadResult readInto(const std::string& filePath, GradeVector& vector) const {
        ReadResult result;

        std::ifstream file(filePath);
        if (!file.is_open()) {
            result.errorMessage = "No se pudo abrir el archivo '" + filePath + "'.";
            return result;
        }

        std::string rawLine;
        while (std::getline(file, rawLine)) {
            const std::string line = detail::trim(rawLine);
            if (line.empty()) {
                // Los renglones vacíos se ignoran silenciosamente.
                continue;
            }

            int grade = 0;
            if (!detail::tryParseInt(line, grade)) {
                ++result.linesSkipped;
                continue;
            }

            try {
                vector.add(grade);
                ++result.linesProcessed;
            } catch (const std::out_of_range&) {
                ++result.linesSkipped;
            }
        }

        result.success = true;
        return result;
    }
};

// Analyzer da formato a los resultados para el usuario.
class Analyzer {
public:
    /**
     * describeFile formatea el resultado de un archivo procesado.
     * @param fileName - nombre del archivo procesado.
     * @param totalCount - cantidad total de calificaciones válidas leídas.
     * @param mode - calificación más común calculada.
     * @param elapsedNanos - tiempo de leer + ordenar + contar (sin imprimir).
     */
    std::string describeFile(const std::string& fileName, std::int64_t totalCount, const Mode& mode,
                              std::int64_t elapsedNanos) const {
        std::ostringstream oss;
        oss << "Archivo: " << fileName << "\n"
            << "Cantidad de datos: " << number_format::withThousandsSeparator(totalCount) << "\n"
            << "Calificación más común: " << mode.grade << "\n"
            << "Frecuencia: " << number_format::withThousandsSeparator(mode.count) << "\n"
            << "Porcentaje: " << number_format::percentage(mode.percentage) << "%\n"
            << "Tiempo de procesamiento: " << time_format::formatDuration(time_format::roundToMicroseconds(elapsedNanos));
        return oss.str();
    }

    /**
     * describeGlobal formatea el resultado acumulado de todos los archivos.
     * @param totalCount - cantidad total de calificaciones válidas de todos los archivos.
     * @param mode - calificación más común global calculada.
     * @param modeElapsedNanos - tiempo de volver a ordenar + calcular la moda global.
     * @param totalElapsedNanos - suma de los tiempos de procesamiento por archivo.
     */
    std::string describeGlobal(std::int64_t totalCount, const Mode& mode, std::int64_t modeElapsedNanos,
                                std::int64_t totalElapsedNanos) const {
        std::ostringstream oss;
        oss << "RESULTADO ACUMULADO DE TODOS LOS ARCHIVOS\n"
            << "-----------------------------------------\n"
            << "Cantidad total de datos: " << number_format::withThousandsSeparator(totalCount) << "\n"
            << "Calificación más común global: " << mode.grade << "\n"
            << "Frecuencia global: " << number_format::withThousandsSeparator(mode.count) << "\n"
            << "Porcentaje global: " << number_format::percentage(mode.percentage) << "%\n"
            << "Tiempo del ordenamiento + cálculo de la moda global: "
            << time_format::formatDuration(time_format::roundToMicroseconds(modeElapsedNanos)) << "\n"
            << "Tiempo total acumulado (todos los archivos): "
            << time_format::formatDuration(time_format::roundToMicroseconds(totalElapsedNanos));
        return oss.str();
    }
};

} // namespace grades
