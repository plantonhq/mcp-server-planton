package stackjob

import (
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/workflow"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
)

var (
	operationTypeResolver   = domains.NewEnumResolver[stackjobv1.StackJobOperationType](stackjobv1.StackJobOperationType_value, "stack job operation type", "stack_job_operation_type_unspecified")
	executionStatusResolver = domains.NewEnumResolver[workflow.WorkflowExecutionStatus](workflow.WorkflowExecutionStatus_value, "execution status", "workflow_execution_status_unspecified")
	executionResultResolver = domains.NewEnumResolver[workflow.WorkflowExecutionResult](workflow.WorkflowExecutionResult_value, "execution result", "workflow_execution_result_unspecified")
)
