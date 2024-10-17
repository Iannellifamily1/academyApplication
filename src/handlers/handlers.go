package handlers

import (
	"academyApplication/src/database"
	"academyApplication/src/models"
	"academyApplication/src/security"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type Handler struct {
	Dao database.DAO
}

func NewHandler(dao database.DAO) *Handler {
	return &Handler{Dao: dao}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var staff models.Staff
	err := json.NewDecoder(r.Body).Decode(&staff)
	if err != nil {
		http.Error(w, "Invalid staff format provided", http.StatusBadRequest)
		return
	}

	staffFromDB, err := h.Dao.GetStaffByEmail(staff.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Staff does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	// Verifica la password
	if !security.CheckPasswordHash(staff.Password, staffFromDB.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Creazione del token e altre operazioni
	// al posto di admin devo prendere il nome del ruolo associato a quell'utente.Inserire la query per fare ciò
	tokenString, err := security.CreateToken(staff.Email, staffFromDB.ID, "User")
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	response := models.StaffResponse{
		ID:    staffFromDB.ID,
		Email: staff.Email,
		Type:  "User", // type that needs to be retrieved from DB
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	//Mapping and sanitization phase
	var staff models.Staff
	err := json.NewDecoder(r.Body).Decode(&staff)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}
	if staff.Email == "" {
		http.Error(w, "Error - email is a mandatory field", http.StatusBadRequest)
		return
	}
	if !IsValidEmail(staff.Email) {
		http.Error(w, "Error - email doesn't match the correct format", http.StatusBadRequest)
		return
	}

	if staff.Name == "" {
		http.Error(w, "Error - name is a mandatory field", http.StatusBadRequest)
		return
	}
	if staff.Password == "" {
		http.Error(w, "Error - password is a mandatory field", http.StatusBadRequest)
		return
	}

	//Check if already exists another staff with the same email.
	_, err = h.Dao.GetStaffByEmail(staff.Email)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, fmt.Sprintf("Error - staff %v already exists", staff.Email))
		return
	}
	if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to retrieve user")
		return
	}

	// Hash della password
	fmt.Println("My password is" + staff.Password)
	hashedPassword, err := security.HashPassword(staff.Password)
	fmt.Println("My hashed password is " + hashedPassword)
	staff.Password = hashedPassword // Sostituisci la password con l'hash

	createdStaff, err := h.Dao.CreateStaff(staff)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert staff: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdStaff)
}

func (h *Handler) GetStaffHandler(w http.ResponseWriter, r *http.Request) {
	staffList, err := h.Dao.GetAllStaff()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get staff: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staffList)
}

func (h *Handler) GetStaffByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	staffID, err := strconv.Atoi(vars["staffID"])
	if err != nil {
		http.Error(w, "Invalid staff ID", http.StatusBadRequest)
		return
	}

	staff, err := h.Dao.GetStaff(staffID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get staff: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staff)
}

func (h *Handler) GetRolesByStaffHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	staffID, err := strconv.Atoi(vars["staffID"])
	if err != nil {
		http.Error(w, "Invalid staff ID", http.StatusBadRequest)
		return
	}

	roles, err := h.Dao.GetRolesByStaff(staffID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get roles: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func (h *Handler) UpdateStaffHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	staffID, err := strconv.Atoi(vars["staffID"])
	if err != nil {
		http.Error(w, "Invalid staff ID", http.StatusBadRequest)
		return
	}
	var staff models.Staff
	err = json.NewDecoder(r.Body).Decode(&staff)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}
	staff.ID = staffID

	updatedStaff, err := h.Dao.UpdateStaff(staff)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update staff: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStaff)
}

func (h *Handler) DeleteStaffHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	staffID, err := strconv.Atoi(vars["staffID"])
	if err != nil {
		http.Error(w, "Invalid staff ID", http.StatusBadRequest)
		return
	}

	err = h.Dao.DeleteStaff(staffID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete staff: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (h *Handler) NewRoleHandler(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}
	if role.Name == "" {
		http.Error(w, "Error - name is a mandatory field", http.StatusBadRequest)
		return
	}

	roleAlreadyExists, err := h.Dao.CheckIfRoleAlreadyExists(role)
	if roleAlreadyExists {
		http.Error(w, fmt.Sprintf("Role with this name already exists: %v", role.Name), http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert role: %v", err), http.StatusInternalServerError)
		return
	}

	createdRole, err := h.Dao.CreateRole(role)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert role: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdRole)
}

