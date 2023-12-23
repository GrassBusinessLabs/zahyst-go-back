package requests

type AddGroupMemberRequest struct {
	AccessCode string `json:"access_code" validate:"required"`
}

type ChangeMemberAccessLevelRequest struct {
	AccessLevel string `json:"access_level" validate:"required"`
}
