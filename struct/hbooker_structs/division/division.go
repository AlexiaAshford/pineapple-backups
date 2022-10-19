package division

var VolumeList = struct {
	Code string       `json:"code"`
	Data DivisionData `json:"data"`
}{}

type DivisionData struct {
	DivisionList []DivisionInfoList `json:"division_list"`
}

type DivisionInfoList struct {
	DivisionID    string `json:"division_id"`
	DivisionIndex string `json:"division_index"`
	DivisionName  string `json:"division_name"`
	Description   string `json:"description"`
}
