package employee

import (
	"context"
	"database/sql"
	"testing"

	"moyo-master-service/config"
	pb "moyo-master-service/pkg/employee/proto"
	"moyo-master-service/pkg/user"
	"moyo-master-service/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockEmployeeRepository is a mock implementation of IEmployeeRepository for testing.
type MockEmployeeRepository struct {
	MockGetEmployee                  func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error)
	MockGetEmployeeById              func(ctx context.Context, id uuid.UUID) (MasterEmployee, error)
	MockCreateEmployee               func(tx *sql.Tx, value MasterEmployee) (MasterEmployee, error)
	MockUpdateEmployee               func(ctx context.Context, value MasterEmployee) (MasterEmployee, error)
	MockDeleteEmployee               func(tx *sql.Tx, id, modifiedBy uuid.UUID) (err error)
	MockGetCountEmployeeByNikOrEmail func(ctx context.Context, employee MasterEmployee) (int, error)
	MockGetEmployeeAll               func() ([]MasterEmployee, error)
}

func (m *MockEmployeeRepository) GetEmployee(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error) {
	if m.MockGetEmployee != nil {
		return m.MockGetEmployee(ctx, skip, take, filter, sort)
	}
	return nil, 0, nil
}

func (m *MockEmployeeRepository) GetEmployeeById(ctx context.Context, id uuid.UUID) (MasterEmployee, error) {
	if m.MockGetEmployeeById != nil {
		return m.MockGetEmployeeById(ctx, id)
	}
	return MasterEmployee{}, nil
}

func (m *MockEmployeeRepository) CreateEmployee(tx *sql.Tx, value MasterEmployee) (MasterEmployee, error) {
	if m.MockCreateEmployee != nil {
		return m.MockCreateEmployee(tx, value)
	}
	return MasterEmployee{}, nil
}

func (m *MockEmployeeRepository) UpdateEmployee(ctx context.Context, value MasterEmployee) (MasterEmployee, error) {
	if m.MockUpdateEmployee != nil {
		return m.MockUpdateEmployee(ctx, value)
	}
	return MasterEmployee{}, nil
}

func (m *MockEmployeeRepository) DeleteEmployee(tx *sql.Tx, id, modifiedBy uuid.UUID) (err error) {
	if m.MockDeleteEmployee != nil {
		return m.MockDeleteEmployee(tx, id, modifiedBy)
	}
	return nil
}

func (m *MockEmployeeRepository) GetCountEmployeeByNikOrEmail(ctx context.Context, employee MasterEmployee) (int, error) {
	if m.MockGetCountEmployeeByNikOrEmail != nil {
		return m.MockGetCountEmployeeByNikOrEmail(ctx, employee)

	}
	return 0, nil
}

func (m *MockEmployeeRepository) GetEmployeeAll() ([]MasterEmployee, error) {
	if m.MockGetEmployeeAll != nil {
		return m.MockGetEmployeeAll()
	}
	return nil, nil
}

// minimal mock for IUserRepository
type mockUserRepo struct{}

