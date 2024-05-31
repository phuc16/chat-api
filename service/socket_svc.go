package service

import (
	"app/dto"
	"app/entity"
	"app/pkg/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ctxFromGin(c *gin.Context) context.Context {
	return c.Request.Context()
}

type WebSocketHandler struct {
	GroupSocketSvc IGroupSocketSvc
	UserSocketSvc  IUserSocketSvc
	ChatSocketSvc  IChatSocketSvc
	sessions       sync.Map // map[string]*sync.Map
}

func NewWebSocketHandler(groupSocketSvc IGroupSocketSvc, userSocketSvc IUserSocketSvc, chatSocketSvc IChatSocketSvc) *WebSocketHandler {
	return &WebSocketHandler{
		GroupSocketSvc: groupSocketSvc,
		UserSocketSvc:  userSocketSvc,
		ChatSocketSvc:  chatSocketSvc,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WebSocketHandler) HandleWebSocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	sessionID := conn.RemoteAddr().String()
	path := ctx.Request.URL.Path
	parts := utils.SplitPath(path)

	clients, _ := s.sessions.LoadOrStore(path, &sync.Map{})
	clients.(*sync.Map).Store(sessionID, conn)

	log.Println(parts)
	log.Println(parts[len(parts)-1])
	if len(parts) > 1 {
		switch parts[1] {
		case "group":
			s.HandleGroup(ctxFromGin(ctx), conn, sessionID, path)
		case "user":
			s.HandleUser(ctxFromGin(ctx), conn, sessionID, path)
		case "chat":
			s.HandleChat(ctxFromGin(ctx), conn, sessionID, parts[len(parts)-1], path)
		default:
			return
		}
	}
}

func (s *WebSocketHandler) HandleGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string) {
	defer func() {
		clients, ok := s.sessions.Load(path)
		if ok {
			log.Printf("Session end: %s", sessionID)
			clients.(*sync.Map).Delete(sessionID)
		}
	}()

	conn.WriteMessage(websocket.TextMessage, []byte("Connect group success"))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var groupDTO dto.GroupMessageDTO
		err = json.Unmarshal(message, &groupDTO)
		if err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}
		log.Printf("** Received message from %v: %v", sessionID, groupDTO)

		switch groupDTO.TGM {
		case dto.TGM01:
			s.HandleCreateGroup(ctx, conn, sessionID, path, message)
		case dto.TGM02:
			s.HandleDeleteGroup(ctx, conn, sessionID, path, message)
		case dto.TGM03:
			s.HandleAppendMemberGroup(ctx, conn, sessionID, path, message)
		case dto.TGM04:
			s.HandleAppendAdminGroup(ctx, conn, sessionID, path, message)
		case dto.TGM05:
			s.HandleRemoveAdminGroup(ctx, conn, sessionID, path, message)
		case dto.TGM06:
			s.HandleRemoveMemberGroup(ctx, conn, sessionID, path, message)
		case dto.TGM07:
			s.HandleChangeOwnerGroup(ctx, conn, sessionID, path, message)
		case dto.TGM08:
			s.HandleUpdateNameChatGroup(ctx, conn, sessionID, path, message)
		case dto.TGM09:
			s.HandleUpdateAvatarGroup(ctx, conn, sessionID, path, message)
		case dto.TGM10:
			s.HandleUpdateSetting_ChangeChatNameAndAvatar(ctx, conn, sessionID, path, message)
		case dto.TGM11:
			s.HandleUpdateSetting_PinMessages(ctx, conn, sessionID, path, message)
		case dto.TGM12:
			s.HandleUpdateSetting_SendMessages(ctx, conn, sessionID, path, message)
		case dto.TGM13:
			s.HandleUpdateSetting_MembershipApproval(ctx, conn, sessionID, path, message)
		case dto.TGM14:
			s.HandleUpdateSetting_CreateNewPolls(ctx, conn, sessionID, path, message)
		default:
			continue
		}
	}
}

