package employee

import (
	"bytes"
	"context"
	"fmt"
	"moyo-master-service/pkg/department"
	"moyo-master-service/pkg/departmentdetail"
	pb "moyo-master-service/pkg/employee/proto"
	"moyo-master-service/pkg/enum"
	"moyo-master-service/pkg/jobdescription"
	"moyo-master-service/pkg/jobfamily"
	"moyo-master-service/pkg/kelurahan"
	"moyo-master-service/pkg/outsource"
	"moyo-master-service/pkg/section"
	"moyo-master-service/utils"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/xuri/excelize/v2"
)

func (s *UseCase) DownloadExcelEmployee(req *pb.GetEmployeeByIDRequest, res *pb.GetFileResponse) error {

	now := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
	filename := fmt.Sprintf("employee_%s.xlsx", now)

	// create excel
	xlsx := excelize.NewFile()
	sheet := "Employee"

	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet)

	headers := []string{
		"NIK", "Name", "Grade", "Service Contract", "Department",
		"Jabatan", "Tempat Lahir", "Tanggal Lahir", "No KTP", "NPWP",
		"Alamat Lengkap", "Email", "No HP", "Tanggal Masuk", "Gender",
		"Status Kawin", "Golongan Darah", "Agama", "Status PTKP",
		"G/L Number", "Cost Center", "Shift Kerja", "Tingkat Pendidikan",
		"Status Kontrak", "Standar Hari Kerja", "Standar Overtime",
		"Sertifikasi", "Posisi",
	}

	headerStyle, _ := xlsx.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Bold: true, Family: "calibri", Size: 12},
	})

	for col, header := range headers {
		colName, _ := excelize.ColumnNumberToName(col + 1)
		cell := colName + "1"
		xlsx.SetCellValue(sheet, cell, header)
		xlsx.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	data, err := s.repository.GetEmployeeAll()
	if err != nil {
		return errorResponse4(res, http.StatusBadRequest, "failed to DownloadExcelEmployee", err.Error())
	}

	for i, row := range data {
		BirthDate := row.Birthdate.Format("2006-01-02 15:04")
		JoinDate := row.JoinDate.Time.Format("2006-01-02 15:04")

		xlsx.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), row.Nik)
		xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), row.EmployeeName)
		xlsx.SetCellValue(sheet, fmt.Sprintf("C%d", i+2), row.Grade.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("D%d", i+2), row.OutsourceName.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("E%d", i+2), row.DepartmentName.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("F%d", i+2), row.Jabatan.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("G%d", i+2), row.Birthplace.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("H%d", i+2), BirthDate)
		xlsx.SetCellValue(sheet, fmt.Sprintf("I%d", i+2), row.NoKtp)
		xlsx.SetCellValue(sheet, fmt.Sprintf("J%d", i+2), row.Npwp.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("K%d", i+2), row.Address.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("L%d", i+2), row.Email)
		xlsx.SetCellValue(sheet, fmt.Sprintf("M%d", i+2), row.Phonenumber.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("N%d", i+2), JoinDate)
		xlsx.SetCellValue(sheet, fmt.Sprintf("O%d", i+2), row.Gender.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("P%d", i+2), row.StatusMarriage.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("Q%d", i+2), row.BloodType.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("R%d", i+2), row.Religion.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("S%d", i+2), row.StatusPtkp.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("T%d", i+2), row.GlNumber.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("U%d", i+2), row.CostCenter.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("V%d", i+2), row.ShiftType.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("W%d", i+2), row.LastEducation.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("X%d", i+2), row.StatusContract.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("Y%d", i+2), row.StandardWorkday.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("Z%d", i+2), row.StandardOvertime.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("AA%d", i+2), row.Certification.String)
		xlsx.SetCellValue(sheet, fmt.Sprintf("AB%d", i+2), row.Position.String)
	}

	// --- simpan ke buffer dan upload ke MinIO ---
	buf := new(bytes.Buffer)
	if err := xlsx.Write(buf); err != nil {
		utils.PushLogf("", "[DownloadExcelEmployee]error writing buffer =>", err.Error())
		return err
	}

	minioClient, err := utils.InitMinioClient(utils.MinioConfig{
		Endpoint:  s.conf.Hosts.Minio.Address,
		AccessKey: s.conf.Hosts.Minio.AccessKey,
		SecretKey: s.conf.Hosts.Minio.SecretKey,
		UseSSL:    s.conf.Hosts.Minio.UseSSL,
	})
	if err != nil {
		utils.PushLogf("", "[DownloadExcelEmployee]init minio client failed =>", err.Error())
		return err
	}

	_, err = minioClient.PutObject(context.Background(), "documents", filename, buf, int64(buf.Len()), minio.PutObjectOptions{
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	})
	if err != nil {
		utils.PushLogf("", "[DownloadExcelEmployee]upload to minio failed =>", err.Error())
		return err
	}

	res.IsError = false
	res.Url = s.conf.Hosts.Minio.Address + "/documents/" + filename
	return nil
}

const (
	MainSheet = "Employee Template"
	CopySheet = "EmployeeCopy"
	DsSheet   = "datasource"
)

