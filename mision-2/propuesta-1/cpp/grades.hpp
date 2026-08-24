#pragma once

// Propuesta #1: arreglo de contadores indexado por calificación (header-only; main.cpp orquesta).
#include <algorithm>
#include <array>
#include <cstdint>
#include <fstream>
#include <iomanip>
#include <sstream>
#include <stdexcept>
#include <string>

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

// FrequencyTable: contador por calificación.
class FrequencyTable {
public:
    static constexpr std::size_t kTableSize = static_cast<std::size_t>(kMaxGrade) + 1;

    /**
     * increment suma 1 al contador de una calificación.
     * @param grade - calificación a incrementar; debe estar en [kMinGrade, kMaxGrade].
     * @throws std::out_of_range si grade está fuera de rango.
     */
    void increment(int grade) {
        if (grade < kMinGrade || grade > kMaxGrade) {
            throw std::out_of_range(
                "La calificación " + std::to_string(grade) +
                " está fuera del rango permitido [" + std::to_string(kMinGrade) +
                ", " + std::to_string(kMaxGrade) + "].");
        }
        ++counters_[static_cast<std::size_t>(grade)];
        ++total_;
    }

    std::int64_t totalCount() const noexcept { return total_; }
    bool isEmpty() const noexcept { return total_ == 0; }

    // mostCommonGrade recorre de 100 a 1; en empate gana la calificación mayor.
    Mode mostCommonGrade() const {
        Mode result;
        if (total_ == 0) {
            return result;
        }

        for (int grade = kMaxGrade; grade >= kMinGrade; --grade) {
            const std::int64_t count = counters_[static_cast<std::size_t>(grade)];
            if (count > result.count) {
                result.grade = grade;
                result.count = count;
            }
        }

        result.percentage =
            (static_cast<double>(result.count) / static_cast<double>(total_)) * 100.0;
        return result;
    }

    /**
     * mergeFrom acumula los contadores de otra tabla en esta.
     * @param other - tabla cuyos contadores se sumarán a esta.
     */
    void mergeFrom(const FrequencyTable& other) {
        for (std::size_t i = static_cast<std::size_t>(kMinGrade); i <= static_cast<std::size_t>(kMaxGrade); ++i) {
            counters_[i] += other.counters_[i];
        }
        total_ += other.total_;
    }

private:
    std::array<std::int64_t, kTableSize> counters_{};
    std::int64_t total_ = 0;
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
     * readInto lee un archivo línea por línea e incrementa la tabla
     * por cada calificación válida.
     * @param filePath - ruta del archivo a leer.
     * @param table - tabla que se incrementará con cada calificación válida.
     */
    ReadResult readInto(const std::string& filePath, FrequencyTable& table) const {
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
                table.increment(grade);
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
     * @param elapsedNanos - tiempo de leer + contar + calcular la moda (sin imprimir).
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
     * @param modeElapsedNanos - tiempo de calcular la moda global.
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
            << "Tiempo del cálculo de la moda global: " << time_format::formatDuration(modeElapsedNanos) << "\n"
            << "Tiempo total acumulado (todos los archivos): "
            << time_format::formatDuration(time_format::roundToMicroseconds(totalElapsedNanos));
        return oss.str();
    }
};

} // namespace grades
