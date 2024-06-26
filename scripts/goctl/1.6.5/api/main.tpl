package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	
	"github.com/qnfnypen/gzocomm/mhandler"

	{{.importPackages}}
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf,generateNoFoundHandler())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func generateNoFoundHandler() rest.RunOption {
	fp := "swagger.json"
	info := &mhandler.SwaggerInfo{BasePath: "//"}
	url := "/swagger"
	swagHandler, err := mhandler.SwaggerHandler(fp, url, info)
	if err != nil {
		return rest.WithNotFoundHandler(nil)
	}
	return rest.WithNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// 适配 swagger
		case strings.HasPrefix(r.URL.Path, url):
			swagHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
}