package str

//Exchange

const (
	MessageExchange = "message_exchange"
	RetryExchange   = "retry_exchange"
	LikeExchange    = "like_exchange"
)

//Queue

const (
	MessageLike       = "message_Like"
	MessageMention    = "message_Mention"
	MessagePrivateMsg = "message_Private_Msg"
	MessageReply      = "message_Reply"
	MessageSystem     = "message_System"
	LikePost          = "like_post"
	LikeComment       = "like_comment"
)

//routing_key

const (
	RoutMessageLike = "message.like"
	RoutMention     = "message.mention"
	RoutPrivateMsg  = "message.private_Msg"
	RoutReply       = "message.reply"
	RoutSystem      = "message.system"
	RoutPost        = "like.post"
	RoutComment     = "like.comment"
)
