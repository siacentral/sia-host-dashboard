package web

import (
	"context"
	"errors"

	"github.com/siacentral/sia-host-dashboard/daemon/web/router"
)

var (
	r *router.APIRouter
)

//Start starts the api router and listens on the specified address
func Start(opts router.APIOptions) error {
	r = router.NewRouter(endpoints, opts)

	return r.ListenAndServe()
}

//Shutdown attempts to gracefully shutdown the started API router
func Shutdown(ctx context.Context) error {
	if r == nil {
		return errors.New("server not started")
	}

	return r.Shutdown(ctx)
}
