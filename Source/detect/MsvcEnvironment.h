#pragma once
#include <JuceHeader.h>

// ============================================================================
// MsvcEnvironment
// ============================================================================
//
// Captures the MSVC environment (PATH, INCLUDE, LIB) needed for cmake to
// invoke cl.exe and vswhere.exe on Windows.
//
// macOS stub — no-op.  Windows implementation deferred to post-MVP sprint.
// Interface is minimal: just enough for the detect pipeline to compile on
// both platforms.
//
// PLAN §Step 3; BLESSED-B (bound to compilation unit — no runtime cost on macOS).

class MsvcEnvironment
{
public:

    // Returns true if the MSVC environment was successfully captured.
    // On macOS: always returns false (no MSVC available).
    // On Windows (post-MVP): runs vswhere.exe, captures vcvarsall output.
    bool captureEnvironment() noexcept;

    // Returns true if MSVC environment has been captured and is valid.
    bool isValid() const noexcept;

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (MsvcEnvironment)
};
