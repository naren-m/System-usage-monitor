package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	// Kapacitor task endpoint
	kapacitorTaskV1Endpoint = "/kapacitor/v1/tasks"
)

// KapacitorError represents errors of HTTP operations.
type KapacitorError struct {
	Message string
}

var (
	// GET:   404 Not Found
	// PATCH: 404 Not Found
	ErrKapacitorTaskNotExist = KapacitorError{"task does not exist"}
	// POST: 400 Bad Request
	ErrKapacitorTaskAlreadyExist = KapacitorError{"task already exists"}
)

var (
	// kapMutex is used to serialize operations.
	kapMutex sync.Mutex

	// kapDebug is used as development debug.
	kapDebug bool = false
)

// KapacitorDBRP holds database retention policy pairs.
// DB: Database name
// RP: Retention policy name
type KapacitorDBRP struct {
	DB string `json:"db"`
	RP string `json:"rp"`
}

// KapacitorTask holds properties for defining and updating a Kapacitor task.
// If an option is not specified, i.e. nil, it is left unmodified in the case
// of update.
//
// ID:     Unique identifier for the task.
// Type:   Kapacitor task type: "stream" or "batch" (not supported).
// DBRPs:  List of DB retention policy pairs for the task.
// Script: Content of the script.
// Status: One of "enabled" or "disabled".
type KapacitorTask struct {
	ID     string          `json:"id"`
	Type   *string         `json:"type,omitempty"`
	DBRPs  []KapacitorDBRP `json:"dbrps,omitempty"`
	Script *string         `json:"script,omitempty"`
	Status *string         `json:"status,omitempty"`
}

