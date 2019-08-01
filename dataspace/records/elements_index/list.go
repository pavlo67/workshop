package elements_index

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/associatio/basis"
	"github.com/pavlo67/associatio/starter/joiner"

	"github.com/pavlo67/associatio/dataspace/records"
)

func All(joinerOp joiner.Operator) ([]Item, basis.Errors) {
	elementsOp := records.Operator(nil)

	var items []Item
	var errs basis.Errors

	for _, component := range joinerOp.ComponentsAllWithInterface(&elementsOp) {
		elementsOp, ok := component.Interface.(records.Operator)
		if ok {
			items = append(items, Item{component.Key, elementsOp})
		} else {
			errs = append(errs, errors.Errorf("incorrect elements.Operator interface (%T) for key '%s'", component.Interface, component.Key))
		}
	}

	return items, errs
}
