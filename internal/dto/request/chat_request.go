package dto

type CreateRoomRequest struct {
	RoomName string `json:"room_name" binding:"required"`
	Password string `json:"password,omitempty"`
}
