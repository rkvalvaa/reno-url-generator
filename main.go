package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// Configuration structures
type Config struct {
	Tenants       map[string]TenantConfig `json:"tenants"`
	PropertyTypes map[string]string       `json:"propertyTypes"`
}

type TenantConfig struct {
	Hostname  string `json:"hostname"`
	UrlScheme string `json:"urlScheme"`
}

var config Config

func main() {
	// Load configuration
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/generate-url/{tenantID}/{propertyTypeID}/{caseID}", generateURLHandler).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func generateURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenantID"]
	propertyTypeID := vars["propertyTypeID"]
	caseID := vars["caseID"]

	tenantConfig, ok1 := config.Tenants[tenantID]
	propertyType, ok2 := config.PropertyTypes[propertyTypeID]

	if !ok1 || !ok2 {
		http.Error(w, "Invalid tenantID or propertyTypeID", http.StatusBadRequest)
		return
	}

	urlScheme := tenantConfig.UrlScheme
	urlScheme = strings.Replace(urlScheme, "{propertyType}", propertyType, -1)
	urlScheme = strings.Replace(urlScheme, "{caseId}", caseID, -1)
	url := fmt.Sprintf("https://%s%s", tenantConfig.Hostname, urlScheme)

	fmt.Fprint(w, url)
}
