package foodpermitdata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GeoSearch(t *testing.T) {
	tests := []struct {
		name   string
		rows   []row
		lat    float64
		lng    float64
		rad    float64
		expect []row
	}{
		{
			name: "Match all",
			rows: []row{
				{Latitude: "1", Longitude: "1"},
			},
			lat: 1,
			lng: 1,
			rad: 1,
			expect: []row{
				{Latitude: "1", Longitude: "1"},
			},
		},
		{
			name: "Match a location about 157410m away",
			rows: []row{
				{Latitude: "1", Longitude: "1"},
			},
			lat: 2,
			lng: 2,
			rad: 157410,
			expect: []row{
				{Latitude: "1", Longitude: "1"},
			},
		},
		{
			name: "Not Match a location about 157410m away",
			rows: []row{
				{Latitude: "1", Longitude: "1"},
			},
			lat:    2,
			lng:    2,
			rad:    157400,
			expect: []row{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd, err := parseGeoData(tt.rows)
			assert.NoError(t, err)
			d := FoodPermitData{geoData: gd, data: tt.rows}
			result := d.GeoSearch(tt.lat, tt.lng, tt.rad)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func Test_AutoComplete(t *testing.T) {
	tests := []struct {
		name   string
		rows   []row
		prefix string
		expect []string
	}{
		{
			name: "Match all",
			rows: []row{
				{Applicant: "1"},
				{Applicant: "12"},
				{Applicant: "123"},
			},
			prefix: "1",
			expect: []string{
				"1",
				"12",
				"123",
			},
		},
		{
			name: "Match some",
			rows: []row{
				{Applicant: "1"},
				{Applicant: "12"},
				{Applicant: "123"},
			},
			prefix: "12",
			expect: []string{
				"12",
				"123",
			},
		},
		{
			name: "Match none",
			rows: []row{
				{Applicant: "1"},
				{Applicant: "12"},
				{Applicant: "123"},
			},
			prefix: "2",
			expect: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "applicant"
			m := createFoodPermitTries(tt.rows, []string{key})
			d := FoodPermitData{foodPermitTries: m}
			result := d.GetSuggestion(tt.prefix, key)
			assert.Equal(t, tt.expect, result)
		})
	}
}
