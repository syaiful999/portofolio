package employee

import (
	"database/sql"
	"time"
)

type MasterEmployee struct {
	ID                    string         `json:"id"`
	EmployeeName          string         `json:"employee_name"`
	Nik                   string         `json:"nik"`
	NoKtp                 string         `json:"noktp"`
	Npwp                  sql.NullString `json:"npwp"`
	Birthplace            sql.NullString `json:"birthplace"`
	Birthdate             time.Time      `json:"birthdate"`
	Email                 string         `json:"email"`
	Phonenumber           sql.NullString `json:"phonenumber"`
	Address               sql.NullString `json:"address"`
	KelurahanId           sql.NullInt32  `json:"kelurahan_id"`
	KelurahanName         sql.NullString `json:"kelurahan_name"`
	KecamatanName         sql.NullString `json:"kecamatan_name"`
	KecamatanId           sql.NullString `json:"kecamatan_id"`
	CityName              sql.NullString `json:"city_name"`
	CityId                sql.NullString `json:"city_id"`
	ProvinceName          sql.NullString `json:"province_name"`
	ProvinceId            sql.NullString `json:"province_id"`
	Domisili              sql.NullString `json:"domisili"`
	DomisiliKelurahanId   sql.NullInt32  `json:"domisili_kelurahan_id"`
	DomisiliKelurahanName sql.NullString `json:"domisili_kelurahan_name"`
	DomisiliKecamatanName sql.NullString `json:"domisili_kecamatan_name"`
	DomisiliKecamatanId   sql.NullString `json:"domisili_kecamatan_id"`
	DomisiliCityName      sql.NullString `json:"domisili_city_name"`
	DomisiliCityId        sql.NullString `json:"domisili_city_id"`
	DomisiliProvinceName  sql.NullString `json:"domisili_province_name"`
	DomisiliProvinceId    sql.NullString `json:"domisili_province_id"`
	Picture               sql.NullString `json:"picture"`
	DepartmentName        sql.NullString `json:"department_name"`
	DepartmentId          sql.NullString `json:"department_id"`
	DepartmentDetailName  sql.NullString `json:"department_detail_name"`
	DepartmentDetailId    sql.NullString `json:"department_detail_id"`
	SectionId             sql.NullString `json:"section_id"`
	SectionName           sql.NullString `json:"section_name"`
	OutsourceId           sql.NullString `json:"outsource_id"`
	OutsourceName         sql.NullString `json:"outsource_name"`
	StandardOvertime      sql.NullString `json:"standard_overtime"`
	StandardWorkday       sql.NullString `json:"standard_workday"`
	Certification         sql.NullString `json:"certification"`
	JoinDate              sql.NullTime   `json:"join_date"`
	StatusContract        sql.NullString `json:"status_contract"`
	StatusContractId      sql.NullString `json:"status_contract_id"`
	StatusPtkp            sql.NullString `json:"status_ptkp"`
	StatusPtkpId          sql.NullString `json:"status_ptkp_id"`
	Religion              sql.NullString `json:"religion"`
	ReligionId            sql.NullString `json:"religion_id"`
	GlNumber              sql.NullString `json:"gl_number"`
	GlNumberId            sql.NullString `json:"gl_number_id"`
	StatusMarriage        sql.NullString `json:"status_marriage"`
	StatusMarriageId      sql.NullString `json:"status_marriage_id"`
	LastEducation         sql.NullString `json:"last_education"`
	LastEducationId       sql.NullString `json:"last_education_id"`
	BloodType             sql.NullString `json:"blood_type"`
	BloodTypeId           sql.NullString `json:"blood_type_id"`
	ShiftType             sql.NullString `json:"shift_type"`
	ShiftTypeId           sql.NullString `json:"shift_type_id"`
	Gender                sql.NullString `json:"gender"`
	GenderId              sql.NullString `json:"gender_id"`
	Grade                 sql.NullString `json:"grade"`
	GradeId               sql.NullString `json:"grade_id"`
	Position              sql.NullString `json:"position"`
	Jabatan               sql.NullString `json:"jabatan"`
	JabatanId             sql.NullString `json:"jabatan_id"`
	CostCenter            sql.NullString `json:"cost_center"`
	CostCenterId          sql.NullString `json:"cost_center_id"`
	JobDescriptionId      sql.NullString `json:"job_description_id"`
	IsActive              bool           `json:"is_active"`
	CreatedDate           time.Time      `json:"created_date"`
	CreatedBy             string         `json:"created_by"`
	ModifiedDate          time.Time      `json:"modified_date"`
	ModifiedBy            string         `json:"modified_by"`
	Status                string         `json:"status"`

	RetirementDate      sql.NullTime          `json:"retirement_date"`
	BpjstkNumber        sql.NullString        `json:"bpjstk_number"`
	BpjskesNumber       sql.NullString        `json:"bpjskes_number"`
	EmailOffice         sql.NullString        `json:"email_office"`
	MinePermitDate      sql.NullTime          `json:"mine_permit_date"`
	ShiftId             sql.NullString        `json:"shift_id"`
	ShiftName           sql.NullString        `json:"shift_name"`
	IjazahUrl           sql.NullString        `json:"ijazah_url"`
	MinePermitUrl       sql.NullString        `json:"mine_permit_url"`
	KtpUrl              sql.NullString        `json:"ktp_url"`
	NpwpUrl             sql.NullString        `json:"npwp_url"`
	BpjstkUrl           sql.NullString        `json:"bpjstk_url"`
	BpjskesUrl          sql.NullString        `json:"bpjskes_url"`
	KkUrl               sql.NullString        `json:"kk_url"`
	ContractStartDate   sql.NullTime          `json:"contract_start_date"`
	ContractEndDate     sql.NullTime          `json:"contract_end_date"`
	ContractRenewal     sql.NullInt32         `json:"contract_renewal"`
	WorkExperiences     []WorkExperienceModel `json:"work_experiences"`
	Performances        []PerformanceModel    `json:"performances"`
	Trainings           []TrainingModel       `json:"trainings"`
	Dependents          []DependentModel      `json:"dependents"`
	Emergencies         []EmergencyModel      `json:"emergencies"`
	Reprimands          []ReprimandModel      `json:"reprimands"`
	HealthRecords       []HealthRecordModel   `json:"health_records"`
	Educations          []EducationModel      `json:"educations"`
	JobDescriptionTitle sql.NullString        `json:"job_description_title"`
	JobDescription      JobDescriptionModel   `json:"job_description"`
}
type MasterEmployeeEvent struct {
	ID             string         `json:"id"`
	NIK            string         `json:"nik"`
	EmployeeName   string         `json:"employee_name"`
	Birthdate      time.Time      `json:"birthdate"`
	DepartmentName string         `json:"department_name"`
	OutsourceName  sql.NullString `json:"outsource_name"`
	AgeOfEvent     int32          `json:"age_of_event"`
}

