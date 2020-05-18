package uc

type UserData struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
	EMail string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UserDataNoPassword struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
	EMail string `json:"email" bson:"email"`
}

type UserIndex struct {
	Id   string `json:"id" bson:"_id"`
	Type string `json:"type" bson:"type"`
	UserId string `json:"user_id" bson:"user_id"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionData struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name"`
}