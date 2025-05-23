package setting

import (
	"context"
	"fmt"
	"wox/common"
	"wox/util"
)

type ResultHash string

type WoxAppData struct {
	QueryHistories  []QueryHistory
	ActionedResults *util.HashMap[ResultHash, []ActionedResult]
	FavoriteResults *util.HashMap[ResultHash, bool]
}

type QueryHistory struct {
	Query     common.PlainQuery
	Timestamp int64
}

type ActionedResult struct {
	Timestamp int64
	Query     string // Record the raw query text when the user performs action on this result
}

func NewResultHash(pluginId string, title, subTitle string) ResultHash {
	return ResultHash(util.Md5([]byte(fmt.Sprintf("%s%s%s", pluginId, title, subTitle))))
}

func GetDefaultWoxAppData(ctx context.Context) WoxAppData {
	return WoxAppData{
		QueryHistories:  []QueryHistory{},
		ActionedResults: util.NewHashMap[ResultHash, []ActionedResult](),
		FavoriteResults: util.NewHashMap[ResultHash, bool](),
	}
}
