#include <JuceHeader.h>
#include "Preferences.h"

// ============================================================================
// File-scope helpers
// ============================================================================
//
// Pure functions on CONFIG VT properties.  No member access, no side effects.
// Static linkage per CONTRACT (no anonymous namespaces).

static juce::AttributedString buildPrefRow (const juce::String& label,
                                            const juce::String& value)
{
    const juce::String text { "  " + label + "  " + value };

    juce::AttributedString result;
    result.append (text, juce::Font (juce::FontOptions{}), juce::Colours::white);
    return result;
}

static juce::AttributedString buildSeparatorRow()
{
    juce::AttributedString result;
    result.append (juce::String::repeatedString ("─", 32),
                   juce::Font (juce::FontOptions{}),
                   juce::Colours::grey);
    return result;
}

// ============================================================================
// Construction / Destruction
// ============================================================================

Preferences::Preferences (juce::ValueTree configSubtree)
    : configTree { configSubtree }
{
    jassert (configTree.isValid());

    configTree.addListener (this);

    // Wire row-action dispatch.
    menu.onItemSelected = [this] (int selectedIndex)
    {
        dispatchRowAction (selectedIndex);
    };

    addChildComponent (menu);
    rebuildRows();
}

Preferences::~Preferences()
{
    configTree.removeListener (this);
}

// ============================================================================
// Component overrides
// ============================================================================

void Preferences::paint (jam::tui::Graphics& g)
{
    const jam::tui::Rectangle bounds { getCellBounds() };

    if (not bounds.isEmpty())
    {
        jam::tui::Graphics menuG { g.clip (bounds) };
        menu.paint (menuG);
    }
}

void Preferences::handleInput (const juce::KeyPress& key)
{
    const int selected { menu.getSelectedIndex() };
    const juce::juce_wchar ch { key.getTextCharacter() };
    const bool isShift { key.getModifiers().isShiftDown() };

    if (selected == ROW_INTERVAL and ch != 0)
    {
        // SPEC §Settings Behavior:
        //   +/=       → decrease by 1 min  (+ = Shift+=, = is bare equals)
        //   -/_       → increase by 1 min  (_ = Shift+-)
        //   Shift+=   → increase by 10 min (character '+', shift=true)
        //   Shift+-   → decrease by 10 min (character '_', shift=true)
        //
        // On a US keyboard: Shift+= produces '+' (shift=true, char='+')
        //                   Shift+- produces '_' (shift=true, char='_')
        //                   bare '=' is shift=false, char='='
        //                   bare '-' is shift=false, char='-'
        if (ch == '+' and isShift)
        {
            adjustInterval (INTERVAL_STEP_LARGE);
        }
        else if (ch == '_' and isShift)
        {
            adjustInterval (-INTERVAL_STEP_LARGE);
        }
        else if (ch == '=' or ch == '+')
        {
            adjustInterval (-INTERVAL_STEP_SMALL);
        }
        else if (ch == '-' or ch == '_')
        {
            adjustInterval (INTERVAL_STEP_SMALL);
        }
        else
        {
            menu.handleInput (key);
        }
    }
    else
    {
        menu.handleInput (key);
    }
}

void Preferences::resized()
{
    const juce::Rectangle<int> bounds { getBounds() };

    if (not bounds.isEmpty())
        menu.setBounds (bounds);
}

// ============================================================================
// ValueTree::Listener
// ============================================================================

void Preferences::valueTreePropertyChanged (juce::ValueTree&        tree,
                                            const juce::Identifier& property)
{
    juce::ignoreUnused (tree);
    juce::ignoreUnused (property);

    rebuildRows();
    repaint();
}

// ============================================================================
// Private helpers
// ============================================================================

void Preferences::rebuildRows()
{
    const bool autoScanEnabled {
        static_cast<bool> (configTree.getProperty (ID::autoScanEnabled))
    };
    const int autoScanInterval {
        static_cast<int> (configTree.getProperty (ID::autoScanInterval))
    };
    const juce::String themeName {
        configTree.getProperty (ID::theme).toString()
    };

    const juce::String autoScanLabel { autoScanEnabled ? "ON" : "OFF" };
    const juce::String intervalLabel { juce::String (autoScanInterval) + " min" };

    juce::Array<juce::AttributedString> items;
    items.add (buildPrefRow ("Auto-update",     autoScanLabel));
    items.add (buildPrefRow ("Update Interval", intervalLabel));
    items.add (buildSeparatorRow());
    items.add (buildPrefRow ("Theme",           themeName));
    items.add (buildPrefRow ("Back to Menu",    ""));

    menu.setItems (items);

    // Row selectability — separator not selectable.
    // Interval row: only selectable when auto-scan is ON (SPEC §Update Interval Adjustment).
    menu.setItemSelectable (ROW_AUTO_UPDATE, true);
    menu.setItemSelectable (ROW_INTERVAL,    autoScanEnabled);
    menu.setItemSelectable (ROW_SEPARATOR,   false);
    menu.setItemSelectable (ROW_THEME,       true);
    menu.setItemSelectable (ROW_BACK,        true);
}

void Preferences::adjustInterval (int delta)
{
    const int current { static_cast<int> (configTree.getProperty (ID::autoScanInterval)) };
    const int next    { juce::jlimit (INTERVAL_MIN, INTERVAL_MAX, current + delta) };

    configTree.setProperty (ID::autoScanInterval, next, nullptr);
}

void Preferences::dispatchRowAction (int selectedIndex)
{
    if (selectedIndex == ROW_AUTO_UPDATE)
    {
        const bool current { static_cast<bool> (configTree.getProperty (ID::autoScanEnabled)) };
        configTree.setProperty (ID::autoScanEnabled, not current, nullptr);
    }
    else if (selectedIndex == ROW_INTERVAL)
    {
        // Interval is adjusted via +/- keys; Enter on this row does nothing.
    }
    else if (selectedIndex == ROW_THEME)
    {
        // SPEC §Theme Cycle: gfx → spring → summer → autumn → winter → gfx.
        static const juce::StringArray THEMES { "gfx", "spring", "summer", "autumn", "winter" };

        const juce::String current { configTree.getProperty (ID::theme).toString() };
        const int currentIndex     { THEMES.indexOf (current) };
        const int nextIndex        { (currentIndex + 1) % THEMES.size() };

        configTree.setProperty (ID::theme, THEMES.getReference (nextIndex), nullptr);
    }
    else if (selectedIndex == ROW_BACK)
    {
        if (onBackToMenu != nullptr)
            onBackToMenu();
    }
}
