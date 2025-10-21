package handler

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/phpdave11/gofpdf"

	"krstenica/internal/dto"
)

const (
	mmPerInch            = 25.4
	mmPerPoint           = 25.4 / 72.0
	pdfCellPaddingMM     = 0.6
	pdfBaselineFactor    = 0.68
	defaultFontSizePt    = 10.0
	defaultTextOffsetXMM = 0.0
	defaultTextOffsetYMM = -0.9
)

var forcedWrapCells = map[string]bool{
	"C58": true,
}

type textOffset struct {
	dx float64
	dy float64
}

var cellOffsets = map[string]textOffset{
	"H11": {dx: 1.6, dy: -0.9},
	"K11": {dx: -3.8, dy: -0.9},
	"F14": {dx: 0.0, dy: -1.4},
	"E20": {dx: 0.0, dy: -1.3},
}

type worksheetLayout struct {
	defaultColWidthMM  float64
	defaultRowHeightMM float64
	colWidthsMM        map[int]float64
	rowHeightsMM       map[int]float64
	scale              float64
	leftMarginMM       float64
	rightMarginMM      float64
	topMarginMM        float64
	bottomMarginMM     float64
	cellStyles         map[string]cellStyle
	maxCol             int
	maxRow             int
	contentWidthMM     float64
	contentHeightMM    float64
}

type cellStyle struct {
	fontSize float64
	wrapText bool
}

type cellRect struct {
	x      float64
	y      float64
	width  float64
	height float64
}

func fillKrstenicaPDFFile(krstenica *dto.Krstenica, templatePath, targetFile, backgroundImage string) error {
	layout, err := loadWorksheetLayout(templatePath)
	if err != nil {
		return fmt.Errorf("load worksheet layout: %w", err)
	}

	values := getKrstenicaCellValues(krstenica)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.SetMargins(0, 0, 0)
	pdf.AddPage()
	pdf.SetTextColor(0, 0, 0)

	if backgroundImage != "" {
		if _, err := os.Stat(backgroundImage); err == nil {
			if err := drawBackgroundImage(pdf, layout, backgroundImage); err != nil {
				return fmt.Errorf("draw background: %w", err)
			}
		}
	}

	cellOrder := []string{
		"C2", "C3", "C4", "C9", "C11", "H11", "K11",
		"F14", "E17", "G17", "I17", "E20", "F24", "H24",
		"D27", "F27", "H27", "E30", "G30", "I30", "E31", "G31",
		"G35", "K35", "E38", "E41", "G44", "F47", "H47",
		"E51", "G51", "I51", "E52", "G52", "E55", "C58", "I58", "K58",
		"F60", "C62",
	}

	paddingScaled := pdfCellPaddingMM * layout.scale

	for _, cell := range cellOrder {
		value, ok := values[cell]
		if !ok || strings.TrimSpace(value) == "" {
			continue
		}

		rect, ok := layout.cellRect(cell)
		if !ok {
			continue
		}

		style := layout.cellStyles[cell]
		fontSize := style.fontSize
		if fontSize <= 0 {
			fontSize = defaultFontSizePt
		}
		fontSizeScaled := fontSize * layout.scale
		if fontSizeScaled < 1 {
			fontSizeScaled = fontSize
		}

		pdf.SetFont("Helvetica", "", fontSizeScaled)

		wrapText := style.wrapText || forcedWrapCells[cell]
		offset := textOffset{dx: defaultTextOffsetXMM, dy: defaultTextOffsetYMM}
		if o, ok := cellOffsets[cell]; ok {
			offset = o
		}
		if wrapText {
			x := layout.leftMarginMM + (rect.x+pdfCellPaddingMM)*layout.scale + offset.dx*layout.scale
			y := layout.topMarginMM + rect.y*layout.scale + paddingScaled + offset.dy*layout.scale
			width := rect.width*layout.scale - 2*paddingScaled
			if width <= 0 {
				width = rect.width * layout.scale
			}
			lineHeight := rect.height * layout.scale
			if lineHeight <= 0 {
				lineHeight = fontSizeScaled * mmPerPoint * 1.2
			}
			pdf.SetXY(x, y)
			pdf.MultiCell(width, lineHeight, value, "", "L", false)
			continue
		}

		baseline := rect.height*pdfBaselineFactor + pdfCellPaddingMM
		x := layout.leftMarginMM + (rect.x+pdfCellPaddingMM)*layout.scale + offset.dx*layout.scale
		y := layout.topMarginMM + (rect.y+baseline)*layout.scale + offset.dy*layout.scale
		pdf.Text(x, y, value)
	}

	if err := pdf.OutputFileAndClose(targetFile); err != nil {
		return fmt.Errorf("write pdf: %w", err)
	}
	return nil
}

