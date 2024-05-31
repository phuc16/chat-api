package dto

type TypeUserMessage string

const (
	TUM00 TypeUserMessage = "TUM00" // notify from root user
	TUM01 TypeUserMessage = "TUM01" // append friend request
	TUM02 TypeUserMessage = "TUM02" // remove friend request
	TUM03 TypeUserMessage = "TUM03" // accept friend request
	TUM04 TypeUserMessage = "TUM04" // unfriend
	TUM05 TypeUserMessage = "TUM05" // create conversation type STRANGER
	// TUM06 TypeUserMessage = "TUM06" // user online

)
