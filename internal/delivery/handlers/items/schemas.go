package items

type JSONAPIItem struct {
	Data JSONAPIItemData `json:"data"`
}

type JSONAPIItemData struct {
	Type          string                    `json:"type" example:"items"`
	ID            string                    `json:"id,omitempty" example:"item_1"`
	Attributes    JSONAPIItemAttributes     `json:"attributes"`
	Relationships *JSONAPIItemRelationships `json:"relationships,omitempty"`
}

type JSONAPIItemAttributes struct {
	Title       string `json:"title" example:"Item Title"`
	Description string `json:"description" example:"Item Description"`
}

type JSONAPIItemRelationships struct {
	ItemProperties *JSONAPIItemPropertiesRel `json:"item_properties,omitempty"`
}

type JSONAPIItemPropertiesRel struct {
	Data []JSONAPIItemPropertyDataIdentifier `json:"data"`
}

type JSONAPIItemPropertyDataIdentifier struct {
	Type string `json:"type" example:"item_properties"`
	ID   string `json:"id" example:"prop_1"`
}

type JSONAPIItemResponse struct {
	Data     JSONAPIItemData       `json:"data"`
	Included []JSONAPIItemProperty `json:"included,omitempty"`
}

type JSONAPIItemListResponse struct {
	Data []JSONAPIItemData `json:"data"`
}

type JSONAPIItemProperty struct {
	Data JSONAPIItemPropertyData `json:"data"`
}

type JSONAPIItemPropertyData struct {
	Type       string                        `json:"type" example:"item_properties"`
	ID         string                        `json:"id,omitempty" example:"prop_1"`
	Attributes JSONAPIItemPropertyAttributes `json:"attributes"`
}

type JSONAPIItemPropertyAttributes struct {
	ItemID string `json:"item_id" example:"item_1"`
	Name   string `json:"name" example:"Property Name"`
	Value  string `json:"value" example:"Property Value"`
}

type JSONAPIItemPropertyResponse struct {
	Data JSONAPIItemPropertyData `json:"data"`
}

type JSONAPIItemPropertyListResponse struct {
	Data []JSONAPIItemPropertyData `json:"data"`
}
