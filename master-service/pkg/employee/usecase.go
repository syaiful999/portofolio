package employee

import (
	"context"
	"database/sql"
	"fmt"
	"moyo-master-service/config"
	"moyo-master-service/pkg/department"
	"moyo-master-service/pkg/departmentdetail"
	pb "moyo-master-service/pkg/employee/proto"
	"moyo-master-service/pkg/enum"
	"moyo-master-service/pkg/jobcompetency"
	"moyo-master-service/pkg/jobdescription"
	"moyo-master-service/pkg/jobfamily"
	"moyo-master-service/pkg/kelurahan"
	"moyo-master-service/pkg/outsource"
	"moyo-master-service/pkg/role"
	"moyo-master-service/pkg/section"
	"moyo-master-service/pkg/shifttemplate"
	"moyo-master-service/pkg/user"
	"moyo-master-service/utils"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IUseCase interface {
	GetEmployees(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeResponse, token *utils.TokenValue) error
	GetEmployeesAdvance(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeResponse, token *utils.TokenValue) error
	GetEmployeesAdvanceExport(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeExportResponse, token *utils.TokenValue) error
	GetEmployeeById(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetEmployeeByIDResponse) error
	GetTrainingRecommendation(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error
	GetTrainingOrganizerRecommendation(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error
	GetEmployeeEventSummary(ctx context.Context, req *pb.GetSummaryRequest, res *pb.GetEmployeeEventSummaryResponse, token *utils.TokenValue) error
	GetEmployeeEventList(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeEventResponse, token *utils.TokenValue) error
	GetEmployeeColumnList(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error
	GetEmployeeTemplateExport(ctx context.Context, req *pb.GetEmployeeTemplateRequest, res *pb.GetExportResponse) error
	GetEmployeeEventListExport(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeExportResponse) error
	GetEmployeeMinePermitSummary(ctx context.Context, req *pb.GetSummaryRequest, res *pb.GetEmployeeMinePermitSummaryResponse) error
	GetEmployeeMinePermitList(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeResponse) error
	GetEmployeeMinePermitListExport(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetExportResponse) error
	GetEmployeeJobDescriptionsCompetencies(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetEmployeeJobDescriptionResponse) error

	CreateEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error

	UpdateEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeBulk(ctx context.Context, req *pb.EmployeeUploadRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeProfile(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeAddress(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeAccount(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeContact(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeePerformance(ctx context.Context, req *pb.MutationEmployeePerformanceRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeEmployment(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeWorkExperience(ctx context.Context, req *pb.MutationEmployeeWorkExperienceRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeTrainingHistory(ctx context.Context, req *pb.MutationEmployeeTrainingHistoryRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeDependent(ctx context.Context, req *pb.MutationEmployeeDependentRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeEmergency(ctx context.Context, req *pb.EmergencyModel, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeReprimand(ctx context.Context, req *pb.MutationEmployeeReprimandRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeEducation(ctx context.Context, req *pb.MutationEmployeeEducationRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeDocument(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeHealthRecord(ctx context.Context, req *pb.MutationEmployeeHealthRecordRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	UpdateEmployeeJobDescription(ctx context.Context, req *pb.MutationEmployeeJobDescriptionRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error
	DeleteEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error

	DownloadExcelEmployee(req *pb.GetEmployeeByIDRequest, res *pb.GetFileResponse) error
}

type UseCase struct {
	repository                 IEmployeeRepository
	repositoryEnum             enum.IEnumRepository
	repositoryUser             user.IUserRepository
	repositoryRole             role.IRoleRepository
	repositoryDepartment       department.IDepartmentRepository
	repositoryDepartmentDetail departmentdetail.IDepartmentDetailRepository
	repositoryOutsource        outsource.IOutsourceRepository
	repositorySection          section.ISectionRepository
	repositoryJobFamily        jobfamily.IJobFamilyRepository
	repositoryJobDescription   jobdescription.IJobDescriptionRepository
	repositoryJobCompetency    jobcompetency.IJobCompetencyRepository
	repositoryShiftTemplate    shifttemplate.IShiftTemplateRepository
	repositoryKelurahan        kelurahan.IKelurahanRepository
	conf                       config.Config
	db                         *sql.DB
}

func NewUseCaseEmployee(
	repo IEmployeeRepository,
	repoEnum enum.IEnumRepository,
	repoUser user.IUserRepository,
	repoRole role.IRoleRepository,
	repoDepartment department.IDepartmentRepository,
	repoDepartmentDetail departmentdetail.IDepartmentDetailRepository,
	repoOutsource outsource.IOutsourceRepository,
	repoSection section.ISectionRepository,
	repoJobFamily jobfamily.IJobFamilyRepository,
	repoJobDescription jobdescription.IJobDescriptionRepository,
	repoJobCompetency jobcompetency.IJobCompetencyRepository,
	repoShiftTemplate shifttemplate.IShiftTemplateRepository,
	repoKelurahan kelurahan.IKelurahanRepository,
	conf config.Config,
	db *sql.DB,
) IUseCase {
	return &UseCase{
		repository:                 repo,
		repositoryEnum:             repoEnum,
		repositoryUser:             repoUser,
		repositoryRole:             repoRole,
		repositoryDepartment:       repoDepartment,
		repositoryDepartmentDetail: repoDepartmentDetail,
		repositoryOutsource:        repoOutsource,
		repositorySection:          repoSection,
		repositoryJobFamily:        repoJobFamily,
		repositoryJobDescription:   repoJobDescription,
		repositoryJobCompetency:    repoJobCompetency,
		repositoryShiftTemplate:    repoShiftTemplate,
		repositoryKelurahan:        repoKelurahan,
		conf:                       conf,
		db:                         db,
	}
}

func errorResponse1(res *pb.GetEmployeeResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task       // Only the high-level task/safe message is shown to users
	utils.PushLogf("", task, msg) // Technical msg is only logged
	return nil
}

func errorResponse2(res *pb.GetEmployeeByIDResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse3(res *pb.MutationEmployeeResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse4(res *pb.GetFileResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}
func errorResponse5(res *pb.RecommendationResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse6(res *pb.GetEmployeeEventSummaryResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse7(res *pb.GetEmployeeEventResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse8(res *pb.GetEmployeeExportResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse9(res *pb.GetEmployeeMinePermitSummaryResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse10(res *pb.RecommendationResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse11(res *pb.GetExportResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

func errorResponse12(res *pb.GetEmployeeJobDescriptionResponse, code int32, task, msg string) error {
	res.IsError = true
	res.ErrorCode = code
	res.ErrorMessage = task
	utils.PushLogf("", task, msg)
	return nil
}

var employeeDirectory = map[string][]string{
	"employment":      {"jabatan_id", "department_id", "department_detail_id", "section_id", "employement_type_id", "grade_id", "shift_id", "outsource_id", "contract_start_date"},
	"job-description": {"job_description_id"},
	"address":         {"kelurahan_id"},
	"profile":         {"blood_type_id", "gender_id", "birthdate"},
}

func (s *UseCase) GetEmployees(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeResponse, token *utils.TokenValue) error {
	//Handle error param take

	employees := make([]*pb.EmployeeModel, 0)

	if req.Take <= 0 {
		return nil
	}

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "Failed to fetch GetEmployees:kendo",
			err.Error(),
		)
	}

	// /* if user with role = 'pic_security' then filter by his department */

	dataRequest, err := s.repositoryUser.GetUserById(ctx, token.ID)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "Get employee failed due to token invalid",
			err.Error(),
		)
	}

	if dataRequest.RoleCode.String != utils.RoleSuperAdmin && dataRequest.RoleCode.String != utils.RolePicOutsourcing {
		if filter == "" {
			filter = fmt.Sprintf("WHERE department_id::TEXT = '%s'::TEXT", dataRequest.DepartmentId.String)
		} else {
			filter += fmt.Sprintf(" AND department_id::TEXT = '%s'::TEXT", dataRequest.DepartmentId.String)
		}
	} else if dataRequest.RoleCode.String == utils.RolePicOutsourcing {
		depts, err := s.repositoryUser.GetDepartementByOutsourceId(dataRequest.OutsourceId.String)
		if err != nil {
			return errorResponse1(
				res, http.StatusInternalServerError, "Get employee failed: outsource invalid",
				err.Error(),
			)
		}
		if len(depts) > 0 {
			deptList := strings.Join(depts, "','")
			if filter == "" {
				filter += fmt.Sprintf(" WHERE department_id::TEXT IN ('%s')", deptList)
			} else {
				filter += fmt.Sprintf(" AND department_id::TEXT IN ('%s')", deptList)
			}
		}

		if filter == "" {
			filter += fmt.Sprintf(" WHERE outsource_id = '%s'", dataRequest.OutsourceId.String)
		} else {
			filter += fmt.Sprintf(" AND outsource_id = '%s'", dataRequest.OutsourceId.String)
		}
	}

	if req.Take > 100 {
		req.Take = 100 // Cap results
	}

	data, count, err := s.repository.GetEmployee(ctx, req.Skip, req.Take, filter, sort)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "An unexpected error occurred while fetching employees",
			err.Error(),
		)
	}

	for _, value := range data {

		jobDescription := pb.JobDescriptionModel{
			JobDescription: value.JobDescriptionTitle.String,
		}
		var employee pb.EmployeeModel
		employee.CreatedDate = timestamppb.New(value.CreatedDate)
		employee.ModifiedBy = value.ModifiedBy
		employee.CreatedBy = value.CreatedBy
		employee.ModifiedDate = timestamppb.New(value.ModifiedDate)

		employee.Id = value.ID
		employee.EmployeeName = value.EmployeeName
		employee.Nik = value.Nik

		employee.DepartmentName = value.DepartmentName.String
		employee.DepartmentId = value.DepartmentId.String
		employee.DepartmentDetailName = value.DepartmentDetailName.String
		employee.DepartmentDetailId = value.DepartmentDetailId.String
		employee.SectionId = value.SectionId.String
		employee.SectionName = value.SectionName.String
		employee.OutsourceName = value.OutsourceName.String
		employee.OutsourceId = value.OutsourceId.String

		employee.Gender = value.Gender.String
		employee.Jabatan = value.Jabatan.String
		employee.Grade = value.Grade.String
		employee.JoinDate = utils.ConvTimeToString(value.JoinDate.Time)
		employee.BloodType = value.BloodType.String
		employee.Picture = value.Picture.String

		employee.DomisiliProvinceName = value.DomisiliProvinceName.String
		employee.DomisiliKecamatanName = value.DomisiliKecamatanName.String
		employee.DomisiliCityName = value.DomisiliCityName.String
		employee.DomisiliKelurahanName = value.DomisiliKelurahanName.String
		employee.JobDescription = &jobDescription

		employee.IsActive = value.IsActive
		employee.CreatedBy = value.CreatedBy
		employee.CreatedDate = timestamppb.New(value.CreatedDate)
		employee.ModifiedBy = value.ModifiedBy
		employee.ModifiedDate = timestamppb.New(value.ModifiedDate)
		employees = append(employees, &employee)
	}
	res.Data = employees
	res.CountData = count
	return nil
}

func (s *UseCase) GetEmployeesAdvance(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeResponse, token *utils.TokenValue) error {
	//Handle error param take

	employees := make([]*pb.EmployeeModel, 0)

	if req.Take <= 0 {
		return errorResponse1(
			res, http.StatusBadRequest,
			fmt.Sprintf("GetEmployeesAdvance:invalid parameter take, Req.Take = %d", req.Take),
			"",
		)
	}

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "Failed to fetch GetEmployeesAdvance:kendo",
			err.Error(),
		)
	}

	data, count, err := s.repository.GetEmployeesAdvanceSearch(ctx, req.Skip, req.Take, filter, sort)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "Failed to fetch GetEmployeesAdvance:",
			err.Error(),
		)
	}

	for _, value := range data {

		jobDescription := pb.JobDescriptionModel{
			JobDescription: value.JobDescriptionTitle.String,
		}

		var employee pb.EmployeeModel
		employee.CreatedDate = timestamppb.New(value.CreatedDate)
		employee.ModifiedBy = value.ModifiedBy
		employee.CreatedBy = value.CreatedBy
		employee.ModifiedDate = timestamppb.New(value.ModifiedDate)

		employee.Id = value.ID
		employee.EmployeeName = value.EmployeeName
		employee.Nik = value.Nik

		employee.Birthdate = utils.ConvTimeToStringDate(value.Birthdate)
		employee.Birthplace = value.Birthplace.String

		employee.DepartmentName = value.DepartmentName.String
		employee.DepartmentId = value.DepartmentId.String
		employee.DepartmentDetailName = value.DepartmentDetailName.String
		employee.DepartmentDetailId = value.DepartmentDetailId.String
		employee.SectionId = value.SectionId.String
		employee.SectionName = value.SectionName.String
		employee.OutsourceName = value.OutsourceName.String
		employee.OutsourceId = value.OutsourceId.String

		employee.Domisili = value.Domisili.String

		employee.ShiftName = value.ShiftName.String
		employee.Gender = value.Gender.String
		employee.Jabatan = value.Jabatan.String
		employee.Grade = value.Grade.String
		employee.JoinDate = utils.ConvTimeToStringDate(value.JoinDate.Time)
		employee.BloodType = value.BloodType.String

		employee.DomisiliProvinceName = value.DomisiliProvinceName.String
		employee.DomisiliKecamatanName = value.DomisiliKecamatanName.String
		employee.DomisiliCityName = value.DomisiliCityName.String
		employee.DomisiliKelurahanName = value.DomisiliKelurahanName.String
		employee.JobDescription = &jobDescription
		employee.StatusContract = value.StatusContract.String

		employee.IsActive = value.IsActive
		employee.CreatedBy = value.CreatedBy
		employee.CreatedDate = timestamppb.New(value.CreatedDate)
		employee.ModifiedBy = value.ModifiedBy
		employee.ModifiedDate = timestamppb.New(value.ModifiedDate)
		employees = append(employees, &employee)
	}
	res.Data = employees
	res.CountData = count
	return nil
}
func (s *UseCase) GetEmployeesAdvanceExport(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeExportResponse, token *utils.TokenValue) error {
	// 1. Define Master Directory
	var employeeDirectory = map[string][]string{
		"employment":      {"jabatan_id", "department_id", "department_detail_id", "section_id", "employement_type_id", "grade_id", "shift_id", "outsource_id", "contract_start_date"},
		"job-description": {"job_description_id"},
		"address":         {"kelurahan_id"},
		"profile":         {"blood_type_id", "gender_id", "birthdate"},
	}

	// 2. Define Column Structure

	// 3. Handle Sorting and Filtering
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse8(nil, http.StatusInternalServerError, "Kendo filter error", err.Error())
	}

	dataRequest, err := s.repositoryUser.GetUserById(ctx, token.ID)
	if err != nil {
		return errorResponse8(res, http.StatusInternalServerError, "User fetch error", err.Error())
	}

	// 4. Role-based Data Filtering Logic
	if dataRequest.RoleCode.String != utils.RoleSuperAdmin && dataRequest.RoleCode.String != utils.RolePicOutsourcing {
		whereClause := fmt.Sprintf("department_id::TEXT = '%s'::TEXT", dataRequest.DepartmentId.String)
		if filter == "" {
			filter = "WHERE " + whereClause
		} else {
			filter += " AND " + whereClause
		}
	} else if dataRequest.RoleCode.String == utils.RolePicOutsourcing {
		depts, err := s.repositoryUser.GetDepartementByOutsourceId(dataRequest.OutsourceId.String)
		if err == nil && len(depts) > 0 {
			deptList := strings.Join(depts, "','")
			whereDept := fmt.Sprintf("department_id::TEXT IN ('%s')", deptList)
			if filter == "" {
				filter = "WHERE " + whereDept
			} else {
				filter += " AND " + whereDept
			}
		}
		whereOutsource := fmt.Sprintf("outsource_id = '%s'", dataRequest.OutsourceId.String)
		if filter == "" {
			filter = "WHERE " + whereOutsource
		} else {
			filter += " AND " + whereOutsource
		}
	}

	// 5. Fetch Data
	data, _, err := s.repository.GetEmployeesAdvanceSearch(ctx, req.Skip, req.Take, filter, sort)
	if err != nil {
		return errorResponse8(nil, http.StatusInternalServerError, "Database fetch error", err.Error())
	}

	// 6. Column Access Logic
	roleMapping, err := s.repositoryRole.GetRoleMappingEmployee(dataRequest.RoleId.String)
	if err != nil {
		return errorResponse8(nil, http.StatusInternalServerError, "Role mapping error", err.Error())
	}

	allowedColumns := make(map[string]bool)
	for _, v := range roleMapping {
		if fields, ok := employeeDirectory[v.MenuCode]; ok {
			for _, field := range fields {
				allowedColumns[field] = true
			}
		}
	}

	// 7. Map all possible columns to their data logic
	allCols := []ColumnDef{
		{ID: "base_nik", Header: "NIK", Value: func(e *MasterEmployee) string { return e.Nik }},
		{ID: "base_name", Header: "Name", Value: func(e *MasterEmployee) string { return e.EmployeeName }},
		{ID: "gender_id", Header: "Jenis Kelamin", Value: func(e *MasterEmployee) string { return e.Gender.String }},
		{ID: "birthdate", Header: "Tanggal Lahir", Value: func(e *MasterEmployee) string { return utils.ConvTimeToStringDate(e.Birthdate) }},
		{ID: "department_id", Header: "Division", Value: func(e *MasterEmployee) string { return e.DepartmentName.String }},
		{ID: "department_detail_id", Header: "Department", Value: func(e *MasterEmployee) string { return e.DepartmentDetailName.String }},
		{ID: "section_id", Header: "Section", Value: func(e *MasterEmployee) string { return e.SectionName.String }},
		{ID: "employement_type_id", Header: "Employee Type", Value: func(e *MasterEmployee) string { return e.StatusContract.String }},
		{ID: "jabatan_id", Header: "Job Family", Value: func(e *MasterEmployee) string { return e.Jabatan.String }},
		{ID: "job_description_id", Header: "Job Description", Value: func(e *MasterEmployee) string { return e.JobDescriptionTitle.String }},
		{ID: "grade_id", Header: "Grade", Value: func(e *MasterEmployee) string { return e.Grade.String }},
		{ID: "shift_id", Header: "Shift", Value: func(e *MasterEmployee) string { return e.ShiftName.String }},
		{ID: "outsource_id", Header: "Service Contract", Value: func(e *MasterEmployee) string { return e.OutsourceName.String }},
		{ID: "contract_start_date", Header: "Date of Joining", Value: func(e *MasterEmployee) string { return utils.ConvTimeToStringDate(e.JoinDate.Time) }},
		{ID: "kelurahan_id", Header: "Domisili", Value: func(e *MasterEmployee) string {
			return fmt.Sprintf("Desa/Kel %s, %s, %s, %s", e.DomisiliKelurahanName.String, e.KecamatanName.String, e.DomisiliCityName.String, e.DomisiliProvinceName.String)
		}},
		{ID: "blood_type_id", Header: "Blood Type", Value: func(e *MasterEmployee) string { return e.BloodType.String }},
	}

	// Filter down to active columns
	var activeCols []ColumnDef
	for _, col := range allCols {
		if strings.HasPrefix(col.ID, "base_") || allowedColumns[col.ID] {
			activeCols = append(activeCols, col)
		}
	}

	// 8. Generate Excel
	f := excelize.NewFile()
	sheet := "Employees"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")
	getSheetActive, _ := f.GetSheetIndex(sheet)
	f.SetActiveSheet(getSheetActive)

	// Write Headers
	for i, colDef := range activeCols {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, colDef.Header)
	}

	// Write Data
	for i, emp := range data {
		rowIdx := i + 2
		for j, colDef := range activeCols {
			cell, _ := excelize.CoordinatesToCellName(j+1, rowIdx)
			f.SetCellValue(sheet, cell, colDef.Value(&emp))
		}
	}

	// Auto-width adjustment
	for i := range activeCols {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, colName, colName, 22)
	}

	// 9. Finalize Buffer and Response
	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorResponse8(nil, http.StatusInternalServerError, "Excel buffer error", err.Error())
	}

	res.Data = &pb.ExportResponse{
		FileName:    fmt.Sprintf("employees_advance_search_%s.xlsx", time.Now().Format("20060102_150405")),
		FileContent: buf.Bytes(),
		MimeType:    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	return nil
}
func (s *UseCase) GetTrainingRecommendation(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error {
	// GetTrainingRecommendation implementation
	reommendationData, err := s.repository.GetTrainingRecommendation(ctx)
	if err != nil {
		return errorResponse5(res, http.StatusInternalServerError, "Failed to fetch GetTrainingRecommendation:", err.Error())
	}
	res.Data = reommendationData
	return nil
}

func (s *UseCase) GetTrainingOrganizerRecommendation(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error {
	// GetTrainingOrganizerRecommendation implementation
	reommendationData, err := s.repository.GetTrainingOrganizerRecommendation(ctx)
	if err != nil {
		return errorResponse5(res, http.StatusInternalServerError, "Failed to fetch GetTrainingOrganizerRecommendation:", err.Error())
	}
	res.Data = reommendationData
	return nil
}

func (s *UseCase) GetEmployeeEventSummary(ctx context.Context, req *pb.GetSummaryRequest, res *pb.GetEmployeeEventSummaryResponse, token *utils.TokenValue) error {
	departmentId := utils.ConvStringToNullString(req.DepartmentId)
	outsourceId := utils.ConvStringToNullString(req.OutsourceId)
	filter := ""
	dataRequest, err := s.repositoryUser.GetUserById(ctx, token.ID)
	if err != nil {
		return errorResponse6(
			res, http.StatusInternalServerError, "Get employee failed due to token invalid",
			err.Error(),
		)
	}

	if dataRequest.RoleCode.String != utils.RoleSuperAdmin && dataRequest.RoleCode.String != utils.RolePicOutsourcing {

		filter += fmt.Sprintf(" AND department_id::TEXT = '%s'::TEXT", dataRequest.DepartmentId.String)
	} else if dataRequest.RoleCode.String == utils.RolePicOutsourcing {
		depts, err := s.repositoryUser.GetDepartementByOutsourceId(dataRequest.OutsourceId.String)
		if err != nil {
			return errorResponse6(
				res, http.StatusInternalServerError, "Get employee failed: outsource invalid",
				err.Error(),
			)
		}
		if len(depts) > 0 {
			deptList := strings.Join(depts, "','")
			filter += fmt.Sprintf(" AND department_id::TEXT IN ('%s')", deptList)
		}

		filter += fmt.Sprintf(" AND outsource_id = '%s'", dataRequest.OutsourceId.String)
	}

	result, err := s.repository.GetEmployeeEventSummary(ctx, filter, departmentId, outsourceId)
	if err != nil {
		return errorResponse6(res, http.StatusInternalServerError, "Failed to fetch GetEmployeeEventSummary:", err.Error())
	}
	dataResult := &pb.EmploeeEventSummaryModel{
		BirthdayToday:     result.BirthdayToday,
		BirthdayThisWeek:  result.BirthdayThisWeek,
		BirthdayThisMonth: result.BirthdayThisMonth,
		BirthdayNextMonth: result.BirthdayNextMonth,
	}
	res.Data = dataResult
	return nil
}

func (s *UseCase) GetEmployeeColumnList(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error {

	var columns []string
	filter := "WHERE enum_type ='employee_columns_bulk_update'"
	data, _, err := s.repositoryEnum.GetEnum(ctx, 0, 1, filter, "")
	if err != nil {
		return errorResponse10(
			res, http.StatusInternalServerError, "Failed to fetch GetEmployeeColumnList: GetEnums:",
			err.Error(),
		)
	}

	for _, value := range data {
		columns = strings.Split(value.EnumValue, "|")
	}

	res.Data = columns
	return nil
}

func (s *UseCase) GetEmployeeEventList(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeEventResponse, token *utils.TokenValue) error {
	dataResult := make([]*pb.EmployeeEventModel, 0)

	table := ""
	switch req.EventType {
	case "birthday_today":
		table = "v_master_employee_birthday_today"
	case "birthday_this_week":
		table = "v_master_employee_birthday_thisweek"
	case "birthday_this_month":
		table = "v_master_employee_birthday_thismonth"
	case "birthday_next_month":
		table = "v_master_employee_birthday_nextmonth"
	default:
		return errorResponse7(res, http.StatusBadRequest, "GetEmployeeEventSummary: invalid event type", "")
	}

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse7(
			res, http.StatusInternalServerError, "Failed to fetch GetEmployeeEventList:kendo",
			err.Error(),
		)
	}

	dataRequest, err := s.repositoryUser.GetUserById(ctx, token.ID)
	if err != nil {
		return errorResponse7(
			res, http.StatusInternalServerError, "Get employee failed due to token invalid",
			err.Error(),
		)
	}

	if dataRequest.RoleCode.String != utils.RoleSuperAdmin && dataRequest.RoleCode.String != utils.RolePicOutsourcing {
		if filter == "" {
			filter = fmt.Sprintf("WHERE department_id::TEXT = '%s'::TEXT", dataRequest.DepartmentId.String)
		} else {
			filter += fmt.Sprintf(" AND department_id::TEXT = '%s'::TEXT", dataRequest.DepartmentId.String)
		}
	} else if dataRequest.RoleCode.String == utils.RolePicOutsourcing {
		depts, err := s.repositoryUser.GetDepartementByOutsourceId(dataRequest.OutsourceId.String)
		if err != nil {
			return errorResponse7(
				res, http.StatusInternalServerError, "Get employee failed: outsource invalid",
				err.Error(),
			)
		}
		if len(depts) > 0 {
			deptList := strings.Join(depts, "','")
			if filter == "" {
				filter += fmt.Sprintf(" WHERE department_id::TEXT IN ('%s')", deptList)
			} else {
				filter += fmt.Sprintf(" AND department_id::TEXT IN ('%s')", deptList)
			}
		}

		if filter == "" {
			filter += fmt.Sprintf(" WHERE outsource_id = '%s'", dataRequest.OutsourceId.String)
		} else {
			filter += fmt.Sprintf(" AND outsource_id = '%s'", dataRequest.OutsourceId.String)
		}
	}
	results, count, err := s.repository.GetEmployeeEventList(ctx, table, filter, sort, req.Skip, req.Take)
	if err != nil {
		return errorResponse7(res, http.StatusInternalServerError, "Failed to fetch GetEmployeeEventList:", err.Error())
	}
	for _, v := range results {
		dataResult = append(dataResult, &pb.EmployeeEventModel{
			EmployeeName:   v.EmployeeName,
			Birthdate:      utils.ConvTimeToString(v.Birthdate),
			Nik:            v.NIK,
			DepartmentName: v.DepartmentName,
			OutsourceName:  v.OutsourceName.String,
			AgeOfEvent:     v.AgeOfEvent,
		})
	}

	res.CountData = count
	res.Data = dataResult
	return nil
}

func (s *UseCase) GetEmployeeEventListExport(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeExportResponse) error {

	table := ""
	switch req.EventType {
	case "birthday_today":
		table = "v_master_employee_birthday_today"
	case "birthday_this_week":
		table = "v_master_employee_birthday_thisweek"
	case "birthday_this_month":
		table = "v_master_employee_birthday_thismonth"
	case "birthday_next_month":
		table = "v_master_employee_birthday_nextmonth"
	default:
		return errorResponse8(res, http.StatusBadRequest, "failed to GetEmployeeEventListExport: invalid event type", "")
	}
	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse8(
			res, http.StatusInternalServerError, "failed to GetEmployeeEventListExport:kendo",
			err.Error(),
		)
	}

	results, _, err := s.repository.GetEmployeeEventList(ctx, table, filter, sort, 0, 999999999)
	if err != nil {
		return errorResponse8(res, http.StatusInternalServerError, "Failed to fetch GetEmployeeEventListExport:", err.Error())
	}

	// Create Excel file
	f := excelize.NewFile()
	sheet := "Employees"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")
	getSheetActive, _ := f.GetSheetIndex(sheet)
	f.SetActiveSheet(getSheetActive)

	// Header
	// nik, name, division, department, section, service contract
	headers := []string{"Nik", "Name", "Division", "Service Contract", "Date of Birth", "Age Of Event"}
	for i, h := range headers {
		col := string(rune('A' + i))
		f.SetCellValue(sheet, fmt.Sprintf("%s1", col), h)
	}

	for i, v := range results {
		row := i + 2

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), v.NIK)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.EmployeeName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.DepartmentName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.OutsourceName.String)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), utils.ConvTimeToStringDate(v.Birthdate))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), v.AgeOfEvent)
	}

	// Optional: adjust column widths
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheet, col, col, 20)
	}

	// Convert Excel file to buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorResponse8(
			nil, http.StatusInternalServerError, "Failed to write buffer in GetEmployeeEventListExport",
			err.Error(),
		)
	}

	// Set response values
	filename := fmt.Sprintf("employees_event_%s.xlsx", time.Now().Format("20060102_150405"))

	response := &pb.ExportResponse{
		FileName:    filename,
		FileContent: buf.Bytes(),
		MimeType:    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}
	res.Data = response

	return nil
}

func (s *UseCase) GetEmployeeById(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetEmployeeByIDResponse) error {
	var err error
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return errorResponse2(res, http.StatusBadRequest, "failed parse GetEmployeeById:uuid", err.Error())
	}
	employee, err := s.repository.GetEmployeeById(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "failed to fetch GetEmployeeById: "+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "GetEmployeeById:", err.Error())
	}

	currentShiftTemplate, err := s.repositoryShiftTemplate.GetCurrentShiftTemplateByEmployeeID(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "GetEmployeeById:GetCurrentShiftTemplateByEmployeeID"+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "GetEmployeeById:GetCurrentShiftTemplateByEmployeeID", err.Error())
	}
	workExperiences := make([]*pb.WorkExperienceModel, 0)
	workExperiencesModel, err := s.repository.GetEmployeeWorkExperience(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "failed to fetch GetEmployeeWorkExperience: "+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "failed to fetch GetEmployeeWorkExperience:", err.Error())
	}
	for _, we := range workExperiencesModel {
		workExperiences = append(workExperiences, &pb.WorkExperienceModel{
			Id:          we.Id,
			IsActive:    true,
			CompanyName: we.CompanyName.String,
			Position:    we.Position.String,
			JoinDate:    utils.ConvTimeToString(we.JoinDate.Time),
			EndDate:     utils.ConvTimeToString(we.EndDate.Time),
		})
	}

	performances := make([]*pb.PerformanceModel, 0)
	performancesModel, err := s.repository.GetEmployeePerformance(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "failed to GetEmployeePerformance: "+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "failed to GetEmployeePerformance:", err.Error())
	}
	for _, p := range performancesModel {
		performances = append(performances, &pb.PerformanceModel{
			Id:          p.Id,
			IsActive:    true,
			Periode:     p.Periode.String,
			Predicate:   p.Predicate.String,
			Score:       p.Score.String,
			Description: p.Description.String,
		})
	}

	trainings := make([]*pb.TrainingModel, 0)
	trainingsModel, err := s.repository.GetEmployeeTraining(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "trainings:"+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "trainings:", err.Error())
	}
	for _, t := range trainingsModel {
		trainingCompetencies, err := s.repository.GetEmployeeTrainingCompetency(ctx, t.Id)
		if err != nil {
			reg := regexp.MustCompile(utils.NoRowsInResultSet)
			if reg.Match([]byte(err.Error())) {
				errStr := fmt.Sprintf("employee id %s not found", ID)
				return errorResponse2(res, http.StatusNotFound, "trainings:"+utils.NoRowsInResultSet, errStr)
			}
			return errorResponse2(res, http.StatusInternalServerError, "trainings:", err.Error())
		}

		competencies := make([]*pb.JobCompetencyModel, 0)
		// trainings = append(trainings, &pb.JobCompetencyModel{
		for _, v := range trainingCompetencies {
			competencies = append(competencies, &pb.JobCompetencyModel{
				Id:                     v.Id,
				JobCompetencyId:        v.JobCompetencyId.String,
				CompetencyName:         v.CompetencyName.String,
				CompetencyCategoryName: v.CompetencyGroup.String,
				IsRequired:             v.IsRequired,
				IsActive:               v.IsActive,
				TrainingId:             v.TrainingId.String,
				Remark:                 v.Remark.String,
				OtherCompetency:        v.OtherCompetency.String,
			})

		}
		trainings = append(trainings, &pb.TrainingModel{
			Id:                t.Id,
			IsActive:          true,
			TrainingTitle:     t.TrainingTitle.String,
			TrainingDate:      utils.ConvTimeToString(t.TrainingDate.Time),
			OrganizerName:     t.OrganizerName.String,
			CertificateNumber: t.CertificateNumber.String,
			CompetenciesName:  t.CompetenciesName.String,
			TrainingMethod:    t.TrainingMethod.String,
			ExpiryDate:        utils.ConvTimeToString(t.ExpiryDate.Time),
			EffectiveDate:     utils.ConvTimeToString(t.EffectiveDate.Time),
			Cost:              t.Cost.String,
			CertificateUrl:    t.CertificateUrl.String,
			Competencies:      competencies,
		})
	}

	dependents := make([]*pb.DependentModel, 0)
	dependentsModel, err := s.repository.GetEmployeeDependent(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "failed to GetEmployeeDependent: "+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "failed to GetEmployeeDependent:", err.Error())
	}
	for _, d := range dependentsModel {
		dependents = append(dependents, &pb.DependentModel{
			Id:            d.Id,
			IsActive:      true,
			DependentName: d.DependentName.String,
			Birthdate:     utils.ConvTimeToString(d.Birthdate.Time),
			Phonenumber:   d.Phonenumber.String,
			RelationId:    d.RelationId.String,
			RelationName:  d.RelationName.String,
		})
	}

	emergenciesModel, err := s.repository.GetEmployeeEmergency(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "failed to GetEmployeeEmergency: "+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "failed to GetEmployeeEmergency:", err.Error())
	}
	emergency := &pb.EmergencyModel{
		Id:            emergenciesModel.Id,
		EmergencyName: emergenciesModel.EmergencyName.String,
		RelationId:    emergenciesModel.RelationId.String,
		RelationName:  emergenciesModel.RelationName.String,
		Phonenumber:   emergenciesModel.Phonenumber.String,
		KelurahanId:   emergenciesModel.KelurahanId.Int32,
		KelurahanName: emergenciesModel.KelurahanName.String,
		KecamatanId:   emergenciesModel.KecamatanId.String,
		KecamatanName: emergenciesModel.KecamatanName.String,
		CityId:        emergenciesModel.CityId.String,
		CityName:      emergenciesModel.CityName.String,
		ProvinceId:    emergenciesModel.ProvinceId.String,
		ProvinceName:  emergenciesModel.ProvinceName.String,
		IsActive:      true,
	}

	reprimands := make([]*pb.ReprimandModel, 0)
	reprimandsModel, err := s.repository.GetEmployeeReprimand(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "failed to GetEmployeeReprimand: "+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "failed to GetEmployeeReprimand:", err.Error())
	}
	for _, r := range reprimandsModel {
		reprimands = append(reprimands, &pb.ReprimandModel{
			Id:                r.Id,
			IsActive:          true,
			AttachmentUrl:     r.AttachmentUrl.String,
			WarningLevelId:    r.WarningLevelId.String,
			WarningLevelValue: r.WarningLevelValue.String,
			StartDate:         utils.ConvTimeToString(r.StartDate.Time),
			EndDate:           utils.ConvTimeToString(r.EndDate.Time),
			DocumentNumber:    r.DocumentNumber.String,
			Description:       r.Description.String,
		})
	}

	healthRecords := make([]*pb.HealthRecordModel, 0)
	healthRecordsModel, err := s.repository.GetEmployeeHealthRecord(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "failed to GetEmployeeHealthRecord: "+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "failed to GetEmployeeHealthRecord:", err.Error())
	}
	for _, hr := range healthRecordsModel {
		healthRecords = append(healthRecords, &pb.HealthRecordModel{
			Id:                hr.Id,
			IsActive:          true,
			McuDate:           utils.ConvTimeToString(hr.McuDate.Time),
			StatusId:          hr.StatusId.String,
			Status:            hr.Status.String,
			HealthDescription: hr.HealthDescription.String,
			McuUrl:            hr.McuUrl.String,
			McuFollowupUrl:    hr.McuFollowupUrl.String,
			StatusFollowupId:  hr.StatusFollowupId.String,
			StatusFollowup:    hr.StatusFollowup.String,
			Periode:           hr.Periode.String,
		})
	}

	educations := make([]*pb.EducationModel, 0)
	educationsModel, err := s.repository.GetEmployeeEducation(ctx, ID)
	if err != nil {
		reg := regexp.MustCompile(utils.NoRowsInResultSet)
		if reg.Match([]byte(err.Error())) {
			errStr := fmt.Sprintf("employee id %s not found", ID)
			return errorResponse2(res, http.StatusNotFound, "failed to GetEmployeeEducation: "+utils.NoRowsInResultSet, errStr)
		}
		return errorResponse2(res, http.StatusInternalServerError, "failed to GetEmployeeEducation:", err.Error())
	}
	for _, ed := range educationsModel {
		educations = append(educations, &pb.EducationModel{
			Id:           ed.Id,
			IsActive:     true,
			InstanceName: ed.InstanceName.String,
			Degree:       ed.Degree.String,
			DegreeId:     ed.DegreeId.String,
			Major:        ed.Major.String,
			Gpa:          ed.GPA,
			GraduateYear: ed.GraduateYear,
		})
	}

	var jobDescriptionModel jobdescription.MasterJobDescription
	competencies := make([]*pb.JobCompetencyModel, 0)
	if employee.JobDescriptionId.String != "" {
		jobDescriptionId, err := uuid.Parse(employee.JobDescriptionId.String)
		if err != nil {
			return errorResponse2(res, http.StatusBadRequest, "failed to GetEmployeeById:uuid", err.Error())
		} else {

			jobDescriptionModel, err = s.repositoryJobDescription.GetJobDescriptionById(ctx, jobDescriptionId)
			if err != nil {
				reg := regexp.MustCompile(utils.NoRowsInResultSet)
				if !reg.Match([]byte(err.Error())) {
					return errorResponse2(res, http.StatusInternalServerError, "employee: failed to GetJobDescriptionById:", err.Error())
				}
			}

			competenciesModel, err := s.repository.GetEmployeeCompetency(ctx, ID)
			if err != nil {
				reg := regexp.MustCompile(utils.NoRowsInResultSet)
				if reg.Match([]byte(err.Error())) {
					errStr := fmt.Sprintf("employee id %s not found", ID)
					return errorResponse2(res, http.StatusNotFound, "failed to GetEmployeeCompetency: "+utils.NoRowsInResultSet, errStr)
				}
				return errorResponse2(res, http.StatusInternalServerError, "failed to GetEmployeeCompetency:", err.Error())
			}
			for _, ed := range competenciesModel {
				competencies = append(competencies, &pb.JobCompetencyModel{
					Id:              ed.Id,
					JobCompetencyId: ed.JobCompetencyId.String,
					CompetencyName:  ed.CompetencyName.String,
					IsRequired:      ed.IsRequired,
					IsActive:        true,
					TotalTraining:   ed.TotalTraining,
				})
			}

		}
	}
	birthDate := utils.ConvTimeToString(employee.Birthdate)
	joinDate := utils.ConvTimeToString(employee.JoinDate.Time)

	res.Data = &pb.EmployeeModel{
		Id:           employee.ID,
		CreatedDate:  timestamppb.New(employee.CreatedDate),
		ModifiedBy:   employee.ModifiedBy,
		CreatedBy:    employee.CreatedBy,
		ModifiedDate: timestamppb.New(employee.ModifiedDate),
		IsActive:     employee.IsActive,

		EmployeeName: employee.EmployeeName,
		Nik:          employee.Nik,
		Npwp:         employee.Npwp.String,
		Noktp:        employee.NoKtp,

		Email:         employee.Email,
		EmailOffice:   employee.EmailOffice.String,
		Phonenumber:   employee.Phonenumber.String,
		Birthplace:    employee.Birthplace.String,
		Birthdate:     birthDate,
		Address:       employee.Address.String,
		KelurahanName: employee.KelurahanName.String,
		KelurahanId:   employee.KelurahanId.Int32,
		KecamatanName: employee.KecamatanName.String,
		KecamatanId:   employee.KecamatanId.String,
		CityName:      employee.CityName.String,
		CityId:        employee.CityId.String,
		ProvinceName:  employee.ProvinceName.String,
		ProvinceId:    employee.ProvinceId.String,

		DomisiliKelurahanId:   employee.DomisiliKelurahanId.Int32,
		DomisiliKelurahanName: employee.DomisiliKelurahanName.String,
		DomisiliKecamatanName: employee.DomisiliKecamatanName.String,
		DomisiliKecamatanId:   employee.DomisiliKecamatanId.String,
		DomisiliCityName:      employee.DomisiliCityName.String,
		DomisiliCityId:        employee.DomisiliCityId.String,
		DomisiliProvinceName:  employee.DomisiliProvinceName.String,
		DomisiliProvinceId:    employee.DomisiliProvinceId.String,

		Picture:          employee.Picture.String,
		StandardOvertime: employee.StandardOvertime.String,
		StandardWorkday:  employee.StandardWorkday.String,
		Certification:    employee.Certification.String,
		JoinDate:         joinDate,

		StatusContract:   employee.StatusContract.String,
		StatusContractId: employee.StatusContractId.String,
		StatusPtkp:       employee.StatusPtkp.String,
		StatusPtkpId:     employee.StatusPtkpId.String,
		Religion:         employee.Religion.String,
		ReligionId:       employee.ReligionId.String,

		GlNumber:         employee.GlNumber.String,
		GlNumberId:       employee.GlNumberId.String,
		StatusMarriage:   employee.StatusMarriage.String,
		StatusMarriageId: employee.StatusMarriageId.String,
		LastEducation:    employee.LastEducation.String,
		LastEducationId:  employee.LastEducationId.String,

		BpjstkNumber:  employee.BpjstkNumber.String,
		BpjskesNumber: employee.BpjskesNumber.String,

		BloodType:   employee.BloodType.String,
		BloodTypeId: employee.BloodTypeId.String,
		ShiftType:   employee.ShiftType.String,
		ShiftTypeId: employee.ShiftTypeId.String,
		Gender:      employee.Gender.String,
		GenderId:    employee.GenderId.String,

		Grade:        employee.Grade.String,
		GradeId:      employee.GradeId.String,
		Position:     employee.Position.String,
		Jabatan:      employee.Jabatan.String,
		JabatanId:    employee.JabatanId.String,
		CostCenter:   employee.CostCenter.String,
		CostCenterId: employee.CostCenterId.String,

		DepartmentName: employee.DepartmentName.String,
		DepartmentId:   employee.DepartmentId.String,
		OutsourceName:  employee.OutsourceName.String,
		OutsourceId:    employee.OutsourceId.String,

		DepartmentDetailName: employee.DepartmentDetailName.String,
		DepartmentDetailId:   employee.DepartmentDetailId.String,
		SectionName:          employee.SectionName.String,
		SectionId:            employee.SectionId.String,

		JobDescription: &pb.JobDescriptionModel{
			Id:                employee.JobDescriptionId.String,
			JobDescription:    jobDescriptionModel.JobDescription.String,
			JobResponsibility: jobDescriptionModel.JobResponsibility.String,
			JobSpecification:  jobDescriptionModel.JobSpecification.String,
			DepartmentId:      jobDescriptionModel.DepartmentId,
			DepartmentName:    jobDescriptionModel.DepartmentName,
			JobFamilyId:       jobDescriptionModel.JobFamilyId,
			JobFamilyName:     jobDescriptionModel.JobFamilyName,
			Competencies:      competencies,
		},

		RetirementDate: utils.ConvTimeToString(employee.RetirementDate.Time),
		MinePermitDate: utils.ConvTimeToString(employee.MinePermitDate.Time),
		ShiftId:        employee.ShiftId.String,
		ShiftName:      currentShiftTemplate,

		IjazahUrl:     employee.IjazahUrl.String,
		MinePermitUrl: employee.MinePermitUrl.String,
		KtpUrl:        employee.KtpUrl.String,
		NpwpUrl:       employee.NpwpUrl.String,
		BpjstkUrl:     employee.BpjstkUrl.String,
		BpjskesUrl:    employee.BpjskesUrl.String,
		KkUrl:         employee.KkUrl.String,

		ContractStartDate: utils.ConvTimeToString(employee.ContractStartDate.Time),
		ContractEndDate:   utils.ConvTimeToString(employee.ContractEndDate.Time),
		ContractRenewal:   employee.ContractRenewal.Int32,

		WorkExperiences: workExperiences,
		Performances:    performances,
		Trainings:       trainings,
		Dependents:      dependents,
		Emergency:       emergency,
		Reprimands:      reprimands,
		HealthRecords:   healthRecords,
		Educations:      educations,
	}

	return nil
}

