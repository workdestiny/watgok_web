package app

import "html/template"

// TemplateFuncs returns template funcs
func (app *App) TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"static": func(s string) string {
			p := app.Static[s]
			return "/-" + p
		},
	}
}
