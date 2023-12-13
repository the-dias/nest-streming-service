package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"nats-service/app/cache"
	"net/http"
	"strconv"
	"sync"
)

type HttpServer struct {
	cache      *cache.Cache
	mutex      *sync.RWMutex
	staticFile string
}

func New(cache *cache.Cache, mutex *sync.RWMutex, staticFile string) *HttpServer {
	return &HttpServer{
		cache:      cache,
		mutex:      mutex,
		staticFile: staticFile,
	}
}

func (hs *HttpServer) StartHttpServer(PORT int, patternServer, patternStatic, staticDir string) {

	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		hs.GetOrderHandler(w, r, hs.staticFile)
	}

	http.HandleFunc(patternServer, httpHandler)
	http.Handle(patternStatic, http.StripPrefix(patternStatic, http.FileServer(http.Dir("static"))))

	fmt.Printf("Server started op port %d:", PORT)

	go http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}

// func HttpHandler(w http.ResponseWriter, r *http.Request) {
// 	GetOrderHandler(w, r, cacheHolder, &cacheMutex)
// }

func (hs *HttpServer) GetOrderHandler(w http.ResponseWriter, r *http.Request, staticFiles string) {
	hs.mutex.RLock()
	defer hs.mutex.RUnlock()

	orderData := struct {
		ID    int
		Order string
		error string
	}{
		ID:    -1,
		Order: "",
		error: "",
	}

	tmpl, err := template.ParseFiles(staticFiles)

	if err != nil {
		http.Error(w, "Can't parse static files", http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")
	fmt.Println(idStr)
	fmt.Println("Order")

	if idStr == "" {
		if err := tmpl.Execute(w, orderData); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	item, ok := hs.cache.Get(id)
	if !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	jsonData, err := json.MarshalIndent(item, " ", " ")
	fmt.Println(string(jsonData))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	orderData.ID = id
	orderData.Order = string(jsonData)

	if err := tmpl.Execute(w, orderData); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
