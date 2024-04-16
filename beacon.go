package main

import (
	"fmt"
	"time"
)

/*
BeaconMeBaby is an amateur radio terminal program to display the live beaconing activity of the
International Beacon Project. Features of the program will be to interactively select a beaconing
schedule by band, make notation of beaconing stations recieved, and to submit a program generated
report to the server.

IBS website : https://www.ncdxf.org/beacon/
*/

type Beacon struct {
	callsign, location string
	position           int
	status             bool
}

// var _20m int = 1
var _17m int = 0

//var _15m int = -1
//var _12m int = -2
//var _10m int = -3

func getStation(tSlot int) {
	beacons := []Beacon{
		{"United Nations", "4U1UN", 1, true},
		{"Canada", "VE8AT", 2, true},
		{"United States", "W6WX", 3, true},
		{"Hawaii", "KH6RS", 4, true},
		{"New Zealand", "ZL6B", 5, true},
		{"Australia", "VK6RBP", 6, true},
		{"Japan", "JA2IGY", 7, true},
		{"Russia", "RR9O", 8, true},
		{"Hong Kong", "VR2B", 9, true},
		{"Sri Lanka", "4S7B", 10, true},
		{"South Africa", "ZS6DN", 11, true},
		{"Kenya", "5Z4B", 12, true},
		{"Israel", "4X6TU", 13, true},
		{"Finland", "OH2B", 14, true},
		{"Madeira", "CS3B", 15, true},
		{"Argentina", "LU4AA", 16, true},
		{"Peru", "OA4B", 17, true},
		{"Venezuela", "YV5B", 18, true},
	}

	// 20m == +1, 17 == 0, 15 == -1, 12 == -2, 10 == -3
	for _, station := range beacons {
		band := _17m
		roughPlace := tSlot + band
		var printStation = func() {
			fmt.Printf("%v - %v \n%v\n\n", station.position, station.callsign, station.location)
		}

		if roughPlace > 18 {
			new_place := (roughPlace) - tSlot
			if new_place == station.position {
				printStation()
			}
		} else if roughPlace < 1 {
			new_place := 18 + tSlot + band
			if new_place == station.position {
				printStation()
			}
		} else {
			if roughPlace == station.position {
				printStation()
			}
		}
	}
}

func main() {
	for {
		now := time.Now()
		if now.Second()%10 != 0 {
			time.Sleep(time.Second * 1)
		} else {
			totalSec := (now.Minute() * 60) + now.Second()
			if totalSec <= 180 {
				tSlot := totalSec / 10
				fmt.Println(now)
				getStation(tSlot)
				time.Sleep(time.Second * 9)

			} else {
				clean_time := totalSec % 180
				tSlot := clean_time / 10
				fmt.Println(now)
				if tSlot == 0 {
					getStation(18)
				} else {
					getStation(tSlot)
				}
				time.Sleep(time.Second * 9)

			}
		}
	}
}
