package base

const (
	OP_NONE     = iota
	OP_LOGIN    = 1 //登录
	OP_REGISTER = 2 //注册
	OP_SYNC     = 3 //同步数据

)

const (
	PLAYER_STATE_NONE  = iota //玩家状态_无
	PLAYER_STATE_LOGIN = 0x01 //玩家状态_登录
	PLAYER_STATE_GAME  = 0x02 //玩家状态_游戏
)
