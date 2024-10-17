package database

import (
	"academyApplication/src/models"
	"database/sql"
)

type DAO interface {
	CreateStaff(staff models.Staff) (models.Staff, error)
	GetAllStaff() ([]models.Staff, error)
	UpdateStaff(staff models.Staff) (models.Staff, error)
	DeleteStaff(staffID int) error
	GetStaff(staffID int) (models.Staff, error)
	GetStaffByEmail(email string) (models.Staff, error)
	CreateRole(role models.Role) (models.Role, error)
	DeleteRole(roleID int) error
	UpdateRole(role models.Role) (models.Role, error)
	GetRole(roleID int) (models.Role, error)
	GetAllRoles() ([]models.Role, error)
	CheckIfRoleAlreadyExists(role models.Role) (bool, error)
	CheckIfStaffAlreadyExists(staff models.Staff) (bool, error)
	AssignRoleToStaff(staffId int, roleId int) (models.StaffWithRoles, error)
	CheckIfStaffRoleAlreadyExists(staffId int, roleId int) (bool, error)
	GetAllHospitals() ([]models.Hospital, error)
	DeleteHospital(id int) error
	UpdateHospital(hospital models.Hospital) (models.Hospital, error)
	AddHospital(hospital models.Hospital) (models.Hospital, error)
}

type DAOImpl struct {
	Db *sql.DB
}

// Staff section
func (dao *DAOImpl) CreateStaff(staff models.Staff) (models.Staff, error) {
	query := `insert into public."Staff" ("Email","Name","Password") values ($1, $2, $3) returning "ID"`
	err := dao.Db.QueryRow(query, staff.Email, staff.Name, staff.Password).Scan(&staff.ID)
	if err != nil {
		return models.Staff{}, err
	}
	return staff, nil
}

func (dao *DAOImpl) UpdateStaff(staff models.Staff) (models.Staff, error) {
	query := `UPDATE public."Staff" SET "Email" = $1, "Name" = $2, "Password" = $3 WHERE "ID" = $4`
	_, err := dao.Db.Exec(query, staff.Email, staff.Name, staff.Password, staff.ID)
	if err != nil {
		return models.Staff{}, err
	}
	return staff, nil
}

func (dao *DAOImpl) DeleteStaff(staffID int) error {
	query := `DELETE FROM public."Staff" WHERE "ID" = $1`
	_, err := dao.Db.Exec(query, staffID)
	return err
}

