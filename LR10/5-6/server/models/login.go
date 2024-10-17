package models

type Login struct {
	ID       string `bson:"_id,omitempty" json:"id,omitempty"`
	Login    string `bson:"login" json:"login"`
	Password string `bson:"password" json:"password"`
	Role     string `bson:"role" json:"role"`
}
