package models

import "time"
type Role string



const(
	RoleAdmin Role = "admin"
	RoleUser Role="user"
	RoleDriver Role="driver"
	
)

func (r Role) IsValid() bool{
	return r==RoleAdmin || r== RoleUser || r==RoleDriver 

}
type User struct {
	ID             string    `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Email          string    `db:"email" json:"email"`
	HashedPassword string    `db:"hashed_password" json:"-"`
	Role           string    `db:"role" json:"role"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