type MasterEmployeeEventSummary struct {
	BirthdayToday     int32 `json:"birthday_today"`
	BirthdayThisWeek  int32 `json:"birthday_this_week"`
	BirthdayThisMonth int32 `json:"birthday_this_month"`
	BirthdayNextMonth int32 `json:"birthday_next_month"`
}

type MasterEmployeeMinePermitSummary struct {
	MinePermitActive          int32 `json:"mine_permit_active"`
	MinePermitExpired         int32 `json:"mine_permit_expired"`
	MinePermitExpiry180Days   int32 `json:"mine_permit_expiry_180_days"`
	MinePermitExpiry90Days    int32 `json:"mine_permit_expiry_90_days"`
	MinePermitEndingThisMonth int32 `json:"mine_permit_ending_this_month"`
}

type WorkExperienceModel struct {
	Id          string         `json:"id"`
	ModifiedBy  string         `json:"modified_by"`
	CreatedBy   string         `json:"created_by"`
	CompanyName sql.NullString `json:"company_name"`
	Position    sql.NullString `json:"position"`
	JoinDate    sql.NullTime   `json:"join_date"`
	EndDate     sql.NullTime   `json:"end_date"`
}

type PerformanceModel struct {
	Id          string         `json:"id"`
	ModifiedBy  string         `json:"modified_by"`
	CreatedBy   string         `json:"created_by"`
	Periode     sql.NullString `json:"periode"`
	Predicate   sql.NullString `json:"predicate"`
	Score       sql.NullString `json:"score"`
	Description sql.NullString `json:"description"`
	EmployeeId  string         `json:"employee_id"`
}

