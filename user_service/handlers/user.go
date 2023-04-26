package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/princeparmar/contact_manager/repositories"
	"github.com/princeparmar/go-helpers/clienthelper"
	"github.com/princeparmar/go-helpers/context"
	"github.com/princeparmar/go-helpers/utils"
)

// User defines a struct for user data.
type User struct {
	ID     int
	Name   string `json:"name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

// createUserModel maps User to User model.
func createUserModel(u *User) *repositories.User {
	return &repositories.User{
		ID:       u.ID,
		UserName: u.Name,
		Mobile:   u.Mobile,
		EmailID:  u.Email,
	}
}

// ParseRequest parses the HTTP request and extracts any relevant data into the User object.
func (u *User) ParseRequest(ctx context.IContext, w http.ResponseWriter, r *http.Request) error {

	if r.Method == http.MethodPost {
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
	}

	// Parse ID from the query parameter
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
	User
	clienthelper.BaseAPIExecutor
	UserRepo repositories.UserRepository
}

// NewCreateUserExecutor returns a new instance of CreateUserExecutor.
func NewCreateUserExecutor(repo repositories.UserRepository) clienthelper.APIExecutor {
	return &CreateUserExecutor{
		UserRepo: repo,
	}
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

// DeleteUserExecutor defines an APIExecutor for deleting a user by ID.
type DeleteUserExecutor struct {
	User
	clienthelper.BaseAPIExecutor
	UserRepo repositories.UserRepository
}

// NewDeleteUserExecutor returns a new instance of DeleteUserExecutor.
func NewDeleteUserExecutor(repo repositories.UserRepository) clienthelper.APIExecutor {
	return &DeleteUserExecutor{
		UserRepo: repo,
	}
}

// Controller executes the business logic for deleting a user by ID and returns any errors that occur during execution.
func (e *DeleteUserExecutor) Controller(ctx context.IContext) (interface{}, error) {
	id := e.User.ID
	err := e.UserRepo.Delete(id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateUserExecutor defines an APIExecutor for updating a user by ID.
type UpdateUserExecutor struct {
	User
	clienthelper.BaseAPIExecutor
	UserRepo repositories.UserRepository
}

// NewUpdateUserExecutor returns a new instance of UpdateUserExecutor.
func NewUpdateUserExecutor(repo repositories.UserRepository) clienthelper.APIExecutor {
	return &UpdateUserExecutor{
		UserRepo: repo,
	}
}

// Controller executes the business logic for updating a user by ID and returns the updated user
// and any errors that occur during execution.
func (e *UpdateUserExecutor) Controller(ctx context.IContext) (interface{}, error) {
	user := createUserModel(&e.User)
	err := e.UserRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserExecutor defines an APIExecutor for getting a user by ID.
type GetUserExecutor struct {
	User
	clienthelper.BaseAPIExecutor
	UserRepo repositories.UserRepository
}

// NewGetUserExecutor returns a new instance of GetUserExecutor.
func NewGetUserExecutor(repo repositories.UserRepository) clienthelper.APIExecutor {
	return &GetUserExecutor{
		UserRepo: repo,
	}
}

// Controller executes the business logic for getting a user by ID and returns the user
// and any errors that occur during execution.
func (e *GetUserExecutor) Controller(ctx context.IContext) (interface{}, error) {
	id := e.User.ID
	user, err := e.UserRepo.Get(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsersExecutor defines an APIExecutor for getting all users.
type GetAllUsersExecutor struct {
	clienthelper.BaseAPIExecutor
	UserRepo repositories.UserRepository
}

// NewGetAllUsersExecutor returns a new instance of GetAllUsersExecutor.
func NewGetAllUsersExecutor(repo repositories.UserRepository) clienthelper.APIExecutor {
	return &GetAllUsersExecutor{
		UserRepo: repo,
	}
}

// Controller executes the business logic for getting all users and returns the users
// and any errors that occur during execution.
func (e *GetAllUsersExecutor) Controller(ctx context.IContext) (interface{}, error) {
	return e.UserRepo.GetAll()
}

// UserAccessExecutor defines an APIExecutor for getting a list of accesses based on the user ID.
type UserAccessExecutor struct {
	User
	clienthelper.BaseAPIExecutor
	userRoleRepository repositories.UserRoleRepository
}

// NewUserAccessExecutor returns a new instance of UserAccessExecutor.
func NewUserAccessExecutor(repo repositories.UserRoleRepository) clienthelper.APIExecutor {
	return &UserAccessExecutor{
		userRoleRepository: repo,
	}
}

// Controller executes the business logic for getting a list of accesses based on the user ID and returns the accesses
// and any errors that occur during execution.
func (e *UserAccessExecutor) Controller(ctx context.IContext) (interface{}, error) {
	accesses, err := e.userRoleRepository.GetAllAccess(e.User.ID)
	if err != nil {
		return nil, err
	}
	if len(accesses) == 0 {
		return nil, errors.New("no accesses found for user")
	}

	return accesses, nil
}
