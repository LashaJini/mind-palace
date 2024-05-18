package config

type User struct {
	UserConfig UserConfig `json:"config"`
}

type UserConfig struct {
	User     string    `json:"user"`
	Commands []Command `json:"commands"`
}

type Command struct {
	Name          string   `json:"name"`
	ResourceTypes []string `json:"resourceTypes"`
	Steps         []string `json:"steps"`
}

func NewUserConfig(user string) *User {
	return &User{
		UserConfig: UserConfig{
			User: user,
		},
	}
}
