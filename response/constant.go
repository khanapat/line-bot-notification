package response

const _externalErrorCode = 4000
const _internalErrorCode = 5000

const (
	SuccessCode uint64 = 2000

	ErrInvalidRequestCode uint64 = _externalErrorCode + 1

	ErrDatabaseCode   uint64 = _internalErrorCode + 1
	ErrThirdPartyCode uint64 = _internalErrorCode + 2
)

const (
	SuccessMessageEn string = "Success."

	ErrInvalidRequestMessageEn string = "Input request error."

	ErrDatabaseMessageEn   string = "Database error."
	ErrThirdPartyMessageEn string = "Third party response error."
)
