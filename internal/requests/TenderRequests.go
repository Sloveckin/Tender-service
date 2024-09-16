package requests

type CreateTenderRequest struct {
	Name            string `json:"name" validate:"required,max=50"`
	Description     string `json:"description" validate:"required"`
	Status          string `json:"status" validate:"required,oneof=Open Closed Published Created"`
	ServiceType     string `json:"serviceType" validate:"required,oneof=Construction Delivery Manufacture"`
	OrganizationId  string `json:"organizationId" validate:"required,uuid"`
	CreatorUserName string `json:"creatorUsername" validate:"required"`
}

type EditTenderRequest struct {
	Name        string `json:"name" validate:"omitempty,max=50"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"omitempty,oneof=Open Closed Published Created"`
}
