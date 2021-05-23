package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

const PRODUCT_IMAGES_DIRECTORY = "./product_images"

type Server struct {
	port      string
	host      string
	router    *Router
	state     *State
	sessions  map[uint64]*User
	trasmisor *websocket.Upgrader
	ok        []byte
}

func (self *Server) broadcastProductChange(caller *User) {
	fmt.Println("Broadcasting product change by caller:", caller.username)
	var messages_sent uint = 0
	for _, user := range self.sessions {
		if user.id != caller.id && user.connection != nil {
			messages_sent++
			user.write(self.composeJson("message", "\"update\""))
		}
	}
	if messages_sent > 0 {
		fmt.Println("Messages sent:", messages_sent)
	}
}

func (self *Server) createSession(user *User, request *http.Request) uint64 {
	var session_key uint64 = shaAsInt64(fmt.Sprintf("%s:%s", user.username, request.RemoteAddr))
	self.sessions[session_key] = user
	return session_key
}

func (self *Server) createProductFromRequest(request *http.Request) *Product {
	var new_product *Product = new(Product)
	new_product.id = self.state.getNewProductId()
	new_product.name = request.FormValue("name")
	new_product.description = request.FormValue("description")
	new_product.stock = stringToInt(request.FormValue("stock"))
	new_product.price = stringToInt(request.FormValue("price"))
	new_product.images = make([]string, 0)
	return new_product
}

func (self *Server) createUserFromRequest(request *http.Request) *User {
	var new_user *User = new(User)
	new_user.id = self.state.getNewUserId()
	new_user.username = request.FormValue("username")
	new_user.name = request.FormValue("name")
	new_user.phone = request.FormValue("phone")
	new_user.email = request.FormValue("email")
	new_user.address = request.FormValue("address")
	new_user.password = shaAsInt64(request.FormValue("password"))
	new_user.cart = make([]*Purchase, 0)
	return new_user
}

func (self *Server) composeJson(key string, value string) []byte {
	return []byte(fmt.Sprintf("{\"%s\": %s}", key, value))
}

func (self *Server) composeResponse(response_value string) []byte {
	return self.composeJson("response", response_value)
}

func (self *Server) composeError(error_value string) []byte {
	return self.composeJson("error", error_value)
}

func (self *Server) composeHost() string {
	return fmt.Sprintf("%s:%s", self.host, self.port)
}

func (self *Server) enableCors(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Headers", "X-sk")
		response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE, CONNECT")
		if request.Method == http.MethodOptions {
			response.WriteHeader(200)
			response.Write(self.ok)
			return
		}

		handler(response, request)
	}
}

func (self *Server) greet(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	fmt.Fprintf(response, "hey")
}

func (self *Server) getSessionKey(request *http.Request) (uint64, error) {
	var sk string = request.Header.Get("X-sk")
	if sk != "" {
		var session_key uint64 = stringToUint64(sk)
		return session_key, nil
	} else {
		return 0, fmt.Errorf("Missing X-sk header")
	}
}

func (self *Server) getUserFromRequest(request *http.Request) (*User, error) {
	var sk string = request.Header.Get("X-sk")
	if sk != "" {
		var session_key uint64 = stringToUint64(sk)
		if user, exists := self.sessions[session_key]; exists {
			return user, nil
		} else {
			return nil, fmt.Errorf("Users doesnt exists")
		}
	} else {
		return nil, fmt.Errorf("The user your looking for is not where you think it is...")
	}
}

func (self *Server) getUserStash(response http.ResponseWriter, request *http.Request) {
	var user *User
	user, _ = self.getUserFromRequest(request)
	var user_stash int = self.state.calculateUserStash(user.id)
	response.WriteHeader(200)
	response.Write(self.composeJson("response", fmt.Sprint(user_stash)))
}

