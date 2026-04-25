#include "State.h"

// ============================================================================
// writeIfChanged — compare / write / snapshot
// ============================================================================

template <typename T>
bool State::writeIfChanged (T&                     snapshotSlot,
                            T                      current,
                            juce::ValueTree&        subtree,
                            const juce::Identifier& key,
                            const juce::var&        asVar) noexcept
{
    bool didWrite { false };

    if (snapshotSlot != current)
    {
        subtree.setProperty (key, asVar, nullptr);
        snapshotSlot = current;
        didWrite     = true;
    }

    return didWrite;
}

// ============================================================================
// flushProject
// ============================================================================

bool State::flushProject() noexcept
{
    juce::ValueTree project { tree.getChildWithName (ID::PROJECT) };
    bool changed { false };

    const int generatorVal { selectedGeneratorAtom.load (std::memory_order_relaxed) };
    if (writeIfChanged (lastFlushed.selectedGenerator, generatorVal, project, ID::selectedGenerator,
                        toString (static_cast<Generator> (generatorVal))))
        changed = true;

    const int configVal { configurationAtom.load (std::memory_order_relaxed) };
    if (writeIfChanged (lastFlushed.configuration, configVal, project, ID::configuration,
                        toString (static_cast<Configuration> (configVal))))
        changed = true;

    const bool hasCmake { hasCMakeListsAtom.load (std::memory_order_relaxed) };
    if (writeIfChanged (lastFlushed.hasCMakeLists, hasCmake, project, ID::hasCMakeLists, juce::var { hasCmake }))
        changed = true;

    return changed;
}

// ============================================================================
// flushAsync
// ============================================================================

bool State::flushAsync() noexcept
{
    juce::ValueTree async { tree.getChildWithName (ID::ASYNC) };
    bool changed { false };

    const bool activeVal { isActiveAtom.load (std::memory_order_relaxed) };
    if (writeIfChanged (lastFlushed.isActive, activeVal, async, ID::isActive, juce::var { activeVal }))
        changed = true;

    const bool abortedVal { isAbortedAtom.load (std::memory_order_relaxed) };
    if (writeIfChanged (lastFlushed.isAborted, abortedVal, async, ID::isAborted, juce::var { abortedVal }))
        changed = true;

    const int opVal { currentOpAtom.load (std::memory_order_relaxed) };
    if (writeIfChanged (lastFlushed.currentOp, opVal, async, ID::currentOp,
                        toString (static_cast<OpType> (opVal))))
        changed = true;

    return changed;
}

// ============================================================================
// flushConfig
// ============================================================================

bool State::flushConfig() noexcept
{
    juce::ValueTree config { tree.getChildWithName (ID::CONFIG) };
    bool changed { false };

    const bool scanEnabled { autoScanEnabledAtom.load (std::memory_order_relaxed) };
    if (writeIfChanged (lastFlushed.autoScanEnabled, scanEnabled, config, ID::autoScanEnabled, juce::var { scanEnabled }))
        changed = true;

    const int scanInterval { autoScanIntervalAtom.load (std::memory_order_relaxed) };
    if (writeIfChanged (lastFlushed.autoScanInterval, scanInterval, config, ID::autoScanInterval, juce::var { scanInterval }))
        changed = true;

    return changed;
}

// ============================================================================
// flushStrings
// ============================================================================

bool State::flushStrings() noexcept
{
    juce::ValueTree project { tree.getChildWithName (ID::PROJECT) };
    juce::ValueTree config  { tree.getChildWithName (ID::CONFIG) };
    bool changed { false };

    if (writeIfChanged (lastFlushed.workingDirectory, workingDirectory, project, ID::workingDirectory,
                        juce::var { workingDirectory }))
        changed = true;

    if (writeIfChanged (lastFlushed.themeName, themeName, config, ID::theme,
                        juce::var { themeName }))
        changed = true;

    return changed;
}