func drawBackgroundImage(pdf *gofpdf.Fpdf, layout *worksheetLayout, imagePath string) error {
	contentW := layout.contentWidthMM * layout.scale
	contentH := layout.contentHeightMM * layout.scale

	pageW, pageH := pdf.GetPageSize()

	if contentW <= 0 {
		contentW = pageW - layout.leftMarginMM - layout.rightMarginMM
	}
	if contentH <= 0 {
		contentH = pageH - layout.topMarginMM - layout.bottomMarginMM
	}
	if contentW <= 0 {
		contentW = pageW
	}
	if contentH <= 0 {
		contentH = pageH
	}

	x := layout.leftMarginMM
	y := layout.topMarginMM

	imgType := strings.ToUpper(strings.TrimPrefix(filepath.Ext(imagePath), "."))
	pdf.ImageOptions(imagePath, x, y, contentW, contentH, false, gofpdf.ImageOptions{ImageType: imgType}, 0, "")
	return nil
}

func (l *worksheetLayout) cellRect(cell string) (cellRect, bool) {
	colLabel, rowIdx, ok := splitCellRef(cell)
	if !ok {
		return cellRect{}, false
	}
	colIdx := columnLetterToIndex(colLabel)
	if colIdx <= 0 || rowIdx <= 0 {
		return cellRect{}, false
	}

	x := 0.0
	for i := 1; i < colIdx; i++ {
		x += l.colWidth(i)
	}
	y := 0.0
	for i := 1; i < rowIdx; i++ {
		y += l.rowHeight(i)
	}

	width := l.colWidth(colIdx)
	height := l.rowHeight(rowIdx)

	return cellRect{x: x, y: y, width: width, height: height}, true
}

func (l *worksheetLayout) colWidth(idx int) float64 {
	if w, ok := l.colWidthsMM[idx]; ok && w > 0 {
		return w
	}
	return l.defaultColWidthMM
}

func (l *worksheetLayout) rowHeight(idx int) float64 {
	if h, ok := l.rowHeightsMM[idx]; ok && h > 0 {
		return h
	}
	return l.defaultRowHeightMM
}

