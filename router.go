package noaweb

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ListAllRoutes lists all routes associated with the given instance.
// Includes all routes in router and all subrouters.
func (i *Instance) ListAllRoutes() {
	i.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})
	fmt.Println()
}

// RegRoute registers a route. Method should be GET or POST.
func (i *Instance) RegRoute(path, method string, handlers ...HandlerFunc) {
	i.router.Path(path).Methods(method).HandlerFunc(chain(handlers...))
}

// RegStatic registers a route to static content like images, css files, etc.
// path is the relative path in URL, and filePath is the path to the static content.
func (i *Instance) RegStatic(path, filePath string) {
	if strings.HasPrefix(filePath, "/") {
		i.router.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir(filePath))))
	} else {
		i.router.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir(i.AssetsDir+"/"+filePath))))
	}
}

// NewSubrouter registers a subrouter with a given name and a pathPrefix.
// name should be used when refering to the subrouter, example when registering routes and such.
func (i *Instance) NewSubrouter(name, pathPrefix string) {
	i.subrouters[name] = i.router.PathPrefix(pathPrefix).Subrouter()
}

// RegSubRoute is the same as RegRoute(), but works on a subroute level on the subroute subName.
func (i *Instance) RegSubRoute(subName, path, method string, handlers ...HandlerFunc) {
	if _, ok := i.subrouters[subName]; ok == false {
		panic("RegSubRoute: subrouteName does not exist!")
	}

	i.subrouters[subName].Path(path).Methods(method).HandlerFunc(chain(handlers...))
}

// RegSubStatic is the same as RegStatic(), but works on a subroute level on the subroute subName.
func (i *Instance) RegSubStatic(subName, path, filePath string) {
	if strings.HasPrefix(filePath, "/") {
		i.subrouters[subName].PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir(filePath))))
	} else {
		i.subrouters[subName].PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir(i.AssetsDir+"/"+filePath))))
	}
}
