#pragma once
#include <JuceHeader.h>
#include <jam_subprocess/jam_subprocess.h>
#include "../state/State.h"
#include "../state/Axis.h"

// Forward declaration to avoid pulling in the full jam::tui header here.
// CmakeRunner operates entirely via appendLine() API — it never renders.
namespace jam::tui { class Console; }

// ============================================================================
// CmakeRunner
// ============================================================================
//
// Subprocess orchestrator for cmake operations.
//
// Owns a single jam::Subprocess (BLESSED-B — one owner, deterministic lifecycle).
// Each public method sets ASYNC atoms on State, spawns a subprocess, streams
// stdout/stderr lines to Console via appendLine() + callAsync, then clears
// ASYNC state on completion.
//
// Threading:
//   generate(), build(), clean(), open() — called on the message thread.
//   jam::Subprocess::Handler::Chunk fires on the Worker thread.
//   Lines are marshalled to Console via appendLine() which posts callAsync.
//   ASYNC atom writes are lock-free from any thread.
//
// Encapsulation (BLESSED-E):
//   CmakeRunner knows State (reads project settings, writes ASYNC atoms)
//   and Console (appends lines). It knows nothing about views or layout.
//
// SPEC §Feature: Generate, Build, Clean, Open IDE/Editor.
// PLAN §Step 6.

class CmakeRunner
{
public:

    // -------------------------------------------------------------------------
    // Construction
    // -------------------------------------------------------------------------

    // Takes non-owning references to State and Console.
    // Both must outlive CmakeRunner — owned by MainComponent (BLESSED-B).
    explicit CmakeRunner (State& state, jam::tui::Console& console);

    ~CmakeRunner();

    // -------------------------------------------------------------------------
    // Operations — each assembles the correct cmake command and streams output.
    // MESSAGE THREAD.
    // -------------------------------------------------------------------------

    // cmake -S . -B Builds/<Gen> -G <GenName>
    // SPEC §Feature: Generate/Regenerate Operation.
    void generate();

    // cmake --build Builds/<Gen> --config <Cfg>
    // SPEC §Feature: Build Operation.
    void build();

    // Deletes Builds/<Gen>/ recursively (no subprocess — file deletion).
    // SPEC §Feature: Clean Operation.
    void clean();

    // Opens the IDE/editor for the selected generator.
    // Xcode: open *.xcodeproj  |  Ninja: nvim Builds/Ninja/
    // SPEC §Feature: Open IDE/Editor.
    void open();

    // -------------------------------------------------------------------------
    // Abort — kills running subprocess, sets isAborted on State.
    // MESSAGE THREAD.
    // -------------------------------------------------------------------------

    void abort();

private:

    // =========================================================================
    // Generator constants — cmake -G argument strings (BLESSED-E: named, not magic)
    // =========================================================================

    static constexpr const char* CMAKE_G_XCODE  { "Xcode" };
    static constexpr const char* CMAKE_G_NINJA  { "Ninja" };
    // VS generators deferred — macOS-first per PLAN §Risks.

    // =========================================================================
    // Build directory name mapping per SPEC §Generator Types
    // =========================================================================

    static constexpr const char* BUILD_DIR_XCODE { "Xcode" };
    static constexpr const char* BUILD_DIR_NINJA { "Ninja" };

    // Root "Builds/" directory per SPEC §Generator Types and PLAN §Step 6.
    static constexpr const char* BUILDS_ROOT { "Builds" };

    // =========================================================================
    // Helpers
    // =========================================================================

    // Returns the build directory relative to the working directory
    // e.g. "Builds/Xcode" for Generator::xcode.
    // SPEC §Build Path Convention.
    juce::String buildDirName (Generator gen) const;

    // Returns the cmake -G argument for the given generator.
    juce::String cmakeGeneratorArg (Generator gen) const;

    // Returns absolute path to the build directory.
    juce::File buildDir (Generator gen) const;

    // Appends info lines (header) to console before launching subprocess.
    void appendOpHeader (const juce::String& label,
                         const juce::String& command) const;

    // Called on completion of a subprocess operation (Worker thread).
    // Sets isActive=false, isAborted reflects abort state.
    void onOperationComplete (int exitCode);

    // =========================================================================
    // State references (non-owning, BLESSED-B — outlives CmakeRunner)
    // =========================================================================

    State&            stateRef;
    jam::tui::Console& consoleRef;

    // =========================================================================
    // Subprocess — owned by CmakeRunner (BLESSED-B)
    // =========================================================================

    jam::Subprocess subprocess;

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (CmakeRunner)
};
