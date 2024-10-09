package main


import (
	"log"
	"net/http"
)


func serveIndex(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" { 
		http.Error(w, "Error :)", http.StatusNotFound)
	}

	if r.Method != "GET" {
		http.Error(w, "Method not authorized", http.StatusNotFound)
	}

	http.ServeFile(w, r, "templates/index.html")
}

func main() {

	hub := NewHub()

	// Start a goroutine where the hub runs
	go hub.run()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// Run the server
	err := http.ListenAndServe(":3000", nil)
	log.Fatal(err)

}