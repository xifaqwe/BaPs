package BaPs

import "github.com/gucooing/BaPs/pkg/logger"

/*
本文件是 为关爱孤儿 基于人道主义 设立的孤儿信息表 欢迎提交
*/

type orphan struct {
	OrphanName    string // 孤儿名称
	Gender        string // 性别
	Address       string // 住址
	OrphanHomeUrl string // 孤儿主页
	OrphanLife    string // 孤儿生平
}

var orphanList = []*orphan{
	{
		OrphanName:    "网络昵称:跨次元的定格的亲叠（来源:闲鱼app） 现实昵称: 不知道也不想知道",
		Gender:        "网络性别:女（来源:闲鱼app） 现实性别:男（来源:支付宝公开数据）",
		Address:       "网络住址:闲鱼app 现实住址:中国江苏",
		OrphanHomeUrl: "https://m.tb.cn/h.6EZTx3u",
		OrphanLife:    "",
	},
}

func logOrphan() {
	logger.Info("--------------------孤儿介绍--------------------\n\n")
	for k, v := range orphanList {
		logger.Warn("孤儿%v号", k+1)
		logger.Error("孤儿名称:%s", v.OrphanName)
		logger.Error("孤儿性别:%s", v.Gender)
		logger.Error("孤儿住址:%s", v.Address)
		logger.Error("孤儿主页:%s", v.OrphanHomeUrl)
		logger.Error("孤儿生平事迹:%s", v.OrphanLife)
	}
	logger.Info("---------------------------------------------")
	logger.Warn("各位爱心人士可以选择相中的孤儿进行帮助,同时我们也非常愿意收录更多孤儿\n\n")
	logger.Info("------------------孤儿介绍结束------------------")
}
