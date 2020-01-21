package tagger_pg

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/selectors"
)

var _ crud.Cleaner = &tagsPg{}

const onSelectToClean = "on tagsPg.SelectToClean(): "

func (taggerOp *tagsPg) SelectToClean(*crud.RemoveOptions) (*selectors.Term, error) {
	return nil, errors.Wrap(common.ErrNotImplemented, onSelectToClean)
}

const onClean = "on tagsPg.Clean(): "

func (taggerOp *tagsPg) Clean(*selectors.Term, *crud.RemoveOptions) error {
	_, err := taggerOp.db.Exec("DELETE FROM " + taggerOp.table)
	if err != nil {
		return errors.Wrap(err, onClean+"can't DELETE FROM "+taggerOp.table)
	}

	return nil
}
