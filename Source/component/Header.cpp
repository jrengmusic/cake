#include <JuceHeader.h>
#include "Header.h"

// ============================================================================
// Public API
// ============================================================================

void Header::setWorkingDirectory (const juce::String& path)
{
    workingDirectory = path;
    repaint();
}

// ============================================================================
// Component overrides
// ============================================================================

void Header::paint (jam::tui::Graphics& g)
{
    const jam::tui::Rectangle bounds { getCellBounds() };

    if (not bounds.isEmpty())
    {
        g.setColour (juce::Colours::grey);
        g.drawText ("Project directory:",
                    bounds.withHeight (jam::literals::Cell { 1 }),
                    juce::Justification::centredLeft);

        g.setColour (juce::Colours::white);
        g.drawText (workingDirectory,
                    bounds.withTrimmedTop (jam::literals::Cell { 1 }).withHeight (jam::literals::Cell { 1 }),
                    juce::Justification::centredLeft);
    }
}
