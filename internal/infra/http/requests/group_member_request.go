package requests

type AddGroupMemberRequest struct {
	AccessCode string `json:"access_code" validate:"required"`
}

type ChangeMemberAccessLevelRequest struct {
	AccessLevel string `json:"access_level" validate:"required"`
}

type FindMembersByAreaRequest struct {
	UpperLeftPoint   map[string]float32 `json:"upper_left_point" validate:"required"`
	BottomRightPoint map[string]float32 `json:"bottom_right_point" validate:"required"`
}
