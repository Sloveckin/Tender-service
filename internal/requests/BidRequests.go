package requests

type BidCreateRequest struct {
	Name            string `json:"name" validate:"required,max=50"`
	Description     string `json:"description" validate:"required"`
	TenderId        string `json:"tenderId" validate:"required,uuid"`
	OrganizationId  string `json:"organizationId" validate:"required,uuid"`
	CreatorUserName string `json:"creatorUsername" validate:"required"`
}
