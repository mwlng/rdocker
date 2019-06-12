package sshtunnel

import (
	"errors"
	"fmt"
)

type IPAddr struct {
	Host string
	Port int
}

type UnixAddr struct {
	SockPath string
}

type Endpoint struct {
	Proto string
	Addr  interface{}
}

func (endpoint *Endpoint) ToString() (string, error) {
	switch endpoint.Proto {
	case "tcp", "udp":
		addr, ok := endpoint.Addr.(IPAddr)
		if !ok {
			return "", errors.New("Invalid endpoint address")
		}
		return fmt.Sprintf("%s:%d", addr.Host, addr.Port), nil
	case "unix":
		addr, ok := endpoint.Addr.(UnixAddr)
		if !ok {
			return "", errors.New("Invalid endpoint address")
		}
		return fmt.Sprintf("%s", addr.SockPath), nil
	}
	return "", errors.New("Unknown protocol type")
}
