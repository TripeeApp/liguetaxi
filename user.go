package liguetaxi

import (
	"context"
	"encoding/json"
	"net/http"
)

// User statuses
const (
	UserStatusInactive userStatus = iota
	UserStatusActive
	UserStatusSynching
)

var (
	// Endpoint for reading user info.
	readUserEndpoint endpoint = `user/check_authorized`

	// Endpoint for editing user status.
	updateUserStatusEndpoint endpoint = `user/status_authorized`

	// Endpoint for editing user info.
	updateUserEndpoint endpoint = `user/edit_authorized`

	// Endpoint for creating user.
	createUserEndpoint endpoint = `user/create_authorized`

	// Endpoint for reading classifier field.
	readClassifierEndpoint endpoint = `user/check_authorized_field`

	// Endpoint for creating classifier field.
	createClassifierEndpoint endpoint = `user/create_authorized_field`
)

// userStatus is the user status.
// Active - 1
// Inactive - 0
type userStatus int

// UnmarshalText implements the TextUnmarshaler interface for
// userStatus type
func (us *userStatus) UnmarshalJSON(t []byte) error {
	switch string(t) {
	case `"24"`:
		*us = UserStatusActive
	case `"46"`:
		*us = UserStatusSynching
	default:
		*us = UserStatusInactive
	}
	return nil
}

// MarshalJSON implements the Marshaler interface for
// userStatus type
func (us *userStatus) MarshalJSON() ([]byte, error) {
	if us == nil {
		return []byte(`null`), nil
	}

	switch *us {
	case UserStatusActive:
		return []byte(`"24"`), nil
	case UserStatusSynching:
		return []byte(`"46"`), nil
	default:
		return []byte(`"25"`), nil
	}
	return nil, nil
}

// New return a pointer to userStatus.
func (us userStatus) New() *userStatus {
	return &us
}

// emptObjToStr is the field that should return an string
// but returns an empty object when string is empty.
type emptyObjToStr string

// UnmarshalJSON implements the Unmarshaler interface for
// emptyObjToStr type
func (e *emptyObjToStr) UnmarshalJSON(b []byte) error {
	var s string
	if token := string(b); token != `{}` {
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*e = emptyObjToStr(s)
	}
	return nil
}

// String returns underlying string for
// emptyObjToStr type.
func (e emptyObjToStr) String() string {
	return string(e)
}

// OperationResponse is the response returned by the API
// for non-idempotent operations on user.
type OperationResponse struct {
	Status reqStatus

	Message string `json:"message"`
}

// ClassifierOperationResponse is the response returned by the API
// for non-idempotent operations on user's classifier fields.
type ClassifierOperationResponse struct {
	OperationResponse

	Data string `json:"data"`
}

// DataUser is the result from check user request.
type DataUser struct {
	ID                string         `json:"authorized_id"`
	Name              string         `json:"client_name"`
	Email             *emptyObjToStr `json:"client_email"`
	Phone             *emptyObjToStr `json:"client_phone"`
	Status            *userStatus    `json:"cod_status"`
	StatusDescription string         `json:"status_description"`
}

// UserResponse is the response returned by the API
// when listing a user info.
type UserResponse struct {
	Status reqStatus

	Data DataUser `json:"data"`
}

// Pulled off for testing
type userFilter struct {
	ID   string `json:"unique_field"`
	Name string `json:"user_name,omitempty"`
}

