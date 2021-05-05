package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/leon332157/replish/netstat"
	server "github.com/leon332157/replish/server"
	toml "github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
)

var (
	dotreplit DotReplit
	port      uint16
	//hasReplishField bool = false
)

type DotReplit struct {
	Run      string
	Language string
	onBoot   string
	packager map[string]interface{}
	Replish  map[string]interface{}
}

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)
	dotreplit = loadDotreplit(loadDotreplitFile())
	//go startHijack()
	startFiber()
	time.Sleep(1 * time.Second) // wait for server to come online
	getPort()
	log.Debugf("Got port: %v\n", port)
	go server.StartForwardServer(port)
	go server.StartReverseProxy(port)
	go server.StartSSHServer()
	for {
		time.Sleep(1 * time.Second)
	}
}
func startFiber() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true, DisableKeepalive: false})

	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendString("haha")
	})

	app.Post("/*", func(c *fiber.Ctx) error {
		return c.SendString("haha")
	})

	go app.Listen("127.0.0.1:7373")
	fmt.Println("fiber started")
}

func getPortAuto() {
	addrs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		if s.Process == nil { // Process can be nil, discard it
			return false
		} else if strings.Contains(s.Process.Name, "System") {
			return false // Discard System process
		}
		return net.IP.IsLoopback(s.LocalAddr.IP) && s.State == netstat.Listen
	})
	if err != nil {
		log.Fatalf("Reading ports failed:%v", err)
	}
	if len(addrs) == 0 {
		log.Fatalf("Looks like we aren't finding any open ports, are you listening on localhost (127.0.0.1)?")
	} else if len(addrs) > 1 {
		fmt.Printf("Multiple ports detected: %v\n", len(addrs))
		for index, sock := range addrs {
			if sock.Process != nil {
				fmt.Printf("%v. %v %v\n", index+1, sock.Process, sock.LocalAddr.Port)
			}
		}
		fmt.Print("Choose port/process: ")
		reader := bufio.NewReader(os.Stdin)
		inp, err := reader.ReadString('\n')
		inp = strings.TrimSuffix(inp, "\r\n")
		if err != nil {
			log.Panic(err)
		}
		sel, err := strconv.ParseInt(inp, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		if sel > int64(len(addrs)) { // Input is a port selection
			port = checkPort(sel)
		} else { // Input is list index
			temp := addrs[sel-1]
			port = temp.LocalAddr.Port
		}
	} else {
		port = addrs[0].LocalAddr.Port
	}
}

func checkPort(p int64) uint16 {
	if p > 65535 || p < 1 {
		log.Fatalf("port %v is out of range(1-65535)", p)
	}
	return uint16(p)
}

func getPort() {
	log.Debug("Getting port")
	rawPort, ok := dotreplit.Replish["port"] // Check if port exist
	if !ok {
		log.Warn("Port is missing, defaulting to auto")
		getPortAuto()
		return
	}
	intPort, ok := rawPort.(int64)
	if ok { // port is int
		port = checkPort(intPort)
		return
	}
	strPort, ok := rawPort.(string)
	if ok {
		// Port is string
		if strPort == "auto" {
			getPortAuto()
			return
		} else {
			intPort, err := strconv.Atoi(strPort)
			temp := int64(intPort)
			if err == nil {
				port = checkPort(temp)
			} else {
				log.Errorf("Error when converting port: %v, defaulting to auto\n", err)
				getPortAuto()
				return
			}
		}
	}
}

func loadDotreplitFile() []byte {
	slug, ok := os.LookupEnv("REPL_SLUG")
	var path string
	if ok {
		path = fmt.Sprintf("/home/runner/%v/.replit", slug)
	} else {
		path = ".replit"
	}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		contents = make([]byte, 0)
	}
	return contents
}

// Use another function to allow for unit tests
func loadDotreplit(contents []byte) DotReplit {
	temp := DotReplit{}
	err := toml.Unmarshal(contents, &temp)
	if err != nil {
		log.Fatalf("failed to unmarshal: %v\n", err)
	}
	if temp.Replish == nil {
		log.Warn("Replish field is empty or doesn't exist! Check for typos in .replit")
		// Write replish field maybe
	}
	return temp
}