// KapacitorTaskInfo represents Kapacitor task information.
type KapacitorTaskInfo struct {
	Link struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"link"`
	ID string `json:"id"`
	//TemplateID string `json:"template-id"`
	Type  string `json:"type"`
	DBRPs []struct {
		DB string `json:"db"`
		RP string `json:"rp"`
	} `json:"dbrps"`
	//Script string `json:"script"`
	//Vars   struct {
	//} `json:"vars"`
	//Dot       string `json:"dot"`
	Status    string `json:"status"`
	Executing bool   `json:"executing"`
	Error     string `json:"error"`
	//Stats     struct {
	//} `json:"stats"`
	Created     string `json:"created"`
	Modified    string `json:"modified"`
	LastEnabled string `json:"last-enabled"`
}

// Error returns the error message.
func (e KapacitorError) Error() string {
	return "kap: " + e.Message
}

// Dump prints the fields of KapacitorTask.
func (t KapacitorTask) Dump() {
	fmt.Printf("\nID: %s", t.ID)

	if t.Type != nil {
		fmt.Printf("\nType: %s", *(t.Type))
	}

	fmt.Printf("\nDBRPs: %v", t.DBRPs)

	if t.Script != nil {
		fmt.Printf("\nScript: %s", *(t.Script))
	}

	if t.Status != nil {
		fmt.Printf("\nStatus: %s", *(t.Status))
	}
}

// kapDebugPrintf prints only if kapDebug is true.
func kapDebugPrintf(format string, args ...interface{}) {
	if !kapDebug {
		return
	}
	fmt.Printf(format, args...)
}

// KapacitorTaskDefine defines a Kapacitor task in Kapacitor. It is an
// error if the task already exists.
//
// It does POST /kapacitor/v1/tasks.
// Ref: https://docs.influxdata.com/kapacitor/v1.2/api/api/
func KapacitorTaskDefine(task *KapacitorTask) error {
	kapMutex.Lock()
	defer kapMutex.Unlock()

	c, err := ConfKapacitor()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d", c.Host, c.Port) +
		kapacitorTaskV1Endpoint

	kapDebugPrintf("\nPOST: URL: %s\n", url)

	_, err = kapTaskHTTP(http.MethodPost, url, task)

	return err
}

// KapacitorTaskUpdate modifies any property of an already existing task.
//
// Note: Setting any DBRP will overwrite all stored DBRPs.
func KapacitorTaskUpdate(task *KapacitorTask) error {
	kapMutex.Lock()
	defer kapMutex.Unlock()

	c, err := ConfKapacitor()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d", c.Host, c.Port) +
		kapacitorTaskV1Endpoint + "/" + task.ID

	kapDebugPrintf("\nPATCH: URL: %s\n", url)

	_, err = kapTaskHTTP(http.MethodPatch, url, task)

	return err
}

// KapacitorTaskDelete deletes specified task in Kapacitor.
//
// It does DELETE /kapacitor/v1/tasks/TASK_ID.
// HTTP: 204 Success
//
// Note: Deleting a non-existent task is not an error and will return a 204
// succes.
func KapacitorTaskDelete(taskID string) error {
	kapMutex.Lock()
	defer kapMutex.Unlock()

	c, err := ConfKapacitor()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d", c.Host, c.Port) +
		kapacitorTaskV1Endpoint + "/" + taskID

	kapDebugPrintf("\nDELETE: URL: %s\n", url)

	_, err = kapTaskHTTP(http.MethodDelete, url, nil)

	return err
}

// KapacitorTaskGet gets information about a task.
func KapacitorTaskGet(taskID string) (*KapacitorTaskInfo, error) {
	c, err := ConfKapacitor()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%d", c.Host, c.Port) +
		kapacitorTaskV1Endpoint + "/" + taskID

	kapDebugPrintf("\nGET: URL: %s\n", url)

	t, err := kapTaskHTTP(http.MethodGet, url, nil)

	return t, err
}

// kapTaskHTTP is a common function for KapacitorTaskDefine,
// KapacitorTaskUpdate, KapacitorTaskDelete, and KapacitorTaskGet.
//
// method: http.MethodDelete
//         http.MethodGet
//         http.MethodPatch
//         http.MethodPost
//
// task:   Specify nil if http.MethodDelete or http.MethodGet
//
// Ref: https://docs.influxdata.com/kapacitor/v1.2/api/api/
//   2xx The request was a success, content is dependent on the request.
//   4xx Invalid request, refer to error for what it wrong with the request.
//       Repeating the request will continue to return the same error.
//   5xx The server was unable to process the request, refer to the error
//       for a reason. Repeating the request may result in a success if the
//       server issue has been resolved.
func kapTaskHTTP(method, url string, task *KapacitorTask) (*KapacitorTaskInfo,
	error) {
	var reader io.Reader
	var taskInfo KapacitorTaskInfo

	if task != nil {
		t, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return &taskInfo,
				fmt.Errorf("kap: json marshal: %v", err)
		}

		kapDebugPrintf("\nKapacitor task:\n%s\n", t)

		reader = bytes.NewReader(t)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return &taskInfo, fmt.Errorf("kap: http new request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &taskInfo, fmt.Errorf("kap: send: %v", err)
	}
	// https://golang.org/pkg/net/http
	// "The http Client and Transport guarantee that Body is always
	//  non-nil, even on responses without a body or responses with
	//  a zero-length body. It is the caller's responsibility to
	//  close Body."
	defer resp.Body.Close()

	kapDebugPrintf("\nResponse:\n%s\n", resp)

	if resp.StatusCode == http.StatusNoContent {
		// RFC-7231
		// "The 204 (No Content) status code indicates that the server
		//  has successfully fulfilled the request and that there is no
		//  additional content to send in the response payload body."
		//
		// https://docs.influxdata.com/kapacitor/v1.2/api/api/
		// http.StatusNoContent 204 == Success (for DELETE)
		return &taskInfo, nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &taskInfo, fmt.Errorf("kap: read: %v", err)
	}

	if err := json.Unmarshal(b, &taskInfo); err != nil {
		return &taskInfo, fmt.Errorf("kap: json unmarshal: %v", err)
	}

	kapDebugPrintf("\nTask info:\n%+v\n", taskInfo)

	// 2xx is success. http.StatusMultipleChoices is 300.
	if resp.StatusCode >= http.StatusMultipleChoices {
		err = kapTaskParseHTTPError(method, resp.StatusCode,
			taskInfo.Error)
		return &taskInfo, err
	}

	return &taskInfo, nil
}

// kapTaskParseHTTPError is a convenience function for kapTaskHTTP.
func kapTaskParseHTTPError(method string, statusCode int,
	errString string) error {
	if strings.Contains(errString, "already exists") {
		// POST: "task <task ID> already exists"
		return ErrKapacitorTaskAlreadyExist
	} else if strings.Contains(errString,
		"task does not exist") ||
		strings.Contains(errString, "no task exists") {
		// PATCH: "task does not exist, cannot update"
		// GET: "no task exists"
		return ErrKapacitorTaskNotExist
	}
	return fmt.Errorf("kap: %s: http code: %d error: %s",
		method, statusCode, errString)
}