func (dao *DAOImpl) CheckIfStaffAlreadyExists(staff models.Staff) (bool, error) {
	query := `select exists (select 1 from public."Staff" where "Email" = $1)`
	var exists bool
	err := dao.Db.QueryRow(query, staff.Email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (dao *DAOImpl) GetStaff(staffID int) (models.Staff, error) {
	query := `SELECT "ID", "Email", "Name", "Password" FROM public."Staff" WHERE "ID" = $1`
	var staff models.Staff
	err := dao.Db.QueryRow(query, staffID).Scan(&staff.ID, &staff.Email, &staff.Name, &staff.Password)
	if err != nil {
		return models.Staff{}, err
	}
	return staff, nil
}

func (dao *DAOImpl) GetStaffByEmail(email string) (models.Staff, error) {
	query := `SELECT "ID", "Email", "Name", "Password" FROM public."Staff" WHERE "Email" = $1`
	var staff models.Staff
	err := dao.Db.QueryRow(query, email).Scan(&staff.ID, &staff.Email, &staff.Name, &staff.Password)
	if err != nil {
		return models.Staff{}, err
	}
	return staff, nil
}

func (dao *DAOImpl) GetAllStaff() ([]models.Staff, error) {
	query := `SELECT "ID", "Email", "Name", "Password" FROM public."Staff"`
	rows, err := dao.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staffList []models.Staff
	for rows.Next() {
		var staff models.Staff
		if err := rows.Scan(&staff.ID, &staff.Email, &staff.Name, &staff.Password); err != nil {
			return nil, err
		}
		staffList = append(staffList, staff)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return staffList, nil
}

//endstaff

// role section
func (dao *DAOImpl) CreateRole(role models.Role) (models.Role, error) {
	query := `insert into public."Role" ("Name") values ($1) returning "ID"`
	err := dao.Db.QueryRow(query, role.Name).Scan(&role.ID)
	if err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (dao *DAOImpl) DeleteRole(roleID int) error {
	query := `DELETE FROM public."Role" WHERE "ID" = $1`
	_, err := dao.Db.Exec(query, roleID)
	return err
}

func (dao *DAOImpl) UpdateRole(role models.Role) (models.Role, error) {
	query := `UPDATE public."Role" SET "Name" = $1 WHERE "ID" = $2`
	_, err := dao.Db.Exec(query, role.Name, role.ID)
	if err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (dao *DAOImpl) GetRole(roleID int) (models.Role, error) {
	query := `SELECT "ID", "Name" FROM public."Role" WHERE "ID" = $1`
	var role models.Role
	err := dao.Db.QueryRow(query, roleID).Scan(&role.ID, &role.Name)
	if err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (dao *DAOImpl) GetAllRoles() ([]models.Role, error) {
	query := `SELECT "ID", "Name" FROM public."Role"`
	rows, err := dao.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

func (dao *DAOImpl) CheckIfRoleAlreadyExists(role models.Role) (bool, error) {
	query := `select exists (select 1 from public."Role" where "Name" = $1)`
	var exists bool
	err := dao.Db.QueryRow(query, role.Name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

//endrole

// staffrole section
func (dao *DAOImpl) CheckIfStaffRoleAlreadyExists(staffID int, roleID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM public."StaffWithRoles" WHERE "staffID" = $1 AND "roleID" = $2)`
	var exists bool
	err := dao.Db.QueryRow(query, staffID, roleID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (dao *DAOImpl) AssignRoleToStaff(staffID int, roleID int) (models.StaffWithRoles, error) {
	queryInsert := `INSERT INTO public."StaffWithRoles" ("staffID", "roleID") VALUES ($1, $2)`
	_, err := dao.Db.Exec(queryInsert, staffID, roleID)
	if err != nil {
		return models.StaffWithRoles{}, err
	}
	// Recupera il record inserito
	querySelect := `
		SELECT s."ID", s."Email", r."ID", r."Name"
		FROM public."Staff" s
		JOIN public."StaffWithRoles" swr ON swr."staffID" = s."ID"
		JOIN public."Role" r ON swr."roleID" = r."ID"
		WHERE swr."staffID" = $1`

	rows, err := dao.Db.Query(querySelect, staffID)
	if err != nil {
		return models.StaffWithRoles{}, err // Gestisci eventuali errori
	}
	defer rows.Close()

	var staffWithRoles models.StaffWithRoles

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&staffWithRoles.Staff.ID, &staffWithRoles.Staff.Email, &role.ID, &role.Name); err != nil {
			return models.StaffWithRoles{}, err
		}
		staffWithRoles.Roles = append(staffWithRoles.Roles, role)
	}

	if err := rows.Err(); err != nil {
		return models.StaffWithRoles{}, err
	}

	return staffWithRoles, nil
}

func (dao *DAOImpl) GetStaffWithRolesByID(staffID int) (models.StaffWithRoles, error) {
	query := `
	SELECT s."ID", s."Email", r."ID", r."Name"
	FROM public."Staff" s
	JOIN public."StaffWithRoles" swr ON swr."staffID" = s."ID"
	JOIN public."Role" r ON swr."roleID" = r."ID"
	WHERE s."ID" = $1`

	rows, err := dao.Db.Query(query, staffID)
	if err != nil {
		return models.StaffWithRoles{}, err
	}
	defer rows.Close()

	var staffWithRoles models.StaffWithRoles

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&staffWithRoles.Staff.ID, &staffWithRoles.Staff.Email, &role.ID, &role.Name); err != nil {
			return models.StaffWithRoles{}, err
		}
		staffWithRoles.Roles = append(staffWithRoles.Roles, role)
	}

	if err := rows.Err(); err != nil {
		return models.StaffWithRoles{}, err
	}

	return staffWithRoles, nil
}

func (dao *DAOImpl) DeleteStaffRoleByStaffID(staffID int, roleID int) error {
	query := `
	DELETE FROM public."StaffWithRoles"
	WHERE "staffID" = $1 AND "roleID" = $2`
	_, err := dao.Db.Exec(query, staffID, roleID)
	return err
}

//end staff role

// hospital section
// GetAllHospitals recupera tutti gli ospedali
func (dao *DAOImpl) GetAllHospitals() ([]models.Hospital, error) {
	query := `SELECT "ID", "name", "address", "phone" FROM public."Hospital"`
	rows, err := dao.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hospitalList []models.Hospital
	for rows.Next() {
		var hospital models.Hospital
		if err := rows.Scan(&hospital.ID, &hospital.Name, &hospital.Address, &hospital.Phone); err != nil {
			return nil, err
		}
		hospitalList = append(hospitalList, hospital)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return hospitalList, nil
}

// DeleteHospital elimina un ospedale in base all'ID
func (dao *DAOImpl) DeleteHospital(id int) error {
	query := `DELETE FROM public."Hospital" WHERE "ID" = $1`
	_, err := dao.Db.Exec(query, id)
	return err
}

// UpdateHospital aggiorna un ospedale esistente
func (dao *DAOImpl) UpdateHospital(hospital models.Hospital) (models.Hospital, error) {
	query := `UPDATE public."Hospital" SET "Name" = $1, "Address" = $2, "Phone" = $3 WHERE "ID" = $4`
	_, err := dao.Db.Exec(query, hospital.Name, hospital.Address, hospital.Phone, hospital.ID)
	if err != nil {
		return models.Hospital{}, err // Restituisci un oggetto vuoto in caso di errore
	}
	return hospital, nil // Restituisci l'ospedale aggiornat
}

// Implementazione di DAOImpl
func (dao *DAOImpl) AddHospital(hospital models.Hospital) (models.Hospital, error) {
	query := `INSERT INTO public."Hospital" ("Name", "Address", "Phone") VALUES ($1, $2, $3) RETURNING "ID"`
	err := dao.Db.QueryRow(query, hospital.Name, hospital.Address, hospital.Phone).Scan(&hospital.ID)
	if err != nil {
		return models.Hospital{}, err // Restituisci un oggetto vuoto in caso di errore
	}
	return hospital, nil // Restituisci l'ospedale creato
}

//end hospital

//department section
