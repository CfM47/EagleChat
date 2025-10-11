package controller

func (c *Controller) setError(msg string) {
	c.error = msg
	c.render()
}
