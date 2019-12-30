package packs

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/types"
)

const InterfaceKey joiner.InterfaceKey = "packs"
const CollectionDefault = "packs"

type Pack struct {
	IdentityKey identity.Key

	From    identity.Key   `bson:",omitempty" json:",omitempty"`
	To      []identity.Key `bson:",omitempty" json:",omitempty"`
	Options common.Map     `bson:",omitempty" json:",omitempty"`

	TypeKey    types.Key   `bson:",omitempty" json:",omitempty"`
	ContentRaw []byte      `bson:"-"          json:"-"`
	Content    interface{} `bson:",omitempty" json:",omitempty"`

	History []crud.Action `bson:",omitempty" json:",omitempty"`
}

type Item struct {
	ID   common.ID `bson:"_id,omitempty" json:",omitempty"`
	Pack `          bson:",inline" json:",inline"`
}

type Operator interface {
	Save(Pack, *crud.SaveOptions) (common.ID, error)
	Remove(common.ID, *crud.RemoveOptions) error
	Read(common.ID, *crud.GetOptions) (*Item, error)
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)

	AddHistory(common.ID, []crud.Action, *crud.SaveOptions) error
}