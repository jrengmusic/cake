# CMakeLists.txt - CAKE-cpp: CMake Project Manager TUI
# ============================================================================
# Context-aware cmake operations (generate/build/clean/open IDE), C++17/JUCE8
# ============================================================================

cmake_minimum_required(VERSION 4.2.0)

# ============================================================================
# JAM BUILD SYSTEM
# ============================================================================

set(JAM_ROOT "$ENV{HOME}/Documents/Poems/dev/jam")
list(APPEND CMAKE_MODULE_PATH "${JAM_ROOT}/cmake")
include(BuildSetup)

# ============================================================================
# PROJECT
# ============================================================================

if(APPLE)
    set(_EXTRA_LANGS OBJC OBJCXX)
else()
    set(_EXTRA_LANGS "")
endif()

project(cake-cpp
    VERSION 0.0.0
    LANGUAGES C CXX ${_EXTRA_LANGS}
    DESCRIPTION "CAKE-cpp: CMake Project Manager TUI"
)

# ============================================================================
# BUILD
# ============================================================================

configure_app(
    TARGET_NAME      cakec
    PRODUCT_NAME     cakec
    VERSION          "${PROJECT_VERSION}"
    COMPANY          jrengmusic
    COMPANY_WEBSITE  "https://jrengmusic.com"
    BUNDLE_ID        "com.jrengmusic.cakec"
    MODULES
        juce_core
        juce_events
        juce_data_structures
        juce_graphics
        juce_gui_basics
    JAM_MODULES
        jam_core
        jam_fonts
        jam_data_structures
        jam_markdown
        jam_tui
        jam_subprocess
    BINARY_FILES
        "${CMAKE_CURRENT_SOURCE_DIR}/Assets/cake-logo.svg"
        "${CMAKE_CURRENT_SOURCE_DIR}/Assets/cake-lie.svg"
    BINARY_NAMESPACE "BinaryData"
    EXTRA_DEFINES
        JUCE_WEB_BROWSER=0
        JUCE_USE_CURL=0
)

message(STATUS "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
message(STATUS "Configuration complete!")
message(STATUS "Generator: ${CMAKE_GENERATOR}")
message(STATUS "Build directory: ${CMAKE_BINARY_DIR}")
message(STATUS "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
