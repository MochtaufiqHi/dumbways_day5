package connection

import (
	"context"
	"fmt"
	"testing"
)

// Mock implementation of pgx.Conn for testing purposes
type mockConn struct{}

func (c *mockConn) Close(context.Context) {}

// TestDatabaseConnect tests the DatabaseConnect function
func TestDatabaseConnect(t *testing.T) {
	// Create a mock connection for testing
	Con := &mockConn{}

	// Call the DatabaseConnect function
	DatabaseConnect()

	// Check if the connection was successful
	if Con == nil {
		t.Error("Connection was not established")
	}

	// Check if the success message was printed
	expectedMsg := ""
	if captured := fmt.Sprint(mockStdout.String()); captured != expectedMsg {
		t.Errorf("Expected message: %s, got: %s", expectedMsg, captured)
	} else {
		a := "Unable connect to database"
		b := "Unable connect to database"
		t.Errorf("Expected message: %s, got: %s", a, b)
	}
}

// Capture stdout for testing
var mockStdout mockWriter

type mockWriter struct {
	captured string
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	w.captured += string(p)
	return len(p), nil
}

func (w *mockWriter) String() string {
	return w.captured
}

// Mock implementation of pgx.Conn for testing purposes
// var Conn pgx.Conn

// func Test2DatabaseConnect() {
// 	databaseUrl := "postgresql://postgres:taufiq97@localhost:5432/db_project"

// 	var err error
// 	Conn, err = pgx.Connect(context.Background(), databaseUrl)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable connect to database: %v\n", err)
// 		os.Exit(1)
// 	}

// 	fmt.Println("Success connect to database!")
// }

// func main() {
// 	// Run the tests
// 	result := testing.Main(matchString([]string{"-test.run=TestDatabaseConnect"}), nil, nil, nil)
// 	os.Exit(result)
// }
