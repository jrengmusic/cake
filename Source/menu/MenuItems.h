#pragma once

// ============================================================================
// MenuItems — compile-time row ID constants for CAKE's 8 fixed menu rows
// ============================================================================
//
// Used by MenuBuilder (row construction) and MainComponent (action dispatch).
// All 8 rows are always present — selectability varies per VT state.
//
// SPEC §Menu Item Selectability Rules; BLESSED-S (SSOT for row IDs).

namespace menu
{
    static constexpr const char* ROW_PROJECT       { "project" };
    static constexpr const char* ROW_GENERATE      { "generate" };
    static constexpr const char* ROW_OPEN          { "open" };
    static constexpr const char* ROW_SEPARATOR     { "separator" };
    static constexpr const char* ROW_CONFIGURATION { "configuration" };
    static constexpr const char* ROW_BUILD         { "build" };
    static constexpr const char* ROW_CLEAN         { "clean" };
    static constexpr const char* ROW_CLEAN_ALL     { "cleanAll" };

} // namespace menu