func (self *Server) handleCart(response http.ResponseWriter, request *http.Request) {
	var buyer *User
	buyer, _ = self.getUserFromRequest(request)
	switch request.Method {
	case http.MethodGet:
		var serialize_cart []string = make([]string, 0)
		for _, p := range buyer.cart {
			serialize_cart = append(serialize_cart, p.toJson())
		}
		response.WriteHeader(200)
		fmt.Fprintf(response, "[%s]", strings.Join(serialize_cart, ","))
	case http.MethodPost:
		var product_id uint = uint(stringToInt(request.FormValue("product_id")))
		var count int = stringToInt(request.FormValue("count"))
		var product *Product = self.state.getProductById(product_id)
		if product == nil {
			response.WriteHeader(400)
			response.Write(self.composeError("product doesnt exists"))
			return
		}
		buyer.cart = append(buyer.cart, &Purchase{product: product, count: count})

		fmt.Printf("%d '%s' was added to cart of %s\n", count, product.name, buyer.username)
		response.WriteHeader(200)
		response.Write(self.ok)
	case http.MethodDelete:
		// clearing cart
		var product_id uint = uint(stringToInt(request.FormValue("id")))
		var new_cart []*Purchase = make([]*Purchase, 0)
		for _, p := range buyer.cart {
			if p.product.id != product_id {
				new_cart = append(new_cart, p)
			}
		}
		buyer.cart = new_cart
		var cart_data []string = make([]string, 0)
		for _, p := range new_cart {
			cart_data = append(cart_data, p.toJson())
		}

		response.WriteHeader(200)
		fmt.Fprintf(response, "[%s]", strings.Join(cart_data, ","))
	case http.MethodPatch:
		// performing transacction
		for _, p := range buyer.cart {
			p.product.solds += p.count
		}
		buyer.cart = buyer.cart[:0]
		self.state.save()
		response.WriteHeader(200)
		response.Write(self.ok)
	case http.MethodOptions:
		response.WriteHeader(200)
		response.Write(self.ok)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write(self.composeError("nope"))
	}
}

func (self *Server) handleProductImages(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		var product_param string = request.FormValue("id")
		if product_param != "" {
			product := self.state.getProductById(uint(stringToInt(product_param)))
			if product != nil {
				var images_serialzed []byte
				images_serialzed, err := json.Marshal(product.images)
				if err != nil {
					fmt.Println("Error:", err.Error())
					response.WriteHeader(500)
					response.Write(self.composeError("server error, sorry for the inconvinece"))
					return
				}
				response.WriteHeader(200)
				response.Write(images_serialzed)
			}
		} else {
			response.WriteHeader(400)
			response.Write(self.composeError("missing product id"))
		}
	} else {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write(self.composeError("nope"))
	}
}

func (self *Server) handleProductsFeedSuscription(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:

		if request.URL.Query().Get("sk") != "" {
			var session_key string = request.URL.Query().Get("sk")
			request.Header.Add("X-sk", session_key)

			var user *User
			user, _ = self.getUserFromRequest(request)
			if user == nil {
				fmt.Printf("Session key %s yielded no valid session\n", session_key)
				response.WriteHeader(http.StatusNonAuthoritativeInfo)
				response.Write(self.composeError("invalid credentials"))

				return
			}

			var connection *websocket.Conn

			connection, err := self.trasmisor.Upgrade(response, request, nil)
			if err != nil {
				fmt.Println("Coulnt stablish a connection with")
			}

			fmt.Printf("user %s suscribed to products feed\n", user.username)
			user.connection = connection

		} else {
			response.WriteHeader(400)
			response.Write(self.composeError("missing credentials"))
		}

		// transmisor.Upgrade already responded
	case http.MethodDelete:
		var user *User
		user, _ = self.getUserFromRequest(request)
		if user.connection != nil {
			user.connection.Close()
			user.connection = nil
			fmt.Println("Unsuscribed user:", user.username)
		} else {
			response.WriteHeader(http.StatusNotAcceptable)
			response.Write(self.composeError("user was not suscribed"))
		}
	case http.MethodOptions:
		response.WriteHeader(200)
		response.Write(self.ok)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write(self.composeError("nope"))
	}
}

func (self *Server) handleRegister(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		var new_user *User = self.createUserFromRequest(request)
		fmt.Printf("New user '%s'\n", new_user.username)
		err := self.state.insertUser(new_user)
		if err == nil {
			response.WriteHeader(200)
			response.Write(self.ok)
		} else {
			fmt.Println("Username already exists")
			response.WriteHeader(http.StatusBadRequest)
			response.Write(self.composeError(err.Error()))
		}
	} else {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write(self.composeError(fmt.Sprintf("Method '%s' is not allowed", request.Method)))
	}
}

