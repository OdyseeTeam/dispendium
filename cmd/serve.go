package cmd

import (
	"net/http"
	"strconv"

	"github.com/lbryio/dispendium/actions"
	"github.com/lbryio/dispendium/config"
	"github.com/lbryio/dispendium/jobs"

	"github.com/pkg/profile"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	serveCmd.PersistentFlags().StringP("host", "", "0.0.0.0", "host to listen on")
	serveCmd.PersistentFlags().IntP("port", "p", 6060, "port binding used for the api server")
	//Bind to Viper
	err := viper.BindPFlag("host", serveCmd.PersistentFlags().Lookup("host"))
	if err != nil {
		logrus.Panic(err)
	}
	err = viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))
	if err != nil {
		logrus.Panic(err)
	}
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the dispendium API server",
	Long:  `Runs the dispendium API server`,
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("codeprofile") {
			defer profile.Start(profile.NoShutdownHook).Stop()
		}
		config.InitializeConfiguration()
		routes := actions.GetRoutes()
		httpServeMux := http.NewServeMux()
		httpServeMux.Handle(promPath, promBasicAuthWrapper(promhttp.Handler()))
		routes.Each(func(pattern string, handler http.Handler) {
			httpServeMux.Handle(pattern, handler)
		})

		actions.ConfigureAPIServer()
		host := viper.GetString("host")
		port := viper.GetInt("port")
		jobs.Start()
		logrus.Infof("API Server started @ %s", "http://"+host+":"+viper.GetString("port")+"/")
		logrus.Fatal(http.ListenAndServe(host+":"+strconv.Itoa(port), httpServeMux))
	},
}

const promPath = "/metrics"

func promBasicAuthWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "authentication required", http.StatusBadRequest)
			return
		}
		if user == "prom" && pass == "prom-dispendium-access" {
			h.ServeHTTP(w, r)
		} else {
			http.Error(w, "invalid username or password", http.StatusForbidden)
		}
	})
}
