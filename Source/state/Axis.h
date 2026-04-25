#pragma once
#include <juce_core/juce_core.h>

// SPEC §State Model — enums for CAKE state axes.
// Mode is computed from VT state, not stored.
// Generator, Configuration, OpType are stored as strings in VT via toString().

    enum class Mode          { menu, preferences, console, invalidProject };
    enum class Generator     { xcode, ninja, vs2026, vs2022 };
    enum class Configuration { debug, release };
    enum class OpType        { none, build, generate, clean, cleanAll, regenerate };

    // String <-> enum bridging (VT stores strings per SPEC §ValueTree Schema)
    juce::String toString (Mode)          noexcept;
    juce::String toString (Generator)     noexcept;
    juce::String toString (Configuration) noexcept;
    juce::String toString (OpType)        noexcept;

    Generator     parseGenerator     (const juce::String&) noexcept;
    Configuration parseConfiguration (const juce::String&) noexcept;
    OpType        parseOpType        (const juce::String&) noexcept;
