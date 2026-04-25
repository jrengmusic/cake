#include <JuceHeader.h>
#include <iostream>
#include "MainComponent.h"
#include "state/State.h"
#include "detect/GeneratorDetector.h"
#include "Identifier.h"

class CakeApp : public juce::JUCEApplication
{
public:
    CakeApp() = default;

    const juce::String getApplicationName() override
    {
        return "cakec";
    }

    const juce::String getApplicationVersion() override
    {
        return "0.0.0";
    }

    bool moreThanOneInstanceAllowed() override
    {
        return true;
    }

    void initialise (const juce::String&) override
    {
        // -------------------------------------------------------------------------
        // 1. State — owns VT schema, atom bridge, flush timer
        // -------------------------------------------------------------------------
        state = std::make_unique<State>();
        state->start();

        // -------------------------------------------------------------------------
        // 1b. Generator detection + initial VT population
        //
        // Run on message thread at startup (before MainComponent attaches listeners).
        // Writes directly to VT — no atom round-trip needed for static detection
        // results that never change during the session.
        // MESSAGE THREAD.
        // -------------------------------------------------------------------------
        populateDetectedState();

        // -------------------------------------------------------------------------
        // 2. MainComponent
        // -------------------------------------------------------------------------
        mainComponent = std::make_unique<MainComponent> (*state);

        // -------------------------------------------------------------------------
        // 2. ansi::Screen — Writer is a value member; Screen holds a reference
        // -------------------------------------------------------------------------
        screen = std::make_unique<jam::tui::Screen> (writer);

        // Size Screen and MainComponent to the current terminal dimensions
        const juce::Rectangle<int> termBounds { jam::tui::getBounds().toJuce() };
        screen->setBounds (termBounds);
        mainComponent->setBounds (termBounds.withPosition (0, 0));

        screen->addAndMakeVisible (*mainComponent);

        // -------------------------------------------------------------------------
        // 3. Input — route keys to MainComponent; route resize to Screen
        // -------------------------------------------------------------------------
        input = std::make_unique<jam::tui::Input>();

        input->start (
            [this] (juce::KeyPress key)
            {
                mainComponent->handleInput (key);
            },
            [this] (juce::String content)
            {
                juce::ignoreUnused (content);
                // No TextBox at app root — paste not consumed here.
            },
            [this]
            {
                const juce::Rectangle<int> newBounds { jam::tui::getBounds().toJuce() };
                screen->setBounds (newBounds);
                mainComponent->setBounds (newBounds.withPosition (0, 0));
                screen->onTerminalResized();
            }
        );

        // -------------------------------------------------------------------------
        // 4. Start render loop
        // -------------------------------------------------------------------------
        screen->start();
    }

    void shutdown() override
    {
        // Stop in reverse construction order (BLESSED-B).
        // Input thread first so no callbacks fire during teardown.
        if (input != nullptr)
            input->stop();

        // std::unique_ptr destructors handle MainComponent, Screen, Input.
        input.reset();
        screen.reset();
        mainComponent.reset();

        // State last — flush timer must stop after all views are gone.
        if (state != nullptr)
            state->stop();

        state.reset();

        // Restore cursor after all rendering is torn down — Screen's destructor
        // may emit a final render that hides the cursor.
        std::cout << ANSI::CURSOR_SHOW << std::flush;
    }

    void anotherInstanceStarted (const juce::String&) override
    {
    }

private:

    // =========================================================================
    // populateDetectedState
    // =========================================================================
    //
    // Runs once at startup (message thread) after State construction.
    // Detects generators and project state, then writes directly to State VT
    // and updates atoms so the flush timer reflects the initial state.
    //
    // Writing to VT directly here is correct: listeners are not yet attached
    // (MainComponent constructed next), and these are one-time static writes.
    // MESSAGE THREAD.
    void populateDetectedState()
    {
        // --- Working directory ---
        const juce::String cwd { juce::File::getCurrentWorkingDirectory().getFullPathName() };
        state->setWorkingDirectory (cwd);

        // --- CMakeLists.txt presence ---
        const juce::File cmakeListsFile { juce::File::getCurrentWorkingDirectory()
                                              .getChildFile ("CMakeLists.txt") };
        state->setHasCMakeLists (cmakeListsFile.existsAsFile());

        // --- Generator detection ---
        const juce::Array<GeneratorDescriptor> detectedGenerators {
            GeneratorDetector::detectAvailableGenerators()
        };

        juce::ValueTree projectSubtree { state->getProjectState() };

        // Remove any pre-existing GENERATOR children (idempotent re-detection).
        projectSubtree.removeAllChildren (nullptr);

        for (const GeneratorDescriptor& descriptor : detectedGenerators)
        {
            juce::ValueTree generatorNode { ID::GENERATOR };
            generatorNode.setProperty (ID::name,  toString (descriptor.generator), nullptr);
            generatorNode.setProperty (ID::isIDE, descriptor.isIDE,                nullptr);
            projectSubtree.addChild (generatorNode, -1, nullptr);
        }

        // Set selectedGenerator to first detected generator (if any detected).
        if (not detectedGenerators.isEmpty())
            state->setSelectedGenerator (detectedGenerators.getFirst().generator);

        // Force flush so the VT reflects detection results before MainComponent reads it.
        state->flush();
    }

    // =========================================================================
    // Members
    // =========================================================================

    jam::tui::Writer                    writer;
    std::unique_ptr<State>              state;
    std::unique_ptr<MainComponent>      mainComponent;
    std::unique_ptr<jam::tui::Screen>   screen;
    std::unique_ptr<jam::tui::Input>    input;

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (CakeApp)
};

START_JUCE_APPLICATION (CakeApp)
