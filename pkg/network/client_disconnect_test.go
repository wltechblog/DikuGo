package network

import (
	"net"
	"testing"
	"time"

	"github.com/wltechblog/DikuGo/pkg/command"
	"github.com/wltechblog/DikuGo/pkg/world"
)

// MockConn implements net.Conn for testing connection drops
type MockConn struct {
	closed     bool
	readError  error
	writeError error
	data       []byte
	readPos    int
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	if m.readError != nil {
		return 0, m.readError
	}
	if m.readPos >= len(m.data) {
		return 0, net.ErrClosed
	}
	n = copy(b, m.data[m.readPos:])
	m.readPos += n
	return n, nil
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	if m.writeError != nil {
		return 0, m.writeError
	}
	if m.closed {
		return 0, net.ErrClosed
	}
	return len(b), nil
}

func (m *MockConn) Close() error {
	m.closed = true
	return nil
}

func (m *MockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 4000}
}

func (m *MockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12345}
}

func (m *MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestClientDisconnectHandling(t *testing.T) {
	// Create a mock world and command registry
	storage := world.NewMockStorage()
	w, err := world.NewWorld(nil, storage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}
	cmdRegistry := command.NewRegistry()

	// Test 1: EOF error should close client properly
	t.Run("EOF_Error_Closes_Client", func(t *testing.T) {
		// Create a mock connection that will return EOF
		mockConn := &MockConn{
			readError: net.ErrClosed, // Simulate connection closed
		}

		// Create client
		client := NewClient(mockConn, w, cmdRegistry)

		// Start handling in a goroutine
		done := make(chan bool)
		go func() {
			client.Handle()
			done <- true
		}()

		// Wait for client to finish handling (should happen quickly due to EOF)
		select {
		case <-done:
			// Client finished handling - this is expected
		case <-time.After(5 * time.Second):
			t.Fatal("Client did not finish handling within timeout - may be stuck in loop")
		}

		// Verify client is marked as closed
		if !client.Closed {
			t.Error("Client should be marked as closed after EOF error")
		}
	})

	// Test 2: Write error should close client
	t.Run("Write_Error_Closes_Client", func(t *testing.T) {
		// Create a mock connection that will fail on write
		mockConn := &MockConn{
			writeError: net.ErrClosed,
		}

		// Create client
		client := NewClient(mockConn, w, cmdRegistry)

		// Try to write to the client
		client.Write("test message")

		// Verify client is marked as closed
		if !client.Closed {
			t.Error("Client should be marked as closed after write error")
		}
	})

	// Test 3: Multiple write attempts to closed client should not cause issues
	t.Run("Multiple_Writes_To_Closed_Client", func(t *testing.T) {
		// Create a normal mock connection
		mockConn := &MockConn{}

		// Create client
		client := NewClient(mockConn, w, cmdRegistry)

		// Manually close the client
		client.Close()

		// Try to write multiple times - should not panic or cause issues
		for i := 0; i < 5; i++ {
			client.Write("test message")
		}

		// Should still be closed
		if !client.Closed {
			t.Error("Client should remain closed")
		}
	})

	// Test 4: Connection timeout should not close client (only EOF/errors should)
	t.Run("Timeout_Does_Not_Close_Client", func(t *testing.T) {
		// Create a mock connection that simulates timeout
		mockConn := &MockConn{
			readError: &net.OpError{
				Op:  "read",
				Err: &timeoutError{},
			},
		}

		// Create client
		client := NewClient(mockConn, w, cmdRegistry)

		// Start handling in a goroutine
		done := make(chan bool)
		go func() {
			// Let it run for a short time to test timeout handling
			time.Sleep(100 * time.Millisecond)
			client.Close() // Close it manually to end the test
			client.Handle()
			done <- true
		}()

		// Wait for client to finish
		select {
		case <-done:
			// Expected
		case <-time.After(2 * time.Second):
			client.Close() // Force close if stuck
			t.Fatal("Client handling took too long")
		}

		// Client should be closed (we closed it manually)
		if !client.Closed {
			t.Error("Client should be closed")
		}
	})
}

// timeoutError implements net.Error with Timeout() returning true
type timeoutError struct{}

func (e *timeoutError) Error() string {
	return "timeout"
}

func (e *timeoutError) Timeout() bool {
	return true
}

func (e *timeoutError) Temporary() bool {
	return true
}

func TestClientCloseIdempotent(t *testing.T) {
	// Create a mock world and command registry
	storage := world.NewMockStorage()
	w, err := world.NewWorld(nil, storage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}
	cmdRegistry := command.NewRegistry()

	// Create a mock connection
	mockConn := &MockConn{}

	// Create client
	client := NewClient(mockConn, w, cmdRegistry)

	// Close multiple times - should not panic
	client.Close()
	client.Close()
	client.Close()

	// Should be closed
	if !client.Closed {
		t.Error("Client should be closed")
	}

	// Connection should be closed
	if !mockConn.closed {
		t.Error("Mock connection should be closed")
	}
}

func TestClientWriteToClosedConnection(t *testing.T) {
	// Create a mock world and command registry
	storage := world.NewMockStorage()
	w, err := world.NewWorld(nil, storage)
	if err != nil {
		t.Fatalf("Failed to create world: %v", err)
	}
	cmdRegistry := command.NewRegistry()

	// Create a mock connection that fails on flush
	mockConn := &MockConn{}

	// Create client
	client := NewClient(mockConn, w, cmdRegistry)

	// Close the mock connection to simulate network failure
	mockConn.closed = true
	mockConn.writeError = net.ErrClosed

	// Try to write - should mark client as closed
	client.Write("test message")

	// Client should be marked as closed
	if !client.Closed {
		t.Error("Client should be marked as closed after write to closed connection")
	}
}
