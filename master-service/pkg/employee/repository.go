package employee

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type IEmployeeRepository interface {
	// returning data master_employee with specific needs
	// | skip is offset & take is limit in postgres
	// | filter is specific query you want to read
	// | sort is order by
	//
	// GetEmployee(ctx, 0, 10, "WHERE is_active=TRUE", "ORDER BY created_date")
	GetEmployee(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error)
	GetEmployeeExcelTemplate(ctx context.Context, departmentId, OutsourceId string) ([]MasterEmployee, error)
	GetEmployeesAdvanceSearch(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error)
	// returning data master_employee with specific id (master_employee.id)
	//
	// GetEmployeeById(ctx, 1)
	GetEmployeeById(ctx context.Context, id uuid.UUID) (MasterEmployee, error)
	GetEmployeeEventSummary(ctx context.Context, filter string, departmentId, outsourceId sql.NullString) (MasterEmployeeEventSummary, error)
	GetEmployeeEventList(ctx context.Context, table, filter, sort string, skip, take int32) ([]MasterEmployeeEvent, int32, error)

	GetEmployeeMinePermitSummary(ctx context.Context, departmentId, outsourceId sql.NullString) (MasterEmployeeMinePermitSummary, error)
	GetEmployeeMinePermit(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error)

	// create new data master_employee
	//
	// newMasterEmployee MasterEmployee{ ... }
	//
	// CreateEmployee(ctx, newMasterEmployee)
	CreateEmployee(tx *sql.Tx, value MasterEmployee) (MasterEmployee, error)
	// update data master_employee by id
	//
	// currentMasterEmployee MasterEmployee{ id, ... }
	//
	// UpdateEmployee(ctx, currentMasterEmployee)
	UpdateEmployee(ctx context.Context, value MasterEmployee) (MasterEmployee, error)
	UpdateBulkEmployees(ctx context.Context, tx *sql.Tx, list []EmployeeUpdateItem) error

	UpdateEmployeeProfile(ctx context.Context, value MasterEmployee) error
	UpdateEmployeeAddress(ctx context.Context, value MasterEmployee) error
	UpdateEmployeeAccount(ctx context.Context, value MasterEmployee) error
	UpdateEmployeeContact(ctx context.Context, value MasterEmployee) error
	UpdateEmployeeEmployment(ctx context.Context, value MasterEmployee) error
	UpdateEmployeeDocument(ctx context.Context, value MasterEmployee) error

	// delete data master_employee by id
	//
	// currentTransactEmployee MasterEmployee{ id, ... }
	//
	// DeleteMasterEmployee(ctx, uuid(123), uuid(123))
	DeleteEmployee(tx *sql.Tx, id, modifiedBy uuid.UUID) (err error)
	// returning data master_employee with specific id (transact_employee.id)
	//
	// GetCountEmployeeByNikOrEmail(ctx, "mymail@mail.com", "123")
	GetCountEmployeeByNikOrEmail(ctx context.Context, employee MasterEmployee) (int, error)
	GetCountEmployeeByNik(ctx context.Context, employee MasterEmployee) (int, error)
	GetEmployeeAll() ([]MasterEmployee, error)

	GetEmployeeWorkExperience(ctx context.Context, employeeId uuid.UUID) ([]WorkExperienceModel, error)
	GetEmployeePerformance(ctx context.Context, employeeId uuid.UUID) ([]PerformanceModel, error)
	GetEmployeeTraining(ctx context.Context, employeeId uuid.UUID) ([]TrainingModel, error)
	GetEmployeeTrainingCompetency(ctx context.Context, trainingId string) ([]TrainingCompetencyModel, error)

	GetEmployeeDependent(ctx context.Context, employeeId uuid.UUID) ([]DependentModel, error)
	GetEmployeeEmergency(ctx context.Context, employeeId uuid.UUID) (EmergencyModel, error)
	GetEmployeeReprimand(ctx context.Context, employeeId uuid.UUID) ([]ReprimandModel, error)
	GetEmployeeHealthRecord(ctx context.Context, employeeId uuid.UUID) ([]HealthRecordModel, error)
	GetEmployeeEducation(ctx context.Context, employeeId uuid.UUID) ([]EducationModel, error)

	InsertEmployeeMultipleWorkExperiences(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []WorkExperienceModel) error
	InsertEmployeeWorkExperience(tx *sql.Tx, userId, employeeId string, value WorkExperienceModel) error
	UpdateEmployeeWorkExperience(tx *sql.Tx, value WorkExperienceModel) error
	DeleteEmployeeWorkExperience(tx *sql.Tx, value WorkExperienceModel) error

	InsertEmployeeMultiplePerformances(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []PerformanceModel) error
	InsertEmployeePerformance(tx *sql.Tx, userId, employeeId string, value PerformanceModel) error
	UpdateEmployeePerformance(tx *sql.Tx, value PerformanceModel) error
	DeleteEmployeePerformance(tx *sql.Tx, value PerformanceModel) error
	DeleteEmployeePerformanceByEmployeeId(tx *sql.Tx, value PerformanceModel) error

	GetTrainingRecommendation(ctx context.Context) ([]string, error)
	GetTrainingOrganizerRecommendation(ctx context.Context) ([]string, error)
	InsertEmployeeMultipleTrainings(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []TrainingModel) error
	InsertEmployeeTraining(tx *sql.Tx, userId, employeeId string, value TrainingModel) (string, error)
	UpdateEmployeeTraining(tx *sql.Tx, value TrainingModel) error
	DeleteEmployeeTraining(tx *sql.Tx, value TrainingModel) error
	DeleteEmployeeTrainingByEmployeeId(tx *sql.Tx, value TrainingModel) error

	InsertEmployeeTrainingCompetency(tx *sql.Tx, userId, employeeId string, value TrainingCompetencyModel) error
	UpdateEmployeeTrainingCompetency(tx *sql.Tx, value TrainingCompetencyModel) error
	DeleteEmployeeTrainingCompetency(tx *sql.Tx, value TrainingCompetencyModel) error
	DeleteEmployeeTrainingCompetencyByTrainingId(tx *sql.Tx, value TrainingCompetencyModel) error

	InsertEmployeeMultipleDependents(tx *sql.Tx, userId, employeeId string, values []DependentModel) error
	InsertEmployeeDependent(tx *sql.Tx, employeeId string, value DependentModel) error
	UpdateEmployeeDependent(tx *sql.Tx, value DependentModel) error
	DeleteEmployeeDependent(tx *sql.Tx, value DependentModel) error
	DeleteEmployeeDependentByEmployeeId(tx *sql.Tx, value DependentModel) error

	InsertEmployeeEmergency(tx *sql.Tx, userId, employeeId string, value EmergencyModel) error
	UpdateEmployeeEmergency(tx *sql.Tx, value EmergencyModel) error
	UpdateEmployeeEmergencyByEmployeeId(tx *sql.Tx, value EmergencyModel) error
	DeleteEmployeeEmergency(tx *sql.Tx, value EmergencyModel) error
	DeleteEmployeeEmergencyByEmployeeId(tx *sql.Tx, value EmergencyModel) error

	InsertEmployeeReprimand(tx *sql.Tx, userId, employeeId string, value ReprimandModel) error
	UpdateEmployeeReprimand(tx *sql.Tx, value ReprimandModel) error
	DeleteEmployeeReprimand(tx *sql.Tx, value ReprimandModel) error
	DeleteEmployeeReprimandByEmployeeId(tx *sql.Tx, value ReprimandModel) error

	InsertEmployeeMultipleHealthRecords(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []HealthRecordModel) error
	InsertEmployeeHealthRecord(tx *sql.Tx, userId, employeeId string, value HealthRecordModel) error
	UpdateEmployeeHealthRecord(tx *sql.Tx, value HealthRecordModel) error
	DeleteEmployeeHealthRecord(tx *sql.Tx, value HealthRecordModel) error
	DeleteEmployeeHealthRecordByEmployeeId(tx *sql.Tx, value HealthRecordModel) error

	InsertEmployeeMultipleEducations(tx *sql.Tx, userId, employeeId string, values []EducationModel) error
	InsertEmployeeEducation(tx *sql.Tx, employeeId string, value EducationModel) error
	UpdateEmployeeEducation(tx *sql.Tx, values EducationModel) error
	DeleteEmployeeEducation(tx *sql.Tx, values EducationModel) error
	DeleteEmployeeEducationByEmployeeId(tx *sql.Tx, value EducationModel) error

	UpdateEmployeeJobDescription(tx *sql.Tx, value MasterEmployee) error

	GetEmployeeCompetency(ctx context.Context, employeeId uuid.UUID) ([]JobCompetencyModel, error)
	InsertEmployeeMultipleCompetencies(tx *sql.Tx, userId, employeeId string, values []JobCompetencyModel) error
	InsertEmployeeCompetency(tx *sql.Tx, employeeId string, value JobCompetencyModel) error
	UpdateEmployeeCompetency(tx *sql.Tx, value JobCompetencyModel) error
	DeleteEmployeeCompetencyByEmployeeId(tx *sql.Tx, value JobCompetencyModel) error

	GetDB() *sql.DB
}

