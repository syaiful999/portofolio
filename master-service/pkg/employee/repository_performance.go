package employee

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *repository) GetEmployeePerformance(ctx context.Context, employeeId uuid.UUID) ([]PerformanceModel, error) {

	var items []PerformanceModel

	query := `SELECT 
				id, periode, predicate, score, description
			FROM transaction.transact_employee_performance
			WHERE is_active = TRUE AND employee_id = $1
			 ORDER BY id`

	rows, err := r.db.Query(query, employeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result PerformanceModel
		if err := rows.Scan(
			&result.Id, &result.Periode, &result.Predicate, &result.Score, &result.Description,
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

func (r *repository) InsertEmployeeMultiplePerformances(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []PerformanceModel) error {
	if len(values) == 0 {
		return nil
	}

	query := `
		INSERT INTO "transaction".transact_employee_performance (
			created_by, modified_by, employee_id, periode, predicate, score, description
		) VALUES `

	// Build the VALUES clause dynamically
	var params []interface{}
	paramIndex := 1

	for i, value := range values {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			paramIndex, paramIndex+1, paramIndex+2, paramIndex+3, paramIndex+4, paramIndex+5, paramIndex+6)

		params = append(params, userId, userId, employeeId, value.Periode, value.Predicate, value.Score, value.Description)
		paramIndex += 7
	}

	_, err := tx.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("failed to insert employee performances: %v", err)
	}

	return nil
}

func (r *repository) InsertEmployeePerformance(tx *sql.Tx, userId, employeeId string, value PerformanceModel) error {
	query := `
		INSERT INTO transaction.transact_employee_performance (
			modified_by, created_by, employee_id, 
			periode, predicate, score,
			description

		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`

	_, err := tx.Exec(query,
		userId, userId, employeeId,
		value.Periode, value.Predicate, value.Score,
		value.Description)

	if err != nil {
		return fmt.Errorf("failed to InsertEmployeePerformance contact: %v", err)
	}

	return nil
}

func (r *repository) UpdateEmployeePerformance(tx *sql.Tx, value PerformanceModel) error {
	query := `UPDATE transaction.transact_employee_performance 
		SET modified_date=now(), 
		 modified_by =  $1,
			periode =  $2,
			predicate =  $3,
			score =  $4,
			description =  $5
		WHERE id = $6`
	_, err := tx.Exec(query,
		value.ModifiedBy,
		value.Periode, value.Predicate, value.Score,
		value.Description,
		value.Id)

	if err != nil {
		return fmt.Errorf("failed to UpdateEmployeePerformance employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeePerformance(tx *sql.Tx, value PerformanceModel) error {
	query := `UPDATE transaction.transact_employee_performance 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.Id)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeePerformance employee dependents: %v", err)
	}

	return nil
}

func (r *repository) DeleteEmployeePerformanceByEmployeeId(tx *sql.Tx, value PerformanceModel) error {
	query := `UPDATE transaction.transact_employee_performance 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE employee_id = $2`
	_, err := tx.Exec(query, value.ModifiedBy, value.EmployeeId)
	if err != nil {
		return fmt.Errorf("failed to DeleteEmployeePerformanceByEmployeeId employee dependents: %v", err)
	}

	return nil
}
