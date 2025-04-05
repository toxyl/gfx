package config

import (
	c "github.com/toxyl/gfx/config/constants"
)

const (
	PATH_TM_LANG_NAME       = c.LANGUAGE_ID + ".tmLanguage.json"
	PATH_CORE_MODULE_IMPORT = "github.com/toxyl/gfx/core"
	PATH_CORE_MODULE        = "./core"
	PATH_MARKDOWN           = "./docs/markdown"
	PATH_HTML               = "./docs/html"
	PATH_CSS                = "./docs/html/style.css"
	PATH_JS                 = "./docs/html/script.js"
	PATH_GFXS_EXAMPLE       = "./docs/gfxs-example.gfxs"
	PATH_VSIX               = "./docs/" + c.LANGUAGE_ID + "-syntax.vsix"
	PATH_VSIX_CODE_AI       = "./docs/vs-code-ai.vsix"
	PATH_TM_LANG            = "./docs/" + PATH_TM_LANG_NAME
)
