#pragma once
#include <juce_core/juce_core.h>

namespace ID
{
    // ---- A. VT root + top-level nodes ----
    const juce::Identifier CAKE         { "CAKE" };
    const juce::Identifier PROJECT      { "PROJECT" };
    const juce::Identifier BUILDS       { "BUILDS" };
    const juce::Identifier BUILD        { "BUILD" };
    const juce::Identifier GENERATOR    { "GENERATOR" };
    const juce::Identifier ASYNC        { "ASYNC" };
    const juce::Identifier CONFIG       { "CONFIG" };
    const juce::Identifier THEME        { "THEME" };

    // ---- B. PROJECT properties ----
    const juce::Identifier workingDirectory     { "workingDirectory" };
    const juce::Identifier hasCMakeLists        { "hasCMakeLists" };
    const juce::Identifier selectedGenerator    { "selectedGenerator" };
    const juce::Identifier configuration        { "configuration" };

    // ---- C. GENERATOR child properties ----
    const juce::Identifier name     { "name" };
    const juce::Identifier isIDE    { "isIDE" };

    // ---- D. BUILD child properties ----
    const juce::Identifier generator    { "generator" };
    const juce::Identifier path         { "path" };
    const juce::Identifier exists       { "exists" };
    const juce::Identifier isConfigured { "isConfigured" };

    // ---- E. ASYNC properties ----
    const juce::Identifier isActive     { "isActive" };
    const juce::Identifier isAborted    { "isAborted" };
    const juce::Identifier currentOp    { "currentOp" };

    // ---- F. CONFIG properties ----
    const juce::Identifier autoScanEnabled  { "autoScanEnabled" };
    const juce::Identifier autoScanInterval { "autoScanInterval" };
    const juce::Identifier theme            { "theme" };

    // ---- G. (reserved — CONSOLE subtree removed in Step 6; Console owns its buffer) ----

    // ---- H. THEME properties (TBD — color properties per theme XML) ----
    // Populated in Step 7 when ThemeLoader parses ~/.config/cake/themes/*.xml

} // namespace ID
