#foodpermit

## Usage

```
docker build -t foodpermit .
docker run -p 8000:8000 --rm -it foodpermit 
```

###Geo queries
```
http://localhost:8000/geosearch?lat=123&lng=222&rad=100
```
Geo expects lat(latitude), lng(longitude) and rad(radius in meters) as URL parameters
Geo queries is done using a O(n) for loop to compare the haversin distance
TODO: Optimization using possible a quad tree for range queries

###Autocomplete
```
http://localhost:8000/autocomplete?val=55&key=address
```
Autocomplete expects val, key as URL parameters
`val` is the prefix to match against
`key` is the data that we want to autocomplete against
Autocomplete is done using a trie for each selected fields
