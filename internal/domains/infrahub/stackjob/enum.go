package stackjob

import (
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/workflow"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
)

// resolveOperationType maps a string (e.g. "update", "destroy") to the
// corresponding StackJobOperationType enum value. Returns a user-friendly
// error listing valid values on mismatch.
func resolveOperationType(s string) (stackjobv1.StackJobOperationType, error) {
	v, ok := stackjobv1.StackJobOperationType_value[s]
	if !ok {
		return 0, fmt.Errorf("unknown stack job operation type %q — valid values: %s",
			s, domains.JoinEnumValues(stackjobv1.StackJobOperationType_value, "stack_job_operation_type_unspecified"))
	}
	return stackjobv1.StackJobOperationType(v), nil
}

// resolveExecutionStatus maps a string (e.g. "running", "completed") to the
// corresponding WorkflowExecutionStatus enum value. Returns a user-friendly
// error listing valid values on mismatch.
func resolveExecutionStatus(s string) (workflow.WorkflowExecutionStatus, error) {
	v, ok := workflow.WorkflowExecutionStatus_value[s]
	if !ok {
		return 0, fmt.Errorf("unknown execution status %q — valid values: %s",
			s, domains.JoinEnumValues(workflow.WorkflowExecutionStatus_value, "workflow_execution_status_unspecified"))
	}
	return workflow.WorkflowExecutionStatus(v), nil
}

// resolveExecutionResult maps a string (e.g. "succeeded", "failed") to the
// corresponding WorkflowExecutionResult enum value. Returns a user-friendly
// error listing valid values on mismatch.
func resolveExecutionResult(s string) (workflow.WorkflowExecutionResult, error) {
	v, ok := workflow.WorkflowExecutionResult_value[s]
	if !ok {
		return 0, fmt.Errorf("unknown execution result %q — valid values: %s",
			s, domains.JoinEnumValues(workflow.WorkflowExecutionResult_value, "workflow_execution_result_unspecified"))
	}
	return workflow.WorkflowExecutionResult(v), nil
}
