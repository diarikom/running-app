package nsql

type Field struct {
	Tracked   bool
	Anonymous bool
	Column    string
	Name      string
}

func (c *Field) setColumn(tagVal string) {
	if tagVal != "" {
		c.Column = tagVal
		return
	}
	c.Column = c.Name
}
