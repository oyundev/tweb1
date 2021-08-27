package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/azureadv2"
	gf "github.com/shareed2k/goth_fiber"
)

const (
	htmlheadsrc = `<!DOCTYPE html><html><head><title>Home</title> 
	<link rel="shortcut icon" type="image/x-icon" href="/favicon.ico" />
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8;">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>ul {list-style-type: none; margin: 0; padding: 0; overflow: hidden; background-color: #333333;} li {float: left;} li a {display: block; color: white; text-align: center; padding: 16px; text-decoration: none; font-family: Ubuntu; font-size: 16pt;} li:hover {background-color: #111111;} body {color: #7a7a7a;background-color: #202020;} table {background-color: #102022; border-collapse: collapse; margin: 25px 0; font-size: 0.9em; font-family: sans-serif; min-width: 400px; box-shadow: 0 0 20px rgba(0, 0, 0, 0.15); }table thead tr { background-color: #91bdbc; color: #000; text-align: left; }table th, td { padding: 12px 15px; text-align: left; word-break:break-all;}table tbody tr { border-bottom: 1px solid #c0c0c0; }table tbody tr:nth-of-type(even) { background-color: #000; }styled-table tbody tr:last-of-type { border-bottom: 1px solid #c0c0c0; } </style>
	</head>`
)

var (
	flagkey       = flag.String("k", "", "simple string parameter")
	secret        = os.Getenv("MICROSOFT_PROVIDER_AUTHENTICATION_SECRET")
	applicationID = os.Getenv("APPLICATION_ID") //"1052423-cccc-4444-2222-cccca4b907"
	redirectUri   = os.Getenv("REDIRECT_URI")   //"https://xxxxxxxxx.azurewebsites.net/auth/azureadv2/callback" debug:"https://127.0.0.1:3000/auth/azureadv2/callback"
	tenantID      = os.Getenv("AD_TENANT_ID")   //"0333033-40bc-4141-8888-1111111111"->"Directory (tenant) ID"
)

func dumpMap(m map[string]interface{}) (dump map[string]string) {
	dump = make(map[string]string)
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			for k2, v2 := range dumpMap(mv) {
				dump[k2] = v2
			}
		} else {
			if v != nil {
				tp := fmt.Sprintf("%T", v)
				if tp != "[]interface {}" {
					dump[k] = v.(string)
				}
			} else {
				dump[k] = ""
			}
		}
	}
	return dump
}

func main() {
	var testparam string
	flag.Parse()
	if len(strings.TrimSpace(*flagkey)) > 1 {
		testparam = strings.TrimSpace(*flagkey)
	}

	options := azureadv2.ProviderOptions{
		Scopes: []azureadv2.ScopeType{azureadv2.UserReadScope},
		Tenant: azureadv2.TenantType(tenantID), //azureadv2.CommonTenant
	}
	goth.UseProviders(
		azureadv2.New(applicationID, secret, redirectUri, options),
	)

	app := fiber.New(fiber.Config{
		ReadBufferSize: 16384,
	})

	app.Get("/auth/:provider/callback", func(ctx *fiber.Ctx) error {
		//app.Get("/", func(ctx *fiber.Ctx) error {

		///////////////////////////////////////////////////////////////
		//providerName := "azureadv2"
		//provider, err := goth.GetProvider(providerName)
		//if err != nil {
		//	return err
		//}
		//sess, err := provider.BeginAuth(gf.SetState(ctx))
		//if err != nil {
		//	return err
		//}

		////_url, err := sess.GetAuthURL()
		////if err != nil {
		////	return err
		////}

		//err = gf.StoreInSession(providerName, sess.Marshal(), ctx)

		//if err != nil {
		//	return err
		//}
		///////////////////////////////////////////////////////////////

		/////////// querystring icinden aliyormus http://127.0.0.1/?provider=azureadv2
		// ctx.Context().Request.SetRequestURI("/?provider=azureadv2")

		user, err := gf.CompleteUserAuth(ctx)
		if err != nil {
			return err
		}
		//ctx.JSON(user)

		//PROVIDER-BUG:https://github.com/markbates/goth/issues/289
		var userPrincipalName string
		for k, v := range dumpMap(user.RawData) {
			if k == "userPrincipalName" {
				userPrincipalName = v
			}
		}

		var b strings.Builder
		b.WriteString(htmlheadsrc)
		b.WriteString(`<body><ul><li><a href="/logout/azureadv2">Logout</a></li></ul><br/><table><thead><tr><th>Name</th><th>Value</th></tr></thead><tbody>`)
		if len(user.Email) < 2 {
			b.WriteString(`<tr><td nowrap>Email</td><td>` + userPrincipalName + `</td></tr>`)
		} else {
			b.WriteString(`<tr><td nowrap>Email</td><td>` + user.Email + `</td></tr>`)
		}
		b.WriteString(`<tr><td nowrap>Name</td><td>` + user.Name + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>FirstName</td><td>` + user.FirstName + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>LastName</td><td>` + user.LastName + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>NickName</td><td>` + user.NickName + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>Description</td><td>` + user.Description + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>UserID</td><td>` + user.UserID + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>AvatarURL</td><td>` + user.AvatarURL + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>Location</td><td>` + user.Location + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>AccessToken</td><td>` + user.AccessToken + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>AccessTokenSecret</td><td>` + user.AccessTokenSecret + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>RefreshToken</td><td>` + user.RefreshToken + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>ExpiresAt</td><td>` + user.ExpiresAt.String() + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>IDToken</td><td>` + user.IDToken + `</td></tr>`)
		b.WriteString(`<tr><td nowrap>userPrincipalName</td><td>` + user.RawData["userPrincipalName"].(string) + `</td></tr>`)
		b.WriteString(`</tbody></table><br/>`)
		b.WriteString(`<p>` + os.Getenv("WEBSITE_HOSTNAME") + `</p>`)
		b.WriteString(`<p>` + testparam + `</p>`)
		b.WriteString(`</body></html>`)

		ctx.Set("Content-Type", "text/html")
		ctx.Send([]byte(b.String()))
		return nil
	})

	app.Get("/logout/:provider", func(ctx *fiber.Ctx) error {
		gf.Logout(ctx)
		ctx.Redirect("/")
		return nil
	})

	app.Get("/auth/:provider", func(ctx *fiber.Ctx) error {
		if authUser, err := gf.CompleteUserAuth(ctx); err == nil {
			ctx.JSON(authUser)
		} else {
			gf.BeginAuthHandler(ctx)
		}
		return nil
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		var b strings.Builder
		b.WriteString(htmlheadsrc)
		b.WriteString(`<body><ul><li><a href="/auth/azureadv2">Login</a></li><li><a>Logout</a></li></ul>`)
		b.WriteString(`<p>` + os.Getenv("WEBSITE_HOSTNAME") + `</p>`)
		b.WriteString(`<p>` + testparam + `</p>`)
		b.WriteString(`</body></html>`)

		ctx.Set("Content-Type", "text/html")
		ctx.Send([]byte(b.String()))
		return nil
	})

	log.Fatal(app.Listen(":80"))
	//log.Fatal(app.ListenTLS(":3000", "./localhost+2.pem", "./localhost+2-key.pem"))
	//app.ListenTLS(":8080", "./cert.pem", "./cert.key")
}
