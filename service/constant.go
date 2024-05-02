package service

type CodeDetail string

const (
	ErrNotFound = "not found"

	//Insert Error_Detail here
	ERR_RECORD_NOT_FOUND CodeDetail = "ERR_RECORD_NOT_FOUND"
	ERR_INTERNAL_ERROR   CodeDetail = "ERR_INTERNAL_ERROR"
	ERR_RECORD_IS_EXIST  CodeDetail = "ERR_RECORD_IS_EXIST"
)
