package masterdata

import (
	"context"
	"moyo-gateway-service/config"
	absence "moyo-gateway-service/proto/transaction/absence"
	adjustment "moyo-gateway-service/proto/transaction/adjustment"
	approval "moyo-gateway-service/proto/transaction/approval"
	assignment "moyo-gateway-service/proto/transaction/assignment"
	attendancemonitoring "moyo-gateway-service/proto/transaction/attendancemonitoring"
	dashboard "moyo-gateway-service/proto/transaction/dashboard"
	employeeLocation "moyo-gateway-service/proto/transaction/employeelocation"
	fileshare "moyo-gateway-service/proto/transaction/fileshare"
	lockingDepartment "moyo-gateway-service/proto/transaction/lockingdepartment"
	lockingEmployee "moyo-gateway-service/proto/transaction/lockingemployee"
	logging "moyo-gateway-service/proto/transaction/logging"
	mcu "moyo-gateway-service/proto/transaction/mcu"
	overview "moyo-gateway-service/proto/transaction/overview"
	performance "moyo-gateway-service/proto/transaction/performance"
	redirectPage "moyo-gateway-service/proto/transaction/redirectpage"
	reprimand "moyo-gateway-service/proto/transaction/reprimand"
	tag "moyo-gateway-service/proto/transaction/tag"
	timecard "moyo-gateway-service/proto/transaction/timecard"
	training "moyo-gateway-service/proto/transaction/training"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func Register(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption, conf config.Config) error {
	var err error
	transactionEndpoint := conf.Services.TransactionURL

	if err = adjustment.RegisterAdjustmentServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = assignment.RegisterAssignmentServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = approval.RegisterApprovalServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = lockingDepartment.RegisterLockingDepartmentServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = fileshare.RegisterFileShareServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = mcu.RegisterMcuServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}

	if err = absence.RegisterAbsenceServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}

	if err = lockingEmployee.RegisterLockingEmployeeServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}

	if err = redirectPage.RegisterRedirectPageServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}

	if err = attendancemonitoring.RegisterAttendanceMonitoringServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = overview.RegisterOverviewServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = timecard.RegisterTimecardServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = tag.RegisterTagServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = performance.RegisterPerformanceServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = dashboard.RegisterDashboardServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = logging.RegisterLogServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = training.RegisterTrainingServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = employeeLocation.RegisterEmployeeLocationServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}
	if err = reprimand.RegisterReprimandServiceHandlerFromEndpoint(ctx, mux, transactionEndpoint, opts); err != nil {
		return err
	}

	return nil

}
