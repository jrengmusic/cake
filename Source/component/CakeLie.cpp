#include <JuceHeader.h>
#include "CakeLie.h"

// ============================================================================
// Construction
// ============================================================================

CakeLie::CakeLie()
{
    // Load SVG from BinaryData — embedded at build time via CMakeLists BINARY_FILES.
    // BinaryData::cakelie_svg is the mangled name for cake-lie.svg.
    svgDocument = juce::XmlDocument::parse (
        juce::String::fromUTF8 (BinaryData::cakelie_svg, BinaryData::cakelie_svgSize)
    );

    jassert (svgDocument != nullptr);
}

// ============================================================================
// Component overrides
// ============================================================================

void CakeLie::paint (jam::tui::Graphics& g)
{
    const jam::tui::Rectangle bounds { getCellBounds() };

    if (not bounds.isEmpty() and brailleGrid.cols > 0 and brailleGrid.rows > 0)
    {
        // Center the braille grid within bounds.
        const int drawCols  { juce::jmin (brailleGrid.cols, bounds.width) };
        const int drawRows  { juce::jmin (brailleGrid.rows, bounds.height) };
        const int offsetCol { (bounds.width  - drawCols) / 2 };
        const int offsetRow { (bounds.height - drawRows) / 2 };

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
                        offsetCol + col, offsetRow + row, 1
                    );
                }
            }
        }
    }
}

void CakeLie::resized()
{
    rasterizeToGrid();
}

// ============================================================================
// Private helpers
// ============================================================================

void CakeLie::rasterizeToGrid() noexcept
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
