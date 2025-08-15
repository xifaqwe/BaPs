package game

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetNoRefreshShopProductList(s *enter.Session, categoryType proto.ShopCategoryType) []*proto.ShopProductDB {
	list := make([]*proto.ShopProductDB, 0)
	products := gdconf.GetShopExcelType(categoryType.String())
	if products == nil {
		return list
	}
	
	for _, product := range products {
		if categoryType == proto.ShopCategoryType_SecretStone {
			if GetCharacterInfo(s, product.Id) == nil {
				continue
			}
		}
		list = append(list, &proto.ShopProductDB{
			EventContentId:     0,
			ShopExcelId:        product.Id,
			Category:           categoryType,
			DisplayOrder:       product.DisplayOrder,
			PurchaseCount:      0,
			PurchaseCountLimit: product.PurchaseCountLimit,
			Price:              0,
			ProductType:        proto.ShopProductType_General,
		})
	}
	return list
}

func GetRefreshShopProductList(categoryType proto.ShopCategoryType) []*proto.ShopProductDB {
	list := make([]*proto.ShopProductDB, 0)
	products := gdconf.GetShopRefreshExcelMap(categoryType.String())
	if products == nil {
		return list
	}
	
	for _, product := range products {
		list = append(list, &proto.ShopProductDB{
			EventContentId:     0,
			ShopExcelId:        product.Id,
			Category:           categoryType,
			DisplayOrder:       product.DisplayOrder,
			PurchaseCount:      0,
			PurchaseCountLimit: 1,
			Price:              0,
			ProductType:        proto.ShopProductType_Refresh,
		})
	}
	return list
}
