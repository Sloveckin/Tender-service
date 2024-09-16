package data

import "time"

type Tender struct {
	Id              string    `json:"tenderId"`
	Name            string    `json:"tenderName"`
	Version         int64     `json:"tenderVersion"`
	Description     string    `json:"tenderDescription"`
	Status          string    `json:"tenderStatus"`
	OrganisationId  string    `json:"organisationId"`
	CreatorUserName string    `json:"creatorUserName"`
	ServiceType     string    `json:"tenderServiceType"`
	CreatedAt       time.Time `json:"createdAt"`
}
