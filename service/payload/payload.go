package payload

type PageRequestBodyValue struct {
	PageNumber int               `json:"pageNumber"`
	PageSize   int               `json:"pageSize"`
	Sort       map[string]string `json:"sort"`
}

type PageRequestBody struct {
	Request PageRequestBodyValue `json:"request"`
}

type CalculateDebitRequestBodyValue struct {
	IsStatisticsAccordingToCurrentUser bool `json:"isStatisticsAccordingToCurrentUser"`
}

type CalculateDebitRequestBody struct {
	Request CalculateDebitRequestBodyValue `json:"request"`
}

type Response struct {
	Trace        string `json:"trace"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Response     any    `json:"response"`
}

type PageResponse struct {
	Trace        string `json:"trace"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	TotalElement int64  `json:"totalElement"`
	TotalPage    int64  `json:"totalPage"`
	Response     any    `json:"response"`
}
