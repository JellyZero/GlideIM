package msgdao

// ChatMessage 一对一聊天全量消息
type ChatMessage struct {
	// MID 消息 ID
	MID int64 `gorm:"primary_key"`
	// ReceiveSeq 接收者全局消息 Seq
	ReceiveSeq int64
	// CliSeq 发送者消息 seq
	CliSeq int64
	// From 发送者ID
	From int64
	// To 接收者ID
	To int64
	// Type 消息类型
	Type int64
	// SendAt 发送时间
	SendAt int64
	// Content 消息内容
	Content string
}

// OfflineMessage 用户不在线, 离线消息
type OfflineMessage struct {
	ID  int64 `gorm:"primary_key"`
	MID int64
	UID int64
}

// GroupMessage 全量群消息
type GroupMessage struct {
	MID int64 `gorm:"primary_key"`
	// Seq 群消息 seq
	Seq int64
	// To 群 ID
	To int64
	// From 发送者 ID
	From    int64
	Type    int64
	SendAt  int64
	Content string
}

// GroupMemberMsgState 群成员确认收到消息记录, 用于计算离线消息的同步量
type GroupMemberMsgState struct {
	// MbID 群成员ID, GID+UID 拼接成
	MbID string `gorm:"primary_key"`
	GID  int64
	UID  int64
	// LastAckMID 最后一次确认收到的消息 id
	LastAckMID int64
	// LastAckSeq 最后一次确认收到的消息 seq
	LastAckSeq int64
}

// GroupMessageState 群消息最新状态 ID 及 seq
type GroupMessageState struct {
	GID int64 `gorm:"primary_key"`
	// LastMID 最后一条消息的ID
	LastMID int64
	// LastSeq	最后一条消息的 seq
	LastSeq int64
	// LastMsgAt 最后一条消息的发送时间
	LastMsgAt int64
}

// GroupMsgSeq 群消息 seq 状态
type GroupMsgSeq struct {
	GID int64 `gorm:"primary_key"`
	// Seq 当前群消息 seq
	Seq int64
	// Step 增长步长
	Step int64
}
