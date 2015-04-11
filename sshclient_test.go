package sshclient

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	c := NewSshClient("localhost", "ariefdarmawan")
	defer c.Close()
	for i := 1; i < 1000; i++ {
		if s, e := c.Run("whoami"); e != nil {
			t.Error(e.Error())
		} else {
			fmt.Println(i, "=", s)
		}
	}
}
