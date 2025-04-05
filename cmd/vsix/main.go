package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/config/constants"
	"github.com/toxyl/gfx/fs"

	"github.com/toxyl/gfx/utils/buildlog"
)

// generateVSIX creates the VS Code extension package (.vsix) for GFXS and returns its path.
func generateVSIX() (string, error) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "vscode-extension-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	// Create required subdirectories.
	syntaxesDir := filepath.Join(tempDir, "syntaxes")
	if err := os.Mkdir(syntaxesDir, 0750); err != nil {
		return "", err
	}

	// Generate language-configuration.json.
	langConfig := map[string]interface{}{
		"comments": map[string]string{
			"lineComment": constants.COMMENT,
		},
		"wordPattern": config.WORD_PATTERN,
		"brackets": [][2]string{
			{constants.LBRACE, constants.RBRACE},
			{constants.LPAREN, constants.RPAREN},
		},
		"autoClosingPairs": []map[string]string{
			{"open": constants.LBRACKET, "close": constants.RBRACKET},
			{"open": constants.LBRACE, "close": constants.RBRACE},
			{"open": constants.LPAREN, "close": constants.RPAREN},
		},
	}
	langConfigJSON, err := json.MarshalIndent(langConfig, "", "    ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(tempDir, "language-configuration.json"), langConfigJSON, 0600); err != nil {
		return "", err
	}

	// Generate tmLanguage JSON using common.GenerateTmLanguage.
	tmLanguageJSON, err := config.GenerateTmLanguage()
	if err != nil {
		return "", err
	}
	tmLanguagePath := filepath.Join(syntaxesDir, config.PATH_TM_LANG_NAME)
	if err := os.WriteFile(tmLanguagePath, []byte(tmLanguageJSON), 0600); err != nil {
		return "", err
	}
	fs.Copy(tmLanguagePath, config.PATH_TM_LANG)

	// Generate package.json.
	packageJSON := map[string]any{
		"name":        constants.PACKAGE_NAME,
		"displayName": constants.DISPLAY_NAME,
		"description": constants.DESCRIPTION,
		"version":     constants.VERSION,
		"engines": map[string]string{
			"vscode": "^1.81.0",
		},
		"categories": []string{"Programming Languages"},
		"contributes": map[string]any{
			"languages": []map[string]any{
				{
					"id":            constants.LANGUAGE_ID,
					"aliases":       []string{constants.LANGUAGE_NAME, constants.LANGUAGE_ID},
					"extensions":    []string{constants.FILE_EXTENSION},
					"configuration": "./language-configuration.json",
				},
			},
			"grammars": []map[string]any{
				{
					"language":  constants.LANGUAGE_ID,
					"scopeName": "source." + constants.LANGUAGE_ID,
					"path":      "./syntaxes/" + config.PATH_TM_LANG_NAME,
				},
			},
		},
	}
	packageJSONBytes, err := json.MarshalIndent(packageJSON, "", "    ")
	if err != nil {
		return "", err
	}
	if err := fs.SaveString(filepath.Join(tempDir, "package.json"), string(packageJSONBytes)); err != nil {
		return "", err
	}

	// Generate README.md.
	if err := fs.SaveString(filepath.Join(tempDir, "README.md"), constants.README); err != nil {
		return "", err
	}

	// Package the extension into a VSIX (.vsix) file.
	vsixPath := filepath.Join(tempDir, "..", constants.LANGUAGE_ID+".vsix")
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

// addFilesToZip recursively adds files from root to the zip archive with the specified prefix.
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

func main() {
	buildlog.Log("Building VSIX",
		func() {
			vsixPath, err := generateVSIX()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating VSIX: %v\n", err)
				os.Exit(1)
			}
			err = fs.Copy(vsixPath, config.PATH_VSIX)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error copying VSIX: %v\n", err)
				os.Exit(1)
			}

		},
	)
}
