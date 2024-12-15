package parser

import (
	"archive/zip"
	"encoding/json"
	"os"
	"path/filepath"
)

// GenerateVSIX creates the VSCode extension package and returns the path to the .vsix file
func GenerateVSIX() (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "vscode-extension-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories
	syntaxesDir := filepath.Join(tempDir, "syntaxes")
	if err := os.Mkdir(syntaxesDir, 0755); err != nil {
		return "", err
	}

	// Generate language-configuration.json
	langConfig := map[string]interface{}{
		"comments": map[string]string{
			"lineComment": STR_COMMENT,
		},
		"wordPattern": WORD_PATTERN,
		"brackets": [][2]string{
			{STR_LBRACE, STR_RBRACE},
			{STR_LPAREN, STR_RPAREN},
		},
		"autoClosingPairs": []map[string]string{
			{"open": STR_LBRACKET, "close": STR_RBRACKET},
			{"open": STR_LBRACE, "close": STR_RBRACE},
			{"open": STR_LPAREN, "close": STR_RPAREN},
		},
	}
	langConfigJSON, err := json.MarshalIndent(langConfig, "", "    ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(tempDir, "language-configuration.json"), langConfigJSON, 0644); err != nil {
		return "", err
	}

	// Generate *.tmLanguage.json
	patterns := []map[string]any{{ // strings
		"name":  "string.quoted." + LANGUAGE_ID,
		"begin": STR_QUOTE,
		"end":   STR_QUOTE,
		"patterns": []map[string]interface{}{
			{
				"name":  "constant.character.escape." + LANGUAGE_ID,
				"match": `\\.`,
			},
		},
	}}
	tmLanguage := map[string]interface{}{
		"name":      LANGUAGE_NAME,
		"fileTypes": []string{FILE_EXTENSION},
		"scopeName": "source." + LANGUAGE_ID,
		"patterns":  &patterns,
	}
	fnAddPattern := func(scope, match string) {
		patterns = append(patterns, map[string]interface{}{ // keywords
			"name":  scope + "." + LANGUAGE_ID,
			"match": match,
		})
	}

	// keywords
	fnAddPattern("keyword.include", KEYWORDS_PATTERN)
	// blendmodes
	fnAddPattern("string.regexp", BLENDMODES_PATTERN)
	// filepaths
	fnAddPattern("entity.name.type", SOURCE_PATTERN)
	// numbers
	fnAddPattern("constant.numeric", NUMBER_PATTERN)
	// composition settings
	fnAddPattern("keyword.other", COMPOSITION_PATTERN+`(?=\s*`+STR_ASSIGN+`)`)
	// layer operations
	fnAddPattern("keyword.other", LAYER_PATTERN+`(?=\s+\d+)`)
	// functions
	fnAddPattern("support.function", WORD_PATTERN+`\s*\`+STR_LPAREN)
	// sections
	fnAddPattern("entity.name.class", SECTION_PATTERN)
	// filter names
	fnAddPattern("support.variable", `\*|(?<=`+COMP_FILTER+`\s*`+STR_ASSIGN+`\s*)`+WORD_PATTERN+`\s*|`+WORD_PATTERN+`\s*(?=\`+STR_LBRACE+`)|(?<=\s*`+BLENDMODES_PATTERN+`\s+`+ALPHA_PATTERN+`\s+)`+WORD_PATTERN+`\s+`)
	// arg names
	fnAddPattern("markup.bold", `(?<=[^\s]+)\s*`+WORD_PATTERN+`\s*(?=`+STR_ASSIGN+`)`)
	// comments
	fnAddPattern("comment.line", STR_COMMENT+".*$")

	tmLanguageJSON, err := json.MarshalIndent(tmLanguage, "", "    ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(syntaxesDir, LANGUAGE_ID+".tmLanguage.json"), tmLanguageJSON, 0644); err != nil {
		return "", err
	}

	// Generate package.json
	packageJSON := map[string]interface{}{
		"name":        PACKAGE_NAME,
		"displayName": DISPLAY_NAME,
		"description": DESCRIPTION,
		"version":     VERSION,
		"engines": map[string]string{
			"vscode": "^1.81.0",
		},
		"categories": []string{"Programming Languages"},
		"contributes": map[string]interface{}{
			"languages": []map[string]interface{}{
				{
					"id":            LANGUAGE_ID,
					"aliases":       []string{LANGUAGE_NAME, LANGUAGE_ID},
					"extensions":    []string{FILE_EXTENSION},
					"configuration": "./language-configuration.json",
				},
			},
			"grammars": []map[string]interface{}{
				{
					"language":  LANGUAGE_ID,
					"scopeName": "source." + LANGUAGE_ID,
					"path":      "./syntaxes/" + LANGUAGE_ID + ".tmLanguage.json",
				},
			},
		},
	}
	packageJSONBytes, err := json.MarshalIndent(packageJSON, "", "    ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(tempDir, "package.json"), packageJSONBytes, 0644); err != nil {
		return "", err
	}

	// Generate README.md
	if err := os.WriteFile(filepath.Join(tempDir, "README.md"), []byte(README), 0644); err != nil {
		return "", err
	}

	// Package the extension into a .vsix file
	vsixPath := filepath.Join(tempDir, "..", LANGUAGE_ID+".vsix")
	zipFile, err := os.Create(vsixPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if err := addFilesToZip(zipWriter, tempDir, "extension/"); err != nil {
		return "", err
	}

	return vsixPath, nil
}

// Helper function to add files to the ZIP archive
func addFilesToZip(zipWriter *zip.Writer, root, prefix string) error {
	files, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := filepath.Join(root, file.Name())
		if file.IsDir() {
			if err := addFilesToZip(zipWriter, filePath, filepath.Join(prefix, file.Name())); err != nil {
				return err
			}
		} else {
			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			zipEntry, err := zipWriter.Create(filepath.Join(prefix, file.Name()))
			if err != nil {
				return err
			}
			if _, err := zipEntry.Write(fileBytes); err != nil {
				return err
			}
		}
	}
	return nil
}
