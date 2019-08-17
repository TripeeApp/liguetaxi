package integration

import (
	"context"
	"fmt"
	"testing"

	"bitbucket.org/mobilitee/liguetaxi"
)

var (
	userID		= fmt.Sprintf("00%s", randString(7, numberBytes))
	userName	= randString(10, letterBytes)
)

// Setup general user data
func TestMain(t *testing.T) {
	// Register functional.
	newUserID := &liguetaxi.Classifier{
		Field: "2",
		Value: userID,
	}

	op, err := ligtaxi.User.CreateClassifier(context.Background(), newUserID)
	if err != nil {
		t.Fatalf("got error calling User.CreateClassifier(%+v): %s; want nil.", newUserID, err.Error())
	}

	if op.Status != liguetaxi.ReqStatusOK {
		t.Errorf("got failed request. Status: %d, message: %s; want %d.", op.Status, op.Message, liguetaxi.ReqStatusOK)
	}

	field, err := ligtaxi.User.ReadClassifier(context.Background(), "2", userID)
	if err != nil {
		t.Fatalf("got error calling User.ReadClassifier(context.Background(), 2, %s): %s; want nil.", userID, err.Error())
	}

	if field.Status != liguetaxi.ReqStatusOK {
		t.Errorf("got failed request. Status: %d; want %d.", field.Status, liguetaxi.ReqStatusOK)
	}

	costCenter := randString(10, letterBytes)

	// Register cost center.
	newCostCenter := &liguetaxi.Classifier{
		Field: "1",
		Value: costCenter,
	}

	op, err = ligtaxi.User.CreateClassifier(context.Background(), newCostCenter)
	if err != nil {
		t.Fatalf("got error calling User.CreateClassifier(%+v): %s; want nil.", newCostCenter, err.Error())
	}

	if op.Status != liguetaxi.ReqStatusOK {
		t.Errorf("got failed request. Status: %d, message: %s; want %d.", op.Status, op.Message, liguetaxi.ReqStatusOK)
	}

	field, err = ligtaxi.User.ReadClassifier(context.Background(), "1", costCenter)
	if err != nil {
		t.Fatalf("got error calling User.ReadClassifier(context.Background(), 2, %s): %s; want nil.", costCenter, err.Error())
	}

	if field.Status != liguetaxi.ReqStatusOK {
		t.Errorf("got failed request. Status: %d; want %d.", field.Status, liguetaxi.ReqStatusOK)
	}

	newUser := &liguetaxi.User{
		Name: userName,
		Email: fmt.Sprintf("%s@gmail.com", randString(5, letterBytes)),
		Phone: "11986548744",
		Password: "test1234",
		Classifier1: costCenter,
		Classifier2: userID,
		Classifier3: "0003",
		Classifier4: "Testing 4",
	}

	uop, err := ligtaxi.User.Create(context.Background(), newUser)
	if err != nil {
		t.Fatalf("got error calling User.Create(%+v): %s; want nil.", newUser, err.Error())
	}

	if uop.Status != liguetaxi.ReqStatusOK {
		t.Errorf("got request. Status: %d, message: '%s'; want %d.", uop.Status, uop.Message, liguetaxi.ReqStatusOK)
	}

	user, err := ligtaxi.User.Read(context.Background(), userID, "")
	if err != nil {
		t.Errorf("got error calling User.Read(%s, %s): %s; want nil.", userID, userName, err.Error())
	}

	if user.Status != liguetaxi.ReqStatusOK {
			t.Errorf("got failed request. Status: %d; want %d.", user.Status, liguetaxi.ReqStatusOK)
	}
}