func (h *Handler) GetAllRolesHandler(w http.ResponseWriter, r *http.Request) {
	roles, err := h.Dao.GetAllRoles() // Assicurati di avere una funzione GetAllRoles nel DAO
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get roles: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func (h *Handler) UpdateRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID, err := strconv.Atoi(vars["roleID"])
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var role models.Role
	err = json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}
	role.ID = roleID

	updatedRole, err := h.Dao.UpdateRole(role)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update role: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedRole)
}

func (h *Handler) DeleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID, err := strconv.Atoi(vars["roleID"])
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.Dao.DeleteRole(roleID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete role: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (h *Handler) AddRolesToStaffHandler(w http.ResponseWriter, r *http.Request) {

	var staffWithRoles models.StaffWithRoles
	err := json.NewDecoder(r.Body).Decode(&staffWithRoles)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Controlla se lo staff esiste
	staffExists, err := h.Dao.CheckIfStaffAlreadyExists(staffWithRoles.Staff)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check staff existence: %v", err), http.StatusInternalServerError)
		return
	}
	if !staffExists {
		http.Error(w, fmt.Sprintf("Staff with email %v does not exist", staffWithRoles.Staff.Email), http.StatusNotFound)
		return
	}

	// Aggiungi ogni ruolo per lo staff
	for _, role := range staffWithRoles.Roles {
		if role.ID == 0 {
			http.Error(w, "Error - role ID is a mandatory field", http.StatusBadRequest)
			return
		}

		// Controlla se il ruolo esiste
		roleExists, err := h.Dao.CheckIfRoleAlreadyExists(role)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to check role existence: %v", err), http.StatusInternalServerError)
			return
		}
		if !roleExists {
			http.Error(w, fmt.Sprintf("Role with ID %d does not exist", role.ID), http.StatusNotFound)
			return
		}

		// Controlla se il ruolo è già assegnato allo staff
		roleAlreadyAssigned, err := h.Dao.CheckIfStaffRoleAlreadyExists(staffWithRoles.Staff.ID, role.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to check existing role assignment: %v", err), http.StatusInternalServerError)
			return
		}
		if roleAlreadyAssigned {
			http.Error(w, fmt.Sprintf("Staff with ID %d already has the role with ID %d", staffWithRoles.Staff.ID, role.ID), http.StatusConflict)
			return
		}

		// Crea la nuova associazione
		if _, err := h.Dao.AssignRoleToStaff(staffWithRoles.Staff.ID, role.ID); err != nil {
			http.Error(w, fmt.Sprintf("Failed to assign role %d to staff %d: %v", role.ID, staffWithRoles.Staff.ID, err), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(staffWithRoles)
}

func IsValidEmail(email string) bool {
	// Regex per validare un'email
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// hospital section
// GetAllHospitals recupera tutti gli ospedali
func (h *Handler) GetAllHospitals(w http.ResponseWriter, r *http.Request) {
	hospitalList, err := h.Dao.GetAllHospitals()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get hospitals: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hospitalList)
}

// DeleteHospital elimina un ospedale in base all'ID
func (h *Handler) DeleteHospital(w http.ResponseWriter, r *http.Request) {
	// Estrai l'ID dall'URL (assumendo che sia passato come parametro)
	vars := mux.Vars(r)
	hospitalID, err := strconv.Atoi(vars["hospitalID"])
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.Dao.DeleteHospital(hospitalID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete hospital: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Restituisci 204 No Content
}

// UpdateHospital aggiorna un ospedale esistente
func (h *Handler) UpdateHospital(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hospitalID, err := strconv.Atoi(vars["hospitalID"]) // Estrai l'ID dall'URL
	if err != nil {
		http.Error(w, "Invalid hospital ID", http.StatusBadRequest)
		return
	}

	var hospital models.Hospital
	err = json.NewDecoder(r.Body).Decode(&hospital) // Decodifica il corpo della richiesta
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}
	hospital.ID = hospitalID // Imposta l'ID dell'ospedale aggiornato

	// Chiama il DAO per aggiornare l'ospedale
	updatedHospital, err := h.Dao.UpdateHospital(hospital)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update hospital: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedHospital) // Restituisci l'ospedale aggiornato
}

// AddHospital aggiunge un nuovo ospedale
func (h *Handler) NewHospital(w http.ResponseWriter, r *http.Request) {
	var hospital models.Hospital
	if err := json.NewDecoder(r.Body).Decode(&hospital); err != nil {
		http.Error(w, fmt.Sprintf("Invalid input: %v", err), http.StatusBadRequest)
		return
	}

	createdHospital, err := h.Dao.AddHospital(hospital)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add hospital: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdHospital) // Restituisci l'ospedale creato
}

