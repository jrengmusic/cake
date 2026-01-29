# Sprint 11 MACHINIST Kickoff - Cleanup Remaining Issues

**Date:** 2026-01-29
**From:** AUDITOR (Amp - Claude Sonnet 4)
**To:** MACHINIST
**Priority:** CRITICAL

---

## Context

Previous MACHINIST task (Sprint 9) was **incomplete**. Verification against codebase found remaining issues.

---

## Issues To Fix

### [FIX-001] CRITICAL: Remove "Ninja Multi-Config" from generators.go

**File:** `internal/utils/generators.go`

**Current (WRONG):**
```go
const (
    GeneratorXcode      = "Xcode"
    GeneratorNinja      = "Ninja"
    GeneratorNinjaMulti = "Ninja Multi-Config"  // ❌ REMOVE
    GeneratorVS2022     = "Visual Studio 17 2022"
    GeneratorVS2019     = "Visual Studio 16 2019"
)

func ValidGenerators() []string {
    return []string{
        GeneratorXcode,
        GeneratorNinja,
        GeneratorNinjaMulti,  // ❌ REMOVE
        GeneratorVS2022,
        GeneratorVS2019,
    }
}
```

**Expected (CORRECT):**
```go
const (
    GeneratorXcode  = "Xcode"
    GeneratorNinja  = "Ninja"
    GeneratorVS2026 = "Visual Studio 18 2026"  // CMake name
    GeneratorVS2022 = "Visual Studio 17 2022"  // CMake name
)

func ValidGenerators() []string {
    return []string{
        GeneratorXcode,
        GeneratorNinja,
        GeneratorVS2026,
        GeneratorVS2022,
    }
}
```

---

### [FIX-002] CRITICAL: Update Visual Studio Versions in project.go

**File:** `internal/state/project.go` (lines 132-139)

**Current (WRONG):**
```go
ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
    Name:  "Visual Studio 17 2022",  // ❌ CMake-style
    IsIDE: true,
})
ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
    Name:  "Visual Studio 16 2019",  // ❌ CMake-style
    IsIDE: true,
})
```

**Expected (CORRECT):**
```go
ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
    Name:  "Visual Studio 18 2026",  // ✅ CMake name, dir mapped to VS2026
    IsIDE: true,
})
ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
    Name:  "Visual Studio 17 2022",  // ✅ CMake name, dir mapped to VS2022
    IsIDE: true,
})
```

---

### [FIX-003] HIGH: Update generator.go switch cases

**File:** `internal/utils/generator.go`

Update `BuildBuildCommand()` switch cases:
- Keep `"Visual Studio 17 2022"` (CMake name)
- Change `"Visual Studio 16 2019"` → `"Visual Studio 18 2026"` (new version)

Update `validGenerators` slice in `BuildCMakeCommand()`:
- Remove `"Ninja Multi-Config"`
- Use CMake names: `"Visual Studio 18 2026"`, `"Visual Studio 17 2022"`

---

### [FIX-004] MEDIUM: Update IsGeneratorIDE() in generators.go

**File:** `internal/utils/generators.go`

Update switch case to use new constants:
```go
func IsGeneratorIDE(generator string) bool {
    switch generator {
    case GeneratorXcode, GeneratorVS2026, GeneratorVS2022:
        return true
    default:
        return false
    }
}
```

---

## Verification Checklist

After fixes, run:
```bash
grep -r "Ninja Multi-Config" internal/
grep -r "Visual Studio 17" internal/
grep -r "Visual Studio 16" internal/
```

All should return **empty**.

---

## Generator Names (SSOT)

| Constant | CMake -G Flag | Directory | Display |
|----------|---------------|-----------|---------|
| GeneratorXcode | Xcode | Xcode | Xcode |
| GeneratorNinja | Ninja | Ninja | Ninja |
| GeneratorVS2026 | Visual Studio 18 2026 | VS2026 | VS2026 |
| GeneratorVS2022 | Visual Studio 17 2022 | VS2022 | VS2022 |

**Approach:** Generator name = CMake name. Mapping function converts to directory name.

---

### [FIX-005] HIGH: Add GetDirectoryName() Mapping Function

**File:** `internal/utils/generators.go` (add new function)

```go
// GetDirectoryName returns the shortened directory name for a generator
// Used for Builds/<dir>/ path construction
func GetDirectoryName(generator string) string {
    switch generator {
    case "Visual Studio 18 2026":
        return "VS2026"
    case "Visual Studio 17 2022":
        return "VS2022"
    case "Visual Studio 16 2019":
        return "VS2019"
    default:
        return generator // Xcode, Ninja unchanged
    }
}
```

---

### [FIX-006] HIGH: Update All Build Path Constructions

**Files to update:**
- `internal/ops/setup.go:26`
- `internal/ops/build.go:17`
- `internal/ops/clean.go:15`
- `internal/ops/open.go:18`
- `internal/state/project.go:227,233,253`
- `internal/app/op_regenerate.go:35`

**Change from:**
```go
buildDir := filepath.Join(workingDir, "Builds", generator)
```

**To:**
```go
buildDir := filepath.Join(workingDir, "Builds", utils.GetDirectoryName(generator))
```

**Note:** This is REF-001 from audit - consider passing buildDir as parameter instead of recalculating everywhere.

---

## Build Verification

```bash
./build.sh
```

Must complete without errors.

---

**End of Kickoff**
