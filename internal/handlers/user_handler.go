package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	. "go-service/internal/models"
	. "go-service/internal/services"
	"net/http"
	"reflect"
	"strings"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	result, err := h.service.Load(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Insert(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("The insert is now running")
	result, er2 := h.service.Insert(r.Context(), &user)
	if er2 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}

	result, er2 := h.service.Update(r.Context(), &user)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	//ids := []string{"id"}

	var user User
	//userType := reflect.TypeOf(user)
	//jsonMap := sv.BuildMapField(userType)

	body, _ := BuildMapAndStruct(r, &user)

	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}

	/*json, er1, _ := BodyToJson(r, user, body, ids, jsonMap, nil)

	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}*/

	result, er2 := h.service.Patch(r.Context(), user.Id, body)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	result, err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, result)
}

func respond(w http.ResponseWriter, result interface{}) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

type CModelBuilder interface {
	BuildToInsert(ctx context.Context, model interface{}) (interface{}, error)
	BuildToUpdate(ctx context.Context, model interface{}) (interface{}, error)
	BuildToPatch(ctx context.Context, model interface{}) (interface{}, error)
	BuildToSave(ctx context.Context, model interface{}) (interface{}, error)
}

func GetValue(model interface{}, index int) (interface{}, string, error) {
	valueObject := reflect.Indirect(reflect.ValueOf(model))
	return reflect.Indirect(valueObject.Field(index)).Interface(), valueObject.Type().Field(index).Name, nil
}

func BuildMapAndStruct(r *http.Request, interfaceBody interface{}) (map[string]interface{}, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	s := buf.String()
	body := make(map[string]interface{})
	er1 := json.NewDecoder(strings.NewReader(s)).Decode(&body)
	if er1 != nil {
		return nil, er1
	}
	er2 := json.NewDecoder(strings.NewReader(s)).Decode(interfaceBody)
	if er2 != nil {
		return nil, er2
	}
	return body, nil
}

func BodyToJson(r *http.Request, structBody interface{}, body map[string]interface{}, jsonIds []string, mapIndex map[string]int, modelBuilder CModelBuilder) (map[string]interface{}, error, error) {
	var controlModel interface{}
	if modelBuilder != nil {
		var er0 error
		controlModel, er0 = modelBuilder.BuildToPatch(r.Context(), structBody)
		if er0 != nil {
			return nil, er0, nil
		}
		inRec, er1 := json.Marshal(controlModel)
		if er1 != nil {
			return nil, nil, er1
		}
		var model map[string]interface{}
		json.Unmarshal(inRec, &model)
		for k, v := range model {
			stringKind := reflect.TypeOf(v).String()
			if (v != nil && stringKind == "float64" && v.(float64) != 0) || (v != nil && stringKind != "float64" && v != "") {
				body[k] = v
			}
		}
	}
	valueOfReq := reflect.ValueOf(structBody)
	if valueOfReq.Kind() == reflect.Ptr {
		valueOfReq = reflect.Indirect(valueOfReq)
	}
	for _, jsonName := range jsonIds {
		if i, ok := mapIndex[jsonName]; ok && i >= 0 {
			v, _, er4 := GetValue(structBody, i)
			if er4 == nil {
				body[jsonName] = v
			}
		}
	}
	result := make(map[string]interface{})
	for keyJsonName, _ := range body {
		v2 := body[keyJsonName]
		if v2 == nil {
			result[keyJsonName] = v2
		} else if i, ok := mapIndex[keyJsonName]; ok && i >= 0 {
			v, _, er4 := GetValue(structBody, i)
			if er4 == nil {
				result[keyJsonName] = v
			}
		}
	}
	return result, nil, nil
}