package state

import (
	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/utils"
	"path/filepath"
	"testing"
)

// --- helpers ---

func makeState(projects []Generator, selected string) *ProjectState {
	return &ProjectState{
		WorkingDirectory:  "/tmp/testproject",
		HasCMakeLists:     true,
		AvailableProjects: projects,
		SelectedProject:   selected,
		Builds:            make(map[string]BuildInfo),
		Configuration:     internal.ConfigDebug,
	}
}

func gens(names ...string) []Generator {
	result := make([]Generator, len(names))
	for i, n := range names {
		result[i] = Generator{Name: n}
	}
	return result
}

// --- CycleToNextProject ---

func TestCycleToNextProject(t *testing.T) {
	tests := []struct {
		name      string
		projects  []Generator
		selected  string
		wantAfter string
	}{
		{
			name:      "zero projects no-ops",
			projects:  []Generator{},
			selected:  "",
			wantAfter: "",
		},
		{
			name:      "one project stays on same",
			projects:  gens("Ninja"),
			selected:  "Ninja",
			wantAfter: "Ninja",
		},
		{
			name:      "two projects advance",
			projects:  gens("Ninja", "Xcode"),
			selected:  "Ninja",
			wantAfter: "Xcode",
		},
		{
			name:      "two projects wrap from last",
			projects:  gens("Ninja", "Xcode"),
			selected:  "Xcode",
			wantAfter: "Ninja",
		},
		{
			name:      "three projects advance mid",
			projects:  gens("Ninja", "Xcode", utils.GeneratorVS2022),
			selected:  "Xcode",
			wantAfter: utils.GeneratorVS2022,
		},
		{
			name:      "three projects wrap from last",
			projects:  gens("Ninja", "Xcode", utils.GeneratorVS2022),
			selected:  utils.GeneratorVS2022,
			wantAfter: "Ninja",
		},
		{
			name:      "unknown selected falls back to first",
			projects:  gens("Ninja", "Xcode"),
			selected:  "Unknown",
			wantAfter: "Ninja",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := makeState(tc.projects, tc.selected)
			ps.CycleToNextProject()
			if ps.SelectedProject != tc.wantAfter {
				t.Errorf("got %q, want %q", ps.SelectedProject, tc.wantAfter)
			}
		})
	}
}

// --- CycleToPrevProject ---

func TestCycleToPrevProject(t *testing.T) {
	tests := []struct {
		name      string
		projects  []Generator
		selected  string
		wantAfter string
	}{
		{
			name:      "zero projects no-ops",
			projects:  []Generator{},
			selected:  "",
			wantAfter: "",
		},
		{
			name:      "one project stays on same",
			projects:  gens("Ninja"),
			selected:  "Ninja",
			wantAfter: "Ninja",
		},
		{
			name:      "two projects go back",
			projects:  gens("Ninja", "Xcode"),
			selected:  "Xcode",
			wantAfter: "Ninja",
		},
		{
			name:      "two projects wrap from first",
			projects:  gens("Ninja", "Xcode"),
			selected:  "Ninja",
			wantAfter: "Xcode",
		},
		{
			name:      "three projects go back mid",
			projects:  gens("Ninja", "Xcode", utils.GeneratorVS2022),
			selected:  "Xcode",
			wantAfter: "Ninja",
		},
		{
			name:      "three projects wrap from first",
			projects:  gens("Ninja", "Xcode", utils.GeneratorVS2022),
			selected:  "Ninja",
			wantAfter: utils.GeneratorVS2022,
		},
		{
			name:      "unknown selected falls back to first",
			projects:  gens("Ninja", "Xcode"),
			selected:  "Unknown",
			wantAfter: "Ninja",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := makeState(tc.projects, tc.selected)
			ps.CycleToPrevProject()
			if ps.SelectedProject != tc.wantAfter {
				t.Errorf("got %q, want %q", ps.SelectedProject, tc.wantAfter)
			}
		})
	}
}

