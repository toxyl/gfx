package parser

import (
	"strings"

	"github.com/toxyl/gfx/color/blend"
)

const (
	LANGUAGE_NAME  = "GFXScript"
	LANGUAGE_ID    = "gfxscript"
	DISPLAY_NAME   = "GFXScript Syntax Highlighter"
	PACKAGE_NAME   = "gfxscript-syntax-highlighter"
	DESCRIPTION    = "Syntax highlighting for GFXScript."
	README         = "# GFXScript Syntax Highlighter\n\nA syntax highlighter for the GFXScript compositing language."
	FILE_EXTENSION = ".gfxs"
	VERSION        = "1.0.0"
)

// char consts
const (
	CHAR_SPACE    = ' '
	CHAR_TAB      = '\t'
	CHAR_ASSIGN   = '='
	CHAR_CLI_ARG  = '$'
	CHAR_COMMENT  = '#'
	CHAR_QUOTE    = '`'
	CHAR_ESCAPE   = '\\'
	CHAR_LPAREN   = '('
	CHAR_RPAREN   = ')'
	CHAR_LBRACE   = '{'
	CHAR_RBRACE   = '}'
	CHAR_LBRACKET = '['
	CHAR_RBRACKET = ']'
	CHAR_COMMA    = ','
)

// section consts
const (
	SECTION_VARS        = "VARS"
	SECTION_FILTERS     = "FILTERS"
	SECTION_COMPOSITION = "COMPOSITION"
	SECTION_LAYERS      = "LAYERS"
)

var (
	SECTIONS = []string{SECTION_VARS, SECTION_FILTERS, SECTION_COMPOSITION, SECTION_LAYERS}
)

// composition consts
const (
	COMP_NAME   = "name"
	COMP_WIDTH  = "width"
	COMP_HEIGHT = "height"
	COMP_COLOR  = "color"
	COMP_FILTER = "filter"
	COMP_CROP   = "crop"
	COMP_RESIZE = "resize"
)

var (
	COMPOSITION = []string{COMP_NAME, COMP_WIDTH, COMP_HEIGHT, COMP_COLOR, COMP_CROP, COMP_FILTER, COMP_RESIZE}
)

// layer consts
const (
	LAYER_CROP   = "crop"
	LAYER_RESIZE = "resize"
	LAYER_OFFSET = "offset"
)

var (
	LAYER = []string{LAYER_CROP, LAYER_RESIZE, LAYER_OFFSET}
)

// keyword consts
const (
	KEYWORD_USE = "use"
)

var (
	KEYWORDS = []string{KEYWORD_USE}
)

// blendmode constants
var (
	BLENDMODES = []string{
		string(blend.ADD),
		string(blend.AVERAGE),
		string(blend.COLOR_BURN),
		string(blend.DARKEN),
		string(blend.DIFFERENCE),
		string(blend.DIVIDE),
		string(blend.ERASE),
		string(blend.EXCLUSION),
		string(blend.HARD_LIGHT),
		string(blend.LIGHTEN),
		string(blend.LINEAR_BURN),
		string(blend.MULTIPLY),
		string(blend.NEGATION),
		string(blend.NORMAL),
		string(blend.OVERLAY),
		string(blend.PIN_LIGHT),
		string(blend.SCREEN),
		string(blend.SOFT_LIGHT),
		string(blend.SUBTRACT),
	}
)

// calculated consts
const (
	STR_SPACE    = string(CHAR_SPACE)
	STR_TAB      = string(CHAR_TAB)
	STR_ASSIGN   = string(CHAR_ASSIGN)
	STR_COMMENT  = string(CHAR_COMMENT)
	STR_QUOTE    = string(CHAR_QUOTE)
	STR_ESCAPE   = string(CHAR_ESCAPE)
	STR_LPAREN   = string(CHAR_LPAREN)
	STR_RPAREN   = string(CHAR_RPAREN)
	STR_LBRACE   = string(CHAR_LBRACE)
	STR_RBRACE   = string(CHAR_RBRACE)
	STR_LBRACKET = string(CHAR_LBRACKET)
	STR_RBRACKET = string(CHAR_RBRACKET)
	STR_COMMA    = string(CHAR_COMMA)
)

// patterns
var (
	WORD_PATTERN        = `\b[a-zA-Z][a-zA-Z0-9-_]+?\b`
	NUMBER_PATTERN      = `\b-?\d+(\.\d+)?\b`
	ALPHA_PATTERN       = `\b(0|1)\.\d+?\b`
	FILE_PATTERN        = `\.\./.*|\./.*|/.*`
	URL_PATTERN         = `\b(http|ftp)s{0,1}://\S*\b`
	CLI_ARG_PATTERN     = `\$\d+`
	SECTIONS_PATTERN    = `\b(` + strings.Join(SECTIONS, "|") + `)\b`
	COMPOSITION_PATTERN = `\b(` + strings.Join(COMPOSITION, "|") + `)\b`
	LAYER_PATTERN       = `\b(` + strings.Join(LAYER, "|") + `)\b`
	KEYWORDS_PATTERN    = `\b(` + strings.Join(KEYWORDS, "|") + `)\b`
	BLENDMODES_PATTERN  = `\b(` + strings.Join(BLENDMODES, "|") + `)\b`
	SECTION_PATTERN     = `\` + STR_LBRACKET + SECTIONS_PATTERN + `\` + STR_RBRACKET
	SOURCE_PATTERN      = `(` + FILE_PATTERN + `|` + URL_PATTERN + `|` + CLI_ARG_PATTERN + `)`
)
