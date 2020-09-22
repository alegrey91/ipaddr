package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Program name
const programName = "IPAddress Microservice"

// Show Help.
func help(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("Received request from: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "%s\n", programName)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "[Description]\n")
	fmt.Fprintf(w, "The following microservice show IP informations about\n")
	fmt.Fprintf(w, "the interfaces present on your system (in YAML format).\n")
	fmt.Fprintf(w, "(Powered by alegrey91)\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "[Usage]\n")
	fmt.Fprintf(w, "To use the %s, just try to curl the following endopoints:\n", programName)
	fmt.Fprintf(w, "http://example-server.com:8080/ipa    → show interfaces with IPv4 and IPv6 addresses\n")
	fmt.Fprintf(w, "http://example-server.com:8080/ipa/4  → show interfaces with IPv4 addresses\n")
	fmt.Fprintf(w, "http://example-server.com:8080/ipa/6  → show interfaces with IPv6 addresses\n")
}

// Determine the address type (IPNet or IPAddr).
func detectAddrType(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	ip = ip.To4()
	return ip
}

// Wrapper of net.Interfaces to manage errors.
func getInterfaces() []net.Interface {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println(fmt.Errorf("localAddresses: %v", err.Error()))
	}
	return ifaces
}

// Wrapper of net.Addrs to manage errors.
func getAddresses(iface net.Interface) []net.Addr {
	addrs, err := iface.Addrs()
	if err != nil {
		log.Println(fmt.Errorf("localAddresses: %v", err.Error()))
	}
	return addrs
}

// Binded to the endpoint /ipa,
// this function list all the IP,
// for each interface present on the system.
func listIP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("Received request from: %s\n", r.RemoteAddr)

	ifaces := getInterfaces()

	for _, iface := range ifaces {
		addrs := getAddresses(iface)
		name := iface.Name
		if len(addrs) > 0 {
			fmt.Fprintf(w, "%s:\n", name)
			if len(addrs) > 0 {
				for _, addr := range addrs {
					fmt.Fprintf(w, "- %s\n", addr.String())
				}
			}
			fmt.Fprintf(w, "\n")
		}
	}
}

// Binded to the endpoint /ipa/:type,
// this function list the IPv4 or IP6
// (depending on the ':type' paramenter),
// for each interface present on the system.
func listIPv46(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if ps.ByName("type") != "4" && ps.ByName("type") != "6" {
		log.Println(fmt.Errorf("Wrong request received at %s", r.URL.Path))
		fmt.Fprintf(w, "Wrong endpoint request\n")
	}
	log.Printf("Received request from: %s\n", r.RemoteAddr)

	ifaces := getInterfaces()

	for _, iface := range ifaces {
		addrs := getAddresses(iface)
		name := iface.Name
		if len(addrs) > 0 {
			for _, addr := range addrs {
				if ps.ByName("type") == "4" {
					ip := detectAddrType(addr)
					// if ip is nil then it's ipv6
					if ip == nil {
						continue
					}
					fmt.Fprintf(w, "%s:\n", name)
					fmt.Fprintf(w, "- %s\n", addr.String())
				}
				if ps.ByName("type") == "6" {
					ip := detectAddrType(addr)
					// if ip is not nil then it's ipv4
					if ip != nil {
						continue
					}
					fmt.Fprintf(w, "%s:\n", name)
					fmt.Fprintf(w, "- %s\n", addr.String())
				}
				fmt.Fprintf(w, "\n")
			}
		}
	}
}

func main() {

	var port int
	flag.IntVar(&port, "port", 8080, "Bind the microservice to the specified port.")
	flag.Parse()

	log.Printf("%s running on port %d\n", programName, port)

	router := httprouter.New()
	router.GET("/", help)
	router.GET("/ipa", listIP)
	router.GET("/ipa/:type", listIPv46)
	log.Fatalf("Server failed with: %s\n", http.ListenAndServe(":"+strconv.Itoa(port), router))
}
