package directory

type Badge struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
}

type badgeResponse struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
}

func (response badgeResponse) Badge() Badge {
	return Badge{
		UserID:      response.UserID,
		DisplayName: response.DisplayName,
		Avatar:      response.Avatar,
	}
}
