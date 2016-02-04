package redis

type Person struct {
	ID               string  `redis-hash:"-"`
	Name             string  `redis-hash:"name"`
	Address          Address `redis-hash:",expand"`
	NonTagged        string
	unexported       string
	unexportedTagged string `redis-hash:"unexported"`
}

type Address struct {
	Data map[string]string `redis-hash:",expand"`
	City *City             `redis-hash:",expand"`
}

type City struct {
	Name string `redis-hash:"city_name"`
}

type InvalidStruct struct {
	Name string `redis-hash:",expand"`
}

type InvalidInnerStruct struct {
	Data map[string]int `redis-hash:",expand"`
}
