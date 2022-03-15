package project

type Project struct {
	Id         uint64  `json:"id"`
	Name       string  `json:"name"`
	ThemeColor string  `json:"theme_color"`
	ParentId   *uint64 `json:"parent_id,omitempty"`
	Pinned     bool    `json:"pinned"`
	Hidden     bool    `json:"hidden"`
}
