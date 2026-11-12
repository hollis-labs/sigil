package server

import (
	"net/http"

	"github.com/hollis-labs/substrate/apierr"
)

// RequireRequestedWith refuses a request that changes state unless it carries a non-empty
// X-Requested-With header. Browsers send that header only from script allowed by CORS, and
// this server answers no preflight, so a page on another site cannot write project files
// through the loopback listener. It is the same rule the Substrate identity module applies
// to cookie-authenticated writes; with auth off there is no identity middleware to do it.
func RequireRequestedWith(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
		default:
			if r.Header.Get("X-Requested-With") == "" {
				apierr.Write(w, http.StatusForbidden, "forbidden", "writes need an X-Requested-With header", nil)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
