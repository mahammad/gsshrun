package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	Ip          string
	Username    string
	Project     string
	Password    string
	Command     string
	Description string
	Port        int64
}

//Connection configuration
type ClientConfig struct {
	Host       string      //ip
	Port       int64       // Port
	Username   string      //Username
	Password   string      //Password
	Client     *ssh.Client //ssh client
	LastResult string      //Result of the last run
}

func main() {

	// var err error

	if len(os.Args) > 1 {
		var server Connection
		project := os.Args[1]
		server = find(project)

		cliConf := new(ClientConfig)
		cliConf.createClient(server.Ip, server.Port, server.Username, server.Password)

		fmt.Println(server.Description)
		fmt.Println("===========================")
		fmt.Println(cliConf.RunShell(server.Command))

	} else {
		fmt.Println("Sunucu ismi belirtiniz...!")
	}
}

// Sunucu bilgileri
func find(project string) Connection {
	connectionlist := read_connections()
	for i := range connectionlist {
		if connectionlist[i].Project == project {
			return connectionlist[i]
		}
	}
	var emty Connection
	return emty
}

// Sunucularımız
func read_connections() []Connection {
	content, err := ioutil.ReadFile("connections.json") // dosya yolunu belirmeniz gerekiyor
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload []Connection
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}

// // SSH Komutu
// func make_ssh(server Connection) string {
// 	return "ssh " + server.Username + "@" + server.Ip + " -p " + server.Port
// }

func (cliConf *ClientConfig) createClient(host string, port int64, username, password string) {
	var (
		client *ssh.Client
		err    error
	)
	cliConf.Host = host
	cliConf.Port = port
	cliConf.Username = username
	cliConf.Password = password
	cliConf.Port = port

	//Generally pass in four parameters: user, []ssh.AuthMethod{ssh.Password(password)}, HostKeyCallback, timeout,
	config := ssh.ClientConfig{
		User: cliConf.Username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)

	//Get client
	if client, err = ssh.Dial("tcp", addr, &config); err != nil {
		log.Fatalln("error occurred:", err)
	}

	cliConf.Client = client
}

func (cliConf *ClientConfig) RunShell(shell string) string {
	var (
		session *ssh.Session
		err     error
	)

	//Get session, this session is used to perform operations remotely
	if session, err = cliConf.Client.NewSession(); err != nil {
		log.Fatalln("error occurred:", err)
	}

	//Execute shell
	if output, err := session.CombinedOutput(shell); err != nil {
		log.Fatalln("error occurred:", err)
	} else {
		cliConf.LastResult = string(output)
	}
	return cliConf.LastResult
}
