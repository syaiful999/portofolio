package entities

type Headers struct {
	X_User_ID       string `json:"x-user-id"`
	X_User_Email    string `json:"x-user-email"`
	X_User_Fullname string `json:"x-user-fullname"`
	X_Company_ID    string `json:"x-company-id"`
	X_Company_Name  string `json:"x-company-name"`
	X_User_Token    string `json:"x-user-token"`
}

func (h *Headers) SetUserId(data string) {
	h.X_User_ID = data
}

func (h *Headers) SetUserEmail(data string) {
	h.X_User_Email = data
}

func (h *Headers) SetUserFullname(data string) {
	h.X_User_Fullname = data
}

func (h *Headers) SetCompanyID(data string) {
	h.X_Company_ID = data
}

func (h *Headers) SetCompanyName(data string) {
	h.X_Company_Name = data
}

func (h *Headers) Generate(userId, userEmail, userFullname, companyId, companyName string) {
	h.SetUserId(userId)
	h.SetUserEmail(userEmail)
	h.SetUserFullname(userFullname)
	h.SetCompanyID(companyId)
	h.SetCompanyName(companyName)
}
