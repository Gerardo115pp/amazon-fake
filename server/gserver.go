package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	port     string
	host     string
	router   *Router
	state    *State
	sessions map[uint64]*User
	ok       []byte
}

func (self *Server) createSession(user *User, request *http.Request) uint64 {
	var session_key uint64 = shaAsInt64(fmt.Sprintf("%s:%s", user.username, request.RemoteAddr))
	self.sessions[session_key] = user
	return session_key
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
		response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
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

	self.router.registerRoute(NewRoute("/", true), self.enableCors(self.greet))
	self.router.registerRoute(NewRoute(`/register-root-[a-z\d]{3,8}-[a-z\d]+`, false), (self.registerUser))
	self.router.registerRoute(NewRoute("/register", true), self.enableCors(self.handleRegister))
	self.router.registerRoute(NewRoute("/user", true), self.enableCors(self.validateSession(self.handleUser)))
	self.router.registerRoute(NewRoute(`/login`, true), self.enableCors(self.login))

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

	new_server.ok = new_server.composeResponse("\"ok\"")

	return new_server
}

func main() {
	var server *Server = createServer(5006)
	server.run()
}
