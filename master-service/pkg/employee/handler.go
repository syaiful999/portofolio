package employee

import (
	"context"
	pb "moyo-master-service/pkg/employee/proto"
	"moyo-master-service/utils"
)

type EmployeeHandler struct {
	srv IUseCase
}

type IEmployeeHandler interface {
	GetEmployees(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeResponse) error
	GetEmployeesAdvance(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeResponse) error
	GetEmployeesAdvanceExport(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeExportResponse) error
	GetEmployeeEventListExport(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeExportResponse) error
	GetEmployeeById(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetEmployeeByIDResponse) error
	GetTrainingOrganizerRecommendation(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error
	GetTrainingRecommendation(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error
	GetEmployeeEventSummary(ctx context.Context, req *pb.GetSummaryRequest, res *pb.GetEmployeeEventSummaryResponse) error
	GetEmployeeEventList(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeEventResponse) error
	GetEmployeeColumnList(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error
	GetEmployeeTemplateExport(ctx context.Context, req *pb.GetEmployeeTemplateRequest, res *pb.GetExportResponse) error
	GetEmployeeMinePermitSummary(ctx context.Context, req *pb.GetSummaryRequest, res *pb.GetEmployeeMinePermitSummaryResponse) error
	GetEmployeeMinePermitList(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeResponse) error
	GetEmployeeMinePermitListExport(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetExportResponse) error
	GetEmployeeJobDescriptionsCompetencies(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetEmployeeJobDescriptionResponse) error

	CreateEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeBulk(ctx context.Context, req *pb.EmployeeUploadRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeProfile(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeAddress(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeAccount(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeContact(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeePerformance(ctx context.Context, req *pb.MutationEmployeePerformanceRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeEmployment(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeWorkExperience(ctx context.Context, req *pb.MutationEmployeeWorkExperienceRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeTrainingHistory(ctx context.Context, req *pb.MutationEmployeeTrainingHistoryRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeDependent(ctx context.Context, req *pb.MutationEmployeeDependentRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeEmergency(ctx context.Context, req *pb.EmergencyModel, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeReprimand(ctx context.Context, req *pb.MutationEmployeeReprimandRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeEducation(ctx context.Context, req *pb.MutationEmployeeEducationRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeDocument(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeJobDescription(ctx context.Context, req *pb.MutationEmployeeJobDescriptionRequest, res *pb.MutationEmployeeResponse) error
	UpdateEmployeeHealthRecord(ctx context.Context, req *pb.MutationEmployeeHealthRecordRequest, res *pb.MutationEmployeeResponse) error

	DeleteEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error
	DownloadExcelEmployee(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetFileResponse) error
}

func NewEmployeeHandler(srv IUseCase) IEmployeeHandler {
	return &EmployeeHandler{
		srv: srv,
	}
}

func (e *EmployeeHandler) GetEmployees(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployees(ctx, req, res, token)
}

func (e *EmployeeHandler) GetEmployeesAdvance(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeesAdvance(ctx, req, res, token)
}

func (e *EmployeeHandler) GetEmployeesAdvanceExport(ctx context.Context, req *pb.GetEmployeeRequest, res *pb.GetEmployeeExportResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeesAdvanceExport(ctx, req, res, token)
}

func (e *EmployeeHandler) GetEmployeeEventListExport(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeExportResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeEventListExport(ctx, req, res)
}
func (e *EmployeeHandler) GetEmployeeById(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetEmployeeByIDResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeById(ctx, req, res)
}

func (e *EmployeeHandler) GetTrainingOrganizerRecommendation(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetTrainingOrganizerRecommendation(ctx, req, res)
}

func (e *EmployeeHandler) GetTrainingRecommendation(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetTrainingRecommendation(ctx, req, res)
}

func (e *EmployeeHandler) GetEmployeeEventSummary(ctx context.Context, req *pb.GetSummaryRequest, res *pb.GetEmployeeEventSummaryResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeEventSummary(ctx, req, res, token)
}

func (e *EmployeeHandler) GetEmployeeEventList(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeEventResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeEventList(ctx, req, res, token)
}

func (e *EmployeeHandler) GetEmployeeColumnList(ctx context.Context, req *pb.EmptyRequest, res *pb.RecommendationResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeColumnList(ctx, req, res)
}

func (e *EmployeeHandler) GetEmployeeTemplateExport(ctx context.Context, req *pb.GetEmployeeTemplateRequest, res *pb.GetExportResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeTemplateExport(ctx, req, res)
}

func (e *EmployeeHandler) UpdateEmployeeBulk(ctx context.Context, req *pb.EmployeeUploadRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeBulk(ctx, req, res, token)
}
func (e *EmployeeHandler) GetEmployeeMinePermitSummary(ctx context.Context, req *pb.GetSummaryRequest, res *pb.GetEmployeeMinePermitSummaryResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeMinePermitSummary(ctx, req, res)
}

func (e *EmployeeHandler) GetEmployeeMinePermitList(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetEmployeeResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeMinePermitList(ctx, req, res)
}

func (e *EmployeeHandler) GetEmployeeMinePermitListExport(ctx context.Context, req *pb.EmployeeEventRequest, res *pb.GetExportResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeMinePermitListExport(ctx, req, res)
}

func (e *EmployeeHandler) GetEmployeeJobDescriptionsCompetencies(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetEmployeeJobDescriptionResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.GetEmployeeJobDescriptionsCompetencies(ctx, req, res)
}

func (e *EmployeeHandler) CreateEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.CreateEmployee(ctx, req, res, token)
}

func (e *EmployeeHandler) UpdateEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployee(ctx, req, res, token)
}

func (e *EmployeeHandler) UpdateEmployeeProfile(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeProfile(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeeAddress(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeAddress(ctx, req, res, token)
}

func (e *EmployeeHandler) UpdateEmployeeAccount(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeAccount(ctx, req, res, token)
}

func (e *EmployeeHandler) UpdateEmployeeContact(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeContact(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeePerformance(ctx context.Context, req *pb.MutationEmployeePerformanceRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeePerformance(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeeEmployment(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeEmployment(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeeWorkExperience(ctx context.Context, req *pb.MutationEmployeeWorkExperienceRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeWorkExperience(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeeTrainingHistory(ctx context.Context, req *pb.MutationEmployeeTrainingHistoryRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeTrainingHistory(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeeDependent(ctx context.Context, req *pb.MutationEmployeeDependentRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeDependent(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeeEmergency(ctx context.Context, req *pb.EmergencyModel, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeEmergency(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeeReprimand(ctx context.Context, req *pb.MutationEmployeeReprimandRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeReprimand(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeeEducation(ctx context.Context, req *pb.MutationEmployeeEducationRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeEducation(ctx, req, res, token)
}

func (e *EmployeeHandler) UpdateEmployeeDocument(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeDocument(ctx, req, res, token)
}
func (e *EmployeeHandler) UpdateEmployeeHealthRecord(ctx context.Context, req *pb.MutationEmployeeHealthRecordRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeHealthRecord(ctx, req, res, token)
}

func (e *EmployeeHandler) UpdateEmployeeJobDescription(ctx context.Context, req *pb.MutationEmployeeJobDescriptionRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.UpdateEmployeeJobDescription(ctx, req, res, token)
}

func (e *EmployeeHandler) DeleteEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse) error {
	token, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.DeleteEmployee(ctx, req, res, token)
}

func (e *EmployeeHandler) DownloadExcelEmployee(ctx context.Context, req *pb.GetEmployeeByIDRequest, res *pb.GetFileResponse) error {
	_, err := utils.HandleToken(ctx)
	if err != nil {
		return err
	}
	return e.srv.DownloadExcelEmployee(req, res)
}
