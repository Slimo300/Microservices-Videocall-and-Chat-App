package models

type Member struct {
	ID       string `mapstructure:"ID"`
	GroupID  string `mapstructure:"groupID"`
	UserID   string `mapstructure:"userID"`
	Username string `mapstructure:"username"`
	// PictureURL string `mapstructure:"pictureURL"`
	Creator bool `mapstructure:"creator"`
	Admin   bool `mapstructure:"admin"`
	Muting  bool `mapstructure:"muting"`
}

func (m *Member) CanMute(mem *Member) bool {
	if m.ID == mem.ID {
		return false
	}
	if m.Creator {
		return true
	}
	if m.Admin && !mem.Creator {
		return true
	}
	if m.Muting && !mem.Creator && !mem.Admin {
		return true
	}
	return false
}
