package models

type UpdateDocumentResponse struct {
	MatchedCount  int     `json:"MatchedCount"`
	ModifiedCount int     `json:"ModifiedCount"`
	UpsertedCount int     `json:"UpsertedCount"`
	UpsertedID    *string `json:"UpsertedID"`
}
type DeleteDocumentResponse struct {
	DeletedCount int `json:"DeletedCount"`
}
type InsertDocumentResponse struct {
	Ok  bool   `json:"ok"`
	Id  string `json:"id"`
	Rev string `json:"rev"`
}
type InsertBulkDocumentResponse []InsertDocumentResponse

//	type FindResponse struct {
//		Docs           []any          `json:"docs"`
//		Bookmark       string         `json:"bookmark"`
//		Warning        string         `json:"warning"`
//		ExecutionStats ExecutionStats `json:"execution_stats"`
//	}
type InsertResponse struct {
	InsertedID string `json:"InsertedID"`
}

// type ExecutionStats struct {
// 	TotalKeysExamined       int64   `json:"total_keys_examined"`
// 	TotalDocsExamined       int64   `json:"total_docs_examined"`
// 	TotalQuorumDocsExamined int64   `json:"total_quorum_docs_examined"`
// 	ResultsReturned         int64   `json:"results_returned"`
// 	ExecutionTimeMS         float64 `json:"execution_time_ms"`
// }

// type SecurityModel struct {
// 	Admins  Admins  `json:"admins"`
// 	Members Members `json:"members"`
// }

// type Admins struct {
// 	Names []string      `json:"names"`
// 	Roles []interface{} `json:"roles"`
// }
// type Members struct {
// 	Names []string      `json:"names"`
// 	Roles []interface{} `json:"roles"`
// }
