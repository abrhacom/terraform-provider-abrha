package go_api_abrha

import (
	"context"
	"net/http"
	"time"
)

// BalanceService is an interface for interfacing with the Balance
// endpoints of the Abrha API
// See: https://docs.parspack.com/api/#operation/balance_get
type BalanceService interface {
	Get(context.Context) (*Balance, *Response, error)
}

// BalanceServiceOp handles communication with the Balance related methods of
// the Abrha API.
type BalanceServiceOp struct {
	client *Client
}

var _ BalanceService = &BalanceServiceOp{}

// Balance represents a Abrha Balance
type Balance struct {
	MonthToDateBalance string    `json:"month_to_date_balance"`
	AccountBalance     string    `json:"account_balance"`
	MonthToDateUsage   string    `json:"month_to_date_usage"`
	GeneratedAt        time.Time `json:"generated_at"`
}

func (r Balance) String() string {
	return Stringify(r)
}

// Get Abrha balance info
func (s *BalanceServiceOp) Get(ctx context.Context) (*Balance, *Response, error) {
	path := "api/public/v1/customers/my/balance"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Balance)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
