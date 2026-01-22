package main

import (
	"github.com/gadz82/go-api-boilerplate/internal/di"
	"go.uber.org/fx"
)

func initDb() {

}

// @title           Go API Boilerplate
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

func main() {
	fx.New(
		di.NewModule(),
	).Run()
}
