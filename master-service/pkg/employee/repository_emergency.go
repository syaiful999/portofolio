package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *repository) GetEmployeeEmergency(ctx context.Context, employeeId uuid.UUID) (EmergencyModel, error) {

	query := `SELECT
				id, emergency_name, relation_id, relation_name, phonenumber, 
				kelurahan_id, kelurahan_name,
				kecamatan_id, kecamatan_name, 
				city_id, city_name,
				province_id, province_name
			FROM master.v_master_employee_emergency
			WHERE is_active = TRUE AND employee_id = $1
			 ORDER BY id
			 LIMIT 1`

	var result EmergencyModel
	err := r.db.QueryRow(query, employeeId).Scan(
		&result.Id, &result.EmergencyName, &result.RelationId, &result.RelationName, &result.Phonenumber,
		&result.KelurahanId, &result.KelurahanName,
		&result.KecamatanId, &result.KecamatanName,
		&result.CityId, &result.CityName,
		&result.ProvinceId, &result.ProvinceName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return EmergencyModel{}, nil
		}
		return EmergencyModel{}, err
	}
	return result, nil
}

func (r *repository) InsertEmployeeEmergency(tx *sql.Tx, userId, employeeId string, value EmergencyModel) error {
	query := `
		INSERT INTO master.master_employee_emergency (
			modified_by, created_by, employee_id, emergency_name, relation_id, phonenumber, kelurahan_id

		) VALUES ($1, $1, $2, $3, $4, $5, $6);`

	_, err := tx.Exec(query,
		userId, employeeId, value.EmergencyName, value.RelationId,
		value.Phonenumber, value.KelurahanId)

	if err != nil {
		return fmt.Errorf("failed to insert employee emergency contact: %v", err)
	}

	return nil
}

func (r *repository) UpdateEmployeeEmergency(tx *sql.Tx, value EmergencyModel) error {
	query := `UPDATE master.master_employee_emergency 
		SET modified_date=now(), 
		 modified_by =  $1,
			emergency_name =  $2,
			relation_id =  $3,
			phonenumber =  $4,
			kelurahan_id =  $5
		WHERE id = $6;`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.EmergencyName,
		value.RelationId,
		value.Phonenumber,
		value.KelurahanId,
		value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeEmergency employee emergencies: %v", err)
	}

	return nil
}
func (r *repository) UpdateEmployeeEmergencyByEmployeeId(tx *sql.Tx, value EmergencyModel) error {
	query := `UPDATE master.master_employee_emergency 
		SET modified_date=now(), 
		 modified_by =  $1,
			emergency_name =  $2,
			relation_id =  $3,
			phonenumber =  $4,
			kelurahan_id =  $5
		WHERE employee_id = $6;`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.EmergencyName,
		value.RelationId,
		value.Phonenumber,
		value.KelurahanId,
		value.EmployeeId)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeEmergencyByEmployeeId employee emergencies: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeEmergency(tx *sql.Tx, value EmergencyModel) error {
	query := `UPDATE master.master_employee_emergency 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.Id)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeEmergency employee emergencies: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeEmergencyByEmployeeId(tx *sql.Tx, value EmergencyModel) error {
	query := `UPDATE master.master_employee_emergency 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE employee_id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.EmployeeId)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeEmergencyByEmployeeId employee emergency: %v", err)
	}

	return nil
}