// User is sent to server when creating or editing user.
type User struct {
	ID           string `json:"unique_field,omitempty"`
	Name         string `json:"user_name"`
	Email        string `json:"user_email"`
	Phone        string `json:"user_phone,omitempty"`
	Password     string `json:"user_password,omitempty"`
	Classifier1  string `json:"classificador1,omitempty"`
	Classifier2  string `json:"classificador2,omitempty"`
	Classifier3  string `json:"classificador3,omitempty"`
	Classifier4  string `json:"classificador4,omitempty"`
	Classifier5  string `json:"classificador5,omitempty"`
	Classifier6  string `json:"classificador6,omitempty"`
	Classifier7  string `json:"classificador7,omitempty"`
	Classifier8  string `json:"classificador8,omitempty"`
	Classifier9  string `json:"classificador9,omitempty"`
	Classifier10 string `json:"classificador10,omitempty"`
	Classifier11 string `json:"classificador11,omitempty"`
	Classifier12 string `json:"classificador12,omitempty"`
	Classifier13 string `json:"classificador13,omitempty"`
	Classifier14 string `json:"classificador14,omitempty"`
	Classifier15 string `json:"classificador15,omitempty"`
	Classifier16 string `json:"classificador16,omitempty"`
	Classifier17 string `json:"classificador17,omitempty"`
	Classifier18 string `json:"classificador18,omitempty"`
	Classifier19 string `json:"classificador19,omitempty"`
	Classifier20 string `json:"classificador20,omitempty"`
}

// UserStatus is the user status infos.
type UserStatus struct {
	ID     string     `json:"authorized_id"`
	Name   string     `json:"user_name,omitempty"`
	Status userStatus `json:"status"`
	Reason string     `json:"reason,omitempty"`
}

type classifierFilter struct {
	Field string `json:"field"`
	Value string `json:"field_value"`
}

// Classifier is the classifier field infos.
type Classifier struct {
	ID              string `json:"field_id,omitempty"`
	Field           string `json:"field,omitempty"`
	Value           string `json:"field_value"`
	AdditionalValue string `json:"field_additional_value,omitempty"`
}

// ClassifierResponse is the response returned by the API
// when reading the classifier field info.
type ClassifierResponse struct {
	Status reqStatus

	Data []Classifier `json:"data"`
}

// UserService handles the requests related to the user.
type UserService service

// Read returns User infos or an error.
func (us *UserService) Read(ctx context.Context, id, name string) (*UserResponse, error) {
	u := &UserResponse{}

	if err := us.client.Request(ctx, http.MethodPost, readUserEndpoint, userFilter{id, name}, u); err != nil {
		return nil, err
	}

	return u, nil
}

// Create returns the status operation for creating a user or an error.
func (us *UserService) Create(ctx context.Context, u *User) (*OperationResponse, error) {
	op := &OperationResponse{}

	if err := us.client.Request(ctx, http.MethodPost, createUserEndpoint, u, op); err != nil {
		return op, err
	}

	return op, nil
}

// Update returns the status operation for updating user or an error.
func (us *UserService) Update(ctx context.Context, u *User) (*OperationResponse, error) {
	op := &OperationResponse{}

	if err := us.client.Request(ctx, http.MethodPost, updateUserEndpoint, u, op); err != nil {
		return op, err
	}

	return op, nil
}

// UpdateStatus returns the status operation for updating the user status or an error.
func (us *UserService) UpdateStatus(ctx context.Context, s *UserStatus) (*OperationResponse, error) {
	op := &OperationResponse{}

	us.client.Request(ctx, http.MethodPost, updateUserStatusEndpoint, s, op)

	return op, nil
}

// ReadClassifier returns the classifier field info.
func (us *UserService) ReadClassifier(ctx context.Context, field, value string) (*ClassifierResponse, error) {
	c := &ClassifierResponse{}

	if err := us.client.Request(ctx, http.MethodPost, readClassifierEndpoint, classifierFilter{field, value}, c); err != nil {
		return c, err
	}

	return c, nil
}

// CreateClassifier returns the status operation for creating classifier field or an error.
func (us *UserService) CreateClassifier(ctx context.Context, c *Classifier) (*ClassifierOperationResponse, error) {
	co := &ClassifierOperationResponse{}

	if err := us.client.Request(ctx, http.MethodPost, createClassifierEndpoint, c, co); err != nil {
		return co, err
	}

	return co, nil
}
