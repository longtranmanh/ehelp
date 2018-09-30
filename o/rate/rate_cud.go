package rate

import (
	"ehelp/x/rest"
)

func (tok *Rate) CrateRate() *Rate {
	rest.AssertNil(RateTable.Create(tok))
	return tok
}