type repository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *repository {

	return &repository{
		db: db,
	}
}

func (r *repository) GetEmployee(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error) {

	var items []MasterEmployee
	var countData int32

	query := fmt.Sprintf(`
	SELECT 
		id, employee_name, nik, 
		department_name, department_id, outsource_name, outsource_id,
		department_detail_name, department_detail_id, section_name, section_id,
		gender  , jabatan , grade , join_date , blood_type , picture,
		domisili_province_name , domisili_kecamatan_name , domisili_city_name , domisili_kelurahan_name , job_description,
		is_active, created_by, created_date, modified_by, modified_date 
	FROM master.v_master_employee %s %s limit $1 offset $2;`,
		filter, sort)

	rows, err := r.db.Query(query, take, skip)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var result MasterEmployee
		if err := rows.Scan(
			&result.ID, &result.EmployeeName, &result.Nik,
			&result.DepartmentName, &result.DepartmentId, &result.OutsourceName, &result.OutsourceId,
			&result.DepartmentDetailName, &result.DepartmentDetailId, &result.SectionName, &result.SectionId,
			&result.Gender, &result.Jabatan, &result.Grade, &result.JoinDate, &result.BloodType, &result.Picture,
			&result.DomisiliProvinceName, &result.DomisiliKecamatanName, &result.DomisiliCityName, &result.DomisiliKelurahanName, &result.JobDescriptionTitle,
			&result.IsActive, &result.CreatedBy, &result.CreatedDate, &result.ModifiedBy, &result.ModifiedDate,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, result)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, err
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	//count data
	queryCount := fmt.Sprintf(`select count(id) from master.v_master_employee %s;`,
		filter,
	)

	row := r.db.QueryRow(queryCount)
	errCount := row.Scan(&countData)
	if errCount != nil {
		return nil, 0, err
	}

	return items, countData, nil
}

