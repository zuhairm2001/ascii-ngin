package ascii

const ASCII_CHARS = " .'\\`^,:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
const LUMINANCE_THRESHOLD = 30
const FONT_RATIO = 0.44
const LUMINANCE_RANGE = 255.0

func CalculateLuminance(r, g, b int) float64 {
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}

func MapLuminanceToASCII(luminance float64) rune {
	index := int((luminance / LUMINANCE_RANGE) * float64(len(ASCII_CHARS)-1))
	return rune(ASCII_CHARS[index])
}

func RecalculateHeight(height int) int {
	return int(float64(height) * FONT_RATIO)
}
