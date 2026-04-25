#include <JuceHeader.h>
#include "Banner.h"

// ============================================================================
// Construction
// ============================================================================

Banner::Banner()
{
    // Load SVG from BinaryData — embedded at build time via CMakeLists BINARY_FILES.
    // BinaryData::cakelogo_svg is the mangled name for cake-logo.svg.
    svgDocument = juce::XmlDocument::parse (
        juce::String::fromUTF8 (BinaryData::cakelogo_svg, BinaryData::cakelogo_svgSize)
    );

    jassert (svgDocument != nullptr);
}

// ============================================================================
// Component overrides
// ============================================================================

void Banner::paint (jam::tui::Graphics& g)
{
    const jam::tui::Rectangle bounds { getCellBounds() };

    if (not bounds.isEmpty() and brailleGrid.cols > 0 and brailleGrid.rows > 0)
    {
        const int drawCols { juce::jmin (brailleGrid.cols, bounds.width) };
        const int drawRows { juce::jmin (brailleGrid.rows, bounds.height) };

        for (int row { 0 }; row < drawRows; ++row)
        {
            for (int col { 0 }; col < drawCols; ++col)
            {
                const jam::tui::Cell& cell { brailleGrid.cells.at (
                    static_cast<std::size_t> (row * brailleGrid.cols + col)
                ) };

                if (cell.codepoint != 0)
                {
                    const juce::Colour fg {
                        cell.fg.red, cell.fg.green, cell.fg.blue
                    };

                    g.setColour (fg);
                    g.drawCellText (
                        juce::String::charToString (static_cast<juce::juce_wchar> (cell.codepoint)),
                        col, row, 1
                    );
                }
            }
        }
    }
}

void Banner::resized()
{
    rasterizeToGrid();
}

// ============================================================================
// Private helpers
// ============================================================================

void Banner::rasterizeToGrid() noexcept
{
    const jam::tui::Rectangle bounds { getCellBounds() };

    if (svgDocument != nullptr and not bounds.isEmpty())
    {
        const int pixelWidth  { bounds.width  * jam::braille::BRAILLE_CELL_WIDTH  };
        const int pixelHeight { bounds.height * jam::braille::BRAILLE_CELL_HEIGHT };

        brailleGrid = jam::braille::renderSvgToBrailleGrid (
            *svgDocument,
            pixelWidth,
            pixelHeight
        );
    }
}
