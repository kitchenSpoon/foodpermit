package foodpermit

import (
	"encoding/json"
	"fmt"
	"foodpermit/internal/foodpermitdata"
	"net/http"
	"strconv"
)

type Service struct {
	data foodpermitdata.FoodPermitData
}

func (s *Service) Geosearch(w http.ResponseWriter, req *http.Request) {
	latStr := req.URL.Query().Get("lat")
	lngStr := req.URL.Query().Get("lng")
	radStr := req.URL.Query().Get("rad")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		WriteInvalidRequest(w, "lat(latitude) param must not be empty")
		return
	}
	if lat < -90 || lat > 90 {
		WriteInvalidRequest(w, "lat(latitude) param must be within -90 and 90")
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		WriteInvalidRequest(w, "lng(longtitude) param must not be empty")
		return
	}
	if lat < -180 || lat > 180 {
		WriteInvalidRequest(w, "lng(longtitude) param must be within -180 and 180")
		return
	}

	rad, err := strconv.ParseFloat(radStr, 64)
	if err != nil {
		WriteInvalidRequest(w, "rad(radius) param must not be empty")
		return
	}
	if rad < 0 {
		WriteInvalidRequest(w, "rad(radius) param must be positive")
		return
	}

	rows := s.data.GeoSearch(lat, lng, rad)
	err = json.NewEncoder(w).Encode(rows)
	if err != nil {
		WriteInternalServerError(w)
	}
}

func (s *Service) GetSuggestion(w http.ResponseWriter, req *http.Request) {
	val := req.URL.Query().Get("val")
	key := req.URL.Query().Get("key")
	if key == "" {
		WriteInvalidRequest(w, "key param must not be empty")
		return
	}
	sug := s.data.GetSuggestion(val, key)
	err := json.NewEncoder(w).Encode(sug)
	if err != nil {
		WriteInternalServerError(w)
	}
}

func (s *Service) Root(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello"))
}

func WriteInvalidRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(400)
	fmt.Fprintf(w, "Invalid request: %s", msg)
}

func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("Internal Server Error"))
}

func NewService() (Service, error) {
	data, err := foodpermitdata.NewFoodPermitData()
	if err != nil {
		return Service{}, fmt.Errorf("fail to create service: %w", err)
	}
	return Service{data: data}, nil
}
