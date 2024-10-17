package models

import (
	"time"
)

type Staff struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type StaffResponse struct {
	ID    int    `json:"id"`
	Email string `json:"username"`
	Type  string `json:"type"`
	Token string `json:"token"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StaffWithRoles struct {
	Staff Staff  `json:"staff"`
	Roles []Role `json:"roles"`
}

type StaffShift struct {
	Staff  Staff   `json:"staff"`  // Dettagli del dipartimento
	Shifts []Shift `json:"shifts"` // Dettagli del turno
}

type Department struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	HospitalID int    `json:"hospital_id"`
}

// Struct per Shift
type Shift struct {
	ID        int       `json:"id"`
	StartTime time.Time `json:"start_time"` // Orario di inizio
	EndTime   time.Time `json:"end_time"`   // Orario di fine
}

type DepartmentShift struct {
	Department Department `json:"department"` // Dettagli del dipartimento
	Shift      Shift      `json:"shift"`      // Dettagli del turno
}

type Hospital struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type WordData struct {
	SearchCount int         `json:"search_count"`
	LastTF      int         `json:"last_tf"`
	LastDF      int         `json:"last_df"`
	History     []TFDFEntry `json:"history"`
}

type TFDFEntry struct {
	TF int `json:"tf"`
	DF int `json:"df"`
}
