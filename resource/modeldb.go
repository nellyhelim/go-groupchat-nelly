package resource

import (
	"database/sql"
	"time"
)

type UserDB struct {
	UserID     sql.NullInt64  `db:"user_id"`
	UserName   sql.NullString `db:"username"`
	ProfilePic sql.NullString `db:"profile_pic"`
	Salt       sql.NullString `db:"salt"`
	Password   sql.NullString `db:"password"`
	CreatedAt  time.Time      `db:"created_at"`
}

// type RoomDB struct {
// 	RoomID      sql.NullInt64  `db:room_id`
// 	Name        sql.NullString `db:name`
// 	AdminUserID sql.NullInt64  `db:admin_user_id`
// 	Description sql.NullString `db:description`
// 	CategoryID  sql.NullInt64  `db:category_id`
// 	CreatedAt   time.Time      `db:"created_at"`
// }
