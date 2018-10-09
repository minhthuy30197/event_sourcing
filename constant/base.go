package constant

type PublishChannelStruct struct {
	RegisterRoutes string
	Logs           string
	NewUser        string
}

var PublishChannel = PublishChannelStruct{
	RegisterRoutes: "register-routes-channel",
	Logs:           "logs-channel",
	NewUser:        "new-user-pubsub-channel",
}
