package bolt

import (
	"github.com/ansible-semaphore/semaphore/db"
	"go.etcd.io/bbolt"
	"reflect"
)

/*
Integrations
*/
func (d *BoltDb) CreateIntegration(integration db.Integration) (db.Integration, error) {
	err := integration.Validate()

	if err != nil {
		return db.Integration{}, err
	}

	newIntegration, err := d.createObject(integration.ProjectID, db.IntegrationProps, integration)
	return newIntegration.(db.Integration), err
}

func (d *BoltDb) GetIntegrations(projectID int, params db.RetrieveQueryParams) (integrations []db.Integration, err error) {
	err = d.getObjects(projectID, db.IntegrationProps, params, nil, &integrations)
	return integrations, err
}

func (d *BoltDb) GetIntegration(projectID int, integrationID int) (integration db.Integration, err error) {
	err = d.getObject(projectID, db.IntegrationProps, intObjectID(integrationID), &integration)
	if err != nil {
		return
	}

	return
}

func (d *BoltDb) UpdateIntegration(integration db.Integration) error {
	err := integration.Validate()

	if err != nil {
		return err
	}

	return d.updateObject(integration.ProjectID, db.IntegrationProps, integration)

}

func (d *BoltDb) GetIntegrationRefs(projectID int, integrationID int) (db.IntegrationReferrers, error) {
	//return d.getObjectRefs(projectID, db.IntegrationProps, integrationID)
	return db.IntegrationReferrers{}, nil
}

func (d *BoltDb) DeleteIntegrationExtractValue(projectID int, valueID int, integrationID int) error {
	return d.deleteObject(projectID, db.IntegrationExtractValueProps, intObjectID(valueID), nil)
}

func (d *BoltDb) CreateIntegrationExtractValue(projectId int, value db.IntegrationExtractValue) (db.IntegrationExtractValue, error) {
	err := value.Validate()

	if err != nil {
		return db.IntegrationExtractValue{}, err
	}

	newValue, err := d.createObject(projectId, db.IntegrationExtractValueProps, value)
	return newValue.(db.IntegrationExtractValue), err

}

func (d *BoltDb) GetIntegrationExtractValues(projectID int, params db.RetrieveQueryParams, integrationID int) (values []db.IntegrationExtractValue, err error) {
	values = make([]db.IntegrationExtractValue, 0)
	var allValues []db.IntegrationExtractValue

	err = d.getObjects(projectID, db.IntegrationExtractValueProps, params, nil, &allValues)

	if err != nil {
		return
	}

	for _, v := range allValues {
		if v.IntegrationID == integrationID {
			values = append(values, v)
		}
	}

	return
}

func (d *BoltDb) GetIntegrationExtractValue(projectID int, valueID int, integrationID int) (value db.IntegrationExtractValue, err error) {
	err = d.getObject(projectID, db.IntegrationExtractValueProps, intObjectID(valueID), &value)
	return value, err
}

func (d *BoltDb) UpdateIntegrationExtractValue(projectID int, integrationExtractValue db.IntegrationExtractValue) error {
	err := integrationExtractValue.Validate()

	if err != nil {
		return err
	}

	return d.updateObject(projectID, db.IntegrationExtractValueProps, integrationExtractValue)
}

func (d *BoltDb) GetIntegrationExtractValueRefs(projectID int, valueID int, integrationID int) (db.IntegrationExtractorChildReferrers, error) {
	return d.getIntegrationExtractorChildrenRefs(projectID, db.IntegrationExtractValueProps, valueID)
}

/*
Integration Matcher
*/
func (d *BoltDb) CreateIntegrationMatcher(projectID int, matcher db.IntegrationMatcher) (db.IntegrationMatcher, error) {
	err := matcher.Validate()

	if err != nil {
		return db.IntegrationMatcher{}, err
	}
	newMatcher, err := d.createObject(projectID, db.IntegrationMatcherProps, matcher)
	return newMatcher.(db.IntegrationMatcher), err
}

func (d *BoltDb) GetIntegrationMatchers(projectID int, params db.RetrieveQueryParams, integrationID int) (matchers []db.IntegrationMatcher, err error) {
	matchers = make([]db.IntegrationMatcher, 0)
	var allMatchers []db.IntegrationMatcher

	err = d.getObjects(projectID, db.IntegrationMatcherProps, db.RetrieveQueryParams{}, nil, &allMatchers)

	if err != nil {
		return
	}

	for _, v := range allMatchers {
		if v.IntegrationID == integrationID {
			matchers = append(matchers, v)
		}
	}

	return
}

