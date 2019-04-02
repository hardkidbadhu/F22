package route

import (
	"F22/datastore"
	"F22/handlers"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"F22/config"
	"F22/models"

	"github.com/globalsign/mgo"
	"github.com/gorilla/context"
)

//Adapter pattern implemented in middleware, please refer the below url for information
//https://go-talks.appspot.com/github.com/matryer/golanguk/building-apis.slide#30
//wrapping up http.handler with addition functionality
type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
// responseWriter for tracking HTTP response status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging middleware
//logs each request in the way the origin method requested statuscode and the end point URL
func loggingHandler(l *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			w := &responseWriter{rw, http.StatusOK}
			t := time.Now()
			h.ServeHTTP(w, r)
			l.Printf("(%s) [%s] %d %q %v\n", getClientIP(r), r.Method, w.statusCode, r.URL.String(), time.Since(t))
		})
	}
}

// Recovery middleware
//If any panic happens in the http request it recovers and maintains the smooth running of server without getting crashed
func recoverHandler() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic: %+v", err)
					debug.PrintStack()
					http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(rw, r)
		})
	}
}

func sessionAuthenticator (dbSession *mgo.Session, cfg *config.Config) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")

			var (
				error error
				userIns *models.User
			)

			log.Print("token", token)

			if userIns, error = datastore.NewUser(dbSession, cfg).FindByToken(token); error != nil && userIns == nil {
				log.Printf("Error - Middleware - %s", error.Error())
				args := make(map[string]interface{})

				args["authenticate"] = fmt.Sprintf("%s/authenticate", cfg.Endpoint)
				args["signUp"] = fmt.Sprintf("%s/signUp", cfg.Endpoint)
				args["name"] = `Guest`
				tempStr := handlers.LoadTemplate(cfg, "views/login.html", args)
				fmt.Fprint(rw, tempStr)
				return
			}

			context.Set(r, "userId", userIns.Id.Hex())
			context.Set(r, "sessionUser", userIns.Name)
			h.ServeHTTP(rw, r)
		})
	}
}

//Function to get the client IP address
func getClientIP(r *http.Request) string {
	if remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return remoteAddr
	}

	// Header X-Forwarded-For
	hdrForwardedFor := http.CanonicalHeaderKey("X-Forwarded-For")
	if fwdFor := strings.TrimSpace(r.Header.Get(hdrForwardedFor)); fwdFor != "" {
		index := strings.Index(fwdFor, ",")
		if index == -1 {
			return fwdFor
		}
		return fwdFor[:index]
	}

	// Header X-Real-Ip
	hdrRealIP := http.CanonicalHeaderKey("X-Real-Ip")
	if realIP := strings.TrimSpace(r.Header.Get(hdrRealIP)); realIP != "" {
		return realIP
	}

	return ""
}
