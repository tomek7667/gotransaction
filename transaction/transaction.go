package transaction

import (
	"errors"
	"slices"
)

// `Transaction` is a struct that allows for smart rollbacking instead of copy-pasting the err handler logic for rollbacking.
type Transaction struct {
	// `Rollbacks` are the funcs triggered when `Transaction` has returned an error at any point. Of coure the rollback action might fail.
	Rollbacks []func() error
}

func New() *Transaction {
	return &Transaction{
		Rollbacks: []func() error{},
	}
}

// `A` is the committed action that will be just called and in case of no error, will add a rollbacking action to the transaction
func (t *Transaction) A(
	action func() error,
	rollback func() error,
) error {
	err := action()
	if err == nil {
		t.Rollbacks = append(t.Rollbacks, rollback)
	}
	return err
}

// `R` is the rollbacking action, first output is the main reason the rollback even happened and the second one is the rollbacking errors if any occurred
func (t *Transaction) R(reason error) (error, error) {
	// We need to reverse the order of rollbacks as the actions should be reversed in the opposite order, so the dependant objects are rolled back first.
	slices.Reverse(t.Rollbacks)
	var rollbackErrs []error
	for _, rollback := range t.Rollbacks {
		err := rollback()
		if err != nil {
			rollbackErrs = append(rollbackErrs, err)
		}
	}
	return reason, errors.Join(rollbackErrs...)
}
