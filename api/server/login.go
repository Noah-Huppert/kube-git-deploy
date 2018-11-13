package server

import (
	"net/http"
	"net/url"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"

	"github.com/Noah-Huppert/golog"
)

// GHLoginURL is the base GitHub login URL
const GHLoginURL string = "https://github.com/login/oauth/authorize"

// GHRedirectURL is the URL users should be redirected to after login
const GHRedirectURL string = "http://localhost:5000/api/v0/github/oauth_callback"

// GHLoginURLHandler returns the GitHub login URL to send the user to
type GHLoginURLHandler struct {
	// logger prints debug information
	logger golog.Logger

	// cfg is configuration
	cfg *config.Config
}

// ServeHTTP implements http.Handler
func (h GHLoginURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create responder
	responder := NewJSONResponder(h.logger, w)

	// Build login URL
	u, err := url.Parse(GHLoginURL)
	if err != nil {
		h.logger.Errorf("error parsing github login URL: %s",
			err.Error())

		responder.Respond(http.StatusInternalServerError,
			map[string]interface{}{
				"ok":    false,
				"error": "failed to build GitHub login URL",
			})

		return
	}
	q := u.Query()
	q.Set("client_id", h.cfg.GitHubClientID)
	q.Set("redirect_url", GHRedirectURL)

	u.RawQuery = q.Encode()

	// Return login URL
	responder.Respond(http.StatusOK, map[string]interface{}{
		"login_url": u.String(),
		"ok":        true,
	})
}
