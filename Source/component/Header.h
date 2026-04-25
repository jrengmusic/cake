#pragma once
#include <JuceHeader.h>

// ============================================================================
// Header
// ============================================================================
//
// Displays the working directory path.
// Stateless View — holds only what is needed for the current paint.
// Working directory is pushed by MainComponent via setWorkingDirectory().
//
// BLESSED-S (stateless View), BLESSED-E (no listener — data pushed by Controller).

class Header : public jam::tui::Component
{
public:
    Header() = default;

    // Sets the path displayed in the header.
    // Called by MainComponent when projectState.workingDirectory changes.
    void setWorkingDirectory (const juce::String& path);

    // -------------------------------------------------------------------------
    // Component overrides
    // -------------------------------------------------------------------------
    void paint (jam::tui::Graphics& g) override;

private:
    juce::String workingDirectory;

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (Header)
};
