package inventory

import (
	"github.com/ShatteredRealms/character-service/pkg/pb"
	"github.com/ShatteredRealms/gamedata-service/pkg/model/item"
)

type ItemInstance struct {
	Id       uint32    `db:"item_id" json:"id"`
	Slot     item.Slot `json:"slot"`
	Quantity uint32    `json:"quantity"`
}

type Inventory []*ItemInstance

func (i *Inventory) ToRows() *[][]any {
	rows := make([][]any, len(*i))
	for idx, item := range *i {
		rows[idx] = []any{
			item.Id,
			item.Slot,
			item.Quantity,
		}
	}
	return &rows
}

func (i *Inventory) ToPb() []*pb.Item {
	items := make([]*pb.Item, len(*i))
	for idx, item := range *i {
		items[idx] = &pb.Item{
			Id:       item.Id,
			Slot:     uint32(item.Slot),
			Quantity: item.Quantity,
		}
	}
	return items
}
