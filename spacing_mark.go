package displaywidth

// spacingMarkWidth measures spacing components after the leading character.
// Grapheme terminals expand a narrow cell to two columns when a spacing mark
// joins it. Further marks in the same cluster do not add more columns. A mark
// standing alone remains one column, and nonspacing/enclosing marks do not
// expand the cell. Thai/Lao SARA AM are spacing components despite category Lo.
func spacingMarkWidth[T ~string | ~[]byte](s T, firstSize int) int {
	for pos := firstSize; pos < len(s); {
		p, size := lookup(s[pos:])
		if property(p).is(_Spacing_Mark) {
			return 2
		}
		if size <= 0 {
			// Width APIs accept invalid UTF-8; always make progress.
			size = 1
		}
		pos += size
	}
	return 1
}