type TrainingModel struct {
	Id                string         `json:"id"`
	ModifiedBy        string         `json:"modified_by"`
	CreatedBy         string         `json:"created_by"`
	TrainingTitle     sql.NullString `json:"training_title"`
	TrainingDate      sql.NullTime   `json:"training_date"`
	OrganizerName     sql.NullString `json:"organizer_name"`
	CertificateNumber sql.NullString `json:"certificate_number"`
	TrainingMethod    sql.NullString `json:"training_method"`
	ExpiryDate        sql.NullTime   `json:"expiry_date"`
	EffectiveDate     sql.NullTime   `json:"effective_date"`
	Cost              sql.NullString `json:"cost"`
	CertificateUrl    sql.NullString `json:"certificate_url"`
	CompetenciesName  sql.NullString `json:"competencies_name"`
	EmployeeId        string         `json:"employee_id"`
}

type TrainingCompetencyModel struct {
	Id              string         `json:"id"`
	ModifiedBy      string         `json:"modified_by"`
	CreatedBy       string         `json:"created_by"`
	TrainingId      sql.NullString `json:"training_id"`
	JobCompetencyId sql.NullString `json:"job_competency_id"`
	Remark          sql.NullString `json:"remark"`
	IsRequired      bool           `json:"is_required"`
	OtherCompetency sql.NullString `json:"other_competency"`
	CompetencyName  sql.NullString `json:"comptency_name"`
	CompetencyGroup sql.NullString `json:"comptency_group"`
	IsActive        bool           `json:"is_active"`
	EmployeeId      string         `json:"employee_id"`
}
type DependentModel struct {
	Id            string         `json:"id"`
	ModifiedBy    string         `json:"modified_by"`
	CreatedBy     string         `json:"created_by"`
	DependentName sql.NullString `json:"dependent_name"`
	Birthdate     sql.NullTime   `json:"birthdate"`
	Birthplace    sql.NullString `json:"birthplace"`
	Phonenumber   sql.NullString `json:"phonenumber"`
	RelationId    sql.NullString `json:"relation_id"`
	RelationName  sql.NullString `json:"relation_name"`
	EmployeeId    string         `json:"employee_id"`
}

type EmergencyModel struct {
	Id            string         `json:"id"`
	ModifiedBy    string         `json:"modified_by"`
	CreatedBy     string         `json:"created_by"`
	EmergencyName sql.NullString `json:"emergency_name"`
	RelationId    sql.NullString `json:"relation_id"`
	RelationName  sql.NullString `json:"relation_name"`
	Phonenumber   sql.NullString `json:"phonenumber"`
	KelurahanId   sql.NullInt32  `json:"kelurahan_id"`
	KelurahanName sql.NullString `json:"kelurahan_name"`
	KecamatanId   sql.NullString `json:"kecamatan_id"`
	KecamatanName sql.NullString `json:"kecamatan_name"`
	CityId        sql.NullString `json:"city_id"`
	CityName      sql.NullString `json:"city_name"`
	ProvinceId    sql.NullString `json:"province_id"`
	ProvinceName  sql.NullString `json:"province_name"`
	EmployeeId    string         `json:"employee_id"`
}

type ReprimandModel struct {
	Id                string         `json:"id"`
	ModifiedBy        string         `json:"modified_by"`
	CreatedBy         string         `json:"created_by"`
	AttachmentUrl     sql.NullString `json:"attachment_url"`
	WarningLevelId    sql.NullString `json:"waning_level_id"`
	WarningLevelValue sql.NullString `json:"warning_level_value"`
	ValidTime         sql.NullTime   `json:"valid_time"`
	StartDate         sql.NullTime   `json:"start_date"`
	EndDate           sql.NullTime   `json:"end_date"`
	DocumentNumber    sql.NullString `json:"document_number"`
	Description       sql.NullString `json:"description"`
	EmployeeId        string         `json:"employee_id"`
}

