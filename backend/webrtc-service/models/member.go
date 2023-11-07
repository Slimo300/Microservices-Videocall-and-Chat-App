package models

type Member struct {
	ID         string
	GroupID    string
	UserID     string
	Username   string
	PictureURL string
	Creator    bool
	Admin      bool
	Muting     bool
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
