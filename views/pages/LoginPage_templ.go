// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.898
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "morethancoder/t3-clone/views/components"
import "context"
import "morethancoder/t3-clone/db"
import "os"

func GetOAuth2Providers(ctx context.Context) []db.OAuth2Provider {
	data, ok := ctx.Value("OAuth2Providers").([]db.OAuth2Provider)
	if !ok {
		return []db.OAuth2Provider{}
	} else {
		return data
	}
}

func LoginPage() templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<main id=\"view\"><div class=\"h-full w-full min-h-screen flex items-center justify-center\"><div class=\"card z-10 w-full max-w-sm aspect-[3/2] mx-auto\nbackdrop-blur\nbg-base-100/50 shadow-md border-2 border-base-content/10\"><div class=\"card-body flex flex-col items-center justify-center gap-4 p-8\"><div class=\"avatar\"></div><div class=\"avatar avatar-placeholder\"><div class=\"bg-neutral text-neutral-content w-10 rounded-box\"><span class=\"\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.IconLockWithHole("size-5").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</span></div></div><h2 class=\"card-title text-3xl\">Welcome Back</h2><p class=\"text-sm opacity-50\">Sign in to continue to your account</p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, provider := range GetOAuth2Providers(ctx) {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 templ.SafeURL
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinURLErrs(templ.SafeURL(provider.AuthURL + os.Getenv("APP_URL") + "/auth-redirect"))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/pages/LoginPage.templ`, Line: 37, Col: 92}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\" class=\"btn btn-wide bg-neutral/15 hover:bg-neutral/30 text-neutral\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			switch provider.Name {
			case "google":
				templ_7745c5c3_Err = components.LogoGoogleColored("size-5").Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			case "github":
				templ_7745c5c3_Err = components.LogoGithub("size-5").Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<span>Sign in with ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(provider.DisplayName)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/pages/LoginPage.templ`, Line: 46, Col: 56}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</span></a>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "<div class=\"divider my-1\"></div><p class=\"text-xs opacity-80 max-w-72\">By continuing, you agree to our <a class=\"link link-primary\">Terms of Service</a></p></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.BackgroundBlobs().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "</div></main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
