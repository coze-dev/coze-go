package coze

import (
	"context"
	"testing"
)

func TestWorkflows_List(t *testing.T) {
	client := NewClient(WithAPIKey("test_key"))
	ctx := context.Background()

	// Test list with pagination
	req := &ListWorkflowsRequest{
		PageSize: 10,
		PageNum:  1,
	}

	paged, err := client.Workflows.List(ctx, req)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	for paged.Next() {
		item := paged.Current()
		if item.ID == "" {
			t.Error("Expected non-empty workflow ID")
		}
	}

	if err := paged.Err(); err != nil {
		t.Errorf("Pagination error: %v", err)
	}
}

func TestWorkflows_List_Filter(t *testing.T) {
	client := NewClient(WithAPIKey("test_key"))
	ctx := context.Background()

	// Test filtering by mode
	req := &ListWorkflowsRequest{
		Mode:     WorkflowModeManual,
		PageSize: 10,
		PageNum:  1,
	}

	paged, err := client.Workflows.List(ctx, req)
	if err != nil {
		t.Fatalf("List with filter failed: %v", err)
	}

	for paged.Next() {
		item := paged.Current()
		if item.Mode != WorkflowModeManual {
			t.Errorf("Expected mode %s, got %s", WorkflowModeManual, item.Mode)
		}
	}

	if err := paged.Err(); err != nil {
		t.Errorf("Pagination error: %v", err)
	}
}