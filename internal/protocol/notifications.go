package protocol

// PushNotificationConfig defines configuration for sending push notifications to an external URL.
type PushNotificationConfig struct {
	ID             *string                             `json:"id,omitempty"`
	URL            string                              `json:"url"`
	Token          *string                             `json:"token,omitempty"`
	Authentication *PushNotificationAuthenticationInfo `json:"authentication,omitempty"`
}

// PushNotificationAuthenticationInfo provides authentication details used in push notification requests.
type PushNotificationAuthenticationInfo struct {
	Schemes     []string `json:"schemes"`
	Credentials *string  `json:"credentials,omitempty"`
}

// TaskPushNotificationConfig links a task with its push notification configuration.
type TaskPushNotificationConfig struct {
	TaskID                 string                 `json:"taskId"`
	PushNotificationConfig PushNotificationConfig `json:"pushNotificationConfig"`
}