// GetEmployeeTemplateExport generates an Excel template for bulk updating employees.
// It includes:
// 1. A 'Main Sheet' with existing employee data (if filters applied) or empty rows.
// 2. Data Validation (dropdowns) for various fields, including dependent dropdowns (e.g., City depends on Province).
// 3. Conditional formatting to highlight changed cells compared to the original data.
func (s *UseCase) GetEmployeeTemplateExport(ctx context.Context, req *pb.GetEmployeeTemplateRequest, res *pb.GetExportResponse) error {
	// ======================================================================================================================================= #region-1 set GET DATA SOURCES

	// #region 0. Handle Columns
	// Parse the requested columns to determine which fields to include in the template.
	listColumns := strings.Split(req.ListColumns, ",")

	columnHandle := HandleReqColumn(listColumns)

	// ================================================region 1. Fetch Data
	// Retrieve all necessary master data lists and existing employee data.
	// We only fetch data for columns that are actually requested (IsCollected is true).
	var err error
	var departmentDetails []departmentdetail.MasterDepartmentDetail

	// department_detail (BE)== department (FE)
	if columnHandle["department_detail_id"].IsCollected {

		departmentDetails, err = s.repositoryDepartmentDetail.GetDepartmentDetailList()
		if err != nil {
			return err
		}
		// Sort Department Details by DepartmentID to ensure grouping for Excel OFFSET formula
		// This sorting is crucial for the dependent dropdown logic (OFFSET + MATCH + COUNTIF) to work correctly.
		sort.Slice(departmentDetails, func(i, j int) bool {
			if departmentDetails[i].DepartmentID.String() == departmentDetails[j].DepartmentID.String() {
				return departmentDetails[i].DepartmentDetailName < departmentDetails[j].DepartmentDetailName
			}
			return departmentDetails[i].DepartmentID.String() < departmentDetails[j].DepartmentID.String()
		})
	}
	// department (BE)== division (FE)
	var departments []department.MasterDepartment

	if columnHandle["department_id"].IsCollected {

		departments, err = s.repositoryDepartment.GetDepartmentList()
		if err != nil {
			return err
		}
	}

	var outsources []outsource.TransactOutsource
	if columnHandle["outsource_id"].IsCollected {
		outsources, err = s.repositoryOutsource.GetOutsourceList()
		if err != nil {
			return err
		}
	}
	var jabatans []jobfamily.MasterJobFamily
	if columnHandle["jabatan_id"].IsCollected {
		jabatans, err = s.repositoryJobFamily.GetJobFamilyList()
		if err != nil {
			return err
		}
	}
	var jobDescriptions []jobdescription.MasterJobDescription
	if columnHandle["job_description_id"].IsCollected {
		jobDescriptions, err = s.repositoryJobDescription.GetJobDescriptionList()
		if err != nil {
			return err
		}
	}
	var sections []section.MasterSection
	if columnHandle["section_id"].IsCollected {

		sections, err = s.repositorySection.GetSectionList()
		if err != nil {
			return err
		}
		// Sort Sections by DepartmentDetailID to ensure grouping for Excel OFFSET formula
		// Grouping is required so that all sections for a specific department detail are contiguous rows.
		sort.Slice(sections, func(i, j int) bool {
			if sections[i].DepartmentDetailID.String() == sections[j].DepartmentDetailID.String() {
				return sections[i].SectionName < sections[j].SectionName
			}
			return sections[i].DepartmentDetailID.String() < sections[j].DepartmentDetailID.String()
		})
	}

	var kelurahans []kelurahan.MasterKelurahan

	if columnHandle["kelurahan_id"].IsCollected || columnHandle["kecamatan_id"].IsCollected || columnHandle["city_id"].IsCollected || columnHandle["province_id"].IsCollected ||
		columnHandle["domisili_kelurahan_id"].IsCollected || columnHandle["domisili_kecamatan_id"].IsCollected || columnHandle["domisili_city_id"].IsCollected || columnHandle["domisili_province_id"].IsCollected {
		kelurahans, err = s.repositoryKelurahan.GetAllKelurahans(ctx)
		if err != nil {
			return err
		}
	}

	// Fetch existing employee data based on request filters (DepartmentId, OutsourceId).
	var employees []MasterEmployee
	employees, err = s.repository.GetEmployeeExcelTemplate(ctx, req.DepartmentId, req.OutsourceId)
	if err != nil {
		return err
		// }
	}

	// Fetch master enum data and map it by EnumType for easy access.
	var enum []enum.MasterEnum
	if columnHandle["grade_id"].IsCollected || columnHandle["status_ptkp_id"].IsCollected ||
		columnHandle["cost_center_id"].IsCollected || columnHandle["gl_number_id"].IsCollected || columnHandle["status_marriage_id"].IsCollected ||
		columnHandle["gender_id"].IsCollected || columnHandle["blood_type_id"].IsCollected || columnHandle["religion_id"].IsCollected {
		enum, err = s.repositoryEnum.GetEnumExcel()
		if err != nil {
			return err
		}
	}
	// HandleEnum is assumed to convert []MasterEnum to map[string][]Option
	enumMap := HandleEnum(enum, columnHandle)
	// ======================================================================================================================================= #endregion

	// ========================#region 2. Create File
	// Initialize a new Excel file and remove the default sheet.
	f := excelize.NewFile()

	defaultSheet := f.GetSheetName(0)
	f.DeleteSheet(defaultSheet)
	// ======================================================================================================================================= #endregion

	// ======================================================================================================================================= #region 3. LIST DATASOURCE SHEET
	// Create the 'datasource' sheet to store reference data for dropdowns (Data Validation).
	// This sheet will be hidden from the user.

	f.NewSheet(DsSheet)

	// --- Populate Department Data ---
	startRow := 2
	headerRow := 1
	f.SetCellValue(DsSheet, fmt.Sprintf("A%d", headerRow), "Department")
	f.SetCellValue(DsSheet, fmt.Sprintf("B%d", headerRow), "Department Id")
	for i, d := range departments {
		f.SetCellValue(DsSheet, fmt.Sprintf("A%d", i+startRow), d.DepartmentName.String)
		f.SetCellValue(DsSheet, fmt.Sprintf("B%d", i+startRow), d.ID.String())
	}

	// outsource
	// --- Populate Outsource Data ---
	f.SetCellValue(DsSheet, fmt.Sprintf("C%d", headerRow), "Outsource")
	f.SetCellValue(DsSheet, fmt.Sprintf("D%d", headerRow), "Outsource Id")

	for i, o := range outsources {
		f.SetCellValue(DsSheet, fmt.Sprintf("C%d", i+startRow), o.OutsourceName)
		f.SetCellValue(DsSheet, fmt.Sprintf("D%d", i+startRow), o.OutsourceId)
	}

	// --- Populate Enum Data (Jabatan, Grade, Status Contract, Cost Center, GL Number) ---
	// Note: The structure for populating enums is repetitive but standard for generating data source lists.
	// enum [1]-start
	f.SetCellValue(DsSheet, fmt.Sprintf("E%d", headerRow), "Jabatan")
	f.SetCellValue(DsSheet, fmt.Sprintf("F%d", headerRow), "Jabatan Id")
	for i, o := range jabatans {
		f.SetCellValue(DsSheet, fmt.Sprintf("E%d", i+startRow), o.JobFamilyName)
		f.SetCellValue(DsSheet, fmt.Sprintf("F%d", i+startRow), o.ID)
	}
	f.SetCellValue(DsSheet, fmt.Sprintf("G%d", headerRow), "Grade")
	f.SetCellValue(DsSheet, fmt.Sprintf("H%d", headerRow), "Grade Id")
	for i, o := range enumMap["grade"] {
		f.SetCellValue(DsSheet, fmt.Sprintf("G%d", i+startRow), o.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("H%d", i+startRow), o.Id)
	}
	f.SetCellValue(DsSheet, fmt.Sprintf("I%d", headerRow), "Status Contract")
	f.SetCellValue(DsSheet, fmt.Sprintf("J%d", headerRow), "Status Contract Id")
	for i, o := range enumMap["status_contract"] {
		f.SetCellValue(DsSheet, fmt.Sprintf("I%d", i+startRow), o.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("J%d", i+startRow), o.Id)
	}
	f.SetCellValue(DsSheet, fmt.Sprintf("K%d", headerRow), "Cost Center")
	f.SetCellValue(DsSheet, fmt.Sprintf("L%d", headerRow), "Cost Center Id")
	for i, o := range enumMap["cost_center"] {
		f.SetCellValue(DsSheet, fmt.Sprintf("K%d", i+startRow), o.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("L%d", i+startRow), o.Id)
	}
	f.SetCellValue(DsSheet, fmt.Sprintf("M%d", headerRow), "GL Number")
	f.SetCellValue(DsSheet, fmt.Sprintf("N%d", headerRow), "GL Number Id")
	for i, o := range enumMap["gl_number"] {
		f.SetCellValue(DsSheet, fmt.Sprintf("M%d", i+startRow), o.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("N%d", i+startRow), o.Id)
	}
	// enum [1]-end

	f.SetCellValue(DsSheet, fmt.Sprintf("O%d", headerRow), "Job Description")
	f.SetCellValue(DsSheet, fmt.Sprintf("P%d", headerRow), "Job Description Id")
	for i, o := range jobDescriptions {
		f.SetCellValue(DsSheet, fmt.Sprintf("O%d", i+startRow), o.JobDescription.String)
		f.SetCellValue(DsSheet, fmt.Sprintf("P%d", i+startRow), o.ID.String())
	}
	// --- Populate Section Data ---
	f.SetCellValue(DsSheet, fmt.Sprintf("Q%d", headerRow), "Section")
	f.SetCellValue(DsSheet, fmt.Sprintf("R%d", headerRow), "Section Id")
	f.SetCellValue(DsSheet, fmt.Sprintf("S%d", headerRow), "Section Departmnent Detail Id")
	for i, o := range sections {
		f.SetCellValue(DsSheet, fmt.Sprintf("Q%d", i+startRow), o.SectionName)
		f.SetCellValue(DsSheet, fmt.Sprintf("R%d", i+startRow), o.ID.String())
		f.SetCellValue(DsSheet, fmt.Sprintf("S%d", i+startRow), o.DepartmentDetailID.String())
	}

	f.SetCellValue(DsSheet, fmt.Sprintf("AN%d", headerRow), "Department Detail")
	f.SetCellValue(DsSheet, fmt.Sprintf("AO%d", headerRow), "Department Detail Id")
	f.SetCellValue(DsSheet, fmt.Sprintf("AP%d", headerRow), "Department Id")
	for i, d := range departmentDetails {
		f.SetCellValue(DsSheet, fmt.Sprintf("AN%d", i+startRow), d.DepartmentDetailName)
		f.SetCellValue(DsSheet, fmt.Sprintf("AO%d", i+startRow), d.ID.String())
		f.SetCellValue(DsSheet, fmt.Sprintf("AP%d", i+startRow), d.DepartmentID)
	}

	// enum [2]-start
	f.SetCellValue(DsSheet, fmt.Sprintf("U%d", headerRow), "Status Marriage")
	f.SetCellValue(DsSheet, fmt.Sprintf("V%d", headerRow), "Status Marriage Id")
	for i, d := range enumMap["status_marriage"] {
		f.SetCellValue(DsSheet, fmt.Sprintf("U%d", i+startRow), d.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("V%d", i+startRow), d.Id)
	}

	f.SetCellValue(DsSheet, fmt.Sprintf("W%d", headerRow), "Gender")
	f.SetCellValue(DsSheet, fmt.Sprintf("X%d", headerRow), "Gender Id")
	for i, d := range enumMap["gender"] {
		f.SetCellValue(DsSheet, fmt.Sprintf("W%d", i+startRow), d.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("X%d", i+startRow), d.Id)
	}

	f.SetCellValue(DsSheet, fmt.Sprintf("Y%d", headerRow), "Blood Type")
	f.SetCellValue(DsSheet, fmt.Sprintf("Z%d", headerRow), "Blood Type Id")
	for i, d := range enumMap["blood_type"] {
		f.SetCellValue(DsSheet, fmt.Sprintf("Y%d", i+startRow), d.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("Z%d", i+startRow), d.Id)
	}

	f.SetCellValue(DsSheet, fmt.Sprintf("AA%d", headerRow), "Religion")
	f.SetCellValue(DsSheet, fmt.Sprintf("AB%d", headerRow), "Religion Id")
	for i, d := range enumMap["religion"] {
		f.SetCellValue(DsSheet, fmt.Sprintf("AA%d", i+startRow), d.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("AB%d", i+startRow), d.Id)
	}

	f.SetCellValue(DsSheet, fmt.Sprintf("AQ%d", headerRow), "Status PTKP")
	f.SetCellValue(DsSheet, fmt.Sprintf("AR%d", headerRow), "Status PTKP Id")
	for i, d := range enumMap["status_ptkp"] {
		f.SetCellValue(DsSheet, fmt.Sprintf("AQ%d", i+startRow), d.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("AR%d", i+startRow), d.Id)
	}
	f.SetCellValue(DsSheet, fmt.Sprintf("AC%d", headerRow), "Kelurahan Id")
	f.SetCellValue(DsSheet, fmt.Sprintf("AD%d", headerRow), "Kelurahan")
	f.SetCellValue(DsSheet, fmt.Sprintf("AE%d", headerRow), "Kelurahan Kecamatan Id")
	f.SetCellValue(DsSheet, fmt.Sprintf("AF%d", headerRow), "Kecamatan")
	f.SetCellValue(DsSheet, fmt.Sprintf("AG%d", headerRow), "Kecamatan Id")
	f.SetCellValue(DsSheet, fmt.Sprintf("AH%d", headerRow), "Kecamatan City Id")
	f.SetCellValue(DsSheet, fmt.Sprintf("AI%d", headerRow), "City")
	f.SetCellValue(DsSheet, fmt.Sprintf("AJ%d", headerRow), "City Id")
	f.SetCellValue(DsSheet, fmt.Sprintf("AK%d", headerRow), "City Province Id")
	f.SetCellValue(DsSheet, fmt.Sprintf("AL%d", headerRow), "Province")
	f.SetCellValue(DsSheet, fmt.Sprintf("AM%d", headerRow), "Province Id")

	var kecamatans = make(map[int32]bool)
	var cities = make(map[int32]bool)
	var provinces = make(map[int32]bool)

	type cityInfo struct {
		Name         string
		ID           int32
		ProvinceID   int32
		ProvinceName string
	}
	type provinceInfo struct {
		Name string
		ID   int32
	}
	type kecamatanInfo struct {
		Name     string
		ID       int32
		CityID   int32
		CityName string
	}
	type kelurahanInfo struct {
		Name          string
		ID            int32
		KecamatanID   int32
		KecamatanName string
	}

	var cityListSlice []cityInfo
	var provinceListSlice []provinceInfo
	var kecamatanListSlice []kecamatanInfo
	var kelurahanListSlice []kelurahanInfo

	// Flatten the hierarchical location data (Kelurahan -> Kecamatan -> City -> Province)
	// into distinct slices for population in the datasource sheet.
	for _, d := range kelurahans {
		kelurahanListSlice = append(kelurahanListSlice, kelurahanInfo{
			Name:          d.Name,
			ID:            d.ID,
			KecamatanID:   d.IDKecamatan,
			KecamatanName: d.NameKecamatan,
		})

		if !kecamatans[d.IDKecamatan] {
			kecamatans[d.IDKecamatan] = true
			kecamatanListSlice = append(kecamatanListSlice, kecamatanInfo{
				Name:     d.NameKecamatan,
				ID:       d.IDKecamatan,
				CityID:   d.IDCity,
				CityName: d.NameCity,
			})
		}

		if !cities[d.IDCity] {
			cities[d.IDCity] = true
			cityListSlice = append(cityListSlice, cityInfo{
				Name:         d.NameCity,
				ID:           d.IDCity,
				ProvinceID:   d.IDProvince,
				ProvinceName: d.NameProvince,
			})
		}

		if !provinces[d.IDProvince] {
			provinces[d.IDProvince] = true
			provinceListSlice = append(provinceListSlice, provinceInfo{Name: d.NameProvince, ID: d.IDProvince})
		}
	}

	// Sort location slices to ensure correct grouping for dependent dropdowns.
	sort.Slice(kelurahanListSlice, func(i, j int) bool {
		if kelurahanListSlice[i].KecamatanName == kelurahanListSlice[j].KecamatanName {
			return kelurahanListSlice[i].Name < kelurahanListSlice[j].Name
		}
		return kelurahanListSlice[i].KecamatanName < kelurahanListSlice[j].KecamatanName
	})

	sort.Slice(kecamatanListSlice, func(i, j int) bool {
		if kecamatanListSlice[i].CityName == kecamatanListSlice[j].CityName {
			return kecamatanListSlice[i].Name < kecamatanListSlice[j].Name
		}
		return kecamatanListSlice[i].CityName < kecamatanListSlice[j].CityName
	})

	sort.Slice(cityListSlice, func(i, j int) bool {
		if cityListSlice[i].ProvinceName == cityListSlice[j].ProvinceName {
			return cityListSlice[i].Name < cityListSlice[j].Name
		}
		return cityListSlice[i].ProvinceName < cityListSlice[j].ProvinceName
	})

	sort.Slice(provinceListSlice, func(i, j int) bool {
		return provinceListSlice[i].Name < provinceListSlice[j].Name
	})

	for i, k := range kelurahanListSlice {
		row := i + startRow
		f.SetCellValue(DsSheet, fmt.Sprintf("AC%d", row), k.ID)
		f.SetCellValue(DsSheet, fmt.Sprintf("AD%d", row), k.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("AE%d", row), k.KecamatanID)
	}

	for i, k := range kecamatanListSlice {
		row := i + startRow
		f.SetCellValue(DsSheet, fmt.Sprintf("AF%d", row), k.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("AG%d", row), k.ID)
		f.SetCellValue(DsSheet, fmt.Sprintf("AH%d", row), k.CityID)
	}

	for i, c := range cityListSlice {
		row := i + startRow
		f.SetCellValue(DsSheet, fmt.Sprintf("AI%d", row), c.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("AJ%d", row), c.ID)
		f.SetCellValue(DsSheet, fmt.Sprintf("AK%d", row), c.ProvinceID)
	}

	for i, p := range provinceListSlice {
		row := i + startRow
		f.SetCellValue(DsSheet, fmt.Sprintf("AL%d", row), p.Name)
		f.SetCellValue(DsSheet, fmt.Sprintf("AM%d", row), p.ID)
	}

	// enum [2]-end
	// --- Define Named Ranges ---
	// Named ranges are crucial for Data Validation dropdowns, referencing the populated data.
	// ======================================================================================================================================= #region 3.1 map dropdown list

	departmentsLen := len(departments)

	departmentList := "departmentList"
	if departmentsLen > 0 {

		deptEnd := len(departments) + 1
		departmentCellList := fmt.Sprintf("%s!$A$%d:$A$%d", DsSheet, startRow, deptEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     departmentList,
			RefersTo: departmentCellList,
		})
	}

	outsourcesLen := len(outsources)
	outsourceList := "outsourceList"
	if outsourcesLen > 0 {

		outEnd := len(outsources) + 1
		outsourceCellList := fmt.Sprintf("%s!$C$%d:$C$%d", DsSheet, startRow, outEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     outsourceList,
			RefersTo: outsourceCellList,
		})
	}

	jabatanList := "jabatanList"
	if len(jabatans) > 0 {

		jabatanEnd := len(jabatans) + 1
		jabatanCellList := fmt.Sprintf("%s!$E$%d:$E$%d", DsSheet, startRow, jabatanEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     jabatanList,
			RefersTo: jabatanCellList,
		})
	}

	enumMapgradeLen := len(enumMap["grade"])
	gradeList := "gradeList"
	if enumMapgradeLen > 0 {

		gradeEnd := len(enumMap["grade"]) + 1
		gradeCellList := fmt.Sprintf("%s!$G$%d:$G$%d", DsSheet, startRow, gradeEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     gradeList,
			RefersTo: gradeCellList,
		})
	}

	enumMapstatus_contractLen := len(enumMap["status_contract"])
	statusContractList := "statusContractList"
	if enumMapstatus_contractLen > 0 {

		statusContractEnd := len(enumMap["status_contract"]) + 1
		statusContractCellList := fmt.Sprintf("%s!$I$%d:$I$%d", DsSheet, startRow, statusContractEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     statusContractList,
			RefersTo: statusContractCellList,
		})
	}

	enumMapcost_centerLen := len(enumMap["cost_center"])
	costCenterList := "costCenterList"
	if enumMapcost_centerLen > 0 {

		costCenterEnd := len(enumMap["cost_center"]) + 1
		costCenterCellList := fmt.Sprintf("%s!$K$%d:$K$%d", DsSheet, startRow, costCenterEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     costCenterList,
			RefersTo: costCenterCellList,
		})
	}

	enumMapgl_numberLen := len(enumMap["gl_number"])
	glNumberList := "glNumberList"
	if enumMapgl_numberLen > 0 {

		glNumberEnd := len(enumMap["gl_number"]) + 1
		glNumberCellList := fmt.Sprintf("%s!$M$%d:$M$%d", DsSheet, startRow, glNumberEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     glNumberList,
			RefersTo: glNumberCellList,
		})
	}

	jobDescriptionsLen := len(jobDescriptions)
	jobDescriptionList := "jobDescriptionList"
	if jobDescriptionsLen > 0 {

		jobDescriptionsEnd := len(jobDescriptions) + 1
		jobDescriptionCellList := fmt.Sprintf("%s!$O$%d:$O$%d", DsSheet, startRow, jobDescriptionsEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     jobDescriptionList,
			RefersTo: jobDescriptionCellList,
		})
	}
	sectionsLen := len(sections)
	sectionList := "sectionList"
	if sectionsLen > 0 {

		sectionsEnd := len(sections) + 1
		sectionCellList := fmt.Sprintf("%s!$Q$%d:$Q$%d", DsSheet, startRow, sectionsEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     sectionList,
			RefersTo: sectionCellList,
		})
	}

	departmentDetailList := "departmentDetailList"
	departmentDetailsLen := len(departmentDetails)
	if departmentDetailsLen > 0 {
		deptDetailEnd := departmentDetailsLen + 1
		departmentDetailCellList := fmt.Sprintf("%s!$AN$%d:$AN$%d", DsSheet, startRow, deptDetailEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     departmentDetailList,
			RefersTo: departmentDetailCellList,
		})
	}

	statusMarriageList := "statusMarriageList"
	if len(enumMap["status_marriage"]) > 0 {
		statusMarriageEnd := len(enumMap["status_marriage"]) + 1
		statusMarriageCellList := fmt.Sprintf("%s!$U$%d:$U$%d", DsSheet, startRow, statusMarriageEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     statusMarriageList,
			RefersTo: statusMarriageCellList,
		})
	}

	enumMapgenderLen := len(enumMap["gender"])
	genderList := "genderList"
	if enumMapgenderLen > 0 {

		genderEnd := enumMapgenderLen + 1
		genderCellList := fmt.Sprintf("%s!$W$%d:$W$%d", DsSheet, startRow, genderEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     genderList,
			RefersTo: genderCellList,
		})
	}

	bloodTypeList := "bloodTypeList"
	if len(enumMap["blood_type"]) > 0 {

		bloodTypeListnd := len(enumMap["blood_type"]) + 1
		bloodTypeCellList := fmt.Sprintf("%s!$Y$%d:$Y$%d", DsSheet, startRow, bloodTypeListnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     bloodTypeList,
			RefersTo: bloodTypeCellList,
		})
	}

	religionList := "religionList"
	if len(enumMap["religion"]) > 0 {

		religionListEnd := len(enumMap["religion"]) + 1
		religionCellList := fmt.Sprintf("%s!$AA$%d:$AA$%d", DsSheet, startRow, religionListEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     religionList,
			RefersTo: religionCellList,
		})
	}

	kelurahanList := "kelurahanList"
	kelurahanListEnd := len(kelurahans) + 1
	if len(kelurahans) > 0 {
		kelurahanListCellList := fmt.Sprintf("%s!$AD$%d:$AD$%d", DsSheet, startRow, kelurahanListEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     kelurahanList,
			RefersTo: kelurahanListCellList,
		})
	}

	kecamatanList := "kecamatanList"
	if len(kecamatanListSlice) > 0 {
		kecamatanListEnd := len(kecamatanListSlice) + 1
		kecamatanListCellList := fmt.Sprintf("%s!$AF$%d:$AF$%d", DsSheet, startRow, kecamatanListEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     kecamatanList,
			RefersTo: kecamatanListCellList,
		})
	}
	cityList := "cityList"
	if len(cities) > 0 {
		cityListEnd := len(cities) + 1
		cityListCellList := fmt.Sprintf("%s!$AI$%d:$AI$%d", DsSheet, startRow, cityListEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     cityList,
			RefersTo: cityListCellList,
		})
	}
	provinceList := "provinceList"
	if len(provinces) > 0 {
		provinceListEnd := len(provinceListSlice) + 1
		provinceListCellList := fmt.Sprintf("%s!$AL$%d:$AL$%d", DsSheet, startRow, provinceListEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     provinceList,
			RefersTo: provinceListCellList,
		})
	}
	statusPTKPList := "statusPTKPList"
	if len(enumMap["status_ptkp"]) > 0 {
		statusPTKPEnd := len(enumMap["status_ptkp"]) + 1
		statusPTKPCellList := fmt.Sprintf("%s!$AQ$%d:$AQ$%d", DsSheet, startRow, statusPTKPEnd)

		f.SetDefinedName(&excelize.DefinedName{
			Name:     statusPTKPList,
			RefersTo: statusPTKPCellList,
		})
	}

	// #endregion

	// ======================================================================================================================================= #region 4. Main Sheet

	// Create the main template sheet where employees will update data.
	// 'MainSheet' is the visible sheet for the user.
	// 'CopySheet' is a hidden sheet containing the original data, used for conditional formatting to highlight changes.
	f.NewSheet(MainSheet)
	f.NewSheet(CopySheet)

	colWidths := make(map[string]int)
	writeCell := func(col string, row int, val string) {
		cell := fmt.Sprintf("%s%d", col, row)
		f.SetCellValue(MainSheet, cell, val)
		f.SetCellValue(CopySheet, cell, val)
		if len(val) > colWidths[col] {
			colWidths[col] = len(val)
		}
	}

	// --- Headers ---
	// columns := []string{"Employee NIK", "Employee Name", "Department", "Outsource", "Jabatan", "Grade", "Status Contract", "Cost Center", "GL Number", "Join Date", "Mine Permit Date", "Contract Start Date", "Contract End Date", "Section", "Job Description", "Standard WorkDay", "Standard Overtime"}

	// Set header row values.
	i := 0
	for _, val := range listColumnsStaticOrder {
		if value, exist := columnHandle[val]; exist {
			i++
			colName, _ := excelize.ColumnNumberToName(i)
			writeCell(colName, 1, value.Title)
		}
	}
	// --- Row Data ---
	// Populate the template with existing employee data for update purposes.
	departmentColNumber, departmentDetailColNumber, outsourceColNumber, jabatanColNumber, gradeColNumber,
		statusContractColNumber, costCenterColNumber, glNumberColNumber, jobDescriptionColNumber, sectionColNumber,
		statusMarriageColNumber, statusPTKPColNumber, genderColNumber, bloodTypeColNumber, religionColNumber,
		kelurahanColNumber, kecamatanColNumber, cityColNumber, provinceColNumber,
		domisiliKelurahanColNumber, domisiliKecamatanColNumber, domisiliCityColNumber, domisiliProvinceColNumber := 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0

	for r, emp := range employees {
		rowNum := r + startRow
		// Convert time fields to string date format before writing to Excel
		colStart := 1

		// Write employee data to respective columns
		// <IMPORTANT!!!> ================================================= this column order should same as "listColumnsStaticOrder"
		writeCell("A", rowNum, emp.Nik)
		writeCell("B", rowNum, emp.EmployeeName)

		if _, exist := columnHandle["gender_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.Gender.String)
			genderColNumber = colStart
			colStart++
		}
		if _, exist := columnHandle["birthdate"]; exist {
			birthDate := utils.ConvTimeToStringDate(emp.Birthdate)
			writeCell(listColHeader[colStart], rowNum, birthDate)
			colStart++
		}

		if _, exist := columnHandle["birthplace"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.Birthplace.String)
			colStart++
		}

		if _, exist := columnHandle["religion_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.Religion.String)
			religionColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["status_marriage_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.StatusMarriage.String)
			statusMarriageColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["status_ptkp_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.StatusPtkp.String)
			statusPTKPColNumber = colStart

			colStart++
		}

		if _, exist := columnHandle["status_contract_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.StatusContract.String)
			statusContractColNumber = colStart

			colStart++
		}
		if _, exist := columnHandle["blood_type_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.BloodType.String)
			bloodTypeColNumber = colStart
			colStart++
		}
		if _, exist := columnHandle["npwp"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.Npwp.String)
			colStart++
		}

		if _, exist := columnHandle["bpjstk_number"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.BpjstkNumber.String)
			colStart++
		}

		if _, exist := columnHandle["bpjskes_number"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.BpjskesNumber.String)
			colStart++
		}
		if _, exist := columnHandle["email"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.Email)
			colStart++
		}

		if _, exist := columnHandle["email_office"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.EmailOffice.String)
			colStart++
		}

		if _, exist := columnHandle["phonenumber"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.Phonenumber.String)
			colStart++
		}

		if _, exist := columnHandle["domisili_province_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.DomisiliProvinceName.String)
			domisiliProvinceColNumber = colStart
			colStart++
		}
		if _, exist := columnHandle["domisili_city_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.DomisiliCityName.String)
			domisiliCityColNumber = colStart
			colStart++
		}
		if _, exist := columnHandle["domisili_kecamatan_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.DomisiliKecamatanName.String)
			domisiliKecamatanColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["domisili_kelurahan_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.DomisiliKelurahanName.String)
			domisiliKelurahanColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["province_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.ProvinceName.String)
			provinceColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["city_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.CityName.String)
			cityColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["kecamatan_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.KecamatanName.String)
			kecamatanColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["kelurahan_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.KelurahanName.String)
			kelurahanColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["jabatan_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.Jabatan.String)
			jabatanColNumber = colStart

			colStart++
		}

		if _, exist := columnHandle["grade_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.Grade.String)
			gradeColNumber = colStart

			colStart++
		}
		if _, exist := columnHandle["join_date"]; exist {
			joinDate := utils.ConvTimeToStringDate(emp.JoinDate.Time)
			writeCell(listColHeader[colStart], rowNum, joinDate)
			colStart++
		}
		if _, exist := columnHandle["mine_permit_date"]; exist {
			minePermitDate := utils.ConvTimeToStringDate(emp.MinePermitDate.Time)
			writeCell(listColHeader[colStart], rowNum, minePermitDate)
			colStart++
		}
		if _, exist := columnHandle["outsource_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.OutsourceName.String)
			outsourceColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["department_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.DepartmentName.String)
			departmentColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["department_detail_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.DepartmentDetailName.String)
			departmentDetailColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["section_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.SectionName.String)
			sectionColNumber = colStart
			colStart++
		}

		if _, exist := columnHandle["job_description_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.JobDescriptionTitle.String)
			jobDescriptionColNumber = colStart
			colStart++
		}
		if _, exist := columnHandle["cost_center_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.CostCenter.String)
			costCenterColNumber = colStart

			colStart++
		}

		if _, exist := columnHandle["gl_number_id"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.GlNumber.String)
			glNumberColNumber = colStart

			colStart++
		}

		if _, exist := columnHandle["standard_workday"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.StandardWorkday.String)
			colStart++
		}

		if _, exist := columnHandle["standard_overtime"]; exist {
			writeCell(listColHeader[colStart], rowNum, emp.StandardOvertime.String)
			colStart++
		}

	}

	for col, width := range colWidths {
		f.SetColWidth(MainSheet, col, col, float64(width)+5)
		f.SetColWidth(CopySheet, col, col, float64(width)+5)
	}

	lastRow := len(employees) + 1

	// --- Conditional Formatting & Data Validation ---
	yellowStyle, err := f.NewConditionalStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}
	lastCol := listColHeader[len(listColHeader)]
	cellRange := fmt.Sprintf("A%d:%s%d", startRow, lastCol, lastRow)

	// Apply conditional formatting to highlight cells that differ from the original data (CopySheet).
	err = f.SetConditionalFormat(
		MainSheet,
		cellRange,
		[]excelize.ConditionalFormatOptions{
			{
				Type:     "formula",
				Criteria: `=INDIRECT("RC",FALSE)<>INDIRECT("EmployeeCopy!RC",FALSE)`,
				Format:   &yellowStyle,
			},
		},
	)
	if err != nil {
		return err
	}

	// department dropdown
	if _, exist := listColHeader[departmentColNumber]; exist {

		dvDept := excelize.NewDataValidation(true)
		departmentCell := fmt.Sprintf("%s%d:%s%d", listColHeader[departmentColNumber], startRow, listColHeader[departmentColNumber], lastRow)
		dvDept.Sqref = departmentCell

		dvDept.SetSqrefDropList("=" + departmentList)
		f.AddDataValidation(MainSheet, dvDept)
	}

	if _, exist := listColHeader[departmentDetailColNumber]; exist {

		dvDeptDetail := excelize.NewDataValidation(true)
		departmentDetailCell := fmt.Sprintf("%s%d:%s%d", listColHeader[departmentDetailColNumber], startRow, listColHeader[departmentDetailColNumber], lastRow)
		dvDeptDetail.Sqref = departmentDetailCell

		if departmentColNumber > 0 {
			deptColLetter := listColHeader[departmentColNumber]
			deptListEnd := len(departments) + 1
			deptDetailListEnd := len(departmentDetails) + 1

			// Dependent Dropdown Formula:
			// Uses OFFSET, MATCH, and COUNTIF to dynamically select the range of Department Details
			// that correspond to the selected Department in the previous column.
			// Formula: =OFFSET(StartCell, MATCH(ParentValue, ParentRange, 0)-1, 0, COUNTIF(ParentRange, ParentValue), 1)
			formula := fmt.Sprintf("=OFFSET('%s'!$AN$2,MATCH(INDEX('%s'!$B$2:$B$%d,MATCH(%s2,'%s'!$A$2:$A$%d,0)),'%s'!$AP$2:$AP$%d,0)-1,0,COUNTIF('%s'!$AP$2:$AP$%d,INDEX('%s'!$B$2:$B$%d,MATCH(%s2,'%s'!$A$2:$A$%d,0))),1)",
				DsSheet,
				DsSheet, deptListEnd,
				deptColLetter, DsSheet, deptListEnd,
				DsSheet, deptDetailListEnd,
				DsSheet, deptDetailListEnd,
				DsSheet, deptListEnd,
				deptColLetter, DsSheet, deptListEnd)
			dvDeptDetail.SetSqrefDropList(formula)
		} else {
			dvDeptDetail.SetSqrefDropList("=" + departmentDetailList)
		}
		f.AddDataValidation(MainSheet, dvDeptDetail)
	}

	// outsource dropdown
	if _, exist := listColHeader[outsourceColNumber]; exist {

		dvOut := excelize.NewDataValidation(true)
		outsourceCell := fmt.Sprintf("%s%d:%s%d", listColHeader[outsourceColNumber], startRow, listColHeader[outsourceColNumber], lastRow)
		dvOut.Sqref = outsourceCell

		dvOut.SetSqrefDropList("=" + outsourceList)
		f.AddDataValidation(MainSheet, dvOut)
	}
	// jabatan dropdown
	if _, exist := listColHeader[jabatanColNumber]; exist {

		jabatanOut := excelize.NewDataValidation(true)
		jabatanCell := fmt.Sprintf("%s%d:%s%d", listColHeader[jabatanColNumber], startRow, listColHeader[jabatanColNumber], lastRow)
		jabatanOut.Sqref = jabatanCell

		jabatanOut.SetSqrefDropList("=" + jabatanList)
		f.AddDataValidation(MainSheet, jabatanOut)
	}
	// gender dropdown
	if _, exist := listColHeader[genderColNumber]; exist {

		genderOut := excelize.NewDataValidation(true)
		genderCell := fmt.Sprintf("%s%d:%s%d", listColHeader[genderColNumber], startRow, listColHeader[genderColNumber], lastRow)
		genderOut.Sqref = genderCell

		genderOut.SetSqrefDropList("=" + genderList)
		f.AddDataValidation(MainSheet, genderOut)
	}
	// grade dropdown
	if _, exist := listColHeader[gradeColNumber]; exist {

		gradeOut := excelize.NewDataValidation(true)
		gradeCell := fmt.Sprintf("%s%d:%s%d", listColHeader[gradeColNumber], startRow, listColHeader[gradeColNumber], lastRow)
		gradeOut.Sqref = gradeCell

		gradeOut.SetSqrefDropList("=" + gradeList)
		f.AddDataValidation(MainSheet, gradeOut)
	}
	// statusContract dropdown
	if _, exist := listColHeader[statusContractColNumber]; exist {

		statusContractOut := excelize.NewDataValidation(true)
		statusContractCell := fmt.Sprintf("%s%d:%s%d", listColHeader[statusContractColNumber], startRow, listColHeader[statusContractColNumber], lastRow)
		statusContractOut.Sqref = statusContractCell

		statusContractOut.SetSqrefDropList("=" + statusContractList)
		f.AddDataValidation(MainSheet, statusContractOut)
	}
	// costCenter dropdown
	if _, exist := listColHeader[costCenterColNumber]; exist {

		costCenterOut := excelize.NewDataValidation(true)
		costCenterCell := fmt.Sprintf("%s%d:%s%d", listColHeader[costCenterColNumber], startRow, listColHeader[costCenterColNumber], lastRow)
		costCenterOut.Sqref = costCenterCell

		costCenterOut.SetSqrefDropList("=" + costCenterList)
		f.AddDataValidation(MainSheet, costCenterOut)
	}
	// glNumber dropdown
	if _, exist := listColHeader[glNumberColNumber]; exist {

		glNumberOut := excelize.NewDataValidation(true)
		glNumberCell := fmt.Sprintf("%s%d:%s%d", listColHeader[glNumberColNumber], startRow, listColHeader[glNumberColNumber], lastRow)
		glNumberOut.Sqref = glNumberCell

		glNumberOut.SetSqrefDropList("=" + glNumberList)
		f.AddDataValidation(MainSheet, glNumberOut)
	}
	// glNumber dropdown
	if _, exist := listColHeader[sectionColNumber]; exist {

		sectionsOut := excelize.NewDataValidation(true)
		sectionCell := fmt.Sprintf("%s%d:%s%d", listColHeader[sectionColNumber], startRow, listColHeader[sectionColNumber], lastRow)
		sectionsOut.Sqref = sectionCell

		if departmentDetailColNumber > 0 {
			deptDetailColLetter := listColHeader[departmentDetailColNumber]
			deptDetailListEnd := len(departmentDetails) + 1
			sectionListEnd := len(sections) + 1

			// Dependent Dropdown for Section based on Department Detail
			formula := fmt.Sprintf("=OFFSET('%s'!$Q$2,MATCH(INDEX('%s'!$AO$2:$AO$%d,MATCH(%s2,'%s'!$AN$2:$AN$%d,0)),'%s'!$S$2:$S$%d,0)-1,0,COUNTIF('%s'!$S$2:$S$%d,INDEX('%s'!$AO$2:$AO$%d,MATCH(%s2,'%s'!$AN$2:$AN$%d,0))),1)",
				DsSheet,
				DsSheet, deptDetailListEnd,
				deptDetailColLetter, DsSheet, deptDetailListEnd,
				DsSheet, sectionListEnd,
				DsSheet, sectionListEnd,
				DsSheet, deptDetailListEnd,
				deptDetailColLetter, DsSheet, deptDetailListEnd)
			sectionsOut.SetSqrefDropList(formula)
		} else {
			sectionsOut.SetSqrefDropList("=" + sectionList)
		}
		f.AddDataValidation(MainSheet, sectionsOut)
	}

	// jobDescription dropdown
	if _, exist := listColHeader[jobDescriptionColNumber]; exist {

		jobDescriptionsOut := excelize.NewDataValidation(true)
		jobDescriptionCell := fmt.Sprintf("%s%d:%s%d", listColHeader[jobDescriptionColNumber], startRow, listColHeader[jobDescriptionColNumber], lastRow)
		jobDescriptionsOut.Sqref = jobDescriptionCell

		jobDescriptionsOut.SetSqrefDropList("=" + jobDescriptionList)
		f.AddDataValidation(MainSheet, jobDescriptionsOut)
	}

	// statusMarriage dropdown
	if _, exist := listColHeader[statusMarriageColNumber]; exist {

		statusMarriageOut := excelize.NewDataValidation(true)
		statusMarriageCell := fmt.Sprintf("%s%d:%s%d", listColHeader[statusMarriageColNumber], startRow, listColHeader[statusMarriageColNumber], lastRow)
		statusMarriageOut.Sqref = statusMarriageCell

		statusMarriageOut.SetSqrefDropList("=" + statusMarriageList)
		f.AddDataValidation(MainSheet, statusMarriageOut)
	}
	// statusPTKP dropdown
	if _, exist := listColHeader[statusPTKPColNumber]; exist {

		statusPTKPOut := excelize.NewDataValidation(true)
		statusPTKPCell := fmt.Sprintf("%s%d:%s%d", listColHeader[statusPTKPColNumber], startRow, listColHeader[statusPTKPColNumber], lastRow)
		statusPTKPOut.Sqref = statusPTKPCell

		statusPTKPOut.SetSqrefDropList("=" + statusPTKPList)
		f.AddDataValidation(MainSheet, statusPTKPOut)
	}

	// gender dropdown
	if _, exist := listColHeader[genderColNumber]; exist {

		genderOut := excelize.NewDataValidation(true)
		genderCell := fmt.Sprintf("%s%d:%s%d", listColHeader[genderColNumber], startRow, listColHeader[genderColNumber], lastRow)
		genderOut.Sqref = genderCell

		genderOut.SetSqrefDropList("=" + genderList)
		f.AddDataValidation(MainSheet, genderOut)
	}
	// gender dropdown
	if _, exist := listColHeader[bloodTypeColNumber]; exist {

		bloodTypeOut := excelize.NewDataValidation(true)
		bloodTypeCell := fmt.Sprintf("%s%d:%s%d", listColHeader[bloodTypeColNumber], startRow, listColHeader[bloodTypeColNumber], lastRow)
		bloodTypeOut.Sqref = bloodTypeCell

		bloodTypeOut.SetSqrefDropList("=" + bloodTypeList)
		f.AddDataValidation(MainSheet, bloodTypeOut)
	}

	// religion dropdown
	if _, exist := listColHeader[religionColNumber]; exist {

		religionOut := excelize.NewDataValidation(true)
		religionCell := fmt.Sprintf("%s%d:%s%d", listColHeader[religionColNumber], startRow, listColHeader[religionColNumber], lastRow)
		religionOut.Sqref = religionCell

		religionOut.SetSqrefDropList("=" + religionList)
		f.AddDataValidation(MainSheet, religionOut)
	}

	// province dropdown
	if _, exist := listColHeader[provinceColNumber]; exist {

		provinceOut := excelize.NewDataValidation(true)
		provinceCell := fmt.Sprintf("%s%d:%s%d", listColHeader[provinceColNumber], startRow, listColHeader[provinceColNumber], lastRow)
		provinceOut.Sqref = provinceCell

		provinceOut.SetSqrefDropList("=" + provinceList)
		f.AddDataValidation(MainSheet, provinceOut)
	}

	// city dropdown
	if _, exist := listColHeader[cityColNumber]; exist {

		cityOut := excelize.NewDataValidation(true)
		cityCell := fmt.Sprintf("%s%d:%s%d", listColHeader[cityColNumber], startRow, listColHeader[cityColNumber], lastRow)
		cityOut.Sqref = cityCell

		if provinceColNumber > 0 {
			provColLetter := listColHeader[provinceColNumber]
			cityListEnd := len(cityListSlice) + 1
			provinceListEnd := len(provinceListSlice) + 1
			// Dependent Dropdown for City based on Province
			formula := fmt.Sprintf("=OFFSET('%s'!$AI$2,MATCH(INDEX('%s'!$AM$2:$AM$%d,MATCH(%s2,'%s'!$AL$2:$AL$%d,0)),'%s'!$AK$2:$AK$%d,0)-1,0,COUNTIF('%s'!$AK$2:$AK$%d,INDEX('%s'!$AM$2:$AM$%d,MATCH(%s2,'%s'!$AL$2:$AL$%d,0))),1)",
				DsSheet,
				DsSheet, provinceListEnd,
				provColLetter, DsSheet, provinceListEnd,
				DsSheet, cityListEnd,
				DsSheet, cityListEnd,
				DsSheet, provinceListEnd,
				provColLetter, DsSheet, provinceListEnd)
			cityOut.SetSqrefDropList(formula)
		} else {
			cityOut.SetSqrefDropList("=" + cityList)
		}
		f.AddDataValidation(MainSheet, cityOut)
	}

	// kecamatan dropdown
	if _, exist := listColHeader[kecamatanColNumber]; exist {

		kecamatanOut := excelize.NewDataValidation(true)
		kecamatanCell := fmt.Sprintf("%s%d:%s%d", listColHeader[kecamatanColNumber], startRow, listColHeader[kecamatanColNumber], lastRow)
		kecamatanOut.Sqref = kecamatanCell

		if cityColNumber > 0 {
			cityColLetter := listColHeader[cityColNumber]
			kecamatanListEnd := len(kecamatanListSlice) + 1
			cityListEnd := len(cityListSlice) + 1
			// Dependent Dropdown for Kecamatan based on City
			formula := fmt.Sprintf("=OFFSET('%s'!$AF$2,MATCH(INDEX('%s'!$AJ$2:$AJ$%d,MATCH(%s2,'%s'!$AI$2:$AI$%d,0)),'%s'!$AH$2:$AH$%d,0)-1,0,COUNTIF('%s'!$AH$2:$AH$%d,INDEX('%s'!$AJ$2:$AJ$%d,MATCH(%s2,'%s'!$AI$2:$AI$%d,0))),1)",
				DsSheet,
				DsSheet, cityListEnd,
				cityColLetter, DsSheet, cityListEnd,
				DsSheet, kecamatanListEnd,
				DsSheet, kecamatanListEnd,
				DsSheet, cityListEnd,
				cityColLetter, DsSheet, cityListEnd)
			kecamatanOut.SetSqrefDropList(formula)
		} else {
			kecamatanOut.SetSqrefDropList("=" + kecamatanList)
		}
		f.AddDataValidation(MainSheet, kecamatanOut)
	}

	// kelurahan dropdown
	if _, exist := listColHeader[kelurahanColNumber]; exist {
		kelurahanOut := excelize.NewDataValidation(true)
		kelurahanCell := fmt.Sprintf("%s%d:%s%d", listColHeader[kelurahanColNumber], startRow, listColHeader[kelurahanColNumber], lastRow)

		kelurahanOut.Sqref = kelurahanCell

		if kecamatanColNumber > 0 {
			kecamatanColLetter := listColHeader[kecamatanColNumber]
			kelurahanListEnd := len(kelurahanListSlice) + 1
			kecamatanListEnd := len(kecamatanListSlice) + 1
			// Dependent Dropdown for Kelurahan based on Kecamatan
			formula := fmt.Sprintf("=OFFSET('%s'!$AD$2,MATCH(INDEX('%s'!$AG$2:$AG$%d,MATCH(%s2,'%s'!$AF$2:$AF$%d,0)),'%s'!$AE$2:$AE$%d,0)-1,0,COUNTIF('%s'!$AE$2:$AE$%d,INDEX('%s'!$AG$2:$AG$%d,MATCH(%s2,'%s'!$AF$2:$AF$%d,0))),1)",
				DsSheet,
				DsSheet, kecamatanListEnd,
				kecamatanColLetter, DsSheet, kecamatanListEnd,
				DsSheet, kelurahanListEnd,
				DsSheet, kelurahanListEnd,
				DsSheet, kecamatanListEnd,
				kecamatanColLetter, DsSheet, kecamatanListEnd)
			kelurahanOut.SetSqrefDropList(formula)
		} else {
			kelurahanOut.SetSqrefDropList("=" + kelurahanList)
		}
		f.AddDataValidation(MainSheet, kelurahanOut)
	}

	// domisiliProvince dropdown
	if _, exist := listColHeader[domisiliProvinceColNumber]; exist {

		domisiliProvinceOut := excelize.NewDataValidation(true)
		domisiliProvinceCell := fmt.Sprintf("%s%d:%s%d", listColHeader[domisiliProvinceColNumber], startRow, listColHeader[domisiliProvinceColNumber], lastRow)
		domisiliProvinceOut.Sqref = domisiliProvinceCell

		domisiliProvinceOut.SetSqrefDropList("=" + provinceList)
		f.AddDataValidation(MainSheet, domisiliProvinceOut)
	}

	// domisiliCity dropdown
	if _, exist := listColHeader[domisiliCityColNumber]; exist {

		domisiliCityOut := excelize.NewDataValidation(true)
		domisiliCityCell := fmt.Sprintf("%s%d:%s%d", listColHeader[domisiliCityColNumber], startRow, listColHeader[domisiliCityColNumber], lastRow)
		domisiliCityOut.Sqref = domisiliCityCell
		if domisiliProvinceColNumber > 0 {
			provColLetter := listColHeader[domisiliProvinceColNumber]
			cityListEnd := len(cityListSlice) + 1
			provinceListEnd := len(provinceListSlice) + 1
			// Dependent Dropdown for Domisili City based on Domisili Province
			formula := fmt.Sprintf("=OFFSET('%s'!$AI$2,MATCH(INDEX('%s'!$AM$2:$AM$%d,MATCH(%s2,'%s'!$AL$2:$AL$%d,0)),'%s'!$AK$2:$AK$%d,0)-1,0,COUNTIF('%s'!$AK$2:$AK$%d,INDEX('%s'!$AM$2:$AM$%d,MATCH(%s2,'%s'!$AL$2:$AL$%d,0))),1)",
				DsSheet,
				DsSheet, provinceListEnd,
				provColLetter, DsSheet, provinceListEnd,
				DsSheet, cityListEnd,
				DsSheet, cityListEnd,
				DsSheet, provinceListEnd,
				provColLetter, DsSheet, provinceListEnd)
			domisiliCityOut.SetSqrefDropList(formula)
		} else {
			domisiliCityOut.SetSqrefDropList("=" + cityList)
		}
		f.AddDataValidation(MainSheet, domisiliCityOut)

	}

	// domisiliKecamatan dropdown
	if _, exist := listColHeader[domisiliKecamatanColNumber]; exist {

		domisiliKecamatanOut := excelize.NewDataValidation(true)
		domisiliKecamatanCell := fmt.Sprintf("%s%d:%s%d", listColHeader[domisiliKecamatanColNumber], startRow, listColHeader[domisiliKecamatanColNumber], lastRow)
		domisiliKecamatanOut.Sqref = domisiliKecamatanCell

		if domisiliCityColNumber > 0 {
			cityColLetter := listColHeader[domisiliCityColNumber]
			kecamatanListEnd := len(kecamatanListSlice) + 1
			cityListEnd := len(cityListSlice) + 1
			// Dependent Dropdown for Domisili Kecamatan based on Domisili City
			formula := fmt.Sprintf("=OFFSET('%s'!$AF$2,MATCH(INDEX('%s'!$AJ$2:$AJ$%d,MATCH(%s2,'%s'!$AI$2:$AI$%d,0)),'%s'!$AH$2:$AH$%d,0)-1,0,COUNTIF('%s'!$AH$2:$AH$%d,INDEX('%s'!$AJ$2:$AJ$%d,MATCH(%s2,'%s'!$AI$2:$AI$%d,0))),1)",
				DsSheet,
				DsSheet, cityListEnd,
				cityColLetter, DsSheet, cityListEnd,
				DsSheet, kecamatanListEnd,
				DsSheet, kecamatanListEnd,
				DsSheet, cityListEnd,
				cityColLetter, DsSheet, cityListEnd)
			domisiliKecamatanOut.SetSqrefDropList(formula)
		} else {
			domisiliKecamatanOut.SetSqrefDropList("=" + kecamatanList)
		}
		f.AddDataValidation(MainSheet, domisiliKecamatanOut)
	}

	// domisiliKelurahan dropdown
	if _, exist := listColHeader[domisiliKelurahanColNumber]; exist {

		domisiliKelurahanOut := excelize.NewDataValidation(true)
		domisiliKelurahanCell := fmt.Sprintf("%s%d:%s%d", listColHeader[domisiliKelurahanColNumber], startRow, listColHeader[domisiliKelurahanColNumber], lastRow)
		domisiliKelurahanOut.Sqref = domisiliKelurahanCell

		if domisiliKecamatanColNumber > 0 {
			kecamatanColLetter := listColHeader[domisiliKecamatanColNumber]
			kelurahanListEnd := len(kelurahanListSlice) + 1
			kecamatanListEnd := len(kecamatanListSlice) + 1
			// Dependent Dropdown for Domisili Kelurahan based on Domisili Kecamatan
			formula := fmt.Sprintf("=OFFSET('%s'!$AD$2,MATCH(INDEX('%s'!$AG$2:$AG$%d,MATCH(%s2,'%s'!$AF$2:$AF$%d,0)),'%s'!$AE$2:$AE$%d,0)-1,0,COUNTIF('%s'!$AE$2:$AE$%d,INDEX('%s'!$AG$2:$AG$%d,MATCH(%s2,'%s'!$AF$2:$AF$%d,0))),1)",
				DsSheet,
				DsSheet, kecamatanListEnd,
				kecamatanColLetter, DsSheet, kecamatanListEnd,
				DsSheet, kelurahanListEnd,
				DsSheet, kelurahanListEnd,
				DsSheet, kecamatanListEnd,
				kecamatanColLetter, DsSheet, kecamatanListEnd)
			domisiliKelurahanOut.SetSqrefDropList(formula)
		} else {
			domisiliKelurahanOut.SetSqrefDropList("=" + kelurahanList)
		}
		f.AddDataValidation(MainSheet, domisiliKelurahanOut)
	}
	// glNumber dropdown
	// 5.5. Protect main sheet with specific cells locked
	// Ensure Column A (Nik) is locked. Unlock columns starting from B (index 2) to last column.

	// Create Locked Style (Explicit)
	lockedStyle, _ := f.NewStyle(&excelize.Style{
		Protection: &excelize.Protection{
			Locked: true,
		},
	})

	// Create Unlocked Style
	unlockedStyle, _ := f.NewStyle(&excelize.Style{
		Protection: &excelize.Protection{
			Locked: false,
		},
	})

	if lastRow >= startRow {
		// Explicitly lock A2:A{LastRow}
		f.SetCellStyle(MainSheet, fmt.Sprintf("A%d", startRow), fmt.Sprintf("A%d", lastRow), lockedStyle)

		if i > 1 {
			lastColName, _ := excelize.ColumnNumberToName(i)
			// Unlock B2 : {LastCol}{LastRow}
			f.SetCellStyle(MainSheet, fmt.Sprintf("B%d", startRow), fmt.Sprintf("%s%d", lastColName, lastRow), unlockedStyle)
		}
	}

	// Protect the MainSheet
	if err := f.ProtectSheet(MainSheet, &excelize.SheetProtectionOptions{
		SelectLockedCells:   true,
		SelectUnlockedCells: true,
	}); err != nil {
		return err
	}

	err = f.DeleteSheet("Sheet1")
	if err != nil {
		fmt.Println("Error deleting sheet:", err)
	}
	// ======================================================================================================================================= #endregion 5. Dropdo

	// Set active sheet to MainSheet
	if idx, err := f.GetSheetIndex(MainSheet); err == nil {
		f.SetActiveSheet(idx)
	}

	// Hide the copy and datasource sheets so users only see the main template.
	f.SetSheetVisible(DsSheet, false)
	f.SetSheetVisible(CopySheet, false)

	// 6. Return
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return err
	}

	filename := fmt.Sprintf("employee_update_template_%s.xlsx", utils.ConvTimeToString(time.Now()))

	res.Data = &pb.ExportResponse{
		FileName:    filename,
		FileContent: buf.Bytes(),
		MimeType:    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	return nil
}
