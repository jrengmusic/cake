#pragma once
#include <JuceHeader.h>
#include "../Identifier.h"

// ============================================================================
// Preferences
// ============================================================================
//
// Stateless View for the preferences screen.  Reads and writes CONFIG
// properties in State VT.  Owns a jam::tui::Menu for the 5-row preferences
// layout.  Rebuilds rows from VT on every CONFIG change via ValueTree::Listener.
//
// Rows (SPEC §Feature: Preferences Screen):
//   0  Auto-update       ON / OFF       (writes CONFIG.autoScanEnabled)
//   1  Update Interval   N min          (writes CONFIG.autoScanInterval)
//   2  Separator         ───            (not selectable)
//   3  Theme             name           (writes CONFIG.theme)
//   4  Back to Menu                     (fires onBackToMenu callback)
//
// +/= and -/_ adjust interval on the interval row.
// Shift+= increases by 10 min; Shift+- decreases by 10 min (SPEC §Keyboard).
//
// BLESSED-S (stateless View), BLESSED-E (listener pattern, callback for action).

class Preferences : public jam::tui::Component,
                    public juce::ValueTree::Listener
{
public:
    // Constructs Preferences with a reference to the CONFIG subtree.
    // configSubtree must be valid and outlive this component.
    explicit Preferences (juce::ValueTree configSubtree);

    ~Preferences() override;

    // -------------------------------------------------------------------------
    // Component overrides
    // -------------------------------------------------------------------------
    void paint       (jam::tui::Graphics& g)   override;
    void handleInput (const juce::KeyPress& key) override;
    void resized     ()                        override;

    // -------------------------------------------------------------------------
    // ValueTree::Listener
    // -------------------------------------------------------------------------
    void valueTreePropertyChanged (juce::ValueTree&        tree,
                                   const juce::Identifier& property) override;

    // -------------------------------------------------------------------------
    // Callback — fired when user selects "Back to Menu" row
    // -------------------------------------------------------------------------
    std::function<void()> onBackToMenu;

private:
    // =========================================================================
    // Internal constants
    // =========================================================================

    // SPEC §Preferences — interval range 1–60 min
    static constexpr int INTERVAL_MIN      { 1 };
    static constexpr int INTERVAL_MAX      { 60 };
    static constexpr int INTERVAL_STEP_SMALL { 1 };
    static constexpr int INTERVAL_STEP_LARGE { 10 };

    // Row indices — 5 fixed rows
    static constexpr int ROW_AUTO_UPDATE { 0 };
    static constexpr int ROW_INTERVAL    { 1 };
    static constexpr int ROW_SEPARATOR   { 2 };
    static constexpr int ROW_THEME       { 3 };
    static constexpr int ROW_BACK        { 4 };

    // =========================================================================
    // Helpers
    // =========================================================================

    // Rebuilds all menu rows from current CONFIG VT state.
    void rebuildRows();

    // Adjusts autoScanInterval by delta, clamped to [INTERVAL_MIN, INTERVAL_MAX].
    void adjustInterval (int delta);

    // Fires onItemSelected action for the currently selected row.
    void dispatchRowAction (int selectedIndex);

    // =========================================================================
    // VT reference — observer into root tree (not owned here)
    // =========================================================================

    juce::ValueTree configTree;

    // =========================================================================
    // Child component
    // =========================================================================

    jam::tui::Menu menu;

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (Preferences)
};
