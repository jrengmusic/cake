#pragma once
#include <JuceHeader.h>
#include <atomic>
#include "Axis.h"
#include "../Identifier.h"

// ============================================================================
// State
// ============================================================================
//
// APVTS-style atomic state bridge for CAKE-cpp.  Mirrors TitState structure
// verbatim, adapted for CAKE's domain (project, builds, async ops, config).
//
// Thread ownership:
//   - Worker threads (cmake subprocess) write exclusively to atoms via set*()
//     setters.  No ValueTree mutations, no allocations, no locks.
//   - The flush Timer (message thread) reads atoms -> writes ValueTree on
//     each tick, firing ValueTree::Listeners.
//   - Message thread reads ValueTree for UI consumption; also owns non-atom
//     fields (workingDirectory, theme) with jassert guards on setters.
//   - Views attach juce::ValueTree::Listener to relevant subtrees.
//
// Zero locks on the hot path.  Zero shadow state.
// Atoms ARE the state; ValueTree REFLECTS it.
//
// SPEC §Architecture; BLESSED B (thread bounds), S (SSOT).

class State : public juce::Timer
{
public:

    // =========================================================================
    // Construction / Destruction
    // =========================================================================

    // Constructs default ValueTree per SPEC §ValueTree Schema.
    // All nodes present; properties initialised to defaults.
    // MESSAGE THREAD — constructed before Timer is started via start().
    State();

    // Stops the Timer and destroys State.
    // MESSAGE THREAD — must be destroyed on the message thread.
    ~State() override;

    // =========================================================================
    // Lifecycle
    // =========================================================================

    // Begins the flush Timer at FLUSH_INTERVAL_MS cadence.
    // MESSAGE THREAD.
    void start() noexcept;

    // Stops the flush Timer.
    // MESSAGE THREAD.
    void stop() noexcept;

    // =========================================================================
    // ValueTree subtree accessors (message-thread observation surface)
    // =========================================================================

    // Each returns the named child of the root CAKE tree.
    // Views attach juce::ValueTree::Listener to the subtree they need.
    // MESSAGE THREAD.
    juce::ValueTree getProjectState() const noexcept;
    juce::ValueTree getBuildsState()  const noexcept;
    juce::ValueTree getAsyncState()   const noexcept;
    juce::ValueTree getConfigState()  const noexcept;
    juce::ValueTree getThemeState()   const noexcept;

    // Forces an immediate flush without waiting for the next Timer tick.
    // Used by tests to drive synchronous flush on the message thread.
    // Returns true if any property was updated.
    // MESSAGE THREAD.
    bool flush() noexcept;

    // =========================================================================
    // Atom setters — worker-thread fast writes (any thread, lock-free)
    // =========================================================================

    void setSelectedGenerator (Generator     v) noexcept;
    void setConfiguration     (Configuration v) noexcept;
    void setCurrentOp         (OpType        v) noexcept;
    void setHasCMakeLists     (bool          v) noexcept;
    void setIsActive          (bool          v) noexcept;
    void setIsAborted         (bool          v) noexcept;
    void setAutoScanEnabled   (bool          v) noexcept;
    void setAutoScanInterval  (int           v) noexcept;

    // =========================================================================
    // Atom getters — any thread, lock-free
    // =========================================================================

    Generator     getSelectedGenerator() const noexcept;
    Configuration getConfiguration()     const noexcept;
    OpType        getCurrentOp()         const noexcept;
    bool          getHasCMakeLists()     const noexcept;
    bool          getIsActive()          const noexcept;
    bool          getIsAborted()         const noexcept;
    bool          getAutoScanEnabled()   const noexcept;
    int           getAutoScanInterval()  const noexcept;

    // =========================================================================
    // Message-thread-only setters for non-atom string fields
    // =========================================================================

    // Sets the working directory and marks flush dirty.
    // MESSAGE THREAD — jassert guards this boundary.
    void setWorkingDirectory (const juce::String& v) noexcept;

    // Sets the theme name and marks flush dirty.
    // MESSAGE THREAD — jassert guards this boundary.
    void setTheme (const juce::String& v) noexcept;

private:

    // =========================================================================
    // Timer
    // =========================================================================

    // Flush cadence: 16 ms ~= 60 Hz.  Named constant per BLESSED-E (no magic numbers).
    static constexpr int FLUSH_INTERVAL_MS { 16 };

    // timerCallback runs on the message thread; delegates to applyDirtyAtoms().
    void timerCallback() override;

    // Copies dirty atoms -> ValueTree properties in one pass.
    // Skip-unchanged optimisation: compares against lastFlushed snapshot.
    // Returns true if any property was updated.
    // MESSAGE THREAD.
    bool applyDirtyAtoms() noexcept;

    // Per-subtree flush helpers — called exclusively from applyDirtyAtoms().
    // Each returns true if any property in its subtree was written.
    // MESSAGE THREAD.
    bool flushProject() noexcept;
    bool flushAsync()   noexcept;
    bool flushConfig()  noexcept;
    bool flushStrings() noexcept;

    // Compare-write-snapshot helper.  If snapshotSlot != current, writes
    // asVar to subtree[key], updates snapshotSlot, and returns true.
    // Called only from flush helpers (message thread).
    template <typename T>
    bool writeIfChanged (T&                     snapshotSlot,
                         T                      current,
                         juce::ValueTree&        subtree,
                         const juce::Identifier& key,
                         const juce::var&        asVar) noexcept;

    // =========================================================================
    // ValueTree — message-thread SSOT (observation surface)
    // =========================================================================

    juce::ValueTree tree;

    // =========================================================================
    // Atoms — cross-thread fast-write primitives
    // =========================================================================

    std::atomic<int>  selectedGeneratorAtom { static_cast<int> (Generator::xcode) };
    std::atomic<int>  configurationAtom     { static_cast<int> (Configuration::debug) };
    std::atomic<int>  currentOpAtom         { static_cast<int> (OpType::none) };
    std::atomic<bool> hasCMakeListsAtom     { false };
    std::atomic<bool> isActiveAtom          { false };
    std::atomic<bool> isAbortedAtom         { false };
    std::atomic<bool> autoScanEnabledAtom   { true };
    std::atomic<int>  autoScanIntervalAtom  { 10 };

    // needsFlush is set by any setter; read+cleared in timerCallback.
    std::atomic<bool> needsFlush { false };

    // =========================================================================
    // Non-atom message-thread fields
    // =========================================================================

    juce::String workingDirectory;
    juce::String themeName;

    // =========================================================================
    // Last-flushed snapshot (message-thread owned, skip-unchanged optimisation)
    // =========================================================================

    struct FlushSnapshot
    {
        int  selectedGenerator { static_cast<int> (Generator::xcode) };
        int  configuration     { static_cast<int> (Configuration::debug) };
        int  currentOp         { static_cast<int> (OpType::none) };
        bool hasCMakeLists     { false };
        bool isActive          { false };
        bool isAborted         { false };
        bool autoScanEnabled   { true };
        int  autoScanInterval  { 10 };
        juce::String workingDirectory;
        juce::String themeName;
    };

    FlushSnapshot lastFlushed;

    JUCE_DECLARE_NON_COPYABLE_WITH_LEAK_DETECTOR (State)
};
