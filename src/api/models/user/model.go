package user

type User struct {
	Email       string `json:"email"`
	UID         string `json:"uid"`
	Disabled    bool   `json:"disabled"`
	Reports     int    `json:"reports"`
	DisplayName string `json:"display_name"`
	PhoneNumber string `json:"phone_number"`
	PhotoURL    string `json:"photo_url"`
	SignUp      int64  `json:"sign_up"`
}

type UserComparator struct {
	Users  []User
	lessFn func(i, j User) bool
}

func NewUserComparator(users []User, lessFn func(i, j User) bool) UserComparator {
	return UserComparator{
		Users:  users,
		lessFn: lessFn,
	}
}

func (a UserComparator) Len() int { return len(a.Users) }

func ByDisabled(i, j User) bool {
	if i.Disabled == j.Disabled {
		return false
	}

	return i.Disabled == false
}

func ByReports(i, j User) bool {
	return i.Reports < j.Reports
}

func BySignUp(i, j User) bool {
	return i.SignUp < j.SignUp
}
func (a UserComparator) Less(i, j int) bool {
	if a.lessFn == nil {
		a.lessFn = BySignUp
	}
	return a.lessFn(a.Users[i], a.Users[j])
}
func (a UserComparator) Swap(i, j int) { a.Users[i], a.Users[j] = a.Users[j], a.Users[i] }
