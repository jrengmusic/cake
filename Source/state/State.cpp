#include "State.h"

// ============================================================================
// Construction / Destruction
// ============================================================================

State::State()
{
    // Build ValueTree skeleton per SPEC §ValueTree Schema — all nodes present,
    // defaults per SPEC §State Model.  Tree is message-thread owned; constructed
    // before Timer starts.

    tree = juce::ValueTree { ID::CAKE };

    // PROJECT node
    juce::ValueTree project { ID::PROJECT };
    project.setProperty (ID::workingDirectory,  juce::String{},                         nullptr);
    project.setProperty (ID::hasCMakeLists,     false,                                  nullptr);
    project.setProperty (ID::selectedGenerator, toString (Generator::xcode),            nullptr);
    project.setProperty (ID::configuration,     toString (Configuration::debug),        nullptr);
    tree.addChild (project, -1, nullptr);

    // BUILDS node (children added dynamically by GeneratorDetector / auto-scan)
    tree.addChild (juce::ValueTree { ID::BUILDS }, -1, nullptr);

    // ASYNC node
    juce::ValueTree async { ID::ASYNC };
    async.setProperty (ID::isActive,   false,                          nullptr);
    async.setProperty (ID::isAborted,  false,                          nullptr);
    async.setProperty (ID::currentOp,  toString (OpType::none),        nullptr);
    tree.addChild (async, -1, nullptr);

    // CONFIG node
    juce::ValueTree config { ID::CONFIG };
    config.setProperty (ID::autoScanEnabled,  true,           nullptr);
    config.setProperty (ID::autoScanInterval, 10,             nullptr);
    config.setProperty (ID::theme,            juce::String{}, nullptr);
    tree.addChild (config, -1, nullptr);

    // THEME node (populated by ThemeLoader in Step 7)
    tree.addChild (juce::ValueTree { ID::THEME }, -1, nullptr);

    // Console owns its own line buffer (Step 6). No CONSOLE node in VT.
}

State::~State()
{
    stop();
}

// ============================================================================
// Lifecycle
// ============================================================================

void State::start() noexcept
{
    startTimer (FLUSH_INTERVAL_MS);
}

void State::stop() noexcept
{
    stopTimer();
}

// ============================================================================
// ValueTree subtree accessors
// ============================================================================

juce::ValueTree State::getProjectState() const noexcept
{
    return tree.getChildWithName (ID::PROJECT);
}

juce::ValueTree State::getBuildsState() const noexcept
{
    return tree.getChildWithName (ID::BUILDS);
}

juce::ValueTree State::getAsyncState() const noexcept
{
    return tree.getChildWithName (ID::ASYNC);
}

juce::ValueTree State::getConfigState() const noexcept
{
    return tree.getChildWithName (ID::CONFIG);
}

juce::ValueTree State::getThemeState() const noexcept
{
    return tree.getChildWithName (ID::THEME);
}

bool State::flush() noexcept
{
    jassert (juce::MessageManager::getInstance()->isThisTheMessageThread());
    needsFlush.store (true, std::memory_order_release);
    return applyDirtyAtoms();
}

// ============================================================================
// Atom setters
// ============================================================================

void State::setSelectedGenerator (Generator v) noexcept
{
    selectedGeneratorAtom.store (static_cast<int> (v), std::memory_order_relaxed);
    needsFlush.store (true, std::memory_order_release);
}

void State::setConfiguration (Configuration v) noexcept
{
    configurationAtom.store (static_cast<int> (v), std::memory_order_relaxed);
    needsFlush.store (true, std::memory_order_release);
}

void State::setCurrentOp (OpType v) noexcept
{
    currentOpAtom.store (static_cast<int> (v), std::memory_order_relaxed);
    needsFlush.store (true, std::memory_order_release);
}

void State::setHasCMakeLists (bool v) noexcept
{
    hasCMakeListsAtom.store (v, std::memory_order_relaxed);
    needsFlush.store (true, std::memory_order_release);
}

void State::setIsActive (bool v) noexcept
{
    isActiveAtom.store (v, std::memory_order_relaxed);
    needsFlush.store (true, std::memory_order_release);
}

void State::setIsAborted (bool v) noexcept
{
    isAbortedAtom.store (v, std::memory_order_relaxed);
    needsFlush.store (true, std::memory_order_release);
}

void State::setAutoScanEnabled (bool v) noexcept
{
    autoScanEnabledAtom.store (v, std::memory_order_relaxed);
    needsFlush.store (true, std::memory_order_release);
}

void State::setAutoScanInterval (int v) noexcept
{
    autoScanIntervalAtom.store (v, std::memory_order_relaxed);
    needsFlush.store (true, std::memory_order_release);
}

// ============================================================================
// Atom getters
// ============================================================================

Generator State::getSelectedGenerator() const noexcept
{
    return static_cast<Generator> (selectedGeneratorAtom.load (std::memory_order_relaxed));
}

Configuration State::getConfiguration() const noexcept
{
    return static_cast<Configuration> (configurationAtom.load (std::memory_order_relaxed));
}

OpType State::getCurrentOp() const noexcept
{
    return static_cast<OpType> (currentOpAtom.load (std::memory_order_relaxed));
}

bool State::getHasCMakeLists() const noexcept
{
    return hasCMakeListsAtom.load (std::memory_order_relaxed);
}

bool State::getIsActive() const noexcept
{
    return isActiveAtom.load (std::memory_order_relaxed);
}

bool State::getIsAborted() const noexcept
{
    return isAbortedAtom.load (std::memory_order_relaxed);
}

bool State::getAutoScanEnabled() const noexcept
{
    return autoScanEnabledAtom.load (std::memory_order_relaxed);
}

int State::getAutoScanInterval() const noexcept
{
    return autoScanIntervalAtom.load (std::memory_order_relaxed);
}

// ============================================================================
// Message-thread-only setters
// ============================================================================

void State::setWorkingDirectory (const juce::String& v) noexcept
{
    jassert (juce::MessageManager::getInstance()->isThisTheMessageThread());
    workingDirectory = v;
    needsFlush.store (true, std::memory_order_release);
}

void State::setTheme (const juce::String& v) noexcept
{
    jassert (juce::MessageManager::getInstance()->isThisTheMessageThread());
    themeName = v;
    needsFlush.store (true, std::memory_order_release);
}

// ============================================================================
// Timer
// ============================================================================

void State::timerCallback()
{
    jassert (juce::MessageManager::getInstance()->isThisTheMessageThread());
    applyDirtyAtoms();
}

bool State::applyDirtyAtoms() noexcept
{
    bool changed { false };

    if (needsFlush.exchange (false, std::memory_order_acquire))
    {
        const bool projectChanged { flushProject() };
        const bool asyncChanged   { flushAsync() };
        const bool configChanged  { flushConfig() };
        const bool stringsChanged { flushStrings() };
        changed = projectChanged or asyncChanged or configChanged or stringsChanged;
    }

    return changed;
}
