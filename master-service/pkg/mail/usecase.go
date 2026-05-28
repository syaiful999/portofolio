package mail

import (
	"context"
	"fmt"
	"moyo-master-service/pkg/mail/proto"
	"moyo-master-service/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createMailServiceClient(serverAddr string) proto.MailServiceClient {

	// Dial the gRPC server with insecure transport credentials (not recommended for production)
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		utils.PushLogf("", "Failed to dial server: ", err.Error())
	}

	return proto.NewMailServiceClient(conn)
}

func sendEmail(client proto.MailServiceClient, request *proto.SendEmailRequest) *proto.SendEmailResponse {
	response, err := client.SendEmail(context.Background(), request)
	if err != nil {
		utils.PushLogf("", "SendEmail failed: ", err.Error())
	}
	return response
}

func handleSendEmailResponse(response *proto.SendEmailResponse) {
	utils.PushLogf("", fmt.Sprintf("SendEmail response: %v\n", response), "")
}

func SendMailActivateUser(address, sendTo, name, email, id, linkActivation string) {
	reciever := []string{sendTo}
	client := createMailServiceClient(address)
	requestParameter := &proto.SendEmailRequest{
		SendTo: reciever,
		// CcTo:           ccTo,
		Title:          "User Activation",
		LinkActivation: linkActivation,
		Param:          "activate-account",
		Employee: &proto.EmployeeModel{
			Name:  name,
			Email: email,
			Id:    id,
		},
	}
	sendEmailResponse := sendEmail(client, requestParameter)
	handleSendEmailResponse(sendEmailResponse)
}

func SendMailForgotPassword(address, sendTo, name, email, id, linkActivation string) {
	reciever := []string{sendTo}
	client := createMailServiceClient(address)
	requestParameter := &proto.SendEmailRequest{
		SendTo: reciever,
		// CcTo:           ccTo,
		Title:          "Forgot Password",
		LinkActivation: linkActivation,
		Param:          "forgot-password",
		Employee: &proto.EmployeeModel{
			Name:  name,
			Email: email,
			Id:    id,
		},
	}
	sendEmailResponse := sendEmail(client, requestParameter)
	handleSendEmailResponse(sendEmailResponse)
}
