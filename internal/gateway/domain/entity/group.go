package entity

type Group struct {
	ID        int64
	Name      string
	MemberIDs []int64
}

type CreateGroupDTO struct {
	Name      string  `form:"name"`
	MemberIDs []int64 `form:"member_ids"`
}
