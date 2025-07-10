package protocol

// TaskState represents the lifecycle state of a task.
type TaskState string

const (
	// TaskStateSubmitted indicates the task was received by the server and acknowledged,
	// but processing has not yet started.
	TaskStateSubmitted TaskState = "submitted"

	// TaskStateWorking indicates the task is actively being processed by the agent.
	// The client may expect further updates or a terminal state.
	TaskStateWorking TaskState = "working"

	// TaskStateInputRequired indicates the agent requires additional input from the client
	// or user to proceed. The task is effectively paused.
	TaskStateInputRequired TaskState = "input-required"

	// TaskStateCompleted indicates the task finished successfully.
	// Results are typically available in Task.artifacts or TaskStatus.message.
	TaskStateCompleted TaskState = "completed"

	// TaskStateCanceled indicates the task was canceled, for example by a tasks/cancel request
	// or a server-side policy.
	TaskStateCanceled TaskState = "canceled"

	// TaskStateFailed indicates the task terminated due to an error during processing.
	// TaskStatus.message may contain error details.
	TaskStateFailed TaskState = "failed"

	// TaskStateRejected indicates the task was rejected by the remote agent.
	// TaskStatus.message may contain error details.
	TaskStateRejected TaskState = "rejected"

	// TaskStateAuthRequired indicates the agent requires additional authentication to proceed.
	// The task is effectively paused.
	TaskStateAuthRequired TaskState = "auth-required"

	// TaskStateUnknown indicates the state of the task cannot be determined.
	// This may occur if the task ID is invalid, unknown, or expired.
	TaskStateUnknown TaskState = "unknown"
)
