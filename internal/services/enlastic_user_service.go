package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	_ "github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"go-service/internal/models"
	_ "go-service/internal/models"
	"reflect"
	"strings"
)

type ElasticUserService struct {
	elastic *elasticsearch.Client
}

func NewEUserService(ela *elasticsearch.Client) *ElasticUserService {
	return &ElasticUserService{elastic: ela}
}

func convertDocToJson(doc interface{}) string {
	jsonString, err := json.Marshal(doc)

	if err != nil {
		fmt.Println("An error is happening when encoded the new user: ", err)
		return ""
	}
	return string(jsonString)
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	/*structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}*/
	val := reflect.ValueOf(value)
	structFieldValue.Set(val)
	return nil
}

func DecodeMapToStruct(m map[string]interface{}) (u *models.User)  {
	for k, v := range m {
		err := SetField(u, k, v)
		if err != nil {
			return nil
		}
		return u
	}
	return nil
}

func transcode(in, out interface{})  {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}

func (e *ElasticUserService) Insert(ctx context.Context, user *models.User) (int64, error) {
	if user == nil {
		fmt.Print("Can not add null user")
		return 0, nil
	}

	userJsonString := convertDocToJson(user)
	request := esapi.IndexRequest{
		Index: "users",
		DocumentID: user.Id,
		Body: strings.NewReader(userJsonString),
		Refresh: "true",
	}
	response, err := request.Do(ctx, e.elastic)

	if err != nil {
		panic(err)
		return 0, nil
	}

	defer response.Body.Close()

	var result map[string]interface{}

	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		panic(err)
		return 0, nil
	}

	fmt.Println("IndexRequest to insert Status: ", response.Status())
	fmt.Println("Result: ", result["result"])

	fmt.Print("The new user %v has been added successfully.", user.Username)
	return 1, nil
}

func (e *ElasticUserService) GetAll(ctx context.Context) (*[]models.User, error) {
	var listUser []models.User
	var mapResponse map[string]interface{}
	var buf bytes.Buffer

	query := `{
  "query": {
    "match_all": {}
  },
  "size": 1
}`

	var queryString = strings.NewReader(query)

	err := json.NewEncoder(&buf).Encode(queryString)
	if err != nil {
		fmt.Print("Error during encoding the query : %s", err.Error())
	}

	result, err := e.elastic.Search(
		e.elastic.Search.WithContext(ctx),
		e.elastic.Search.WithIndex("users"),
		e.elastic.Search.WithBody(queryString),
		e.elastic.Search.WithTrackTotalHits(true),
		e.elastic.Search.WithPretty(),
		)
	defer result.Body.Close()

	err = json.NewDecoder(result.Body).Decode(&mapResponse)
	fmt.Println("This is map response: ", mapResponse)
	if err != nil {
		fmt.Println("Error parsing the result to User type:", err.Error())
	}

	var u = &models.User{}
	//var decodedUser models.User{}
	for _, hit := range mapResponse["hits"].(map[string]interface{})["hits"].([]interface{}) {
		user := hit.(map[string]interface{})

		source := user["_source"]
		//userId := user["_id"]
		//u = DecodeMapToStruct(user)
		fmt.Println("This is the source:")

		fmt.Println(source)
		listUser = append(listUser, *u)

	}
	return &listUser, nil
}


func (e *ElasticUserService) Load(ctx context.Context, id string) (*models.User, error) {
	var listUser []models.User
	var mapResponse map[string]interface{}
	var buf bytes.Buffer

	query := `{
  "query": {
    "match": { "id": "{0}" }
  },
  "size": 1
}`

	query = strings.Replace(query, "{0}", id, 1)
	var queryString = strings.NewReader(query)

	err := json.NewEncoder(&buf).Encode(queryString)
	if err != nil {
		fmt.Print("Error during encoding the query : %s", err.Error())
	}

	result, err := e.elastic.Search(
		e.elastic.Search.WithContext(ctx),
		e.elastic.Search.WithIndex("users"),
		e.elastic.Search.WithBody(queryString),
		e.elastic.Search.WithTrackTotalHits(true),
		e.elastic.Search.WithPretty(),
	)
	defer result.Body.Close()

	err = json.NewDecoder(result.Body).Decode(&mapResponse)
	fmt.Println("This is map response: ", mapResponse)
	if err != nil {
		fmt.Println("Error parsing the result to User type:", err.Error())
	}

	var u = &models.User{}
	//var decodedUser models.User{}
	for _, hit := range mapResponse["hits"].(map[string]interface{})["hits"].([]interface{}) {
		user := hit.(map[string]interface{})

		source := user["_source"]
		//userId := user["_id"]
		//u = DecodeMapToStruct(user)
		fmt.Println("This is the source:")

		fmt.Println(source)
		listUser = append(listUser, *u)

	}
	return u, nil
}

func (e *ElasticUserService) Update(ctx context.Context, user *models.User) (int64, error) {
	userJsonString := convertDocToJson(user)
	request := esapi.UpdateRequest{
		Index: "users",
		DocumentID: user.Id,
		Body: strings.NewReader(userJsonString),
		Refresh: "true",
	}
	response, err := request.Do(ctx, e.elastic)

	if err != nil {
		panic(err)
		return 0, nil
	}

	defer response.Body.Close()

	var result map[string]interface{}

	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		panic(err)
		return 0, nil
	}

	fmt.Println("IndexRequest to update Status: ", response.Status())
	fmt.Println("Result: ", result["result"])

	fmt.Print("The user %v has been updated successfully.", user.Username)
	return 1, nil
}

func (e *ElasticUserService) Patch(ctx context.Context, id string, user map[string]interface{}) (int64, error) {

	userJsonString := convertDocToJson(user)
	var userid = reflect.ValueOf(user["_id"])
	delete(user, "_id")
	request := esapi.UpdateRequest{
		Index: "users",
		DocumentID: userid.String(),
		Body: strings.NewReader(userJsonString),
		Refresh: "true",
	}
	response, err := request.Do(ctx, e.elastic)

	if err != nil {
		panic(err)
		return 0, nil
	}

	defer response.Body.Close()

	var result map[string]interface{}

	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		panic(err)
		return 0, nil
	}

	fmt.Println("IndexRequest to update Status: ", response.Status())
	fmt.Println("Result: ", result["result"])

	fmt.Print("The user %v has been updated successfully.", userid.String())
	return 1, nil
}

func (e *ElasticUserService) Delete(ctx context.Context, id string) (int64, error) {
	request := esapi.DeleteRequest{
		Index: "users",
		DocumentID: id,
	}
	response, err := request.Do(ctx, e.elastic)

	if err != nil {
		panic(err)
		return 0, nil
	}

	defer response.Body.Close()

	var result map[string]interface{}

	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		panic(err)
		return 0, nil
	}

	fmt.Println("IndexRequest to update Status: ", response.Status())
	fmt.Println("Result: ", result["result"])

	fmt.Print("Delete user: %s successfully.", id)
	return 1, nil
}