func (s *UseCase) GetEmployeeMinePermitSummary(ctx context.Context, req *pb.GetSummaryRequest, res *pb.GetEmployeeMinePermitSummaryResponse) error {
	departmentId := utils.ConvStringToNullString(req.DepartmentId)
	outsourceId := utils.ConvStringToNullString(req.OutsourceId)
	result, err := s.repository.GetEmployeeMinePermitSummary(ctx, departmentId, outsourceId)
	if err != nil {
		return errorResponse9(res, http.StatusInternalServerError, "GetEmployeeMinePermitSummary:", err.Error())
	}
	dataResult := &pb.EmployeeMinePermitSummaryModel{
		MinePermitActive:          result.MinePermitActive,
		MinePermitExpired:         result.MinePermitExpired,
		MinePermitExpiry_180Days:  result.MinePermitExpiry180Days,
		MinePermitExpiry_90Days:   result.MinePermitExpiry90Days,
		MinePermitEndingThisMonth: result.MinePermitEndingThisMonth,
	}
	res.Data = dataResult
	return nil
}

func (s *UseCase) GetEmployeeMinePermitList(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeResponse) error {
	var filterMinePermit string
	switch req.EventType {
	case "mine_permit_active":
		filterMinePermit = "mine_permit_date > CURRENT_DATE"
	case "mine_permit_expired":
		filterMinePermit = "mine_permit_date <= CURRENT_DATE"
	case "mine_permit_expiry_180_days":
		filterMinePermit = "mine_permit_date >= CURRENT_DATE AND mine_permit_date <= (CURRENT_DATE + '180 days'::interval)"
	case "mine_permit_expiry_90_days":
		filterMinePermit = "mine_permit_date >= CURRENT_DATE AND mine_permit_date <= (CURRENT_DATE + '90 days'::interval)"
	case "mine_permit_ending_this_month":
		filterMinePermit = "to_char(mine_permit_date::timestamp with time zone, 'YYYYMM'::text) = to_char(CURRENT_DATE::timestamp with time zone, 'YYYYMM'::text)"
	default:
		return errorResponse1(res, http.StatusBadRequest, "GetEmployeeMinePermitList:", "invalid event type")
	}

	employees := make([]*pb.EmployeeModel, 0)

	if req.Take <= 0 {
		return errorResponse1(
			res, http.StatusBadRequest, "GetEmployees:",
			fmt.Sprintf("invalid parameter take, Req.Take = %d", req.Take),
		)
	}

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse1(
			res, http.StatusInternalServerError, "GetEmployees:kendo",
			err.Error(),
		)
	}

	if len(filter) > 0 {
		filter = fmt.Sprintf("%s AND (%s)", filter, filterMinePermit)
	} else {
		filter = fmt.Sprintf(" WHERE %s", filterMinePermit)
	}

	result, count, err := s.repository.GetEmployeeMinePermit(ctx, req.Skip, req.Take, filter, sort)
	if err != nil {
		return errorResponse1(res, http.StatusInternalServerError, "failed to GetEmployeeMinePermitList:", err.Error())
	}

	for _, v := range result {
		employees = append(employees, &pb.EmployeeModel{
			Nik:            v.Nik,
			EmployeeName:   v.EmployeeName,
			DepartmentName: v.DepartmentName.String,
			OutsourceName:  v.OutsourceName.String,
			MinePermitDate: utils.ConvTimeToStringDate(v.MinePermitDate.Time),
			Status:         v.Status,
		})
	}

	res.Data = employees
	res.CountData = count
	return nil
}

