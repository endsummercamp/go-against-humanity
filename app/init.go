package app

import (
	"github.com/ESCah/go-against-humanity/app/models"
	gc_log "github.com/denysvitali/gc_log"
	"github.com/revel/revel"
	"math/rand"
	"strings"
	"time"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		// revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.TemplateFuncs["replace"] = func(input, from, to string) string {
		return strings.Replace(input, from, to, -1)
	}

	revel.TemplateFuncs["card_text"] = func(input models.Card) string {
		return input.GetText()
	}

	revel.TemplateFuncs["card_dash"] = func(input models.Card) string {
		return strings.Replace(input.GetText(), "_", "<div class=\"long-dash\"></div>", -1)
	}

	revel.TemplateFuncs["card_black"] = func(input models.Card) bool {
		return input.GetColor() == models.BLACK_CARD
	}

	revel.TemplateFuncs["long_text"] = func(input string) bool {
		return len(input) > 100
	}

	revel.TemplateFuncs["is_player"] = func(user models.User) bool {
		return user.UserType == models.PlayerType
	}

	revel.TemplateFuncs["is_admin"] = func(user models.User) bool {
		return user.IsAdmin()
	}

	revel.TemplateFuncs["format_date"] = func(date time.Time) string {
		return date.Format("2 Jan 2006, 15:04:01")
	}

	rand.Seed(time.Now().Unix())

	gc_log.SetDebug(true)

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