func loadWorksheetLayout(templatePath string) (*worksheetLayout, error) {
	reader, err := zip.OpenReader(templatePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var sheetData []byte
	var stylesData []byte

	for _, file := range reader.File {
		switch file.Name {
		case "xl/worksheets/sheet1.xml":
			sheetData, err = readZipFile(file)
			if err != nil {
				return nil, err
			}
		case "xl/styles.xml":
			stylesData, err = readZipFile(file)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(sheetData) == 0 {
		return nil, fmt.Errorf("sheet1.xml not found in template")
	}

	var worksheet worksheetXML
	if err := xml.Unmarshal(sheetData, &worksheet); err != nil {
		return nil, err
	}

	var styles stylesXML
	if len(stylesData) > 0 {
		if err := xml.Unmarshal(stylesData, &styles); err != nil {
			return nil, err
		}
	}

	layout := &worksheetLayout{
		defaultColWidthMM:  widthExcelToMM(8.43),
		defaultRowHeightMM: rowHeightToMM(18.55),
		colWidthsMM:        make(map[int]float64),
		rowHeightsMM:       make(map[int]float64),
		scale:              1.0,
		cellStyles:         make(map[string]cellStyle),
	}

	if worksheet.SheetFormatPr != nil {
		if worksheet.SheetFormatPr.DefaultColWidth > 0 {
			layout.defaultColWidthMM = widthExcelToMM(worksheet.SheetFormatPr.DefaultColWidth)
		}
		if worksheet.SheetFormatPr.DefaultRowHeight > 0 {
			layout.defaultRowHeightMM = rowHeightToMM(worksheet.SheetFormatPr.DefaultRowHeight)
		}
	}

	for _, col := range worksheet.Cols {
		width := widthExcelToMM(col.Width)
		if width <= 0 {
			width = layout.defaultColWidthMM
		}
		for i := col.Min; i <= col.Max; i++ {
			layout.colWidthsMM[i] = width
		}
	}

	for _, row := range worksheet.SheetData.Rows {
		height := layout.defaultRowHeightMM
		if row.Height > 0 {
			height = rowHeightToMM(row.Height)
		}
		layout.rowHeightsMM[row.Index] = height
	}

	if worksheet.PageSetup != nil && worksheet.PageSetup.Scale > 0 {
		layout.scale = worksheet.PageSetup.Scale / 100.0
	}
	if layout.scale <= 0 {
		layout.scale = 1.0
	}

	if worksheet.PageMargins != nil {
		layout.leftMarginMM = worksheet.PageMargins.Left * mmPerInch
		layout.rightMarginMM = worksheet.PageMargins.Right * mmPerInch
		layout.topMarginMM = worksheet.PageMargins.Top * mmPerInch
		layout.bottomMarginMM = worksheet.PageMargins.Bottom * mmPerInch
	}

	fonts := extractFontSizes(styles)
	xfStyles := extractXfStyles(styles)

	for _, row := range worksheet.SheetData.Rows {
		for _, cell := range row.Cells {
			styleIdx := cell.Style
			style := cellStyle{fontSize: defaultFontSizePt}
			if styleIdx >= 0 && styleIdx < len(xfStyles) {
				xf := xfStyles[styleIdx]
				if xf.fontID >= 0 && xf.fontID < len(fonts) && fonts[xf.fontID] > 0 {
					style.fontSize = fonts[xf.fontID]
				}
				style.wrapText = xf.wrap
			}
			layout.cellStyles[cell.Ref] = style
		}
	}

	if worksheet.Dimension != nil && worksheet.Dimension.Ref != "" {
		if _, _, maxCol, maxRow, ok := parseDimension(worksheet.Dimension.Ref); ok {
			layout.maxCol = maxCol
			layout.maxRow = maxRow
		}
	}

	if layout.maxRow == 0 {
		for idx := range layout.rowHeightsMM {
			if idx > layout.maxRow {
				layout.maxRow = idx
			}
		}
	}
	if layout.maxRow == 0 {
		layout.maxRow = len(layout.rowHeightsMM)
	}

	if layout.maxCol == 0 {
		for _, row := range worksheet.SheetData.Rows {
			for _, cell := range row.Cells {
				if colLabel, _, ok := splitCellRef(cell.Ref); ok {
					colIdx := columnLetterToIndex(colLabel)
					if colIdx > layout.maxCol {
						layout.maxCol = colIdx
					}
				}
			}
		}
	}
	if layout.maxCol == 0 {
		layout.maxCol = len(layout.colWidthsMM)
	}
	if layout.maxCol == 0 {
		layout.maxCol = 13
	}
	if layout.maxRow == 0 {
		layout.maxRow = 66
	}

	for i := 1; i <= layout.maxCol; i++ {
		layout.contentWidthMM += layout.colWidth(i)
	}
	for i := 1; i <= layout.maxRow; i++ {
		layout.contentHeightMM += layout.rowHeight(i)
	}

	return layout, nil
}

func extractFontSizes(styles stylesXML) []float64 {
	fonts := make([]float64, len(styles.Fonts.Fonts))
	for i, font := range styles.Fonts.Fonts {
		if font.Size.Val > 0 {
			fonts[i] = font.Size.Val
		} else {
			fonts[i] = defaultFontSizePt
		}
	}
	return fonts
}

type xfStyle struct {
	fontID int
	wrap   bool
}

func extractXfStyles(styles stylesXML) []xfStyle {
	xfs := make([]xfStyle, len(styles.CellXfs.Xfs))
	for i, xf := range styles.CellXfs.Xfs {
		wrap := false
		if xf.Alignment != nil && xf.Alignment.WrapText {
			wrap = true
		}
		xfs[i] = xfStyle{fontID: xf.FontID, wrap: wrap}
	}
	return xfs
}

func readZipFile(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func columnLetterToIndex(col string) int {
	result := 0
	for i := 0; i < len(col); i++ {
		ch := col[i]
		if ch < 'A' || ch > 'Z' {
			continue
		}
		result = result*26 + int(ch-'A'+1)
	}
	return result
}

func splitCellRef(cell string) (string, int, bool) {
	if cell == "" {
		return "", 0, false
	}
	var idx int
	for idx = 0; idx < len(cell); idx++ {
		if cell[idx] >= '0' && cell[idx] <= '9' {
			break
		}
	}
	if idx == 0 || idx == len(cell) {
		return "", 0, false
	}
	row, err := strconv.Atoi(cell[idx:])
	if err != nil {
		return "", 0, false
	}
	return strings.ToUpper(cell[:idx]), row, true
}

func parseDimension(ref string) (minCol, minRow, maxCol, maxRow int, ok bool) {
	if ref == "" {
		return 0, 0, 0, 0, false
	}
	parts := strings.Split(ref, ":")
	switch len(parts) {
	case 1:
		colLabel, row, valid := splitCellRef(parts[0])
		if !valid {
			return 0, 0, 0, 0, false
		}
		idx := columnLetterToIndex(colLabel)
		return idx, row, idx, row, true
	case 2:
		colLabelStart, rowStart, validStart := splitCellRef(parts[0])
		colLabelEnd, rowEnd, validEnd := splitCellRef(parts[1])
		if !validStart || !validEnd {
			return 0, 0, 0, 0, false
		}
		startCol := columnLetterToIndex(colLabelStart)
		endCol := columnLetterToIndex(colLabelEnd)
		minCol, maxCol = startCol, endCol
		if endCol < startCol {
			minCol, maxCol = endCol, startCol
		}
		minRow, maxRow = rowStart, rowEnd
		if rowEnd < rowStart {
			minRow, maxRow = rowEnd, rowStart
		}
		return minCol, minRow, maxCol, maxRow, true
	default:
		return 0, 0, 0, 0, false
	}
}

func widthExcelToMM(width float64) float64 {
	if width <= 0 {
		return 0
	}
	pixels := 0.0
	if width < 1 {
		pixels = math.Floor(width*12 + 0.5)
	} else {
		pixels = math.Floor((256*width + math.Floor(128.0/7.0)) / 256.0 * 7.0)
	}
	return pixels * mmPerInch / 96.0
}

func rowHeightToMM(points float64) float64 {
	if points <= 0 {
		return 0
	}
	return points * mmPerPoint
}

// Structures for worksheet and styles XML parsing

type worksheetXML struct {
	SheetFormatPr *sheetFormatPrXML `xml:"sheetFormatPr"`
	Cols          []colXML          `xml:"cols>col"`
	SheetData     struct {
		Rows []rowXML `xml:"row"`
	} `xml:"sheetData"`
	Dimension   *dimensionXML   `xml:"dimension"`
	PageMargins *pageMarginsXML `xml:"pageMargins"`
	PageSetup   *pageSetupXML   `xml:"pageSetup"`
}

type sheetFormatPrXML struct {
	DefaultColWidth  float64 `xml:"defaultColWidth,attr"`
	DefaultRowHeight float64 `xml:"defaultRowHeight,attr"`
}

type colXML struct {
	Min   int     `xml:"min,attr"`
	Max   int     `xml:"max,attr"`
	Width float64 `xml:"width,attr"`
}

type rowXML struct {
	Index  int       `xml:"r,attr"`
	Height float64   `xml:"ht,attr"`
	Cells  []cellXML `xml:"c"`
}

type cellXML struct {
	Ref   string `xml:"r,attr"`
	Style int    `xml:"s,attr"`
}

type pageMarginsXML struct {
	Left   float64 `xml:"left,attr"`
	Right  float64 `xml:"right,attr"`
	Top    float64 `xml:"top,attr"`
	Bottom float64 `xml:"bottom,attr"`
}

type dimensionXML struct {
	Ref string `xml:"ref,attr"`
}

type pageSetupXML struct {
	Scale float64 `xml:"scale,attr"`
}

type stylesXML struct {
	Fonts   fontsXML   `xml:"fonts"`
	CellXfs cellXfsXML `xml:"cellXfs"`
}

type fontsXML struct {
	Fonts []fontXML `xml:"font"`
}

type fontXML struct {
	Size fontSizeXML `xml:"sz"`
}

type fontSizeXML struct {
	Val float64 `xml:"val,attr"`
}

type cellXfsXML struct {
	Xfs []xfXML `xml:"xf"`
}

type xfXML struct {
	FontID    int           `xml:"fontId,attr"`
	Alignment *alignmentXML `xml:"alignment"`
}

type alignmentXML struct {
	WrapText bool `xml:"wrapText,attr"`
}
