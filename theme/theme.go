package theme

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v12/color"
	"github.com/essentialkaos/ek/v12/jsonutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Theme struct {
	Name      string            `json:"name"`
	Author    string            `json:"author"`
	Variables map[string]string `json:"variables"`
	Globals   Map               `json:"globals"`
	Rules     []*Rule           `json:"rules"`
}

type Rule struct {
	Name                string `json:"name"`
	Scope               string `json:"scope"`
	FontStyle           string `json:"font_style"`
	Foreground          string `json:"foreground"`
	Background          string `json:"background"`
	SelectionForeground string `json:"selection_foreground"`

	theme *Theme
}

type Map struct {
	Keys []string
	Data map[string]string

	theme *Theme
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	colorRegex = regexp.MustCompile(`color\((#[0-9a-fA-F]{3,8}) alpha\(([0-9\.]+)\)\)`)
	rgbRegex   = regexp.MustCompile(`rgb\(\d+, *\d+, *\d+\)`)
	rgbaRegex  = regexp.MustCompile(`rgba\((\d+), *(\d+), *(\d+)\, *([0-9\.]+)\)`)
	hslRegex   = regexp.MustCompile(`hsl\((\d+), *(\d+)%, *(\d+)%\)`)
	hslaRegex  = regexp.MustCompile(`hsla\((\d+), *(\d+)%, *(\d+)%, *([0-9\.]+)\)`)
	varRegex   = regexp.MustCompile(`var\(([^)]+)\)`)
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Read reads and parses JSON-encoded .sublime-color-scheme file
func Read(file string) (*Theme, error) {
	result := &Theme{}

	err := jsonutil.Read(file, result)

	if err != nil {
		return nil, err
	}

	// Inject pointer to theme into every rule
	for _, rule := range result.Rules {
		rule.theme = result
	}

	result.Globals.theme = result

	return result, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func (m *Map) UnmarshalJSON(data []byte) error {
	m.Data = make(map[string]string)

	d := json.NewDecoder(bytes.NewReader(data))

	var key, value string

	for {
		t, err := d.Token()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch v := t.(type) {
		case string:
			value = v
		default:
			continue
		}

		if key == "" {
			key = value
			continue
		}

		m.Data[key] = value
		m.Keys = append(m.Keys, key)

		key = ""
	}

	return nil
}

func (m *Map) Get(key string) string {
	return m.theme.eval(m.Data[key])
}

// ////////////////////////////////////////////////////////////////////////////////// //

// GetFontStyle evaluates variables and colors in rule and returns final value
func (r *Rule) GetFontStyle() string {
	return r.theme.eval(r.FontStyle)
}

// GetForeground evaluates variables and colors in rule and returns final value
func (r *Rule) GetForeground() string {
	return r.theme.eval(r.Foreground)
}

// GetBackground evaluates variables and colors in rule and returns final value
func (r *Rule) GetBackground() string {
	return r.theme.eval(r.Background)
}

// GetSelectionForeground evaluates variables and colors in rule and returns final value
func (r *Rule) GetSelectionForeground() string {
	return r.theme.eval(r.SelectionForeground)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// evalVariable evaluates given data into final value
func (t *Theme) eval(data string) string {
	if strings.Contains(data, "var(") {
		data = replaceVariables(data, t)
	}

	if strings.Contains(data, "color(") {
		data = convertColorToHex(data)
	}

	if strings.Contains(data, "rgb(") {
		data = convertRGBToHex(data)
	}

	if strings.Contains(data, "rgba(") {
		data = convertRGBAToHex(data)
	}

	if strings.Contains(data, "hsl(") {
		data = convertHSLToHex(data)
	}

	if strings.Contains(data, "hsla(") {
		data = convertHSLAToHex(data)
	}

	if strings.HasPrefix(data, "#") && len(data) == 4 {
		data = "#" + data[1:2] + data[1:2] + data[2:3] + data[2:3] + data[3:] + data[3:]
	}

	if strings.HasPrefix(data, "#") && len(data) == 5 {
		data = "#" + data[1:2] + data[1:2] + data[2:3] + data[2:3] +
			data[3:4] + data[3:4] + data[4:] + data[4:]
	}

	return data
}

// ////////////////////////////////////////////////////////////////////////////////// //

// replaceVariables replaces all variables in given data
func replaceVariables(data string, t *Theme) string {
	for i := 0; i < 16; i++ {
		if !varRegex.MatchString(data) {
			break
		}

		foundVar := varRegex.FindStringSubmatch(data)
		varValue, ok := t.Variables[foundVar[1]]

		if !ok {
			varValue = "[UNKNOWN-VAR:" + foundVar[1] + "]"
		}

		data = strings.ReplaceAll(data, foundVar[0], varValue)
	}

	return data
}

// convertColor converts color notation into hex
func convertColorToHex(data string) string {
	for i := 0; i < 16; i++ {
		if !colorRegex.MatchString(data) {
			break
		}

		foundColor := colorRegex.FindStringSubmatch(data)
		baseColor, alpha := foundColor[1], foundColor[2]
		colorHex, err := color.Parse(baseColor)

		if err != nil {
			data = strings.ReplaceAll(data, foundColor[0], "[COLOR-VALUE]")
			continue
		}

		alphaChannel, _ := strconv.ParseFloat(alpha, 64)

		if err != nil {
			data = strings.ReplaceAll(data, foundColor[0], "[COLOR-ALPHA]")
			continue
		}

		finalColor := colorHex.ToRGBA().WithAlpha(alphaChannel).ToHex().ToWeb(true, false)

		data = strings.ReplaceAll(data, foundColor[0], finalColor)
	}

	return data
}

// convertRGBToHex converts RGB color to Hex
func convertRGBToHex(data string) string {
	for i := 0; i < 16; i++ {
		if !rgbRegex.MatchString(data) {
			break
		}

		foundColor := rgbRegex.FindStringSubmatch(data)
		r, g, b := foundColor[1], foundColor[2], foundColor[3]

		ri, _ := strconv.Atoi(r)
		gi, _ := strconv.Atoi(g)
		bi, _ := strconv.Atoi(b)

		finalColor := color.RGB{R: uint8(ri), G: uint8(gi), B: uint8(bi)}.ToHex().ToWeb(true, false)

		data = strings.ReplaceAll(data, foundColor[0], finalColor)
	}

	return data
}

// convertRGBAToHex converts RGBA color to Hex
func convertRGBAToHex(data string) string {
	for i := 0; i < 16; i++ {
		if !rgbaRegex.MatchString(data) {
			break
		}

		foundColor := rgbaRegex.FindStringSubmatch(data)
		r, g, b, a := foundColor[1], foundColor[2], foundColor[3], foundColor[4]

		ri, _ := strconv.Atoi(r)
		gi, _ := strconv.Atoi(g)
		bi, _ := strconv.Atoi(b)
		af, _ := strconv.ParseFloat(a, 64)

		finalColor := color.RGBA{R: uint8(ri), G: uint8(gi), B: uint8(bi)}.WithAlpha(af).ToHex().ToWeb(true, false)

		data = strings.ReplaceAll(data, foundColor[0], finalColor)
	}

	return data
}

// convertHSLToHex converts HSL color to Hex
func convertHSLToHex(data string) string {
	for i := 0; i < 16; i++ {
		if !hslRegex.MatchString(data) {
			break
		}

		foundColor := hslRegex.FindStringSubmatch(data)
		h, s, l := foundColor[1], foundColor[2], foundColor[3]

		hf, _ := strconv.ParseFloat(h, 64)
		sf, _ := strconv.ParseFloat(s, 64)
		lf, _ := strconv.ParseFloat(l, 64)

		finalColor := color.HSL{H: hf, S: sf, L: lf}.ToRGB().ToHex().ToWeb(true, false)

		data = strings.ReplaceAll(data, foundColor[0], finalColor)
	}

	return data
}

// convertHSLAToHex converts HSLA color to Hex
func convertHSLAToHex(data string) string {
	for i := 0; i < 16; i++ {
		if !hslaRegex.MatchString(data) {
			break
		}

		foundColor := hslaRegex.FindStringSubmatch(data)
		h, s, l, a := foundColor[1], foundColor[2], foundColor[3], foundColor[4]

		hf, _ := strconv.ParseFloat(h, 64)
		sf, _ := strconv.ParseFloat(s, 64)
		lf, _ := strconv.ParseFloat(l, 64)
		af, _ := strconv.ParseFloat(a, 64)

		finalColor := color.HSL{H: hf, S: sf, L: lf, A: af}.ToRGBA().ToHex().ToWeb(true, false)

		data = strings.ReplaceAll(data, foundColor[0], finalColor)
	}

	return data
}
