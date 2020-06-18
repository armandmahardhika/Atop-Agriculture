package httpserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/austinjan/AtopIOTServer/httpserver/router"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// func handler(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Set("Content-Type", "text/plain")
// 	w.Write([]byte("This is an example server.\n"))
// }

// Run running http server
func Run(ctx context.Context) {
	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://localhost:8000", "https://localhost:3001", "https://localhost:3000", "https://localhost:4000", "ws://localhost:3000", "ws://localhost:3001"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	srv := &http.Server{
		Addr:         ":3001",
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      handler,
		// TLSConfig:    cfg,
		// TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	router.InitRouter(r)
	go func() {
		log.Println("Https server start at port 3001")
		// if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	log.Println("Http server shoutdown!")

}
