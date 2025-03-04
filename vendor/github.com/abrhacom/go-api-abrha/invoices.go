package go_api_abrha

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

const invoicesBasePath = "api/public/v1/customers/my/invoices"

// InvoicesService is an interface for interfacing with the Invoice
// endpoints of the Abrha API
// See: https://docs.parspack.com/api/#tag/Billing
type InvoicesService interface {
	Get(context.Context, string, *ListOptions) (*Invoice, *Response, error)
	GetPDF(context.Context, string) ([]byte, *Response, error)
	GetCSV(context.Context, string) ([]byte, *Response, error)
	List(context.Context, *ListOptions) (*InvoiceList, *Response, error)
	GetSummary(context.Context, string) (*InvoiceSummary, *Response, error)
}

// InvoicesServiceOp handles communication with the Invoice related methods of
// the Abrha API.
type InvoicesServiceOp struct {
	client *Client
}

var _ InvoicesService = &InvoicesServiceOp{}

// Invoice represents a Abrha Invoice
type Invoice struct {
	InvoiceItems []InvoiceItem `json:"invoice_items"`
	Links        *Links        `json:"links"`
	Meta         *Meta         `json:"meta"`
}

// InvoiceItem represents a line-item on a Abrha Invoice
type InvoiceItem struct {
	Product          string    `json:"product"`
	ResourceID       string    `json:"resource_id"`
	ResourceUUID     string    `json:"resource_uuid"`
	GroupDescription string    `json:"group_description"`
	Description      string    `json:"description"`
	Amount           string    `json:"amount"`
	Duration         string    `json:"duration"`
	DurationUnit     string    `json:"duration_unit"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	ProjectName      string    `json:"project_name"`
	Category         string    `json:"category"`
}

// InvoiceList contains a paginated list of all of a customer's invoices.
// The InvoicePreview is the month-to-date usage generated by Abrha.
type InvoiceList struct {
	Invoices       []InvoiceListItem `json:"invoices"`
	InvoicePreview InvoiceListItem   `json:"invoice_preview"`
	Links          *Links            `json:"links"`
	Meta           *Meta             `json:"meta"`
}

// InvoiceListItem contains a small list of information about a customer's invoice.
// More information can be found in the Invoice or InvoiceSummary
type InvoiceListItem struct {
	InvoiceUUID   string    `json:"invoice_uuid"`
	Amount        string    `json:"amount"`
	InvoicePeriod string    `json:"invoice_period"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// InvoiceSummary contains metadata and summarized usage for an invoice generated by Abrha
type InvoiceSummary struct {
	InvoiceUUID           string                  `json:"invoice_uuid"`
	BillingPeriod         string                  `json:"billing_period"`
	Amount                string                  `json:"amount"`
	UserName              string                  `json:"user_name"`
	UserBillingAddress    Address                 `json:"user_billing_address"`
	UserCompany           string                  `json:"user_company"`
	UserEmail             string                  `json:"user_email"`
	ProductCharges        InvoiceSummaryBreakdown `json:"product_charges"`
	Overages              InvoiceSummaryBreakdown `json:"overages"`
	Taxes                 InvoiceSummaryBreakdown `json:"taxes"`
	CreditsAndAdjustments InvoiceSummaryBreakdown `json:"credits_and_adjustments"`
}

// Address represents the billing address of a customer
type Address struct {
	AddressLine1    string    `json:"address_line1"`
	AddressLine2    string    `json:"address_line2"`
	City            string    `json:"city"`
	Region          string    `json:"region"`
	PostalCode      string    `json:"postal_code"`
	CountryISO2Code string    `json:"country_iso2_code"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// InvoiceSummaryBreakdown is a grouped set of InvoiceItems from an invoice
type InvoiceSummaryBreakdown struct {
	Name   string                        `json:"name"`
	Amount string                        `json:"amount"`
	Items  []InvoiceSummaryBreakdownItem `json:"items"`
}

// InvoiceSummaryBreakdownItem further breaks down the InvoiceSummary by product
type InvoiceSummaryBreakdownItem struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
	Count  string `json:"count"`
}

func (i Invoice) String() string {
	return Stringify(i)
}

// Get detailed invoice items for an Invoice
func (s *InvoicesServiceOp) Get(ctx context.Context, invoiceUUID string, opt *ListOptions) (*Invoice, *Response, error) {
	path := fmt.Sprintf("%s/%s", invoicesBasePath, invoiceUUID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Invoice)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root, resp, err
}

// List invoices for a customer
func (s *InvoicesServiceOp) List(ctx context.Context, opt *ListOptions) (*InvoiceList, *Response, error) {
	path := invoicesBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(InvoiceList)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root, resp, err
}

// GetSummary returns a summary of metadata and summarized usage for an Invoice
func (s *InvoicesServiceOp) GetSummary(ctx context.Context, invoiceUUID string) (*InvoiceSummary, *Response, error) {
	path := fmt.Sprintf("%s/%s/summary", invoicesBasePath, invoiceUUID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(InvoiceSummary)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetPDF returns the pdf for an Invoice
func (s *InvoicesServiceOp) GetPDF(ctx context.Context, invoiceUUID string) ([]byte, *Response, error) {
	path := fmt.Sprintf("%s/%s/pdf", invoicesBasePath, invoiceUUID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root bytes.Buffer
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Bytes(), resp, err
}

// GetCSV returns the csv for an Invoice
func (s *InvoicesServiceOp) GetCSV(ctx context.Context, invoiceUUID string) ([]byte, *Response, error) {
	path := fmt.Sprintf("%s/%s/csv", invoicesBasePath, invoiceUUID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root bytes.Buffer
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root.Bytes(), resp, err
}
