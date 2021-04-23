package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	//"fmt"
	"github.com/olivere/elastic/v7"
	//"strings"

	_ "github.com/olivere/elastic/v7"
	. "go-service/internal/models"
)

type ElasticUserService struct {
	elastic *elastic.Client
}

func NewEUserService(ela *elastic.Client) *ElasticUserService {
	return &ElasticUserService{elastic: ela}
}

func (e *ElasticUserService) GetAll(ctx context.Context) (*[]User, error) {
	var result []User

	query := elastic.MatchAllQuery{}

	searchResult, err := e.elastic.Search().Index("users").Query(query).Do(ctx)
	if err != nil {
		fmt.Print("Error during execution GetAll : %s", err.Error())
	}

	for _, hit := range searchResult.Hits.Hits {
		var user User
		err := json.Unmarshal(hit.Source, &user)

		if err != nil {
			fmt.Print("An error when read the user: ", err)
		}
		result = append(result, user)
	}
	return &result, nil
}


func (e *ElasticUserService) Load(ctx context.Context, id string) (*User, error) {
	var user User
	fmt.Println("Try to find id := " + id)

	query := elastic.NewBoolQuery()
	searchConditions := []elastic.Query{elastic.NewTermQuery("id", id)}
	query = query.Must(searchConditions...)

	searchResult, err := e.elastic.Search().Index("users").Query(query).Do(ctx)

	if err != nil {
		fmt.Print("Error during the Finding user: %s", err.Error())
	}

	fmt.Print("Found a total of %d user", searchResult.Hits.TotalHits.Value)
	for _, hit := range searchResult.Hits.Hits {
		err := json.Unmarshal(hit.Source, & user)
		if err != nil {
			fmt.Print("Error when reading user: %s", err.Error())
		}
		return &user, nil
	}

	return &user, nil
}

func (e *ElasticUserService) Insert(ctx context.Context, user *User) (int64, error) {
	if user == nil {
		fmt.Print("Can not add null user")
		return 0, nil
	}

	_, err := e.elastic.Index().Index("users").BodyJson(&user).Do(ctx)

	if err != nil {
		panic(err)
		return 0, nil
	}

	fmt.Print("The new user %s has been added successfully.", user.Username)
	return 1, nil
}

func (e *ElasticUserService) Update(ctx context.Context, user *User) (int64, error) {
	fmt.Println("Update user Id = ", user.Id)
	query := elastic.NewMatchQuery("id", user.Id)

	// This function must use Id of the Entry as param to search for update
	//update, err := e.elastic.Update().Index("users").Id(user.Id).Doc(map[string]interface{}{"Id": user.Id, "username": user.Username, "email": user.Email, "phone": user.Phone, "dateOfBirth": user.DateOfBirth}).Do(ctx)//.Doc(map[string]interface{}{"id": user.Id, "username": user.Username, "email": user.Email, "phone": user.Phone, "dateOfBirth": user.DateOfBirth}).Do(ctx)
	update, err := e.elastic.UpdateByQuery().Index("users").Query(query).Script(elastic.NewScriptInline("ctx._source.username = '" + user.Username + "';ctx._source.email = '" + user.Email + "';ctx._source.phone = '" + user.Phone + "';ctx._source.dateOfBirth = '" + user.DateOfBirth.Format(time.RFC3339) +"'")).Do(ctx)
	fmt.Println("Update date time = ", user.DateOfBirth.Format("yyyy-MM-dd'T'HH:mm:ss.SSSZ"))
	if err != nil {
		fmt.Print("Error when updating user: %s", err.Error())
		return 0, nil
	}

	fmt.Print("The user ID = %s has been updated successfully", update.Updated)
	return 1, nil
}

func (e *ElasticUserService) Delete(ctx context.Context, id string) (int64, error) {
	query := elastic.NewBoolQuery()
	searchConditions := []elastic.Query{elastic.NewTermQuery("id", id)}
	query = query.Must(searchConditions...)

	_, err := elastic.NewDeleteByQueryService(e.elastic).Index("users").Query(query).Do(ctx)

	if err != nil {
		fmt.Print("Error when delete user: %s", err.Error())
		return 0, nil
	}

	e.elastic.Flush().Index("user").Do(ctx)
	fmt.Print("Delete user: %s successfully", id)
	return 1, nil
}
