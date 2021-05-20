package main

import (
	"fmt"
	"strings"
)

type User struct {
	id       uint
	username string
	name     string
	phone    string
	email    string
	address  string
	password uint64
}

func (self *User) String() string {
	return self.toString()
}

func (self *User) getId() uint {
	return self.id
}

func (self *User) compareString(other string) bool {
	return self.username == other
}

func (self *User) toString() string {
	return self.username
}

func (self *User) toRstring() string {
	return fmt.Sprintf("%d|%s|%s|%s|%s|%s|%d", self.id, self.username, self.name, self.phone, self.email, self.address, self.password)
}

func (self *User) load(rstring string) error {
	var values []string = strings.Split(rstring, "|")
	if len(values) == 7 {
		self.id = uint(stringToInt(values[0]))
		self.username = values[1]
		self.name = values[2]
		self.phone = values[3]
		self.email = values[4]
		self.address = values[5]
		self.password = stringToUint64(values[6])
		return nil
	} else {
		return fmt.Errorf("User requires 7 values but rstring '%s' had only %d", rstring, len(values))
	}
}

func (self *User) toJson() string {
	return fmt.Sprintf("{ \"username\": \"%s\", \"name\": \"%s\", \"phone\": \"%s\", \"email\": \"%s\", \"address\": \"%s\"}", self.username, self.name, self.phone, self.email, self.address)
}
