package ui

type TableSearchComponent struct {
	ViewName    string              `json:"viewName"`
	SearchItems []*TableSearchItem  `json:"searchItems"`
	Sorting     *TableSearchSorting `json:"sorting"`
	Paging      *TableSearchPaging  `json:"paging"`
}

type TableSearchItem struct {
	FieldName string `json:"fieldName"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

type TableSearchSorting struct {
	Direction  string   `json:"direction"`
	FieldNames []string `json:"fieldNames"`
}

type TableSearchPaging struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
}

type TableListComponent struct {
	Columns []*TableListColumn `json:"columns"`
	Rows    []*TableListRow    `json:"rows"`
}

type TableListColumn struct {
	FieldName string `json:"fieldName"`
	FieldType string `json:"fieldType"`
	Label     string `json:"label"`
	Index     int    `json:"index"`
}

type TableListField struct {
	FieldName string `json:"fieldName"`
	Text      string `json:"text"`
	Value     string `json:"value"`
}

type TableListRow struct {
	Data []*TableListField `json:"data"`
}

type TableListPaging struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
}
