package entity

type LoginRequest struct {
	Device   int64  `json:"device"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Token    string `json:"token"`
	DeviceId int64  `json:"device_id"`
	Uid      int64  `json:"uid"`
}

type RegisterRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// AuthorResponse login or register result
type AuthorResponse struct {
	Token string
	Uid   int64
}

type UserInfoRequest struct {
	Uid []int64
}

type UserInfoResponse struct {
	Uid      int64
	Nickname string
	Avatar   string
}

type UserInfoListResponse struct {
	UserInfo []*UserInfoResponse
}

type UserNewChatRequest struct {
	Id   int64
	Type int8
}

type GroupResponse struct {
	Gid    int64
	Name   string
	Avatar string
}

type RelationResponse struct {
	Groups  []GroupResponse
	Friends []int64
}

type ChatHistoryRequest struct {
	Cid  int64
	Time int64
	Type int8
}

type ChatInfoRequest struct {
	Cid int64
}

type CreateGroupRequest struct {
	Name string
}

type JoinGroupRequest struct {
	Gid int64
}

type ExitGroupRequest struct {
	Gid int64
}

type GetGroupMemberRequest struct {
	Gid int64
}

type GroupMemberResponse struct {
	Uid        int64
	Nickname   string
	RemarkName string
	Type       int8
	Online     bool
}

type AddMemberRequest struct {
	Gid int64
	Uid []int64
}

type RemoveMemberRequest struct {
	Gid int64
	Uid []int64
}
