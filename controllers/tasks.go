package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nattatorn-dev/log-manager/common"
	"github.com/nattatorn-dev/log-manager/models"
	"gopkg.in/mgo.v2/bson"
)

// gets all tasks for a given user
func GetAllTasks(d models.TaskStore, w http.ResponseWriter, r *http.Request) {

	userClaims := r.Context().Value("userContext").(*common.UserClaims)

	tasks, err := d.GetAllTasksByUserId(userClaims.UserId)
	if err != nil {
		common.DisplayAppError(w, err, common.FetchError, http.StatusInternalServerError)
		return
	}

	common.WriteJson(w, "success", &tasks, http.StatusOK)
}

// Get a specific task for a given user
func GetTask(d models.TaskStore, w http.ResponseWriter, r *http.Request) {
	// get task_id
	taskId := mux.Vars(r)["id"]
	if !bson.IsObjectIdHex(taskId) {
		common.DisplayAppError(w, errors.New("taskId is not objectId"), common.FetchError, http.StatusInternalServerError)
		return
	}

	userClaims := r.Context().Value("userContext").(*common.UserClaims)
	task, err := d.GetTaskByUserIdAndTaskId(userClaims.UserId, bson.ObjectIdHex(taskId))
	if err != nil {
		common.DisplayAppError(w, err, common.FetchError, http.StatusInternalServerError)
		return
	}

	common.WriteJson(w, "success", &task, http.StatusOK)
}

func CreateTask(d models.TaskStore, w http.ResponseWriter, r *http.Request) {
	// get userId
	userClaims := r.Context().Value("userContext").(*common.UserClaims)
	userId := userClaims.UserId

	task := models.Task{} // dependent on models.Task, but doesn't come in through params :/
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		common.DisplayAppError(w, err, common.InvalidData, http.StatusInternalServerError)
		return
	}

	// creete the task
	task.UserId = userId
	createdTask, err := d.CreateTask(task)
	if err != nil {
		common.DisplayAppError(w, err, common.InvalidData, http.StatusInternalServerError)
		return
	}
	common.WriteJson(w, "success", &createdTask, http.StatusOK)

}

func DeleteTask(d models.TaskStore, w http.ResponseWriter, r *http.Request) {
	// get userId
	userClaims := r.Context().Value("userContext").(*common.UserClaims)
	userId := userClaims.UserId
	// get taskId
	taskId := mux.Vars(r)["id"]
	if !bson.IsObjectIdHex(taskId) {
		common.DisplayAppError(w, errors.New("taskId is not objectId"), common.FetchError, http.StatusInternalServerError)
		return
	}
	deletedTask := models.Task{}
	deletedTask, err := d.DeleteTaskByUserIdAndTaskId(userId, bson.ObjectIdHex(taskId))
	if err != nil {
		common.DisplayAppError(w, err, common.InvalidData, http.StatusInternalServerError)
		return
	}
	common.WriteJson(w, "success", &deletedTask, http.StatusOK)
}

func UpdateTask(d models.TaskStore, w http.ResponseWriter, r *http.Request) {
	// get userId
	userClaims := r.Context().Value("userContext").(*common.UserClaims)
	userId := userClaims.UserId
	// get taskId
	taskId := mux.Vars(r)["id"]
	if !bson.IsObjectIdHex(taskId) {
		common.DisplayAppError(w, errors.New("taskId is not objectId"), common.FetchError, http.StatusInternalServerError)
		return
	}
	taskIdObj := bson.ObjectIdHex(taskId)

	newTask := models.Task{} // dependent on models.Task, but doesn't come in through params :/
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		common.DisplayAppError(w, err, common.InvalidData, http.StatusInternalServerError)
		return
	}

	updatedTask := models.Task{}
	updatedTask, err = d.UpdateTaskByUserIdAndTaskID(newTask, userId, taskIdObj)
	if err != nil {
		fmt.Println("ERRORAAAAA")
		common.DisplayAppError(w, err, common.InvalidData, http.StatusInternalServerError)
		return
	}
	common.WriteJson(w, "success", &updatedTask, http.StatusOK)
}
