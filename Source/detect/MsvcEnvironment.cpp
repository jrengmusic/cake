#include "MsvcEnvironment.h"

// ============================================================================
// MsvcEnvironment — macOS stub
// ============================================================================
//
// Windows implementation deferred to post-MVP sprint per PLAN §Risks.
// All methods are no-ops on macOS and Linux.

bool MsvcEnvironment::captureEnvironment() noexcept
{
    return false;
}

bool MsvcEnvironment::isValid() const noexcept
{
    return false;
}
