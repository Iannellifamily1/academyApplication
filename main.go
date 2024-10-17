package main

import (
	"academyApplication/src/database"
	"academyApplication/src/handlers"
	"academyApplication/src/middlewares"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer db.Close()

	router := mux.NewRouter()

	daoImpl := &database.DAOImpl{Db: db}
	h := handlers.NewHandler(daoImpl)

	//endpoint to handle login
	router.HandleFunc("/login", h.LoginHandler).Methods("POST")
	router.HandleFunc("/register", h.RegisterHandler).Methods("POST")
	router.HandleFunc("/staff", middlewares.IsAdmin(h.GetStaffHandler)).Methods("GET")
	router.HandleFunc("/staff/{staffID}", middlewares.IsAdminOrSelf(daoImpl, h.GetStaffByIDHandler)).Methods("GET")
	router.HandleFunc("/staff/{staffID}", h.DeleteStaffHandler).Methods("DELETE")
	router.HandleFunc("/staff/{staffID}", h.UpdateStaffHandler).Methods("UPDATE")
	router.HandleFunc("/staff/{staffID}/roles", h.GetRolesByStaffHandler).Methods("GET")
	router.HandleFunc("/staff/{staffID}/roles/{roleID}", h.NewStaffWithRolesHandler).Methods("POST")
	router.HandleFunc("/role", middlewares.IsAdmin(h.NewRoleHandler)).Methods("POST")
	router.HandleFunc("/role", middlewares.IsAdmin(h.GetAllRolesHandler)).Methods("GET")
	router.HandleFunc("/role/{roleID}", middlewares.IsAdmin(h.DeleteRoleHandler)).Methods("DELETE")
	router.HandleFunc("/role/{roleID}", middlewares.IsAdmin(h.UpdateRoleHandler)).Methods("UPDATE")

	router.HandleFunc("/staffWithRoles", h.NewStaffWithRolesHandler).Methods("POST")

	//endpoint for concurrency task
	router.HandleFunc("/wordfrequency", handlers.WordFrequencyHandler).Methods("POST")

	fmt.Println("Server started. Listening on port 3001")
	err = http.ListenAndServe(":3001", router)
	if err != nil {
		log.Fatalf("Could not start the server: %v", err)
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
