package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tylerconlee/zendesk-go/zendesk/sideload"
)

type CustomField struct {
	ID int64 `json:"id"`
	// Valid types are string or []string.
	Value interface{} `json:"value"`
}

// Custom Unmarshal function required because a custom field's value can be
// a string or array of strings.
func (cf *CustomField) UnmarshalJSON(data []byte) error {
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	cf.ID = int64(temp["id"].(float64))

	switch v := temp["value"].(type) {
	case string, nil, bool:
		cf.Value = v
	case []interface{}:
		var list []string

		for _, v := range temp["value"].([]interface{}) {
			if s, ok := v.(string); ok {
				list = append(list, s)
			} else {
				return fmt.Errorf("%T is an invalid type for custom field value", v)
			}
		}

		cf.Value = list
	default:
		return fmt.Errorf("%T is an invalid type for custom field value", v)
	}

	return nil
}

type Ticket struct {
	ID              int64         `json:"id,omitempty"`
	URL             string        `json:"url,omitempty"`
	ExternalID      string        `json:"external_id,omitempty"`
	Type            string        `json:"type,omitempty"`
	Subject         string        `json:"subject,omitempty"`
	RawSubject      string        `json:"raw_subject,omitempty"`
	Description     string        `json:"description,omitempty"`
	Priority        string        `json:"priority,omitempty"`
	Status          string        `json:"status,omitempty"`
	Recipient       string        `json:"recipient,omitempty"`
	RequesterID     int64         `json:"requester_id,omitempty"`
	SubmitterID     int64         `json:"submitter_id,omitempty"`
	AssigneeID      int64         `json:"assignee_id,omitempty"`
	OrganizationID  int64         `json:"organization_id,omitempty"`
	GroupID         int64         `json:"group_id,omitempty"`
	CollaboratorIDs []int64       `json:"collaborator_ids,omitempty"`
	FollowerIDs     []int64       `json:"follower_ids,omitempty"`
	EmailCCIDs      []int64       `json:"email_cc_ids,omitempty"`
	ForumTopicID    int64         `json:"forum_topic_id,omitempty"`
	ProblemID       int64         `json:"problem_id,omitempty"`
	HasIncidents    bool          `json:"has_incidents,omitempty"`
	DueAt           time.Time     `json:"due_at,omitempty"`
	Tags            []string      `json:"tags,omitempty"`
	CustomFields    []CustomField `json:"custom_fields,omitempty"`

	// TODO: Via          #123

	SatisfactionRating struct {
		ID      int64  `json:"id"`
		Score   string `json:"score"`
		Comment string `json:"comment"`
	} `json:"satisfaction_rating,omitempty"`

	SharingAgreementIDs []int64   `json:"sharing_agreement_ids,omitempty"`
	FollowupIDs         []int64   `json:"followup_ids,omitempty"`
	ViaFollowupSourceID int64     `json:"via_followup_source_id,omitempty"`
	MacroIDs            []int64   `json:"macro_ids,omitempty"`
	TicketFormID        int64     `json:"ticket_form_id,omitempty"`
	BrandID             int64     `json:"brand_id,omitempty"`
	AllowChannelback    bool      `json:"allow_channelback,omitempty"`
	AllowAttachments    bool      `json:"allow_attachments,omitempty"`
	IsPublic            bool      `json:"is_public,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`

	// Collaborators is POST only
	Collaborators Collaborators `json:"collaborators,omitempty"`

	// Comment is POST only and required
	Comment TicketComment `json:"comment,omitempty"`
	Slas    struct {
		PolicyMetrics []interface{} `json:"policy_metrics,omitempty"`
	} `json:"slas,omitempty"`
	MetricEvents struct {
		PeriodicUpdateTime []struct {
			ID         int64     `json:"id,omitempty"`
			TicketID   int       `json:"ticket_id,omitempty"`
			Metric     string    `json:"metric,omitempty"`
			InstanceID int       `json:"instance_id,omitempty"`
			Type       string    `json:"type,omitempty"`
			Time       time.Time `json:"time,omitempty"`
			Status     struct {
				Calendar int `json:"calendar,omitempty"`
				Business int `json:"business,omitempty"`
			} `json:"status,omitempty"`
		} `json:"periodic_update_time,omitempty"`
		RequesterWaitTime []struct {
			ID         int64     `json:"id,omitempty"`
			TicketID   int       `json:"ticket_id,omitempty"`
			Metric     string    `json:"metric,omitempty"`
			InstanceID int       `json:"instance_id,omitempty"`
			Type       string    `json:"type,omitempty"`
			Time       time.Time `json:"time,omitempty"`
		} `json:"requester_wait_time,omitempty"`
		ResolutionTime []struct {
			ID         int64     `json:"id,omitempty"`
			TicketID   int       `json:"ticket_id,omitempty"`
			Metric     string    `json:"metric,omitempty"`
			InstanceID int       `json:"instance_id",omitempty`
			Type       string    `json:"type,omitempty"`
			Time       time.Time `json:"time,omitempty"`
		} `json:"resolution_time,omitempty"`
		PausableUpdateTime []struct {
			ID         int64     `json:"id,omitempty"`
			TicketID   int       `json:"ticket_id,omitempty"`
			Metric     string    `json:"metric,omitempty"`
			InstanceID int       `json:"instance_id,omitempty"`
			Type       string    `json:"type,omitempty"`
			Time       time.Time `json:"time,omitempty"`
			Status     struct {
				Calendar int `json:"calendar,omitempty"`
				Business int `json:"business,omitempty"`
			} `json:"status,omitempty"`
		} `json:"pausable_update_time,omitempty"`
		AgentWorkTime []struct {
			ID         int64     `json:"id,omitempty"`
			TicketID   int       `json:"ticket_id,omitempty"`
			Metric     string    `json:"metric,omitempty"`
			InstanceID int       `json:"instance_id,omitempty"`
			Type       string    `json:"type,omitempty"`
			Time       time.Time `json:"time,omitempty"`
		} `json:"agent_work_time,omitempty"`
		ReplyTime []struct {
			ID         int64     `json:"id,omitempty"`
			TicketID   int       `json:"ticket_id,omitempty"`
			Metric     string    `json:"metric,omitempty"`
			InstanceID int       `json:"instance_id,omitempty"`
			Type       string    `json:"type,omitempty"`
			Time       time.Time `json:"time,omitempty"`
			SLA        struct {
				Target        int  `json:"target,omitempty"`
				BusinessHours bool `json:"business_hours,omitempty"`
				Policy        struct {
					ID          int         `json:"id,omitempty"`
					Title       string      `json:"title,omitempty"`
					Description interface{} `json:"description,omitempty"`
				} `json:"policy,omitempty"`
			} `json:"sla,omitempty"`
			Deleted bool `json:"deleted,omitempty"`
			Status  struct {
				Calendar int `json:"calendar,omitempty"`
				Business int `json:"business,omitempty"`
			} `json:"status,omitempty"`
		} `json:"reply_time,omitempty"`
	} `json:"metric_events,omitempty"`

	// TODO: TicketAudit (POST only) #126
}

