#include <JuceHeader.h>
#include "MenuBuilder.h"
#include "../Identifier.h"
#include "../state/Axis.h"

// ============================================================================
// File-scope helpers
// ============================================================================
//
// Each helper is a pure function on VT subtrees.  No member access, no side
// effects, no early returns.  Named namespace and static linkage per CONTRACT.

namespace menu
{

// ---------------------------------------------------------------------------
// buildGenerateLabel
// ---------------------------------------------------------------------------
// Go menu.go: regenerateLabel based on hasBuild.
// Returns "Generate" when no build exists, "Regenerate" otherwise.

static juce::String buildGenerateLabel (bool hasBuild)
{
    juce::String result { "Generate" };

    if (hasBuild)
        result = "Regenerate";

    return result;
}

// ---------------------------------------------------------------------------
// buildOpenLabel
// ---------------------------------------------------------------------------
// Go menu.go openIdeLabel: label depends on isIDEGenerator.

static juce::String buildOpenLabel (bool isIDEGenerator)
{
    juce::String result { "Open Editor" };

    if (isIDEGenerator)
        result = "Open IDE";

    return result;
}

// ---------------------------------------------------------------------------
// findBuildForGenerator
// ---------------------------------------------------------------------------
// Returns the first BUILD child whose generator property matches generatorName.
// Returns an invalid ValueTree when none found.

static juce::ValueTree findBuildForGenerator (const juce::ValueTree& buildsSubtree,
                                              const juce::String&    generatorName)
{
    juce::ValueTree result;

    for (int i { 0 }; i < buildsSubtree.getNumChildren(); ++i)
    {
        const juce::ValueTree child { buildsSubtree.getChild (i) };

        if (child.getProperty (ID::generator).toString() == generatorName)
        {
            result = child;
            break;
        }
    }

    return result;
}

// ---------------------------------------------------------------------------
// hasAnyExistingBuild
// ---------------------------------------------------------------------------
// Returns true when any BUILD child has exists == true.
// SPEC §Menu Item Selectability Rules: "Clean All: selectable if any build
// directory exists".

static bool hasAnyExistingBuild (const juce::ValueTree& buildsSubtree)
{
    bool result { false };

    for (int i { 0 }; i < buildsSubtree.getNumChildren(); ++i)
    {
        const juce::ValueTree child { buildsSubtree.getChild (i) };

        if (static_cast<bool> (child.getProperty (ID::exists)))
        {
            result = true;
            break;
        }
    }

    return result;
}

// ---------------------------------------------------------------------------
// buildAttributedRow
// ---------------------------------------------------------------------------
// Produces a single-row AttributedString: "  <label><pad><value>  <shortcut>"
// Plain white text — color theming is Step 7.

static juce::AttributedString buildAttributedRow (const juce::String& label,
                                                  const juce::String& value,
                                                  const juce::String& shortcut)
{
    juce::String text { "  " + label };

    if (value.isNotEmpty())
    {
        text += "  " + value;
    }

    if (shortcut.isNotEmpty())
    {
        text += "  [" + shortcut + "]";
    }

    juce::AttributedString result;
    result.append (text, juce::Font (juce::FontOptions{}), juce::Colours::white);
    return result;
}

// ---------------------------------------------------------------------------
// buildSeparatorRow
// ---------------------------------------------------------------------------
// Produces the visual separator row (horizontal line characters).

static juce::AttributedString buildSeparatorRow()
{
    juce::AttributedString result;
    result.append (juce::String::repeatedString ("─", 32),
                   juce::Font (juce::FontOptions{}),
                   juce::Colours::grey);
    return result;
}

// ============================================================================
// buildMenu — public API
// ============================================================================

BuiltMenu buildMenu (const juce::ValueTree& projectSubtree,
                     const juce::ValueTree& buildsSubtree)
{
    jassert (projectSubtree.isValid());
    jassert (buildsSubtree.isValid());

    // --- Read project properties ---
    const juce::String selectedGeneratorName {
        projectSubtree.getProperty (ID::selectedGenerator).toString()
    };
    const juce::String configurationName {
        projectSubtree.getProperty (ID::configuration).toString()
    };

    // --- Resolve selected generator's isIDE flag ---
    bool isIDEGenerator { false };

    for (int i { 0 }; i < projectSubtree.getNumChildren(); ++i)
    {
        const juce::ValueTree child { projectSubtree.getChild (i) };

        if (child.getType() == ID::GENERATOR
            and child.getProperty (ID::name).toString() == selectedGeneratorName)
        {
            isIDEGenerator = static_cast<bool> (child.getProperty (ID::isIDE));
            break;
        }
    }

    // --- Resolve build state for selected generator ---
    const juce::ValueTree selectedBuild { findBuildForGenerator (buildsSubtree, selectedGeneratorName) };
    const bool hasBuild       { selectedBuild.isValid() and static_cast<bool> (selectedBuild.getProperty (ID::exists)) };
    const bool isConfigured   { hasBuild and static_cast<bool> (selectedBuild.getProperty (ID::isConfigured)) };
    const bool hasGeneratorSelected { selectedGeneratorName.isNotEmpty() };
    const bool canOpenIDE     { hasGeneratorSelected and hasBuild };
    const bool canClean       { hasBuild };
    const bool hasBuildsToClean { hasAnyExistingBuild (buildsSubtree) };

    // --- SPEC §Menu Item Selectability Rules: 8 fixed rows ---
    // Row 0: Project — ALWAYS selectable
    // Row 1: Generate/Regenerate — ALWAYS selectable
    // Row 2: Open IDE / Open Editor — selectable if canOpenIDE
    // Row 3: Separator — NOT selectable
    // Row 4: Configuration — ALWAYS selectable
    // Row 5: Build — selectable if hasBuild AND isConfigured
    // Row 6: Clean — selectable if canClean
    // Row 7: Clean All — selectable if hasBuildsToClean

    BuiltMenu result;

    // Row 0: Project
    result.items.add (buildAttributedRow ("Project", selectedGeneratorName, ""));
    result.selectable.add (true);
    result.rowIds.add (ROW_PROJECT);

    // Row 1: Generate / Regenerate
    result.items.add (buildAttributedRow (buildGenerateLabel (hasBuild), "", "g"));
    result.selectable.add (true);
    result.rowIds.add (ROW_GENERATE);

    // Row 2: Open IDE / Open Editor
    result.items.add (buildAttributedRow (buildOpenLabel (isIDEGenerator), "", "o"));
    result.selectable.add (canOpenIDE);
    result.rowIds.add (ROW_OPEN);

    // Row 3: Separator
    result.items.add (buildSeparatorRow());
    result.selectable.add (false);
    result.rowIds.add (ROW_SEPARATOR);

    // Row 4: Configuration
    result.items.add (buildAttributedRow ("Configuration", configurationName, ""));
    result.selectable.add (true);
    result.rowIds.add (ROW_CONFIGURATION);

    // Row 5: Build
    result.items.add (buildAttributedRow ("Build", "", "b"));
    result.selectable.add (hasBuild and isConfigured);
    result.rowIds.add (ROW_BUILD);

    // Row 6: Clean
    result.items.add (buildAttributedRow ("Clean", "", "c"));
    result.selectable.add (canClean);
    result.rowIds.add (ROW_CLEAN);

    // Row 7: Clean All
    result.items.add (buildAttributedRow ("Clean All", "", "x"));
    result.selectable.add (hasBuildsToClean);
    result.rowIds.add (ROW_CLEAN_ALL);

    return result;
}

} // namespace menu
