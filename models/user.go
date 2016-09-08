package models

// User struct contains the user response from MongoDB
type User struct {
	ID        int    `bson:"id"`
	FirstName string `bson:"firstName"`
}

// GetUserByID gets a single stored user from MongoDB
func GetUserByID(id int) (*User, error) {
	var u User

	return &u, nil
}
