package employee

import (
	"context"
	pb "moyo-master-service/pkg/employee/proto"
	"moyo-master-service/utils"
)

// UpdateEmployeeBulk handles bulk employee update from Excel URL.
func (s *UseCase) UpdateEmployeeBulk(
	ctx context.Context,
	req *pb.EmployeeUploadRequest,
	res *pb.MutationEmployeeResponse,
	token *utils.TokenValue,
) error {

	// =====================================================
	// 1. Download & Open Excel file
	// =====================================================
	/* local directory
	filePath := "/Users/jhdkml/Documents/xapiens/moyo/moyo-master-service/pkg/employee/employee_update_template_08.xlsx"
	f, err := excelizev2.OpenFile(filePath) //localtest
	if err != nil {
		// Jika gagal, coba path alternatif atau log path absolut untuk debugging
		return errorResponse3(res, 400, "failed to open local excel file", err.Error())
	}

	// Pastikan menutup file setelah selesai digunakan untuk mencegah memory leak
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Error closing excel file:", err)
		}
	}()
	// */
	// /* path url
	f, err := utils.OpenExcelFromURL(req.FileName)
	if err != nil {
		return errorResponse3(res, 400, "invalid file url", err.Error())
	}
	// */

	// Read header row (row 1)
	header, err := f.GetRows(MainSheet)
	if err != nil {
		return errorResponse3(res, 400, "read excel", err.Error())
	}
	if len(header) == 0 {
		return errorResponse3(res, 400, "empty excel", "")
	}

	headerRow := header[0] // first row = header

	// =====================================================
	// 2. Validate header vs static column list
	// =====================================================
	columnKeys, err := ValidateExcelHeaders(headerRow, listColumnTitle)
	if err != nil {
		return errorResponse3(res, 400, "validate header", err.Error())
	}

	// =====================================================
	// 3. Parse each employee row into updates[]
	// =====================================================
	lookupMaps, err := ParseEmployeeDatasource(f, DsSheet)
	if err != nil {
		return errorResponse3(res, 400, "parse datasource", err.Error())
	}
	updates, err := ParseEmployeeBulkRows(
		f,
		MainSheet,
		headerRow,
		columnKeys,
		lookupMaps,
	)
	if err != nil {
		return errorResponse3(res, 400, "parse employee rows "+err.Error(), "")
	}

	if len(updates) == 0 {
		return errorResponse3(res, 400, "no data found", "")
	}

	// =====================================================
	// 4. Begin transaction (moved from repo → usecase)
	// =====================================================
	tx, err := s.db.Begin()
	if err != nil {
		return errorResponse3(res, 500, "DB begin", err.Error())
	}

	// =====================================================
	// 5. Call repository (no tx mgmt in repo)
	// =====================================================
	err = s.repository.UpdateBulkEmployees(ctx, tx, updates)
	if err != nil {
		tx.Rollback()
		return errorResponse3(res, 500, "bulk update failed", err.Error())
	}

	tx.Commit()

	// =====================================================
	// 6. Response OK
	// =====================================================

	return nil
}
