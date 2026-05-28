package employee

import (
	"fmt"
	"moyo-master-service/pkg/enum"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

/*
// Helper function to map strings to struct fields
func getStructValue(v MasterEmployee, colName string) string {
	switch colName {
	case "employee_name":
		return v.EmployeeName
	case "nik":
		return v.Nik
	case "email":
		return v.Email
	case "department_name":
		return v.DepartmentName.String
	case "outsource_name":
		return v.OutsourceName.String
	case "gender":
		return v.Gender.String
	case "birthdate":
		return utils.ConvTimeToStringDate(v.Birthdate)
	case "birthplace":
		return v.Birthplace.String
	case "religion":
		return v.Religion.String
	case "status_marriage":
		return v.StatusMarriage.String
	case "status_ptkp":
		return v.StatusPtkp.String
	case "gl_number":
		return v.GlNumber.String
	case "status_contract":
		return v.StatusContract.String
	case "blood_type":
		return v.BloodType.String
	case "bpjstk_number":
		return v.BpjstkNumber.String
	case "bpjskes_number":
		return v.BpjskesNumber.String
	case "email_office":
		return v.EmailOffice.String
	case "phonenumber":
		return v.Phonenumber.String
	default:
	}
	return "" // Return empty string if column name doesn't match

}
// */

// handle master_enum where enum type in ('jabatan', 'grade', 'status_contract', 'cost_center', 'gl_number', 'status_marriage',  'blood_type', 'status_ptkp', 'religion')
func HandleEnum(enum []enum.MasterEnum, columnHandle map[string]ColumnModel) map[string][]Option {
	// Initialize the map that will hold the grouped results.
	// The key is the EnumType (string) and the value is a slice of Option structs.
	result := make(map[string][]Option)

	// Iterate through each MasterEnum struct in the input slice.
	for _, item := range enum {

		// 1. Create the Option struct from the MasterEnum fields.
		// We convert the UUID to a string for the 'id' field.
		option := Option{
			Id:   item.ID.String(), // Use the UUID as a string for the ID
			Name: item.EnumValue,   // Use EnumValue as the display name
		}

		// 2. Group the Option by its EnumType.
		// The append function works safely here:
		// - If 'item.EnumType' is not yet a key in 'result', a new slice will be created.
		// - If it already exists, the new 'option' will be appended to the existing slice.
		if (item.EnumType == "gender" && columnHandle["gender_id"].IsCollected) ||
			(item.EnumType == "grade" && columnHandle["grade_id"].IsCollected) ||
			(item.EnumType == "status_contract" && columnHandle["status_contract_id"].IsCollected) ||
			(item.EnumType == "status_ptkp" && columnHandle["status_ptkp_id"].IsCollected) ||
			(item.EnumType == "cost_center" && columnHandle["cost_center_id"].IsCollected) ||
			(item.EnumType == "gl_number" && columnHandle["gl_number_id"].IsCollected) ||
			(item.EnumType == "status_marriage" && columnHandle["status_marriage_id"].IsCollected) ||
			(item.EnumType == "blood_type" && columnHandle["blood_type_id"].IsCollected) ||
			(item.EnumType == "religion" && columnHandle["religion_id"].IsCollected) {
			result[item.EnumType] = append(result[item.EnumType], option)
		}
	}

	return result
}

type ColumnModel struct {
	Title                  string `json:"title"`
	SourceDropdownName     string `json:"source_dropdown_name"`
	SourceDropdownPosition string `json:"source_dropdown_position"`
	ActualColumnPosition   string `json:"actual_column_position"`
	ActualDropdownPosition string `json:"actual_dropdown_position"`
	IsCollected            bool
}

var listColHeader = map[int]string{
	1:  "C",
	2:  "D",
	3:  "E",
	4:  "F",
	5:  "G",
	6:  "H",
	7:  "I",
	8:  "J",
	9:  "K",
	10: "L",
	11: "M",
	12: "N",
	13: "O",
	14: "P",
	15: "Q",
	16: "R",
	17: "S",
	18: "T",
	19: "U",
	20: "V",
	21: "W",
	22: "X",
	23: "Y",
	24: "Z",
	25: "AA",
	26: "AB",
	27: "AC",
	28: "AD",
	29: "AE",
	30: "AF",
	31: "AG",
	32: "AH",
	33: "AI",
	34: "AJ",
}

// |npwp|bpjstk_number|bpjskes_number|email|email_office|phonenumber|domisili_kelurahan_id|domisili_kecamatan_id|domisili_city_id|domisili_province_id|kelurahan_id|kecamatan_id|city_id|province_id|jabatan_id|grade_id|join_date|mine_permit_date|status_ptkp_id|status_contract_id|outsource_id|department_id|department_detail_id|section_id|job_description_id|cost_center_id|gl_number_id|standard_workday|standard_overtime
var listColumnsStaticOrder = []string{
	"nik",
	"employee_name",
	"gender_id",
	"birthdate",
	"birthplace",
	"religion_id",
	"status_marriage_id",
	"status_ptkp_id",
	"status_contract_id",
	"blood_type_id",
	"npwp",
	"bpjstk_number",
	"bpjskes_number",
	"email",
	"email_office",
	"phonenumber",
	"domisili_province_id",
	"domisili_city_id",
	"domisili_kecamatan_id",
	"domisili_kelurahan_id",
	"province_id",
	"city_id",
	"kecamatan_id",
	"kelurahan_id",
	"jabatan_id",
	"grade_id",
	"join_date",
	"mine_permit_date",
	"outsource_id",
	"department_id",
	"department_detail_id",
	"section_id",
	"job_description_id",
	"cost_center_id",
	"gl_number_id",
	"standard_workday",
	"standard_overtime",
}

var listColumnsStatic = map[string]ColumnModel{
	"nik": {
		Title: "Employee NIK",
	},
	"employee_name": {
		Title: "Employee Name",
	},
	"birthdate": {
		Title: "Birthdate",
	},
	"birthplace": {
		Title: "Birthplace",
	},
	"email": {
		Title: "Email",
	},
	"email_office": {
		Title: "Email Office",
	},
	"phonenumber": {
		Title: "Phonenumber",
	},
	"npwp": {
		Title: "NPWP",
	},
	"department_id": {
		Title:              "Division",
		SourceDropdownName: "departmentList",
	},
	"department_detail_id": {
		Title:              "Department",
		SourceDropdownName: "departmentDetailList",
	},
	"outsource_id": {
		Title:              "Outsource",
		SourceDropdownName: "outsourceList",
	},
	"jabatan_id": {
		Title:              "Jabatan",
		SourceDropdownName: "jabatanList",
	},
	"grade_id": {
		Title:              "Grade",
		SourceDropdownName: "gradeList",
	},
	"status_ptkp_id": {
		Title:              "Status PTKP",
		SourceDropdownName: "statusPTKPList",
	},
	"status_contract_id": {
		Title:              "Status Contract",
		SourceDropdownName: "statusContractList",
	},
	"cost_center_id": {
		Title:              "Cost Center",
		SourceDropdownName: "costCenterList",
	},
	"gl_number_id": {
		Title:              "GL Number",
		SourceDropdownName: "glNumberList",
	},
	"join_date": {
		Title: "Join Date",
	},
	"mine_permit_date": {
		Title: "Mine Permit Date",
	},
	"contract_start_date": {
		Title: "Contract Start Date",
	},
	"contract_end_date": {
		Title: "Contract End Date",
	},
	"section_id": {
		Title:              "Section",
		SourceDropdownName: "sectionList",
	},
	"job_description_id": {
		Title:              "Job Description",
		SourceDropdownName: "jobDescriptionList",
	},
	"standard_workday": {
		Title: "Standard WorkDay",
	},
	"standard_overtime": {
		Title: "Standard Overtime",
	},
	"bpjstk_number": {
		Title: "BPJSTK Number",
	},
	"bpjskes_number": {
		Title: "BPJSKES Number",
	},
	"status_marriage_id": {
		Title: "Status Marriage",
	},
	"gender_id": {
		Title: "Gender",
	},
	"blood_type_id": {
		Title: "Blood Type",
	},
	"religion_id": {
		Title: "Religion",
	},
	"kelurahan_id": {
		Title: "Kelurahan",
	},
	"kecamatan_id": {
		Title: "Kecamatan",
	},
	"city_id": {
		Title: "City",
	},
	"province_id": {
		Title: "Province",
	},
	"domisili_kelurahan_id": {
		Title: "Domisili Kelurahan",
	},
	"domisili_kecamatan_id": {
		Title: "Domisili Kecamatan",
	},
	"domisili_city_id": {
		Title: "Domisili City",
	},
	"domisili_province_id": {
		Title: "Domisili Province",
	},
}

var listColumnTitle = map[string]string{
	"Employee NIK":        "nik",
	"Employee Name":       "employee_name",
	"Birthdate":           "birthdate",
	"Birthplace":          "birthplace",
	"Email":               "email",
	"Email Office":        "email_office",
	"Phonenumber":         "phonenumber",
	"Division":            "department_id",
	"Department":          "department_detail_id",
	"Outsource":           "outsource_id",
	"Jabatan":             "jabatan_id",
	"Grade":               "grade_id",
	"Gender":              "gender_id",
	"Status Contract":     "status_contract_id",
	"Status PTKP":         "status_ptkp_id",
	"Cost Center":         "cost_center_id",
	"GL Number":           "gl_number_id",
	"Join Date":           "join_date",
	"Mine Permit Date":    "mine_permit_date",
	"Contract Start Date": "contract_start_date",
	"Contract End Date":   "contract_end_date",
	"Section":             "section_id",
	"Job Description":     "job_description_id",
	"Standard WorkDay":    "standard_workday",
	"Standard Overtime":   "standard_overtime",
	"Status Marriage":     "status_marriage_id",
	"NPWP":                "npwp",
	"Blood Type":          "blood_type_id",
	"Religion":            "religion_id",
	"Kelurahan":           "kelurahan_id",
	"Kecamatan":           "kecamatan_id",
	"City":                "city_id",
	"Province":            "province_id",
	"Domisili Kelurahan":  "domisili_kelurahan_id",
	"Domisili Kecamatan":  "domisili_kecamatan_id",
	"Domisili City":       "domisili_city_id",
	"Domisili Province":   "domisili_province_id",
	"BPJSTK Number":       "bpjstk_number",
	"BPJSKES Number":      "bpjskes_number",
}

func HandleReqColumn(listColumns []string) map[string]ColumnModel {
	// --- Step 1: clone static columns (avoid mutating global map) ---
	result := make(map[string]ColumnModel)
	for k, v := range listColumnsStatic {
		result[k] = v
	}

	// convert listColumns to map for fast lookup
	cols := make(map[string]bool)
	for _, c := range listColumns {
		cols[c] = true
	}

	// --- Step 2: update IsCollected for used columns ---
	for col := range result {
		if cols[col] {
			val := result[col]
			val.IsCollected = true
			result[col] = val
		}
	}

	// --- Step 3: remove unused columns except mandatory ones ---
	for key := range result {
		if key != "nik" && key != "employee_name" {
			if !cols[key] {
				delete(result, key)
			}
		}
	}

	return result
}

// ValidateExcelHeaders checks if incoming excel header titles are valid.
func ValidateExcelHeaders(headerRow []string, titleMap map[string]string) ([]string, error) {
	result := []string{}
	for _, h := range headerRow {
		key, ok := titleMap[h]
		if !ok {
			return nil, fmt.Errorf("unknown column: %s", h)
		}
		result = append(result, key)
	}
	return result, nil
}

func ParseEmployeeBulkRows(
	f *excelize.File,
	sheet string,
	headerRow []string,
	columnKeys []string,
	lookup *LookupMaps,
) ([]EmployeeUpdateItem, error) {

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	// Retrieve original data from CopySheet to identify changes
	copyRows, _ := f.GetRows(CopySheet)

	updates := []EmployeeUpdateItem{}

	// iterate row 2 → last
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// Get corresponding row from CopySheet
		var copyRow []string
		if i < len(copyRows) {
			copyRow = copyRows[i]
		}

		item := EmployeeUpdateItem{
			Nik:    "",
			Fields: map[string]interface{}{},
		}

		for c := range headerRow {
			if c >= len(row) {
				continue
			}

			key := columnKeys[c]
			val := strings.TrimSpace(row[c])

			// Set nik as identifier, it's always needed.
			if key == "nik" {
				item.Nik = val
				continue
			}

			// Determine if the cell value has changed compared to the original data
			var originalVal string
			if c < len(copyRow) {
				originalVal = strings.TrimSpace(copyRow[c])
			}

			// If values are the same, skip (no update needed)
			if val == originalVal {
				continue
			}

			if val == "" {
				continue
			}

			// --- DATE COLUMNS ---
			if IsDateColumn(key) {

				// 1. Coba parse sebagai excel serial float (45723)
				if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
					if t, err := excelize.ExcelDateToTime(floatVal, false); err == nil {
						item.Fields[key] = t
						continue
					}
				}

				// 2. Coba parse sebagai normal date string
				layouts := []string{
					"2006-01-02",
					"02-01-2006",
					"1/2/2006",
					"02/01/2006",
					"2/1/06",
					"2006/01/02",
					"January 2, 2006",
					"2 Jan 2006",
				}

				var parsed time.Time
				var valid bool
				for _, layout := range layouts {
					t, err := time.Parse(layout, val)
					if err == nil {
						parsed = t
						valid = true
						break
					}
				}

				if !valid {
					return nil, fmt.Errorf("invalid date '%s' at row %d col %d", val, i+1, c+1)
				}

				item.Fields[key] = parsed
				continue
			}

			// --- VALUE → ID MAPPING ---
			switch key {
			case "department_id":
				id := lookup.Department[val]
				if id == "" {
					return nil, fmt.Errorf("invalid department '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "department_detail_id":
				id := lookup.DepartmentDetail[val]
				if id == "" {
					return nil, fmt.Errorf("invalid department detail '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "outsource_id":
				id := lookup.Outsource[val]
				if id == "" {
					return nil, fmt.Errorf("invalid outsource '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "jabatan_id":
				id := lookup.Jabatan[val]
				if id == "" {
					return nil, fmt.Errorf("invalid jabatan '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "grade_id":
				id := lookup.Grade[val]
				if id == "" {
					return nil, fmt.Errorf("invalid grade '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue
			case "gender_id":
				id := lookup.Gender[val]
				if id == "" {
					return nil, fmt.Errorf("invalid gender '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue
			case "grade":
				id := lookup.Grade[val]
				if id == "" {
					return nil, fmt.Errorf("invalid grade '%s' at row %d", val, i+1)
				}
				item.Fields["grade_id"] = id
				continue

			case "status_ptkp_id":
				id := lookup.StatusPTKP[val]
				if id == "" {
					return nil, fmt.Errorf("invalid status ptkp '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "status_contract_id":
				id := lookup.StatusContract[val]
				if id == "" {
					return nil, fmt.Errorf("invalid status contract '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "cost_center_id":
				id := lookup.CostCenter[val]
				if id == "" {
					return nil, fmt.Errorf("invalid cost center '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "gl_number_id":
				id := lookup.GLNumber[val]
				if id == "" {
					return nil, fmt.Errorf("invalid gl number '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "job_description_id":
				id := lookup.JobDescription[val]
				if id == "" {
					return nil, fmt.Errorf("invalid job description '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "section_id":
				id := lookup.Section[val]
				if id == "" {
					return nil, fmt.Errorf("invalid section '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "status_marriage_id":
				id := lookup.StatusMarriage[val]
				if id == "" {
					return nil, fmt.Errorf("invalid StatusMarriage '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "blood_type_id":
				id := lookup.BloodType[val]
				if id == "" {
					return nil, fmt.Errorf("invalid BloodType '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "religion_id":
				id := lookup.Religion[val]
				if id == "" {
					return nil, fmt.Errorf("invalid Religion '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "kelurahan_id":
				id := lookup.Kelurahan[val]
				if id == "" {
					return nil, fmt.Errorf("invalid Kelurahan '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "kecamatan_id":
				id := lookup.Kecamatan[val]
				if id == "" {
					return nil, fmt.Errorf("invalid Kecamatan '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "city_id":
				id := lookup.City[val]
				if id == "" {
					return nil, fmt.Errorf("invalid City '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue

			case "province_id":
				id := lookup.Province[val]
				if id == "" {
					return nil, fmt.Errorf("invalid Province '%s' at row %d", val, i+1)
				}
				item.Fields[key] = id
				continue
			}

			// --- NORMAL TEXT FIELD ---
			item.Fields[key] = val
		}

		if item.Nik == "" {
			continue
		}

		updates = append(updates, item)
	}

	return updates, nil
}

func IsDateColumn(key string) bool {
	return key == "birthdate" ||
		key == "join_date" ||
		key == "mine_permit_date" ||
		key == "contract_start_date" ||
		key == "contract_end_date"
}

func ParseEmployeeDatasource(f *excelize.File, sheet string) (*LookupMaps, error) {

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed get datasource rows: %w", err)
	}

	lookup := &LookupMaps{
		Department:       map[string]string{},
		DepartmentDetail: map[string]string{},
		Outsource:        map[string]string{},
		Jabatan:          map[string]string{},
		Grade:            map[string]string{},
		Gender:           map[string]string{},
		StatusPTKP:       map[string]string{},
		StatusContract:   map[string]string{},
		CostCenter:       map[string]string{},
		GLNumber:         map[string]string{},
		JobDescription:   map[string]string{},
		Section:          map[string]string{},

		StatusMarriage: map[string]string{},
		BloodType:      map[string]string{},
		Religion:       map[string]string{},
		Kelurahan:      map[string]string{},
		Kecamatan:      map[string]string{},
		City:           map[string]string{},
		Province:       map[string]string{},
	}

	/* datasource header format:
	0	 A: Department
	1	 B: Department Id
	2	 C: Outsource
	3	 D: Outsource Id
	4	 E: Jabatan
	5	 F: Jabatan Id
	6	 G: Grade
	7	 H: Grade Id
	8	 I: Status Contract
	9	 J: Status Contract Id
	10	 K: Cost Center
	11	 L: Cost Center Id
	12	 M: GL Number
	13	 N: GL Number Id
	14	 O: Job Description
	15	 P: Job Description Id
	16	 Q: Section
	17	 R: Section Id
	18	 S: Section Department Detail Id
	19	 T: -
	20	 U: Status Marriage
	21	 V: Status Marriage Id
	22	 W: Gender
	23	 X: Gender Id
	24	 Y: Blood Type
	25	 Z: Blood Type Id
	26	AA: Religion
	27	AB: Religion Id
	28	AC: Kelurahan Id
	29	AD: Kelurahan
	30	AE: Kelurahan Kecamatan Id
	31	AF: Kecamatan
	32	AG: Kecamatan Id
	33	AH: Kecamatan City Id
	34	AI: City
	35	AJ: City Id
	36	AK: City Province Id
	37	AL: Province
	38	AM: Province Id
	39	AN: Department Detail Name
	40	AO: Department Detail Id
	41	AP: Department Detail (Department Id)
	42	AQ: Status PTKP
	43	AR: Status PTKP Id
		//  */

	for i := 1; i < len(rows); i++ { // skip header row (0)
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// Helper to safely read cell
		get := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		// Department
		name := get(0)
		id := get(1)
		if name != "" && id != "" {
			lookup.Department[name] = id
		}

		// Outsource
		name = get(2)
		id = get(3)
		if name != "" && id != "" {
			lookup.Outsource[name] = id
		}

		// Jabatan
		name = get(4)
		id = get(5)
		if name != "" && id != "" {
			lookup.Jabatan[name] = id
		}

		// Grade
		name = get(6)
		id = get(7)
		if name != "" && id != "" {
			lookup.Grade[name] = id
		}

		// Status Contract
		name = get(8)
		id = get(9)
		if name != "" && id != "" {
			lookup.StatusContract[name] = id
		}

		// Cost Center
		name = get(10)
		id = get(11)
		if name != "" && id != "" {
			lookup.CostCenter[name] = id
		}

		// GL Number
		name = get(12)
		id = get(13)
		if name != "" && id != "" {
			lookup.GLNumber[name] = id
		}

		// Job Description
		name = get(14)
		id = get(15)
		if name != "" && id != "" {
			lookup.JobDescription[name] = id
		}

		// Section
		name = get(16)
		id = get(17)
		if name != "" && id != "" {
			lookup.Section[name] = id
		}

		// Status Marriage
		name = get(20)
		id = get(21)
		if name != "" && id != "" {
			lookup.StatusMarriage[name] = id
		}

		// Gender
		name = get(22)
		id = get(23)
		if name != "" && id != "" {
			lookup.Gender[name] = id
		}

		// Blood Type
		name = get(24)
		id = get(25)
		if name != "" && id != "" {
			lookup.BloodType[name] = id
		}

		// Religion
		name = get(26)
		id = get(27)
		if name != "" && id != "" {
			lookup.Religion[name] = id
		}

		// Kelurahan
		name = get(29)
		id = get(30)
		if name != "" && id != "" {
			lookup.Kelurahan[name] = id
		}

		// Kecamatan
		name = get(31)
		id = get(32)
		if name != "" && id != "" {
			lookup.Kecamatan[name] = id
		}

		// City
		name = get(34)
		id = get(35)
		if name != "" && id != "" {
			lookup.City[name] = id
		}

		// Province
		name = get(37)
		id = get(38)
		if name != "" && id != "" {
			lookup.Province[name] = id
		}

		// Department Detail
		name = get(39)
		id = get(40)
		if name != "" && id != "" {
			lookup.DepartmentDetail[name] = id
		}

		// Status PTKP
		name = get(42)
		id = get(43)
		if name != "" && id != "" {
			lookup.StatusPTKP[name] = id
		}
	}

	return lookup, nil
}
