package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/toxyl/gfx/config"
	"github.com/toxyl/gfx/core/blendmodes"
	"github.com/toxyl/gfx/core/color"
	"github.com/toxyl/gfx/core/fx"
	"github.com/toxyl/gfx/core/projections"
	"github.com/toxyl/gfx/fs"
	"github.com/toxyl/gfx/utils/buildlog"
)

func build(path, title, prefix, suffix string, content func() string) {
	buildlog.Log("Building Markdown docs ("+title+")",
		func() {
			prefix = strings.TrimSpace(prefix)
			suffix = strings.TrimSpace(suffix)
			sb := strings.Builder{}
			sb.WriteString("# ")
			sb.WriteString(title)
			sb.WriteString("\n")
			if prefix != "" {
				sb.WriteString(prefix)
				sb.WriteString("\n")
			}
			if content != nil {
				sb.WriteString(content())
				sb.WriteString("\n")
			}
			if suffix != "" {
				sb.WriteString(suffix)
				sb.WriteString("\n")
			}
			if err := fs.SaveString(path, sb.String()); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating documentation (Markdown): %v\n", err)
				os.Exit(1)
			}
		})
}

func formatFunctionDoc(meta *fx.FunctionMeta) string {
	sb := strings.Builder{}
	sb.WriteString("## ")
	sb.WriteString(meta.Name)
	sb.WriteString("\n")
	if meta.Description != "" {
		sb.WriteString(meta.Description)
		sb.WriteString("\n\n")
	}
	if len(meta.Args) > 0 {
		sb.WriteString("```gfxs\n")
		sb.WriteString(meta.Name)
		sb.WriteString("(")
		for i, arg := range meta.Args {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(arg.Name)
			sb.WriteString("=")
			if arg.Default != nil {
				sb.WriteString(fmt.Sprintf("%v", arg.Default))
			}
		}
		sb.WriteString(")\n```\n\n")
		sb.WriteString("| Argument | Type | Unit | Default | Min | Max | Description |\n")
		sb.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
		for _, arg := range meta.Args {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %v | %v | %v | %s |\n",
				arg.Name, arg.Type, "", arg.Default, arg.Min, arg.Max, arg.Description))
		}
	} else {
		sb.WriteString("```gfxs\n")
		sb.WriteString(meta.Name)
		sb.WriteString("()\n```\n")
	}
	return sb.String()
}

func main() {
	build(config.PATH_MARKDOWN+"/filters.md", "Filters", "", "", func() string {
		list := []string{}
		for _, name := range fx.DefaultRegistry.List() {
			meta := fx.DefaultRegistry.GetMeta(name)
			if meta != nil {
				list = append(list, formatFunctionDoc(meta))
			}
		}
		sort.Strings(list)
		return strings.Join(list, "\n\n") + "\n"
	})
	build(config.PATH_MARKDOWN+"/colormodels.md", "Color Models", "", "", func() string {
		list := []string{}
		for _, name := range color.DefaultColorModelRegistry.List() {
			model, _ := color.DefaultColorModelRegistry.Get(name)
			if model != nil {
				list = append(list, model.Meta.Doc())
			}
		}
		sort.Strings(list)
		return strings.Join(list, "\n\n") + "\n"
	})
	build(config.PATH_MARKDOWN+"/blendmodes.md", "Blend Modes", "| Name | Description |\n| --- | --- |\n", "", func() string {
		list := []string{}
		for _, name := range blendmodes.List() {
			mode, err := blendmodes.Get(name)
			if err == nil && mode != nil {
				list = append(list, mode.Meta().Doc())
			}
		}
		sort.Strings(list)
		return strings.Join(list, "\n") + "\n"
	})
	build(config.PATH_MARKDOWN+"/projections.md", "Projections", "", "", func() string {
		list := []string{}
		for _, p := range projections.Default.List() {
			if p != nil && p.Meta != nil {
				list = append(list, p.Meta.Doc())
			}
		}
		sort.Strings(list)
		return strings.Join(list, "\n") + "\n"
	})
}
