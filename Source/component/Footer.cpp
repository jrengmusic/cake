#include <JuceHeader.h>
#include "Footer.h"

// ============================================================================
// File-scope constants
// ============================================================================

static constexpr juce::uint32 COLOR_FOOTER_TEXT { 0xff888888 };

// ============================================================================
// Content
// ============================================================================

void Footer::setHint (const juce::String& hint)
{
    hintText = hint;
    repaint();
}

// ============================================================================
// Component overrides
// ============================================================================

void Footer::paint (jam::tui::Graphics& g)
{
    const jam::tui::Rectangle bounds { getCellBounds() };

    if (not bounds.isEmpty())
    {
        g.setColour (juce::Colour { COLOR_FOOTER_TEXT });
        g.drawText (hintText, bounds, juce::Justification::centredLeft);
    }
}
