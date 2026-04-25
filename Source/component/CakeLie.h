#pragma once
#include <JuceHeader.h>

// ============================================================================
// CakeLie
// ============================================================================
//
// Full-screen braille banner for invalidProject mode.
// Renders cake-lie.svg as braille art, centered in the available bounds.
//
// SVG is loaded once from BinaryData on construction.  BrailleGrid is
// rasterized on every resize, rendered from cache on paint.
//
// Stateless View — holds only cached render state (braille grid).
// No VT listener — static display, no data dependency.
//
// SPEC §Edge Case 1: No CMakeLists.txt — Mode::invalidProject.
// BLESSED-S (stateless View), BLESSED-B (BrailleGrid owned here, lifecycle clear).

class CakeLie : public jam::tui::Component
{
public:
    CakeLie();
    ~CakeLie() override = default;

    // -------------------------------------------------------------------------
    // Component overrides
    // -------------------------------------------------------------------------
    void paint   (jam::tui::Graphics& g) override;
    void resized ()                      override;

private:
    // -------------------------------------------------------------------------
    // Helpers
    // -------------------------------------------------------------------------
    void rasterizeToGrid() noexcept;

    // -------------------------------------------------------------------------
    // SVG source — loaded once from BinaryData
    // -------------------------------------------------------------------------
    std::unique_ptr<juce::XmlElement> svgDocument;

    // -------------------------------------------------------------------------
    // Cached braille grid — rebuilt on resize
    // -------------------------------------------------------------------------
    jam::braille::BrailleGrid brailleGrid;

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (CakeLie)
};
