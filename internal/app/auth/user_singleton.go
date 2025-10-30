package auth

import "Lab1/internal/app/ds"

var currentUser *ds.User

func SetCurrentUser(user *ds.User) {
	currentUser = user
}

func GetCurrentUser() *ds.User {
	if currentUser == nil {
		// Возвращаем фиксированного пользователя (как требуется)
		return &ds.User{
			ID:    1,
			Login: "fixed_user",
		}
	}
	return currentUser
}
