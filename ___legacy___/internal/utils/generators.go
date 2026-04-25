package utils

// Generator name constants - SSOT for all generator references
// These are the CMake generator names passed to -G flag
const (
	GeneratorXcode  = "Xcode"
	GeneratorNinja  = "Ninja"
	GeneratorVS2026 = "Visual Studio 18 2026"
	GeneratorVS2022 = "Visual Studio 17 2022"
)

// ValidGenerators returns all supported generator names
func ValidGenerators() []string {
	return []string{
		GeneratorXcode,
		GeneratorNinja,
		GeneratorVS2026,
		GeneratorVS2022,
	}
}

// GetDirectoryName returns shortened directory name for a generator
// Used for Builds/<dir>/ path construction
func GetDirectoryName(generator string) string {
	switch generator {
	case GeneratorVS2026:
		return "VS2026"
	case GeneratorVS2022:
		return "VS2022"
	default:
		return generator // Xcode, Ninja unchanged
	}
}

// GetGeneratorNameFromDirectory returns CMake generator name from directory name
// Used when scanning existing build directories
func GetGeneratorNameFromDirectory(dirName string) string {
	switch dirName {
	case "VS2026":
		return GeneratorVS2026
	case "VS2022":
		return GeneratorVS2022
	default:
		return dirName // Xcode, Ninja unchanged
	}
}

// IsGeneratorIDE returns true if generator creates an IDE project
func IsGeneratorIDE(generator string) bool {
	switch generator {
	case GeneratorXcode, GeneratorVS2026, GeneratorVS2022:
		return true
	default:
		return false
	}
}
