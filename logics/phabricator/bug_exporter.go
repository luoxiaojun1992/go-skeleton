package phabricator

import (
	"encoding/json"
	"github.com/lithammer/shortuuid/v3"
	"github.com/luoxiaojun1992/go-skeleton/models"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"github.com/luoxiaojun1992/go-skeleton/services/phabricator"
	"github.com/uber/gonduit/responses"
	"log"
	"time"
)

type BugExporter struct {
}

func (be *BugExporter) FetchUsersInfo(userIds []string) map[string]string {
	var usersInfo map[string]string
	usersInfo = make(map[string]string)

	phabricator.Phabricator.SearchUsers(userIds, 10, func(data []*phabricator.UserSearchResponseItem) {
		for _, user := range data {
			usersInfo[user.PHID] = user.Fields.Username
		}
	})

	return usersInfo
}

func (be *BugExporter) FetchTagsInfo(tagIds []string) map[string]string {
	var tagsInfo map[string]string
	tagsInfo = make(map[string]string)

	phabricator.Phabricator.SearchTags(tagIds, 10, func(data []*responses.ProjectSearchResponseItem) {
		for _, tag := range data {
			tagsInfo[tag.PHID] = tag.Fields.Slug
		}
	})

	return tagsInfo
}

func (be *BugExporter) AddTaskRecords(tasks []*phabricator.ManiphestSearchResponseItem) {
	userIds := []string{}
	tagIds := []string{}

	for _, task := range tasks {
		if len(task.Fields.AuthorPHID) > 0 {
			if len(userIds) < 200 {
				userIds = append(userIds, task.Fields.AuthorPHID)
			} else {
				log.Println("too many phabricator task users")
			}
		}
		if len(task.Fields.OwnerPHID) > 0 {
			if len(userIds) < 200 {
				userIds = append(userIds, task.Fields.OwnerPHID)
			} else {
				log.Println("too many phabricator task users")
			}
		}
		if len(task.Fields.CloserPHID) > 0 {
			if len(userIds) < 200 {
				userIds = append(userIds, task.Fields.CloserPHID)
			} else {
				log.Println("too many phabricator task users")
			}
		}
		if (len(tagIds) + len(task.Attachments.Projects.ProjectPHIDs)) <= 200 {
			tagIds = append(tagIds, task.Attachments.Projects.ProjectPHIDs...)
		} else {
			log.Println("too many phabricator task tags")
		}
	}

	usersInfo := be.FetchUsersInfo(userIds)
	tagsInfo := be.FetchTagsInfo(tagIds)

	var taskRecords []models.PhabricatorBugTasks

	for _, task := range tasks {
		taskRecord := models.PhabricatorBugTasks{
			JingUUID:         shortuuid.New(),
			TaskID:           int64(task.ID),
			TaskPHID:         task.PHID,
			TaskName:         task.Fields.Name,
			TaskStatus:       task.Fields.Status.Value,
			TaskPriority:     float64(task.Fields.Priority.Value),
			TaskSubPriority:  task.Fields.Priority.Subpriority,
			TaskAuthorPHID:   task.Fields.AuthorPHID,
			TaskAuthorName:   "",
			TaskOwnerPHID:    task.Fields.OwnerPHID,
			TaskOwnerName:    "",
			TaskCloserPHID:   task.Fields.CloserPHID,
			TaskCloserName:   "",
			TaskClosed:       0,
			TaskClosedDate:   time.Time{},
			TaskCreated:      0,
			TaskCreatedDate:  time.Time{},
			TaskModified:     0,
			TaskModifiedDate: time.Time{},
			TaskTags:         "",
		}

		if len(task.Fields.AuthorPHID) > 0 {
			authorName, hasAuthorName := usersInfo[task.Fields.AuthorPHID]
			if hasAuthorName {
				taskRecord.TaskAuthorName = authorName
			}
		}

		if len(task.Fields.OwnerPHID) > 0 {
			ownerName, hasOwnerName := usersInfo[task.Fields.OwnerPHID]
			if hasOwnerName {
				taskRecord.TaskOwnerName = ownerName
			}
		}

		if len(task.Fields.CloserPHID) > 0 {
			closerName, hasCloserName := usersInfo[task.Fields.CloserPHID]
			if hasCloserName {
				taskRecord.TaskCloserName = closerName
			}
		}

		if task.Fields.DateCreated != nil {
			dateCreated := time.Time(*task.Fields.DateCreated)
			taskRecord.TaskCreated = dateCreated.Unix()
			taskRecord.TaskCreatedDate = dateCreated
		}

		if task.Fields.DateModified != nil {
			dateModified := time.Time(*task.Fields.DateModified)
			taskRecord.TaskModified = dateModified.Unix()
			taskRecord.TaskModifiedDate = dateModified
		}

		if task.Fields.DateClosed != nil {
			dateClosed := time.Time(*task.Fields.DateClosed)
			taskRecord.TaskClosed = dateClosed.Unix()
			taskRecord.TaskClosedDate = dateClosed
		}

		tagNames := []string{}
		for _, tagId := range task.Attachments.Projects.ProjectPHIDs {
			tagName, hasTagName := tagsInfo[tagId]
			if hasTagName {
				tagNames = append(tagNames, tagName)
			}
		}
		taskTags, errEncodeTags := json.Marshal(&tagNames)
		helper.CheckErrThenPanic("failed to encode tags json", errEncodeTags)
		taskRecord.TaskTags = string(taskTags)

		taskRecords = append(taskRecords, taskRecord)
	}

	errUpsertAll := (&models.PhabricatorBugTasks{}).BatchInsert(taskRecords, true)
	helper.CheckErrThenPanic("failed to upsert tasks", errUpsertAll)
}

func (be *BugExporter) Export(start string, end string) {
	location, errLoadLocation := time.LoadLocation("Asia/Shanghai")
	helper.CheckErrThenPanic("failed to load time location", errLoadLocation)

	modifiedStart, errModifiedStart := time.ParseInLocation("2006-01-02 15:04:05", start, location)
	helper.CheckErrThenPanic("failed to parse task modified start", errModifiedStart)

	modifiedEnd, errModifiedEnd := time.ParseInLocation("2006-01-02 15:04:05", end, location)
	helper.CheckErrThenPanic("failed to parse task modified end", errModifiedEnd)

	tasks := []*phabricator.ManiphestSearchResponseItem{}
	phabricator.Phabricator.SearchLiveBugs(modifiedStart, modifiedEnd, 10, func(data []*phabricator.ManiphestSearchResponseItem) {
		tasks = append(tasks, data...)

		tasksLen := len(tasks)
		if tasksLen >= 10 {
			if tasksLen <= 200 {
				be.AddTaskRecords(tasks)
				tasks = []*phabricator.ManiphestSearchResponseItem{}
			} else {
				log.Panic("Too many phabricator tasks")
			}
		}
	})

	tasksLen := len(tasks)
	if tasksLen > 0 {
		if tasksLen <= 200 {
			be.AddTaskRecords(tasks)
		} else {
			log.Panic("Too many phabricator tasks")
		}
	}
}
