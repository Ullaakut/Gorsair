// Package symbol contains a few cherry-picked UTF-8 symbols to be
// used to build user-friendly command-line interfaces.
// They are all colorless so that they can be used along
// with `disgo/logger`'s formatting helpers.
package symbol

const (
	// Check displays ✔
	Check = "\xe2\x9c\x94"

	// Cross displays ✖
	Cross = "\xe2\x9c\x96"

	// LeftArrow displays ❮
	LeftArrow = "\xe2\x9d\xae"

	// RightArrow displays ❯
	RightArrow = "\xe2\x9d\xaf"

	// LeftTriangle displays ◀
	LeftTriangle = "\xe2\x97\x80"

	// RightTriangle displays ▶
	RightTriangle = "\xe2\x96\xb6"
)
