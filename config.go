package main

import (
	"flag"
	"github.com/ghodss/yaml"
	"github.com/micro/mdns"
	"github.com/valyala/gorpc"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type Node struct {
	Client  *gorpc.Client
	Service *mdns.ServiceEntry
}

type Config struct {
	Host string
	Self string

	Node      string `json:"node"`
	Join      string `json:"join"`
	IFace     string `json:"iface"`
	Discovery int    `json:"discovery"`
	Port      int    `json:"port"`

	Net *net.Interface
	IPs []net.IP

	Nodes map[string]*Node
}

func defaultConfig() Config {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	config := Config{
		Node:      host,
		Host:      host,
		Port:      8181,
		Discovery: 8801,
		IFace:     "eth0",
		Nodes:     make(map[string]*Node),
	}
	return config
}

func LoadConfig() *Config {
	config := Config{}

	defaultConfig := defaultConfig()

	config.Nodes = defaultConfig.Nodes

	if len(os.Args) > 1 {
		configFile := os.Args[len(os.Args)-1]

		_, err := os.Stat(configFile)
		var fileConfig Config
		if !os.IsNotExist(err) {
			content, err0 := ioutil.ReadFile(configFile)
			if err0 != nil {
				panic(err0)
			}

			err1 := yaml.Unmarshal(content, &fileConfig)
			if err1 != nil {
				panic(err1)
			}

			if fileConfig.Node != "" {
				defaultConfig.Node = fileConfig.Node
			}
			if fileConfig.Join != "" {
				defaultConfig.Join = fileConfig.Join
			}
			if fileConfig.IFace != "" {
				defaultConfig.IFace = fileConfig.IFace
			}
			if fileConfig.Discovery != 0 {
				defaultConfig.Discovery = fileConfig.Discovery
			}
			if fileConfig.Port != 0 {
				defaultConfig.Port = fileConfig.Port
			}
		}
	}

	flag.StringVar(&config.Node, "node", defaultConfig.Node, "Name of this node")
	flag.StringVar(&config.Join, "join", defaultConfig.Join, "Address of node to join")
	flag.StringVar(&config.IFace, "iface", defaultConfig.IFace, "Network Interface to bind to")
	flag.IntVar(&config.Discovery, "discovery", defaultConfig.Discovery, "Port for network discovery")
	flag.IntVar(&config.Port, "port", defaultConfig.Port, "Port for cluster conns")
	flag.Parse()

	iface, errFace := net.InterfaceByName(config.IFace)
	if errFace != nil {
		panic(errFace)
	}
	config.Net = iface

	addrs, aerr := config.Net.Addrs()
	if aerr != nil {
		panic(aerr)
	}
	for _, addr := range addrs {
		ip, _, iperr := net.ParseCIDR(addr.String())
		if iperr != nil {
			log.Println(iperr)
		} else {
			log.Println(ip)
			config.IPs = append(config.IPs, ip)
		}
	}

	return &config
}
