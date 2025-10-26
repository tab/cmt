package commit

// FetchSuccessMsg indicates successful initial data fetch
type FetchSuccessMsg struct {
	Status  string
	Diff    string
	Message string
}

// FetchErrorMsg indicates initial data fetch failure
type FetchErrorMsg struct {
	Err error
}

// RegenerateMsg triggers commit message regeneration
type RegenerateMsg struct {
	Message string
	Err     error
}

// CommitSuccessMsg indicates successful commit
type CommitSuccessMsg struct {
	Output string
}

// CommitErrorMsg indicates commit failure
type CommitErrorMsg struct {
	Err error
}
