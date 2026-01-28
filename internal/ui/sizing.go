package ui

const (
	MinWidth         = 69
	MinHeight        = 19
	HeaderHeight     = 3
	FooterHeight     = 1
	MinContentHeight = 4
	HorizontalMargin = 2
)

type DynamicSizing struct {
	TerminalWidth     int
	TerminalHeight    int
	ContentHeight     int
	ContentInnerWidth int
	HeaderInnerWidth  int
	FooterInnerWidth  int
	IsTooSmall        bool
}

func CalculateDynamicSizing(termWidth, termHeight int) DynamicSizing {
	isTooSmall := termWidth < MinWidth || termHeight < MinHeight

	contentHeight := termHeight - HeaderHeight - FooterHeight
	if contentHeight < MinContentHeight {
		contentHeight = MinContentHeight
	}

	innerWidth := termWidth - (HorizontalMargin * 2)
	if innerWidth < 20 {
		innerWidth = 20
	}

	return DynamicSizing{
		TerminalWidth:     termWidth,
		TerminalHeight:    termHeight,
		ContentHeight:     contentHeight,
		ContentInnerWidth: innerWidth,
		HeaderInnerWidth:  innerWidth,
		FooterInnerWidth:  termWidth,
		IsTooSmall:        isTooSmall,
	}
}

func NewDynamicSizing() DynamicSizing {
	return CalculateDynamicSizing(80, 24)
}
