package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Kirito1029/mmt-flight/utils"
)

func main() {
	res, err := utils.ReadCsvFile("./ivtest-sched.csv")
	// res, err := utils.ReadCsvFile("./test.csv")
	if err != nil {
		log.Fatalf("Unable to read csv %s", err)
	}
	AirportNodes = make(map[string]*AirPort)
	for _, row := range res {
		// fmt.Println(row)
		if _, ok := AirportNodes[row[1]]; !ok {
			AirportNodes[row[1]] = NewAirPort(row[1])
		}
		if _, ok := AirportNodes[row[2]]; !ok {
			AirportNodes[row[2]] = NewAirPort(row[2])
		}
		code, _ := strconv.Atoi(string(row[0]))
		AirportNodes[row[1]].AddDestination(row[2], code, row[3], row[4])
	}

	// AirportNodes["ATQ"].GetPath("BLR")

	// AirportNodes["IXC"].GetPath("COK")

	// AirportNodes["IXC"].GetPath("GAU")

	http.HandleFunc("/getFlights", serverSchedule)

	http.ListenAndServe(":80", nil)

}

func serverSchedule(w http.ResponseWriter, req *http.Request) {
	src := req.URL.Query().Get("src")
	dest := req.URL.Query().Get("dest")
	if _, ok := AirportNodes[src]; ok {
		if _, ok := AirportNodes[dest]; ok {
			data := AirportNodes[src].GetPath(dest)
			res, _ := json.Marshal(data)
			w.Write(res)
		} else {
			http.Error(w, "No Flight available to this destination", http.StatusNotFound)
		}
	} else {
		http.Error(w, "No Flight available from this source", http.StatusNotFound)
	}
}

/*

FlightNo,From,To,start time,end time HHMM

A -> B
direct
A->...->B


list of 5 fastest

ATQ -> BLR

<id><src><dst><start><end>

*/
