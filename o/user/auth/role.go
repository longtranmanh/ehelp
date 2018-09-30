package auth

type Role int

const (
	RoleCustomer = Role(1)
	RoleEmployee = Role(2)
)

func (r Role) IsValid() bool {
	if r == RoleEmployee || r == RoleCustomer {
		return true
	}
	return false
}

func (r Role) IsCustomer() bool {
	return r == RoleCustomer
}

func (r Role) IsEmployee() bool {
	return r == RoleEmployee
}

type Rate struct {
	EmpID string `json:"emp_id"`
	Rate  int    `json:"rate"`
}
