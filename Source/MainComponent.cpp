#include <JuceHeader.h>
#include "MainComponent.h"
#include "menu/MenuItems.h"

// ============================================================================
// File-scope constants
// ============================================================================

const juce::String MainComponent::menuHint {
    "up/down navigate | Enter select | Ctrl+C quit"
};

// ============================================================================
// Construction / Destruction
// ============================================================================

MainComponent::MainComponent (State& state)
    : stateRef     { state }
    , projectState { state.getProjectState() }
    , buildsState  { state.getBuildsState()  }
{
    jassert (projectState.isValid());
    jassert (buildsState.isValid());

    setLookAndFeel (&lookAndFeel);

    // Register as ValueTree listener on PROJECT and BUILDS subtrees.
    projectState.addListener (this);
    buildsState.addListener  (this);

    // Wire menu action dispatch.
    menu.onItemSelected = [this] (int selectedIndex)
    {
        jassert (selectedIndex >= 0 and selectedIndex < lastBuiltMenu.rowIds.size());
        dispatchMenuAction (lastBuiltMenu.rowIds.getReference (selectedIndex));
    };

    footer.setHint (menuHint);
    header.setWorkingDirectory (projectState.getProperty (ID::workingDirectory).toString());

    addChildComponent (header);
    addChildComponent (menu);
    addChildComponent (footer);
    addChildComponent (banner);

    rebuildMenu();
}

MainComponent::~MainComponent()
{
    setLookAndFeel (nullptr);
    buildsState.removeListener  (this);
    projectState.removeListener (this);
}

// ============================================================================
// Component overrides
// ============================================================================

MainComponent::PaintRegions MainComponent::computePaintRegions (const jam::tui::Rectangle& bounds)
{
    const int leftWidth     { (bounds.width * leftPanePercent) / 100 };
    const int rightWidth    { bounds.width - leftWidth };
    const int contentHeight { bounds.height - headerRows - footerRows };

    return {
        { jam::literals::Cell { 0 },         jam::literals::Cell { 0 },
          jam::literals::Cell { leftWidth },  jam::literals::Cell { headerRows } },
        { jam::literals::Cell { 0 },         jam::literals::Cell { headerRows },
          jam::literals::Cell { leftWidth },  jam::literals::Cell { contentHeight } },
        { jam::literals::Cell { 0 },         jam::literals::Cell { bounds.height - footerRows },
          jam::literals::Cell { leftWidth },  jam::literals::Cell { footerRows } },
        { jam::literals::Cell { leftWidth }, jam::literals::Cell { 0 },
          jam::literals::Cell { rightWidth }, jam::literals::Cell { bounds.height } }
    };
}

void MainComponent::paintChrome (jam::tui::Graphics& g, const PaintRegions& regions)
{
    jam::tui::Graphics headerG { g.clip (regions.header) };
    header.paint (headerG);

    jam::tui::Graphics footerG { g.clip (regions.footer) };
    footer.paint (footerG);
}

void MainComponent::paint (jam::tui::Graphics& g)
{
    const jam::tui::Rectangle bounds { getCellBounds() };

    if (not bounds.isEmpty())
    {
        const PaintRegions regions { computePaintRegions (bounds) };

        paintChrome (g, regions);

        jam::tui::Graphics menuG { g.clip (regions.content) };
        menu.paint (menuG);

        jam::tui::Graphics bannerG { g.clip (regions.banner) };
        banner.paint (bannerG);
    }
}

// ============================================================================
// handleInput
// ============================================================================

void MainComponent::handleInput (const juce::KeyPress& key)
{
    if (key == juce::KeyPress ('c', juce::ModifierKeys::ctrlModifier, 0))
    {
        juce::JUCEApplication::getInstance()->systemRequestedQuit();
    }
    else
    {
        handleMenuInput (key);
    }
}

void MainComponent::handleMenuInput (const juce::KeyPress& key)
{
    const juce::juce_wchar ch { key.getTextCharacter() };
    const bool noModifiers { not key.getModifiers().isAnyModifierKeyDown() };

    if (noModifiers and ch != 0)
    {
        static const std::unordered_map<juce::juce_wchar, juce::String> menuShortcuts {
            { 'g', menu::ROW_GENERATE  },
            { 'b', menu::ROW_BUILD     },
            { 'o', menu::ROW_OPEN      },
            { 'c', menu::ROW_CLEAN     },
            { 'x', menu::ROW_CLEAN_ALL }
        };

        const auto it { menuShortcuts.find (ch) };

        if (it != menuShortcuts.end())
            dispatchMenuAction (it->second);
        else
            menu.handleInput (key);
    }
    else
    {
        menu.handleInput (key);
    }
}

