package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fmt"
	"reflect"

	"github.com/bfbarry/CollabSource/back-end/model"
	"github.com/bfbarry/CollabSource/back-end/repository"
	"github.com/bfbarry/CollabSource/back-end/responseEntity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const USER_COLLECTION = "users"

type UserController struct {
	repository *repository.Repository
}

var defaultUserController *UserController

func GetUserController() *UserController {
	return defaultUserController
}

func init() {
	defaultUserController = &UserController{repository: repository.GetMongoRepository()}
}

func (self *UserController) Register(w http.ResponseWriter, r *http.Request) {
	userEntity := model.User{}

	if err := json.NewDecoder(r.Body).Decode(&userEntity); err != nil {
		responseEntity.SendRequest(w, http.StatusBadRequest, []byte("Invalid JSON"))
		return
	}

	// TODO check skills array is empty
	// TODO create function for business logic readability e.g., SatisfiesRecommender()
	if userEntity.Name == "" || userEntity.Description == "" ||
		userEntity.Email == "" || userEntity.Password == "" {
		responseEntity.SendRequest(w, http.StatusUnprocessableEntity, []byte("Invalid payload"))
		return
	}

	emailFilter := bson.M{"email": userEntity.Email}
	userDummy := &model.User{}
	if err := self.repository.FindOne(USER_COLLECTION, emailFilter, userDummy); err == nil {
		responseEntity.SendRequest(w, http.StatusUnprocessableEntity, []byte("email exists"))
		return
	}

	id, err := self.repository.Insert(USER_COLLECTION, userEntity)
	if err != nil {
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("server error on insert"))
		return
	}

	responseEntity.SendRequest(w, http.StatusOK, []byte(id.Hex()))
}

func (self *UserController) Login(w http.ResponseWriter, r *http.Request) {
	// TODO model field constraints
	loginFields := model.LoginFields{}
	user := &model.User{}

	if err := json.NewDecoder(r.Body).Decode(&loginFields); err != nil {
		responseEntity.SendRequest(w, http.StatusBadRequest, []byte("Invalid JSON"))
		return
	}

	filter := bson.M{"email": loginFields.Email}
	err := self.repository.FindOne(USER_COLLECTION, filter, user)
	if err != nil {
		responseEntity.SendRequest(w, http.StatusUnauthorized, []byte("incorrect email"))
		fmt.Println("bad email", loginFields.Email)
		return
	}

	if user.Password != loginFields.Password {
		fmt.Println(user.Password, loginFields.Password)
		responseEntity.SendRequest(w, http.StatusUnauthorized, []byte("incorrect password"))
		return
	}

	responseEntity.SendRequest(w, http.StatusOK, []byte(user.Id.Hex()))
}

func (self *UserController) GetUserByID(w http.ResponseWriter, userUUID string, id string) {
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseEntity.SendRequest(w, http.StatusUnprocessableEntity, []byte("Invalid Object ID"))
		return
	}
	var userEntity interface{}
	var reflector interface{}
	if userUUID == id {
		userEntity = &model.User{}
		reflector = &model.User{}
	} else {
		userEntity = &model.PublicUser{}
		reflector = &model.PublicUser{}
	}

	mongoErr := self.repository.FindByID(USER_COLLECTION, ObjId, userEntity)
	if reflect.DeepEqual(userEntity, reflector) {
		responseEntity.SendRequest(w, http.StatusNotFound, []byte("not found"))
		return
	}
	if mongoErr != nil {
		fmt.Println(mongoErr)
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	jsonRes, jsonerr := json.Marshal(userEntity)
	if jsonerr != nil {
		fmt.Println("json err")
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}
	responseEntity.SendRequest(w, http.StatusOK, jsonRes)
}

