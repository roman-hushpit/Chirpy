package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

type bodyContent struct {
	Body string `json:"body"`
}		
type errorBody struct {
	Error string `json:"error"`
}

type success struct {
	CleanedBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits = cfg.fileserverHits + 1
		next.ServeHTTP(w, r)
	})
}


func (cfg *apiConfig) addcountHeader(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	errorReponse := errorBody{
		Error: msg,
	}
	w.WriteHeader(code)
	data, _ := json.Marshal(errorReponse)
	w.Write(data)
}

var badWords = []string{"kerfuffle", "sharbert", "fornax"}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(code)
    w.Write(response)
    return nil
}

func cleanString(message string) string{
	words := strings.Split(message, " ")
	for i, word := range words {
		if contains(badWords, word) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == strings.ToLower(e) {
            return true
        }
    }
    return false
}

func main() {
	apiCfg := apiConfig{fileserverHits: 0} 
	fs := http.FileServer(http.Dir("."))
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		params := bodyContent{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, 400, "Something went wrong")
			return
		}
		requestBody := params.Body
		if len(requestBody) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}

		successResponse := success{
			CleanedBody: cleanString(requestBody),
		}
		respondWithJSON(w, 200, successResponse)
	})

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		//apiCfg.addcountHeader(w)
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(`<html> 
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>

		</html>`, apiCfg.fileserverHits)))
	})
	mux.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits = 0
	})
	srv := &http.Server {
		Addr: ":8080",
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Server failed:", err)
	}
}
