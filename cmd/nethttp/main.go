package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var serverStartTime time.Time

func Authorization(req *http.Request, w http.ResponseWriter) error {
	auth := req.Header.Get("Authorization")
	api_key := os.Getenv("API_KEY")
	if api_key == "" {
		http.Error(w, "Server Missing Authorization Key", http.StatusInternalServerError)
		return fmt.Errorf("No Authorization Key in Server")
	}
	if auth == "" {
		http.Error(w, "Missing Authorization Key", http.StatusBadRequest)
		return fmt.Errorf("No Authorization Key in Header")
	}
	token, found := strings.CutPrefix(auth, "Bearer ")
	if !found {
		http.Error(w, "Authorization Header Misformatted", http.StatusBadRequest)
		return fmt.Errorf("Impropert API Format")
	}
	if token != api_key {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return fmt.Errorf("Unauthorized")
	}
	return nil
}

//Age set to String to simplify mapping

type User struct {
	Name  string `json:"name"`
	Age   string `json:"age"`
	Email string `json:"email"`
}

func main() {

	godotenv.Load()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /address", func(w http.ResponseWriter, req *http.Request) {
		err := Authorization(req, w)
		if err != nil {
			return
		}
		io.WriteString(w, "request allowed\n")
		fmt.Println("Address ran")
	})
	mux.HandleFunc("GET /uptime", func(w http.ResponseWriter, req *http.Request) {
		err := Authorization(req, w)
		if err != nil {
			return
		}
		uptime := time.Since(serverStartTime)
		truncate := uptime.Truncate(time.Second)
		fmt.Println("The Server has been running for " + truncate.String())
		io.WriteString(w, "The Server has been running for "+truncate.String())
	})
	mux.HandleFunc("POST /server", func(w http.ResponseWriter, req *http.Request) {
		err := Authorization(req, w)
		if err != nil {
			return
		}
		contenttype := req.Header.Get("Content-Type")
		if contenttype != "application/json" {
			http.Error(w, "Content Incorrect Format", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error Reading Body", http.StatusBadRequest)
			return
		}
		var user User

		jsonerr := json.Unmarshal(body, &user)
		if jsonerr != nil {
			http.Error(w, "Json Format Error", http.StatusBadRequest)
			return
		}
		fmt.Printf("Mapped struct: %+v\n", user)
		io.WriteString(w, "The Users name is "+user.Name)
	})
	serverStartTime = time.Now()
	srv := http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		//Use CrossOrigin Protection.Handler to block all non-safe cross-origin
		//browser request to mux
		Handler: http.NewCrossOriginProtection().Handler(mux),
	}
	log.Fatal(srv.ListenAndServe())

}