func (s *UseCase) GetEmployeeMinePermitListExport(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetExportResponse) error {
	var filterMinePermit string
	switch req.EventType {
	case "mine_permit_active":
		filterMinePermit = "mine_permit_date > CURRENT_DATE"
	case "mine_permit_expired":
		filterMinePermit = "mine_permit_date <= CURRENT_DATE"
	case "mine_permit_expiry_180_days":
		filterMinePermit = "mine_permit_date >= CURRENT_DATE AND mine_permit_date <= (CURRENT_DATE + '180 days'::interval)"
	case "mine_permit_expiry_90_days":
		filterMinePermit = "mine_permit_date >= CURRENT_DATE AND mine_permit_date <= (CURRENT_DATE + '90 days'::interval)"
	case "mine_permit_ending_this_month":
		filterMinePermit = "to_char(mine_permit_date::timestamp with time zone, 'YYYYMM'::text) = to_char(CURRENT_DATE::timestamp with time zone, 'YYYYMM'::text)"
	default:
		return errorResponse11(res, http.StatusBadRequest, "GetEmployeeMinePermitList:"+utils.InvalidEventType, "")
	}

	if req.Take <= 0 {
		return errorResponse11(
			res, http.StatusBadRequest, "GetEmployees:",
			fmt.Sprintf("invalid parameter take, Req.Take = %d", req.Take),
		)
	}

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter(req.Sort, req.Filter)
	if err != nil {
		return errorResponse11(
			res, http.StatusInternalServerError, "GetEmployees:kendo",
			err.Error(),
		)
	}

	if len(filter) > 0 {
		filter = fmt.Sprintf("%s AND (%s)", filter, filterMinePermit)
	} else {
		filter = fmt.Sprintf(" WHERE %s", filterMinePermit)
	}

	results, _, err := s.repository.GetEmployeeMinePermit(ctx, req.Skip, req.Take, filter, sort)
	if err != nil {
		return errorResponse11(res, http.StatusInternalServerError, "failed to GetEmployeeMinePermitListExport:", err.Error())
	}

	// Create Excel file
	f := excelize.NewFile()
	sheet := "Employees"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")
	getSheetActive, _ := f.GetSheetIndex(sheet)
	f.SetActiveSheet(getSheetActive)

	// Header
	headers := []string{"NIK", "Name", "Division", "Service Contract", "Valid Until", "Status"}
	for i, h := range headers {
		col := string(rune('A' + i))
		f.SetCellValue(sheet, fmt.Sprintf("%s1", col), h)
	}
	var nik string
	for i, v := range results {
		row := i + 2

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), v.Nik)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.EmployeeName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.DepartmentName.String)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.OutsourceName.String)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), utils.ConvTimeToStringDate(v.MinePermitDate.Time))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), v.Status)
		nik = v.Nik
	}

	// Optional: adjust column widths
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheet, col, col, 20)
	}

	// Convert Excel file to buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorResponse3(
			nil, http.StatusInternalServerError, "failed to write buffer",
			err.Error(),
		)
	}

	// Set response values
	filename := fmt.Sprintf("MOYO_MINE_PERMIT_%s_%s.xlsx", nik, time.Now().Format("20060102_150405"))

	response := &pb.ExportResponse{
		FileName:    filename,
		FileContent: buf.Bytes(),
		MimeType:    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}
	res.Data = response

	return nil
}