//end hospital section

// concurrent task
var wordStats = map[string]*models.WordData{}

const filesDirectory = "./files/"

func WordFrequencyHandler(w http.ResponseWriter, r *http.Request) {
	var wordsRequest []string
	err := json.NewDecoder(r.Body).Decode(&wordsRequest)
	if err != nil {
		http.Error(w, "Invalid request payload - please provide a list of strings", http.StatusBadRequest)
		return
	}

	result := make(map[string]models.TFDFEntry)

	files, err := os.ReadDir(filesDirectory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading files - %v", err), http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup //sincronizza l'esecuzione di tutte le goroutine - prima di proseguire attende che tutte abbiano finito
	tfChan := make(chan map[string]int, len(files))
	dfChan := make(chan map[string]int, len(files))
	errorChan := make(chan error, len(files))

	for _, file := range files {
		//for each file i start a new goroutine
		wg.Add(1)
		go func(file fs.DirEntry) {
			defer wg.Done() //-1 nel counter del WaitGroup

			fileInfo, err := file.Info()
			if err != nil {
				errorChan <- err
				return
			}

			tfMap := make(map[string]int)
			dfMap := make(map[string]int)

			content, err := os.ReadFile(filesDirectory + fileInfo.Name())
			if err != nil {
				errorChan <- err
				return
			}

			fileText := strings.ToLower(string(content))
			wordsSlice := strings.Fields(fileText)

			fileSeen := make(map[string]bool)

			for _, word := range wordsSlice {
				for _, searchWord := range wordsRequest {
					searchWord = strings.ToLower(searchWord)

					if word == searchWord {
						tfMap[searchWord]++
						if !fileSeen[searchWord] {
							dfMap[searchWord]++
							fileSeen[searchWord] = true
						}
					}
				}
			}

			tfChan <- tfMap
			dfChan <- dfMap
		}(file)
	}

	go func() {
		wg.Wait() //viene eseguita quando il counter è zero
		close(tfChan)
		close(dfChan)
		close(errorChan)
	}()

	totalTF := make(map[string]int)
	totalDF := make(map[string]int)

	//i get the results from the channels
	//term frequency channel
	for tf := range tfChan {
		for word, count := range tf {
			totalTF[word] += count
		}
	}
	//document frequency channel
	for df := range dfChan {
		for word, count := range df {
			totalDF[word] += count
		}
	}

	if len(errorChan) > 0 {
		http.Error(w, "Error processing files", http.StatusInternalServerError)
		return
	}

	// Update word metadata and store results
	for _, word := range wordsRequest {
		tf := totalTF[word]
		df := totalDF[word]

		metadata, exists := wordStats[word]
		if !exists {
			// i init a new metadata word in the stats map
			wordStats[word] = &models.WordData{
				SearchCount: 0,
				History:     []models.TFDFEntry{},
			}
			metadata = wordStats[word]
		}
		metadata.SearchCount++
		metadata.LastTF = tf
		metadata.LastDF = df
		metadata.History = append(metadata.History, models.TFDFEntry{TF: tf, DF: df})

		result[word] = models.TFDFEntry{TF: tf, DF: df}
	}
	for w, stat := range wordStats {
		fmt.Printf("%v was searched %v times\nHistory: %v\n\n", w, stat.SearchCount, stat.History)
	}
	fmt.Println("----------------------------------------------------------------")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