// --- CycleConfiguration ---

func TestCycleConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		before string
		after  string
	}{
		{"Debug to Release", internal.ConfigDebug, internal.ConfigRelease},
		{"Release to Debug", internal.ConfigRelease, internal.ConfigDebug},
		{"unknown stays to Debug", "Unknown", internal.ConfigDebug},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := makeState(nil, "")
			ps.Configuration = tc.before
			ps.CycleConfiguration()
			if ps.Configuration != tc.after {
				t.Errorf("got %q, want %q", ps.Configuration, tc.after)
			}
		})
	}
}

// --- SetConfiguration ---

func TestSetConfiguration(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue string
	}{
		{"Debug accepted", internal.ConfigDebug, internal.ConfigDebug},
		{"Release accepted", internal.ConfigRelease, internal.ConfigRelease},
		{"invalid ignored", "MinSizeRel", internal.ConfigDebug},
		{"empty ignored", "", internal.ConfigDebug},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := makeState(nil, "")
			ps.Configuration = internal.ConfigDebug
			ps.SetConfiguration(tc.input)
			if ps.Configuration != tc.wantValue {
				t.Errorf("got %q, want %q", ps.Configuration, tc.wantValue)
			}
		})
	}
}

// --- SetSelectedProject ---

func TestSetSelectedProject(t *testing.T) {
	available := gens("Ninja", "Xcode")

	tests := []struct {
		name         string
		input        string
		wantSelected string
	}{
		{"valid name Ninja", "Ninja", "Ninja"},
		{"valid name Xcode", "Xcode", "Xcode"},
		{"invalid name ignored", "Unknown", ""},
		{"empty name ignored", "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := makeState(available, "")
			ps.SetSelectedProject(tc.input)
			if ps.SelectedProject != tc.wantSelected {
				t.Errorf("got %q, want %q", ps.SelectedProject, tc.wantSelected)
			}
		})
	}
}

// --- CanGenerate ---

func TestCanGenerate(t *testing.T) {
	tests := []struct {
		name          string
		selected      string
		hasCMake      bool
		wantCanGen    bool
	}{
		{"selected + cmake = true", "Ninja", true, true},
		{"no selected = false", "", true, false},
		{"no cmake = false", "Ninja", false, false},
		{"neither = false", "", false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := makeState(gens("Ninja"), tc.selected)
			ps.HasCMakeLists = tc.hasCMake
			if got := ps.CanGenerate(); got != tc.wantCanGen {
				t.Errorf("got %v, want %v", got, tc.wantCanGen)
			}
		})
	}
}

// --- CanBuild ---

func TestCanBuild(t *testing.T) {
	tests := []struct {
		name         string
		selected     string
		buildInfo    BuildInfo
		wantCanBuild bool
	}{
		{
			name:     "exists and configured = true",
			selected: "Ninja",
			buildInfo: BuildInfo{
				Generator:    "Ninja",
				Exists:       true,
				IsConfigured: true,
			},
			wantCanBuild: true,
		},
		{
			name:     "exists but not configured = false",
			selected: "Ninja",
			buildInfo: BuildInfo{
				Generator:    "Ninja",
				Exists:       true,
				IsConfigured: false,
			},
			wantCanBuild: false,
		},
		{
			name:         "no build info in map = false",
			selected:     "Xcode",
			buildInfo:    BuildInfo{},
			wantCanBuild: false,
		},
		{
			name:         "no selected project = false",
			selected:     "",
			buildInfo:    BuildInfo{},
			wantCanBuild: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := makeState(gens("Ninja", "Xcode"), tc.selected)
			if tc.buildInfo.Generator != "" {
				ps.Builds[tc.selected] = tc.buildInfo
			}
			if got := ps.CanBuild(); got != tc.wantCanBuild {
				t.Errorf("got %v, want %v", got, tc.wantCanBuild)
			}
		})
	}
}

