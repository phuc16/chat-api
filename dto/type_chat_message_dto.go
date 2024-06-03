package dto

type TypeChatMessage string

const (
	TCM00 TypeChatMessage = "TCM00" // notify from root user
	TCM01 TypeChatMessage = "TCM01" // send message
	TCM02 TypeChatMessage = "TCM02" // message delivery
	TCM03 TypeChatMessage = "TCM03" // message read
	TCM04 TypeChatMessage = "TCM04" // message hidden
	TCM05 TypeChatMessage = "TCM05" // message recall
	TCM06 TypeChatMessage = "TCM06" // user typing a text message
	TCM07 TypeChatMessage = "TCM07" // append voter
	TCM08 TypeChatMessage = "TCM08" // change voter
	TCM09 TypeChatMessage = "TCM09" // lock voting
)
