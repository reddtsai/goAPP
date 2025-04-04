package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	svc := New(context.Background())
	got := svc.GetUser("abc123")
	want := "User(ID=abc123)"

	assert.Equal(t, want, got)
}
