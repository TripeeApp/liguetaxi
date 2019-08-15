package liguetaxi

import "context"

// ReqStatus is the request status.
// Success = 1
// Fail = 0
type ReqStatus int

// UserStatus is the user status.
// Active - 1
// Inactive - 0
type UserStatus int

func (us *UserStatus) UnmarshalJSON(b []byte) error {
	*us = UserStatusActive
	return nil
}

// Request Status
const (
	ReqStatusFail ReqStatus = iota
	ReqStatusOK
)

// User status
const (
	UserStatusInactive UserStatus = iota
	UserStatusActive
)

var (
	// Endpoint for reading user info.
	getUserEndpoint endpoint = `user/check_authorized`

	// Endpoint for editing user status.
	editUserStatusEndpoint endpoint = `user/status_authorized`

	// Endpoint for editing user info.
	editUserEndpoint endpoint = `user/edit_authorized`

	// Endpoint for creating user.
	createUserEndpoint endpoint = `user/create_authorized`

	// Endpoint for reading classifier field.
	getUserFieldEndpoint endpoint = `user/check_authorized_field`

	// Endpoint for creating classifier field.
	createUserFieldEndpoint endpoint = `user/create_authorized_field`
)

// status is the request status.
type status struct {
	Status ReqStatus `json:"status"`
}

// OperationResponse is the response returned by the API
// for non-idempotent operations.
type OperationResponse struct {
	status

	Message string `json:"message"`
}

type ClassifierOperationResponse struct {
	OperationResponse

	Data string `json:"data"`
}

// UserResponse is the response returned by the API
// when listing a user info.
type UserResponse struct {
	status

	Data struct {
		ID	string		`json:"authorized_id"`
		Name	string		`json:"client_name"`
		Email	string		`json:"client_email"`
		Phone	string		`json:"client_phone"`
		Status	UserStatus	`json:"cod_status"`
	} `json:"data"`
}

// User is sent to server when creating or editing user.
type User struct {
	ID		string `json:"unique_field,omitempty"`
	Name		string `json:"user_name"`
	Email		string `json:"user_email"`
	Phone		string `json:"user_phone,omitempty"`
	Password	string `json:"user_password,omitempty"`
	Classifier1	string `json:"classificador1,omitempty"`
	Classifier2	string `json:"classificador2,omitempty"`
	Classifier3	string `json:"classificador3,omitempty"`
	Classifier4	string `json:"classificador4,omitempty"`
	Classifier5	string `json:"classificador5,omitempty"`
	Classifier6	string `json:"classificador6,omitempty"`
	Classifier7	string `json:"classificador7,omitempty"`
	Classifier8	string `json:"classificador8,omitempty"`
	Classifier9	string `json:"classificador9,omitempty"`
	Classifier10	string `json:"classificador10,omitempty"`
	Classifier11	string `json:"classificador11,omitempty"`
	Classifier12	string `json:"classificador12,omitempty"`
	Classifier13	string `json:"classificador13,omitempty"`
	Classifier14	string `json:"classificador14,omitempty"`
	Classifier15	string `json:"classificador15,omitempty"`
	Classifier16	string `json:"classificador16,omitempty"`
	Classifier17	string `json:"classificador17,omitempty"`
	Classifier18	string `json:"classificador18,omitempty"`
	Classifier19	string `json:"classificador19,omitempty"`
	Classifier20	string `json:"classificador20,omitempty"`
}

type Classifier struct {
	ID		string `json:"field_id"`
	Value		string `json:"field_value"`
	AdditionalValue string `json:"field_additional_value"`
}

type ClassifierResponse struct {
	status

	Data []Classifier `json:"data"`
}

// UserService handles the requests related to the user.
type UserService struct {
	client requester
}

func (us *UserService) Read(ctx context.Context, name, id string) (*UserResponse, error) {
	return nil, nil
}

func (us *UserService) Create(ctx context.Context, u *User) (*OperationResponse, error) {
	return nil, nil
}

func (us *UserService) Update(ctx context.Context, u *User) (*OperationResponse, error) {
	return nil, nil
}

func (us *UserService) ReadClassifier(ctx context.Context, field string, value string) (*ClassifierResponse, error) {
	return nil, nil
}

func (us *UserService) CreateClassifier(ctx context.Context) (*ClassifierOperationResponse, error) {
	return nil, nil
}
