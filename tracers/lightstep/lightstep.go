package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	lightstep "github.com/lightstep/lightstep-tracer-go"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	defComponentName = "skipper"
)

func InitTracer(opts []string) (opentracing.Tracer, error) {
	componentName := defComponentName
	var port int
	var host, token string

	for _, o := range opts {
		parts := strings.SplitN(o, "=", 2)
		switch parts[0] {
		case "component-name":
			if len(parts) > 1 {
				componentName = parts[1]
			}
		case "token":
			token = parts[1]
		case "collector":
			var err error
			var sport string

			host, sport, err = net.SplitHostPort(parts[1])
			if err != nil {
				return nil, err
			}

			port, err = strconv.Atoi(sport)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s as int: %s", sport, err)
			}
		}
	}
	if token == "" {
		return nil, errors.New("missing token= option")
	}

	// Set defaults.
	if host == "" {
		host = lightstep.DefaultGRPCCollectorHost
		port = lightstep.DefaultSecurePort
	}

	opts2 := lightstep.Options{
		AccessToken: token,
		Collector: lightstep.Endpoint{
			Host: host,
			Port: port,
		},
		UseGRPC: true,
		Tags: map[string]interface{}{
			lightstep.ComponentNameKey: componentName,
		},
	}
	fmt.Printf("%#v", opts2)
	return lightstep.NewTracer(opts2), nil
}
