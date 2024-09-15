// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func AdminPage(isAdmin bool) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<main><h1>Admin</h1><div><h2>Anmelden</h2><form hx-post=\"/auth/login\" hx-indicator=\"#loading\" hx-target=\"#content\"><label for=\"password\">Passwort:</label> <input id=\"password\" placeholder=\"Passwort\" type=\"password\" name=\"password\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isAdmin {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" disabled=\"true\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> <input type=\"submit\" value=\"Anmelden\" name=\"submit\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isAdmin {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" class=\"button-disabled\" disabled=\"true\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("></form></div><div><h2>Passwort ändern</h2><form hx-put=\"/auth/password\" hx-indicator=\"#loading\" hx-target=\"#content\"><label for=\"old-password\">Altes Passwort:</label> <input id=\"old-password\" placeholder=\"Altes Passwort\" type=\"password\" name=\"oldPassword\"> <label for=\"new-password\">Neues Passwort:</label> <input id=\"new-password\" placeholder=\"Neues Passwort\" type=\"password\" name=\"newPassword\"> <input type=\"submit\" value=\"Bestätigen\" name=\"submit\"></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isAdmin {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"logout-section\"><button class=\"danger-button\" hx-post=\"/auth/logout\" hx-trigger=\"click\" hx-target=\"#content\">Abmelden</button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
