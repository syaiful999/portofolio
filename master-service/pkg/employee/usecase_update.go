package employee

import (
	"context"
	"database/sql"
	pb "moyo-master-service/pkg/employee/proto"
	"moyo-master-service/utils"
	"net/http"
)

func (s *UseCase) UpdateEmployee(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error
	var result MasterEmployee

	birthDate, _ := utils.ConvStringToTime(req.Birthdate)
	joinDate, _ := utils.ConvStringToNullTime(req.JoinDate)

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
		ID:         req.Id,
		ModifiedBy: token.ID.String(),

		EmployeeName: req.EmployeeName,
		Nik:          req.Nik,
		NoKtp:        req.Noktp,
		Npwp:         utils.ConvStringToNullString(req.Npwp),

		Email:       req.Email,
		Phonenumber: utils.ConvStringToNullString(req.Phonenumber),
		Birthplace:  utils.ConvStringToNullString(req.Birthplace),
		Birthdate:   birthDate,
		Address:     utils.ConvStringToNullString(req.Address),
		KelurahanId: sql.NullInt32{Int32: req.KelurahanId, Valid: true},

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

		GradeId:      utils.ConvStringToNullString(req.GradeId),
		JabatanId:    utils.ConvStringToNullString(req.JabatanId),
		Position:     utils.ConvStringToNullString(req.Position),
		CostCenterId: utils.ConvStringToNullString(req.CostCenterId),
		DepartmentId: utils.ConvStringToNullString(req.DepartmentId),
		OutsourceId:  utils.ConvStringToNullString(req.OutsourceId),
	}

	//check if exist
	countData, err := s.repository.GetCountEmployeeByNikOrEmail(ctx, data)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "UpdateEmployee:GetCountEmployeeByNikOrEmail", err.Error())
	}
	if countData > 0 {
		errStr := "Nik, no ktp, & email must be unique"
		return errorResponse3(res, http.StatusBadRequest, "UpdateEmployee:"+errStr, "")
	}

	if result, err = s.repository.UpdateEmployee(ctx, data); err != nil {
		return errorResponse3(res, http.StatusBadRequest, "UpdateEmployee:", err.Error())
	}

	res.Data = &pb.GetEmployeeByIDRequest{Id: result.ID}
	return nil
}

func (s *UseCase) UpdateEmployeeProfile(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error

	birthDate, _ := utils.ConvStringToTime(req.Birthdate)
	retirementDate, _ := utils.ConvStringToNullTime(req.RetirementDate)

	data := MasterEmployee{
		ID:         req.Id,
		ModifiedBy: token.ID.String(),

		EmployeeName: req.EmployeeName,
		Nik:          req.Nik,
		GenderId:     utils.ConvStringToNullString(req.GenderId),

		Birthdate:        birthDate,
		Birthplace:       utils.ConvStringToNullString(req.Birthplace),
		ReligionId:       utils.ConvStringToNullString(req.ReligionId),
		StatusMarriageId: utils.ConvStringToNullString(req.StatusMarriageId),

		Picture:       utils.ConvStringToNullString(req.Picture),
		BloodTypeId:   utils.ConvStringToNullString(req.BloodTypeId),
		Npwp:          utils.ConvStringToNullString(req.Npwp),
		NoKtp:         req.Noktp,
		BpjskesNumber: utils.ConvStringToNullString(req.BpjskesNumber),
		BpjstkNumber:  utils.ConvStringToNullString(req.BpjstkNumber),

		RetirementDate: retirementDate,
		StatusPtkpId:   utils.ConvStringToNullString(req.StatusPtkpId),
	}

	//check if exist
	countData, err := s.repository.GetCountEmployeeByNik(ctx, data)
	if err != nil {
		return errorResponse3(res, http.StatusBadRequest, "UpdateEmployee:GetCountEmployeeByNikOrEmail", err.Error())
	}
	if countData > 0 {
		errStr := "Nik, no ktp, & email must be unique"
		return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployee:"+errStr, "")
	}

	if err = s.repository.UpdateEmployeeProfile(ctx, data); err != nil {
		return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeProfile:", err.Error())
	}

	res.Data = &pb.GetEmployeeByIDRequest{Id: data.ID}
	return nil
}

