package foodpermitdata

import (
	"encoding/csv"
	"errors"
	"fmt"
	trie "github.com/vivekn/autocomplete"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type row struct {
	Locationid          string
	Applicant           string
	FacilityType        string
	Cnn                 string
	LocationDescription string
	Address             string
	Blocklot            string
	Block               string
	Lot                 string
	Permit              string
	Status              string
	FoodItems           string
	X                   string
	Y                   string
	Latitude            string
	Longitude           string
	Schedule            string
	NOISent             string
	Approved            string
	Received            string
	PriorPermit         string
	ExpirationDate      string
	Location            string
}

type latLng struct {
	lat float64
	lng float64
}

type FoodPermitData struct {
	geoData         []latLng
	data            []row
	foodPermitTries map[string]*trie.Trie
}

//https://gist.github.com/cdipaolo/d3f8db3848278b49db68
// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func calculateDist(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func (d FoodPermitData) GeoSearch(lat, long, rad float64) []row {
	rows := make([]row, 0, len(d.geoData))
	for idx, gd := range d.geoData {
		if calculateDist(gd.lat, gd.lng, lat, long) < rad {
			rows = append(rows, d.data[idx])
		}
	}
	return rows
}

func parseFile(f *os.File) ([]row, error) {
	d := make([]byte, 348814)
	n, err := f.Read(d)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error reading file"))
	}

	fileString := string(d[:n])
	fileString = strings.TrimSuffix(fileString, "\r")

	r := csv.NewReader(strings.NewReader(fileString))
	records, err := r.ReadAll()
	if err != nil {
		fmt.Printf("err = %v\n", err)
	}
	//Remove first record as it is the header
	records = records[1:]

	rows := make([]row, 0, len(records))
	for _, values := range records {
		if len(values) != 23 {
			continue
		}
		r := row{
			Locationid:          values[0],
			Applicant:           values[1],
			FacilityType:        values[2],
			Cnn:                 values[3],
			LocationDescription: values[4],
			Address:             values[5],
			Blocklot:            values[6],
			Block:               values[7],
			Lot:                 values[8],
			Permit:              values[9],
			Status:              values[10],
			FoodItems:           values[11],
			X:                   values[12],
			Y:                   values[13],
			Latitude:            values[14],
			Longitude:           values[15],
			Schedule:            values[16],
			NOISent:             values[17],
			Approved:            values[18],
			Received:            values[19],
			PriorPermit:         values[20],
			ExpirationDate:      values[21],
			Location:            values[22],
		}
		rows = append(rows, r)
	}
	return rows, nil
}

func parseGeoData(rows []row) ([]latLng, error) {
	lls := make([]latLng, 0, len(rows))
	for _, r := range rows {
		lat, err := strconv.ParseFloat(r.Latitude, 64)
		if err != nil {
			fmt.Printf("fail to parse row: %v, err: %v\n", r, err)
			continue
		}
		lng, err := strconv.ParseFloat(r.Longitude, 64)
		if err != nil {
			fmt.Printf("fail to parse row: %v, err: %v\n", r, err)
			continue
		}
		ll := latLng{
			lat: lat,
			lng: lng,
		}
		lls = append(lls, ll)
	}
	return lls, nil
}

func NewFoodPermitData() (FoodPermitData, error) {
	file, err := os.Open("Mobile_Food_Permit_Map.csv")
	if err != nil {
		return FoodPermitData{}, errors.New("fail to load data")
	}
	defer file.Close()

	rows, err := parseFile(file)
	if err != nil {
		return FoodPermitData{}, errors.New("fail to parse data")
	}
	geoData, err := parseGeoData(rows)
	if err != nil {
		return FoodPermitData{}, errors.New("fail to parse data")
	}

	foodPermitTries := createFoodPermitTries(rows, []string{"applicant", "address", "LocationDescription"})
	return FoodPermitData{geoData: geoData, data: rows, foodPermitTries: foodPermitTries}, nil
}

func (d FoodPermitData) GetSuggestion(value, key string) []string {
	trie := d.getTrie(key)
	//Only returns error when there is empty result
	suggestion, _ := trie.AutoComplete(strings.ToLower(value))
	return suggestion
}

func (d FoodPermitData) getTrie(key string) *trie.Trie {
	return d.foodPermitTries[key]
}

func createTrie(rows []row, key string) *trie.Trie {
	t := trie.NewTrie()
	for _, r := range rows {
		switch key {
		case "address":
			t.Insert(strings.ToLower(r.Address))
		case "LocationDescription":
			t.Insert(strings.ToLower(r.Address))
		case "applicant":
			fallthrough
		default:
			t.Insert(strings.ToLower(r.Applicant))
		}
	}
	return t
}

func createFoodPermitTries(rows []row, keys []string) map[string]*trie.Trie {
	begin := time.Now()

	var wg sync.WaitGroup
	foodPermitTries := make(map[string]*trie.Trie)
	for _, k := range keys {
		wg.Add(1)
		k := k
		var t *trie.Trie
		go func() {
			t = createTrie(rows, k)
			foodPermitTries[k] = t
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Printf("tries created in: %v\n", time.Now().Sub(begin))
	return foodPermitTries
}
