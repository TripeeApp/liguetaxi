package integration

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"bitbucket.org/mobilitee/liguetaxi"
)

var (
	userID		= flag.String("id", fmt.Sprintf("00%s", randString(9, numberBytes)), "Define user ID to be created or searched")
	retries		= 8
	delay		= 5 * time.Second

)

// Setup general user data
func TestMain(t *testing.T) {
	// Register functional.
	newUserID := &liguetaxi.Classifier{
		Field: "2",
		Value: *userID,
	}

	op, err := ligtaxi.User.CreateClassifier(context.Background(), newUserID)
	if err != nil {
		t.Fatalf("got error calling User.CreateClassifier(%+v): %s; want nil.", newUserID, err.Error())
	}

	if op.Status != liguetaxi.ReqStatusOK {
		t.Errorf("got failed request. Status: %d, message: %s; want %d.", op.Status, op.Message, liguetaxi.ReqStatusOK)
	}

	ok, err := checkOperation(delay, retries, func() error {
		field, err := ligtaxi.User.ReadClassifier(context.Background(), "2", *userID)
		if err != nil {
			return err
		}

		if field.Status != liguetaxi.ReqStatusOK {
			return fmt.Errorf("got failed request. Status: %d; want %d.", field.Status, liguetaxi.ReqStatusOK)
		}

		return nil
	})
	if !ok {
		t.Fatalf("got error calling User.ReadClassifier(context.Background(), 2, %s): %s; want nil.", *userID, err.Error())
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

	ok, err = checkOperation(delay, retries, func() error {
		field, err := ligtaxi.User.ReadClassifier(context.Background(), "1", costCenter)
		if err != nil {
			return err
		}

		if field.Status != liguetaxi.ReqStatusOK {
			return fmt.Errorf("got failed request. Status: %d; want %d.", field.Status, liguetaxi.ReqStatusOK)
		}

		return nil
	})
	if !ok {
		t.Fatalf("got error calling User.ReadClassifier(context.Background(), 1, %s): %s; want nil.", costCenter, err.Error())
	}

	newUser := &liguetaxi.User{
		Name: randString(10, letterBytes),
		Email: fmt.Sprintf("%s@gmail.com", randString(5, letterBytes)),
		Phone: "11986548744",
		Password: "test1234",
		Classifier1: costCenter,
		Classifier2: *userID,
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

	ok, err = checkOperation(delay, retries, func() error {
		user, err := ligtaxi.User.Read(context.Background(), *userID, "")
		if err != nil {
			return err
		}

		if user.Status != liguetaxi.ReqStatusOK {
			return fmt.Errorf("got failed request. Status: %d; want %d.", user.Status, liguetaxi.ReqStatusOK)
		}

		return nil
	})

	if !ok {
		t.Fatalf("got error calling User.Read(%s, ''): %s; want nil.", *userID, err.Error())
	}
}

func checkUpdateStatus(u *liguetaxi.UserResponse) func() error {
	// Check if user was indeed updated.
	return func() error {
		user, err := ligtaxi.User.Read(context.Background(), *userID, "")
		if err != nil {
			return fmt.Errorf("got error calling User.Read(%s): %s; want nil.", *userID, err.Error())
		}

		if want := liguetaxi.ReqStatusOK; user.Status != want {
			fmt.Errorf("got failed request. Status: %d; want %d.", user.Status, want)
		}

		if *user.Data.Status != liguetaxi.UserStatusSynching {
			*u = *user
			return nil
		}

		return fmt.Errorf("Last status: %d.", *user.Data.Status)
	}
}

func TestUserUpdateStatus(t *testing.T) {
	u, err := ligtaxi.User.Read(context.Background(), *userID, "")
	if err != nil {
		t.Errorf("got error calling User.Read(%s): %s; want nil.", *userID, err.Error())
	}

	if want := liguetaxi.ReqStatusOK; u.Status != want {
		t.Errorf("got failed request. Status: %d; want %d.", u.Status, want)
	}

	newUserStatus := &liguetaxi.UserStatus{
		ID: u.Data.ID,
		Status: liguetaxi.UserStatusInactive,
	}

	op, err := ligtaxi.User.UpdateStatus(context.Background(), newUserStatus)
	if err != nil {
		t.Fatalf("got error calling User.UpdateStatus(%+v): %s; want nil.", newUserStatus, err.Error())
	}

	if want := liguetaxi.ReqStatusOK; op.Status != want {
		t.Fatalf("got failed request. Status: %d; want %d.", op.Status, want)
	}

	user := &liguetaxi.UserResponse{}
	if ok, err := checkOperation(delay, retries, checkUpdateStatus(user)); !ok {
		t.Fatalf("got error calling User.Read(%s): %s; want nil.", *userID, err.Error())
	}

	if want := liguetaxi.UserStatusInactive; user.Data.ID != ""  && *user.Data.Status != want {
		t.Errorf("got user.Status: %d; want %d.", *user.Data.Status, want)
	}
}

func TestUserUpdate(t *testing.T) {
	newEmail := fmt.Sprintf("%s@gmail.com", randString(5, letterBytes))
	newUserInfo := &liguetaxi.User{
		ID: *userID,
		Email: newEmail,
	}

	op, err := ligtaxi.User.Update(context.Background(), newUserInfo)
	if err != nil {
		t.Fatalf("got error calling User.Update(%+v): %s; want nil.", newUserInfo, err.Error())
	}

	if want := liguetaxi.ReqStatusOK; op.Status != want {
		t.Errorf("got status: %d; want %d.", op.Status, want)
	}

	user := &liguetaxi.UserResponse{}

	if ok, err := checkOperation(delay, retries, checkUpdateStatus(user)); !ok {
		t.Fatalf("got error calling User.Read(%s, ''): %s; want nil.", *userID, err.Error())
	}

	if user.Data.Email.String() != newEmail {
		t.Errorf("got email: %s; want: %s.", user.Data.Email, newEmail)
	}
}