func (s *UseCase) UpdateEmployeeAccount(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error

	data := MasterEmployee{
		ID:         req.Id,
		ModifiedBy: token.ID.String(),

		CostCenterId:     utils.ConvStringToNullString(req.CostCenterId),
		GlNumberId:       utils.ConvStringToNullString(req.GlNumberId),
		StandardOvertime: utils.ConvStringToNullString(req.StandardOvertime),
		StandardWorkday:  utils.ConvStringToNullString(req.StandardWorkday),
	}

	if err = s.repository.UpdateEmployeeAccount(ctx, data); err != nil {
		return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeAccount:", err.Error())
	}

	res.Data = &pb.GetEmployeeByIDRequest{Id: data.ID}
	return nil
}
func (s *UseCase) UpdateEmployeeAddress(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error

	data := MasterEmployee{
		ID:         req.Id,
		ModifiedBy: token.ID.String(),

		DomisiliKelurahanId: utils.ConvInt32toNullInt32(req.DomisiliKelurahanId),
		KelurahanId:         utils.ConvInt32toNullInt32(req.KelurahanId),
	}

	if err = s.repository.UpdateEmployeeAddress(ctx, data); err != nil {
		return errorResponse3(res, http.StatusBadRequest, "UpdateEmployeeAddress:", err.Error())
	}

	res.Data = &pb.GetEmployeeByIDRequest{Id: data.ID}
	return nil
}

func (s *UseCase) UpdateEmployeeContact(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error

	data := MasterEmployee{
		ID:         req.Id,
		ModifiedBy: token.ID.String(),

		EmailOffice: utils.ConvStringToNullString(req.EmailOffice),
		Email:       req.Email,
		Phonenumber: utils.ConvStringToNullString(req.Phonenumber),
	}

	if err = s.repository.UpdateEmployeeContact(ctx, data); err != nil {
		return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeContact:", err.Error())
	}

	res.Data = &pb.GetEmployeeByIDRequest{Id: data.ID}
	return nil
}

func (s *UseCase) UpdateEmployeePerformance(ctx context.Context, req *pb.MutationEmployeePerformanceRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	modifiedBy := token.ID.String()
	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "failed to UpdateEmployeePerformance:db Begin", err.Error())
	}
	for _, v := range req.Performances {

		data := PerformanceModel{
			ModifiedBy:  modifiedBy,
			Id:          v.Id,
			Periode:     utils.ConvStringToNullString(v.Periode),
			Predicate:   utils.ConvStringToNullString(v.Predicate),
			Score:       utils.ConvStringToNullString(v.Score),
			Description: utils.ConvStringToNullString(v.Description),
		}
		if v.IsActive == false {
			// delete
			if err = s.repository.DeleteEmployeePerformance(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeePerformance:", err.Error())
			}
		} else if v.Id == "" {
			// create
			if err = s.repository.InsertEmployeePerformance(tx, modifiedBy, req.Id, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeePerformance:", err.Error())
			}
		} else {
			// update
			if err = s.repository.UpdateEmployeePerformance(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeePerformance:", err.Error())
			}
		}
	}
	tx.Commit()

	res.Data = &pb.GetEmployeeByIDRequest{Id: ""}
	return nil
}

func (s *UseCase) UpdateEmployeeEmployment(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error

	joinDate, _ := utils.ConvStringToNullTime(req.JoinDate)
	minePermitDate, _ := utils.ConvStringToNullTime(req.MinePermitDate)
	contractEndDate, _ := utils.ConvStringToNullTime(req.ContractEndDate)
	contractStartDate, _ := utils.ConvStringToNullTime(req.ContractStartDate)

	data := MasterEmployee{
		ID:         req.Id,
		ModifiedBy: token.ID.String(),

		JoinDate:           joinDate,
		JabatanId:          utils.ConvStringToNullString(req.JabatanId),
		GradeId:            utils.ConvStringToNullString(req.GradeId),
		MinePermitDate:     minePermitDate,
		ContractStartDate:  contractStartDate,
		ContractEndDate:    contractEndDate,
		ContractRenewal:    utils.ConvInt32toNullInt32(req.ContractRenewal),
		OutsourceId:        utils.ConvStringToNullString(req.OutsourceId),
		SectionId:          utils.ConvStringToNullString(req.SectionId),
		StatusContractId:   utils.ConvStringToNullString(req.StatusContractId),
		DepartmentDetailId: utils.ConvStringToNullString(req.DepartmentDetailId),
		DepartmentId:       utils.ConvStringToNullString(req.DepartmentId),
	}

	if err = s.repository.UpdateEmployeeEmployment(ctx, data); err != nil {
		return errorResponse3(res, http.StatusBadRequest, "UpdateEmployeeEmployment:", err.Error())
	}

	res.Data = &pb.GetEmployeeByIDRequest{Id: data.ID}
	return nil
}

