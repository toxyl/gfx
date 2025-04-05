package main

import (
	"fmt"
	"os"

	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/config/constants"
	"github.com/toxyl/gfx/fs"
	"github.com/toxyl/gfx/utils/buildlog"
)

func main() {
	buildlog.Log("Building JS",
		func() {
			// Generate the tmLanguage JSON for GFXS.
			grammarJSON, err := config.GenerateTmLanguage()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating grammar: %v\n", err)
				os.Exit(1)
			}
			if err := fs.SaveString(config.PATH_JS, fmt.Sprintf(`(function(){
                // Dynamically generated tmLanguage grammar for GFXScript.
                var tmGrammar = %s;
            
                // A simple conversion function from a tmLanguage grammar to a Prism language definition.
                // This only handles patterns with "match" and "name" properties.
                function convertTmToPrism(grammar) {
                    var prismLang = {};
                    if (grammar.patterns && grammar.patterns.length) {
                        grammar.patterns.forEach(function(pattern) {
                            if (pattern.match && pattern.name) {
                                // Use the first segment of the scope name as the Prism token key.
                                var tokenKey = pattern.name.replaceAll(".gfxscript", "").replaceAll(".", " ");
                                if (!prismLang[tokenKey]) {
                                    prismLang[tokenKey] = [];
                                }
                                // Create a global RegExp from the pattern.
                                prismLang[tokenKey].push(new RegExp(pattern.match, "g"));
                            }
                        });
                    }
                    // Combine multiple regexes for the same token if necessary.
                    for (var key in prismLang) {
                        if (prismLang[key].length === 1) {
                            prismLang[key] = prismLang[key][0];
                        } else {
                            var combined = prismLang[key].map(function(re) { return re.source; }).join("|");
                            prismLang[key] = new RegExp(combined, "g");
                        }
                    }
                    return prismLang;
                }
            
                // Convert the tmLanguage grammar to a Prism language definition.
                var prismDefinition = convertTmToPrism(tmGrammar);
                // Register the language with Prism using your language ID.
                Prism.languages.%s = prismDefinition;
            })();`, grammarJSON, constants.LANGUAGE_ID)); err != nil {
				fmt.Fprintf(os.Stderr, "Error generating documentation (JS): %v\n", err)
				os.Exit(1)
			}
		},
	)
}