func (self *Server) handleProduct(response http.ResponseWriter, request *http.Request) {
	var vendor *User
	vendor, _ = self.getUserFromRequest(request)

	switch request.Method {
	case http.MethodGet:
		var user_selector string = request.FormValue("id")
		if user_selector == "*" {
			response.WriteHeader(200)
			fmt.Fprint(response, self.state.getAllProducts(vendor.id))
		} else if user_selector == "" {
			var vendor_products string = self.state.getProductsByVendor(vendor.id)
			response.WriteHeader(200)
			fmt.Fprint(response, vendor_products)
		}
	case http.MethodPost:
		var new_product *Product = self.createProductFromRequest(request)
		var filename string

		new_product.vendor = vendor.id
		var file_counter int = 0

		for {
			file, header, err := request.FormFile(fmt.Sprintf("image-%d", file_counter))
			if err != nil {
				break
			}
			filename = fmt.Sprintf("%s_%d_%d%s", new_product.name, vendor.id, file_counter, filepath.Ext(header.Filename))
			if err = self.saveFile(file, filename); err != nil {
				fmt.Println("Warning:", err.Error())
				continue
			}
			new_product.images = append(new_product.images, filename)
			file_counter++
		}

		self.state.insertProduct(new_product)
		response.WriteHeader(200)
		response.Write(self.ok)
		self.broadcastProductChange(vendor)
	case http.MethodPatch:
		if request.FormValue("product_id") != "" && request.FormValue("restock_by") != "" {
			var product_id uint = uint(stringToInt(request.FormValue("product_id")))
			var restock_by int = stringToInt(request.FormValue("restock_by"))
			var target_product *Product = self.state.getProductById(product_id)
			if target_product.vendor == vendor.id {
				fmt.Printf("restocking '%s' by %d\n", target_product.name, restock_by)
				target_product.stock += restock_by

				self.state.save()

				fmt.Println("succes!, stock is now:", target_product.stock)
				response.WriteHeader(200)
				response.Write(self.ok)
			} else {
				response.WriteHeader(http.StatusUnauthorized)
				response.Write(self.composeError("nope"))
			}

		} else {
			response.WriteHeader(http.StatusBadRequest)
			response.Write(self.composeError("missing parameters"))
		}
		self.broadcastProductChange(vendor)
	case http.MethodOptions:
		response.WriteHeader(200)
		response.Write(self.ok)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write(self.composeError("not allowed"))
	}
}

func (self *Server) handleUser(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		var session_key uint64
		session_key, _ = self.getSessionKey(request)
		user, _ := self.sessions[session_key]
		response.WriteHeader(200)
		response.Write([]byte(user.toJson()))
	case http.MethodOptions:
		response.WriteHeader(200)
		response.Write(self.ok)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(response, "Method '%s' not allowed", request.Method)
	}
}

func (self *Server) handleFileSystem(response http.ResponseWriter, request *http.Request) {
	var file_path string = request.URL.RequestURI()
	fmt.Println(file_path)
	response.WriteHeader(200)
	response.Write(self.ok)
}

func (self *Server) login(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		var username string = request.URL.Query().Get("username")
		var passwd string = request.URL.Query().Get("password")
		if username != "" && passwd != "" {
			var target *User = self.state.getUserByUsername(username)
			if target != nil {
				if target.password == shaAsInt64(passwd) {
					// correct login
					var session_key uint64 = self.createSession(target, request)
					fmt.Printf("Session created for user '%s' on '%s': %d\n", username, request.RemoteAddr, session_key)

					response.WriteHeader(200)
					response.Write(self.composeResponse(fmt.Sprintf("\"%d\"", session_key)))
				} else {
					fmt.Println("wrong password", target.password, "!==", shaAsInt64(passwd))
					//wrong password
					response.WriteHeader(http.StatusBadRequest)
					response.Write(self.composeError("\"wrong password\""))
				}
			} else {
				//wrong username
				response.WriteHeader(http.StatusNotFound)
				response.Write(self.composeError("\"user doesnt exists\""))
			}
		} else {
			// missing information
			fmt.Printf("Incomplete login credentials username='%s' password='%s'\n", username, passwd)
			response.WriteHeader(http.StatusBadRequest)
			response.Write(self.composeError("\"missing information\""))
		}
	}
}

func (self *Server) logout(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPatch {
		var session_key uint64
		session_key, _ = self.getSessionKey(request)
		delete(self.sessions, session_key)
		fmt.Printf("Session %d was finished\n", session_key)
		response.WriteHeader(200)
		response.Write(self.ok)
	} else {
		fmt.Println("Wrong method, user was not logged out")
		response.WriteHeader(http.StatusNotModified)
		response.Write(self.composeError("nope"))
	}
}

