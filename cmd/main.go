package main

import (
	"flag"
	"github.com/FRahimov84/FileService/cmd/app"
	"github.com/FRahimov84/FileService/pkg/core/file"
	"github.com/FRahimov84/Mux/pkg/mux"
	"github.com/FRahimov84/di/pkg/di"
	"github.com/FRahimov84/myJwt/pkg/jwt"
	"net"
	"net/http"
	"os"
)

var (
	host = flag.String("host", "", "Server host")
	port = flag.String("port", "", "Server port")
)

const (
	envHost = "HOST"
	envPort = "PORT"
)

func main() {
	flag.Parse()
	serverHost := checkENV(envHost, *host)
	serverPort := checkENV(envPort, *port)
	addr := net.JoinHostPort(serverHost, serverPort)
	secret := jwt.Secret("secret")
	start(addr, secret)
}

func checkENV(env string, loc string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		return loc
	}
	return str
}
func start(addr string, secret jwt.Secret) {
	container := di.NewContainer()
	container.Provide(
		func() string {return "files"},
		func() jwt.Secret { return secret },
		app.NewServer,
		file.NewService,
		mux.NewExactMux,
	)
	container.Start()
	var appServer *app.Server
	container.Component(&appServer)

	panic(http.ListenAndServe(addr, appServer))
}