void MainComponent::resized()
{
    const juce::Rectangle<int> bounds { getBounds() };

    if (not bounds.isEmpty())
    {
        const int leftWidth     { (bounds.getWidth() * leftPanePercent) / 100 };
        const int totalHeight   { bounds.getHeight() };
        const int contentHeight { totalHeight - headerRows - footerRows };
        const int rightWidth    { bounds.getWidth() - leftWidth };

        header.setBounds (0, 0,                        leftWidth,  headerRows);
        footer.setBounds (0, totalHeight - footerRows, leftWidth,  footerRows);
        menu.setBounds   (0, headerRows,               leftWidth,  contentHeight);
        banner.setBounds (leftWidth, 0,                rightWidth, totalHeight);
    }
}

// ============================================================================
// ValueTree::Listener
// ============================================================================

void MainComponent::valueTreePropertyChanged (juce::ValueTree& tree,
                                              const juce::Identifier& property)
{
    juce::ignoreUnused (tree);
    juce::ignoreUnused (property);

    header.setWorkingDirectory (projectState.getProperty (ID::workingDirectory).toString());
    rebuildMenu();
    repaint();
}

void MainComponent::valueTreeChildAdded (juce::ValueTree& parent,
                                         juce::ValueTree& child)
{
    juce::ignoreUnused (parent);
    juce::ignoreUnused (child);

    rebuildMenu();
    repaint();
}

void MainComponent::valueTreeChildRemoved (juce::ValueTree& parent,
                                           juce::ValueTree& child,
                                           int              oldIndex)
{
    juce::ignoreUnused (parent);
    juce::ignoreUnused (child);
    juce::ignoreUnused (oldIndex);

    rebuildMenu();
    repaint();
}

// ============================================================================
// Private helpers
// ============================================================================

void MainComponent::rebuildMenu()
{
    lastBuiltMenu = menu::buildMenu (projectState, buildsState);

    menu.setItems (lastBuiltMenu.items);

    for (int i { 0 }; i < lastBuiltMenu.selectable.size(); ++i)
        menu.setItemSelectable (i, lastBuiltMenu.selectable.getUnchecked (i));
}

void MainComponent::dispatchMenuAction (const juce::String& rowId)
{
    using Action = void (MainComponent::*)();

    static const std::unordered_map<std::string, Action> actionMap {
        { menu::ROW_PROJECT,       &MainComponent::cycleSelectedGenerator },
        { menu::ROW_CONFIGURATION, &MainComponent::toggleConfiguration    }
    };

    const auto it { actionMap.find (rowId.toStdString()) };

    if (it != actionMap.end())
        (this->*it->second)();
}

int MainComponent::findCurrentGeneratorIndex() const
{
    const juce::String current { projectState.getProperty (ID::selectedGenerator).toString() };
    const int generatorCount   { projectState.getNumChildren() };
    int result                 { 0 };

    for (int i { 0 }; i < generatorCount; ++i)
    {
        const juce::ValueTree child { projectState.getChild (i) };

        if (child.getType() == ID::GENERATOR
            and child.getProperty (ID::name).toString() == current)
        {
            result = i;
            break;
        }
    }

    return result;
}

void MainComponent::cycleSelectedGenerator()
{
    const int generatorCount { projectState.getNumChildren() };

    if (generatorCount > 1)
    {
        const int nextIndex                  { (findCurrentGeneratorIndex() + 1) % generatorCount };
        const juce::ValueTree nextGenerator  { projectState.getChild (nextIndex) };

        if (nextGenerator.isValid())
        {
            projectState.setProperty (
                ID::selectedGenerator,
                nextGenerator.getProperty (ID::name),
                nullptr
            );
        }
    }
}

void MainComponent::toggleConfiguration()
{
    const juce::String current { projectState.getProperty (ID::configuration).toString() };
    const juce::String next    {
        (current == toString (Configuration::debug))
            ? toString (Configuration::release)
            : toString (Configuration::debug)
    };

    projectState.setProperty (ID::configuration, next, nullptr);
}