func (s *UseCase) GetEmployeeJobDescriptionsCompetencies(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetEmployeeJobDescriptionResponse) error {
	//Handle error param take

	jobDescriptions := make([]*pb.JobDescriptionModel, 0)
	var skip int32 = 0
	var take int32 = 99999999

	//generate condition where
	sort, filter, _, err := utils.KendoSortFilter("", "")
	if err != nil {
		return errorResponse12(
			res, http.StatusInternalServerError,
			"GetJobDescriptions:KendoSortFilter",
			err.Error(),
		)
	}
	// repoJobDescription
	dataCompetenciesTraining, err := s.repositoryJobCompetency.GetJobCompetencyTrainingEmployee(ctx, req.Id)
	if err != nil {
		return errorResponse12(
			res, http.StatusInternalServerError,
			"failed to GetJobCompetencyGroupByDescription",
			err.Error(),
		)
	}

	// repoJobCompetency
	data, count, err := s.repositoryJobDescription.GetJobDescription(ctx, skip, take, filter, sort)
	if err != nil {
		return errorResponse12(
			res, http.StatusInternalServerError,
			"GetJobDescriptions",
			err.Error(),
		)
	}

	for _, value := range data {

		dataCompetencies, err := s.repositoryJobCompetency.GetJobCompetencyGroupByDescription(ctx, value.ID)
		if err != nil {
			return errorResponse12(
				res, http.StatusInternalServerError,
				"failed to GetJobCompetencyGroupByDescription",
				err.Error(),
			)
		}

		// /*
		var jobCompetencies []*pb.JobCompetencyModel
		for _, v := range dataCompetencies {
			var totalTraining int32 = 0
			if _, exists := dataCompetenciesTraining[v.JobCompetencyID.String]; exists {
				totalTraining = dataCompetenciesTraining[v.JobCompetencyID.String]
				// if competency id used, then set it to 0; in order to calculate other training
				dataCompetenciesTraining[v.JobCompetencyID.String] = 0
			}

			jobCompetencies = append(jobCompetencies, &pb.JobCompetencyModel{
				Id:                     v.ID.String(),
				JobCompetencyId:        v.JobCompetencyID.String,
				CompetencyName:         v.CompetencyName.String,
				CompetencyCategoryName: v.CompetencyCategoryName,
				IsRequired:             v.IsRequired,
				TotalTraining:          totalTraining,
				IsActive:               true,
			})
		}
		// */
		var otherTraining int32
		for _, v := range dataCompetenciesTraining {
			otherTraining += v
		}
		var jobDescription pb.JobDescriptionModel

		jobDescription.Id = value.ID.String()
		jobDescription.JobDescription = value.JobDescription.String
		jobDescription.JobResponsibility = value.JobResponsibility.String
		jobDescription.JobSpecification = value.JobSpecification.String
		jobDescription.JobFamilyId = value.JobFamilyId
		jobDescription.JobFamilyName = value.JobFamilyName
		jobDescription.Competencies = jobCompetencies
		jobDescription.TotalOtherTraining = otherTraining
		jobDescriptions = append(jobDescriptions, &jobDescription)
	}

	res.Data = jobDescriptions
	res.CountData = count
	return nil
}

