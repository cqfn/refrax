package protocol

type PushNotificationConfig struct {
	ID             *string                             `json:"id,omitempty"`             // Optional
	URL            string                              `json:"url"`                      // Required
	Token          *string                             `json:"token,omitempty"`          // Optional
	Authentication *PushNotificationAuthenticationInfo `json:"authentication,omitempty"` // Optional
}

type PushNotificationAuthenticationInfo struct {
	Schemes     []string `json:"schemes"`               // Required
	Credentials *string  `json:"credentials,omitempty"` // Optional
}

type TaskPushNotificationConfig struct {
	TaskID                 string                 `json:"taskId"`                 // Required: task to configure or query
	PushNotificationConfig PushNotificationConfig `json:"pushNotificationConfig"` // Required: config to set or return
}
