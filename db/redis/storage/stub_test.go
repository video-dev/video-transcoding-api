package storage

import "time"

type Person struct {
	ID               string    `redis-hash:"-"`
	Name             string    `redis-hash:"name"`
	Address          Address   `redis-hash:"address,expand"`
	Age              uint      `redis-hash:"age"`
	BirthTime        time.Time `redis-hash:"birth"`
	NonTagged        string
	unexported       string
	unexportedTagged string `redis-hash:"unexported"`
}

type Address struct {
	Data   map[string]string `redis-hash:"data,expand"`
	Number int               `redis-hash:"number"`
	Main   bool              `redis-hash:"main"`
	City   *City             `redis-hash:"city,expand"`
}

type City struct {
	Name string `redis-hash:"name"`
}

type InvalidStruct struct {
	Name string `redis-hash:"name,expand"`
}

type InvalidInnerStruct struct {
	Data map[string]int `redis-hash:"data,expand"`
}
