package git_commit

// GraphQLRequest represents a GraphQL API request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// GraphQLError represents a GraphQL API error
type GraphQLError struct {
	Message string `json:"message"`
}

// RepositoryData represents repository information from GraphQL
type RepositoryData struct {
	ID               string        `json:"id"`
	DefaultBranchRef BranchRefData `json:"defaultBranchRef"`
}

// BranchRefData represents branch reference information
type BranchRefData struct {
	Target TargetData `json:"target"`
}

// TargetData represents the target of a branch reference
type TargetData struct {
	OID string `json:"oid"`
}

// CommitData represents commit information from GraphQL
type CommitData struct {
	Commit CommitInfo `json:"commit"`
}

// CommitInfo represents commit details
type CommitInfo struct {
	OID string `json:"oid"`
}
