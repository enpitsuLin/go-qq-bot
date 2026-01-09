package model

// WebhookPayload 标准 Webhook 请求格式
type WebhookPayload struct {
	ID uint32      `json:"id"` // 事件 ID
	Op uint32      `json:"op"` // 操作码：0=事件，13=验证
	D  interface{} `json:"d"`  // 数据对象
	S  uint32      `json:"s"`  // 序列号
	T  string      `json:"t"`  // 事件类型
}

// VerificationPayload 验证请求数据（op=13）
type VerificationPayload struct {
	PlainToken string `json:"plain_token"` // 明文 token
	EventTs    string `json:"event_ts"`    // 事件时间戳
}

// VerificationResponse 验证响应
type VerificationResponse struct {
	PlainToken   string `json:"plain_token"`    // 原始 token
	MsgSig       string `json:"msg_sig"`        // 签名
	ResponceTime string `json:"responsce_time"` // 响应时间（注意拼写是官方格式）
}

// EventPayload 事件数据
type EventPayload struct {
	GroupID     string  `json:"group_id,omitempty"`     // 群组 ID
	GroupOpenID string  `json:"group_openid,omitempty"` // 群组 Open ID
	Author      *Author `json:"author,omitempty"`       // 作者信息
	Content     string  `json:"content,omitempty"`      // 消息内容
	ID          string  `json:"id,omitempty"`           // 消息 ID
	Timestamp   string  `json:"timestamp,omitempty"`    // 时间戳
}

// Author 消息作者信息
type Author struct {
	ID           string `json:"id,omitempty"`            // 用户 ID
	UserOpenID   string `json:"user_openid,omitempty"`   // 用户 Open ID
	MemberOpenID string `json:"member_openid,omitempty"` // 成员 Open ID（群组场景）
}

// EventType 事件类型
type EventType string

// 事件类型常量
const (
	// 群聊和私聊事件
	GroupAtMessageCreate EventType = "GROUP_AT_MESSAGE_CREATE" // 群组 @消息
	C2cMessageCreate     EventType = "C2C_MESSAGE_CREATE"      // 私聊消息
	FriendAdd            EventType = "FRIEND_ADD"              // 添加好友
	FriendDel            EventType = "FRIEND_DEL"              // 删除好友
	GroupAddRobot        EventType = "GROUP_ADD_ROBOT"         // 机器人被添加到群组
	GroupDelRobot        EventType = "GROUP_DEL_ROBOT"         // 机器人被移出群组

	// 频道事件
	MessageCreate       EventType = "MESSAGE_CREATE"        // 频道消息创建
	MessageDelete       EventType = "MESSAGE_DELETE"        // 频道消息删除
	DirectMessageCreate EventType = "DIRECT_MESSAGE_CREATE" // 频道私信
	AtMessageCreate     EventType = "AT_MESSAGE_CREATE"     // 频道 @消息

	// 其他事件
	GuildCreate EventType = "GUILD_CREATE" // 频道创建
	GuildUpdate EventType = "GUILD_UPDATE" // 频道更新
	GuildDelete EventType = "GUILD_DELETE" // 频道删除
)
