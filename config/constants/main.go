package constants

// syntax highlighter consts
const (
	LANGUAGE_NAME  = "GFXScript"
	LANGUAGE_ID    = "gfxs"
	DISPLAY_NAME   = "GFXScript Syntax Highlighter"
	PACKAGE_NAME   = "gfxscript-syntax-highlighter"
	DESCRIPTION    = "Syntax highlighting for GFXScript."
	README         = "# GFXScript Syntax Highlighter\n\nA syntax highlighter for the GFXScript compositing language."
	FILE_EXTENSION = ".gfxs"
	VERSION        = "1.0.1"
	PUBLISHER      = "tox"
)

// char consts
const (
	CHAR_SPACE    = ' '
	CHAR_TAB      = '\t'
	CHAR_ASSIGN   = '='
	CHAR_CLI_ARG  = '$'
	CHAR_COMMENT  = '#'
	CHAR_QUOTE    = '"'
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

// composition consts
const (
	COMP_WIDTH  = "width"
	COMP_HEIGHT = "height"
	COMP_COLOR  = "color"
	COMP_FILTER = "filter"
	COMP_CROP   = "crop"
	COMP_RESIZE = "resize"
)

// keyword consts
const (
	KEYWORD_USE  = "use"
	KEYWORD_HSLA = "hsla"
)

// calculated consts
const (
	SPACE    = string(CHAR_SPACE)
	TAB      = string(CHAR_TAB)
	ASSIGN   = string(CHAR_ASSIGN)
	COMMENT  = string(CHAR_COMMENT)
	QUOTE    = string(CHAR_QUOTE)
	ESCAPE   = string(CHAR_ESCAPE)
	LPAREN   = string(CHAR_LPAREN)
	RPAREN   = string(CHAR_RPAREN)
	LBRACE   = string(CHAR_LBRACE)
	RBRACE   = string(CHAR_RBRACE)
	LBRACKET = string(CHAR_LBRACKET)
	RBRACKET = string(CHAR_RBRACKET)
	COMMA    = string(CHAR_COMMA)
)
