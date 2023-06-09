package main

import (
	"log"

	config "github.com/stebinsabu13/ecommerce-api/pkg/config"
	di "github.com/stebinsabu13/ecommerce-api/pkg/di"
)

//	@title			SPORTZONE_E-COMMERCE REST API
//	@version		2.0
//	@description	SPORTZONE_E-COMMERCE REST API built using Go, PSQL, REST API following Clean Architecture.

//	@contact
// name: Stebin Sabu
// url: https://github.com/stebinsabu13
// email: stebinsabu369@gmail.com

//	@license
// name: MIT
// url: https://opensource.org/licenses/MIT

//	@host	sportzone.cloud

// @Basepath	/
// @Accept		json
// @Produce	json
// @Router		/ [get]
func main() {
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}
	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
