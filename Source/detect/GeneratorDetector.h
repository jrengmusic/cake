#pragma once
#include <JuceHeader.h>
#include "../state/Axis.h"

// ============================================================================
// GeneratorDetector
// ============================================================================
//
// Detects which CMake generators are available on the current system by
// testing tool presence in PATH.  Returns an ordered list of detected
// generators per SPEC §Available Generators Detection.
//
// Detection is synchronous and cheap — only called once at startup and on
// auto-scan ticks.  Each check spawns a short-lived child process (`which`)
// and checks the exit code.  No disk scanning.
//
// Platform:
//   macOS  — Xcode (xcodebuild), Ninja (ninja)
//   Linux  — Ninja (ninja)
//   Windows — Visual Studio (vswhere.exe), Ninja (ninja) — STUB, post-MVP
//
// SPEC §Generator Types; BLESSED-E (no magic strings), BLESSED-S (SSOT for
// detection logic), BLESSED-S2 (pure function, no stored state).

struct GeneratorDescriptor
{
    Generator    generator;
    bool         isIDE;
};

class GeneratorDetector
{
public:

    // Detects all available generators on the current platform.
    // Returns them in priority order (IDE generators first, per Go reference).
    // Called on the message thread — synchronous but fast (which/where exits quickly).
    // MESSAGE THREAD.
    static juce::Array<GeneratorDescriptor> detectAvailableGenerators() noexcept;

private:

    // Spawns `which <toolName>` (macOS/Linux) and returns true if the tool
    // was found in PATH (exit code 0).
    // Blocks until the child exits; which is near-instant.
    static bool isToolInPath (const juce::String& toolName) noexcept;

    JUCE_DECLARE_NON_COPYABLE (GeneratorDetector)
};
