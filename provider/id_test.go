package provider

import (
	"crypto/rand"
	"fmt"
)

func generateID(n int) string {
	id := make([]byte, n)
	rand.Read(id)
	return fmt.Sprintf("%x", id)
}
