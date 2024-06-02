package entity

const (
	ALLOW_FRIENDS  string = "FRIENDS" // allow only friend send message
	ALLOW_EVERYONE string = "EVERYONE" // allow everyone send message
)

const (
	SHOW_BIRTHDAY_NO  string = "NO" // don't show birthday
	SHOW_BIRTHDAY_DMY string = "DMY" // show birthday
	SHOW_BIRTHDAY_DM  string = "DM" // show birthday only with month and day
)

type Setting struct {
	AllowMessaging string `bson:"allow_messaging" json:"allowMessaging"`
	ShowBirthday   string `bson:"show_birthday" json:"showBirthday"`
}
