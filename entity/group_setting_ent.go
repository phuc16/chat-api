package entity

type GroupSetting struct {
	ChangeChatNameAndAvatar bool `bson:"change_chat_name_and_avatar" json:"changeChatNameAndAvatar"`
	PinMessages             bool `bson:"pin_messages" json:"pinMessages"`
	SendMessages            bool `bson:"send_messages" json:"sendMessages"`
	MembershipApproval      bool `bson:"membership_approval" json:"membershipApproval"`
	CreateNewPolls          bool `bson:"create_new_polls" json:"createNewPolls"`
}
