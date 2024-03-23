/*
 * Sample REST Server
 *
 * TODO
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"time"
)

type Task struct {
	// Unique identifier for the task
	Id int32 `json:"id,omitempty"`
	// Description of the task
	Text string `json:"text,omitempty"`
	// Tags associated with the task
	Tags []string `json:"tags,omitempty"`
	// The date the task should be completed by
	Due time.Time `json:"due,omitempty"`
}
