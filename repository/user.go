package repository

import "github.com/workdestiny/watgok_web/entity"

//GetUser get user data
func GetUser(q Queryer, userID string) (*entity.UserModel, error) {
	var u entity.UserModel
	err := q.QueryRow(`
			SELECT id, name, display, role
			  FROM users
			 WHERE id = $1
	`, userID).Scan(&u.ID, &u.Name, &u.Display, &u.Role)

	return &u, err
}

//GetUserByFBID get by facebook id
func GetUserByFBID(q Queryer, fbID string) (*entity.UserModel, error) {
	var u entity.UserModel
	err := q.QueryRow(`
			SELECT id, name, display, role
			  FROM users
			 WHERE fb_id = $1
	`, fbID).Scan(&u.ID, &u.Name, &u.Display, &u.Role)

	return &u, err
}

//CreateUser create new
func CreateUser(q Queryer, name, display, fbID string, role entity.Role) (string, error) {
	var id string
	err := q.QueryRow(`
		INSERT INTO users
			        (name, display, fb_id, role)
			 VALUES ($1, $2, $3, $4)
		  RETURNING id;
	`, name, display, fbID, role).Scan(&id)
	return id, err
}
