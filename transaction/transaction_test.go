package transaction_test

import (
	"errors"
	"testing"

	"github.com/tomek7667/gotransaction/transaction"
)

var (
	ErrSampleActionFailure = errors.New("simulated action error")
	ErrRollbackFailed      = errors.New("transaction rollback failed")
)

func TestSuccessfulRollbackAfterServeralChangesTransac(t *testing.T) {
	statefulResourceA := 0
	statefulResourceB := 0
	statefulResourceC := 0
	errResourceCFailed := errors.New("stateful resource C have failed")

	tx := transaction.New()

	err := tx.A(func() error {
		statefulResourceA++
		return nil
	}, func() error {
		statefulResourceA--
		return nil
	})
	if err != nil {
		t.Errorf("there shouldn't be an error at this point")
	}
	err = tx.A(func() error {
		statefulResourceA++
		return nil
	}, func() error {
		statefulResourceA--
		return nil
	})
	if err != nil {
		t.Errorf("there shouldn't be an error at this point")
	}
	err = tx.A(func() error {
		statefulResourceA++
		return nil
	}, func() error {
		statefulResourceA--
		return nil
	})
	if err != nil {
		t.Errorf("there shouldn't be an error at this point")
	}
	err = tx.A(func() error {
		statefulResourceB++
		return nil
	}, func() error {
		statefulResourceB--
		return nil
	})
	if err != nil {
		t.Errorf("there shouldn't be an error at this point")
	}
	err = tx.A(func() error {
		return errResourceCFailed
	}, func() error {
		statefulResourceC--
		return nil
	})
	if err != nil {
		reason, rollbackErrors := tx.R(err)
		if reason != err {
			t.Errorf("reason should be the same as action error")
		}
		if rollbackErrors != nil {
			t.Errorf("the rollback should have been succesful, as stated above, the success value for the rollback is \"abc\"")
		}
	} else {
		t.Errorf("an error should have appeared here")
	}
	if statefulResourceA != 0 {
		t.Errorf("resource A has not been rolled back")
	}
	if statefulResourceB != 0 {
		t.Errorf("resource B has not been rolled back")
	}
	if statefulResourceC != 0 {
		t.Errorf("resource C is not in the initial state")
	}
}
