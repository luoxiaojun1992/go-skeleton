package phabricator

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"github.com/uber/gonduit"
	"github.com/uber/gonduit/core"
	"github.com/uber/gonduit/entities"
	"github.com/uber/gonduit/requests"
	"github.com/uber/gonduit/responses"
	"strconv"
	"strings"
	"time"
)

// ManiphestSearchResponseItem contains information about a particular search result.
type ManiphestSearchResponseItem struct {
	// ID is task identifier.
	ID int `json:"id"`
	// Type is task type.
	Type string `json:"type"`
	// PHID is PHID of the task.
	PHID string `json:"phid"`
	// Fields contains task data.
	Fields struct {
		// Name is task name.
		Name string `json:"name"`
		// Description is detailed task description.
		Description *responses.TaskDescription `json:"description"`
		// AuthorPHID is PHID of task submitter.
		AuthorPHID string `json:"authorPHID"`
		// OwnerPHID is PHID of the person who currently assigned to task.
		OwnerPHID string `json:"ownerPHID"`
		// Status is task status.
		Status responses.ManiphestSearchResultStatus `json:"status"`
		// Priority is task priority.
		Priority responses.ManiphestSearchResultPriority `json:"priority"`
		// Points is point value of the task.
		Points json.Number `json:"points"`
		// Subtype of the task.
		Subtype string `json:"subtype"`
		// CloserPHID is user who closed the task, if the task is closed.
		CloserPHID string `json:"closerPHID"`
		// SpacePHID is PHID of the policy space this object is part of.
		SpacePHID string `json:"spacePHID"`
		// Date created is epoch timestamp when the object was created.
		DateCreated *UnixTimestamp `json:"dateCreated"`
		// DateModified is epoch timestamp when the object was last updated.
		DateModified *UnixTimestamp `json:"dateModified"`
		// DateClosed is epoch timestamp when the object was closed.
		DateClosed *UnixTimestamp `json:"dateClosed"`
		// Policy is map of capabilities to current policies.
		Policy responses.SearchResultPolicy `json:"policy"`
		// CustomTaskType is custom task type.
		CustomTaskType string `json:"custom.task_type"`
		// CustomSeverity is task severity custom value.
		CustomSeverity string `json:"custom.severity"`
	} `json:"fields"`
	Attachments struct {
		// Columns contains columnt data if requested.
		Columns struct {
			// Boards is ???.
			Boards *responses.ManiphestSearchAttachmentColumnBoards `json:"boards"`
		} `json:"columns"`
		// Subscribers contains subscribers attachment data.
		Subscribers struct {
			// SubscriberPHIDs is a collection of PHIDs of persons subscribed to a task.
			SubscriberPHIDs []string `json:"subscriberPHIDs"`
			// SubscriberCount is number of subscribers.
			SubscriberCount int `json:"subscriberCount"`
			// ViewerIsSubscribed specifies if request is subscribed to this task.
			ViewerIsSubscribed bool `json:"viewerIsSubscribed"`
		} `json:"subscribers"`
		// Projects contains project attachment data.
		Projects struct {
			// ProjectPHIDs is collection of PHIDs of projects that this task is tagged with.
			ProjectPHIDs []string `json:"projectPHIDs"`
		} `json:"projects"`
	} `json:"attachments"`
}

// ManiphestSearchResponse contains fields that are in server response to maniphest.search.
type ManiphestSearchResponse struct {
	// Data contains search results.
	Data []*ManiphestSearchResponseItem `json:"data"`
	// Curson contains paging data.
	Cursor struct {
		Limit  uint64 `json:"limit"`
		After  string `json:"after"`
		Before string `json:"before"`
	} `json:"cursor,omitempty"`
}

// UnixTimestamp represents a
type UnixTimestamp time.Time

// MarshalJSON implements the json.Marshaler interface.
func (t UnixTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).Unix())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *UnixTimestamp) UnmarshalJSON(data []byte) (err error) {
	seconds, err := strconv.Atoi(
		strings.Trim(string(data), "\\\""),
	)
	if err != nil {
		return err
	}

	*t = UnixTimestamp(time.Unix(int64(seconds), 0))

	return nil
}

