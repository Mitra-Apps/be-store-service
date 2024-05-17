package service

// Code detail for error response
type CodeDetail string

const (
	ErrNotFound                               = "not found"
	ErrIncorrectSqlSyntax                     = "SQLSTATE 42601"
	ErrInvalidColumnName                      = "SQLSTATE 42703"
	ERR_UNKNOWN                    CodeDetail = "ERR_UNKNOWN"
	ERR_PRODUCT_CATEGORY_NOT_FOUND CodeDetail = "ERR_PRODUCT_CATEGORY_NOT_FOUND"
	ERR_PRODUCT_CATEGORY_IS_EXIST  CodeDetail = "ERR_PRODUCT_CATEGORY_IS_EXIST"
	ERR_PRODUCT_NOT_FOUND          CodeDetail = "ERR_PRODUCT_NOT_FOUND"
	ERR_PRODUCT_TYPE_NOT_FOUND     CodeDetail = "ERR_PRODUCT_TYPE_NOT_FOUND"
)
