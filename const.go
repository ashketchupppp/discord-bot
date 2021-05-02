package main

const (
	ConfigPathFlagName = "configPath"

	DefaultChangeMe               = "CHANGEME"
	DefaultBotToken               = DefaultChangeMe
	DefaultConfigPath             = "./.bot.conf.json"
	DefaultMongoDBName            = "discordbot"
	DefaultMongoDBQuoteCollection = "quotes"
	DefaultMongoDBConnStr         = DefaultChangeMe
	DefaultQuoteChannel           = DefaultChangeMe
	DefaultLeaveMessageChannel    = DefaultChangeMe
	DefaultCommandPrefix          = "$"

	HelpCmdName     = "help"
	AddQuoteCmdName = "addquote"
	GetQuoteCmdName = "getquote"

	LeaveMessageFeatureName = "leavemessage"

	QuoteChannelSettingName = "quotechannel"
	LeaveChannelSettingName = "leavemessagechannel"
)
