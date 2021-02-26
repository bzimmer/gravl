package rwgps_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/bzimmer/gravl/pkg/providers/activity/rwgps"
	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var u activity.Upload = &rwgps.Upload{TaskID: 2302}
	a.Equal(activity.UploadID(2302), u.Identifier())

	tests := []struct {
		name   string
		done   bool
		upload *rwgps.Upload
	}{
		// no tasks
		{name: "only task id - success: 0", done: false, upload: &rwgps.Upload{Success: 0}},
		{name: "only task id - success: -1", done: true, upload: &rwgps.Upload{Success: -1}},
		{name: "only task id - success: 1", done: true, upload: &rwgps.Upload{Success: 1}},
		// one task
		{name: "one task - success: 0, status: 1", done: true, upload: &rwgps.Upload{
			Success: 0, Tasks: []*rwgps.Task{{Status: 1}}}},
		{name: "one task - success: -1, status: 0", done: false, upload: &rwgps.Upload{
			Success: 0, Tasks: []*rwgps.Task{{Status: 0}}}},
		{name: "one task - success: 1, status: -1", done: true, upload: &rwgps.Upload{
			Success: 0, Tasks: []*rwgps.Task{{Status: -1}}}},
		// more than one task
		{name: "more than one task - success: 0, status: 1,-1", done: true, upload: &rwgps.Upload{
			Success: 0, Tasks: []*rwgps.Task{{Status: 1}, {Status: -1}}}},
		{name: "more than one task - success: -1, status: 0,1", done: false, upload: &rwgps.Upload{
			Success: 0, Tasks: []*rwgps.Task{{Status: 0}, {Status: 1}}}},
		{name: "more than one task - success: 1, status: -1,0", done: false, upload: &rwgps.Upload{
			Success: 0, Tasks: []*rwgps.Task{{Status: -1}, {Status: 0}}}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a.Equal(tt.done, tt.upload.Done())
		})
	}
}
