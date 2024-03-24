package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func Logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("handled %s %s in %s\n", r.Method, r.RequestURI, time.Since(startTime))
	})
}

func PanicRecovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Printf("WARNING: Panic occurred: %s", err)
				log.Println(string(debug.Stack()))
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
