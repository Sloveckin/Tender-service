package data

import "time"

type Bid struct {
	Id              string    `json:"bidId"`
	Name            string    `json:"bidName"`
	TenderId        string    `json:"tenderId"`
	Description     string    `json:"bidDescription"`
	Status          string    `json:"bidStatus"`
	OrganizationId  string    `json:"organizationId"`
	CreatorUserName string    `json:"bidAuthorId"`
	CreatedTime     time.Time `json:"createdAt"`
}
