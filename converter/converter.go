package converter

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/essentialkaos/detheme/theme"
	"github.com/essentialkaos/ek/v12/uuid"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Convert converts given theme to thTheme
func Convert(t *theme.Theme) ([]byte, error) {
	buf := &bytes.Buffer{}

	writeHeader(buf, t)
	writeBasicInfo(buf, t)
	writeBody(buf, t)
	writeFooter(buf, t)

	return buf.Bytes(), nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// writeHeader writes header data
func writeHeader(buf *bytes.Buffer, t *theme.Theme) {
	fmt.Fprintln(buf, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(buf, `<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">`)
	fmt.Fprintln(buf, "")
	fmt.Fprintln(buf, "<!--")
	fmt.Fprintln(buf, "  This theme converted from sublime-color-scheme by detheme (https://kaos.sh/detheme)")
	fmt.Fprintln(buf, "-->")
	fmt.Fprintln(buf, "")
	fmt.Fprintln(buf, `<plist version="1.0">`)
	fmt.Fprintln(buf, `  <dict>`)
}

// writeFooter writes footer data
func writeFooter(buf *bytes.Buffer, t *theme.Theme) {
	fmt.Fprintln(buf, `  </dict>`)
	fmt.Fprintln(buf, `</plist>`)
}

// writeBasicInfo writes basic theme info (name, author, color space, uuid)
func writeBasicInfo(buf *bytes.Buffer, t *theme.Theme) {
	writeNode(buf, 2, "name", t.Name)

	if t.Author != "" {
		writeNode(buf, 2, "author", t.Author)
	}

	writeNode(buf, 2, "colorSpaceName", "sRGB")
	writeNode(buf, 2, "uuid", uuid.UUID4().String())
}

// writeBody writes settings and rules
func writeBody(buf *bytes.Buffer, t *theme.Theme) {
	fmt.Fprintln(buf, `    <key>settings</key>`)
	fmt.Fprintln(buf, `    <array>`)
	fmt.Fprintln(buf, `      <dict>`)
	fmt.Fprintln(buf, `        <key>settings</key>`)
	fmt.Fprintln(buf, `        <dict>`)

	for _, key := range t.Globals.Keys {
		writeNode(buf, 5, key, t.Globals.Get(key))
	}

	fmt.Fprintln(buf, `        </dict>`)
	fmt.Fprintln(buf, `      </dict>`)

	for _, rule := range t.Rules {
		writeRule(buf, rule)
	}

	fmt.Fprintln(buf, `    </array>`)
}

// writeRule writes rule data
func writeRule(buf *bytes.Buffer, r *theme.Rule) {
	fmt.Fprintln(buf, `      <dict>`)
	writeNode(buf, 4, "name", r.Name)
	writeNode(buf, 4, "scope", r.Scope)
	fmt.Fprintln(buf, `        <key>settings</key>`)
	fmt.Fprintln(buf, `        <dict>`)

	if r.FontStyle != "" {
		writeNode(buf, 5, "fontStyle", r.GetFontStyle())
	}

	if r.Foreground != "" {
		writeNode(buf, 5, "foreground", r.GetForeground())
	}

	if r.Background != "" {
		writeNode(buf, 5, "background", r.GetBackground())
	}

	if r.SelectionForeground != "" {
		writeNode(buf, 5, "selectionForeground", r.GetSelectionForeground())
	}

	fmt.Fprintln(buf, `        </dict>`)
	fmt.Fprintln(buf, `      </dict>`)
}

// writeNode write node info (key + string)
func writeNode(buf *bytes.Buffer, indent int, key, value string) {
	spaces := strings.Repeat("  ", indent)

	key = escape(key)
	value = escape(value)

	fmt.Fprintf(buf, "%s<key>%s</key>\n", spaces, key)
	fmt.Fprintf(buf, "%s<string>%s</string>\n", spaces, value)
}

func escape(data string) string {
	data = strings.ReplaceAll(data, "'", "&apos;")
	data = strings.ReplaceAll(data, "\"", "&quot;")
	data = strings.ReplaceAll(data, "&", "&amp;")
	data = strings.ReplaceAll(data, "<", "&lt;")
	data = strings.ReplaceAll(data, ">", "&gt;")

	return data
}
