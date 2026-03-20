package config

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	s, _ := GetAccessToken()
	fmt.Println(s)
}