func (self *Server) registerUser(response http.ResponseWriter, request *http.Request) {
	if self.state.users.length == 0 {
		user_data := strings.Split(request.URL.Path, "-")
		var new_user *User = new(User)
		new_user.id = 0
		new_user.username = user_data[1]
		new_user.password = shaAsInt64(user_data[2])
		new_user.name = "root"
		fmt.Println("New root user:", new_user.toString())
		self.state.startState(new_user)
		response.WriteHeader(200)
		fmt.Fprintln(response, new_user.toJson())
	} else {
		response.WriteHeader(http.StatusForbidden)
		response.Write(self.composeError("\"forbidden\""))
	}
}

func (self *Server) sessionExists(sk string) bool {
	var session_key uint64 = stringToUint64(sk)
	_, exists := self.sessions[session_key]
	return exists
}

func (self *Server) setTrasmisor() *websocket.Upgrader {
	var new_transmisor *websocket.Upgrader = new(websocket.Upgrader)
	new_transmisor.ReadBufferSize = 512
	new_transmisor.WriteBufferSize = 512
	new_transmisor.CheckOrigin = func(r *http.Request) bool {
		fmt.Printf("\n%s %s%s %v\n", r.Method, r.Host, r.RequestURI, r.Proto)
		return true
	}

	return new_transmisor
}

func (self *Server) saveFile(file multipart.File, filename string) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if !pathExists(PRODUCT_IMAGES_DIRECTORY) {
		if err = os.Mkdir(PRODUCT_IMAGES_DIRECTORY, 0755); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(fmt.Sprintf("%s/%s", PRODUCT_IMAGES_DIRECTORY, filename), data, 0666)
}

func (self *Server) handleUnsuscribe() {

}

func (self *Server) validateSession(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		var session_key string = request.Header.Get("X-sk")
		fmt.Println("Session:", session_key)
		if session_key != "" && self.sessionExists(session_key) {
			handler(response, request)
		} else {
			response.WriteHeader(http.StatusForbidden)
			response.Write(self.composeError("Invalid or missing session key"))
		}
	}
}

func (self *Server) run() {
	if err := self.state.loadState(); err != nil {
		logFatal(err)
	}

	var file_system http.Handler = http.FileServer(http.Dir("./product_images"))

	self.router.registerRoute(NewRoute("/", true), self.enableCors(self.greet))
	self.router.registerRoute(NewRoute(`/register-root-[a-z\d]{3,8}-[a-z\d]+`, false), (self.registerUser))
	self.router.registerRoute(NewRoute("/register", true), self.enableCors(self.handleRegister))
	self.router.registerRoute(NewRoute("/products", true), self.enableCors(self.validateSession(self.handleProduct)))
	self.router.registerRoute(NewRoute("/products-feed", true), self.enableCors(self.handleProductsFeedSuscription))
	self.router.registerRoute(NewRoute("/products-images", true), self.enableCors(self.validateSession(self.handleProductImages)))
	self.router.registerRoute(NewRoute(`/cart`, true), self.enableCors(self.validateSession(self.handleCart)))
	self.router.registerRoute(NewRoute("/user", true), self.enableCors(self.validateSession(self.handleUser)))
	self.router.registerRoute(NewRoute("/user-stash", true), self.enableCors(self.validateSession(self.getUserStash)))
	self.router.registerRoute(NewRoute(`/login`, true), self.enableCors(self.login))
	self.router.registerRoute(NewRoute(`/logout`, true), self.enableCors(self.validateSession(self.logout)))
	self.router.registerRoute(NewRoute(`/static/.+`, false), self.enableCors(http.StripPrefix("/static/", file_system).ServeHTTP))

	fmt.Println("Lisiting on '", self.composeHost(), "'")
	http.ListenAndServe(self.composeHost(), self.router)
}

func createServer(port int) *Server {
	var new_server *Server = new(Server)
	var server_port string = os.Getenv("GSERVER_PORT")
	if server_port == "" {
		server_port = "5006"
	}
	var server_host string = os.Getenv("GSERVER_HOST")
	if server_host == "" {
		server_host = "127.0.0.1"
	}

	new_server.router = createRouter()
	new_server.state = createState()
	new_server.sessions = make(map[uint64]*User)
	new_server.host = server_host
	new_server.port = server_port
	new_server.trasmisor = new_server.setTrasmisor()
	new_server.ok = new_server.composeResponse("\"ok\"")

	return new_server
}

func main() {
	var server *Server = createServer(5006)
	server.run()
}
