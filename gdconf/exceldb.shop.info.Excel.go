package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadShopInfoExcel() {
	g.GetExcel().ShopInfoExcel = make([]*sro.ShopInfoExcel, 0)
	name := "ShopInfoExcel.json"
	mx.LoadExcelJson(g.excelDbPath+name, &g.GetExcel().ShopInfoExcel)
}

type ShopInfoExcel struct {
	ShopInfoExcelMap map[string]*sro.ShopInfoExcel
}

func (g *GameConfig) gppShopInfoExcel() {
	g.GetGPP().ShopInfoExcel = &ShopInfoExcel{
		ShopInfoExcelMap: make(map[string]*sro.ShopInfoExcel),
	}
	for _, v := range g.GetExcel().GetShopInfoExcel() {
		g.GetGPP().ShopInfoExcel.ShopInfoExcelMap[v.CategoryType] = v
	}
	logger.Info("处理商店配置完成,商店类型:%v个", len(g.GetGPP().ShopInfoExcel.ShopInfoExcelMap))
}

func GetShopInfoExcel(categoryType string) *sro.ShopInfoExcel {
	return GC.GetGPP().ShopInfoExcel.ShopInfoExcelMap[categoryType]
}
