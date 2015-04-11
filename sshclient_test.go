package sshclient

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	c := NewSshClient("localhost", "ariefdarmawan")
	defer c.Close()
	if s, e := c.Run("ls -al"); e != nil {
		t.Error(e.Error())
	} else {
		fmt.Println(s)
	}
}
