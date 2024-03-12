package main

import (
	"fmt"
	"miras/internal/controllers"
	"miras/internal/services"
	"miras/internal/transport/rest"
	"net/http"
)

func main() {

	db, err := controllers.NewDB()
	if err != nil {
		panic(err)
	}
	cache, err := controllers.NewRedis()
	if err != nil {
		panic(err)
	}
	repo := controllers.NewRepository(db, cache)
	service := services.NewService(repo)
	handlers := rest.NewHandler(service)

	router := rest.NewRouter(handlers)

	http.Handle("/", router)
	fmt.Println("Server working...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}

}
