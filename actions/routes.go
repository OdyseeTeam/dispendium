package actions

import (
	"net/http"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/orderedmap"
)

// Routes holds a map of api handlers with the key being the route
type Routes struct {
	m *orderedmap.Map
}

func (r *Routes) set(key string, h api.Handler) {
	if r.m == nil {
		r.m = orderedmap.New()
	}
	r.m.Set(key, h)
}

// Each applies to a function of the type specified to each of the Routes
func (r *Routes) Each(f func(string, http.Handler)) {
	if r.m == nil {
		return
	}
	for _, k := range r.m.Keys() {
		a, _ := r.m.Get(k)
		f(k, a.(http.Handler))
	}
}

// GetRoutes returns the set of Routes specified for the API server
func GetRoutes() *Routes {
	routes := Routes{}
	routes.set("/", Root)
	routes.set("/test", Test)

	routes.set("/send", Send)
	routes.set("/balance", Balance)

	return &routes
}
