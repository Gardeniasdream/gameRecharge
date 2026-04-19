package gamesrv

import "log"

// SendDiamond 调用游戏服发放钻石
func SendDiamond(serverId int, roleId string, diamond int) {
	log.Printf("【游戏服对接】区服：%d, 角色：%s, 发放钻石：%d", serverId, roleId, diamond)
}
