package example

import "fmt"

type StripeClient struct{}

func (c *StripeClient) CreateAccount(account string) error {
	return fmt.Errorf("failed to create an account")
}

func (c *StripeClient) RemoveAccount(account string) error {
	return nil
}
