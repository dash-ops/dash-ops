package models

// BaseResponse provides a common structure for API responses
type BaseResponse struct {
	Success  bool                   `json:"success"`
	Message  string                 `json:"message,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// BaseMetadata provides common metadata fields
type BaseMetadata struct {
	Total      int  `json:"total,omitempty"`
	HasMore    bool `json:"has_more,omitempty"`
	NextOffset int  `json:"next_offset,omitempty"`
}
