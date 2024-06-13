package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	// "github.com/bfbarry/CollabSource/back-end/errors"
	"github.com/bfbarry/CollabSource/back-end/model"
	"github.com/bfbarry/CollabSource/back-end/repository"
	"github.com/bfbarry/CollabSource/back-end/responseEntity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PROJECT_COLLECTION = "projects"

type ProjectController struct {
	repository *repository.Repository
}

var defaultProjectController *ProjectController

func GetProjectController() *ProjectController {
	return defaultProjectController
}

func init() {
	defaultProjectController = &ProjectController{repository: repository.GetMongoRepository()}
}

func (self *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {

	streamObj := r.Body
	projectEntity := model.Project{}
	if err := json.NewDecoder(streamObj).Decode(&projectEntity); err != nil {
		responseEntity.ResponseEntity(w, http.StatusBadRequest, []byte("Invalid JSON"))
		return
	}

	if projectEntity.Name == "" || projectEntity.Description == "" {
		responseEntity.ResponseEntity(w, http.StatusUnprocessableEntity, []byte("Invalid payload"))
		return
	}

	if err := self.repository.Insert(PROJECT_COLLECTION, projectEntity); err != nil {
		responseEntity.ResponseEntity(w, http.StatusInternalServerError, []byte("Error"))
		return
	}

	responseEntity.ResponseEntity(w, http.StatusOK, []byte("Success"))
}

func (self *ProjectController) GetProjectByID(w http.ResponseWriter, id string) {
	// var op errors.Op = "controllers.GetProjectByID"
	projectEntity := &model.Project{}

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseEntity.ResponseEntity(w, http.StatusUnprocessableEntity, []byte("Invalid Object ID"))
		return
	}

	result, mongoErr := self.repository.FindByID(PROJECT_COLLECTION, ObjId, projectEntity)
	if mongoErr != nil {
		responseEntity.ResponseEntity(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}
	if result == nil {
		responseEntity.ResponseEntity(w, http.StatusBadRequest, []byte("ID Not Found"))
		return
	}

	jsonResponse, jsonerr := json.Marshal(result)
	if jsonerr != nil { 
		responseEntity.ResponseEntity(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	responseEntity.ResponseEntity(w, http.StatusOK, jsonResponse)
}

func (self *ProjectController) UpdateProject(w http.ResponseWriter, id string, r *http.Request) {

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseEntity.ResponseEntity(w, http.StatusUnprocessableEntity, []byte("Invalid Object ID"))
		return
	}

	streamObj := r.Body
	projectEntity := model.Project{}
	if err := json.NewDecoder(streamObj).Decode(&projectEntity); err != nil {
		responseEntity.ResponseEntity(w, http.StatusBadRequest, []byte("Invalid JSON"))
		return
	}

	updatedCount, mongoErr := self.repository.Update(PROJECT_COLLECTION, ObjId, projectEntity)
	if mongoErr != nil {
		responseEntity.ResponseEntity(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	if updatedCount == 0 {
		responseEntity.ResponseEntity(w, http.StatusBadRequest, []byte("ID Not Found"))
		return
	}
	
	responseEntity.ResponseEntity(w, http.StatusOK, []byte("success"))
}

func (self *ProjectController) DeleteProject(w http.ResponseWriter, id string) {
	// TODO pass in reader to get URL param

	ObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseEntity.ResponseEntity(w, http.StatusUnprocessableEntity, []byte("Invalid Object ID"))
		return
	}

	// deleteModeStr := r.URL.Query().Get("mode") // TODO separate hard and soft delete in repository.go
	// deleteMode := repository.Str2Enum(deleteModeStr)
	deletedCount, mongoErr := self.repository.Delete(PROJECT_COLLECTION, ObjId)
	if mongoErr != nil {
		responseEntity.ResponseEntity(w, http.StatusInternalServerError, []byte("Something went wrong"))
		return
	}

	if deletedCount == 0{
		responseEntity.ResponseEntity(w, http.StatusBadRequest, []byte("ID Not Found"))
		return
	}
	
	responseEntity.ResponseEntity(w, http.StatusOK, []byte("Success"))
}

func (self *ProjectController) GetProjectsByFilter(w http.ResponseWriter, streamFilterObj *io.ReadCloser, pageNumber int64, pageSize int64) {
	// var op errors.Op = "controllers.GetProjectsByFilter"
	results, err := self.repository.Find(PROJECT_COLLECTION, streamFilterObj, 0, 10)
	if err != nil {
		// return nil, err
	}
	// var jsonErr error
	jsonResponse, _ := json.Marshal(results)
	if err != nil { // TODO: 500
		// return nil, errors.E(jsonErr, http.StatusInternalServerError, op, "json marshall error")
	}
	responseEntity.ResponseEntity(w, http.StatusOK, jsonResponse)
	// return jsonResponse, nil
}