func (s *WebSocketHandler) HandleUser(ctx context.Context, conn *websocket.Conn, sessionID, path string) {
	defer func() {
		clients, ok := s.sessions.Load(path)
		if ok {
			log.Printf("Session end: %s", sessionID)
			clients.(*sync.Map).Delete(sessionID)
		}
	}()

	conn.WriteMessage(websocket.TextMessage, []byte("Connect user success"))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var userDTO dto.UserMessageDTO
		err = json.Unmarshal(message, &userDTO)
		if err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}
		log.Printf("** Received message from %v: %v", sessionID, userDTO)

		switch userDTO.TUM {
		case dto.TUM01:
			s.HandleAppendFriendRequests(ctx, conn, sessionID, path, message)
		case dto.TUM02:
			s.HandleRemoveFriendRequests(ctx, conn, sessionID, path, message)
		case dto.TUM03:
			s.HandleAcceptFriendRequests(ctx, conn, sessionID, path, message)
		case dto.TUM04:
			s.HandleUnfriend(ctx, conn, sessionID, path, message)
		case dto.TUM05:
			s.HandleAppendConversations(ctx, conn, sessionID, path, message)
		default:
			continue
		}
	}
}

func (s *WebSocketHandler) HandleChat(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string) {
	defer func() {
		clients, ok := s.sessions.Load(path)
		if ok {
			log.Printf("Session end: %s", sessionID)
			chat, err := s.ChatSocketSvc.GetChatTop10(ctx, chatID)
			if err != nil {
				log.Printf("GetChatTop10 err: %v", err)
			}
			if err == nil {
				err = s.UserSocketSvc.UpdateConversations(ctx, *chat)
				if err != nil {
					log.Printf("UpdateConversations err: %v", err)
				}
			}
			clients.(*sync.Map).Delete(sessionID)
		}
	}()

	conn.WriteMessage(websocket.TextMessage, []byte("Connect chat success"))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var chatDTO dto.ChatMessageDTO
		err = json.Unmarshal(message, &chatDTO)
		if err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}
		log.Printf("** Received message from %v: %v", sessionID, chatDTO)

		switch chatDTO.TCM {
		case dto.TCM01:
			s.HandleAppendChat(ctx, conn, sessionID, chatID, path, message)
		case dto.TCM02:
			s.HandleChangeDeliveryChat(ctx, conn, sessionID, chatID, path, message)
		case dto.TCM03:
			s.HandleChangeReadChat(ctx, conn, sessionID, chatID, path, message)
		case dto.TCM04:
			s.HandleAppendHiddenMessage(ctx, conn, sessionID, chatID, path, message)
		case dto.TCM05:
			s.HandleRecallMessage(ctx, conn, sessionID, chatID, path, message)
		case dto.TCM06:
			s.HandleUserTypingMessage(ctx, conn, sessionID, chatID, path, message)
		case dto.TCM07:
			s.HandleAppendVoter(ctx, conn, sessionID, chatID, path, message)
		case dto.TCM08:
			s.HandleChangeVoter(ctx, conn, sessionID, chatID, path, message)
		case dto.TCM09:
			s.HandleLockVoting(ctx, conn, sessionID, chatID, path, message)

		default:
			continue
		}
	}
}

// group
func (s *WebSocketHandler) HandleCreateGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.CreateGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	listID := []string{req.Owner.UserID}
	for _, member := range req.Members {
		listID = append(listID, member.UserID)
	}

	_, err = s.GroupSocketSvc.Create(ctx, listID, req)
	if err != nil {
		log.Printf("** %v \n", err)
		s.SendMessageToClientUser(path, sessionID,
			dto.NotifyUser{
				UserMessageDTO: dto.UserMessageDTO{
					ID:  req.ID,
					TUM: dto.TUM00,
				},
				TypeNotify: dto.TYPE_NOTIFY_FAILED,
			}, "Failed | Create Group")
		log.Println("Create Group error:", err)
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS,
	}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | Create Group")
	s.SendMessageToAllClientsGroup(listID, req.Owner.UserID, req)
}

