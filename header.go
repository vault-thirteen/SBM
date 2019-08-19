package sbm

// ASCII Symbols.
const (
	CR = "\r"
	LF = "\n"
	NL = CR + LF
)

// Meta-Data Header Parameters.
const (
	HeaderPartsSeparator    = " "
	HeaderPartsBracketLeft  = "("
	HeaderPartsBracketRight = ")"
	HeaderPartsPlus         = "+"

	Header_FormatName = "SBM (SIMPLE BIT MAP)" + NL

	HeaderPrefix_Version = "VERSION"
	HeaderFormat_Version = HeaderPrefix_Version +
		HeaderPartsSeparator +
		"%v" +
		NL

	HeaderPrefix_Width = "WIDTH"
	HeaderFormat_Width = HeaderPrefix_Width +
		HeaderPartsSeparator +
		"%v" +
		HeaderPartsSeparator +
		HeaderPartsBracketLeft +
		"%v" +
		HeaderPartsSeparator +
		HeaderPartsPlus +
		HeaderPartsSeparator +
		"%v" +
		HeaderPartsBracketRight +
		NL

	HeaderPrefix_Height = "HEIGHT"
	HeaderFormat_Height = HeaderPrefix_Height +
		HeaderPartsSeparator +
		"%v" +
		HeaderPartsSeparator +
		HeaderPartsBracketLeft +
		"%v" +
		HeaderPartsSeparator +
		HeaderPartsPlus +
		HeaderPartsSeparator +
		"%v" +
		HeaderPartsBracketRight +
		NL

	HeaderPrefix_Area = "AREA"
	HeaderFormat_Area = HeaderPrefix_Area +
		HeaderPartsSeparator +
		"%v" +
		HeaderPartsSeparator +
		HeaderPartsBracketLeft +
		"%v" +
		HeaderPartsSeparator +
		HeaderPartsPlus +
		HeaderPartsSeparator +
		"%v" +
		HeaderPartsBracketRight +
		NL
)

// MIME Type.
// This MIME Type is un-official, not registered in IANA.
const (
	MimeType = "image/x-portable-bitmap"
)
