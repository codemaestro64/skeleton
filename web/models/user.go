package models

type User struct {
	*session
}

const UserModel = "userModel"

func (d *Database) registerUserModel() {
	d.registerModelOnce(UserModel, &User{})
}