type TicketListOptions struct {
	PageOptions

	// SortBy can take "assignee", "assignee.name", "created_at", "group", "id",
	// "locale", "requester", "requester.name", "status", "subject", "updated_at"
	SortBy string `url:"sort_by,omitempty"`

	// SortOrder can take "asc" or "desc"
	SortOrder string `url:"sort_order,omitempty"`

	// StartTime is a UNIX timestamp of when an incremental export should begin
	StartTime string `url:"start_time,omitempty"`

	// Cursor determines "page" in an incremental export
	// https://developer.zendesk.com/rest_api/docs/support/incremental_export#cursor-based-incremental-exports
	Cursor string `url:"cursor,omitempty"`

	// Sideload includes additional endpoints
	Sideload string `url:"include,omitempty"`
}

// TicketAPI an interface containing all ticket related methods
type TicketAPI interface {
	GetTickets(ctx context.Context, opts *TicketListOptions) ([]Ticket, Page, error)
	GetTicket(ctx context.Context, id int64, sideload ...sideload.SideLoader) (Ticket, error)
	GetMultipleTickets(ctx context.Context, ticketIDs []int64) ([]Ticket, error)
	CreateTicket(ctx context.Context, ticket Ticket) (Ticket, error)
}

// GetTickets get ticket list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tickets#list-tickets
func (z *Client) GetTickets(ctx context.Context, opts *TicketListOptions) ([]Ticket, Page, error) {
	var data struct {
		Tickets []Ticket `json:"tickets"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &TicketListOptions{}
	}

	u, err := addOptions("/tickets.json", tmp)
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
	return data.Tickets, data.Page, nil
}

// GetTickets get ticket list
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tickets#list-tickets
func (z *Client) GetIncrementalTickets(ctx context.Context, opts *TicketListOptions) ([]Ticket, string, bool, error) {
	var data struct {
		Tickets []Ticket `json:"tickets"`
		URL     string   `json:"after_url"`
		EoS     bool     `json:"end_of_stream"`
	}

	tmp := opts
	if tmp == nil {
		tmp = &TicketListOptions{}
	}

	u, err := addOptions("/incremental/tickets.json", tmp)
	if err != nil {
		return nil, "", true, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, "", true, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, "", true, err
	}
	return data.Tickets, data.URL, data.EoS, nil
}

// GetTicket gets a specified ticket
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tickets#show-ticket
func (z *Client) GetTicket(ctx context.Context, ticketID int64, sideLoad ...sideload.SideLoader) (Ticket, error) {
	var result struct {
		Ticket Ticket `json:"ticket"`
	}

	var builder includeBuilder

	for _, v := range sideLoad {
		builder.addKey(v.Key())
	}

	u, err := builder.path(fmt.Sprintf("/tickets/%d.json", ticketID))
	if err != nil {
		return Ticket{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return Ticket{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Ticket{}, err
	}

	for _, sideLoader := range sideLoad {
		err = sideLoader.Unmarshal(body)
		if err != nil {
			return Ticket{}, err
		}
	}

	return result.Ticket, nil
}

// GetMultipleTickets gets multiple specified tickets
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tickets#show-multiple-tickets
func (z *Client) GetMultipleTickets(ctx context.Context, ticketIDs []int64) ([]Ticket, error) {
	var result struct {
		Tickets []Ticket `json:"tickets"`
	}

	var req struct {
		IDs string `url:"ids,omitempty"`
	}
	idStrs := make([]string, len(ticketIDs))
	for i := 0; i < len(ticketIDs); i++ {
		idStrs[i] = strconv.FormatInt(ticketIDs[i], 10)
	}
	req.IDs = strings.Join(idStrs, ",")

	u, err := addOptions("/tickets/show_many.json", req)
	if err != nil {
		return nil, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Tickets, nil
}

// CreateTicket create a new ticket
//
// ref: https://developer.zendesk.com/rest_api/docs/support/tickets#create-ticket
func (z *Client) CreateTicket(ctx context.Context, ticket Ticket) (Ticket, error) {
	var data, result struct {
		Ticket Ticket `json:"ticket"`
	}
	data.Ticket = ticket

	body, err := z.post(ctx, "/tickets.json", data)
	if err != nil {
		return Ticket{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Ticket{}, err
	}
	return result.Ticket, nil
}
