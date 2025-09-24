package network

// ClientInterface defines the interface for client operations
type ClientInterface interface {
	Write(message string)
	SetState(state int)
}

// SetState sets the state of the client
func (c *Client) SetState(state int) {
	c.State = ConnectionState(state)
}