func (d *BoltDb) GetIntegrationMatcher(projectID int, matcherID int, integrationID int) (matcher db.IntegrationMatcher, err error) {
	var matchers []db.IntegrationMatcher
	matchers, err = d.GetIntegrationMatchers(projectID, db.RetrieveQueryParams{}, integrationID)

	for _, v := range matchers {
		if v.ID == matcherID {
			matcher = v
		}
	}

	return
}

func (d *BoltDb) UpdateIntegrationMatcher(projectID int, integrationMatcher db.IntegrationMatcher) error {
	err := integrationMatcher.Validate()

	if err != nil {
		return err
	}

	return d.updateObject(projectID, db.IntegrationMatcherProps, integrationMatcher)
}

func (d *BoltDb) DeleteIntegrationMatcher(projectID int, matcherID int, integrationID int) error {
	return d.deleteObject(projectID, db.IntegrationMatcherProps, intObjectID(matcherID), nil)
}
func (d *BoltDb) DeleteIntegration(projectID int, integrationID int) error {
	matchers, err := d.GetIntegrationMatchers(projectID, db.RetrieveQueryParams{}, integrationID)

	if err != nil {
		return err
	}

	for m := range matchers {
		d.DeleteIntegrationMatcher(projectID, matchers[m].ID, integrationID)
	}

	return d.deleteObject(projectID, db.IntegrationProps, intObjectID(integrationID), nil)
}

func (d *BoltDb) GetIntegrationMatcherRefs(projectID int, matcherID int, integrationID int) (db.IntegrationExtractorChildReferrers, error) {
	return d.getIntegrationExtractorChildrenRefs(projectID, db.IntegrationMatcherProps, matcherID)
}

var integrationAliasProps = db.ObjectProps{
	TableName:         "integration_alias",
	Type:              reflect.TypeOf(db.IntegrationAlias{}),
	PrimaryColumnName: "alias",
}

var projectLevelIntegrationId = -1

func (d *BoltDb) CreateIntegrationAlias(alias db.IntegrationAlias) (res db.IntegrationAlias, err error) {
	newAlias, err := d.createObject(alias.ProjectID, db.IntegrationAliasProps, alias)

	if err != nil {
		return
	}

	res = newAlias.(db.IntegrationAlias)

	_, err = d.createObject(-1, integrationAliasProps, alias)

	if err != nil {
		_ = d.DeleteIntegrationAlias(alias.ProjectID, alias.IntegrationID)
		return
	}

	return
}

func (d *BoltDb) GetIntegrationAlias(projectID int, integrationID *int) (res db.IntegrationAlias, err error) {
	if integrationID == nil {
		integrationID = &projectLevelIntegrationId
	}
	err = d.getObject(projectID, db.IntegrationAliasProps, intObjectID(*integrationID), &res)
	return
}

func (d *BoltDb) GetIntegrationAliasByAlias(alias string) (res db.IntegrationAlias, err error) {

	err = d.getObject(-1, integrationAliasProps, strObjectID(alias), &res)

	return
}

func (d *BoltDb) UpdateIntegrationAlias(alias db.IntegrationAlias) error {

	var integrationID int
	if alias.IntegrationID == nil {
		integrationID = projectLevelIntegrationId
	} else {
		integrationID = *alias.IntegrationID
	}

	oldAlias, err := d.GetIntegrationAlias(alias.ProjectID, &integrationID)
	if err != nil {
		return err
	}

	err = d.db.Update(func(tx *bbolt.Tx) error {
		err := d.updateObjectTx(tx, alias.ProjectID, db.IntegrationAliasProps, alias)
		if err != nil {
			return err
		}

		err = d.deleteObject(-1, integrationAliasProps, strObjectID(oldAlias.Alias), tx)
		if err != nil {
			return err
		}

		_, err = d.createObjectTx(tx, -1, integrationAliasProps, strObjectID(alias.Alias))

		return err
	})

	return err
}

func (d *BoltDb) DeleteIntegrationAlias(projectID int, integrationID *int) (err error) {
	if integrationID == nil {
		integrationID = &projectLevelIntegrationId
	}

	alias, err := d.GetIntegrationAlias(projectID, integrationID)
	if err != nil {
		return
	}

	err = d.deleteObject(projectID, db.IntegrationAliasProps, intObjectID(*integrationID), nil)
	if err != nil {
		return
	}

	err = d.deleteObject(-1, integrationAliasProps, strObjectID(alias.Alias), nil)
	if err != nil {
		return
	}

	return
}