func (r *repository) GetEmployeeExcelTemplate(ctx context.Context, departmentId, outsourceId string) ([]MasterEmployee, error) {

	var items []MasterEmployee

	var rows *sql.Rows
	var err error
	columns := `SELECT 
		employee_name, nik, npwp, email,
		department_name, department_detail_name , outsource_name,
		kelurahan_name, kecamatan_name, city_name, province_name,
		domisili_kelurahan_name, domisili_kecamatan_name, domisili_city_name, domisili_province_name, 
		gender, birthdate, birthplace, religion,
		status_marriage, status_ptkp, 
		gl_number, status_contract, last_education, blood_type, 
		cost_center , standard_workday, standard_overtime, job_description,
		jabatan, section_name, contract_start_date, contract_end_date,
		grade, join_date, mine_permit_date,
		bpjskes_number, bpjstk_number, email_office, phonenumber`
	query := ""
	if departmentId != "" {

		query = fmt.Sprintf(` %s FROM master.v_master_employee WHERE department_id=$1 `, columns)
		rows, err = r.db.Query(query, departmentId)

		if err != nil {
			return nil, err
		}

	} else if outsourceId != "" {

		query = fmt.Sprintf(` %s FROM master.v_master_employee WHERE outsource_id=$1 `, columns)
		rows, err = r.db.Query(query, outsourceId)

		if err != nil {
			return nil, err
		}

	}

	defer rows.Close()
	for rows.Next() {
		var result MasterEmployee
		if err := rows.Scan(
			&result.EmployeeName, &result.Nik, &result.Npwp, &result.Email, &result.DepartmentName, &result.DepartmentDetailName, &result.OutsourceName,
			&result.KelurahanName, &result.KecamatanName, &result.CityName, &result.ProvinceName,
			&result.DomisiliKelurahanName, &result.DomisiliKecamatanName, &result.DomisiliCityName, &result.DomisiliProvinceName,
			&result.Gender, &result.Birthdate, &result.Birthplace, &result.Religion,
			&result.StatusMarriage, &result.StatusPtkp,
			&result.GlNumber, &result.StatusContract, &result.LastEducation, &result.BloodType,
			&result.CostCenter, &result.StandardWorkday, &result.StandardOvertime, &result.JobDescriptionTitle,
			&result.Jabatan, &result.SectionName, &result.ContractStartDate, &result.ContractEndDate,
			&result.Grade, &result.JoinDate, &result.MinePermitDate,
			&result.BpjskesNumber, &result.BpjstkNumber, &result.EmailOffice, &result.Phonenumber,
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

func (r *repository) GetEmployeesAdvanceSearch(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error) {

	var items []MasterEmployee
	var countData int32

	query := fmt.Sprintf(`
	SELECT 
		id, employee_name, nik, domisili,
		department_name, department_id, outsource_name, outsource_id,
		department_detail_name, department_detail_id, section_name, section_id,
		gender, status_contract, jabatan, grade, join_date, blood_type, 
		birthdate, birthplace,
		domisili_province_name, domisili_kecamatan_name, domisili_city_name, domisili_kelurahan_name, job_description,
		shift_type,
		is_active, created_by, created_date, modified_by, modified_date 
	FROM master.v_master_employee_advance %s %s limit $1 offset $2;`,
		filter, sort)

	rows, err := r.db.Query(query, take, skip)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var result MasterEmployee
		if err := rows.Scan(
			&result.ID, &result.EmployeeName, &result.Nik, &result.Domisili,
			&result.DepartmentName, &result.DepartmentId, &result.OutsourceName, &result.OutsourceId,
			&result.DepartmentDetailName, &result.DepartmentDetailId, &result.SectionName, &result.SectionId,
			&result.Gender, &result.StatusContract, &result.Jabatan, &result.Grade, &result.JoinDate, &result.BloodType,
			&result.Birthdate, &result.Birthplace,
			&result.DomisiliProvinceName, &result.DomisiliKecamatanName, &result.DomisiliCityName, &result.DomisiliKelurahanName, &result.JobDescriptionTitle,
			&result.ShiftName,
			&result.IsActive, &result.CreatedBy, &result.CreatedDate, &result.ModifiedBy, &result.ModifiedDate,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, result)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, err
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	//count data
	queryCount := fmt.Sprintf(`select count(id) from master.v_master_employee_advance %s;`,
		filter,
	)

	row := r.db.QueryRow(queryCount)
	errCount := row.Scan(&countData)
	if errCount != nil {
		return nil, 0, err
	}

	return items, countData, nil
}

func (r *repository) GetEmployeeEventSummary(ctx context.Context, filter string, departmentId, outsourceId sql.NullString) (MasterEmployeeEventSummary, error) {
	var result MasterEmployeeEventSummary
	query := fmt.Sprintf(`SELECT
    COUNT(*) FILTER (
        WHERE to_char(me.birthdate, 'MM-DD')
              = to_char(CURRENT_DATE, 'MM-DD')
    ) AS birthday_today,
	 COUNT(*) FILTER (
        WHERE to_char(me.birthdate, 'MM-DD') >=
              to_char(date_trunc('week', CURRENT_DATE), 'MM-DD')
          AND to_char(me.birthdate, 'MM-DD') <=
              to_char(date_trunc('week', CURRENT_DATE) + INTERVAL '6 days', 'MM-DD')
    ) AS birthday_thisweek,
	 COUNT(*) FILTER (
        WHERE EXTRACT(month FROM me.birthdate)
              = EXTRACT(month FROM CURRENT_DATE)
    ) AS birthday_thismonth,
	 COUNT(*) FILTER (
        WHERE EXTRACT(month FROM me.birthdate)
              = EXTRACT(month FROM CURRENT_DATE + INTERVAL '1 month')
    ) AS birthday_nextmonth
	FROM master.master_employee me
	WHERE
		me.is_active
		AND me.birthdate <> DATE '0001-01-01'
		AND ($1::uuid is NULL OR me.department_id = $1)
		AND ($2::uuid is NULL OR me.outsource_id = $2)
		%s`, filter)
	row := r.db.QueryRow(query, departmentId, outsourceId)
	err := row.Scan(&result.BirthdayToday, &result.BirthdayThisWeek, &result.BirthdayThisMonth, &result.BirthdayNextMonth)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *repository) GetEmployeeEventList(ctx context.Context, table, filter, sort string, skip, take int32) ([]MasterEmployeeEvent, int32, error) {
	var result []MasterEmployeeEvent

	var countData int32

	query := fmt.Sprintf(`SELECT 
	nik, employee_name , department_name, outsource_name, birthdate, age_of_event
	FROM master.%s
	%s %s
	limit $1 offset $2;`, table, filter, sort)
	rows, err := r.db.Query(query, take, skip)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var item MasterEmployeeEvent
		if err := rows.Scan(
			&item.NIK, &item.EmployeeName,
			&item.DepartmentName, &item.OutsourceName, &item.Birthdate, &item.AgeOfEvent,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, item)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, err
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// //count data
	queryCount := fmt.Sprintf(`SELECT count(*) FROM master.%s;`,
		table,
	)

	row := r.db.QueryRow(queryCount)
	errCount := row.Scan(&countData)
	if errCount != nil {
		return nil, 0, err
	}

	return result, countData, nil
}

func (r *repository) GetEmployeeMinePermit(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error) {
	var result []MasterEmployee

	var countData int32

	query := fmt.Sprintf(`SELECT nik, employee_name, department_name, 
		outsource_name, mine_permit_date, status
		FROM master.v_master_employee_mine_permit %s %s limit $1 offset $2;`,
		filter, sort)

	rows, err := r.db.Query(query, take, skip)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var item MasterEmployee
		if err := rows.Scan(
			&item.Nik, &item.EmployeeName, &item.DepartmentName,
			&item.OutsourceName, &item.MinePermitDate, &item.Status,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, item)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, err
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	//count data
	queryCount := fmt.Sprintf(`SELECT count(*) FROM master.v_master_employee_mine_permit %s;`,
		filter,
	)

	row := r.db.QueryRow(queryCount)
	errCount := row.Scan(&countData)
	if errCount != nil {
		return nil, 0, err
	}

	return result, countData, nil
}

func (r *repository) GetEmployeeMinePermitSummary(ctx context.Context, departmentId, outsourceId sql.NullString) (MasterEmployeeMinePermitSummary, error) {
	var result MasterEmployeeMinePermitSummary

	row := r.db.QueryRow(`
	SELECT
		COUNT(*) FILTER (
			WHERE me.mine_permit_date > CURRENT_DATE
		) AS mine_permit_active,

		COUNT(*) FILTER (
			WHERE me.mine_permit_date <= CURRENT_DATE
		) AS mine_permit_expired,

		COUNT(*) FILTER (
			WHERE me.mine_permit_date BETWEEN CURRENT_DATE
				AND CURRENT_DATE + INTERVAL '180 days'
		) AS mine_permit_expiry_180_days,

		COUNT(*) FILTER (
			WHERE me.mine_permit_date BETWEEN CURRENT_DATE
				AND CURRENT_DATE + INTERVAL '90 days'
		) AS mine_permit_expiry_90_days,

		COUNT(*) FILTER (
			WHERE date_trunc('month', me.mine_permit_date)
				= date_trunc('month', CURRENT_DATE)
		) AS mine_permit_ending_this_month
	FROM master.master_employee me
	WHERE me.is_active 
	AND ($1::uuid is NULL OR me.department_id = $1)
    AND ($2::uuid is NULL OR me.outsource_id  = $2)`, departmentId, outsourceId)
	err := row.Scan(&result.MinePermitActive, &result.MinePermitExpired, &result.MinePermitExpiry180Days, &result.MinePermitExpiry90Days, &result.MinePermitEndingThisMonth)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *repository) GetEmployeeById(ctx context.Context, id uuid.UUID) (MasterEmployee, error) {

	var result MasterEmployee

	query := `
		SELECT 
			id, employee_name, nik, noktp, npwp, email, phonenumber, birthplace, birthdate, 
			address, kelurahan_name, kelurahan_id,
			kecamatan_name, kecamatan_id, city_name, city_id, province_name, province_id, 
			domisili_kelurahan_id, domisili_kelurahan_name, domisili_kecamatan_name, domisili_kecamatan_id,
			 domisili_city_name, domisili_city_id, domisili_province_name, domisili_province_id,
			picture, standard_overtime, standard_workday, certification, join_date, 
			status_contract, status_contract_id, status_ptkp, status_ptkp_id, religion, religion_id, 
			gl_number, gl_number_id, status_marriage, status_marriage_id, last_education, last_education_id,
			blood_type, blood_type_id, shift_type, shift_type_id, gender, gender_id, 
			grade, grade_id, "position", jabatan, jabatan_id, cost_center, cost_center_id, 
			department_name, department_id, outsource_name, outsource_id, 
			department_detail_name, department_detail_id, section_name, section_id,
			job_description_id, retirement_date, bpjstk_number, bpjskes_number, 
			email_office, mine_permit_date,
			ijazah_url, mine_permit_url, ktp_url, npwp_url, bpjstk_url, bpjskes_url, kk_url,
			contract_start_date, contract_end_date, contract_renewal,
			is_active, created_by, created_date, modified_by, modified_date 
		FROM master.v_master_employee
		WHERE id=$1`

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&result.ID, &result.EmployeeName, &result.Nik, &result.NoKtp, &result.Npwp, &result.Email, &result.Phonenumber, &result.Birthplace, &result.Birthdate,
		&result.Address, &result.KelurahanName, &result.KelurahanId,
		&result.KecamatanName, &result.KecamatanId, &result.CityName, &result.CityId, &result.ProvinceName, &result.ProvinceId,
		&result.DomisiliKelurahanId, &result.DomisiliKelurahanName, &result.DomisiliKecamatanName, &result.DomisiliKecamatanId,
		&result.DomisiliCityName, &result.DomisiliCityId, &result.DomisiliProvinceName, &result.DomisiliProvinceId,
		&result.Picture, &result.StandardOvertime, &result.StandardWorkday, &result.Certification, &result.JoinDate,
		&result.StatusContract, &result.StatusContractId, &result.StatusPtkp, &result.StatusPtkpId, &result.Religion, &result.ReligionId,
		&result.GlNumber, &result.GlNumberId, &result.StatusMarriage, &result.StatusMarriageId, &result.LastEducation, &result.LastEducationId,
		&result.BloodType, &result.BloodTypeId, &result.ShiftType, &result.ShiftTypeId, &result.Gender, &result.GenderId,
		&result.Grade, &result.GradeId, &result.Position, &result.Jabatan, &result.JabatanId, &result.CostCenter, &result.CostCenterId,
		&result.DepartmentName, &result.DepartmentId, &result.OutsourceName, &result.OutsourceId,
		&result.DepartmentDetailName, &result.DepartmentDetailId, &result.SectionName, &result.SectionId,
		&result.JobDescriptionId, &result.RetirementDate, &result.BpjstkNumber, &result.BpjskesNumber,
		&result.EmailOffice, &result.MinePermitDate,
		&result.IjazahUrl, &result.MinePermitUrl, &result.KtpUrl, &result.NpwpUrl, &result.BpjstkUrl, &result.BpjskesUrl, &result.KkUrl,
		&result.ContractStartDate, &result.ContractEndDate, &result.ContractRenewal,
		&result.IsActive, &result.CreatedBy, &result.CreatedDate, &result.ModifiedBy, &result.ModifiedDate,
	)
	return result, err
}

func (r *repository) CreateEmployee(tx *sql.Tx, value MasterEmployee) (MasterEmployee, error) {

	var result MasterEmployee
	result = value

	mutation := `
		INSERT INTO master.master_employee 
			(employee_name, nik, noktp, npwp, email, email_office, phonenumber, 
			birthplace, birthdate, address, kelurahan_id, domisili_kelurahan_id,
			picture, standard_overtime, standard_workday, certification, join_date,
			status_contract_id, status_ptkp_id, religion_id,
			gl_number_id,  status_marriage_id, last_education_id, 
			blood_type_id, shift_type_id, gender_id, 
			grade_id, jabatan_id, cost_center_id, department_id, outsource_id,
			section_id, department_detail_id, job_description_id,
			bpjstk_number, bpjskes_number,
			contract_start_date, contract_end_date, mine_permit_date,
			position, created_by, modified_by,
			kk_url, ktp_url, npwp_url, bpjstk_url,
			ijazah_url, bpjskes_url, mine_permit_url
			)
		VALUES($1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15,$16, $17,
			$18, $19, $20,
			$21, $22, $23,
			$24, $25, $26,
			$27, $28, $29, $30, $31,
			$32,$33, $34,
			$35, $36,
			$37, $38, $39,
			$40, $41, $42,
			$43, $44, $45, $46,
			$47, $48, $49
			)
		RETURNING id, created_date, modified_date`

	row := tx.QueryRow(mutation,
		value.EmployeeName, value.Nik, value.NoKtp, value.Npwp, value.Email, value.EmailOffice, value.Phonenumber.String,
		value.Birthplace.String, value.Birthdate, value.Address.String, value.KelurahanId, value.DomisiliKelurahanId,
		value.Picture, value.StandardOvertime.String, value.StandardWorkday.String, value.Certification.String, value.JoinDate,
		value.StatusContractId, value.StatusPtkpId, value.ReligionId,
		value.GlNumberId, value.StatusMarriageId, value.LastEducationId,
		value.BloodTypeId, value.ShiftTypeId, value.GenderId,
		value.GradeId, value.JabatanId, value.CostCenterId, value.DepartmentId, value.OutsourceId,
		value.SectionId, value.DepartmentDetailId, value.JobDescriptionId,
		value.BpjstkNumber, value.BpjskesNumber,
		value.ContractStartDate, value.ContractEndDate, value.MinePermitDate,
		value.Position, value.CreatedBy, value.ModifiedBy,
		value.KkUrl, value.KtpUrl, value.NpwpUrl, value.BpjstkUrl,
		value.IjazahUrl, value.BpjskesUrl, value.MinePermitUrl,
	)

	err := row.Scan(
		&result.ID,
		&result.CreatedDate,
		&result.ModifiedDate,
	)
	return result, err

}

func (r *repository) UpdateEmployee(ctx context.Context, value MasterEmployee) (MasterEmployee, error) {

	var result MasterEmployee
	result = value

	mutation := `UPDATE master.master_employee
				SET
					employee_name=$1,
					nik=$2,
					email=$3,
					phonenumber=$4,
					birthplace=$5,
					birthdate=$6,
					address=$7,
					kelurahan_id=$8,
					picture=$9,
					standard_overtime=$10,
					standard_workday=$11,
					certification=$12,
					join_date=$13,
					status_contract_id=$14,
					status_ptkp_id=$15,
					religion_id=$16,
					gl_number_id=$17,
					status_marriage_id=$18,
					last_education_id=$19,
					blood_type_id=$20,
					shift_type_id=$21,
					gender_id=$22,
					grade_id=$23,
					jabatan_id=$24,
					cost_center_id=$25,
					department_id=$26,
					outsource_id=$27,					
					position=$28,				
					noktp=$29,
					npwp=$30,
					modified_by=$31,
					modified_date=now()		
				WHERE
					id = $32
				RETURNING modified_date`
	_, err := r.db.Exec(mutation,
		value.EmployeeName,
		value.Nik,
		value.Email,
		value.Phonenumber,
		value.Birthplace,
		value.Birthdate,
		value.Address,
		value.KelurahanId,
		value.Picture,
		value.StandardOvertime,
		value.StandardWorkday,
		value.Certification,
		value.JoinDate,
		value.StatusContractId,
		value.StatusPtkpId,
		value.ReligionId,
		value.GlNumberId,
		value.StatusMarriageId,
		value.LastEducationId,
		value.BloodTypeId,
		value.ShiftTypeId,
		value.GenderId,
		value.GradeId,
		value.JabatanId,
		value.CostCenterId,
		value.DepartmentId,
		value.OutsourceId,
		value.Position,
		value.NoKtp,
		value.Npwp,
		value.ModifiedBy,
		value.ID,
	)

	return result, err
}

func (r *repository) UpdateBulkEmployees(ctx context.Context, tx *sql.Tx, list []EmployeeUpdateItem) error {

	for _, item := range list {
		if len(item.Fields) == 0 {
			continue
		}

		setParts := []string{}
		args := []interface{}{}
		argIdx := 1

		for key, val := range item.Fields {
			if unusedField(key) {
				continue
			}
			setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argIdx))
			args = append(args, val)
			argIdx++
		}

		// where nik
		args = append(args, item.Nik)

		query := fmt.Sprintf(`
            UPDATE master.master_employee
            SET %s, modified_date = NOW()
            WHERE nik = $%d
        `, strings.Join(setParts, ", "), argIdx)

		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

// this field can't be changed in via bulk update
func unusedField(keyText string) bool {
	return keyText == "domisili_kecamatan_id" ||
		keyText == "domisili_city_id" ||
		keyText == "domisili_province_id" ||
		keyText == "kecamatan_id" ||
		keyText == "city_id" ||
		keyText == "province_id" ||
		keyText == "nik"
}
func (r *repository) UpdateEmployeeProfile(ctx context.Context, value MasterEmployee) error {

	mutation := `UPDATE master.master_employee
				SET
					employee_name=$1,
					nik=$2,
					gender_id=$3,
					birthdate=$4,
					birthplace=$5,
					religion_id=$6,
					status_marriage_id=$7,
					picture=$8,
					blood_type_id=$9,
					npwp=$10,
					bpjstk_number=$11,
					bpjskes_number=$12,
					retirement_date=$13,
					status_ptkp_id=$14,
					noktp =$15,
					modified_by=$16,
					modified_date=now()		
				WHERE
					id = $17
				RETURNING modified_date`
	_, err := r.db.Exec(mutation,
		value.EmployeeName,
		value.Nik,
		value.GenderId,
		value.Birthdate,
		value.Birthplace,
		value.ReligionId,
		value.StatusMarriageId,
		value.Picture,
		value.BloodTypeId,
		value.Npwp,
		value.BpjstkNumber,
		value.BpjskesNumber,
		value.RetirementDate,
		value.StatusPtkpId,
		value.NoKtp,
		value.ModifiedBy,
		value.ID,
	)

	return err
}

func (r *repository) UpdateEmployeeAddress(ctx context.Context, value MasterEmployee) error {

	mutation := `UPDATE master.master_employee
				SET
					kelurahan_id=$1,
					domisili_kelurahan_id=$2,
					modified_by=$3,
					modified_date=now()		
				WHERE
					id = $4
				RETURNING modified_date`
	_, err := r.db.Exec(mutation,
		value.KelurahanId,
		value.DomisiliKelurahanId,
		value.ModifiedBy,
		value.ID,
	)

	return err
}
func (r *repository) UpdateEmployeeAccount(ctx context.Context, value MasterEmployee) error {

	mutation := `UPDATE master.master_employee
				SET
					cost_center_id=$1,
					gl_number_id=$2,
					standard_overtime=$3,
					standard_workday=$4,
					modified_by=$5,
					modified_date=now()		
				WHERE
					id = $6
				RETURNING modified_date`
	_, err := r.db.Exec(mutation,
		value.CostCenterId,
		value.GlNumberId,
		value.StandardOvertime,
		value.StandardWorkday,
		value.ModifiedBy,
		value.ID,
	)

	return err
}
func (r *repository) UpdateEmployeeContact(ctx context.Context, value MasterEmployee) error {

	mutation := `UPDATE master.master_employee
				SET
					email_office=$1,
					email=$2,
					phonenumber=$3,
					modified_by=$4,
					modified_date=now()		
				WHERE
					id = $5
				RETURNING modified_date`
	_, err := r.db.Exec(mutation,
		value.EmailOffice,
		value.Email,
		value.Phonenumber,
		value.ModifiedBy,
		value.ID,
	)

	return err
}
func (r *repository) UpdateEmployeeEmployment(ctx context.Context, value MasterEmployee) error {

	mutation := `UPDATE master.master_employee
				SET
					jabatan_id=$1,
					grade_id=$2,
					join_date=$3,
					mine_permit_date=$4,
					contract_start_date=$5,
					contract_end_date=$6,
					contract_renewal=$7,
					outsource_id=$8,
					section_id=$9,
					department_id=$10,
					department_detail_id=$11,
					status_contract_id=$12,
					modified_by=$13,
					modified_date=now()
				WHERE	id = $14
				RETURNING modified_date`
	_, err := r.db.Exec(mutation,
		value.JabatanId,
		value.GradeId,
		value.JoinDate,
		value.MinePermitDate,
		value.ContractStartDate,
		value.ContractEndDate,
		value.ContractRenewal,
		value.OutsourceId,
		value.SectionId,
		value.DepartmentId,
		value.DepartmentDetailId,
		value.StatusContractId,
		value.ModifiedBy,
		value.ID,
	)

	return err
}
func (r *repository) UpdateEmployeeDocument(ctx context.Context, value MasterEmployee) error {

	mutation := `UPDATE master.master_employee
				SET
					ijazah_url=$1,
					mine_permit_url=$2,
					ktp_url=$3,
					npwp_url=$4,
					bpjstk_url=$5,
					bpjskes_url=$6,
					kk_url=$7,
					modified_by=$8,
					modified_date=now()		
				WHERE
					id = $9
				RETURNING modified_date`
	_, err := r.db.Exec(mutation,
		value.IjazahUrl,
		value.MinePermitUrl,
		value.KtpUrl,
		value.NpwpUrl,
		value.BpjstkUrl,
		value.BpjskesUrl,
		value.KkUrl,
		value.ModifiedBy,
		value.ID,
	)

	return err
}

func (r *repository) DeleteEmployee(tx *sql.Tx, id, modifiedBy uuid.UUID) error {

	query := `UPDATE master.master_employee 
		SET modified_date=now(), 
			modified_by = $1, 
			is_active=FALSE 
		WHERE id = $2`
	_, err := tx.Exec(query, modifiedBy, id)
	return err
}

func (r *repository) GetCountEmployeeByNikOrEmail(ctx context.Context, employee MasterEmployee) (int, error) {

	var count int

	query := `SELECT count(*) from master.master_employee 
		WHERE (lower(nik) = lower($1) or lower(email) = lower($2) or lower(noktp) = lower($3))
			AND id!=$4
			AND is_active = TRUE
			AND email != ''`
	row := r.db.QueryRow(query, employee.Nik, employee.Email, employee.NoKtp, employee.ID)
	err := row.Scan(&count)

	return count, err
}

func (r *repository) GetCountEmployeeByNik(ctx context.Context, employee MasterEmployee) (int, error) {

	var count int
	query := `SELECT count(*) from master.master_employee 
		WHERE lower(nik) = lower($1)
			AND id!=$2
			AND is_active = TRUE
			AND email != ''`
	row := r.db.QueryRow(query, employee.Nik, employee.ID)
	err := row.Scan(&count)

	return count, err
}

func (r *repository) GetEmployeeAll() ([]MasterEmployee, error) {

	var items []MasterEmployee

	query := `SELECT 
				id, employee_name, nik, noktp, npwp, email, phonenumber, birthplace, birthdate, 
				address, kelurahan_name, kelurahan_id,
				picture, standard_overtime, standard_workday, certification, join_date, 
				status_contract, status_contract_id, status_ptkp, status_ptkp_id, religion, religion_id, 
				gl_number, gl_number_id, status_marriage, status_marriage_id, last_education, last_education_id,
				blood_type, blood_type_id, shift_type, shift_type_id, gender, gender_id, 
				grade, grade_id, "position", jabatan, jabatan_id, cost_center, cost_center_id, 
				department_name, department_id, outsource_name, outsource_id,
				is_active, created_by, created_date, modified_by, modified_date 
			FROM master.v_master_employee where is_active = true order by employee_name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result MasterEmployee
		if err := rows.Scan(
			&result.ID, &result.EmployeeName, &result.Nik, &result.NoKtp, &result.Npwp, &result.Email, &result.Phonenumber, &result.Birthplace, &result.Birthdate,
			&result.Address, &result.KelurahanName, &result.KelurahanId,
			&result.Picture, &result.StandardOvertime, &result.StandardWorkday, &result.Certification, &result.JoinDate,
			&result.StatusContract, &result.StatusContractId, &result.StatusPtkp, &result.StatusPtkpId, &result.Religion, &result.ReligionId,
			&result.GlNumber, &result.GlNumberId, &result.StatusMarriage, &result.StatusMarriageId, &result.LastEducation, &result.LastEducationId,
			&result.BloodType, &result.BloodTypeId, &result.ShiftType, &result.ShiftTypeId, &result.Gender, &result.GenderId,
			&result.Grade, &result.GradeId, &result.Position, &result.Jabatan, &result.JabatanId, &result.CostCenter, &result.CostCenterId,
			&result.DepartmentName, &result.DepartmentId, &result.OutsourceName, &result.OutsourceId,
			&result.IsActive, &result.CreatedBy, &result.CreatedDate, &result.ModifiedBy, &result.ModifiedDate,
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

func (r *repository) GetDB() *sql.DB {
	return r.db
}
