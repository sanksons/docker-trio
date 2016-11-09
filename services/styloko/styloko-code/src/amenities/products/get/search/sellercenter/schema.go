package sellercenter

import (
	"fmt"
)

type SellerSearchRequest struct {
	Offset       int
	Limit        int
	SellerIds    []int
	ResetCounter bool
	LastSCId     int
}

func (s SellerSearchRequest) ToString() string {
	return fmt.Sprintf(
		"Offset:%d, Limit:%d, Seller:%v, ResetCounter:%v, LastScId:%v",
		s.Offset, s.Limit, s.SellerIds, s.ResetCounter, s.LastSCId,
	)
}
