package userprofile

type CreateUser struct {
	FullName string `json:"fullname,omitempty"`
	Age      int    `json:"age,omitempty"`
	Gender   string `json:"gender,omitempty"`
}

type GetUserParams struct {
}
