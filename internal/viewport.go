package internal

// Viewport manages scrollable content display within a fixed-size area.
// It handles content that may be larger than the available display space,
// providing scrolling functionality and viewport management.
type Viewport struct {
	offset  int      // Current scroll offset (top line being displayed)
	width   int      // Width of the viewport in characters
	height  int      // Height of the viewport in lines
	content []string // All content lines available for display
}

// NewViewport creates a new Viewport with default dimensions.
// The viewport starts with no content and default size of 80x24.
func NewViewport() *Viewport {
	return &Viewport{
		offset:  0,
		width:   80,
		height:  24,
		content: make([]string, 0),
	}
}

// SetSize updates the viewport's display dimensions.
// This affects how much content can be visible at once.
func (v *Viewport) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// SetContent replaces all content in the viewport with the provided lines.
// The scroll offset is reset to 0 (top of content).
func (v *Viewport) SetContent(content []string) {
	v.content = content
	v.offset = 0
}

// Scroll moves the viewport by the specified number of lines.
// Positive delta scrolls down, negative delta scrolls up.
// The viewport will not scroll beyond the content boundaries.
func (v *Viewport) Scroll(delta int) {
	maxOffset := max(0, len(v.content)-v.height)
	v.offset = max(0, min(maxOffset, v.offset+delta))
}

// VisibleContent returns the lines currently visible in the viewport.
// This is a subset of the total content based on the current scroll position.
func (v *Viewport) VisibleContent() []string {
	if len(v.content) == 0 {
		return []string{}
	}

	start := v.offset
	end := min(start+v.height, len(v.content))

	if start >= len(v.content) {
		return []string{}
	}

	return v.content[start:end]
}

// CanScrollUp returns true if the viewport can scroll up (show earlier content).
// This is true when the current offset is greater than 0.
func (v *Viewport) CanScrollUp() bool {
	return v.offset > 0
}

// CanScrollDown returns true if the viewport can scroll down (show later content).
// This is true when there is more content below the current visible area.
func (v *Viewport) CanScrollDown() bool {
	return v.offset < len(v.content)-v.height
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
