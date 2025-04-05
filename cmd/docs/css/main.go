package main

import (
	"fmt"
	"os"

	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/fs"
	"github.com/toxyl/gfx/utils/buildlog"
)

func main() {
	buildlog.Log("Building CSS",
		func() {
			if err := fs.SaveString(config.PATH_CSS, `
body {
	margin: 0;
	padding: 0;
	background-color: #0d1117;
}

button {
    background: hsla(0, 0%, 20%, 0.8);
    border: none;
    color: #fff;
    font-size: 14px;
    padding: 5px 10px;
    margin: 5px;
	cursor: pointer;
}

button:hover {
    background: hsla(0, 0%, 30%, 0.8);
}

.menu {
    position: fixed;
    left: 0;
    right: 90vw;
    bottom: 0;
    top: 0;
    padding: 10px;
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    justify-content: center;
}

.markdown-body {
    position: fixed;
    top: 0;
    left: 20vw;
    right: 0;
    bottom: 0;
    box-sizing: border-box;
    margin: 0 auto;
    padding: 45px;
	padding-right: 20vw;
    color: #c9d1d9;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    line-height: 1.6;
    overflow-y: scroll;
}

/* Headings */
.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4,
.markdown-body h5,
.markdown-body h6 {
	margin-top: 1.5em;
	margin-bottom: 0.5em;
	font-weight: 600;
	line-height: 1.25;
}

.markdown-body h1 { font-size: 2em; }
.markdown-body h2 { font-size: 1.75em; }
.markdown-body h3 { font-size: 1.5em; }
.markdown-body h4 { font-size: 1.25em; }
.markdown-body h5 { font-size: 1em; }
.markdown-body h6 { font-size: 0.875em; }

/* Paragraphs and Links */
.markdown-body p {
	margin: 0.8em 0;
}

.markdown-body a {
	color: #58a6ff;
	text-decoration: none;
}

.markdown-body a:hover {
	text-decoration: underline;
}

/* Blockquotes */
.markdown-body blockquote {
	padding: 0 1em;
	color: #8b949e;
	border-left: 0.25em solid #30363d;
	margin: 1em 0;
}

/* Code (inline and blocks) */
.markdown-body code {
	background-color: #161b22;
	color: #c9d1d9;
	padding: 0.2em 0.4em;
	border-radius: 3px;
	font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace;
	font-size: 85%;
}

.markdown-body pre {
	background-color: #161b22;
	padding: 1em;
	overflow: auto;
	border-radius: 6px;
	margin: 1.5em 0;
}

.markdown-body pre code {
	background-color: transparent;
	padding: 0;
	border: 0;
	font-size: 100%;
}

/* Lists */
.markdown-body ul,
.markdown-body ol {
	padding-left: 2em;
	margin: 1em 0;
}

.markdown-body li {
	margin-bottom: 0.25em;
}

/* Horizontal Rules */
.markdown-body hr {
	border: 0;
	border-top: 1px solid #30363d;
	margin: 2em 0;
}

/* Tables */
.markdown-body table {
	width: 100%;
	border-collapse: collapse;
	margin: 1em 0;
}

.markdown-body th,
.markdown-body td {
	border: 1px solid #30363d;
	padding: 0.6em 1em;
}

.markdown-body th {
	background-color: #161b22;
}

/* Images */
.markdown-body img {
	max-width: 100%;
	display: block;
	margin: 0.5em 0;
	border-radius: 5px;
}

/* Task Lists (checkboxes) */
.markdown-body input[type="checkbox"] {
	margin: 0 .2em .25em -1.6em;
	vertical-align: middle;
}
	
code[class*="language-gfxs"],
pre[class*="language-gfxs"] {
	color: #d4d4d4;
	background: #1e1e1e;
	padding: 0.5em;
	overflow: auto;
	border-radius: 0.3em;
	font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
	font-size: 1em;
	line-height: 1.5em;
}

/* token styles */
.token.constant.string {
	color: #CE9178;
}
.token.string.quoted.double {
	color: #CE9178;
}
.token.keyword {
	color: #569CD6;
}
.token.constant.number {
	color: #B5CEA8;
}
.token.constant.numeric {
	color: #B5CEA8;
}
.token.support.function {
	color: #DCDCAA;
}
.token.support.variable {
	color: #9cdcfe;
}
.token.comment.line {
	color: #6a9955;
}
.token.entity.name.type {
	color:#37A98B;
}
.token.entity.name.class {
	color:#4EC9B0;
}
.token.string.regexp {
	color:#D16969;
}
.token.markup.bold {
	color:#569CD6;
}
.token.punctuation {
	color:#C586C0;
}
`); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating documentation (CSS): %v\n", err)
				os.Exit(1)
			}
		},
	)
}
