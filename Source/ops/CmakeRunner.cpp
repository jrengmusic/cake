#include "CmakeRunner.h"

// Pull Console header here — CmakeRunner.h forward-declared it to avoid
// pulling jam_tui into every file that includes CmakeRunner.h.
#include <jam_tui/component/jam_tui_console.h>

#include <unistd.h>

// ============================================================================
// Construction / Destruction
// ============================================================================

CmakeRunner::CmakeRunner (State& state, jam::tui::Console& console)
    : stateRef   { state }
    , consoleRef { console }
{
}

CmakeRunner::~CmakeRunner()
{
    subprocess.kill();
}

// ============================================================================
// Operations
// ============================================================================

void CmakeRunner::generate()
{
    jassert (juce::MessageManager::getInstance()->isThisTheMessageThread());

    const Generator gen      { stateRef.getSelectedGenerator() };
    const juce::String dir   { buildDirName (gen) };
    const juce::String genArg { cmakeGeneratorArg (gen) };
    const juce::String workDir { stateRef.getProjectState()
                                          .getProperty (ID::workingDirectory).toString() };

    jassert (workDir.isNotEmpty());

    consoleRef.clear();

    const juce::String command {
        juce::String ("cmake -G ") + genArg
        + " -S . -B " + dir
    };

    appendOpHeader ("GENERATING", command);

    stateRef.setCurrentOp (OpType::generate);
    stateRef.setIsActive  (true);
    stateRef.setIsAborted (false);

    const juce::StringArray argv {
        "cmake",
        "-G",   genArg,
        "-S",   ".",
        "-B",   dir
    };

    subprocess.launch (
        argv,
        juce::File (workDir),
        [this] (int exitCode,
                const juce::String& /*stdoutCapture*/,
                const juce::String& /*stderrCapture*/)
        {
            onOperationComplete (exitCode);
        },
        [this] (juce::StringRef chunk, bool /*isReplace*/)
        {
            consoleRef.appendLine (chunk.text);
        }
    );
}

void CmakeRunner::build()
{
    jassert (juce::MessageManager::getInstance()->isThisTheMessageThread());

    const Generator     gen    { stateRef.getSelectedGenerator() };
    const Configuration cfg    { stateRef.getConfiguration() };
    const juce::String  dir    { buildDirName (gen) };
    const juce::String  cfgStr { toString (cfg) };
    const juce::String  workDir { stateRef.getProjectState()
                                           .getProperty (ID::workingDirectory).toString() };

    jassert (workDir.isNotEmpty());

    consoleRef.clear();

    const juce::String command {
        juce::String ("cmake --build ") + dir + " --config " + cfgStr
    };

    appendOpHeader ("BUILDING", command);

    stateRef.setCurrentOp (OpType::build);
    stateRef.setIsActive  (true);
    stateRef.setIsAborted (false);

    const juce::StringArray argv {
        "cmake",
        "--build", dir,
        "--config", cfgStr
    };

    subprocess.launch (
        argv,
        juce::File (workDir),
        [this] (int exitCode,
                const juce::String& /*stdoutCapture*/,
                const juce::String& /*stderrCapture*/)
        {
            onOperationComplete (exitCode);
        },
        [this] (juce::StringRef chunk, bool /*isReplace*/)
        {
            consoleRef.appendLine (chunk.text);
        }
    );
}

void CmakeRunner::clean()
{
    jassert (juce::MessageManager::getInstance()->isThisTheMessageThread());

    const Generator gen   { stateRef.getSelectedGenerator() };
    const juce::File dir  { buildDir (gen) };

    consoleRef.clear();
    consoleRef.appendLine ("Cleaning: " + dir.getFullPathName());

    stateRef.setCurrentOp (OpType::clean);
    stateRef.setIsActive  (true);
    stateRef.setIsAborted (false);

    if (not dir.exists())
    {
        consoleRef.appendLine ("Project directory clean.");
        consoleRef.appendLine ("Press ESC to return to menu.");
    }
    else
    {
        if (dir.deleteRecursively())
        {
            consoleRef.appendLine ("ok");
            consoleRef.appendLine ("Project directory clean.");
            consoleRef.appendLine ("Press ESC to return to menu.");
        }
        else
        {
            consoleRef.appendLine ("ERROR: Failed to remove directory: " + dir.getFullPathName());
        }
    }

    // Clean is synchronous — no subprocess. Clear ASYNC immediately.
    stateRef.setIsActive  (false);
    stateRef.setCurrentOp (OpType::none);
}

