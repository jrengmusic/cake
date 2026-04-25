#include "Axis.h"

// ============================================================================
// toString
// ============================================================================

juce::String toString (Mode v) noexcept
{
    if (v == Mode::menu)            return "menu";
    if (v == Mode::preferences)     return "preferences";
    if (v == Mode::console)         return "console";
    return "invalidProject";
}

juce::String toString (Generator v) noexcept
{
    if (v == Generator::xcode)   return "Xcode";
    if (v == Generator::ninja)   return "Ninja";
    if (v == Generator::vs2026)  return "VS2026";
    return "VS2022";
}

juce::String toString (Configuration v) noexcept
{
    if (v == Configuration::debug)  return "Debug";
    return "Release";
}

juce::String toString (OpType v) noexcept
{
    if (v == OpType::none)          return "none";
    if (v == OpType::build)         return "build";
    if (v == OpType::generate)      return "generate";
    if (v == OpType::clean)         return "clean";
    if (v == OpType::cleanAll)      return "cleanAll";
    return "regenerate";
}

// ============================================================================
// parse*
// ============================================================================

Generator parseGenerator (const juce::String& s) noexcept
{
    if (s == "Xcode")   return Generator::xcode;
    if (s == "Ninja")   return Generator::ninja;
    if (s == "VS2026")  return Generator::vs2026;
    return Generator::vs2022;
}

Configuration parseConfiguration (const juce::String& s) noexcept
{
    if (s == "Debug")   return Configuration::debug;
    return Configuration::release;
}

OpType parseOpType (const juce::String& s) noexcept
{
    if (s == "none")        return OpType::none;
    if (s == "build")       return OpType::build;
    if (s == "generate")    return OpType::generate;
    if (s == "clean")       return OpType::clean;
    if (s == "cleanAll")    return OpType::cleanAll;
    return OpType::regenerate;
}
