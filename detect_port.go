package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/smtp"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BaseConfig struct {
	To       []string `json:"to"`
	From     string   `json:"from"`
	Password string   `json:"password"`
}

var config BaseConfig
var dpArray = map[string]interface{}{}
var sendtext string
var waitGroup sync.WaitGroup
var gameWait *sync.WaitGroup

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	if len(os.Args) < 3 {
		fmt.Println("opt is missing. [exp:] detect_port 5(sec for test connection time out) 10(min for test loop time)")
		os.Exit(1)
	}

	configfile, e := ioutil.ReadFile("dp.config")
	if e != nil {
		fmt.Printf("dp.config File read error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(configfile, &config)
	fmt.Printf("dp.config load finish.\n")

	dpfile, e := ioutil.ReadFile("dpconfig.json")
	if e != nil {
		fmt.Printf("dpconfig.json File read error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(dpfile, &dpArray)
	fmt.Printf("dpconfig.json load finish.\n")
	min, _ := strconv.Atoi(os.Args[2])
	TaskTimer(min)
}

func TaskTimer(min int) {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
	detect_port()
	ticker := time.NewTicker(time.Duration(min) * 60 * time.Second)
	for _ = range ticker.C {
		defer ticker.Stop()
		timestamp := time.Now().Unix()
		tm := time.Unix(timestamp, 0)
		fmt.Println(tm.Format("2006-01-02 15:04:05"))
		go detect_port()
	}
}

func detect_port() {
	for k := range dpArray {
		value := ToSlice(dpArray[k])
		for _, v := range value {
			waitGroup.Add(1)
			fmt.Printf("Ip [%s] Port [%s]\n", k, v)
			iptrue := validIP4(k)
			if iptrue {
				go TestConnStatus(k, v.(string), &waitGroup)
			}
		}
	}
	waitGroup.Wait()
	if sendtext != "" {
		send(sendtext)
		sendtext = ""
	}
}

func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

func TestConnStatus(server string, port string, wait *sync.WaitGroup) {
	timeout, _ := strconv.Atoi(os.Args[1])
	conn, err := net.DialTimeout("tcp", server+":"+port, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Println("Connection error:", err)
		sendtext += server + ":" + port + " Unreachable" + "\n"
	} else {
		log.Println(server + ":" + port + " Online")
		defer conn.Close()
	}
	wait.Done()
	gameWait = wait
}

func send(body string) {
	from := config.From
	pass := config.Password
	to := config.To

	err := SendToMail(from, pass, "smtp.gmail.com:587", to, "Server STATUS", body, "")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}

	log.Print("sent mail!")
	sendtext = ""
}

func SendToMail(user string, password string, host string, to []string, subject string, body string, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to[0] + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail(host, auth, user, to, msg)
	return err
}

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}