// ManiphestSearchConstraints describes search criteria for request.
type ManiphestSearchConstraints struct {
	// IDs - search for objects with specific IDs.
	IDs []int `json:"ids,omitempty"`
	// PHIDs - search for objects with specific PHIDs.
	PHIDs []string `json:"phids,omitempty"`
	// AssignedTo - search for tasks owned by a user from a list.
	AssignedTo []string `json:"assigned,omitempty"`
	// Authors - search for tasks with given authors.
	Authors []string `json:"authorPHIDs,omitempty"`
	// Statuses - search for tasks with given statuses.
	Statuses []string `json:"statuses,omitempty"`
	// Priorities - search for tasks with given priorities.
	Priorities []int `json:"priorities,omitempty"`
	// Subtypes - search for tasks with given subtypes.
	Subtypes []string `json:"subtypes,omitempty"`
	// Column PHIDs ??? - no doc on phab site.
	ColumnPHIDs []string `json:"columnPHIDs,omitempty"`
	// OpenParents - search for tasks that have parents in open state.
	OpenParents *bool `json:"hasParents,omitempty"`
	// OpenSubtasks - search for tasks that have child tasks in open state.
	OpenSubtasks *bool `json:"hasSubtasks,omitempty"`
	// ParentIDs - search for children of these parents.
	ParentIDs []int `json:"parentIDs,omitempty"`
	// SubtaskIDs - Search for tasks that have these children.
	SubtaskIDs []int `json:"subtaskIDs,omitempty"`
	// CreatedAfter - search for tasks created after given date.
	CreatedAfter *UnixTimestamp `json:"createdStart,omitempty"`
	// CreatedBefore - search for tasks created before given date.
	CreatedBefore *UnixTimestamp `json:"createdEnd,omitempty"`
	// ModifiedAfter - search for tasks modified after given date.
	ModifiedAfter *UnixTimestamp `json:"modifiedStart,omitempty"`
	// ModifiedBefore - search for tasks modified before given date.
	ModifiedBefore *UnixTimestamp `json:"modifiedEnd,omitempty"`
	// ClosedAfter - search for tasks closed after given date.
	ClosedAfter *UnixTimestamp `json:"closedStart,omitempty"`
	// ClosedBefore - search for tasks closed before given date.
	ClosedBefore *UnixTimestamp `json:"closedEnd,omitempty"`
	// ClosedBy - search for tasks closed by people with given PHIDs.
	ClosedBy []string `json:"closerPHIDs,omitempty"`
	// Query - find objects matching a fulltext search query.
	Query string `json:"query,omitempty"`
	// Subscribers - search for objects with certain subscribers.
	Subscribers []string `json:"subscribers,omitempty"`
	// Projects - search for objects tagged with given projects.
	Projects []string `json:"projects,omitempty"`
	// Spaces - search for objects in certain spaces.
	Spaces []string `json:"spaces,omitempty"`
}

// ManiphestSearchRequest represents a request to maniphest.search API method.
type ManiphestSearchRequest struct {
	// QueryKey is builtin or saved query to use. It is optional and sets initial constraints.
	QueryKey string `json:"queryKey,omitempty"`
	// Constraints contains additional filters for results. Applied on top of query if provided.
	Constraints *ManiphestSearchConstraints `json:"constraints,omitempty"`
	// Attachments specified what additional data should be returned with each result.
	Attachments *requests.ManiphestSearchAttachments `json:"attachments,omitempty"`

	*entities.Cursor
	requests.Request
}

type UserSearchConstraints struct {
	PHIDs []string `json:"phids,omitempty"`
}

type UserSearchRequest struct {
	Constraints *UserSearchConstraints `json:"constraints,omitempty"`

	*entities.Cursor
	requests.Request
}

type UserSearchResponseItem struct {
	PHID   string `json:"phid"`
	Fields struct {
		Username string `json:"username"`
	} `json:"fields"`
}

type UserSearchResponse struct {
	// Data contains search results.
	Data []*UserSearchResponseItem `json:"data"`
	// Curson contains paging data.
	Cursor struct {
		Limit  uint64 `json:"limit"`
		After  string `json:"after"`
		Before string `json:"before"`
	} `json:"cursor,omitempty"`
}

var Phabricator *Client

type Client struct {
	ApiGateway string
	ApiToken   string
}

func Setup() {
	Phabricator = &Client{
		ApiGateway: config.String("phabricator.api_gateway"),
		ApiToken:   config.String("phabricator.api_token"),
	}
}

func (c *Client) Connect() *gonduit.Conn {
	apiClient, errDial := gonduit.Dial(c.ApiGateway, &core.ClientOptions{APIToken: c.ApiToken, Timeout: 5 * time.Second})
	helper.CheckErrThenPanic("failed to dial phabricator api gateway", errDial)

	return apiClient
}

