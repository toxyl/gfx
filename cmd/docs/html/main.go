package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/fs"
	"github.com/toxyl/gfx/utils/buildlog"
	"github.com/toxyl/gfx/utils/markdown"
)

func build(slug, file string) string {
	html, err := markdown.ToHTML(fs.LoadString(file))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating documentation (HTML): %v\n", err)
		os.Exit(1)
	}

	sb := strings.Builder{}
	sb.WriteString(`<a name="`)
	sb.WriteString(slug)
	sb.WriteString(`"></a>`)
	sb.WriteString(html)
	sb.WriteString(`<br>`)
	return sb.String()
}

func main() {
	buildlog.Log("Building HTML docs",
		func() {
			sb := strings.Builder{}
			sb.WriteString(`<!-- auto-generated -->
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>GFXScript Docs</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-okaidia.min.css">
	<style>`)
			sb.WriteString(fs.LoadString(config.PATH_CSS))
			sb.WriteString(`</style>
</head>
<body>
	<div class="menu">
		<button onclick="window.location.href='#filters'">Filters</button>
		<button onclick="window.location.href='#colormodels'">Color Models</button>
		<button onclick="window.location.href='#blendmodes'">Blend Modes</button>
		<button onclick="window.location.href='#projections'">Projections</button>
		<button onclick="window.location.href='#example'">Example Script</button>
	</div>
	<div class="markdown-body">`)
			sb.WriteString(build("filters", config.PATH_MARKDOWN+"/filters.md"))
			sb.WriteString(build("colormodels", config.PATH_MARKDOWN+"/colormodels.md"))
			sb.WriteString(build("blendmodes", config.PATH_MARKDOWN+"/blendmodes.md"))
			sb.WriteString(build("projections", config.PATH_MARKDOWN+"/projections.md"))
			sb.WriteString(`
		<a name="example"></a><h1>Example GFXScript</h1>
		<pre><code class="language-gfxs">` + fs.LoadString(config.PATH_GFXS_EXAMPLE) + `</code></pre>
	</div>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/prism.min.js"></script>
	<script>` + fs.LoadString(config.PATH_JS) + `</script>
</body>
</html>`)
			if err := fs.SaveString(config.PATH_HTML+"/index.html", sb.String()); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating documentation (HTML): %v\n", err)
				os.Exit(1)
			}
		},
	)
}
