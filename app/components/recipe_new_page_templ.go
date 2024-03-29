// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.598
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func RecipeNewPage() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = Header().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<main><h2>Neues Rezept</h2><form hx-post=\"/recipe\" hx-target=\"#content\"><label for=\"title\">Titel:</label> <input id=\"title\" placeholder=\"Titel\" type=\"text\" name=\"title\"> <label for=\"description\">Beschreibung:</label> <input id=\"description\" placeholder=\"Beschreibung\" type=\"text\" name=\"description\"> <label for=\"duration\">Zubereitungszeit (Minuten):</label> <input id=\"duration\" placeholder=\"Zubereitungszeit (Minuten)\" type=\"number\" name=\"duration\"> <label for=\"tags\">Tags:</label> <input id=\"tags\" placeholder=\"Tags\" type=\"text\" name=\"tags\"> <label for=\"ingredient\">Zutaten:</label> <textarea id=\"ingredients\" placeholder=\"Zutaten\" name=\"ingredients\"></textarea> <label for=\"instructions\">Anleitung:</label> <textarea id=\"instructions\" placeholder=\"Anleitung\" name=\"instructions\"></textarea> <input type=\"submit\" value=\"Rezept erstellen\" name=\"submit\"></form></main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
