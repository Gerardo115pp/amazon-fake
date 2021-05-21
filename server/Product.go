package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Product struct {
	id          uint
	name        string
	description string
	stock       int
	solds       int
	price       int
	vendor      uint
	images      []string
}

func (self *Product) String() string {
	return self.toString()
}

func (self *Product) getId() uint {
	return self.id
}

func (self *Product) compareString(other string) bool {
	return self.name == other
}

func (self *Product) isAvaliable() bool {
	return self.stock > self.solds
}

func (self *Product) toString() string {
	return self.name
}

func (self *Product) toRstring() string {
	product_images, err := json.Marshal(self.images)
	if err != nil {
		logFatal(err)
	}
	return fmt.Sprintf("%d|%s|%s|%d|%d|%d|%d|%s", self.id, self.name, self.description, self.stock, self.price, self.solds, self.vendor, string(product_images))
}

func (self *Product) load(rstring string) error {
	var values []string = strings.Split(rstring, "|")
	if len(values) == 8 {
		self.id = uint(stringToInt(values[0]))
		self.name = values[1]
		self.description = values[2]
		self.stock = stringToInt(values[3])
		self.price = stringToInt(values[4])
		self.solds = stringToInt(values[5])
		self.vendor = uint(stringToInt(values[6]))
		err := json.Unmarshal([]byte(values[7]), &(self.images))
		if err != nil {
			logFatal(err)
		}
		return nil
	} else {
		return fmt.Errorf("User requires 7 values but rstring '%s' had only %d", rstring, len(values))
	}
}

func (self *Product) toJson() string {
	return fmt.Sprintf("{ \"id\": \"%d\", \"name\": \"%s\", \"description\": \"%s\", \"stock\": \"%d\", \"solds\": \"%d\", \"price\": \"%d\"}", self.id, self.name, self.description, self.stock, self.solds, self.price)
}
