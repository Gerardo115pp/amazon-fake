package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const USERS_TAG = "users"
const PRODUCTS_TAG = "products"
const SERVER_DATA_TAG = "server_data"

type State struct {
	users           *List
	products        *List
	storage_file    *RandomAccessFile
	last_user_id    uint
	last_product_id uint
}

func (self *State) clearState() {
	self.users.clear()
}

func (self *State) composeTag(content string, tag_name string) string {
	return fmt.Sprintf("<%s>%s</%s>", tag_name, content, tag_name)
}

func (self *State) calculateUserStash(user_id uint) int {
	var stash int = 0
	var current_product *Product
	for current_node := self.products.root; current_node != nil; current_node = current_node.Next {
		current_product = current_node.NodeContent.(*Product)
		if current_product.vendor == user_id {
			stash += current_product.solds * current_product.price
		}
	}
	return stash
}

func (self *State) getAllProducts(vendor_id uint) string {
	var all_products []string
	for currrent_node := self.products.root; currrent_node != nil; currrent_node = currrent_node.Next {
		if currrent_node.NodeContent.(*Product).isAvaliable() && currrent_node.NodeContent.(*Product).vendor != vendor_id {
			all_products = append(all_products, currrent_node.NodeContent.toJson())
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(all_products, ","))
}

func (self *State) getPairTags(tag_name string) (*MatchRange, *MatchRange) {
	open_tag := self.storage_file.getRange(fmt.Sprintf("<%s>", tag_name))
	close_tag := self.storage_file.getRange(fmt.Sprintf("</%s>", tag_name))
	return open_tag, close_tag
}

func (self *State) getNewUserId() uint {
	self.last_user_id++
	return self.last_user_id
}

func (self *State) getNewProductId() uint {
	self.last_product_id++
	return self.last_product_id
}

func (self *State) getTagContent(tag_name string) string {
	var open_tag, close_tag *MatchRange = self.getPairTags(tag_name)
	var tag_content string = string(self.storage_file.readFrom(open_tag.right+1, close_tag.left))
	return tag_content
}

func (self *State) getProductById(product_id uint) *Product {
	for current_node := self.products.root; current_node != nil; current_node = current_node.Next {
		if current_node.NodeContent.getId() == product_id {
			return current_node.NodeContent.(*Product)
		}
	}
	return nil
}

func (self *State) getUserByUsername(username string) *User {
	c := self.users.exists(username)
	if c == nil {
		return nil
	}
	return c.(*User)
}

func (self *State) getProductsByVendor(vendor_id uint) string {
	var vendor_products []string
	for current_node := self.products.root; current_node != nil; current_node = current_node.Next {
		if current_node.NodeContent.(*Product).vendor == vendor_id {
			vendor_products = append(vendor_products, current_node.NodeContent.toJson())
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(vendor_products, ","))
}

func (self *State) insertUser(new_user *User) error {
	if other := self.users.exists(new_user.toString()); other == nil {
		self.users.append(new_user)
		if err := self.saveTag(self.users, USERS_TAG); err != nil {
			logFatal(err)
		}
		return nil
	} else {
		return fmt.Errorf("User '%s' already exists", new_user.toString())
	}
}

func (self *State) insertProduct(new_product *Product) {
	self.products.append(new_product)
	if err := self.saveTag(self.products, PRODUCTS_TAG); err != nil {
		logFatal(err)
	}
}

func (self *State) loadState() error {
	var err error
	if err = self.loadUsers(); err != nil {
		return err
	} else if err = self.loadServerData(); err != nil {
		return err
	} else if err = self.loadProducts(); err != nil {
		logFatal(err)
	}
	return err
}

func (self *State) loadServerData() error {
	var server_data string = self.getTagContent(SERVER_DATA_TAG)
	server_data_json := &struct {
		LastUserID int `json:"last_user_id"`
		LastItemID int `json:"last_item_id"`
	}{}
	err := json.Unmarshal([]byte(server_data), server_data_json)
	self.last_user_id = uint(server_data_json.LastUserID)
	self.last_product_id = uint(server_data_json.LastItemID)
	return err
}

func (self *State) loadUsers() (err error) {
	var users_data string = self.getTagContent(USERS_TAG)
	var current_user *User
	if users_data != "" {
		for _, rstring := range strings.Split(users_data, "*") {
			current_user = new(User)
			if err = current_user.load(rstring); err != nil {
				fmt.Printf("Warning rstring '%s' couldnt be loaded...\n", rstring)
				continue
			}
			self.users.append(current_user)
		}
	}
	fmt.Println("Users loaded:", self.users.length)
	return
}

func (self *State) loadProducts() (err error) {
	var users_data string = self.getTagContent(PRODUCTS_TAG)
	var current_product *Product
	if users_data != "" {
		for _, rstring := range strings.Split(users_data, "*") {
			current_product = new(Product)
			if err = current_product.load(rstring); err != nil {
				fmt.Printf("Warning rstring '%s' couldnt be loaded...\n", rstring)
				continue
			}
			self.products.append(current_product)
		}
	}
	fmt.Println("Products loaded:", self.products.length)
	return
}

func (self *State) startState(root_user *User) {
	self.clearState()
	self.users.append(root_user)
	self.storage_file.clear()
}

func (self *State) save() {
	// saveing users
	if err := self.saveTag(self.users, USERS_TAG); err != nil {
		logFatal(err)
	}
	// saveing products
	if err := self.saveTag(self.products, PRODUCTS_TAG); err != nil {
		logFatal(err)
	}
}

func (self *State) saveServerData() error {
	var open_tag, close_tag *MatchRange = self.getPairTags(SERVER_DATA_TAG)
	server_data := &struct {
		LastUserID int `json:"last_user_id"`
		LastItemID int `json:"last_item_id"`
	}{LastUserID: int(self.last_user_id), LastItemID: int(self.last_product_id)}
	json_data, err := json.Marshal(server_data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	self.saveContentToRAF(open_tag, close_tag, self.composeTag(string(json_data), SERVER_DATA_TAG))
	return nil
}

func (self *State) saveTag(state_list *List, tag_name string) error {
	fmt.Printf("Saving %d %s\n", self.users.length, tag_name)
	var open_tag, close_tag *MatchRange = self.getPairTags(tag_name)
	if open_tag.right != -1 && close_tag.right != -1 {
		var tag_serialized []string = state_list.mapFunc(func(ln *ListNode) string { return ln.NodeContent.toRstring() })

		self.saveServerData()
		return self.saveContentToRAF(open_tag, close_tag, self.composeTag(strings.Join(tag_serialized, "*"), tag_name))
	} else {
		return fmt.Errorf("tag '%s' doesnt exists", tag_name)
	}
}

func (self *State) saveContentToRAF(open_tag *MatchRange, close_tag *MatchRange, content string) error {
	var previous_content string
	var remaing_content string
	//getting previous content
	if open_tag.left != 0 {
		previous_content = string(self.storage_file.readFrom(0, open_tag.left))
	} else {
		previous_content = ""
	}

	// getting the content on the right of content
	if close_tag.right < self.storage_file.Size() {
		remaing_content = string(self.storage_file.readFrom(close_tag.right+1, self.storage_file.Size()))
	} else {
		remaing_content = ""
	}

	//saveing all the content
	var total_content string = fmt.Sprintf("%s%s%s", previous_content, content, remaing_content)
	self.storage_file.truncate(int64(len(total_content)), true)
	self.storage_file.seek(0)
	return self.storage_file.write(total_content)
}

func createState() *State {
	var new_state *State = new(State)
	new_state.users = new(List)
	new_state.products = new(List)
	new_state.storage_file = createRAF("store")

	return new_state
}