func (m *mockUserRepo) GetUser(ctx context.Context, skip, take int32, filter, sort string) ([]user.MasterUser, int32, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) GetUserById(ctx context.Context, id uuid.UUID) (user.MasterUser, error) {
	// return a user with no special role
	return user.MasterUser{ID: id, RoleCode: sql.NullString{String: "", Valid: false}, DepartmentId: sql.NullString{String: "", Valid: false}}, nil
}
func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (user.MasterUser, string, error) {
	return user.MasterUser{}, "", nil
}
func (m *mockUserRepo) CreateUser(ctx context.Context, value user.MasterUser) (user.MasterUser, string, error) {
	return user.MasterUser{}, "", nil
}
func (m *mockUserRepo) UpdateUser(ctx context.Context, value user.MasterUser) (user.MasterUser, string, error) {
	return user.MasterUser{}, "", nil
}
func (m *mockUserRepo) UpdateUserPassword(ctx context.Context, value user.MasterUser, redirectId string) (user.MasterUser, error) {
	return user.MasterUser{}, nil
}
func (m *mockUserRepo) DeleteUser(ctx context.Context, id, modifiedBy uuid.UUID) error {
	return nil
}
func (m *mockUserRepo) GetCountUniqueUser(ctx context.Context, data user.MasterUser) (int, error) {
	return 0, nil
}
func (m *mockUserRepo) GetUserGroupbyRole(ctx context.Context, skip, take int32, filter, sort string) ([]user.MasterUserGroupbyRole, error) {
	return nil, nil
}
func (m *mockUserRepo) ActivateUser(ctx context.Context, value user.MasterUser) (user.MasterUser, error) {
	return user.MasterUser{}, nil
}
func (m *mockUserRepo) GetDepartementByOutsourceId(id string) ([]string, error) { return nil, nil }
func (m *mockUserRepo) GetUserPasswordHash(ctx context.Context, userId uuid.UUID) (string, error) {
	return "", nil
}

func TestGetEmployees(t *testing.T) {
	mockRepo := new(MockEmployeeRepository)
	mockUser := &mockUserRepo{}
	var conf config.Config
	// wrap mockRepo with a fake that implements the full IEmployeeRepository
	fakeRepo := &fakeEmployeeRepo{mock: mockRepo}
	useCase := NewUseCaseEmployee(fakeRepo, nil, mockUser, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, conf, nil)

	// Mock implementation for GetEmployee
	mockRepo.MockGetEmployee = func(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error) {
		return []MasterEmployee{{ID: "uuid.New()", IsActive: true, CreatedBy: "user", EmployeeName: "Employee1"}}, 1, nil
	}

	req := &pb.GetEmployeeRequest{Skip: 0, Take: 1, Filter: "", Sort: ""}

	res := &pb.GetEmployeeResponse{}
	err := useCase.GetEmployees(context.Background(), req, res, &utils.TokenValue{})

	assert.NoError(t, err)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, int32(1), res.CountData)
}

// fakeEmployeeRepo implements IEmployeeRepository by delegating a few methods
// to MockEmployeeRepository and stubbing the rest with zero values.
type fakeEmployeeRepo struct {
	mock *MockEmployeeRepository
}

