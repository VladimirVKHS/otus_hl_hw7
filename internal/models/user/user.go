package user

type User struct {
	Id        int
	Login     string `validate:"lte=255,required"`
	FirstName string `validate:"lte=255,required"`
	LastName  string `validate:"lte=255,required"`
	Password  string `validate:"lte=255,required"`
	Age       int    `validate:"min=18,max=150,required"`
	Interests string `validate:"lte=4096"`
	City      string `validate:"lte=255,required"`
	Sex       string `validate:"oneof=male female,required"`
	CreatedAt string
	IsPublic  bool
}

func (u *User) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"Id":        u.Id,
		"FirstName": u.FirstName,
		"LastName":  u.LastName,
		"Login":     u.Login,
		"Age":       u.Age,
		"Interests": u.Interests,
		"City":      u.City,
		"IsPublic":  u.IsPublic,
		"Sex":       u.Sex,
		"CreatedAt": u.CreatedAt,
	}
}

func UsersToResponse(users []*User) []map[string]interface{} {
	var result []map[string]interface{}
	for _, u := range users {
		result = append(result, u.ToResponse())
	}
	return result
}
