package dto


type TypeGroupMessage string

const (
	TGM00 TypeGroupMessage = "TGM00" // notify from root user
	TGM01 TypeGroupMessage = "TGM01" // create group
	TGM02 TypeGroupMessage = "TGM02" // delete group
	TGM03 TypeGroupMessage = "TGM03" // append member
	TGM04 TypeGroupMessage = "TGM04" // append admin
	TGM05 TypeGroupMessage = "TGM05" // remove admin
	TGM06 TypeGroupMessage = "TGM06" // remove member
	TGM07 TypeGroupMessage = "TGM07" // change owner
	TGM08 TypeGroupMessage = "TGM08" // change name chat
	TGM09 TypeGroupMessage = "TGM09" // change avatar
	TGM10 TypeGroupMessage = "TGM10" // update setting change chat name and avatar
	TGM11 TypeGroupMessage = "TGM11" // update setting pin messages
	TGM12 TypeGroupMessage = "TGM12" // update setting send messages
	TGM13 TypeGroupMessage = "TGM13" // update setting membership approval
	TGM14 TypeGroupMessage = "TGM14" // update setting create new polls
)
