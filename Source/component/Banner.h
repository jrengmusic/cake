#pragma once
#include <JuceHeader.h>

// ============================================================================
// Banner
// ============================================================================
//
// Renders cake-logo.svg as braille art in the 35% right pane.
// SVG is loaded once from BinaryData on construction and rasterized to a
// BrailleGrid whenever the component is resized.  paint() draws the cached grid.
//
// Stateless View — holds only cached render state (braille grid).
// Re-rasterizes on resize; renders from cache on paint.
//
// BLESSED-S (stateless View), BLESSED-B (BrailleGrid owned here, lifecycle clear).

class Banner : public jam::tui::Component
{
public:
    Banner();
    ~Banner() override = default;

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

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (Banner)
};