func (c *Client) SearchUsers(userIds []string, batch uint64, callback func(data []*UserSearchResponseItem)) {
	start := 0
	end := batch
	userIdsLen := uint64(len(userIds))
	if end > userIdsLen {
		end = userIdsLen
	}

	for true {
		var res UserSearchResponse
		apiClient := c.Connect()
		doSearch := func() error {
			return apiClient.Call("user.search", &UserSearchRequest{
				Constraints: &UserSearchConstraints{PHIDs: userIds[start:end]},
				Cursor: &entities.Cursor{
					Limit: batch,
				},
			}, &res)
		}
		errSearch := doSearch()
		if helper.CheckErr(errSearch) {
			errSearch := doSearch()
			helper.CheckErrThenPanic("failed to search user of phabricator", errSearch)
		}

		response := &res
		if len(response.Data) > 0 {
			callback(response.Data)
		}

		start = start + int(batch)
		end = uint64(start) + batch

		if uint64(start) >= userIdsLen {
			break
		}
		if end > userIdsLen {
			end = userIdsLen
		}
	}
}

func (c *Client) SearchTags(tagIds []string, batch uint64, callback func(data []*responses.ProjectSearchResponseItem)) {
	start := 0
	end := batch
	tagIdsLen := uint64(len(tagIds))
	if end > tagIdsLen {
		end = tagIdsLen
	}

	for true {
		apiClient := c.Connect()
		doSearch := func() (*responses.ProjectSearchResponse, error) {
			return apiClient.ProjectSearch(requests.ProjectSearchRequest{
				Constraints: &requests.ProjectSearchConstraints{PHIDs: tagIds[start:end]},
				Cursor: &entities.Cursor{
					Limit: batch,
				},
			})
		}
		response, errSearch := doSearch()
		if helper.CheckErr(errSearch) {
			response, errSearch = doSearch()
			helper.CheckErrThenPanic("failed to search tag of phabricator", errSearch)
		}

		if len(response.Data) > 0 {
			callback(response.Data)
		}

		start = start + int(batch)
		end = uint64(start) + batch

		if uint64(start) >= tagIdsLen {
			break
		}
		if end > tagIdsLen {
			end = tagIdsLen
		}
	}
}

func (c *Client) SearchLiveBugs(start time.Time, end time.Time, batch uint64, callback func(data []*ManiphestSearchResponseItem)) {
	modifiedAfter := UnixTimestamp(start)
	modifiedBefore := UnixTimestamp(end)

	var cursorAfter uint64
	cursorAfter = 0
	var cursorBefore uint64
	cursorBefore = 0

	for true {
		apiClient := c.Connect()
		var res ManiphestSearchResponse
		cursor := &entities.Cursor{
			Limit: batch,
		}
		if cursorAfter > 0 {
			cursor.After = cursorAfter
		}
		if cursorBefore > 0 {
			cursor.Before = cursorBefore
		}
		doSearch := func() error {
			return apiClient.Call("maniphest.search", &ManiphestSearchRequest{
				Constraints: &ManiphestSearchConstraints{
					ModifiedAfter:  &modifiedAfter,
					ModifiedBefore: &modifiedBefore,
					Projects:       []string{"bug"},
				},
				Attachments: &requests.ManiphestSearchAttachments{
					Projects: true,
				},
				Cursor: cursor,
			}, &res)
		}
		errSearch := doSearch()
		if helper.CheckErr(errSearch) {
			errSearch := doSearch()
			helper.CheckErrThenPanic("failed to search maniphest of phabricator", errSearch)
		}

		response := &res
		if len(response.Data) <= 0 {
			break
		}

		callback(response.Data)

		resCursor := response.Cursor
		if len(resCursor.After) <= 0 {
			break
		} else {
			cursorAfterInt, errCursorAfterInt := strconv.Atoi(resCursor.After)
			helper.CheckErrThenPanic("failed to parse int cursor after", errCursorAfterInt)
			cursorAfter = uint64(cursorAfterInt)
		}
		if len(resCursor.Before) > 0 {
			cursorBeforeInt, errCursorBeforeInt := strconv.Atoi(resCursor.Before)
			helper.CheckErrThenPanic("failed to parse int cursor before", errCursorBeforeInt)
			cursorBefore = uint64(cursorBeforeInt)
		}
	}
}