void CmakeRunner::open()
{
    jassert (juce::MessageManager::getInstance()->isThisTheMessageThread());

    const Generator gen     { stateRef.getSelectedGenerator() };
    const juce::File dir    { buildDir (gen) };

    consoleRef.clear();

    if (gen == Generator::xcode)
    {
        // Find *.xcodeproj directory inside the build dir.
        juce::File projectFile;
        bool found { false };

        for (const auto& entry : dir.findChildFiles (juce::File::findDirectories, false, "*.xcodeproj"))
        {
            projectFile = entry;
            found = true;
            break;
        }

        if (not found)
        {
            consoleRef.appendLine ("Failed to open IDE: Xcode project file not found in " + dir.getFullPathName());
        }
        else
        {
            consoleRef.appendLine ("Opening Xcode project: " + projectFile.getFullPathName());

            const juce::StringArray argv { "open", projectFile.getFullPathName() };

            stateRef.setCurrentOp (OpType::none);
            stateRef.setIsActive  (false);

            subprocess.launch (
                argv,
                dir,
                nullptr,
                nullptr
            );
        }
    }
    else if (gen == Generator::ninja)
    {
        // Replace the cakec process with nvim — exec() never returns on success.
        // The OS reclaims all resources; no TUI teardown needed.
        // On failure, fall through to quit so cakec does not hang.
        execlp ("nvim", "nvim", dir.getFullPathName().toRawUTF8(), nullptr);

        // exec failed — nvim not found or not executable. Exit cleanly.
        juce::JUCEApplication::getInstance()->systemRequestedQuit();
    }
    else
    {
        // VS generators deferred — macOS-first per PLAN §Risks.
        consoleRef.appendLine ("Open IDE: Visual Studio support is post-MVP.");
        stateRef.setCurrentOp (OpType::none);
        stateRef.setIsActive  (false);
    }
}

void CmakeRunner::abort()
{
    jassert (juce::MessageManager::getInstance()->isThisTheMessageThread());

    subprocess.kill();
    stateRef.setIsAborted (true);
    stateRef.setIsActive  (false);
    stateRef.setCurrentOp (OpType::none);
}

// ============================================================================
// Private helpers
// ============================================================================

juce::String CmakeRunner::buildDirName (Generator gen) const
{
    juce::String result { juce::String (BUILDS_ROOT) + juce::File::getSeparatorString() };

    if      (gen == Generator::xcode) result += BUILD_DIR_XCODE;
    else if (gen == Generator::ninja) result += BUILD_DIR_NINJA;
    else                              result += BUILD_DIR_NINJA; // fallback — deferred generators

    return result;
}

juce::String CmakeRunner::cmakeGeneratorArg (Generator gen) const
{
    juce::String result { CMAKE_G_NINJA };

    if (gen == Generator::xcode)
        result = juce::String { CMAKE_G_XCODE };

    return result;
}

juce::File CmakeRunner::buildDir (Generator gen) const
{
    const juce::String workDir { stateRef.getProjectState()
                                          .getProperty (ID::workingDirectory).toString() };
    return juce::File (workDir).getChildFile (buildDirName (gen));
}

void CmakeRunner::appendOpHeader (const juce::String& label,
                                  const juce::String& command) const
{
    consoleRef.appendLine (label);
    consoleRef.appendLine ("Running: " + command);
    consoleRef.appendLine ({});
}

void CmakeRunner::onOperationComplete (int exitCode)
{
    // Fires on Worker thread — only atom writes permitted here.
    const bool succeeded { exitCode == 0 };

    juce::MessageManager::callAsync ([this, succeeded, exitCode]
    {
        if (succeeded)
        {
            consoleRef.appendLine ({});
            consoleRef.appendLine ("Completed successfully.");
            consoleRef.appendLine ("Press ESC to return to menu.");
        }
        else
        {
            consoleRef.appendLine ({});
            consoleRef.appendLine ("ERROR: Operation failed (exit code "
                                   + juce::String (exitCode) + ").");
            consoleRef.appendLine ("Press ESC to return to menu.");
        }

        stateRef.setIsActive  (false);
        stateRef.setCurrentOp (OpType::none);
    });
}
