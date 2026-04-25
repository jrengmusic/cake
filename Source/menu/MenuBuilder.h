#pragma once
#include <JuceHeader.h>
#include "MenuItems.h"

// ============================================================================
// MenuBuilder — pure function from VT state to menu item data
// ============================================================================
//
// Mirrors Go GenerateMenuRows(): reads PROJECT + BUILDS subtrees, returns
// exactly 8 rows with their labels, selectability, and row IDs.
// No stored state — stateless per BLESSED-S.
// Called by MainComponent whenever PROJECT or BUILDS VT subtrees change.
//
// SPEC §Menu Item Selectability Rules; BLESSED-D (same VT → same output).

namespace menu
{

// ============================================================================
// BuiltMenu
// ============================================================================
//
// Output of buildMenu(). Three parallel arrays, one entry per row (always 8).

struct BuiltMenu
{
    juce::Array<juce::AttributedString> items;      // row text, one per row
    juce::Array<bool>                   selectable; // per-row selectability
    juce::Array<juce::String>           rowIds;     // row ID for action dispatch
};

// ============================================================================
// buildMenu
// ============================================================================
//
// Pure function. Reads projectSubtree (PROJECT node) and buildsSubtree
// (BUILDS node) and returns exactly 8 rows.
//
// Precondition: projectSubtree and buildsSubtree must be valid ValueTree nodes.
// SPEC §Menu Item Selectability Rules — selectability logic lives here only.

BuiltMenu buildMenu (const juce::ValueTree& projectSubtree,
                     const juce::ValueTree& buildsSubtree);

} // namespace menu
