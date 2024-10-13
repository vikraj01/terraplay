package config

type GameConfig struct {
	Name       string
	VolumePath string
}

var AllGameConfigs = []GameConfig{
	{Name: "minecraft", VolumePath: "/opt/minecraft/data"},
	{Name: "csgo", VolumePath: "/opt/csgo/data"},
	{Name: "rust", VolumePath: "/opt/rust/data"},
	{Name: "factorio", VolumePath: "/opt/factorio/data"},
	{Name: "unturned", VolumePath: "/opt/unturned/data"},
	{Name: "terraria", VolumePath: "/opt/terraria/data"},
	{Name: "valheim", VolumePath: "/opt/valheim/data"},
	{Name: "minetest", VolumePath: "/opt/minetest/data"},
	{Name: "openttd", VolumePath: "/opt/openttd/data"},
	{Name: "teeworlds", VolumePath: "/opt/teeworlds/data"},
	{Name: "xonotic", VolumePath: "/opt/xonotic/data"},
	{Name: "freeciv", VolumePath: "/opt/freeciv/data"},
	{Name: "wesnoth", VolumePath: "/opt/wesnoth/data"},
	{Name: "openra", VolumePath: "/opt/openra/data"},
	{Name: "veloren", VolumePath: "/opt/veloren/data"},
	{Name: "hedgewars", VolumePath: "/opt/hedgewars/data"},
}

func (gc *GameConfig) GetVolumePath() string {
	return gc.VolumePath
}

func (gc *GameConfig) GetName() string {
	return gc.Name
}

func FindGameConfig(name string) *GameConfig {
	for _, config := range AllGameConfigs {
		if config.Name == name {
			return &config
		}
	}
	return nil
}

func IsValidGame(name string) bool {
	return FindGameConfig(name) != nil
}

func GetAllGames() []string {
	games := make([]string, len(AllGameConfigs))
	for i, config := range AllGameConfigs {
		games[i] = config.Name
	}
	return games
}
