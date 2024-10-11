package config

type Command string

const (
	CreateGame  Command = "create-game"
	StopGame    Command = "stop-game"
	RestartGame Command = "restart-game"
	BackupGame  Command = "backup-game"
	RestoreGame Command = "restore-game"
	DeleteGame  Command = "delete-game"

	ListSessions Command = "list-sessions"
	ListGames    Command = "list-games"
	GetLogs      Command = "get-logs"
	GameStatus   Command = "game-status"
	Expenses         Command = "expenses"

	SessionSettings  Command = "session-settings"
	ScaleGame        Command = "scale-game"
	UpdateGame       Command = "update-game"
	ChangeGameConfig Command = "change-game-config"

	MonitorResources Command = "monitor-resources"
	FetchMetrics     Command = "fetch-metrics"
	SetNotifications Command = "set-notifications"
	ListBackups      Command = "list-backups"
	RestorePoint     Command = "restore-point"

	EnableMaintenance  Command = "enable-maintenance"
	DisableMaintenance Command = "disable-maintenance"

	Login  Command = "login"
	Logout Command = "logout"
)

var ValidCommands = []Command{
	CreateGame, StopGame, RestartGame, BackupGame, RestoreGame, DeleteGame,
	ListSessions, ListGames, GetLogs, GameStatus,
	SessionSettings, Expenses, ScaleGame, UpdateGame, ChangeGameConfig,
	MonitorResources, FetchMetrics, SetNotifications, ListBackups, RestorePoint,
	EnableMaintenance, DisableMaintenance,
	Login, Logout,
}
