package go_api_abrha

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ImageActionsService is an interface for interfacing with the image actions
// endpoints of the Abrha API
// See: https://docs.parspack.com/api/#tag/Image-Actions
type ImageActionsService interface {
	Get(context.Context, int, int) (*Action, *Response, error)
	GetByURI(context.Context, string) (*Action, *Response, error)
	Transfer(context.Context, int, *ActionRequest) (*Action, *Response, error)
	Convert(context.Context, int) (*Action, *Response, error)
}

// ImageActionsServiceOp handles communication with the image action related methods of the
// Abrha API.
type ImageActionsServiceOp struct {
	client *Client
}

var _ ImageActionsService = &ImageActionsServiceOp{}

// Transfer an image
func (i *ImageActionsServiceOp) Transfer(ctx context.Context, imageID int, transferRequest *ActionRequest) (*Action, *Response, error) {
	if imageID < 1 {
		return nil, nil, NewArgError("imageID", "cannot be less than 1")
	}

	if transferRequest == nil {
		return nil, nil, NewArgError("transferRequest", "cannot be nil")
	}

	path := fmt.Sprintf("api/public/v1/images/%d/actions", imageID)

	req, err := i.client.NewRequest(ctx, http.MethodPost, path, transferRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := i.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

// Convert an image to a snapshot
func (i *ImageActionsServiceOp) Convert(ctx context.Context, imageID int) (*Action, *Response, error) {
	if imageID < 1 {
		return nil, nil, NewArgError("imageID", "cannont be less than 1")
	}

	path := fmt.Sprintf("api/public/v1/images/%d/actions", imageID)

	convertRequest := &ActionRequest{
		"type": "convert",
	}

	req, err := i.client.NewRequest(ctx, http.MethodPost, path, convertRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := i.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

// Get an action for a particular image by id.
func (i *ImageActionsServiceOp) Get(ctx context.Context, imageID, actionID int) (*Action, *Response, error) {
	if imageID < 1 {
		return nil, nil, NewArgError("imageID", "cannot be less than 1")
	}

	if actionID < 1 {
		return nil, nil, NewArgError("actionID", "cannot be less than 1")
	}

	path := fmt.Sprintf("api/public/v1/images/%d/actions/%d", imageID, actionID)
	return i.get(ctx, path)
}

// GetByURI gets an action for a particular image by URI.
func (i *ImageActionsServiceOp) GetByURI(ctx context.Context, rawurl string) (*Action, *Response, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, nil, err
	}

	return i.get(ctx, u.Path)
}

func (i *ImageActionsServiceOp) get(ctx context.Context, path string) (*Action, *Response, error) {
	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := i.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}
