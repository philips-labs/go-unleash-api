package api

import (
	"bytes"
	"net/http"
)

type FeatureToggle struct {
	Archived     bool          `json:"archived"`
	CreatedAt    string        `json:"createdAt,omitempty"`
	LastSeenAt   string        `json:"lastSeenAt,omitempty"`
	Description  string        `json:"description"`
	Name         string        `json:"name"`
	Project      string        `json:"project"`
	Stale        bool          `json:"stale"`
	Type         string        `json:"type"`
	Environments []Environment `json:"environments"`
	Variants     []Variant     `json:"variants"`
}

type FeatureStrategy struct {
	ID          string               `json:"id,omitempty"`
	Name        string               `json:"name"`
	Constraints []StrategyConstraint `json:"constraints,omitempty"`
	Parameters  interface{}          `json:"parameters,omitempty"`
	SortOrder   int                  `json:"sortOrder"`
}

type StrategyConstraint struct {
	ContextName     string   `json:"contextName"`
	Operator        string   `json:"operator"`
	Value           string   `json:"value"`
	Values          []string `json:"values"`
	Inverted        bool     `json:"inverted"`
	CaseInsensitive bool     `json:"caseInsensitive"`
}

type Variant struct {
	Name       string            `json:"name"`
	Stickiness string            `json:"stickiness"`
	Weight     int               `json:"weight"`
	WeightType string            `json:"weightType"`
	Overrides  []VariantOverride `json:"overrides,omitempty"`
	Payload    *VariantPayload   `json:"payload,omitempty"`
}

type VariantOverride struct {
	ContextName string   `json:"contextName"`
	Values      []string `json:"values"`
}

type VariantPayload struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type Environment struct {
	Enabled    bool              `json:"enabled"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Strategies []FeatureStrategy `json:"strategies"`
}

type FeatureTogglesService struct {
	client *ApiClient
}

func (p *FeatureTogglesService) GetFeatureByName(projectId string, featureName string) (*FeatureToggle, *Response, error) {
	req, _ := p.client.newRequest("admin/projects/"+projectId+"/features/"+featureName, "GET", nil)

	var feature FeatureToggle

	resp, err := p.client.do(req, &feature)
	if err != nil {
		return nil, resp, err
	}
	return &feature, resp, err
}

func (p *FeatureTogglesService) CreateFeature(projectId string, feature FeatureToggle) (*FeatureToggle, *Response, error) {
	req, err := p.client.newRequest("admin/projects/"+projectId+"/features", "POST", feature)
	if err != nil {
		return nil, nil, err
	}

	var createdFeature FeatureToggle

	resp, err := p.client.do(req, &createdFeature)
	if err != nil {
		return nil, resp, err
	}
	return &createdFeature, resp, err
}

func (p *FeatureTogglesService) UpdateFeature(projectId string, feature FeatureToggle) (*FeatureToggle, *Response, error) {
	req, err := p.client.newRequest("admin/projects/"+projectId+"/features/"+feature.Name, "PUT", feature)
	if err != nil {
		return nil, nil, err
	}

	var updatedFeature FeatureToggle

	resp, err := p.client.do(req, &updatedFeature)
	if err != nil {
		return nil, resp, err
	}
	return &updatedFeature, resp, err
}

func (p *FeatureTogglesService) ArchiveFeature(projectId string, featureName string) (bool, *Response, error) {
	req, _ := p.client.newRequest("admin/projects/"+projectId+"/features/"+featureName, "DELETE", nil)

	var deleteResponse bytes.Buffer

	resp, err := p.client.do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusAccepted {
		return false, resp, err
	}
	return true, resp, nil
}

func (p *FeatureTogglesService) DeleteArchivedFeature(featureName string) (bool, *Response, error) {
	req, _ := p.client.newRequest("admin/archive/"+featureName, "DELETE", nil)

	var deleteResponse bytes.Buffer

	resp, err := p.client.do(req, &deleteResponse)
	if resp == nil || resp.StatusCode != http.StatusAccepted {
		return false, resp, err
	}
	return true, resp, nil
}

func (p *FeatureTogglesService) GetFeaturesByProject(projectId string) (*[]FeatureToggle, *Response, error) {
	req, _ := p.client.newRequest("admin/projects/"+projectId+"/features", "GET", nil)

	var features []FeatureToggle

	resp, err := p.client.do(req, &features)
	if err != nil {
		return nil, resp, err
	}
	return &features, resp, err
}

// Adds a strategy to a feature toggle in a given environment
func (p *FeatureTogglesService) AddStrategyToFeature(projectId string, featureName string, environment string, featureStrategy FeatureStrategy) (*FeatureStrategy, *Response, error) {
	req, err := p.client.newRequest("admin/projects/"+projectId+"/features/"+featureName+"/environments/"+environment+"/strategies", "POST", featureStrategy)
	if err != nil {
		return nil, nil, err
	}

	var addedStrategy FeatureStrategy

	resp, err := p.client.do(req, &addedStrategy)
	if err != nil {
		return nil, resp, err
	}
	return &addedStrategy, resp, err
}

func (p *FeatureTogglesService) UpdateFeatureStrategy(projectId string, featureName string, environment string, featureStrategy FeatureStrategy) (*FeatureStrategy, *Response, error) {
	req, err := p.client.newRequest("admin/projects/"+projectId+"/features/"+featureName+"/environments/"+environment+"/strategies/"+featureStrategy.ID, "PUT", featureStrategy)
	if err != nil {
		return nil, nil, err
	}

	var updatedStrategy FeatureStrategy

	resp, err := p.client.do(req, &updatedStrategy)
	if err != nil {
		return nil, resp, err
	}
	return &updatedStrategy, resp, err
}

// Deletes a strategy from a feature toggle in a given environment
func (p *FeatureTogglesService) DeleteStrategyFromFeature(projectId string, featureName string, environment string, strategyId string) (bool, *Response, error) {
	req, _ := p.client.newRequest("admin/projects/"+projectId+"/features/"+featureName+"/environments/"+environment+"/strategies/"+strategyId, "DELETE", nil)

	var deleteResponse bytes.Buffer

	resp, err := p.client.do(req, &deleteResponse)
	if resp == nil {
		return false, resp, err
	}
	return true, resp, nil
}

func (p *FeatureTogglesService) EnableFeatureOnEnvironment(projectId string, featureName string, environment string, enabled bool) (bool, *Response, error) {
	var featureState string
	if enabled {
		featureState = "on"
	} else {
		featureState = "off"
	}
	req, _ := p.client.newRequest("admin/projects/"+projectId+"/features/"+featureName+"/environments/"+environment+"/"+featureState, "POST", FeatureToggle{})

	var response bytes.Buffer

	resp, err := p.client.do(req, &response)
	if err != nil {
		return false, resp, err
	}
	return true, resp, nil
}
