package models

import "time"

type Ride struct {
	ID   uint32    `json:"id"`
	Type string    `json:"type"`
	Date time.Time `json:"date"`
	From string    `json:"from"`
	To   string    `json:"to"`
}

type RideDTO struct {
	Type string `json:"type"`
	Date int64  `json:"date"`
	From string `json:"from"`
	To   string `json:"to"`
}

var Rides = []Ride{
	{ID: 1, Type: "taxi", Date: time.Now(), From: "Campus", To: "Sabiha Gokcen Airport"},
	{ID: 2, Type: "offer", Date: time.Now(), From: "Campus", To: "Kadikoy"},
	{ID: 3, Type: "request", Date: time.Now(), From: "Taksim", To: "Kampus"},
}

var RidesId uint32 = 3
