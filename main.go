package main

// TODO: Add proper html error/success codes along with responses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Asset struct {
	ID      string `json:"id"`
	Details string `json:"details"`
}

var assetDB = []Asset{
	{"2", "It worked"},
	{"8", "Even better"},
}

func SearchAssetDB(assetID string) (Asset, bool) {
	for _, assetFromDB := range assetDB {
		// Check assetID against assetFromDB
		if assetFromDB.ID == assetID {
			return assetFromDB, true
		}
	}
	// Return an empty asset and false if no match to assetID was found
	return Asset{}, false
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", FetchAssetDB)
		r.Get("/{id}", FetchAssetDetailsByID)
	})

	r.Group(func(r chi.Router) {
		// TODO: Add AuthMiddleware

		r.Post("/manage", CreateAsset)
	})

	http.ListenAndServe(":3000", r)
}

func FetchAssetDB(w http.ResponseWriter, r *http.Request) {
	// Encode assetDB with JSON
	assetByte, err := json.Marshal(assetDB)
	if err != nil {
		fmt.Fprintln(w, "An issue occured when marshaling")
		return
	}

	w.Write(assetByte)
}

func FetchAssetDetailsByID(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id") // Get assetID param

	// Search for asset
	asset, ok := SearchAssetDB(assetID)
	if !ok {
		fmt.Fprintln(w, "An issue occured when searching DB for ID")
		return
	}

	// Encode asset with JSON
	assetByte, err := json.Marshal(asset)
	if err != nil {
		fmt.Fprintln(w, "An issue occured when marshaling")
		return
	}

	w.Write(assetByte)
}

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var asset Asset

	// Decode incomming JSON into asset
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&asset); err != nil {
		fmt.Fprintln(w, "An issue occured when decoding asset")
		return
	}

	// Check if asset with assetID already exists
	_, ok := SearchAssetDB(asset.ID)
	if ok {
		fmt.Fprintln(w, "Asset with provided ID already exists")
		return
	}

	// Save asset to assetDB
	assetDB = append(assetDB, asset)

	fmt.Fprintln(w, "Successfully saved asset to DB")
}
