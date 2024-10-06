package model

type Tender struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	OrganizationID  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}
