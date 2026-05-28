package entities

import "time"

// type RBAC struct {
// 	IsSuccess bool   `json:"IsSuccess"`
// 	Status    int    `json:"Status"`
// 	Message   string `json:"Message"`
// 	Data      struct {
// 		Page      int    `json:"Page"`
// 		Limit     int    `json:"Limit"`
// 		Sort      string `json:"Sort"`
// 		SortBy    string `json:"SortBy"`
// 		LastPage  string `json:"LastPage"`
// 		Keyword   string `json:"Keyword"`
// 		TotalRows int    `json:"TotalRows"`
// 		Items     []item `json:"Items"`
// 	} `json:"Data"`
// }

// type item struct {
// 	ID           string       `json:"ID"`
// 	PermissionID string       `json:"PermissionID"`
// 	RoleID       string       `json:"RoleID"`
// 	RoleName     string       `json:"RoleName"`
// 	Permission   []Permission `json:"Permission"`
// }

// type Permission struct {
// 	ID         string      `json:"ID"`
// 	RoleID     string      `json:"RoleID"`
// 	Name       string      `json:"Name"`
// 	Module     Module      `json:"Module"`
// 	Operations []Operation `json:"Operations"`
// }

// type Module struct {
// 	ID        string    `json:"ID"`
// 	Name      string    `json:"Name"`
// 	CreatedBy string    `json:"CreatedBy"`
// 	UpdatedBy string    `json:"UpdatedBy"`
// 	CreatedAt time.Time `json:"CreatedAt"`
// 	UpdatedAt time.Time `json:"UpdatedAt"`
// }

// type Operation struct {
// 	ID        string    `json:"ID"`
// 	Name      string    `json:"Name"`
// 	URL       string    `json:"URL"`
// 	Path      string    `json:"Path"`
// 	Method    string    `json:"Method"`
// 	CreatedBy string    `json:"CreatedBy"`
// 	UpdatedBy string    `json:"UpdatedBy"`
// 	CreatedAt time.Time `json:"CreatedAt"`
// 	UpdatedAt time.Time `json:"UpdatedAt"`
// }

type RBAC struct {
	ID       string   `json:"ID"`
	Type     string   `json:"Type"`
	Issuer   string   `json:"Issuer"`
	Metadata Metadata `json:"Metadata"`
	Secret   string   `json:"Secret"`
	Token    Token    `json:"Token"`
	User     User     `json:"User"`
	Rbac     struct {
		Application []struct {
			ID            string        `json:"ID"`
			Name          string        `json:"Name"`
			Icon          string        `json:"Icon"`
			Status        string        `json:"Status"`
			AccessControl AccessControl `json:"AccessControl"`
		} `json:"Application"`
	} `json:"RBAC"`
}

type Metadata struct {
	ID          string    `json:"ID"`
	CompanyID   string    `json:"CompanyID"`
	CompanyName string    `json:"CompanyName"`
	CompanyCode string    `json:"CompanyCode"`
	ClientKey   string    `json:"ClientKey"`
	SecretKey   string    `json:"SecretKey"`
	Platform    string    `json:"Platform"`
	Scope       string    `json:"Scope"`
	Status      bool      `json:"Status"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdateAt    time.Time `json:"UpdateAt"`
}

type Token struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
	Created      int    `json:"Created"`
	Expiry       int    `json:"Expiry"`
}

type User struct {
	Role    Role       `json:"Role"`
	Company Company    `json:"Company"`
	Account Account    `json:"Account"`
	User    UserDetail `json:"User"`
}

type Role struct {
	ID        string      `json:"ID"`
	Name      string      `json:"Name"`
	RoleGroup interface{} `json:"RoleGroup"`
}

type Company struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	Code string `json:"Code"`
}

type Account struct {
	Email      string `json:"Email"`
	FcmID      string `json:"FcmID"`
	GoogleID   string `json:"GoogleID"`
	FacebookID string `json:"FacebookID"`
	AppleID    string `json:"AppleID"`
	IsAD       bool   `json:"IsAD"`
	Status     bool   `json:"Status"`
}

type UserDetail struct {
	ID          string `json:"ID"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	IDNumber    string `json:"IDNumber"`
	Phone       string `json:"Phone"`
	Departement string `json:"Departement"`
	Position    string `json:"Position"`
	RoleID      string `json:"RoleID"`
	RoleName    string `json:"RoleName"`
	Grade       string `json:"Grade"`
}

type AccessControl struct {
	ID         string       `json:"ID"`
	RoleID     string       `json:"RoleID"`
	RoleName   string       `json:"RoleName"`
	CreatedBy  string       `json:"CreatedBy"`
	UpdateBy   string       `json:"UpdateBy"`
	CreatedAt  string       `json:"CreatedAt"`
	UpdateAt   string       `json:"UpdateAt"`
	Permission []Permission `json:"Permission"`
}

type Permission struct {
	ID            string       `json:"ID"`
	ParentID      string       `json:"ParentID"`
	RoleID        string       `json:"RoleID"`
	RoleName      string       `json:"RoleName"`
	ApplicationID string       `json:"ApplicationID"`
	Name          string       `json:"Name"`
	Module        Module       `json:"Module"`
	Operations    []Operation  `json:"Operations"`
	Permission    []Permission `json:"Permission"`
}

type Module struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	Path string `json:"Path"`
}

type Operation struct {
	ID           string `json:"ID"`
	PermissionID string `json:"PermissionID"`
	Name         string `json:"Name"`
	URL          string `json:"URL"`
	Path         string `json:"Path"`
	Method       string `json:"Method"`
}
