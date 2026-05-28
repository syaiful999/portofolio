package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *repository) GetEmployeeEducation(ctx context.Context, employeeId uuid.UUID) ([]EducationModel, error) {

	var items []EducationModel

	query := `SELECT 
				id, instance_name, major, degree, degree_id, gpa, graduate_year
			FROM master.v_master_employee_education
			WHERE is_active = TRUE AND employee_id = $1
			 ORDER BY id`

	rows, err := r.db.Query(query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result EducationModel
		if err := rows.Scan(
			&result.Id, &result.InstanceName, &result.Major, &result.Degree, &result.DegreeId, &result.GPA, &result.GraduateYear,
		); err != nil {
			return nil, err
		}
		items = append(items, result)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repository) InsertEmployeeMultipleEducations(tx *sql.Tx, userId, employeeId string, values []EducationModel) error {
	if len(values) == 0 {
		return nil
	}
	query := `
		INSERT INTO master.master_employee_education (
			modified_by, created_by, employee_id, 
			instance_name, major, degree_id, gpa, graduate_year
		) VALUES `
	params := []interface{}{}
	for i, v := range values {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8)

		params = append(params,
			userId,
			userId,
			employeeId,
			v.InstanceName,
			v.Major,
			v.DegreeId,
			v.GPA,
			v.GraduateYear,
		)
	}

	// Execute the query without RETURNING clause
	_, err := tx.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to InsertEmployeeMultipleEducations employee education: %v", err)
	}

	return nil
}

func (r *repository) InsertEmployeeEducation(tx *sql.Tx, employeeId string, value EducationModel) error {
	query := `INSERT INTO master.master_employee_education 
	(modified_by, created_by, employee_id, instance_name, major, degree_id, gpa, graduate_year)
			VALUES($1, $1, $2, 
			$3, $4, $5, $6, $7)`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		employeeId,
		value.InstanceName,
		value.Major,
		value.DegreeId,
		value.GPA,
		value.GraduateYear)

	if err != nil {
		return fmt.Errorf("failed to InsertEmployeeEducation employee education: %v", err)
	}

	return nil
}
func (r *repository) UpdateEmployeeEducation(tx *sql.Tx, value EducationModel) error {
	query := `UPDATE master.master_employee_education 
		SET modified_date=now(), 
		 modified_by =  $1,
			instance_name =  $2,
			major =  $3,
			degree_id =  $4,
			gpa =  $5,
			graduate_year =  $6
		WHERE id = $7`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.InstanceName,
		value.Major,
		value.DegreeId,
		value.GPA,
		value.GraduateYear,
		value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeeEducation employee education: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeEducation(tx *sql.Tx, value EducationModel) error {
	query := `UPDATE master.master_employee_education 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.Id)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeEducation employee education: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeeEducationByEmployeeId(tx *sql.Tx, value EducationModel) error {
	query := `UPDATE master.master_employee_education 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE employee_id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.EmployeeId)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeeEducationByEmployeeId employee education: %v", err)
	}

	return nil
}
