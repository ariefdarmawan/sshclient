# sshclient
SSH Client written in Go. It is more like wrapper of crypto/ssh package

This package is inspired from easyssh from hypersleep. I modify so it could be authenticated using password and manage one persistent SSH session client

## Load from repo
```
import "github.com/juragan360/sshclient"
```

## Connect
Connect to SSH Server, by default it will use port 22 and PublicKey authentication method
```
c := sshclient.NewSshClient("localhost","ariefdarmawan")
```

or
```
c := sshclient.NewSshClient("localhost","ariefdarmawan")
c.Password = "mypassword"
c.AuthType = "Password"
```

## Run Command
```
if s, e := c.Run("ls -al /Users/ariefdarmawan"); e != nil {
	t.Error(e.Error())
} else {
	fmt.Println(s)
}
```

## Close
Dont forget to defer close the connection once done
```
c := sshclient.NewSshClient("localhost","ariefdarmawan")
defer c.Close()
```

