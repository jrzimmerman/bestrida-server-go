package models

// Friend struct handles the MongoDB schema for each of a users friends
type Friend struct {
	ID             int    `bson:"id"`
	UserName       string `bson:"userName"`
	FirstName      string `bson:"firstName"`
	LastName       string `bson:"lastName"`
	FullName       string `bson:"fullName"`
	Photo          string `bson:"photo"`
	ChallengeCount int    `bson:"challengeCount"`
	Wins           int    `bson:"wins"`
	Losses         int    `bson:"losses"`
}

// Segment struct handles the MongoDB schema for each of a users segments
type Segment struct {
	ID    int    `bson:"id"`
	Name  string `bson:"name"`
	Count int    `bson:"count"`
}

// User struct handles the MongoDB schema for a user
type User struct {
	ID        int       `bson:"id"`
	FirstName string    `bson:"firstName"`
	LastName  string    `bson:"lastName"`
	FullName  string    `bson:"fullName"`
	Token     string    `bson:"token"`
	Photo     string    `bson:"photo"`
	Email     string    `bson:"email"`
	Friends   []Friend  `bson:"friends"`
	Segments  []Segment `bson:"segments"`
	Wins      int       `bson:"wins"`
	Losses    int       `bson:"losses"`
}

// GetUserByID gets a single stored user from MongoDB
func GetUserByID(id int) (*User, error) {
	var u User

	return &u, nil
}
