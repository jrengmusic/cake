#pragma once
#include <JuceHeader.h>

// ============================================================================
// Footer
// ============================================================================
//
// Displays context hints at the bottom of the main layout.
// Stateless View — holds only the hint string it is told to display.
// Caller sets hint via setHint() when mode changes.
//
// BLESSED-S (stateless), BLESSED-E (no listener — hint pushed by Controller).

class Footer : public jam::tui::Component
{
public:
    Footer() = default;

    // Sets the hint text displayed in the footer row.
    // Called by MainComponent on mode transitions.
    void setHint (const juce::String& hint);

    // -------------------------------------------------------------------------
    // Component overrides
    // -------------------------------------------------------------------------
    void paint (jam::tui::Graphics& g) override;

private:
    juce::String hintText;

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (Footer)
};
