package ui

import (
	"sync"
	"time"
)

// CharAnimationState represents the animation state of a single character
type CharAnimationState int

const (
	CharStateStable CharAnimationState = iota
	CharStateBlinking
	CharStateComplete
)

// CharAnimation tracks the animation state for a character position
type CharAnimation struct {
	OldChar      rune
	NewChar      rune
	State        CharAnimationState
	BlinkPhase   int // 0 or 1 for blinking
	StartTime    time.Time
	mu           sync.Mutex
}

// AnimatedText manages character-by-character animations for text
type AnimatedText struct {
	OldText      string
	NewText      string
	Chars        []*CharAnimation
	MaxLength    int
	mu           sync.Mutex
}

// NewAnimatedText creates a new animated text with the given max length
func NewAnimatedText(maxLength int) *AnimatedText {
	return &AnimatedText{
		Chars:     make([]*CharAnimation, maxLength),
		MaxLength: maxLength,
	}
}

// Update sets new text and initiates animations for changed characters
func (at *AnimatedText) Update(newText string) {
	at.mu.Lock()
	defer at.mu.Unlock()

	// Pad or truncate to max length
	if len(newText) > at.MaxLength {
		newText = newText[:at.MaxLength]
	} else {
		for len(newText) < at.MaxLength {
			newText += " "
		}
	}

	oldText := at.OldText
	if len(oldText) > at.MaxLength {
		oldText = oldText[:at.MaxLength]
	} else {
		for len(oldText) < at.MaxLength {
			oldText += " "
		}
	}

	at.NewText = newText

	// Initialize or update character animations
	for i := 0; i < at.MaxLength; i++ {
		var oldChar, newChar rune
		if i < len(oldText) {
			oldChar = rune(oldText[i])
		} else {
			oldChar = ' '
		}
		if i < len(newText) {
			newChar = rune(newText[i])
		} else {
			newChar = ' '
		}

		if at.Chars[i] == nil {
			at.Chars[i] = &CharAnimation{
				OldChar: oldChar,
				NewChar: newChar,
				State:   CharStateStable,
			}
		}

		// If character changed, start animation
		if oldChar != newChar {
			at.Chars[i].OldChar = oldChar
			at.Chars[i].NewChar = newChar
			at.Chars[i].State = CharStateBlinking
			at.Chars[i].BlinkPhase = 0
			at.Chars[i].StartTime = time.Now()
		}
	}

	at.OldText = newText
}

// Tick updates animation states (call this periodically)
func (at *AnimatedText) Tick() {
	at.mu.Lock()
	defer at.mu.Unlock()

	for _, char := range at.Chars {
		if char == nil {
			continue
		}

		char.mu.Lock()
		if char.State == CharStateBlinking {
			// Toggle blink phase every ~100ms
			elapsed := time.Since(char.StartTime)
			char.BlinkPhase = int(elapsed.Milliseconds() / 100) % 2

			// After ~300ms, complete the animation
			if elapsed > 300*time.Millisecond {
				char.State = CharStateComplete
			}
		} else if char.State == CharStateComplete {
			// Mark as stable after completion
			char.State = CharStateStable
		}
		char.mu.Unlock()
	}
}

// Render returns the current display string with animations applied
func (at *AnimatedText) Render() string {
	at.mu.Lock()
	defer at.mu.Unlock()

	result := make([]rune, at.MaxLength)
	for i, char := range at.Chars {
		if char == nil {
			result[i] = ' '
			continue
		}

		char.mu.Lock()
		switch char.State {
		case CharStateBlinking:
			// Show blinking cursor during animation
			if char.BlinkPhase == 0 {
				result[i] = '█' // Solid block
			} else {
				result[i] = '░' // Light block
			}
		case CharStateComplete, CharStateStable:
			result[i] = char.NewChar
		default:
			result[i] = char.NewChar
		}
		char.mu.Unlock()
	}

	return string(result)
}

// IsAnimating returns true if any character is currently animating
func (at *AnimatedText) IsAnimating() bool {
	at.mu.Lock()
	defer at.mu.Unlock()

	for _, char := range at.Chars {
		if char != nil {
			char.mu.Lock()
			if char.State == CharStateBlinking {
				char.mu.Unlock()
				return true
			}
			char.mu.Unlock()
		}
	}
	return false
}

