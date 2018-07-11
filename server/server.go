package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WhoSV/codestack-api/database"
	"github.com/WhoSV/codestack-api/router"
)

// ExecuteServer ...
func ExecuteServer() {
	database.ConnectDB()

	// after func finishies
	defer database.CloseDB()

	fmt.Println("API server version", "1.2.0", "is listening on port", "8000")

	log.Fatal(http.ListenAndServe("localhost:8000", router.GetRouter()))

}
