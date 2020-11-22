package actions

import (
	"net/http"
	"strconv"

	"github.com/lbryio/dispendium/util"

	"github.com/lbryio/lbry.go/extras/api"
	v "github.com/lbryio/ozzo-validation"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

func init() {
	// make validation always return json-style names for fields
	f := func(str string) string {
		return util.Underscore(str)
	}
	v.ErrorTagFunc = &f
}

//ConfigureAPIServer handles all of the configuration of the api server
func ConfigureAPIServer() {
	api.TraceEnabled = util.Debugging

	hs := make(map[string]string)

	hs["Server"] = "lbry.com"
	hs["Content-Type"] = "application/json; charset=utf-8"

	hs["Access-Control-Allow-Methods"] = "GET, PUT, POST, DELETE, OPTIONS"
	hs["Access-Control-Allow-Origin"] = "*"

	hs["X-Content-Type-Options"] = "nosniff"
	hs["X-Frame-Options"] = "deny"
	hs["Content-Security-Policy"] = "default-src 'none'"
	hs["X-XSS-Protection"] = "1; mode=block"
	hs["Referrer-Policy"] = "same-origin"
	if !util.Debugging {
		hs["Strict-Transport-Security"] = "max-age=31536000; preload"
	}

	api.ResponseHeaders = hs

	api.Log = func(request *http.Request, response *api.Response, err error) {
		consoleText := request.RemoteAddr + " [" + strconv.Itoa(response.Status) + "]: " + request.Method + " " + request.URL.Path
		if err == nil {
			log.Debug(color.GreenString(consoleText))
		} else {
			log.Warning(color.YellowString(consoleText + ": " + err.Error()))
			if response.Status >= http.StatusInternalServerError {
				log.Error(color.RedString(consoleText + ": " + err.Error()))
			}
		}
	}
}
