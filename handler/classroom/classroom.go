package classroom

type ListResponse struct {
	Count int              `json:"count"`
	List  []*ClassroomItem `json:"list"`
}

type ClassroomItem struct {
	Time  int8     `json:"time"` // 节次
	Rooms []string `json:"rooms"`
}

/*
{
	"count" : 12,
	"list": [
		{
			"time": 1,
			"rooms": []
		},
		...
	]
}
*/
