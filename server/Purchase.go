package main

import "fmt"

type Purchase struct {
	product *Product
	count   int
}

func (self *Purchase) toJson() string {
	var product_thumbnail string = ""
	if len(self.product.images) > 0 {
		product_thumbnail = self.product.images[0]
	}
	return fmt.Sprintf("{ \"product_id\": %d,\"thumbnail\": \"%s\", \"product_name\": \"%s\", \"unit_price\": %d, \"count\": %d }", self.product.id, product_thumbnail, self.product.name, self.product.price, self.count)
}
