package config

import (
	"encoding/json"
	"regexp"
	"strings"

	c "github.com/toxyl/gfx/config/constants"
)

var (
	SECTIONS    = []string{c.SECTION_VARS, c.SECTION_FILTERS, c.SECTION_COMPOSITION, c.SECTION_LAYERS}
	COMPOSITION = []string{c.COMP_WIDTH, c.COMP_HEIGHT, c.COMP_COLOR, c.COMP_CROP, c.COMP_FILTER, c.COMP_RESIZE}
	KEYWORDS    = []string{c.KEYWORD_USE, c.KEYWORD_HSLA}
	FUNCTIONS   = []string{}
	BLENDMODES  = []string{}
	COLORMODELS = []string{}
	PROJECTIONS = []string{}
)

// patterns
var (
	KEYWORDS_PATTERN    = func() string { return `(` + orWords(KEYWORDS...) + `)` }
	FUNCTIONS_PATTERN   = func() string { return `\b(` + or(FUNCTIONS...) + `)` + lookAhead(`\s*`+quote(c.LPAREN)) }
	BLENDMODES_PATTERN  = func() string { return `(` + orWords(BLENDMODES...) + `)` }
	WORD_PATTERN        = `\b[a-zA-Z][a-zA-Z0-9\-_]*\b`
	NUMBER_PATTERN      = `\b-?\d+(\.\d+)?\b`
	STRING_PATTERN      = quote(c.QUOTE) + `[^` + quote(c.QUOTE) + `]*` + quote(c.QUOTE)
	ALPHA_PATTERN       = `\b(0|1)\.\d+?\b`
	FILE_PATTERN        = `\.\./.*|\./.*|/.*`
	URL_PATTERN         = `\b(http|ftp)s{0,1}://\S*\b`
	CLI_ARG_PATTERN     = `\$IMG\d+`
	SECTIONS_PATTERN    = `(` + orWords(SECTIONS...) + `)`
	COMPOSITION_PATTERN = `(` + orWords(COMPOSITION...) + `)`
	SECTION_PATTERN     = quote(c.LBRACKET) + SECTIONS_PATTERN + quote(c.RBRACKET)
	SOURCE_PATTERN      = `(` + FILE_PATTERN + `|` + URL_PATTERN + `|` + CLI_ARG_PATTERN + `)`
	PUNCTUATION_PATTERN = `(` + or(c.LPAREN, c.RPAREN, c.LBRACE, c.RBRACE, c.LBRACKET, c.RBRACKET) + `)`
)

// GenerateTmLanguage generates the complete TextMate grammar JSON for the GFXS language.
func GenerateTmLanguage() (string, error) {
	// Start with a pattern for double-quoted strings.
	patterns := []map[string]any{
		{
			"name":  "string.quoted.double." + c.LANGUAGE_ID,
			"begin": c.QUOTE,
			"end":   c.QUOTE,
			"patterns": []map[string]any{
				{
					"name":  "constant.character.escape." + c.LANGUAGE_ID,
					"match": `\\.`,
				},
			},
		},
	}

	// Create the base tmLanguage structure.
	tmLanguage := map[string]any{
		"name":      c.LANGUAGE_NAME,
		"fileTypes": []string{c.FILE_EXTENSION},
		"scopeName": "source." + c.LANGUAGE_ID,
		"patterns":  patterns,
	}

	// Helper function to add a pattern.
	fnAddPattern := func(scope, match string) {
		patterns = append(patterns, map[string]any{
			"name":  scope + "." + c.LANGUAGE_ID,
			"match": match,
		})
	}

	// Add additional patterns using real regexes from
	fnAddPattern("comment.line", c.COMMENT+`.*`+lookAhead(`$|\r\n|\n`))
	fnAddPattern("string.regexp", lookBehind(`^|\s*`)+BLENDMODES_PATTERN()+`\s+`+lookAhead(ALPHA_PATTERN))
	fnAddPattern("support.variable",
		lookBehind(ALPHA_PATTERN)+`\s+`+WORD_PATTERN+`\s+`+lookAhead(SOURCE_PATTERN)+`|`+ // word between alpha value and source (ie filter name)
			lookBehind(c.COMP_FILTER+`\s*`+c.ASSIGN+`\s*`)+WORD_PATTERN+`\s*|`+ // filter = name
			`\b`+WORD_PATTERN+`\b\s*`+lookAhead(quote(c.LBRACE))+`|`+ // word followed by left brace (ie function name)
			`\*`+
			``,
	)
	fnAddPattern("entity.name.class", SECTION_PATTERN)
	fnAddPattern("entity.name.type", SOURCE_PATTERN)
	fnAddPattern("constant.numeric", NUMBER_PATTERN)
	fnAddPattern("constant.string", STRING_PATTERN)
	fnAddPattern("keyword.include", KEYWORDS_PATTERN())
	fnAddPattern("keyword.other", COMPOSITION_PATTERN+lookAhead(`\s*`+c.ASSIGN))
	fnAddPattern("support.function", FUNCTIONS_PATTERN())
	fnAddPattern("punctuation", PUNCTUATION_PATTERN)
	fnAddPattern("markup.bold", `\b`+WORD_PATTERN+`\b\s*`+lookAhead(c.ASSIGN))

	// Update the patterns in the tmLanguage.
	tmLanguage["patterns"] = patterns

	// Marshal the grammar to JSON.
	b, err := json.MarshalIndent(tmLanguage, "", "    ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var quote = regexp.QuoteMeta

func or(strs ...string) string {
	for i, s := range strs {
		strs[i] = regexp.QuoteMeta(s)
	}
	return strings.Join(strs, "|")
}

func orWords(strs ...string) string {
	for i, s := range strs {
		strs[i] = `\b` + regexp.QuoteMeta(s) + `\b`
	}
	return strings.Join(strs, "|")
}

func lookBehind(str string) string { return `(?<=` + str + `)` }
func lookAhead(str string) string  { return `(?=` + str + `)` }
