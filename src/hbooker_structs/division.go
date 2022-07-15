package hbooker_structs

type DivisionStruct struct {
	Code string       `json:"code"`
	Data DivisionData `json:"data"`
}

type DivisionData struct {
	DivisionList []DivisionList `json:"division_list"`
}

type DivisionList struct {
	DivisionID    string `json:"division_id"`
	DivisionIndex string `json:"division_index"`
	DivisionName  string `json:"division_name"`
	Description   string `json:"description"`
}