// --- GetProjectLabel ---

func TestGetProjectLabel(t *testing.T) {
	tests := []struct {
		name      string
		selected  string
		wantLabel string
	}{
		{"Ninja returned as-is", "Ninja", "Ninja"},
		{"Xcode returned as-is", "Xcode", "Xcode"},
		{"VS2026 abbreviated", utils.GeneratorVS2026, "VS 2026"},
		{"VS2022 abbreviated", utils.GeneratorVS2022, "VS 2022"},
		{"empty returns empty", "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := makeState(nil, tc.selected)
			if got := ps.GetProjectLabel(); got != tc.wantLabel {
				t.Errorf("got %q, want %q", got, tc.wantLabel)
			}
		})
	}
}

// --- GetBuildDirectory ---

func TestGetBuildDirectory(t *testing.T) {
	tests := []struct {
		name      string
		generator string
		wantDir   string
	}{
		{
			name:      "Ninja",
			generator: "Ninja",
			wantDir:   filepath.Join("/tmp/testproject", internal.BuildsDirName, "Ninja"),
		},
		{
			name:      "Xcode",
			generator: "Xcode",
			wantDir:   filepath.Join("/tmp/testproject", internal.BuildsDirName, "Xcode"),
		},
		{
			name:      "VS2026 uses short dir",
			generator: utils.GeneratorVS2026,
			wantDir:   filepath.Join("/tmp/testproject", internal.BuildsDirName, "VS2026"),
		},
		{
			name:      "VS2022 uses short dir",
			generator: utils.GeneratorVS2022,
			wantDir:   filepath.Join("/tmp/testproject", internal.BuildsDirName, "VS2022"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ps := makeState(nil, "")
			if got := ps.GetBuildDirectory(tc.generator); got != tc.wantDir {
				t.Errorf("got %q, want %q", got, tc.wantDir)
			}
		})
	}
}

// --- GetSelectedBuildInfo ---

func TestGetSelectedBuildInfo(t *testing.T) {
	t.Run("returns BuildInfo from map when present", func(t *testing.T) {
		ps := makeState(gens("Ninja"), "Ninja")
		ps.Builds["Ninja"] = BuildInfo{
			Generator:    "Ninja",
			Path:         "/tmp/testproject/Builds/Ninja",
			Exists:       true,
			IsConfigured: true,
		}
		info := ps.GetSelectedBuildInfo()
		if !info.Exists || !info.IsConfigured || info.Generator != "Ninja" {
			t.Errorf("unexpected BuildInfo: %+v", info)
		}
	})

	t.Run("returns zero BuildInfo with Exists=false when not in map", func(t *testing.T) {
		ps := makeState(gens("Xcode"), "Xcode")
		info := ps.GetSelectedBuildInfo()
		if info.Exists {
			t.Errorf("expected Exists=false, got true")
		}
		if info.Generator != "Xcode" {
			t.Errorf("expected generator %q, got %q", "Xcode", info.Generator)
		}
	})

	t.Run("empty selected returns zero BuildInfo", func(t *testing.T) {
		ps := makeState(nil, "")
		info := ps.GetSelectedBuildInfo()
		if info.Exists {
			t.Errorf("expected Exists=false, got true")
		}
	})
}

// --- parseProjectCallName ---

func TestParseProjectCallName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain name", "MyProject)", "MyProject"},
		{"plain name with space", "MyProject VERSION 1.0)", "MyProject"},
		{"quoted name", `"MyProject" VERSION 1.0)`, "MyProject"},
		{"quoted name only", `"MyProject")`, "MyProject"},
		{"variable reference rejected", "${PROJECT_NAME})", ""},
		{"variable in quoted rejected", `"${PROJECT_NAME}")`, ""},
		{"empty string", "", ""},
		{"name with tab separator", "MyProject\tVERSION)", "MyProject"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseProjectCallName(tc.input); got != tc.want {
				t.Errorf("parseProjectCallName(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