func (s *UseCase) CreateEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error
	var result MasterEmployee

	ID := uuid.New()

	//Handle mandatori field
	if utils.HandleMandatoryField(req.EmployeeName) {
		return errorResponse3(res, http.StatusBadRequest, "Employee Name is required", "")
	}

	birthDate, _ := utils.ConvStringToTime(req.Birthdate)
	joinDate, _ := utils.ConvStringToNullTime(req.JoinDate)
	contractStartDate, _ := utils.ConvStringToNullTime(req.ContractStartDate)
	contractEndDate, _ := utils.ConvStringToNullTime(req.ContractEndDate)
	minePermitDate, _ := utils.ConvStringToNullTime(req.MinePermitDate)

	// Input Validation
	if req.Email != "" && !utils.ValidateInput(req.Email, utils.RegexEmail) {
		return errorResponse3(res, http.StatusBadRequest, "Invalid email format", "Validation failed for email: "+req.Email)
	}
	if req.Phonenumber != "" && !utils.ValidateInput(req.Phonenumber, utils.RegexPhone) {
		return errorResponse3(res, http.StatusBadRequest, "Invalid phone number format", "Validation failed for phone: "+req.Phonenumber)
	}
	if req.Noktp != "" && !utils.ValidateInput(req.Noktp, utils.RegexNIK) {
		return errorResponse3(res, http.StatusBadRequest, "Invalid KTP format (must be 16 digits)", "Validation failed for KTP: "+req.Noktp)
	}
	data := MasterEmployee{
		ID:         ID.String(),
		IsActive:   true,
		CreatedBy:  token.ID.String(),
		ModifiedBy: token.ID.String(),

		EmployeeName: req.EmployeeName,
		Nik:          req.Nik,
		Npwp:         utils.ConvStringToNullString(req.Npwp),
		NoKtp:        req.Noktp,

		Email:               req.Email,
		EmailOffice:         utils.ConvStringToNullString(req.EmailOffice),
		Phonenumber:         utils.ConvStringToNullString(req.Phonenumber),
		Birthplace:          utils.ConvStringToNullString(req.Birthplace),
		Birthdate:           birthDate,
		Address:             utils.ConvStringToNullString(req.Address),
		KelurahanId:         utils.ConvInt32toNullInt32(req.KelurahanId),
		DomisiliKelurahanId: utils.ConvInt32toNullInt32(req.DomisiliKelurahanId),

		Picture:          utils.ConvStringToNullString(req.Picture),
		StandardOvertime: utils.ConvStringToNullString(req.StandardOvertime),
		StandardWorkday:  utils.ConvStringToNullString(req.StandardWorkday),
		Certification:    utils.ConvStringToNullString(req.Certification),
		JoinDate:         joinDate,

		StatusContractId: utils.ConvStringToNullString(req.StatusContractId),
		StatusPtkpId:     utils.ConvStringToNullString(req.StatusPtkpId),
		ReligionId:       utils.ConvStringToNullString(req.ReligionId),

		GlNumberId:       utils.ConvStringToNullString(req.GlNumberId),
		StatusMarriageId: utils.ConvStringToNullString(req.StatusMarriageId),
		LastEducationId:  utils.ConvStringToNullString(req.LastEducationId),

		BloodTypeId: utils.ConvStringToNullString(req.BloodTypeId),
		ShiftTypeId: utils.ConvStringToNullString(req.ShiftTypeId),
		GenderId:    utils.ConvStringToNullString(req.GenderId),

		GradeId:            utils.ConvStringToNullString(req.GradeId),
		JabatanId:          utils.ConvStringToNullString(req.JabatanId),
		Position:           utils.ConvStringToNullString(req.Position),
		CostCenterId:       utils.ConvStringToNullString(req.CostCenterId),
		DepartmentId:       utils.ConvStringToNullString(req.DepartmentId),
		OutsourceId:        utils.ConvStringToNullString(req.OutsourceId),
		DepartmentDetailId: utils.ConvStringToNullString(req.DepartmentDetailId),
		SectionId:          utils.ConvStringToNullString(req.SectionId),
		JobDescriptionId:   utils.ConvStringToNullString(req.JobDescriptionId),

		ContractStartDate: contractStartDate,
		ContractEndDate:   contractEndDate,
		MinePermitDate:    minePermitDate,
		BpjskesNumber:     utils.ConvStringToNullString(req.BpjskesNumber),
		BpjstkNumber:      utils.ConvStringToNullString(req.BpjstkNumber),
		KkUrl:             utils.ConvStringToNullString(req.KkUrl),
		KtpUrl:            utils.ConvStringToNullString(req.KtpUrl),
		NpwpUrl:           utils.ConvStringToNullString(req.NpwpUrl),
		BpjstkUrl:         utils.ConvStringToNullString(req.BpjstkUrl),
		IjazahUrl:         utils.ConvStringToNullString(req.IjazahUrl),
		BpjskesUrl:        utils.ConvStringToNullString(req.BpjskesUrl),
		MinePermitUrl:     utils.ConvStringToNullString(req.MinePermitUrl),
	}

	//check if exist
	countData, err := s.repository.GetCountEmployeeByNikOrEmail(ctx, data)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "Create employee failed", err.Error())
	}
	if countData > 0 {
		errStr := "Nik, no ktp, & email must be unique"
		return errorResponse3(res, http.StatusBadRequest, errStr, "")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "Create employee failed", err.Error())
	}

	if result, err = s.repository.CreateEmployee(tx, data); err != nil {
		return errorResponse3(res, http.StatusBadRequest, "Create employee failed", err.Error())
	}

	var employeeWorkExperiences []WorkExperienceModel
	for _, v := range req.WorkExperiences {
		joinDate, _ := utils.ConvStringToNullTime(v.JoinDate)
		endDate, _ := utils.ConvStringToNullTime(v.EndDate)
		employeeWorkExperiences = append(employeeWorkExperiences, WorkExperienceModel{
			CompanyName: utils.ConvStringToNullString(v.CompanyName),
			Position:    utils.ConvStringToNullString(v.Position),
			JoinDate:    joinDate,
			EndDate:     endDate,
		})
	}
	if err = s.repository.InsertEmployeeMultipleWorkExperiences(tx, ctx, data.ModifiedBy, result.ID, employeeWorkExperiences); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Create employee failed to InsertEmployeeMultipleWorkExperiences:", err.Error())
	}
	var employeePerformances []PerformanceModel
	for _, v := range req.Performances {

		employeePerformances = append(employeePerformances, PerformanceModel{
			Periode:     utils.ConvStringToNullString(v.Periode),
			Predicate:   utils.ConvStringToNullString(v.Predicate),
			Score:       utils.ConvStringToNullString(v.Score),
			Description: utils.ConvStringToNullString(v.Description),
		})
	}

	if err = s.repository.InsertEmployeeMultiplePerformances(tx, ctx, data.ModifiedBy, result.ID, employeePerformances); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Create employee failed to InsertEmployeeMultiplePerformances:", err.Error())
	}
	var employeeTraining []TrainingModel
	for _, v := range req.Trainings {
		trainingDate, _ := utils.ConvStringToNullTime(v.TrainingDate)
		expiryDate, _ := utils.ConvStringToNullTime(v.ExpiryDate)
		effectiveDate, _ := utils.ConvStringToNullTime(v.EffectiveDate)
		employeeTraining = append(employeeTraining, TrainingModel{
			TrainingTitle:     utils.ConvStringToNullString(v.TrainingTitle),
			TrainingDate:      trainingDate,
			OrganizerName:     utils.ConvStringToNullString(v.OrganizerName),
			CertificateNumber: utils.ConvStringToNullString(v.CertificateNumber),
			TrainingMethod:    utils.ConvStringToNullString(v.TrainingMethod),
			ExpiryDate:        expiryDate,
			EffectiveDate:     effectiveDate,
			Cost:              utils.ConvStringToNullString(v.Cost),
			CertificateUrl:    utils.ConvStringToNullString(v.CertificateUrl),
		})
	}

	if err = s.repository.InsertEmployeeMultipleTrainings(tx, ctx, data.ModifiedBy, result.ID, employeeTraining); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Create employee failed to InsertEmployeeMultipleTrainings:", err.Error())
	}
	var employeeDependent []DependentModel
	for _, v := range req.Dependents {
		birthdate, _ := utils.ConvStringToNullTime(v.Birthdate)
		employeeDependent = append(employeeDependent, DependentModel{
			DependentName: utils.ConvStringToNullString(v.DependentName),
			Birthdate:     birthdate,
			// Birthplace:    utils.ConvStringToNullString(v.Birthplace),
			Phonenumber:  utils.ConvStringToNullString(v.Phonenumber),
			RelationId:   utils.ConvStringToNullString(v.RelationId),
			RelationName: utils.ConvStringToNullString(v.RelationName),
		})
	}

	if err = s.repository.InsertEmployeeMultipleDependents(tx, data.ModifiedBy, result.ID, employeeDependent); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Create employee failed to InsertEmployeeMultipleDependents:", err.Error())
	}
	employeeEmergency := EmergencyModel{
		EmergencyName: utils.ConvStringToNullString(req.Emergency.EmergencyName),
		RelationId:    utils.ConvStringToNullString(req.Emergency.RelationId),
		RelationName:  utils.ConvStringToNullString(req.Emergency.RelationName),
		Phonenumber:   utils.ConvStringToNullString(req.Emergency.Phonenumber),
		KelurahanId:   utils.ConvInt32toNullInt32(req.Emergency.KelurahanId),
		KelurahanName: utils.ConvStringToNullString(req.Emergency.KelurahanName),
		KecamatanName: utils.ConvStringToNullString(req.Emergency.KecamatanName),
		CityName:      utils.ConvStringToNullString(req.Emergency.CityName),
		ProvinceName:  utils.ConvStringToNullString(req.Emergency.ProvinceName),
	}

	if err = s.repository.InsertEmployeeEmergency(tx, data.ModifiedBy, result.ID, employeeEmergency); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to InsertEmployeeEmergency:", err.Error())
	}

	var employeeHealthRecords []HealthRecordModel
	for _, v := range req.HealthRecords {
		mcuDate, _ := utils.ConvStringToNullTime(v.McuDate)

		employeeHealthRecords = append(employeeHealthRecords, HealthRecordModel{
			McuDate:           mcuDate,
			Periode:           utils.ConvStringToNullString(v.Periode),
			StatusId:          utils.ConvStringToNullString(v.StatusId),
			HealthDescription: utils.ConvStringToNullString(v.HealthDescription),
			McuUrl:            utils.ConvStringToNullString(v.McuUrl),
			McuFollowupUrl:    utils.ConvStringToNullString(v.McuFollowupUrl),
			StatusFollowupId:  utils.ConvStringToNullString(v.StatusFollowupId),
		})
	}

	if err = s.repository.InsertEmployeeMultipleHealthRecords(tx, ctx, data.ModifiedBy, result.ID, employeeHealthRecords); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to InsertEmployeeMultipleHealthRecords:", err.Error())
	}

	var employeeEducations []EducationModel
	for _, v := range req.Educations {

		employeeEducations = append(employeeEducations, EducationModel{
			InstanceName: utils.ConvStringToNullString(v.InstanceName),
			Major:        utils.ConvStringToNullString(v.Major),
			Degree:       utils.ConvStringToNullString(v.Degree),
			DegreeId:     utils.ConvStringToNullString(v.DegreeId),
			GPA:          v.Gpa,
			GraduateYear: v.GraduateYear,
		})
	}

	if err = s.repository.InsertEmployeeMultipleEducations(tx, data.ModifiedBy, result.ID, employeeEducations); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to InsertEmployeeMultipleEducations:", err.Error())
	}

	var employeeCompetencies []JobCompetencyModel
	for _, v := range req.Competencies {

		employeeCompetencies = append(employeeCompetencies, JobCompetencyModel{
			JobCompetencyId:        utils.ConvStringToNullString(v.JobCompetencyId),
			CompetencyName:         utils.ConvStringToNullString(v.CompetencyName),
			CompetencyCategoryName: utils.ConvStringToNullString(v.CompetencyCategoryName),
			JobDescriptionId:       utils.ConvStringToNullString(req.JobDescriptionId),
			IsRequired:             v.IsRequired,
		})
	}

	if err = s.repository.InsertEmployeeMultipleCompetencies(tx, data.ModifiedBy, result.ID, employeeCompetencies); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to InsertEmployeeMultipleEducations:", err.Error())
	}

	tx.Commit()
	res.Data = &pb.GetEmployeeByIDRequest{Id: result.ID}
	return nil
}