func (s *WebSocketHandler) HandleDeleteGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.DeleteGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.Delete(ctx, req.IDChat)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | Create Group")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | Delete Group")
	s.SendMessageToAllClientsGroup(arrayID, arrayID[0], req)
}

func (s *WebSocketHandler) HandleAppendMemberGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.AppendMemberGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.AppendMember(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | append Member Group")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | append Member Group")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleAppendAdminGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.AppendMemberGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.AppendAdmin(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | append admin Group")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | append admin Group")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleRemoveAdminGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.AppendMemberGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.RemoveAdmin(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | remove admin Group")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | remove admin Group")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleRemoveMemberGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.AppendMemberGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.RemoveMember(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | remove member Group")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | remove member Group")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleChangeOwnerGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.AppendMemberGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.ChangeOwner(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | change owner Group")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | change owner Group")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleUpdateNameChatGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.ChangeNameChatGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.UpdateNameChat(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | change chat name Group")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | change chat name Group")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleUpdateAvatarGroup(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.ChangeAvatarGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.UpdateAvatar(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | change avatar Group")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | change avatar Group")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleUpdateSetting_ChangeChatNameAndAvatar(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.UpdateSettingGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.UpdateSettingChangeChatNameAndAvatar(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | update setting change chat name and avatar Group")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | update setting change chat name and avatar Group")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleUpdateSetting_PinMessages(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.UpdateSettingGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.UpdateSettingPinMessages(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | update setting pin messages")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | update setting pin messages")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleUpdateSetting_SendMessages(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.UpdateSettingGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.UpdateSettingSendMessages(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | update setting send messages")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | update setting send messages")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleUpdateSetting_MembershipApproval(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.UpdateSettingGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.UpdateSettingMembershipApproval(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | update setting membership approval")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | update setting membership approval")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

func (s *WebSocketHandler) HandleUpdateSetting_CreateNewPolls(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.UpdateSettingGroupDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrayID, err := s.GroupSocketSvc.UpdateSettingCreateNewPolls(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyGroup{
			GroupMessageDTO: dto.GroupMessageDTO{
				ID:  req.ID,
				TGM: dto.TGM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientGroup(path, sessionID, notify, "Failed | update setting create new polls")
		return
	}

	notify := dto.NotifyGroup{
		GroupMessageDTO: dto.GroupMessageDTO{
			ID:  req.ID,
			TGM: dto.TGM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientGroup(path, sessionID, notify, "Pass | update setting create new polls")
	s.SendMessageToAllClientsGroup(arrayID, "", req)
}

// user
func (s *WebSocketHandler) HandleAppendFriendRequests(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.FriendRequestAddDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}
	req.SendAt = time.Now()
	err = s.UserSocketSvc.AppendFriendRequests(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyUser{
			UserMessageDTO: dto.UserMessageDTO{
				ID:  req.ID,
				TUM: dto.TUM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientUser(path, sessionID, notify, "Failed | append friend requests")
		return
	}

	notify := dto.NotifyUser{
		UserMessageDTO: dto.UserMessageDTO{
			ID:  req.ID,
			TUM: dto.TUM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientUser(path, sessionID, notify, "Pass | append friend requests")
	s.SendMessageToAllClientsUser(path, sessionID, req, "append friend requests")
}

func (s *WebSocketHandler) HandleRemoveFriendRequests(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.FriendRequestRemoveDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.UserSocketSvc.RemoveFriendRequests(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyUser{
			UserMessageDTO: dto.UserMessageDTO{
				ID:  req.ID,
				TUM: dto.TUM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientUser(path, sessionID, notify, "Failed | remove friend requests")
		return
	}

	notify := dto.NotifyUser{
		UserMessageDTO: dto.UserMessageDTO{
			ID:  req.ID,
			TUM: dto.TUM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientUser(path, sessionID, notify, "Pass | remove friend requests")
	s.SendMessageToAllClientsUser(path, sessionID, req, "remove friend requests")
}

func (s *WebSocketHandler) HandleAcceptFriendRequests(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.FriendRequestAcceptDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.UserSocketSvc.AcceptFriendRequests(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyUser{
			UserMessageDTO: dto.UserMessageDTO{
				ID:  req.ID,
				TUM: dto.TUM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientUser(path, sessionID, notify, "Failed | accept friend requests")
		return
	}

	notify := dto.NotifyUser{
		UserMessageDTO: dto.UserMessageDTO{
			ID:  req.ID,
			TUM: dto.TUM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientUser(path, sessionID, notify, "Pass | accept friend requests")
	s.SendMessageToAllClientsUser(path, sessionID, req, "accept friend requests")
}

func (s *WebSocketHandler) HandleUnfriend(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.UnfriendDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.UserSocketSvc.Unfriend(ctx, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyUser{
			UserMessageDTO: dto.UserMessageDTO{
				ID:  req.ID,
				TUM: dto.TUM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientUser(path, sessionID, notify, "Failed | unfriend")
		return
	}

	notify := dto.NotifyUser{
		UserMessageDTO: dto.UserMessageDTO{
			ID:  req.ID,
			TUM: dto.TUM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientUser(path, sessionID, notify, "Pass | unfriend")
	s.SendMessageToAllClientsUser(path, sessionID, req, "unfriend")
}

func (s *WebSocketHandler) HandleAppendConversations(ctx context.Context, conn *websocket.Conn, sessionID, path string, message []byte) {
	var req dto.AppendConversationDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.UserSocketSvc.AppendConversations(ctx, req, entity.TYPE_STRANGER)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyUser{
			UserMessageDTO: dto.UserMessageDTO{
				ID:  req.ID,
				TUM: dto.TUM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientUser(path, sessionID, notify, "Failed | append conversations")
		return
	}

	notify := dto.NotifyUser{
		UserMessageDTO: dto.UserMessageDTO{
			ID:  req.ID,
			TUM: dto.TUM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientUser(path, sessionID, notify, "Pass | append conversations")
	s.SendMessageToAllClientsUser(path, sessionID, req, "append conversations")
}

// chat
func (s *WebSocketHandler) HandleAppendChat(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string, message []byte) {
	var req dto.MessageAppendDTO
	log.Println("MessageAppendDTO", string(message))
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.ChatSocketSvc.AppendChat(ctx, chatID, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyChat{
			ChatMessageDTO: dto.ChatMessageDTO{
				ID:  req.ID,
				TCM: dto.TCM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientChat(path, sessionID, notify, "Failed | append chat")
		return
	}

	notify := dto.NotifyChat{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: dto.TCM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientChat(path, sessionID, notify, "Pass | append chat")
	s.SendMessageToAllClientsChat(path, sessionID, req, "append chat")
}

func (s *WebSocketHandler) HandleChangeDeliveryChat(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string, message []byte) {
	var req dto.MessageDeliveryDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.ChatSocketSvc.ChangeDeliveryChat(ctx, chatID, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyChat{
			ChatMessageDTO: dto.ChatMessageDTO{
				ID:  req.ID,
				TCM: dto.TCM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientChat(path, sessionID, notify, "Failed | change delivery chat")
		return
	}

	notify := dto.NotifyChat{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: dto.TCM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientChat(path, sessionID, notify, "Pass | change delivery chat")
	s.SendMessageToAllClientsChat(path, sessionID, req, "change delivery chat")
}

func (s *WebSocketHandler) HandleChangeReadChat(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string, message []byte) {
	var req dto.MessageDeliveryDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.ChatSocketSvc.ChangeReadChat(ctx, chatID, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyChat{
			ChatMessageDTO: dto.ChatMessageDTO{
				ID:  req.ID,
				TCM: dto.TCM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientChat(path, sessionID, notify, "Failed | change read chat")
		return
	}

	notify := dto.NotifyChat{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: dto.TCM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientChat(path, sessionID, notify, "Pass | change read chat")
	s.SendMessageToAllClientsChat(path, sessionID, req, "change read chat")
}

func (s *WebSocketHandler) HandleAppendHiddenMessage(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string, message []byte) {
	var req dto.MessageHiddenDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.ChatSocketSvc.AppendHiddenMessage(ctx, chatID, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyChat{
			ChatMessageDTO: dto.ChatMessageDTO{
				ID:  req.ID,
				TCM: dto.TCM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientChat(path, sessionID, notify, "Failed | change hidden chat")
		return
	}

	notify := dto.NotifyChat{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: dto.TCM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientChat(path, sessionID, notify, "Pass | change hidden chat")
	s.SendMessageToAllClientsChat(path, sessionID, req, "change hidden chat")
}

func (s *WebSocketHandler) HandleRecallMessage(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string, message []byte) {
	var req dto.MessageHiddenDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.ChatSocketSvc.RecallMessage(ctx, chatID, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyChat{
			ChatMessageDTO: dto.ChatMessageDTO{
				ID:  req.ID,
				TCM: dto.TCM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientChat(path, sessionID, notify, "Failed | recall message")
		return
	}

	notify := dto.NotifyChat{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: dto.TCM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientChat(path, sessionID, notify, "Pass | recall message")
	s.SendMessageToAllClientsChat(path, sessionID, req, "recall message")
}

func (s *WebSocketHandler) HandleUserTypingMessage(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string, message []byte) {
	var req dto.TypingTextMessageDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}
	s.SendMessageToAllClientsChat(path, sessionID, req, "user typing a text message")
	return
}

func (s *WebSocketHandler) HandleAppendVoter(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string, message []byte) {
	var req dto.MessageAppendDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrStr := strings.Split(req.Contents[0].Value, "@")
	appendVoterDTO := dto.AppendVoterDTO{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: req.TCM,
		},
		VotingID: arrStr[0],
		Name:     arrStr[1],
		Voter: entity.PersonInfo{
			UserID:     req.UserID,
			UserName:   req.UserName,
			UserAvatar: req.UserAvatar,
		},
	}

	err = s.ChatSocketSvc.AppendVoter(ctx, appendVoterDTO, chatID, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyChat{
			ChatMessageDTO: dto.ChatMessageDTO{
				ID:  req.ID,
				TCM: dto.TCM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientChat(path, sessionID, notify, "Failed | append voter")
		return
	}

	notify := dto.NotifyChat{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: dto.TCM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientChat(path, sessionID, notify, "Pass | append voter")
	s.SendMessageToAllClientsChat(path, sessionID, req, "append voter")
}

func (s *WebSocketHandler) HandleChangeVoter(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string, message []byte) {
	var req dto.MessageAppendDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	arrStr := strings.Split(req.Contents[0].Value, "@")
	changeVoterDTO := dto.ChangeVoterDTO{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: req.TCM,
		},
		VotingID: arrStr[0],
		OldName:  arrStr[1],
		NewName:  arrStr[2],
		Voter: entity.PersonInfo{
			UserID:     req.UserID,
			UserName:   req.UserName,
			UserAvatar: req.UserAvatar,
		},
	}

	err = s.ChatSocketSvc.ChangeVoting(ctx, changeVoterDTO, chatID, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyChat{
			ChatMessageDTO: dto.ChatMessageDTO{
				ID:  req.ID,
				TCM: dto.TCM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientChat(path, sessionID, notify, "Failed | change voter")
		return
	}

	notify := dto.NotifyChat{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: dto.TCM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientChat(path, sessionID, notify, "Pass | change voter")
	s.SendMessageToAllClientsChat(path, sessionID, req, "change voter")
}

func (s *WebSocketHandler) HandleLockVoting(ctx context.Context, conn *websocket.Conn, sessionID, chatID, path string, message []byte) {
	var req dto.MessageAppendDTO
	err := json.Unmarshal(message, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return
	}

	err = s.ChatSocketSvc.LockVoting(ctx, chatID, req)
	if err != nil {
		log.Printf("** %v \n", err)
		notify := dto.NotifyChat{
			ChatMessageDTO: dto.ChatMessageDTO{
				ID:  req.ID,
				TCM: dto.TCM00,
			},
			TypeNotify: dto.TYPE_NOTIFY_FAILED}
		if err.Error() == "CONFLICT" {
			notify.TypeNotify = dto.TYPE_NOTIFY_CONFLICT
		}
		s.SendMessageToClientChat(path, sessionID, notify, "Failed | lock voting")
		return
	}

	notify := dto.NotifyChat{
		ChatMessageDTO: dto.ChatMessageDTO{
			ID:  req.ID,
			TCM: dto.TCM00,
		},
		TypeNotify: dto.TYPE_NOTIFY_SUCCESS}
	s.SendMessageToClientChat(path, sessionID, notify, "Pass | lock voting")
	s.SendMessageToAllClientsChat(path, sessionID, req, "lock voting")
}

// notify group
func (s *WebSocketHandler) SendMessageToClientGroup(path, sessionID string, notify dto.NotifyGroup, logStr string) {
	log.Println("sendMessageToClient Group", logStr)
	message, err := json.Marshal(notify)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}

	clients, ok := s.sessions.Load(path)
	if ok {
		clients.(*sync.Map).Range(func(key, value interface{}) bool {
			if key.(string) == sessionID {
				err = value.(*websocket.Conn).WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error sending message: %s", err)
				}
			}
			return true
		})
	}
}

func (s *WebSocketHandler) SendMessageToAllClientsGroup(arrayID []string, ignore string, obj interface{}) {
	log.Println("sendMessageToAllClients group")
	message, err := json.Marshal(obj)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}

	for _, id := range arrayID {
		if id == ignore {
			continue
		}
		path := "/ws/user/" + id
		clients, ok := s.sessions.Load(path)
		if ok {
			clients.(*sync.Map).Range(func(key, value interface{}) bool {
				err = value.(*websocket.Conn).WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error sending message: %s", err)
				}
				return true
			})
		}
	}
}

// notify chat
func (s *WebSocketHandler) SendMessageToClientChat(path, sessionID string, notify dto.NotifyChat, logStr string) {
	log.Printf("** sendMessageToClient Chat %s", logStr)
	message, err := json.Marshal(notify)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}

	clients, ok := s.sessions.Load(path)
	if ok {
		clients.(*sync.Map).Range(func(key, value interface{}) bool {
			if key.(string) == sessionID {
				err = value.(*websocket.Conn).WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error sending message: %s", err)
				}
			}
			return true
		})
	}
}

func (s *WebSocketHandler) SendMessageToAllClientsChat(path, sessionID string, obj interface{}, logStr string) {
	log.Printf("** sendMessageToAllClients Chat %s", logStr)
	message, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Error marshalling ChatMessageDTO: %s", err)
		return
	}

	clients, ok := s.sessions.Load(path)
	if ok {
		clients.(*sync.Map).Range(func(key, value interface{}) bool {
			if key.(string) != sessionID {
				err = value.(*websocket.Conn).WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error sending message: %s", err)
				}
			}
			return true
		})
	}
}

// notify user
func (s *WebSocketHandler) SendMessageToClientUser(path, sessionID string, notify dto.NotifyUser, logStr string) {
	log.Printf("** sendMessageToClient User %s", logStr)
	message, err := json.Marshal(notify)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}

	clients, ok := s.sessions.Load(path)
	if ok {
		clients.(*sync.Map).Range(func(key, value interface{}) bool {
			if key.(string) == sessionID {
				err = value.(*websocket.Conn).WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error sending message: %s", err)
				}
			}
			return true
		})
	}
}

func (s *WebSocketHandler) SendMessageToAllClientsUser(path, sessionID string, obj interface{}, logStr string) {
	log.Printf("** sendMessageToAllClients User %s", logStr)
	message, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Error marshalling ChatMessageDTO: %s", err)
		return
	}

	clients, ok := s.sessions.Load(path)
	if ok {
		clients.(*sync.Map).Range(func(key, value interface{}) bool {
			if key.(string) != sessionID {
				err = value.(*websocket.Conn).WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error sending message: %s", err)
				}
			}
			return true
		})
	}
}
