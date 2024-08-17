package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := os.Getenv("RESPONSE")
	if len(response) == 0 {
		response = "You are using the following network interfaces: \n\n"
		response += getAllNetworkInterfaces()
		response += "\n\n"
		response += "Your current IP address is: " + getCurrentIPAddress() + "\n"
	}

	fmt.Fprintln(w, response)
	fmt.Println("Servicing an impatient beginner's request.")
}

func listenAndServe(port string) {
	fmt.Printf("serving on %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func getCurrentIPAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range interfaces {
		if strings.HasPrefix(iface.Name, "Ethernet") {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && ipNet.IP.To4() != nil {
					logrus.Infof("Found Ethernet interface: %s with IP: %s", iface.Name, addr.String())
					return ipNet.IP.String()
				}
			}
		}
	}
	return ""
}

func sayHello() {
	fmt.Println("Hallo Ute")
	fmt.Println("The current time is:", time.Now().Format("15:04:05"))
}

func getAllNetworkInterfaces() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var result string
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.To4() != nil {
				result += fmt.Sprintf("Interface: %s, IP: %s\n", iface.Name, ipNet.IP.String())
			}
		}
	}
	return result
}

func main() {
	http.HandleFunc("/", helloHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	go listenAndServe(port)

	select {}
}