func (s *UseCase) DeleteEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error

	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed: invalid ID ", err.Error())
	}
	db := s.repository.GetDB()
	tx, err := db.Begin()

	if err = s.repository.DeleteEmployee(tx, ID, token.ID); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed", err.Error())
	}
	valueCompetency := JobCompetencyModel{
		EmployeeID: ID.String(),
		ModifiedBy: token.ID.String(),
	}
	if err = s.repository.DeleteEmployeeCompetencyByEmployeeId(tx, valueCompetency); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to DeleteEmployeeCompetencyByEmployeeId", err.Error())
	}

	valueDependent := DependentModel{
		EmployeeId: req.Id,
		ModifiedBy: token.ID.String(),
	}
	if err = s.repository.DeleteEmployeeDependentByEmployeeId(tx, valueDependent); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to DeleteEmployeeDependentByEmployeeId", err.Error())
	}

	valueEducation := EducationModel{
		EmployeeId: req.Id,
		ModifiedBy: token.ID.String(),
	}
	if err = s.repository.DeleteEmployeeEducationByEmployeeId(tx, valueEducation); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to DeleteEmployeeEducationByEmployeeId", err.Error())
	}

	valueEmergency := EmergencyModel{
		EmployeeId: req.Id,
		ModifiedBy: token.ID.String(),
	}
	if err = s.repository.DeleteEmployeeEmergencyByEmployeeId(tx, valueEmergency); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to DeleteEmployeeEmergency", err.Error())
	}

	valueHealthRecord := HealthRecordModel{
		EmployeeId: req.Id,
		ModifiedBy: token.ID.String(),
	}
	if err = s.repository.DeleteEmployeeHealthRecordByEmployeeId(tx, valueHealthRecord); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to DeleteEmployeeHealthRecord", err.Error())
	}

	valuePerformance := PerformanceModel{
		EmployeeId: req.Id,
		ModifiedBy: token.ID.String(),
	}
	if err = s.repository.DeleteEmployeePerformanceByEmployeeId(tx, valuePerformance); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to DeleteEmployeePerformanceByEmployeeId", err.Error())
	}

	valueReprimand := ReprimandModel{
		EmployeeId: req.Id,
		ModifiedBy: token.ID.String(),
	}
	if err = s.repository.DeleteEmployeeReprimandByEmployeeId(tx, valueReprimand); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to DeleteEmployeeReprimandByEmployeeId", err.Error())
	}

	valueTraining := TrainingModel{
		EmployeeId: req.Id,
		ModifiedBy: token.ID.String(),
	}
	if err = s.repository.DeleteEmployeeTrainingByEmployeeId(tx, valueTraining); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "Delete employee failed to DeleteEmployeeTrainingByEmployeeId", err.Error())
	}

	tx.Commit()

	res.Data = &pb.GetEmployeeByIDRequest{Id: req.Id}
	return nil
}
