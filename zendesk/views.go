package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// View represents the views within Zendesk where tickets are grouped and
// observed from.
//  https://developer.zendesk.com/rest_api/docs/support/views#json-format
type View struct {
	ID          int64  `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Active      bool   `json:"active,omitempty"`
	Restriction struct {
		Type string `json:"type,omitempty"`
		ID   int64  `json:"id"`
	} `json:"restriction,omitempty"`
	Position  int64 `json:"position,omitempty"`
	Execution struct {
		GroupBy    string `json:"group_by,omitempty"`
		SortBy     string `json:"sort_by,omitempty"`
		GroupOrder string `json:"group_order,omitempty"`
		SortOrder  string `json:"sort_order,omitempty"`
		Columns    []struct {
			ID    string `json:"-"`
			Title string `json:"-"`
		} `json:"-"`
		Group struct {
			ID    string `json:"-"`
			Title string `json:"-"`
			Order string `json:"-"`
		} `json:"-"`
		Sort struct {
			ID    string `json:"-"`
			Title string `json:"-"`
			Order string `json:"-"`
		} `-"`
	} `json:"execution,omitempty"`
	Conditions struct {
		All []struct {
			Field    string `json:"field,omitempty"`
			Operator string `json:"operator,omitempty"`
			Value    string `json:"value,omitempty"`
		} `json:"all,omitempty"`
		Any []struct {
			Field    string `json:"field,omitempty"`
			Operator string `json:"operator,omitempty"`
			Value    string `json:"value,omitempty"`
		} `json:"any,omitempty"`
	} `json:"conditions,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// ViewCount represents the return from the `count` endpoints
type ViewCount struct {
	ViewID int64  `json:"view_id,omitempty"`
	URL    string `json:"url,omitempty"`
	Value  int64  `json:"value,omitempty"`
	Pretty string `json:"pretty,omitempty"`
	Fresh  bool   `json:"fresh,omitempty"`
}

// ViewAPI is an interface containing all view related methods
type ViewAPI interface {
	GetViews(ctx context.Context) ([]View, Page, error)
	GetActiveViews(ctx context.Context) ([]View, Page, error)
	GetViewCount(ctx context.Context, viewID int) (ViewCount, error)
	GetView(ctx context.Context, viewID int) (View, error)
	CreateView(ctx context.Context, view View) (View, error)
	UpdateView(ctx context.Context, viewID int, view View) (View, error)
}

// GetViews gets a list of all of the current views (active & deactivated)
// Endpoint: GET /api/v2/views.json
// https://developer.zendesk.com/rest_api/docs/support/views#list-views
func (z *Client) GetViews(ctx context.Context) ([]View, Page, error) {
	var data struct {
		Views []View `json:"views"`
		Page
	}
	u, err := addOptions("/views.json", nil)
	if err != nil {
		return nil, Page{}, err
	}
	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.Views, data.Page, nil
}

// GetActiveViews gets a list of all of the current active views
// Endpoint: GET /api/v2/views/active.json
// https://developer.zendesk.com/rest_api/docs/support/views#list-active-views
func (z *Client) GetActiveViews(ctx context.Context) ([]View, Page, error) {
	var data struct {
		Views []View `json:"views"`
		Page
	}
	u, err := addOptions("/views/active.json", nil)
	if err != nil {
		return nil, Page{}, err
	}
	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.Views, data.Page, nil
}

// GetViewCount gets the count of tickets in a given view.
// Endpoint: GET /api/v2/views/{id}/count.json
// https://developer.zendesk.com/rest_api/docs/support/views#get-view-count
func (z *Client) GetViewCount(ctx context.Context, viewID int64) (ViewCount, error) {
	var result struct {
		ViewCount ViewCount `json:"view_count"`
	}

	var builder includeBuilder

	u, err := builder.path(fmt.Sprintf("/views/%d/count.json", viewID))

	if err != nil {
		return ViewCount{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return ViewCount{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ViewCount{}, err
	}
	return result.ViewCount, nil
}

// GetView gets the details of a specified view
// Endpoint: GET /api/v2/views/{ID}.json
// https://developer.zendesk.com/rest_api/docs/support/views#show-view
func (z *Client) GetView(ctx context.Context, viewID int64) (View, error) {
	var result struct {
		View View `json:"view"`
	}

	var builder includeBuilder

	u, err := builder.path(fmt.Sprintf("/views/%d.json", viewID))

	if err != nil {
		return View{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return View{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return View{}, err
	}
	return result.View, nil
}

// CreateView takes a View instance and saves it as a new view in Zendesk
// Endpoint: POST /api/v2/views.json
// https://developer.zendesk.com/rest_api/docs/support/views#create-view
func (z *Client) CreateView(ctx context.Context, view View) (View, error) {
	var data, result struct {
		View View `json:"View"`
	}
	data.View = view

	body, err := z.post(ctx, "/views.json", data)
	if err != nil {
		return View{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return View{}, err
	}
	return result.View, nil
}

// UpdateView takes a View instance and saves it as a new view in Zendesk
// Endpoint: PUT /api/v2/views/{ID}.json
// https://developer.zendesk.com/rest_api/docs/support/views#update-view
func (z *Client) UpdateView(ctx context.Context, viewID int64, view View) (View, error) {
	var data, result struct {
		View View `json:"View"`
	}
	data.View = view
	var builder includeBuilder

	u, err := builder.path(fmt.Sprintf("/views/%d.json", viewID))

	body, err := z.put(ctx, u, data)
	if err != nil {
		return View{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return View{}, err
	}
	return result.View, nil
}
