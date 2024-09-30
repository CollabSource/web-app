package routes

import (
	"fmt"
	"net/http"

	"github.com/bfbarry/CollabSource/back-end/controllers"
	"github.com/bfbarry/CollabSource/back-end/responseEntity"
)

type UserRouter struct {
	Router
	controller *controllers.UserController
}

var defaultUserRouter *UserRouter

func GetDefaultUserRouter() *UserRouter {
	return defaultUserRouter
}

func init() {
	defaultUserRouter = &UserRouter{}
	defaultUserRouter.controller = controllers.GetUserController()
	defaultUserRouter.initiateUserRoutes()
}

func (self *UserRouter) initiateUserRoutes() {
	endpoints := []Route{}
	endpoints = append(endpoints, Route{Path: fmt.Sprintf("%s/user/{id}", BASE_URL), Handler: self.user})
	endpoints = append(endpoints, Route{Path: fmt.Sprintf("%s/users", BASE_URL), Handler: self.users})
	endpoints = append(endpoints, Route{Path: fmt.Sprintf("%s/login", BASE_URL), Handler: self.login})
	endpoints = append(endpoints, Route{Path: fmt.Sprintf("%s/register", BASE_URL), Handler: self.register})
	endpoints = append(endpoints, Route{Path: fmt.Sprintf("%s/user_to_project/{id}", BASE_URL), Handler: self.userToProject})
	self.routes = endpoints
}

func (self *UserRouter) user(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	UUID := r.Header.Get("UUID")

	switch r.Method {
	case http.MethodGet:
		fmt.Println("GET /user")
		self.controller.GetUserByID(w, UUID, id)
	case http.MethodPatch:
		fmt.Println("PATCH /user")
		self.controller.UpdateUser(w, UUID, id, r)
	case http.MethodDelete:
		fmt.Println("DELETE /user")
		self.controller.DeleteUser(w, UUID, id)
	default:
		responseEntity.SendRequest(w, http.StatusMethodNotAllowed, []byte("Method Not Allowed"))
	}
}

func (self *UserRouter) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Println("GET /users")
		self.controller.GetUsers(w, r)
	case http.MethodPost:
		fmt.Println("POST /users")
		self.controller.GetUsersByQuery(w, r)
		//TODO Methods and structs for getting data by query
		// implementation: pass in multiple user IDs
		break
	default:
		responseEntity.SendRequest(w, http.StatusMethodNotAllowed, []byte("Method not allowed"))
	}
}

func (self *UserRouter) userToProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		fmt.Println("POST /user_to_project")
		self.controller.GetUserProjects(w, r, id)
	default:
		responseEntity.SendRequest(w, http.StatusMethodNotAllowed, []byte("Method not allowed"))
	}
}

func (self *UserRouter) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("POST /login")
		self.controller.Login(w, r)
	default:
		responseEntity.SendRequest(w, http.StatusMethodNotAllowed, []byte("Method not allowed"))
	}
}

func (self *UserRouter) register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("POST /register")
		self.controller.Register(w, r)
	default:
		responseEntity.SendRequest(w, http.StatusMethodNotAllowed, []byte("Method not allowed"))
	}
}

