package leopard

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
	"strings"
)

type MiddlewareFunc mux.MiddlewareFunc

// GET handler register
func (a *LeopardApp) GET(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodGet, p, name, h)
}

// POST register a route with the method POST
func (a *LeopardApp) POST(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodPost, p, name, h)
}

// PUT register a route with the method PUT
func (a *LeopardApp) PUT(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodPut, p, name, h)
}

// DELETE register a route with the method DELETE
func (a *LeopardApp) DELETE(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodDelete, p, name, h)
}

// PATCH register a reoute with the method PATCH
func (a *LeopardApp) PATCH(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodPatch, p, name, h)
}

// AddRoute adds a route to the route manager
// This is mainly called by methods as GET, POST, PUT, DELETE and PATCH
// However if needed a user could register a custom method name (or one we did not include)
func (a *LeopardApp) AddRoute(method string, p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	r := a.router.NewRoute()

	r.Methods(method)
	r.Path(p)
	r.Name(name)
	r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		context := NewContext(w, r, a)

		defer func() {
			//if r := recover(); r != nil {
			//	err := context.Error(fmt.Errorf("%v", r))
			//
			//	if err != nil {
			//		return
			//	}
			//}
		}()
		h(context)
	})
}

// StaticDir register a static directory
func (a *LeopardApp) StaticDir(p string, root string) {
	h := a.fileServer(root, p)
	a.router.PathPrefix(p).Handler(h)
}

// fileServer creates a file server and returns its handler
func (a *LeopardApp) fileServer(rootDir string, p string) http.Handler {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(path.Clean(strings.TrimPrefix(r.URL.Path, p)))
		f, err := os.Open(path.Join(rootDir, strings.TrimPrefix(r.URL.Path, p)))

		if os.IsNotExist(err) {
			w.WriteHeader(404)
			return
		}

		stat, _ := f.Stat()
		maxAge := "31536000"
		w.Header().Add("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%s", maxAge))
		w.Header().Add("Content-Length", fmt.Sprintf("%d", stat.Size()))

		http.ServeContent(w, r, f.Name(), stat.ModTime(), f)
	})
	//
	//if a.CompressFiles {
	//	return handlers.CompressHandler(baseHandler)
	//}

	return baseHandler
}
