package example

type FilesClient struct {
	Files map[string]string
}

func (c *FilesClient) InitFs() {
	c.Files = map[string]string{}
}

func (c *FilesClient) CreateFile(path, contents string) error {
	c.Files[path] = contents
	return nil
}

func (c *FilesClient) RemoveFile(path string) error {
	delete(c.Files, path)
	return nil
}
