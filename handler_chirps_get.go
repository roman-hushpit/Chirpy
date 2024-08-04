package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
		AuthorId: dbChirp.AuthorId,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	
	authorId := r.URL.Query().Get("author_id")


	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if authorId != "" {
			numericAuthorId, _ := strconv.Atoi(authorId)
			if numericAuthorId == dbChirp.AuthorId {
				chirps = append(chirps, Chirp{
					ID:   dbChirp.ID,
					Body: dbChirp.Body,
					AuthorId: dbChirp.AuthorId,
				})
			}
		} else {
			chirps = append(chirps, Chirp{
				ID:   dbChirp.ID,
				Body: dbChirp.Body,
				AuthorId: dbChirp.AuthorId,
			})
		}
	}
	
	sortParameter := getSortParameter(r.URL.Query().Get("sort"))
	if sortParameter == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func getSortParameter(sort string) string {
	if sort == "" || sort == "asc"{
		return "asc" 
	}
	if sort == "desc"{
		return "desc"
	}
	return "asc"
}
