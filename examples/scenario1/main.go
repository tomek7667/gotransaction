package main

import (
	"fmt"
	"log/slog"

	"github.com/tomek7667/gotransaction/examples/scenario1/example"
	"github.com/tomek7667/gotransaction/transaction"
)

func main() {
	var (
		err            error
		stripeClient   example.StripeClient
		databaseClient example.DatabaseClient
		filesClient    example.FilesClient
	)
	stripeClient = example.StripeClient{}
	databaseClient = example.DatabaseClient{}
	filesClient = example.FilesClient{}
	filesClient.InitFs()
	// as you can see we have different services, all of which store state, and all of which might fail at some point.
	// in this scenario:
	//
	// Stripe-connected app with storing files.
	// 1. adds a row in database
	// 2. uploads a file
	// 3. registers in stripe

	// if the 3rd step fails, rollbacks for 2. and 1. will be issued.

	tx := transaction.New()

	// 1.
	var recordIdx int
	err = tx.A(func() error {
		recordIdx, err = databaseClient.Insert("customer")
		return err
	}, func() error {
		return databaseClient.Delete(recordIdx)
	})
	if err != nil {
		originalErr, rollbackErr := tx.R(err)
		fmt.Printf("something went wrong with DB. Rollbacked. Err: %v", originalErr)
		if rollbackErr != nil {
			// the action that should revert changes failed. That's deserving a panic
			panic(rollbackErr)
		}
		return
	}

	// 2.
	err = tx.A(func() error {
		return filesClient.CreateFile("/path/to/file", "hello world")
	}, func() error {
		return filesClient.RemoveFile("/path/to/file")
	})
	if err != nil {
		originalErr, rollbackErr := tx.R(err)
		fmt.Printf("something went wrong with Files. Rollbacked. Err: %v", originalErr)
		if rollbackErr != nil {
			panic(rollbackErr)
		}
		return
	}

	slog.Info(
		"state of previous clients before the fail",
		"database", databaseClient.Records,
		"files", filesClient.Files,
	)
	// 3.
	err = tx.A(func() error {
		return stripeClient.CreateAccount("user-account")
	}, func() error {
		return stripeClient.RemoveAccount("user-account")
	})
	if err != nil {
		_, rollbackErr := tx.R(err)
		if rollbackErr != nil {
			panic(rollbackErr)
		} else {
			slog.Info(
				"state of previous clients after the fail",
				"database", databaseClient.Records,
				"files", filesClient.Files,
			)
		}
		return
	}
	// 2025/07/15 02:17:57 INFO state of previous clients before the fail database=[customer] files="map[/path/to/file:hello world]"
	// 2025/07/15 02:17:57 INFO state of previous clients after the fail database=[] files=map[]
	// As you can see from above, both of the actions were rolled back when the StripeClient errored.
}
