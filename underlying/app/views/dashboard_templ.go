// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/nefarius/cornelian/underlying/app"

func Dashboard(person app.Person, questions []app.Question) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container\"><div class=\"user-info\"><strong>Logged in as:</strong> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(person.Email)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `underlying/app/views/dashboard.templ`, Line: 8, Col: 49}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <a href=\"/logout\" hx-get=\"/logout\" hx-target=\"html\">Log out</a></div><div class=\"jumbotron\"><h2 class=\"text-success text-center\"><i class=\"bi bi-graph-up\"></i> Вопросы</h2><div class=\"button-group\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if person.Admin {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-info\" hx-get=\"/all\" hx-trigger=\"click\" hx-target=\"#questions\">Все вопросы <span hx-get=\"/countall\" hx-trigger=\"every 5s\" hx-target=\"this\"><span hx-get=\"/countall\" hx-trigger=\"load\" hx-target=\"this\"></span></span></button> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-secondary\" hx-get=\"/mine\" hx-trigger=\"click\" hx-target=\"#questions\">Мои вопросы <span hx-get=\"/countmine\" hx-trigger=\"every 5s\" hx-target=\"this\"><span hx-get=\"/countmine\" hx-trigger=\"load\" hx-target=\"this\"></span></span></button> <button class=\"btn btn-secondary\" hx-get=\"/current-quiz\" hx-trigger=\"click\" hx-target=\"body\">Квиз <span hx-get=\"/countmine\" hx-trigger=\"every 5s\" hx-target=\"this\"><span hx-get=\"/countmine\" hx-trigger=\"load\" hx-target=\"this\"></span></span></button> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if person.Admin {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn btn-primary\" hx-get=\"/quizzes-panel\" hx-trigger=\"click\" hx-target=\"body\">Квизы</button>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Questions(questions).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
