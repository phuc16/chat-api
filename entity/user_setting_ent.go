package entity

const (
	ALLOW_FRIENDS  string = "FRIENDS"
	ALLOW_EVERYONE string = "EVERYONE"
)

const (
	SHOW_BIRTHDAY_NO  string = "NO"
	SHOW_BIRTHDAY_DMY string = "DMY"
	SHOW_BIRTHDAY_DM  string = "DM"
)

type Setting struct {
	AllowMessaging string `bson:"allow_messaging" json:"allowMessaging"`
	ShowBirthday   string `bson:"show_birthday" json:"showBirthday"`
}
