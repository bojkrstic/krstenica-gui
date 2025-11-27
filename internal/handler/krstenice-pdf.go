package handler

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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
	pdfFontScaleFactor   = 4.0 / 3.0
)

var forcedWrapCells = map[string]bool{
	"C58": true,
}

var boldCells = map[string]bool{
	"F27": true,
}

type textOffset struct {
	dx float64
	dy float64
}

var cellOffsets = map[string]textOffset{
	"C1":  {dx: 0.0, dy: -8.0},
	"C2":  {dx: 0.0, dy: -7.0},
	"C3":  {dx: 0.0, dy: -3.5},
	"F8":  {dx: 12.0, dy: -5.0},
	"C10": {dx: 0.0, dy: -6.0},
	"I10": {dx: 10.6, dy: -6.0},
	"N10": {dx: 2.0, dy: -6.0},
	"F13": {dx: 12.0, dy: -5.0},
	"E16": {dx: 12.0, dy: -3.0},
	"F16": {dx: 12.0, dy: -3.0},
	"G16": {dx: 12.0, dy: -3.0},
	"F19": {dx: 12.0, dy: 3.0},
	"G24": {dx: 12.0, dy: -6.0},
	"I24": {dx: 12.0, dy: -6.0},
	"D27": {dx: 12.0, dy: -8.0},
	"F27": {dx: 12.0, dy: -8.0},
	"I27": {dx: 0.0, dy: -8.0},
	"F30": {dx: 12.0, dy: -8.0},
	"I30": {dx: 0.0, dy: -8.0},
	"K30": {dx: 0.0, dy: -8.0},
	"F31": {dx: 12.0, dy: -7.0},
	"I31": {dx: 0.0, dy: -7.0},
	"I32": {dx: -10.0, dy: 7.0},
	// "K32": {dx: 7.0, dy: 3.0},
	"I36": {dx: -10.0, dy: -3.0},
	"I38": {dx: -10.0, dy: 3.0},
	"I41": {dx: 5.0, dy: 1.0},
	"F43": {dx: 12.0, dy: 9.0},
	"H43": {dx: 10.0, dy: 9.0},
	"K43": {dx: 0.0, dy: 9.0},
	"E48": {dx: 10.0, dy: -2.0},
	"G48": {dx: 10.0, dy: -2.0},
	// "K48": {dx: 0.0, dy: -2.0},
	"E49": {dx: 10.0, dy: 0.0},
	"G49": {dx: 10.0, dy: 0.0},
	"E51": {dx: 30.0, dy: 7.0},
	"C54": {dx: 20.0, dy: 1.0},
	"B62": {dx: 10.0, dy: -2.0},
	"B63": {dx: 0.0, dy: 4.0},
	"C63": {dx: 8.0, dy: 4.0},
	"B65": {dx: 2.0, dy: 2.0},
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

func formatCyrillicIzCity(city string) string {
	trimmed := strings.TrimSpace(city)
	if trimmed == "" {
		return ""
	}
	cityText := fmt.Sprintf("из %s", trimmed)
	runes := []rune(cityText)
	if len(runes) == 0 {
		return cityText
	}
	lastIdx := len(runes) - 1
	switch runes[lastIdx] {
	case 'а':
		return cityText
	case 'a':
		runes[lastIdx] = 'а'
	default:
		runes = append(runes, 'а')
	}
	return string(runes)
}

func fillKrstenicaPDFFile(krstenica *dto.Krstenica, templatePath, targetFile, backgroundImage string, fullBleed bool) error {
	layout, err := loadWorksheetLayout(templatePath)
	if err != nil {
		return fmt.Errorf("load worksheet layout: %w", err)
	}

	values := getKrstenicaCellValues(krstenica)
	values["G32"] = strings.TrimSpace(krstenica.BirthOrder)
	priestFirst := strings.TrimSpace(values["F43"])
	priestLast := strings.TrimSpace(values["H43"])
	priestTitle := strings.TrimSpace(values["K43"])
	var priestParts []string
	if priestFirst != "" {
		priestParts = append(priestParts, priestFirst)
	}
	if priestLast != "" {
		priestParts = append(priestParts, priestLast)
	}
	priestText := strings.Join(priestParts, " ")
	if priestTitle != "" {
		if priestText != "" {
			priestText = fmt.Sprintf("%s, %s", priestText, priestTitle)
		} else {
			priestText = priestTitle
		}
	}
	if priestText != "" {
		values["F43"] = priestText
		values["H43"] = ""
		values["K43"] = ""
	}
	parentFirst := strings.TrimSpace(values["F30"])
	parentLast := strings.TrimSpace(values["I30"])
	parentOccupation := strings.TrimSpace(values["K30"])
	var parentParts []string
	if parentFirst != "" {
		parentParts = append(parentParts, parentFirst)
	}
	if parentLast != "" {
		parentParts = append(parentParts, parentLast)
	}
	parentText := strings.Join(parentParts, " ")
	if parentOccupation != "" {
		if parentText != "" {
			parentText = fmt.Sprintf("%s, %s", parentText, parentOccupation)
		} else {
			parentText = parentOccupation
		}
	}
	if parentText != "" {
		values["F30"] = parentText
		values["I30"] = ""
		values["K30"] = ""
	}
	godfatherFirst := strings.TrimSpace(values["E48"])
	godfatherLast := strings.TrimSpace(values["G48"])
	godfatherOccupation := strings.TrimSpace(values["K48"])
	var godfatherParts []string
	if godfatherFirst != "" {
		godfatherParts = append(godfatherParts, godfatherFirst)
	}
	if godfatherLast != "" {
		godfatherParts = append(godfatherParts, godfatherLast)
	}
	godfatherText := strings.Join(godfatherParts, " ")
	if godfatherText != "" {
		if godfatherOccupation != "" {
			godfatherText = fmt.Sprintf("%s, %s", godfatherText, godfatherOccupation)
		} else {
			godfatherText = godfatherText + ","
		}
		values["E48"] = godfatherText
		values["G48"] = ""
		values["K48"] = ""
	}
	templeCity := strings.TrimSpace(values["G24"])
	templeName := strings.TrimSpace(values["I24"])
	var templeParts []string
	if templeCity != "" {
		cityText := strings.TrimSuffix(templeCity, ",")
		cityText = strings.TrimSpace(cityText)
		if cityText != "" {
			templeParts = append(templeParts, cityText+",")
		}
	}
	if templeName != "" {
		templeParts = append(templeParts, templeName)
	}
	if len(templeParts) > 0 {
		values["G24"] = strings.Join(templeParts, " ")
		values["I24"] = ""
	}
	if cityText := formatCyrillicIzCity(values["F31"]); cityText != "" {
		values["F31"] = cityText
	}
	if religion, ok := values["I31"]; ok {
		religion = strings.TrimSpace(religion)
		if religion != "" {
			religion = religion + " "
		}
		values["I31"] = religion
	}
	if godfatherCity := formatCyrillicIzCity(values["E49"]); godfatherCity != "" {
		values["E49"] = godfatherCity
	}
	if godfatherReligion, ok := values["G49"]; ok {
		godfatherReligion = strings.TrimSpace(godfatherReligion)
		if godfatherReligion != "" {
			godfatherReligion = godfatherReligion + " "
		}
		values["G49"] = godfatherReligion
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.SetMargins(0, 0, 0)
	pdf.AddPage()
	pdf.SetTextColor(0, 0, 0)

	pdf.AddUTF8FontFromBytes("DejaVuSans", "", dejavuSansFont)
	if err := pdf.Error(); err != nil {
		return fmt.Errorf("register utf-8 font: %w", err)
	}
	pdf.AddUTF8FontFromBytes("DejaVuSans", "B", dejavuSansFontBold)
	if err := pdf.Error(); err != nil {
		return fmt.Errorf("register bold utf-8 font: %w", err)
	}

	if backgroundImage != "" {
		if _, err := os.Stat(backgroundImage); err == nil {
			if err := drawBackgroundImage(pdf, layout, backgroundImage, fullBleed); err != nil {
				return fmt.Errorf("draw background: %w", err)
			}
		}
	}

	cellOrder := []string{
		"C1", "C2", "C3", "F8", "C10", "I10", "N10",
		"F13", "E16", "F16", "G16", "F19", "G24", "I24",
		"D27", "F27", "I27", "F30", "I30", "K30", "F31", "I31",
		"I32", "K32", "I36", "I38", "I41", "F43", "H43", "K43",
		"E48", "G48", "E49", "G49", "E51", "C54", "B62", "B63", "C63", "B65",
	}

	paddingScaled := pdfCellPaddingMM * layout.scale
	uniformFontSizePt := defaultFontSizePt
	if refStyle, ok := layout.cellStyles["C9"]; ok && refStyle.fontSize > 0 {
		uniformFontSizePt = refStyle.fontSize
	}
	uniformFontSizePt *= pdfFontScaleFactor

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
		fontSize := uniformFontSizePt
		fontSizeScaled := fontSize * layout.scale
		if fontSizeScaled < 1 {
			fontSizeScaled = fontSize
		}

		fontStyle := ""
		if boldCells[cell] {
			fontStyle = "B"
		}
		pdf.SetFont("DejaVuSans", fontStyle, fontSizeScaled)

		wrapText := style.wrapText || forcedWrapCells[cell]
		offset := textOffset{dx: defaultTextOffsetXMM, dy: defaultTextOffsetYMM}
		if o, ok := cellOffsets[cell]; ok {
			offset = o
		}
		if cell == "N10" {
			if trimmed := strings.TrimSpace(value); len(trimmed) == 4 {
				offset.dx -= 8.0
			}
		}
		if wrapText {
			x := layout.leftMarginMM + (rect.x+pdfCellPaddingMM)*layout.scale + offset.dx*layout.scale
			y := layout.topMarginMM + rect.y*layout.scale + paddingScaled + offset.dy*layout.scale
			width := rect.width*layout.scale - 2*paddingScaled
			if width <= 0 {
				width = rect.width * layout.scale
			}
			lineHeight := rect.height * layout.scale
			minLineHeight := fontSizeScaled * mmPerPoint * 1.2
			if lineHeight < minLineHeight {
				lineHeight = minLineHeight
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

func drawBackgroundImage(pdf *gofpdf.Fpdf, layout *worksheetLayout, imagePath string, fullBleed bool) error {
	contentW := layout.contentWidthMM * layout.scale
	contentH := layout.contentHeightMM * layout.scale

	pageW, pageH := pdf.GetPageSize()
	imgType := strings.ToUpper(strings.TrimPrefix(filepath.Ext(imagePath), "."))

	if fullBleed {
		targetW := pageW
		targetH := pageH
		offsetX := 0.0
		offsetY := 0.0

		if imgW, imgH, err := imageSize(imagePath); err == nil && imgW > 0 && imgH > 0 {
			imgRatio := float64(imgW) / float64(imgH)
			pageRatio := pageW / pageH
			if imgRatio > 0 {
				if imgRatio > pageRatio {
					targetH = pageH
					targetW = pageH * imgRatio
					offsetX = (pageW - targetW) / 2
				} else {
					targetW = pageW
					targetH = pageW / imgRatio
					offsetY = (pageH - targetH) / 2
				}
			}
		}

		pdf.ImageOptions(imagePath, offsetX, offsetY, targetW, targetH, false, gofpdf.ImageOptions{ImageType: imgType}, 0, "")
		return nil
	}

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

	pdf.ImageOptions(imagePath, x, y, contentW, contentH, false, gofpdf.ImageOptions{ImageType: imgType}, 0, "")
	return nil
}

func imageSize(path string) (int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
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
