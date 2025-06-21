package protocol

type TaskState string

const (
	// Task received by the server and acknowledged, but processing has not yet actively started.
	TaskStateSubmitted TaskState = "submitted"

	// Task is actively being processed by the agent.
	// Client may expect further updates or a terminal state.
	TaskStateWorking TaskState = "working"

	// Agent requires additional input from the client/user to proceed.
	// The task is effectively paused.
	TaskStateInputRequired TaskState = "input-required"

	// Task finished successfully.
	// Results are typically available in Task.artifacts or TaskStatus.message.
	TaskStateCompleted TaskState = "completed"

	// Task was canceled (e.g., by a tasks/cancel request or server-side policy).
	TaskStateCanceled TaskState = "canceled"

	// Task terminated due to an error during processing.
	// TaskStatus.message may contain error details.
	TaskStateFailed TaskState = "failed"

	// Task terminated due to rejection by remote agent.
	// TaskStatus.message may contain error details.
	TaskStateRejected TaskState = "rejected"

	// Agent requires additional authentication from the client/user to proceed.
	// The task is effectively paused.
	TaskStateAuthRequired TaskState = "auth-required"

	// TaskStateUnknown:
	// The state of the task cannot be determined (e.g., task ID is invalid, unknown, or has expired).
	TaskStateUnknown TaskState = "unknown"
)