func (s *UseCase) UpdateEmployeeWorkExperience(ctx context.Context, req *pb.MutationEmployeeWorkExperienceRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {

	employeeId := req.Id
	userId := token.ID.String()

	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "failed to UpdateEmployeeWorkExperience:db Begin", err.Error())
	}
	for _, v := range req.WorkExperiences {

		joinDate, _ := utils.ConvStringToNullTime(v.JoinDate)
		endDate, _ := utils.ConvStringToNullTime(v.EndDate)

		data := WorkExperienceModel{
			Id:          v.Id,
			ModifiedBy:  token.ID.String(),
			CompanyName: utils.ConvStringToNullString(v.CompanyName),
			Position:    utils.ConvStringToNullString(v.Position),
			JoinDate:    joinDate,
			EndDate:     endDate,
		}
		if v.IsActive == false {

			if err = s.repository.DeleteEmployeeWorkExperience(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeWorkExperience:", err.Error())
			}
		} else if data.Id == "" {

			if err = s.repository.InsertEmployeeWorkExperience(tx, userId, employeeId, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeWorkExperience:", err.Error())
			}
		} else {

			if err = s.repository.UpdateEmployeeWorkExperience(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeWorkExperience:", err.Error())
			}

		}
	}
	tx.Commit()
	res.Data = &pb.GetEmployeeByIDRequest{Id: ""}

	return nil
}
func (s *UseCase) UpdateEmployeeTrainingHistory(ctx context.Context, req *pb.MutationEmployeeTrainingHistoryRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	employeeId := req.Id
	userId := token.ID.String()

	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "failed to UpdateEmployeeWorkExperience:db Begin", err.Error())
	}
	for _, v := range req.Trainings {

		trainingDate, _ := utils.ConvStringToNullTime(v.TrainingDate)
		expiryDate, _ := utils.ConvStringToNullTime(v.ExpiryDate)
		effectiveDate, _ := utils.ConvStringToNullTime(v.EffectiveDate)

		data := TrainingModel{
			Id:                v.Id,
			ModifiedBy:        token.ID.String(),
			TrainingTitle:     utils.ConvStringToNullString(v.TrainingTitle),
			TrainingDate:      trainingDate,
			OrganizerName:     utils.ConvStringToNullString(v.OrganizerName),
			CertificateNumber: utils.ConvStringToNullString(v.CertificateNumber),
			TrainingMethod:    utils.ConvStringToNullString(v.TrainingMethod),
			ExpiryDate:        expiryDate,
			EffectiveDate:     effectiveDate,
			Cost:              utils.ConvStringToNullString(v.Cost),
			CertificateUrl:    utils.ConvStringToNullString(v.CertificateUrl),
		}
		if v.IsActive == false {

			if err = s.repository.DeleteEmployeeTraining(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeTraining:", err.Error())
			}
			competency := TrainingCompetencyModel{
				TrainingId: utils.ConvStringToNullString(data.Id),
				ModifiedBy: userId,
			}

			if err = s.repository.DeleteEmployeeTrainingCompetencyByTrainingId(tx, competency); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeTrainingCompetencyByTrainingId", err.Error())
			}

		} else if data.Id == "" {
			var trainingId string
			if trainingId, err = s.repository.InsertEmployeeTraining(tx, userId, employeeId, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeTraining:", err.Error())
			}

			for _, element := range v.Competencies {
				competency := TrainingCompetencyModel{
					TrainingId:      utils.ConvStringToNullString(trainingId),
					JobCompetencyId: utils.ConvStringToNullString(element.JobCompetencyId),
					Remark:          utils.ConvStringToNullString(element.Remark),
					IsRequired:      element.IsRequired,
					OtherCompetency: utils.ConvStringToNullString(element.OtherCompetency),
				}
				if err = s.repository.InsertEmployeeTrainingCompetency(tx, userId, employeeId, competency); err != nil {
					tx.Rollback()
					return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeTrainingCompetency:", err.Error())
				}
			}
		} else {
			if err = s.repository.UpdateEmployeeTraining(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeTraining:", err.Error())
			}
			for _, element := range v.Competencies {
				competency := TrainingCompetencyModel{
					Id:              element.Id,
					TrainingId:      utils.ConvStringToNullString(data.Id),
					JobCompetencyId: utils.ConvStringToNullString(element.JobCompetencyId),
					Remark:          utils.ConvStringToNullString(element.Remark),
					IsRequired:      element.IsRequired,
					IsActive:        element.IsActive,
					OtherCompetency: utils.ConvStringToNullString(element.OtherCompetency),
					ModifiedBy:      userId,
				}

				if competency.Id != "" {
					if err = s.repository.UpdateEmployeeTrainingCompetency(tx, competency); err != nil {
						tx.Rollback()
						return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeTrainingCompetency:", err.Error())
					}
					continue
				} else if element.IsActive == false {
					if err = s.repository.DeleteEmployeeTrainingCompetencyByTrainingId(tx, competency); err != nil {
						tx.Rollback()
						return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeTrainingCompetencyByTrainingId:", err.Error())
					}
					continue
				} else {
					if err = s.repository.InsertEmployeeTrainingCompetency(tx, userId, employeeId, competency); err != nil {
						tx.Rollback()
						return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeTrainingCompetency:", err.Error())
					}
				}
			}
		}
	}
	tx.Commit()
	res.Data = &pb.GetEmployeeByIDRequest{Id: ""}

	return nil
}
func (s *UseCase) UpdateEmployeeDependent(ctx context.Context, req *pb.MutationEmployeeDependentRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	employeeId := req.Id

	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "failed to UpdateEmployeeDependent:db Begin", err.Error())
	}
	for _, v := range req.Dependents {

		birthdate, _ := utils.ConvStringToNullTime(v.Birthdate)

		data := DependentModel{
			Id:            v.Id,
			ModifiedBy:    token.ID.String(),
			DependentName: utils.ConvStringToNullString(v.DependentName),
			Birthdate:     birthdate,
			Birthplace:    utils.ConvStringToNullString(v.Birthplace),
			Phonenumber:   utils.ConvStringToNullString(v.Phonenumber),
			RelationId:    utils.ConvStringToNullString(v.RelationId),
		}
		if v.IsActive == false {
			if err = s.repository.DeleteEmployeeDependent(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeDependent:", err.Error())
			}
		} else if data.Id == "" {
			if err = s.repository.InsertEmployeeDependent(tx, employeeId, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeDependent:", err.Error())
			}
		} else {
			if err = s.repository.UpdateEmployeeDependent(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeDependent:", err.Error())
			}

		}
	}
	tx.Commit()
	res.Data = &pb.GetEmployeeByIDRequest{Id: ""}

	return nil
}

