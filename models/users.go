package models

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/strava/go.strava"
)

// Friend struct handles the MongoDB schema for each users friends
type Friend struct {
	ID             int64  `bson:"_id" json:"id"`
	FirstName      string `bson:"firstname" json:"firstName"`
	LastName       string `bson:"lastname" json:"lastName"`
	FullName       string `bson:"fullName" json:"fullName"`
	Photo          string `bson:"photo" json:"photo"`
	ChallengeCount int    `bson:"challengeCount" json:"challengeCount"`
	Wins           int    `bson:"wins" json:"wins"`
	Losses         int    `bson:"losses" json:"losses"`
}

// UserSegment struct handles the MongoDB schema for each users segments
type UserSegment struct {
	ID           int64  `bson:"_id" json:"id"`
	Name         string `bson:"name" json:"name"`
	Count        int    `bson:"count" json:"count"`
	ActivityType string `bson:"activityType" json:"activityType"`
}

// User struct handles the MongoDB schema for a user
type User struct {
	ID             int64          `bson:"_id" json:"id"`
	FirstName      string         `bson:"firstname" json:"firstName"`
	LastName       string         `bson:"lastname" json:"lastName"`
	FullName       string         `bson:"fullname" json:"fullName"`
	City           string         `bson:"city" json:"city"`
	State          string         `bson:"state" json:"state"`
	Country        string         `bson:"country" json:"country"`
	Gender         string         `bson:"gender" json:"gender"`
	Token          string         `bson:"token" json:"token"`
	Photo          string         `bson:"photo" json:"photo"`
	Email          string         `bson:"email" json:"email"`
	Friends        []*Friend      `bson:"friends" json:"friends"`
	Segments       []*UserSegment `bson:"segments" json:"segments"`
	Wins           int            `bson:"wins" json:"wins"`
	Losses         int            `bson:"losses" json:"losses"`
	ChallengeCount int            `bson:"challengeCount" json:"challengeCount"`
	CreatedAt      time.Time      `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt      time.Time      `bson:"updatedAt" json:"updatedAt,omitempty"`
	DeletedAt      *time.Time     `bson:"deletedAt" json:"deletedAt,omitempty"`
}

// GetUserByID gets a single stored user from MongoDB
func GetUserByID(id int64) (*User, error) {
	s := session.Copy()
	defer s.Close()

	var u User

	if err := s.DB(name).C("users").FindId(id).One(&u); err != nil {
		log.WithField("USER ID", id).Errorf("Unable to find user with id:\n%v", err)
		return nil, err
	}

	log.WithField("USER ID", u.ID).Infof("user found %d", u.ID)

	return &u, nil
}

// CreateUser creates user in MongoDB
func CreateUser(auth *strava.AuthorizationResponse) (*User, error) {
	s := session.Copy()
	defer s.Close()

	user := User{
		ID:        auth.Athlete.Id,
		FirstName: auth.Athlete.FirstName,
		LastName:  auth.Athlete.LastName,
		FullName:  auth.Athlete.FirstName + " " + auth.Athlete.LastName,
		City:      auth.Athlete.City,
		State:     auth.Athlete.State,
		Country:   auth.Athlete.Country,
		Gender:    string(auth.Athlete.Gender),
		Token:     auth.AccessToken,
		Photo:     auth.Athlete.Profile,
		Email:     auth.Athlete.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.DB(name).C("users").Insert(&user); err != nil {
		log.WithField("ID", user.ID).Errorf("Unable to create user with id:\n %v", err)
		return nil, err
	}

	return &user, nil
}

// RemoveFriendsSegmentsFromUser modifies user in MongoDB
func (u User) RemoveFriendsSegmentsFromUser() error {
	s := session.Copy()
	defer s.Close()

	u.Friends = u.Friends[:0]
	log.Infof("%d friends", len(u.Friends))
	u.Segments = u.Segments[:0]
	log.Infof("%d segments", len(u.Segments))
	u.UpdatedAt = time.Now()

	if err := s.DB(name).C("users").UpdateId(u.ID, &u); err != nil {
		log.WithField("USER ID", u.ID).Errorf("Unable to remove segments and friends from user:\n %v", err)
		return err
	}
	log.Infof("segments and friends removed from user %d", u.ID)
	return nil
}

// UpdateUser updates user in MongoDB
func (u User) UpdateUser(auth *strava.AuthorizationResponse) (*User, error) {
	s := session.Copy()
	defer s.Close()

	u.ID = auth.Athlete.Id
	u.FirstName = auth.Athlete.FirstName
	u.LastName = auth.Athlete.LastName
	u.FullName = auth.Athlete.FirstName + " " + auth.Athlete.LastName
	u.City = auth.Athlete.City
	u.State = auth.Athlete.State
	u.Country = auth.Athlete.Country
	u.Gender = string(auth.Athlete.Gender)
	u.Token = auth.AccessToken
	u.Photo = auth.Athlete.Profile
	u.Email = auth.Athlete.Email
	u.UpdatedAt = time.Now()

	if err := s.DB(name).C("users").UpdateId(u.ID, &u); err != nil {
		log.WithField("USER ID", u.ID).Errorf("Unable to update user:\n %v", err)
		return nil, err
	}

	return &u, nil
}

// RegisterUser creates a user in MongoDB
func RegisterUser(auth *strava.AuthorizationResponse) (*User, error) {
	u, err := GetUserByID(auth.Athlete.Id)
	if err != nil {
		log.WithField("USER ID", auth.Athlete.Id).Infof("Unable to find user with id %v creating user", auth.Athlete.Id)
		user, err := CreateUser(auth)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	log.WithField("USER ID", u.ID).Infof("Found user with id %v updating user", u.ID)
	user, err := u.UpdateUser(auth)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateAthlete updates user in MongoDB
func (u User) UpdateAthlete(athlete *strava.AthleteDetailed) (*User, error) {
	s := session.Copy()
	defer s.Close()

	u.ID = athlete.Id
	u.FirstName = athlete.FirstName
	u.LastName = athlete.LastName
	u.FullName = athlete.FirstName + " " + athlete.LastName
	u.City = athlete.City
	u.State = athlete.State
	u.Country = athlete.Country
	u.Gender = string(athlete.Gender)
	u.Photo = athlete.Profile
	u.Email = athlete.Email
	u.UpdatedAt = time.Now()

	if err := s.DB(name).C("users").UpdateId(u.ID, &u); err != nil {
		log.WithField("USER ID", u.ID).Errorf("Unable to update user %v:\n %v", u.ID, err)
		return nil, err
	}
	log.WithField("USER ID", u.ID).Infof("user %d updated from Strava", u.ID)
	return &u, nil
}

// SaveUserFriends save user friends
func (u User) SaveUserFriends(friends []*Friend) error {
	s := session.Copy()
	defer s.Close()
	u.Friends = friends
	u.UpdatedAt = time.Now()

	if err := s.DB(name).C("users").UpdateId(u.ID, &u); err != nil {
		log.Error("unable to save user friends")
		return err
	}
	log.WithField("USER ID", u.ID).Infof("stored %v friends", len(friends))
	return nil
}

// SaveUserSegments save user segments
func (u User) SaveUserSegments(segments []*UserSegment) error {
	s := session.Copy()
	defer s.Close()
	u.Segments = segments
	u.UpdatedAt = time.Now()

	if err := s.DB(name).C("users").UpdateId(u.ID, &u); err != nil {
		log.WithField("USER ID", u.ID).Error("unable to save user segments")
		return err
	}
	log.WithField("USER ID", u.ID).Infof("stored %v segments in db for user %v", len(segments), u.ID)
	return nil
}

// IncrementWins increment wins and challenge count for a particular user by id
func (u *User) IncrementWins(id int64) error {
	s := session.Copy()
	defer s.Close()
	u.Wins++
	u.ChallengeCount++

	// loop over friends for user to find id
	for _, friend := range u.Friends {
		if friend.ID == id {
			// found friend.. increment count and wins
			log.Infof("incrementing count and wins for friend %d", friend.ID)
			friend.ChallengeCount = friend.ChallengeCount + 1
			friend.Wins = friend.Wins + 1
			break
		}
	}

	u.UpdatedAt = time.Now()

	if err := s.DB(name).C("users").UpdateId(u.ID, &u); err != nil {
		log.Error("unable to save user increment wins")
		return err
	}
	log.WithField("USER ID", u.ID).Infof("incremented wins for challenge %d", id)
	return nil
}

// IncrementLosses increment losses and challenge count for a particular user by id
func (u *User) IncrementLosses(id int64) error {
	s := session.Copy()
	defer s.Close()
	u.Losses++
	u.ChallengeCount++

	// loop over friends for user to find id
	for _, friend := range u.Friends {
		if friend.ID == id {
			// found friend.. increment count and losses
			log.Infof("incrementing count and losses for friend %d", friend.ID)
			friend.ChallengeCount = friend.ChallengeCount + 1
			friend.Losses = friend.Losses + 1
			break
		}
	}

	u.UpdatedAt = time.Now()

	if err := s.DB(name).C("users").UpdateId(u.ID, &u); err != nil {
		log.Error("unable to save user increment losses")
		return err
	}
	log.WithField("USER ID", u.ID).Infof("incremented losses for challenge %d", id)
	return nil
}

// IncrementSegments increment segment count for a particular user by id
func (u *User) IncrementSegments(id int64) error {
	s := session.Copy()
	defer s.Close()

	// loop over segments for user to find id
	for _, segment := range u.Segments {
		if segment.ID == id {
			// found segment.. increment count
			log.Infof("incrementing count for segment %d", segment.ID)
			segment.Count = segment.Count + 1
			break
		}
	}

	u.UpdatedAt = time.Now()

	if err := s.DB(name).C("users").UpdateId(u.ID, &u); err != nil {
		log.Error("unable to save user segments increment")
		return err
	}
	log.WithField("USER ID", u.ID).Infof("incremented count for segment %d", id)
	return nil
}

// GetAllUsers returns all users from the DB
func GetAllUsers() ([]User, error) {
	s := session.Copy()
	defer s.Close()

	var users []User

	if err := s.DB(name).C("users").Find(nil).Sort("updatedAt").All(&users); err != nil {
		log.WithError(err).Error("Unable to return users")
		return nil, err
	}

	return users, nil
}

// RemoveUser deletes user from DB
func RemoveUser(ID int64) error {
	sess := session.Copy()
	defer sess.Close()

	if err := sess.DB(name).C("users").RemoveId(ID); err != nil {
		log.WithField("USER ID", ID).Errorf("Unable to remove user:\n %v", err)
		return err
	}

	return nil
}
