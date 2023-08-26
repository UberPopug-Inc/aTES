package events

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type TaskDataV2 struct {
	TaskID     string `json:"task_id"`
	TaskTitle  string `json:"task_title"`
	JiraID     string `json:"jira_id"`
	AssignUUID string `json:"assign_uuid"`
}

type TaskV2 struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`

	Properties struct {
		EventID      string     `json:"event_id"`
		EventVersion []int      `json:"event_version"`
		EventName    string     `json:"event_name"`
		EventTime    string     `json:"event_time"`
		Data         TaskDataV2 `json:"data"`
	} `json:"properties"`
}

func (t *TaskV2) validate() error {
	sl := gojsonschema.NewReferenceLoader("schemas/events/task/2.json")
	dl := gojsonschema.NewBytesLoader([]byte(t.string()))

	result, err := gojsonschema.Validate(sl, dl)
	if err != nil {
		panic(err.Error())
	}

	if !result.Valid() {
		err := fmt.Errorf("the document is not valid. see errors")

		for _, desc := range result.Errors() {
			err = errors.Join(err, fmt.Errorf("- %s\n", desc))
		}

		return err
	}

	return nil
}

func (t *TaskV2) string() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}

	return string(bytes)
}
