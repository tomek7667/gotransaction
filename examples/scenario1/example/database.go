package example

// just some database with state
type DatabaseClient struct {
	Records []any
}

// returns the index of inserted item
func (c *DatabaseClient) Insert(data any) (int, error) {
	c.Records = append(c.Records, data)
	return len(c.Records) - 1, nil
}

func (c *DatabaseClient) Delete(i int) error {
	c.Records = append(c.Records[:i], c.Records[i+1:]...)
	return nil
}
