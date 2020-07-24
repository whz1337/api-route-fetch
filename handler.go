package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux" //using gorilla mux to handle the custom URLs easier
)

type route struct {
	Destination string  `json:"destination"`
	Duration    float64 `json:"duration"`
	Distance    float64 `json:"distance"`
}

type correctOutput struct {
	//data to be put out
	Source string  `json:"source"`
	Routes []route `json:"routes"`
}

//correctOutput is split into []route slice and route struct in order to avoid having to deal w/ anonymous structs

type correctInput struct {
	//expected data from osrm
	Routes []struct {
		Legs []struct {
			Summary  string        `json:"summary"`
			Weight   int           `json:"weight"`
			Duration float64       `json:"duration"`
			Steps    []interface{} `json:"steps"`
			Distance float64       `json:"distance"`
		} `json:"legs"`
		WeightName string  `json:"weight_name"`
		Weight     int     `json:"weight"`
		Duration   float64 `json:"duration"`
		Distance   float64 `json:"distance"`
	} `json:"routes"`
	Waypoints []struct {
		Hint     string    `json:"hint"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"waypoints"`
	Code string `json:"code"`
}

//main is mostly boilerplate code, setting up a server on port 55532 and a handler for the proper queries.
//all other queries should automatically be handled by the default 404 golang handler.
func main() {
	r := mux.NewRouter()
	port := os.Getenv("PORT") //Heroku forces its own dynamic port on each boot, therefore this is necessary, doesn't affect localhost
	if port == "" {
		port = "55532"
	}

	fmt.Printf("attempting to init a server on port "+ port)

	r.Path("/routes").
		HandlerFunc(routesHandler).
		Queries("src", "{src}", "dst", "{dst}").
		Name("routes")

	server := &http.Server{
		Handler: r,
		Addr:    (":" + port),
	}
	log.Fatal(server.ListenAndServe())
}

func routesHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	/*
	0 - Distance
	1 - Duration
 	*/
	var parsedValues [2]float64

	data := correctOutput{}
	currentRoute := route{}
	if (floatValidator(r.Form["src"][0])) { //doesn't call the API for nothing if destination input is wrong, reduces time and network load
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}

	data.Source = r.Form["src"][0]

	for i := 0; i < len(r.Form["dst"]); i++ {
		if (floatValidator(r.Form["dst"][i])) { //same check as before
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad Request"))
			return
		}
		parsedValues = dataParser(data.Source, r.Form["dst"][i]) //makes it so the API doesn't have to be called twice per cycle for a small memory load increase
		if (parsedValues[0] != 400.0 && parsedValues[1] != -400.0) {
			currentRoute = route{r.Form["dst"][i], parsedValues[1], parsedValues[0]}
			data.Routes = append(data.Routes, currentRoute)
		}
	}

	if (parsedValues[0] == 400.0 && parsedValues[1] == -400.0) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
	} else {
		data.Routes = sortData(data.Routes)
		finalJsonStr, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(finalJsonStr))
	}
}

func dataParser(src string, dst string) [2]float64 {
	var str strings.Builder
	var parsedValues [2]float64
	time.Sleep(500) // osrm ratelimits really heavily, this should subsidise it

	//concat of the url for the project-osrm site
	str.WriteString("http://router.project-osrm.org/route/v1/driving/")
	str.WriteString(src)
	str.WriteString(";")
	str.WriteString(dst)
	str.WriteString("?overview=false")

	resp, err := http.Get(str.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()


	buffer := new(bytes.Buffer)
	buffer.ReadFrom(resp.Body)

	//this check checks if InvalidQuery or InvalidValue are present in the input, as those were the errors I was able to generate by hand
	if (!strings.Contains("buffer.String()", "InvalidQuery") && !strings.Contains("buffer.String()", "InvalidValue")) {

		data := correctInput{}
		json.Unmarshal([]byte(buffer.String()), &data)

		/*
		0 - Distance
		1 - Duration
 		*/
		parsedValues[0] = data.Routes[0].Distance
		parsedValues[1] = data.Routes[0].Duration

		return parsedValues
	} else {
		parsedValues[0] = 400.0 //400 - Bad Request
		parsedValues[1] = -400.0

		return parsedValues
	}
}

func floatValidator(subject string) bool {
	if (!strings.Contains(subject, ",")) {
		return true
	}
	_, err := strconv.ParseFloat(subject[:strings.Index(subject, ",")], 64) //checks if number before a comma is a float64
	if err != nil {
		return true
	}
	_, err = strconv.ParseFloat(subject[(strings.Index(subject, ","))+1:], 64) //same check, after a comma
	if err != nil {
		return true
	}
	return false
}

func sortData(data []route) []route {
	sort.Slice(data, func(i, j int) bool {
		if data[i].Duration < data[j].Duration {
			return true
		}
		if data[i].Duration > data[j].Duration {
			return false
		}
		return data[i].Distance < data[j].Distance
	})
	return data
}
