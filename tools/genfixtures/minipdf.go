package main

import (
	"bytes"
	"fmt"
	"strings"
)

// MiniPDF writes a small, valid PDF with one Helvetica text block per page. It exists
// so the book's PDF fixtures are generated here, not copied from anywhere.
func MiniPDF(pages []string) []byte {
	var objs []string // 1-based object bodies, in order
	add := func(s string) int { objs = append(objs, s); return len(objs) }

	catalog := add("") // filled once the pages object exists
	pagesObj := add("")
	font := add("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>")

	var kids []string
	for _, text := range pages {
		var content strings.Builder
		content.WriteString("BT /F1 16 Tf 72 720 Td 20 TL\n")
		for _, line := range strings.Split(text, "\n") {
			fmt.Fprintf(&content, "(%s) Tj T*\n", escape(line))
		}
		content.WriteString("ET\n")
		stream := add(fmt.Sprintf("<< /Length %d >>\nstream\n%sendstream", content.Len(), content.String()))
		page := add(fmt.Sprintf("<< /Type /Page /Parent %d 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 %d 0 R >> >> /Contents %d 0 R >>", pagesObj, font, stream))
		kids = append(kids, fmt.Sprintf("%d 0 R", page))
	}
	objs[catalog-1] = fmt.Sprintf("<< /Type /Catalog /Pages %d 0 R >>", pagesObj)
	objs[pagesObj-1] = fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>", strings.Join(kids, " "), len(kids))

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n%\xe2\xe3\xcf\xd3\n")
	offsets := make([]int, len(objs))
	for i, body := range objs {
		offsets[i] = buf.Len()
		fmt.Fprintf(&buf, "%d 0 obj\n%s\nendobj\n", i+1, body)
	}
	xref := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 %d\n0000000000 65535 f \n", len(objs)+1)
	for _, off := range offsets {
		fmt.Fprintf(&buf, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(&buf, "trailer\n<< /Size %d /Root %d 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objs)+1, catalog, xref)
	return buf.Bytes()
}

func escape(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `(`, `\(`, `)`, `\)`)
	return r.Replace(s)
}
