package helper

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
)

type SSH struct {

}

func (vssh *SSH) NewSSHClient(){
	config := ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			vssh.ParsePublicKey("/Users/singaravelannandakumar/github/src/github.com/coordination-service/resource_manager/velan-kubespray.pem"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", "3.15.140.234:22",&config)
	if err != nil{
		fmt.Println(err)
		return
	}
	//"sudo docker run --publish 50052:50052 --network host --detach --name rm singaravelan21/resourcemanager:v1.2.0"
	vssh.RunCommand([]string{"sudo adduser velan"}, client)
}

func (vssh *SSH) ParsePublicKey(path string)ssh.AuthMethod{
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func (vssh *SSH) RunCommand(cmd []string, conn *ssh.Client) {
	for _, v := range cmd{
		sess, err := conn.NewSession()
		if err != nil {
			panic(err)
		}
		defer sess.Close()
		sessStdOut, err := sess.StdoutPipe()
		if err != nil {
			panic(err)
		}
		go io.Copy(os.Stdout, sessStdOut)
		sessStderr, err := sess.StderrPipe()
		if err != nil {
			panic(err)
		}
		go io.Copy(os.Stderr, sessStderr)
		err = sess.Run(v) // eg., /usr/bin/whoami
		if err != nil {
			panic(err)
		}
		fmt.Println("Suce")
	}
}

