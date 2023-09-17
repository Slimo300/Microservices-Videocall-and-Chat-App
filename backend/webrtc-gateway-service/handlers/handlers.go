package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Slimo300/MicroservicesChatApp/backend/webrtc-gateway-service/database"
	"github.com/gorilla/mux"
)

func Setup(db database.DBLayer, origin string) http.Handler {
	mux := mux.NewRouter()

	mux.Handle("/video-call/{callID}/websocket", NewReverseProxy(db)).Methods("GET")

	corsMux := CORSMiddleware(mux, origin)

	return corsMux
}

func NewReverseProxy(db database.DBLayer) http.Handler {

	return &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetXForwarded()

			params := mux.Vars(pr.In)
			callID := params["callID"]

			var domainName string
			var err error

			domainName, err = db.GetCallInstanceDomainName(callID)
			if err != nil {
				domainName, err = db.GetLeastUsedInstanceDomainName()
				if err != nil {
					panic(err)
				}
			}

			if err := db.AddConnection(callID, domainName); err != nil {
				panic(err)
			}

			pr.Out.URL, err = url.Parse(fmt.Sprintf("http://%s/video-call/%s/ws", domainName, callID))
			if err != nil {
				panic(err)
			}
			pr.Out.URL.RawQuery = pr.In.URL.RawQuery

			go func() {
				<-pr.Out.Context().Done()

				if err := db.DeleteConnection(callID, domainName); err != nil {
					log.Printf("Error when trying to delete connection: %v", err)
				}
			}()

		},
	}
}

func CORSMiddleware(handler http.Handler, origin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
