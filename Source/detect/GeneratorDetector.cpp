#include "GeneratorDetector.h"

// ============================================================================
// Tool name constants — BLESSED-E: no magic strings
// ============================================================================

static const juce::String TOOL_XCODEBUILD { "xcodebuild" };
static const juce::String TOOL_NINJA      { "ninja" };
static const juce::String TOOL_VSWHERE    { "vswhere.exe" };

static const juce::String CMD_WHICH       { "which" };
static const juce::String CMD_WHERE       { "where" };    // Windows; unused on macOS

// Timeout for each `which` child process — 2 seconds is more than enough.
static constexpr int TOOL_CHECK_TIMEOUT_MS { 2000 };

// ============================================================================
// isToolInPath
// ============================================================================

bool GeneratorDetector::isToolInPath (const juce::String& toolName) noexcept
{
    juce::ChildProcess probe;

    const juce::StringArray args { CMD_WHICH, toolName };
    const bool started { probe.start (args, juce::ChildProcess::wantStdOut) };

    bool found { false };

    if (started)
    {
        probe.waitForProcessToFinish (TOOL_CHECK_TIMEOUT_MS);
        found = (probe.getExitCode() == 0);
    }

    return found;
}

// ============================================================================
// detectAvailableGenerators
// ============================================================================

juce::Array<GeneratorDescriptor> GeneratorDetector::detectAvailableGenerators() noexcept
{
    juce::Array<GeneratorDescriptor> detected;

#if JUCE_MAC

    // Xcode — macOS only, IDE generator.  Priority 1.
    if (isToolInPath (TOOL_XCODEBUILD))
        detected.add ({ Generator::xcode, true });

    // Ninja — cross-platform, CLI generator.  Priority 2 on macOS.
    if (isToolInPath (TOOL_NINJA))
        detected.add ({ Generator::ninja, false });

#elif JUCE_LINUX

    // Ninja — only detected generator on Linux.
    if (isToolInPath (TOOL_NINJA))
        detected.add ({ Generator::ninja, false });

#elif JUCE_WINDOWS

    // Windows detection is post-MVP — MsvcEnvironment stub handles VS.
    // Ninja on Windows also requires VS env capture; deferred to post-MVP sprint.
    // Both are intentionally empty stubs here.
    juce::ignoreUnused (TOOL_VSWHERE);
    juce::ignoreUnused (CMD_WHERE);

#endif

    return detected;
}
