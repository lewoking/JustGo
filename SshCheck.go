package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"
)

type Website struct {
	Host    string `xml:"name,attr"`
	User    string
	Pwd     string
	Type    string
	KeyPath string
	Port    int
	Rcmd    string
}

func main() {
	filePtr, err := os.Open("./config.json")
	if err != nil {
		fmt.Println("文件打开失败 [Err:%s]", err.Error())
		return
	}
	defer filePtr.Close()
	var config []Website
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("解码失败", err.Error())
	} else {
		fmt.Println("解码成功")
		// fmt.Println(config)
	}
	for _, cfg := range config {
		// fmt.Println(cfg)
		v := reflect.ValueOf(cfg)
		fmt.Println(v.Field(0).Interface())
		sshHost := v.Field(0).Interface()
		sshUser := v.Field(1).Interface().(string)
		sshPassword := v.Field(2).Interface().(string)
		sshType := v.Field(3).Interface()             //password 或者 key
		sshKeyPath := v.Field(4).Interface().(string) //ssh id_rsa.id 路径"
		sshPort := v.Field(5).Interface()

		//创建sshp登陆配置
		config := &ssh.ClientConfig{
			Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
			User:            sshUser,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
		if sshType == "password" {
			config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
		} else {
			config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
		}

		//dial 获取ssh client
		addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
		sshClient, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			log.Fatal("创建ssh client 失败", err)
		}
		defer sshClient.Close()

		//创建ssh-session
		session, err := sshClient.NewSession()
		if err != nil {
			log.Fatal("创建ssh session 失败", err)
		}
		defer session.Close()
		//执行远程命令
		combo, err := session.CombinedOutput(v.Field(6).Interface().(string))
		if err != nil {
			log.Fatal("远程执行cmd 失败", err)
		}
		log.Println("命令输出:", string(combo))

	}
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