func (s *UseCase) UpdateEmployeeEmergency(ctx context.Context, req *pb.EmergencyModel, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	userId := token.ID.String()

	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "UpdateEmployeeWorkExperience:db Begin", err.Error())
	}

	data := EmergencyModel{
		EmployeeId:    req.Id,
		EmergencyName: utils.ConvStringToNullString(req.EmergencyName),
		RelationId:    utils.ConvStringToNullString(req.RelationId),
		Phonenumber:   utils.ConvStringToNullString(req.Phonenumber),
		KelurahanId:   utils.ConvInt32toNullInt32(req.KelurahanId),
		ModifiedBy:    userId,
	}
	// currently 1 employee only has 1 emergency
	if err = s.repository.UpdateEmployeeEmergencyByEmployeeId(tx, data); err != nil {
		tx.Rollback()
		return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeEmergency:", err.Error())
	}

	/*
		if req.EmergencyId == "" {
			if err = s.repository.InsertEmployeeEmergency(tx, data.ModifiedBy, req.Id, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeEmergency:", err.Error())
			}
		} else if req.IsActive == false {

			if err = s.repository.DeleteEmployeeEmergency(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeEmergency:", err.Error())
			}
		} else {

			if err = s.repository.UpdateEmployeeEmergency(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeEmergency:", err.Error())
			}
		// }*/

	tx.Commit()
	res.Data = &pb.GetEmployeeByIDRequest{Id: ""}

	return nil
}

