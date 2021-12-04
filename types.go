package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var data [][]int

const LIST_SIZE = 5

type (
	AirPort struct {
		Name         string
		Destinations []Path
	}
	Path struct {
		Name    string
		Code    int
		Minutes float64
		Start   string
		End     string
		AirPort *AirPort
	}
	Response struct {
		Res map[string]TracePath
	}
	TracePath struct {
		PathCode string
		Minutes  float64
	}
)

func NewAirPort(Name string) (res *AirPort) {
	res = &AirPort{
		Name:         Name,
		Destinations: []Path{},
	}
	return
}

var AirportNodes map[string]*AirPort

func (a *AirPort) AddDestination(Name string, code int, FromTime string, ToTime string) {
	minutes := TimeBetween(FromTime, ToTime)
	if minutes == 0 {
		fmt.Println("HERERE", Name, code)
	}
	a.Destinations = append(a.Destinations, NewPath(Name, code, minutes, FromTime, ToTime))
}

func NewPath(Name string, code int, minutes float64, FromTime, ToTime string) Path {
	return Path{
		Name:    Name,
		Code:    code,
		Minutes: minutes,
		AirPort: AirportNodes[Name],
		Start:   FromTime,
		End:     ToTime,
	}
}

func TimeBetween(start string, end string) float64 {

	if start == "" {
		return 0
	}
	endI, _ := strconv.Atoi(end)
	startI, _ := strconv.Atoi(start)
	var startTime, endTime time.Time
	startTime = time.Date(0, 0, 0, startI/100, startI%100, 0, 0, time.UTC)
	if startI > endI {
		endTime = time.Date(0, 0, 1, endI/100, endI%100, 0, 0, time.UTC)
	} else {
		endTime = time.Date(0, 0, 0, endI/100, endI%100, 0, 0, time.UTC)
	}
	res := endTime.Sub(startTime)
	res1 := res.Minutes()
	return res1
}

func GetPath(from string, to string) error {
	srcAirport, ok := AirportNodes[from]
	if !ok {
		return errors.New("No Path")
	}
	_ = srcAirport
	return nil
}

type PathX struct {
	reachAirport *AirPort
	// Path          Path
	TravelMinutes float64
	CurrentTime   string
	FullPath      string
	codePath      string
}

func (a *AirPort) GetPath(Name string) (responsePath *Response) {
	var Q []PathX
	Q = append(Q, PathX{
		reachAirport:  a,
		TravelMinutes: 0,
		FullPath:      a.Name,
	})
	responsePath = &Response{
		Res: make(map[string]TracePath),
	}
	maxMinutes := -1.0
	count := 0
	for len(Q) > 0 {
		cur := Q[0]
		if len(Q) > 0 {
			Q = Q[1:] // pop
		}
		dest := cur.reachAirport.Destinations
		for _, res := range dest {
			if res.Name == Name {
				if count >= LIST_SIZE {
					if maxMinutes > cur.TravelMinutes+res.Minutes+TimeBetween(cur.CurrentTime, res.Start) {
						var pathCoderes string
						if cur.codePath == "" {
							pathCoderes = strconv.Itoa(res.Code)
						} else {
							pathCoderes = cur.codePath + "_" + strconv.Itoa(res.Code)
						}
						responsePath.Res[cur.FullPath+"_"+res.Name+"||"+pathCoderes] = TracePath{
							PathCode: pathCoderes,
							Minutes:  cur.TravelMinutes + res.Minutes + TimeBetween(cur.CurrentTime, res.Start),
						}
						maxMinutes = responsePath.RemoveMax()
						count++
					}
				} else {
					var pathCoderes string
					if cur.codePath == "" {
						pathCoderes = strconv.Itoa(res.Code)
					} else {
						pathCoderes = cur.codePath + "_" + strconv.Itoa(res.Code)
					}
					responsePath.Res[cur.FullPath+"_"+res.Name+"||"+pathCoderes] = TracePath{
						PathCode: pathCoderes,
						Minutes:  cur.TravelMinutes + res.Minutes + TimeBetween(cur.CurrentTime, res.Start),
					}
					if maxMinutes < cur.TravelMinutes+res.Minutes+TimeBetween(cur.CurrentTime, res.Start) {
						maxMinutes = cur.TravelMinutes + res.Minutes + TimeBetween(cur.CurrentTime, res.Start)
					}
					count++
				}
			} else {
				if !strings.Contains(cur.FullPath, "_") {
					if !strings.Contains(cur.FullPath, res.AirPort.Name) {
						// fmt.Println(cur.FullPath, res.Name)
						var restime, startTime string
						if cur.CurrentTime == "" {
							restime = AddTime(res.End, time.Minute*120)
							startTime = res.Start
						} else {
							restime = AddTime(cur.CurrentTime, time.Minute*(120+time.Duration(res.Minutes)))
							startTime = cur.CurrentTime
						}
						diffTime := TimeBetween(startTime, restime)
						if count >= LIST_SIZE && cur.TravelMinutes+diffTime > maxMinutes {
							// fmt.Println("Skipping", cur.FullPath, res.Name, maxMinutes)
							continue
						}
						var pathCoderes string
						if cur.codePath == "" {
							pathCoderes = strconv.Itoa(res.Code)
						} else {
							pathCoderes = cur.codePath + "_" + strconv.Itoa(res.Code)
						}
						Q = append(Q, PathX{
							reachAirport:  res.AirPort,
							TravelMinutes: cur.TravelMinutes + diffTime,
							CurrentTime:   restime,
							FullPath:      cur.FullPath + "_" + res.AirPort.Name,
							codePath:      pathCoderes,
						})
					}
				}
			}
		}
	}
	return responsePath
}

func (r *Response) RemoveMax() float64 {
	max := -1.0
	maxKey := ""
	for k, v := range r.Res {
		if max < v.Minutes {
			max = v.Minutes
			maxKey = k
		}
	}
	delete(r.Res, maxKey)
	return max
}

func AddTime(currentTime string, addTime time.Duration) string {
	a, _ := strconv.Atoi(currentTime)
	res := time.Date(0, 0, 0, a/100, a%100, 0, 0, time.UTC)
	res = res.Add(addTime)
	return res.Format("1504")
}
