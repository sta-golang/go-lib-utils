package err

type Error struct {
	Code int `json:"code"`
	Err error `json:"err"`
}