func (s *UseCase) UpdateEmployeeReprimand(ctx context.Context, req *pb.MutationEmployeeReprimandRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	employeeId := req.Id
	userId := token.ID.String()

	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "UpdateEmployeeWorkExperience:db Begin", err.Error())
	}
	for _, v := range req.Reprimands {

		validTime, _ := utils.ConvStringToNullTime(v.ValidTime)
		startDate, _ := utils.ConvStringToNullTime(v.StartDate)
		endDate, _ := utils.ConvStringToNullTime(v.EndDate)

		data := ReprimandModel{
			Id:                v.Id,
			ModifiedBy:        token.ID.String(),
			AttachmentUrl:     utils.ConvStringToNullString(v.AttachmentUrl),
			WarningLevelId:    utils.ConvStringToNullString(v.WarningLevelId),
			WarningLevelValue: utils.ConvStringToNullString(v.WarningLevelValue),
			ValidTime:         validTime,
			StartDate:         startDate,
			EndDate:           endDate,
			DocumentNumber:    utils.ConvStringToNullString(v.DocumentNumber),
			Description:       utils.ConvStringToNullString(v.Description),
		}
		if v.IsActive == false {

			if err = s.repository.DeleteEmployeeReprimand(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeReprimand:", err.Error())
			}
		} else if data.Id == "" {

			if err = s.repository.InsertEmployeeReprimand(tx, userId, employeeId, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeReprimand:", err.Error())
			}
		} else {

			if err = s.repository.UpdateEmployeeReprimand(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeReprimand:", err.Error())
			}

		}
	}
	tx.Commit()
	res.Data = &pb.GetEmployeeByIDRequest{Id: ""}

	return nil
}

func (s *UseCase) UpdateEmployeeEducation(ctx context.Context, req *pb.MutationEmployeeEducationRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	employeeId := req.Id

	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "failed to UpdateEmployeeEducation:db Begin", err.Error())
	}
	for _, v := range req.Educations {
		data := EducationModel{
			Id:           v.Id,
			InstanceName: utils.ConvStringToNullString(v.InstanceName),
			Major:        utils.ConvStringToNullString(v.Major),
			Degree:       utils.ConvStringToNullString(v.Degree),
			DegreeId:     utils.ConvStringToNullString(v.DegreeId),
			GPA:          v.Gpa,
			GraduateYear: v.GraduateYear,
			ModifiedBy:   token.ID.String(),
		}
		if v.IsActive == false {

			if err = s.repository.DeleteEmployeeEducation(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeEducation:", err.Error())
			}
		} else if data.Id == "" {

			if err = s.repository.InsertEmployeeEducation(tx, employeeId, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeEducation:", err.Error())
			}
		} else {

			if err = s.repository.UpdateEmployeeEducation(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeEducation:", err.Error())
			}

		}
	}
	tx.Commit()
	res.Data = &pb.GetEmployeeByIDRequest{Id: ""}

	return nil
}

func (s *UseCase) UpdateEmployeeHealthRecord(ctx context.Context, req *pb.MutationEmployeeHealthRecordRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	employeeId := req.Id

	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "UpdateEmployeeWorkExperience:db Begin", err.Error())
	}

	for _, v := range req.HealthRecords {
		mcuDate, _ := utils.ConvStringToNullTime(v.McuDate)

		data := HealthRecordModel{
			Id:         v.Id,
			ModifiedBy: token.ID.String(),
			// Map other fields

			McuDate:           mcuDate,
			Periode:           utils.ConvStringToNullString(v.Periode),
			StatusId:          utils.ConvStringToNullString(v.StatusId),
			HealthDescription: utils.ConvStringToNullString(v.HealthDescription),
			McuUrl:            utils.ConvStringToNullString(v.McuUrl),
			McuFollowupUrl:    utils.ConvStringToNullString(v.McuFollowupUrl),
			StatusFollowupId:  utils.ConvStringToNullString(v.StatusFollowupId),
		}
		if v.IsActive == false {

			if err = s.repository.DeleteEmployeeHealthRecord(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeHealthRecord:", err.Error())
			}
		} else if data.Id == "" {

			if err = s.repository.InsertEmployeeHealthRecord(tx, data.ModifiedBy, employeeId, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeHealthRecord:", err.Error())
			}
		} else {

			if err = s.repository.UpdateEmployeeHealthRecord(tx, data); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeHealthRecord:", err.Error())
			}

		}
	}
	tx.Commit()
	res.Data = &pb.GetEmployeeByIDRequest{Id: ""}

	return nil
}

