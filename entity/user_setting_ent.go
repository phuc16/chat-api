package entity

const (
	ALLOW_FRIENDS  string = "friends"
	ALLOW_EVERYONE string = "everyone"
)

const (
	SHOW_BIRTHDAY_NO  string = "no"
	SHOW_BIRTHDAY_DMY string = "dmy"
	SHOW_BIRTHDAY_DM  string = "dm"
)

type Setting struct {
	AllowMessaging string `bson:"allow_messaging" json:"allowMessaging"`
	ShowBirthday   string `bson:"show_birthday" json:"showBirthday"`
}