func (self *UserController) GetUsersByQuery(w http.ResponseWriter, r *http.Request) {
	// TODO check for nonexistent IDs
	defaultPageNum := 1
	defaultPageSize := 10
	var err error
	query := model.UserPostQuery{}

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		responseEntity.SendRequest(w, http.StatusBadRequest, []byte("Invalid JSON"))
		return
	}
	if len(query.IDs) < 1 {
		responseEntity.SendRequest(w, http.StatusBadRequest, []byte("IDs cannot be empty"))
		return
	}
	var filter bson.M
	filter, err = postQueryToBsonM(query)
	if err != nil {
		responseEntity.SendRequest(w, http.StatusBadRequest, []byte("Invalid JSON"))
		return
	}

	var pageNum int
	var pageSize int

	queryParams := r.URL.Query()

	if pageNum, err = strconv.Atoi(queryParams.Get("page")); err != nil {
		pageNum = defaultPageNum
	}

	if pageSize, err = strconv.Atoi(queryParams.Get("size")); err != nil {
		pageSize = defaultPageSize
	}

	var userEntities []model.PublicUser

	hasNext, mongoErr := self.repository.FindManyByPage(USER_COLLECTION, &userEntities, pageNum, pageSize, filter)
	if mongoErr != nil {
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	if err != nil {
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	res := responseEntity.PaginatedResponseBody[model.PublicUser]{
		Items: userEntities,
		Page: pageNum,
		HasNext: hasNext,
	}

	jsonRes, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	responseEntity.SendRequest(w, http.StatusOK, jsonRes)
}

func (self *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {

	// Why do we have this?  Should just need GetUsersByQuery

	defaultPageNum := 1
	defaultPageSize := 10
	var err error

	var pageNum int
	var pageSize int

	queryParams := r.URL.Query()

	if pageNum, err = strconv.Atoi(queryParams.Get("page")); err != nil {
		pageNum = defaultPageNum
	}

	if pageSize, err = strconv.Atoi(queryParams.Get("size")); err != nil {
		pageNum = defaultPageSize
	}

	var userEntities []model.PublicUser
	filter := bson.M{}

	hasNext, mongoErr := self.repository.FindManyByPage(USER_COLLECTION, &userEntities, pageNum, pageSize, filter)
	if mongoErr != nil {
		fmt.Println(mongoErr)
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}
	
	if err != nil {
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	res := responseEntity.PaginatedResponseBody[model.PublicUser]{
		Items: userEntities,
		Page: pageNum,
		HasNext: hasNext,
	}

	jsonRes, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	responseEntity.SendRequest(w, http.StatusOK, jsonRes)
}

func (self *UserController) UpdateUser(w http.ResponseWriter, userUUID string, id string, r *http.Request) {
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseEntity.SendRequest(w, http.StatusUnprocessableEntity, []byte("Invalid Object ID"))
		return
	}
	
	userCheck := &model.IdObj{}
	var mongoErr error
	mongoErr = self.repository.FindByID(USER_COLLECTION, ObjId, userCheck)
	if reflect.DeepEqual(*userCheck, model.IdObj{}) {
		responseEntity.SendRequest(w, http.StatusNotFound, []byte("not found"))
		return
	}
	if mongoErr != nil {
		fmt.Println(mongoErr)
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	if userCheck.Id.Hex() != userUUID {
		responseEntity.SendRequest(w, http.StatusUnauthorized, []byte("unauthorized"))
		return
	}
	userEntity := model.User{}

	if err := json.NewDecoder(r.Body).Decode(&userEntity); err != nil {
		responseEntity.SendRequest(w, http.StatusBadRequest, []byte("Invalid JSON"))
		return
	}

	// TODO other controller for LoginFields
	if userEntity.Password != "" {
		fmt.Println("cannot change password in UpdateUser")
		responseEntity.SendRequest(w, http.StatusBadRequest, []byte("Invalid Json"))
		return
	}
	var updatedCount int64
	updatedCount, mongoErr = self.repository.Update(USER_COLLECTION, ObjId, userEntity)
	if mongoErr != nil {
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	if updatedCount == 0 {
		responseEntity.SendRequest(w, http.StatusNoContent, []byte("no change"))
		return
	}

	responseEntity.SendRequest(w, http.StatusOK, []byte("success"))
}

func (self *UserController) DeleteUser(w http.ResponseWriter, userUUID string, id string) {
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseEntity.SendRequest(w, http.StatusUnprocessableEntity, []byte("Invalid Object ID"))
		return
	}

	userCheck := &model.IdObj{}
	var mongoErr error
	mongoErr = self.repository.FindByID(USER_COLLECTION, ObjId, userCheck)
	if reflect.DeepEqual(*userCheck, model.IdObj{}) {
		responseEntity.SendRequest(w, http.StatusNotFound, []byte("not found"))
		return
	}
	if mongoErr != nil {
		fmt.Println(mongoErr)
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}
	if userCheck.Id.String() != userUUID {
		responseEntity.SendRequest(w, http.StatusUnauthorized, []byte("unauthorized"))
		return
	}
	// deleteModeStr := r.URL.Query().Get("mode") // TODO separate hard and soft delete in repository.go
	// deleteMode := repository.Str2Enum(deleteModeStr)
	deletedCount, mongoErr := self.repository.Delete(USER_COLLECTION, ObjId)
	if mongoErr != nil {
		fmt.Println(mongoErr)
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	if deletedCount == 0 {
		responseEntity.SendRequest(w, http.StatusBadRequest, []byte("ID Not Found"))
		return
	}

	responseEntity.SendRequest(w, http.StatusOK, []byte("Success"))
}

func (self *UserController) GetUserProjects(w http.ResponseWriter, r *http.Request, id string) {
	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseEntity.SendRequest(w, http.StatusUnprocessableEntity, []byte("Invalid Object ID"))
		return
	}
	var pageNum int
	var pageSize int

	queryParams := r.URL.Query()
	defaultPageNum := 1
	defaultPageSize := 10
	if pageNum, err = strconv.Atoi(queryParams.Get("page")); err != nil {
		pageNum = defaultPageNum
	}

	if pageSize, err = strconv.Atoi(queryParams.Get("size")); err != nil {
		pageSize = defaultPageSize
	}

	var entities []model.Project
	hasNext, mongoErr := self.repository.FindManyByJunction(USER_PROJECT_COLLECTION, "user_id", ObjId, 
															"project_id", "projects", pageNum, pageSize, &entities)
	if mongoErr != nil {
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	res := responseEntity.PaginatedResponseBody[model.Project]{
		Items: entities,
		Page: pageNum,
		HasNext: hasNext,
	}

	jsonRes, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		responseEntity.SendRequest(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	responseEntity.SendRequest(w, http.StatusOK, jsonRes)
}