func (f *fakeEmployeeRepo) GetEmployee(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error) {
	if f.mock.MockGetEmployee != nil {
		return f.mock.MockGetEmployee(ctx, skip, take, filter, sort)
	}
	return nil, 0, nil
}
func (f *fakeEmployeeRepo) GetEmployeeExcelTemplate(ctx context.Context, departmentId, OutsourceId string) ([]MasterEmployee, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) GetEmployeesAdvanceSearch(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error) {
	return nil, 0, nil
}
func (f *fakeEmployeeRepo) GetEmployeeById(ctx context.Context, id uuid.UUID) (MasterEmployee, error) {
	if f.mock.MockGetEmployeeById != nil {
		return f.mock.MockGetEmployeeById(ctx, id)
	}
	return MasterEmployee{}, nil
}
func (f *fakeEmployeeRepo) GetEmployeeEventSummary(ctx context.Context, companyId string, startDate, endDate sql.NullString) (MasterEmployeeEventSummary, error) {
	return MasterEmployeeEventSummary{}, nil
}
func (f *fakeEmployeeRepo) GetEmployeeEventList(ctx context.Context, table, filter, sort string, skip, take int32) ([]MasterEmployeeEvent, int32, error) {
	return nil, 0, nil
}
func (f *fakeEmployeeRepo) GetEmployeeMinePermitSummary(ctx context.Context, departmentId, outsourceId sql.NullString) (MasterEmployeeMinePermitSummary, error) {
	return MasterEmployeeMinePermitSummary{}, nil
}
func (f *fakeEmployeeRepo) GetEmployeeMinePermit(ctx context.Context, skip, take int32, filter, sort string) ([]MasterEmployee, int32, error) {
	return nil, 0, nil
}
func (f *fakeEmployeeRepo) CreateEmployee(tx *sql.Tx, value MasterEmployee) (MasterEmployee, error) {
	if f.mock.MockCreateEmployee != nil {
		return f.mock.MockCreateEmployee(tx, value)
	}
	return MasterEmployee{}, nil
}
func (f *fakeEmployeeRepo) UpdateEmployee(ctx context.Context, value MasterEmployee) (MasterEmployee, error) {
	if f.mock.MockUpdateEmployee != nil {
		return f.mock.MockUpdateEmployee(ctx, value)
	}
	return MasterEmployee{}, nil
}
func (f *fakeEmployeeRepo) UpdateBulkEmployees(ctx context.Context, tx *sql.Tx, list []EmployeeUpdateItem) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeProfile(ctx context.Context, value MasterEmployee) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeAddress(ctx context.Context, value MasterEmployee) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeAccount(ctx context.Context, value MasterEmployee) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeContact(ctx context.Context, value MasterEmployee) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeEmployment(ctx context.Context, value MasterEmployee) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeDocument(ctx context.Context, value MasterEmployee) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployee(tx *sql.Tx, id, modifiedBy uuid.UUID) (err error) {
	if f.mock.MockDeleteEmployee != nil {
		return f.mock.MockDeleteEmployee(tx, id, modifiedBy)
	}
	return nil
}
func (f *fakeEmployeeRepo) GetCountEmployeeByNikOrEmail(ctx context.Context, employee MasterEmployee) (int, error) {
	if f.mock.MockGetCountEmployeeByNikOrEmail != nil {
		return f.mock.MockGetCountEmployeeByNikOrEmail(ctx, employee)
	}
	return 0, nil
}
func (f *fakeEmployeeRepo) GetCountEmployeeByNik(ctx context.Context, employee MasterEmployee) (int, error) {
	return 0, nil
}
func (f *fakeEmployeeRepo) GetEmployeeAll() ([]MasterEmployee, error) {
	if f.mock.MockGetEmployeeAll != nil {
		return f.mock.MockGetEmployeeAll()
	}
	return nil, nil
}
func (f *fakeEmployeeRepo) GetEmployeeWorkExperience(ctx context.Context, employeeId uuid.UUID) ([]WorkExperienceModel, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) GetEmployeePerformance(ctx context.Context, employeeId uuid.UUID) ([]PerformanceModel, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) GetEmployeeTraining(ctx context.Context, employeeId uuid.UUID) ([]TrainingModel, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) GetEmployeeTrainingCompetency(ctx context.Context, trainingId string) ([]TrainingCompetencyModel, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) GetEmployeeDependent(ctx context.Context, employeeId uuid.UUID) ([]DependentModel, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) GetEmployeeEmergency(ctx context.Context, employeeId uuid.UUID) (EmergencyModel, error) {
	return EmergencyModel{}, nil
}
func (f *fakeEmployeeRepo) GetEmployeeReprimand(ctx context.Context, employeeId uuid.UUID) ([]ReprimandModel, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) GetEmployeeHealthRecord(ctx context.Context, employeeId uuid.UUID) ([]HealthRecordModel, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) GetEmployeeEducation(ctx context.Context, employeeId uuid.UUID) ([]EducationModel, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) InsertEmployeeMultipleWorkExperiences(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []WorkExperienceModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeWorkExperience(tx *sql.Tx, userId, employeeId string, value WorkExperienceModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeWorkExperience(tx *sql.Tx, value WorkExperienceModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeWorkExperience(tx *sql.Tx, value WorkExperienceModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeMultiplePerformances(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []PerformanceModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeePerformance(tx *sql.Tx, userId, employeeId string, value PerformanceModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeePerformance(tx *sql.Tx, value PerformanceModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeePerformance(tx *sql.Tx, value PerformanceModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeePerformanceByEmployeeId(tx *sql.Tx, value PerformanceModel) error {
	return nil
}
func (f *fakeEmployeeRepo) GetTrainingRecommendation(ctx context.Context) ([]string, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) GetTrainingOrganizerRecommendation(ctx context.Context) ([]string, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) InsertEmployeeMultipleTrainings(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []TrainingModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeTraining(tx *sql.Tx, userId, employeeId string, value TrainingModel) (string, error) {
	return "", nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeTraining(tx *sql.Tx, value TrainingModel) error { return nil }
func (f *fakeEmployeeRepo) DeleteEmployeeTraining(tx *sql.Tx, value TrainingModel) error { return nil }
func (f *fakeEmployeeRepo) DeleteEmployeeTrainingByEmployeeId(tx *sql.Tx, value TrainingModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeTrainingCompetency(tx *sql.Tx, userId, employeeId string, value TrainingCompetencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeTrainingCompetency(tx *sql.Tx, value TrainingCompetencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeTrainingCompetency(tx *sql.Tx, value TrainingCompetencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeTrainingCompetencyByTrainingId(tx *sql.Tx, value TrainingCompetencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeMultipleDependents(tx *sql.Tx, userId, employeeId string, values []DependentModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeDependent(tx *sql.Tx, employeeId string, value DependentModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeDependent(tx *sql.Tx, value DependentModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeDependent(tx *sql.Tx, value DependentModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeDependentByEmployeeId(tx *sql.Tx, value DependentModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeEmergency(tx *sql.Tx, userId, employeeId string, value EmergencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeEmergency(tx *sql.Tx, value EmergencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeEmergency(tx *sql.Tx, value EmergencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeEmergencyByEmployeeId(tx *sql.Tx, value EmergencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeReprimand(tx *sql.Tx, userId, employeeId string, value ReprimandModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeReprimand(tx *sql.Tx, value ReprimandModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeReprimand(tx *sql.Tx, value ReprimandModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeReprimandByEmployeeId(tx *sql.Tx, value ReprimandModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeMultipleHealthRecords(tx *sql.Tx, ctx context.Context, userId, employeeId string, values []HealthRecordModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeHealthRecord(tx *sql.Tx, userId, employeeId string, value HealthRecordModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeHealthRecord(tx *sql.Tx, value HealthRecordModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeHealthRecord(tx *sql.Tx, value HealthRecordModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeHealthRecordByEmployeeId(tx *sql.Tx, value HealthRecordModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeMultipleEducations(tx *sql.Tx, userId, employeeId string, values []EducationModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeEducation(tx *sql.Tx, employeeId string, value EducationModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeEducation(tx *sql.Tx, values EducationModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeEducation(tx *sql.Tx, values EducationModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeEducationByEmployeeId(tx *sql.Tx, value EducationModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeJobDescription(tx *sql.Tx, value MasterEmployee) error {
	return nil
}
func (f *fakeEmployeeRepo) GetEmployeeCompetency(ctx context.Context, employeeId uuid.UUID) ([]JobCompetencyModel, error) {
	return nil, nil
}
func (f *fakeEmployeeRepo) InsertEmployeeMultipleCompetencies(tx *sql.Tx, userId, employeeId string, values []JobCompetencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) InsertEmployeeCompetency(tx *sql.Tx, employeeId string, value JobCompetencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) UpdateEmployeeCompetency(tx *sql.Tx, value JobCompetencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) DeleteEmployeeCompetencyByEmployeeId(tx *sql.Tx, value JobCompetencyModel) error {
	return nil
}
func (f *fakeEmployeeRepo) GetDB() *sql.DB { return nil }
func (f *fakeEmployeeRepo) UpdateEmployeeEmergencyByEmployeeId(tx *sql.Tx, value EmergencyModel) error {
	return nil
}
