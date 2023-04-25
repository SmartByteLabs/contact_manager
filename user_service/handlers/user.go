package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/princeparmar/contact_manager/repositories"
	"github.com/princeparmar/go-helpers/clienthelper"
	"github.com/princeparmar/go-helpers/context"
	"github.com/princeparmar/go-helpers/utils"
)

func createUserModel(u *User) *repositories.User {
	return &repositories.User{
		ID:       u.ID,
		UserName: u.Name,
		Mobile:   u.Mobile,
		EmailID:  u.Email,
	}
}

// User defines a struct for user data.
type User struct {
	ID     int
	Name   string `json:"name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

// ParseRequest parses the HTTP request and extracts any relevant data into the User object.
func (u *User) ParseRequest(ctx context.IContext, w http.ResponseWriter, r *http.Request) error {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Unmarshal the request body into the User object
	err = json.Unmarshal(body, u)
	if err != nil {
		return err
	}

	id := r.URL.Query().Get("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid id in query")
	}

	u.ID = i

	return nil

}

// ValidateRequest validates the data in the User object and returns any errors that occur during validation.
func (u *User) ValidateRequest(ctx context.IContext) error {

	// Validate email format
	if u.Email != "" && !utils.ValidateEmail(u.Email) {
		return errors.New("email format is invalid")
	}

	// Validate mobile format
	if u.Mobile != "" && !utils.ValidateMobile(u.Mobile) {
		return errors.New("mobile format is invalid")
	}

	return nil
}

// CreateUserExecutor defines an APIExecutor for creating a new user.
type CreateUserExecutor struct {
	User     User
	UserRepo repositories.UserRepository
}

// NewCreateUserExecutor returns a new instance of CreateUserExecutor.
func NewCreateUserExecutor(repo repositories.UserRepository) clienthelper.APIExecutor {
	return &CreateUserExecutor{
		UserRepo: repo,
	}
}

// ParseRequest parses the HTTP request and extracts any relevant data into the User object.
func (e *CreateUserExecutor) ParseRequest(ctx context.IContext, w http.ResponseWriter, r *http.Request) error {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Unmarshal the request body into the User object
	err = json.Unmarshal(body, &e.User)
	if err != nil {
		return err
	}

	return nil
}

// ValidateRequest validates the data in the User object and returns any errors that occur during validation.
func (e *CreateUserExecutor) ValidateRequest(ctx context.IContext) error {
	// Validate email format
	if e.User.Email != "" && !utils.ValidateEmail(e.User.Email) {
		return errors.New("email format is invalid")
	}

	// Validate mobile format
	if e.User.Mobile != "" && !utils.ValidateMobile(e.User.Mobile) {
		return errors.New("mobile format is invalid")
	}

	return nil
}

// Controller executes the business logic for creating a new user and returns the created user
// and any errors that occur during execution.
func (e *CreateUserExecutor) Controller(ctx context.IContext) (interface{}, error) {
	user := createUserModel(&e.User)
	err := e.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// HandleResponse takes in the result of the Controller method, any errors that occur,
// and the APIResponse object to generate the HTTP response. It takes in an interface,
// an error, and an APIResponse pointer and does not return anything.
func (e *CreateUserExecutor) HandleResponse(ctx context.IContext, response interface{}, err error, apiResponse *clienthelper.APIResponse) {
	if err != nil {
		// Handle error response
		apiResponse.SetStatusCode(http.StatusInternalServerError).
			GetResponse().SetStatusCode(http.StatusInternalServerError).
			GetErrors().SetClientMessage("an error occurred while processing the request").AddError(clienthelper.NewErrorFromErr(err))
	} else {
		// Handle success response
		apiResponse.SetStatusCode(http.StatusOK).
			GetResponse().SetStatusCode(http.StatusOK).SetData(response)
	}

	// Send the API response
	if err := apiResponse.Send(); err != nil {
		// Handle send response error
		log.Printf("Failed to send API response: %v", err)
	}
}
