package main

import (
 "encoding/json"
 "fmt"
 "net/http"

 "github.com/paulmach/orb"
 "github.com/paulmach/orb/geojson"
)

type Line struct {
 ID string `json:"id"`
 Start orb.Point
 End orb.Point
}

type Intersection struct {
 ID string `json:"id"`
 Location orb.Point `json:"location"`
}

func main() {
 // define the 50 lines
 lines := []Line{
  {ID: "L01", Start: orb.Point{-10, -10}, End: orb.Point{10, 10}},
  {ID: "L02", Start: orb.Point{-10, 0}, End: orb.Point{10, 0}},
  // ... define the remaining 48 lines
 }

 // create the HTTP handler function
 handler := func(w http.ResponseWriter, r *http.Request) {
  // check the authorization header
  authHeader := r.Header.Get("Authorization")
  if authHeader != "my-secret-token" {
   http.Error(w, "Unauthorized", http.StatusUnauthorized)
   return
  }

  // parse the GeoJSON linestring from the request body
  var linestring geojson.LineString
  err := json.NewDecoder(r.Body).Decode(&linestring)
  if err != nil {
   http.Error(w, err.Error(), http.StatusBadRequest)
   return
  }

  // check that the linestring is valid
  if !linestring.Valid() {
   http.Error(w, "Invalid linestring", http.StatusInternalServerError)
   return
  }

  // check for intersections with each line
  intersections := []Intersection{}
  for _, line := range lines {
   // convert the line to a GeoJSON feature
   lineFeature := geojson.NewLineStringFeature([]orb.Point{line.Start, line.End})
   if !lineFeature.Valid() {
    http.Error(w, "Invalid line", http.StatusInternalServerError)
    return
   }

   // check for intersections
   if linestring.Intersects(lineFeature.Geometry) {
    // compute the intersection point
    location, ok := linestring.Intersection(lineFeature.Geometry)
    if !ok {
     http.Error(w, "Failed to compute intersection", http.StatusInternalServerError)
     return
    }

    // add the intersection to the list
    intersections = append(intersections, Intersection{
     ID: line.ID,
     Location: location,
    })
   }
  }

  // encode the result as JSON and write to the response
  json.NewEncoder(w).Encode(intersections)
 }

 // start the HTTP server
 http.HandleFunc("/", handler)
 err := http.ListenAndServe(":8080", nil)
 if err != nil {
  fmt.Println(err)
 }
}
