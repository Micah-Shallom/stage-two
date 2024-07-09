package validator

type RegisterReq struct {
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

type OrgRegisterReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OrgReq struct {
	Name        string `json:"name"`
	OrgId       string `json:"orgId"`
	Description string `json:"description"`
}

type OrgAddUserReq struct {
	UserId string `json:"userId"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