type HealthRecordModel struct {
	Id                string         `json:"id"`
	ModifiedBy        string         `json:"modified_by"`
	CreatedBy         string         `json:"created_by"`
	McuDate           sql.NullTime   `json:"mcu_date"`
	Periode           sql.NullString `json:"periode"`
	StatusId          sql.NullString `json:"status_id"`
	Status            sql.NullString `json:"status"`
	HealthDescription sql.NullString `json:"health_description"`
	McuUrl            sql.NullString `json:"mcu_url"`
	McuFollowupUrl    sql.NullString `json:"mcu_followup__url"`
	StatusFollowupId  sql.NullString `json:"status_followup_id"`
	StatusFollowup    sql.NullString `json:"status_followup"`
	EmployeeId        string         `json:"employee_id"`
}

type EducationModel struct {
	Id           string         `json:"id"`
	ModifiedBy   string         `json:"modified_by"`
	CreatedBy    string         `json:"created_by"`
	InstanceName sql.NullString `json:"instance_name"`
	Major        sql.NullString `json:"major"`
	Degree       sql.NullString `json:"degree"`
	DegreeId     sql.NullString `json:"degree_id"`
	GPA          float64        `json:"gpa"`
	GraduateYear int32          `json:"graduate_year"`
	EmployeeId   string         `json:"employee_id"`
}

type JobDescriptionModel struct {
	Id                string               `json:"id"`
	ModifiedBy        string               `json:"modified_by"`
	CreatedBy         string               `json:"created_by"`
	JobDescription    sql.NullString       `json:"job_description"`
	JobResponsibility sql.NullString       `json:"job_responsibility"`
	JobSpecification  sql.NullString       `json:"job_specification"`
	DepartmentId      sql.NullString       `json:"department_id"`
	DepartmentName    sql.NullString       `json:"department_name"`
	JobFamilyId       sql.NullString       `json:"job_family_id"`
	JobFamilyName     sql.NullString       `json:"job_family_name"`
	CompetenciesName  sql.NullString       `json:"competencies_name"`
	Competencies      []JobCompetencyModel `json:"competencies"`
}

type JobCompetencyModel struct {
	Id                     string         `json:"id"`
	ModifiedBy             string         `json:"modified_by"`
	EmployeeID             string         `json:"employee_id"`
	JobCompetencyId        sql.NullString `json:"job_competency_id"`
	JobDescriptionId       sql.NullString `json:"job_description_id"`
	JobDescription         sql.NullString `json:"job_description"`
	CompetencyName         sql.NullString `json:"competency_name"`
	CompetencyCategoryName sql.NullString `json:"competency_category_name"`
	IsRequired             bool           `json:"is_required"`
	TotalTraining          int32          `json:"total_training"`
	IsActive               bool           `json:"is_active"`
}

type Option struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type EmployeeUpdateItem struct {
	Nik    string
	Fields map[string]interface{} // e.g {"birthdate": "2024-01-01", "department_id": "uuid"}
}

type LookupMaps struct {
	Department       map[string]string
	DepartmentDetail map[string]string
	Outsource        map[string]string
	Jabatan          map[string]string
	Grade            map[string]string
	Gender           map[string]string
	StatusPTKP       map[string]string
	StatusContract   map[string]string
	CostCenter       map[string]string
	GLNumber         map[string]string
	JobDescription   map[string]string
	Section          map[string]string
	StatusMarriage   map[string]string
	BloodType        map[string]string
	Religion         map[string]string
	Kelurahan        map[string]string
	Kecamatan        map[string]string
	City             map[string]string
	Province         map[string]string
}

// 2. Define your Column Definitions
// We map the "ID" from your employeeDirectory to the Header Name and the Data Logic
type ColumnDef struct {
	ID     string
	Header string
	Value  func(emp *MasterEmployee) string // Adjust 'EmployeeEntity' to your actual struct type
}