func (s *UseCase) UpdateEmployeeDocument(ctx context.Context, req *pb.MutationEmployeeRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error

	data := MasterEmployee{
		ID:         req.Id,
		ModifiedBy: token.ID.String(),

		IjazahUrl:     utils.ConvStringToNullString(req.IjazahUrl),
		MinePermitUrl: utils.ConvStringToNullString(req.MinePermitUrl),
		KtpUrl:        utils.ConvStringToNullString(req.KtpUrl),
		NpwpUrl:       utils.ConvStringToNullString(req.NpwpUrl),
		BpjstkUrl:     utils.ConvStringToNullString(req.BpjstkUrl),
		BpjskesUrl:    utils.ConvStringToNullString(req.BpjskesUrl),
		KkUrl:         utils.ConvStringToNullString(req.KkUrl),
	}

	if err = s.repository.UpdateEmployeeDocument(ctx, data); err != nil {
		return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeDocument:", err.Error())
	}

	res.Data = &pb.GetEmployeeByIDRequest{Id: data.ID}
	return nil
}

func (s *UseCase) UpdateEmployeeJobDescription(ctx context.Context, req *pb.MutationEmployeeJobDescriptionRequest, res *pb.MutationEmployeeResponse, token *utils.TokenValue) error {
	var err error
	var employeeCompetencies []JobCompetencyModel

	for _, v := range req.Competencies {
		employeeCompetencies = append(employeeCompetencies, JobCompetencyModel{
			IsActive:         v.IsActive,
			IsRequired:       v.IsRequired,
			Id:               v.Id,
			JobCompetencyId:  utils.ConvStringToNullString(v.JobCompetencyId),
			JobDescriptionId: utils.ConvStringToNullString(req.JobDescriptionId),
		})
	}

	data := MasterEmployee{
		ID:               req.Id,
		ModifiedBy:       token.ID.String(),
		JobDescriptionId: utils.ConvStringToNullString(req.JobDescriptionId),
	}
	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, http.StatusInternalServerError, "UpdateEmployeeJobDescription:db Begin", err.Error())
	}

	if req.JobDescriptionId != req.JobDescriptionPreviousId {
		// 1. update job description
		if err = s.repository.UpdateEmployeeJobDescription(tx, data); err != nil {
			tx.Rollback()
			return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeDocument:", err.Error())
		}

		// 2. delete current employee competencies
		dataCompetency := JobCompetencyModel{
			EmployeeID: req.Id,
			ModifiedBy: token.ID.String(),
		}

		if err = s.repository.DeleteEmployeeCompetencyByEmployeeId(tx, dataCompetency); err != nil {
			tx.Rollback()
			return errorResponse3(res, http.StatusBadRequest, "failed to DeleteEmployeeCompetencyByEmployeeId:", err.Error())
		}
		// 3. insert new employee competencies

		if err = s.repository.InsertEmployeeMultipleCompetencies(tx, token.ID.String(), req.Id, employeeCompetencies); err != nil {
			tx.Rollback()
			return errorResponse3(res, http.StatusBadRequest, "failed to InsertEmployeeMultipleCompetencies:", err.Error())
		}
	} else {
		// update current employee competencies
		for _, v := range employeeCompetencies {
			if err = s.repository.UpdateEmployeeCompetency(tx, v); err != nil {
				tx.Rollback()
				return errorResponse3(res, http.StatusBadRequest, "failed to UpdateEmployeeCompetency:", err.Error())
			}
		}
	}
	tx.Commit()

	res.Data = &pb.GetEmployeeByIDRequest{Id: data.ID}
	return nil
}
