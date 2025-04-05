(function(){
                // Dynamically generated tmLanguage grammar for GFXScript.
                var tmGrammar = {
    "fileTypes": [
        ".gfxs"
    ],
    "name": "GFXScript",
    "patterns": [
        {
            "begin": "\"",
            "end": "\"",
            "name": "string.quoted.double.gfxs",
            "patterns": [
                {
                    "match": "\\\\.",
                    "name": "constant.character.escape.gfxs"
                }
            ]
        },
        {
            "match": "#.*(?=$|\\r\\n|\\n)",
            "name": "comment.line.gfxs"
        },
        {
            "match": "(?\u003c=^|\\s*)(\\badd\\b|\\baverage\\b|\\bcolor\\b|\\bcolor-burn\\b|\\bcolor-dodge\\b|\\bcontrast-negate\\b|\\bdarken\\b|\\bdarker-color\\b|\\bdifference\\b|\\bdivide\\b|\\berase\\b|\\bexclusion\\b|\\bglow\\b|\\bhard-light\\b|\\bhard-mix\\b|\\bhue\\b|\\blighten\\b|\\blighter-color\\b|\\blinear-burn\\b|\\blinear-light\\b|\\bluminosity\\b|\\bmultiply\\b|\\bnegation\\b|\\bnormal\\b|\\boverlay\\b|\\bpin-light\\b|\\breflect\\b|\\bsaturation\\b|\\bscreen\\b|\\bsoft-light\\b|\\bsubtract\\b|\\bvivid-light\\b)\\s+(?=\\b(0|1)\\.\\d+?\\b)",
            "name": "string.regexp.gfxs"
        },
        {
            "match": "(?\u003c=\\b(0|1)\\.\\d+?\\b)\\s+\\b[a-zA-Z][a-zA-Z0-9\\-_]*\\b\\s+(?=(\\.\\./.*|\\./.*|/.*|\\b(http|ftp)s{0,1}://\\S*\\b|\\$IMG\\d+))|(?\u003c=filter\\s*=\\s*)\\b[a-zA-Z][a-zA-Z0-9\\-_]*\\b\\s*|\\b\\b[a-zA-Z][a-zA-Z0-9\\-_]*\\b\\b\\s*(?=\\{)|\\*",
            "name": "support.variable.gfxs"
        },
        {
            "match": "\\[(\\bVARS\\b|\\bFILTERS\\b|\\bCOMPOSITION\\b|\\bLAYERS\\b)\\]",
            "name": "entity.name.class.gfxs"
        },
        {
            "match": "(\\.\\./.*|\\./.*|/.*|\\b(http|ftp)s{0,1}://\\S*\\b|\\$IMG\\d+)",
            "name": "entity.name.type.gfxs"
        },
        {
            "match": "\\b-?\\d+(\\.\\d+)?\\b",
            "name": "constant.numeric.gfxs"
        },
        {
            "match": "\"[^\"]*\"",
            "name": "constant.string.gfxs"
        },
        {
            "match": "(\\buse\\b|\\bhsla\\b)",
            "name": "keyword.include.gfxs"
        },
        {
            "match": "(\\bwidth\\b|\\bheight\\b|\\bcolor\\b|\\bcrop\\b|\\bfilter\\b|\\bresize\\b)(?=\\s*=)",
            "name": "keyword.other.gfxs"
        },
        {
            "match": "\\b(alpha-map|brightness|color-shift|contrast|convolution|blur|crop|crop-circle|edge-detect|emboss|enhance|extract|flip-h|flip-v|gamma|gray|hue|hue-contrast|invert|lum|lum-contrast|pastelize|remove|rotate|sat|sat-contrast|scale|sepia|sharpen|threshold|to-polar|transform|translate|translate-wrap|vibrance|project-as-equirectangular|project-as-mercator|project-as-polar-to-rect|project-as-rect-to-polar|project-as-sinusoidal|project-as-stereographic)(?=\\s*\\()",
            "name": "support.function.gfxs"
        },
        {
            "match": "(\\(|\\)|\\{|\\}|\\[|\\])",
            "name": "punctuation.gfxs"
        },
        {
            "match": "\\b\\b[a-zA-Z][a-zA-Z0-9\\-_]*\\b\\b\\s*(?==)",
            "name": "markup.bold.gfxs"
        }
    ],
    "scopeName": "source.gfxs"
};
            
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
                Prism.languages.gfxs = prismDefinition;
            })();