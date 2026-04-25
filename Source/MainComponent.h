#pragma once
#include <JuceHeader.h>
#include "Identifier.h"
#include "state/State.h"
#include "menu/MenuBuilder.h"
#include "component/Header.h"
#include "component/Footer.h"
#include "component/Banner.h"

// ============================================================================
// MainComponent
// ============================================================================
//
// Root component and Controller for CAKE.  Owns the four front-view components
// and orchestrates data flow:
//
//   Input → handleInput() → VT mutation
//   VT change → valueTreeChildAdded / valueTreePropertyChanged → rebuildMenu()
//   rebuildMenu() → jam::tui::Menu::setItems() + setItemSelectable()
//
// Owns (value members — deterministic lifecycle, BLESSED-B):
//   Header           — project directory
//   Footer           — hint text
//   Banner           — braille logo in right pane
//   jam::tui::Menu   — 8-row menu primitive
//
// Layout: 65% left (Header + Menu + Footer, stacked),
//         35% right (Banner).
//
// BLESSED-E (Controller, never View), BLESSED-S2 (reads/writes VT only).
// SPEC §Architecture: "MainComponent is sole orchestrator".

class MainComponent : public jam::tui::Component,
                      public juce::ValueTree::Listener
{
public:
    explicit MainComponent (State& state);
    ~MainComponent() override;

    // -------------------------------------------------------------------------
    // Component overrides
    // -------------------------------------------------------------------------
    void paint       (jam::tui::Graphics& g)             override;
    void handleInput (const juce::KeyPress& key)         override;
    void resized     ()                                  override;

    // -------------------------------------------------------------------------
    // ValueTree::Listener
    // -------------------------------------------------------------------------
    void valueTreePropertyChanged (juce::ValueTree&         tree,
                                   const juce::Identifier&  property) override;
    void valueTreeChildAdded      (juce::ValueTree& parent,
                                   juce::ValueTree& child)             override;
    void valueTreeChildRemoved    (juce::ValueTree& parent,
                                   juce::ValueTree& child,
                                   int              oldIndex)          override;

private:
    // =========================================================================
    // Internal constants
    // =========================================================================

    // SPEC §Layout Structure: 65% left, 35% right.
    static constexpr int leftPanePercent { 65 };

    // SPEC §Layout Structure: Header occupies top 3 rows, Footer occupies 1 row.
    static constexpr int headerRows { 3 };
    static constexpr int footerRows { 1 };

    // Default footer hint per SPEC §Keyboard Shortcuts.
    static const juce::String menuHint;

    // =========================================================================
    // State reference (non-owning, BLESSED-B)
    // =========================================================================

    State& stateRef;

    // =========================================================================
    // VT state references — observer references into the root tree (not owned here)
    // =========================================================================

    juce::ValueTree projectState;
    juce::ValueTree buildsState;

    // =========================================================================
    // Child components — value members, deterministic lifecycle
    // =========================================================================

    Header             header;
    Footer             footer;
    Banner             banner;
    jam::tui::Menu     menu;

    // =========================================================================
    // LookAndFeel
    // =========================================================================

    jam::tui::LookAndFeel lookAndFeel;

    // =========================================================================
    // Cached menu data — updated by rebuildMenu(), used by dispatchMenuAction()
    // =========================================================================

    menu::BuiltMenu lastBuiltMenu;

    // =========================================================================
    // Helpers
    // =========================================================================

    // Rebuilds the jam::tui::Menu from current VT state.
    // Called from VT listener callbacks and on construction.
    void rebuildMenu();

    // Dispatches action for the selected menu row.
    // rowId comes from menu::BuiltMenu::rowIds.
    void dispatchMenuAction (const juce::String& rowId);

    // Cycles selected generator to next in GENERATOR list.
    void cycleSelectedGenerator();

    // Toggles Configuration between Debug and Release.
    void toggleConfiguration();

    // Finds the index of the currently selected generator in projectState children.
    // Returns 0 if not found.
    int findCurrentGeneratorIndex() const;

    // Computes the four paint regions from given bounds.
    struct PaintRegions
    {
        jam::tui::Rectangle header;
        jam::tui::Rectangle content;
        jam::tui::Rectangle footer;
        jam::tui::Rectangle banner;
    };
    static PaintRegions computePaintRegions (const jam::tui::Rectangle& bounds);

    // Paints header and footer chrome.
    void paintChrome (jam::tui::Graphics& g, const PaintRegions& regions);

    // Menu-mode input handler.
    void handleMenuInput (const juce::KeyPress& key);

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (MainComponent)
};
