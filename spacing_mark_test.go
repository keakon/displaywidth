package displaywidth

import (
	"testing"
	"unicode"
)

func TestSpacingMarkWidths(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		want       int
	}{
		{"ascii", "plain text", 10},
		{"devanagari", "\u092a\u0949", 2},
		{"gurmukhi_on_ascii", "a\u0a3e", 2},
		{"multiple_spacing_marks", "a\u093e\u093e", 2},
		{"standalone_spacing_mark", "\u093e", 1},
		{"wide_base", "\u4e00\u093e", 2},
		{"nonspacing", "e\u0301", 1},
		{"enclosing", "a\u20dd", 1},
		{"thai", "\u0e19\u0e33", 2},
		{"thai_tone", "\u0e19\u0e49\u0e33", 2},
		{"lao", "\u0e99\u0eb3", 2},
		{"standalone_thai", "\u0e33", 1},
		{"hangul", "\u1100\u1161\u11a8", 2},
		{"emoji_modifier", "\U0001f44b\U0001f3ff", 2},
		{"text_emoji_modifier", "\U0001f575\U0001f3ff", 1},
		{"emoji_zwj", "\U0001f469\u200d\U0001f4bb", 2},
		{"flag", "\U0001f1ef\U0001f1f5", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := String(tc.text); got != tc.want {
				t.Fatalf("String(%q) = %d, want %d", tc.text, got, tc.want)
			}
			if got := Bytes([]byte(tc.text)); got != tc.want {
				t.Fatalf("Bytes(%q) = %d, want %d", tc.text, got, tc.want)
			}
			var total int
			g := StringGraphemes(tc.text)
			for g.Next() {
				total += g.Width()
			}
			if total != tc.want {
				t.Fatalf("iterator width = %d, want %d", total, tc.want)
			}
			total = 0
			gb := BytesGraphemes([]byte(tc.text))
			for gb.Next() {
				total += gb.Width()
			}
			if total != tc.want {
				t.Fatalf("byte iterator width = %d, want %d", total, tc.want)
			}
		})
	}
}

func TestSpacingMarkCategory(t *testing.T) {
	for r := rune(0); r <= unicode.MaxRune; r++ {
		if !unicode.Is(unicode.Mc, r) {
			continue
		}
		text := "a" + string(r)
		if got := String(text); got != 2 {
			t.Errorf("U+%04X spacing mark: width of %q = %d, want 2", r, text, got)
		}
	}
}

func TestSpacingMarkControlSequences(t *testing.T) {
	options := Options{ControlSequences: true}
	const text = "\x1b[31ma\u093e\x1b[0mx"
	if got := options.String(text); got != 3 {
		t.Fatalf("styled width = %d, want 3", got)
	}
	if got := options.TruncateString(text, 2, ""); got != "\x1b[31ma\u093e\x1b[0m" {
		t.Fatalf("styled truncation = %q", got)
	}
}

func TestSpacingMarkTruncation(t *testing.T) {
	const text = "a\u093eb"
	if got := TruncateString(text, 1, ""); got != "" {
		t.Fatalf("one-column truncation = %q", got)
	}
	if got := TruncateString(text, 2, ""); got != "a\u093e" {
		t.Fatalf("two-column truncation = %q", got)
	}
	if got := string(TruncateBytes([]byte(text), 2, nil)); got != "a\u093e" {
		t.Fatalf("byte truncation = %q", got)
	}
	if got := Rune('\u093e'); got != 1 {
		t.Fatalf("standalone spacing rune = %d", got)
	}
}
