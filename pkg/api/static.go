package api

import (
	"net/http"

	"github.com/ncarlier/webhookd/pkg/config"
)

func static(prefix string) HandlerFunc {
	return func(conf *config.Config) http.Handler {
		if conf.StaticDir != "" {
			fs := http.FileServer(http.Dir(conf.StaticDir))
			return http.StripPrefix(prefix, fs)
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "404 page not found", http.StatusNotFound)
		})
	}
}
