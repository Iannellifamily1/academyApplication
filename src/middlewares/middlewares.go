package middlewares

import (
	"academyApplication/src/database"
	"academyApplication/src/security"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func IsAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		staffClaims, err := security.VerifyToken(tokenString)
		fmt.Println("my User type:" + staffClaims.Type)
		if err != nil || staffClaims.Type != "Admin" {
			http.Error(w, "Invalid token or insufficient permissions", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func IsAdminOrSelf(daoImpl *database.DAOImpl, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		staffID, err := strconv.Atoi(vars["staffID"])
		if err != nil {
			http.Error(w, "Invalid staff ID", http.StatusBadRequest)
			return
		}
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		staffClaims, err := security.VerifyToken(tokenString)
		fmt.Println("my User type:" + staffClaims.Type)
		//If there was an error or registered staff is different than the staff that have made the request
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		staffFromDB, err := daoImpl.GetStaffByEmail(staffClaims.Email)

		if staffClaims.Type != "Admin" && staffFromDB.ID != staffID {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
