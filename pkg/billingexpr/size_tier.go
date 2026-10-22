package billingexpr

import (
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// sizeTierPresets mirrors infinite-canvas image size presets: each tier
// targets a fixed pixel budget (1k ~1.05MP, 2k ~4.19MP, 4k ~8.29MP), so
// wide ratios carry fewer pixels than the square preset.
var sizeTierPresets = map[string]string{
	"1024x1024": "1k", "1024x1536": "1k", "1536x1024": "1k", "1024x768": "1k",
	"768x1024": "1k", "1536x864": "1k", "864x1536": "1k", "2016x864": "1k", "864x2016": "1k",
	"2048x2048": "2k", "1360x2048": "2k", "2048x1360": "2k", "2048x1536": "2k",
	"1536x2048": "2k", "2048x1152": "2k", "1152x2048": "2k", "2688x1152": "2k", "1152x2688": "2k",
	"2880x2880": "4k", "2336x3520": "4k", "3520x2336": "4k", "3312x2480": "4k",
	"2480x3312": "4k", "3840x2160": "4k", "2160x3840": "4k", "3840x1648": "4k", "1648x3840": "4k",
}

// sizeTierFromBody derives the 1k/2k/4k tier from a JSON request body: the
// preset table wins, otherwise the long side thresholds match infinite-canvas.
func sizeTierFromBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	size := strings.ToLower(strings.TrimSpace(gjson.GetBytes(body, "size").String()))
	if size == "" {
		size = strings.ToLower(strings.TrimSpace(gjson.GetBytes(body, "image_size").String()))
	}
	if size == "" {
		return ""
	}
	if tier, ok := sizeTierPresets[size]; ok {
		return tier
	}
	parts := strings.SplitN(size, "x", 2)
	if len(parts) != 2 {
		return ""
	}
	w, errW := strconv.Atoi(strings.TrimSpace(parts[0]))
	h, errH := strconv.Atoi(strings.TrimSpace(parts[1]))
	if errW != nil || errH != nil || w <= 0 || h <= 0 {
		return ""
	}
	longSide := w
	if h > longSide {
		longSide = h
	}
	switch {
	case longSide >= 3072:
		return "4k"
	case longSide >= 1536:
		return "2k"
	default:
		return "1k"
	}
}
