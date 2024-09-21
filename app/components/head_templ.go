// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func head() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<head><title>Let Him Cook</title><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><script src=\"/static/js/htmx.min.js\" defer></script><script src=\"/static/js/main.js\" defer></script><link rel=\"stylesheet\" href=\"/static/css/styles.css\"><link rel=\"stylesheet\" href=\"/static/css/fonts.css\"><link href=\"/static/fa/css/fontawesome.css\" rel=\"stylesheet\"><link href=\"/static/fa/css/solid.css\" rel=\"stylesheet\"><link rel=\"icon\" type=\"image/x-icon\" href=\"/static/favicon.ico\"><meta name=\"htmx-config\" content=\"{\n            &#34;responseHandling&#34;:[\n                {&#34;code&#34;:&#34;204&#34;, &#34;swap&#34;: false},\n                {&#34;code&#34;:&#34;[23]..&#34;, &#34;swap&#34;: true},\n                {&#34;code&#34;:&#34;[45]..&#34;, &#34;swap&#34;: true, &#34;error&#34;: true},\n                {&#34;code&#34;:&#34;...&#34;, &#34;swap&#34;: true}\n            ]}\"></head>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
