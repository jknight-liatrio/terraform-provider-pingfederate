package pingfederate

import (
	"bytes"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pf "github.com/iwarapter/pingfederate-sdk-go/pingfederate/models"
)

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func hashcodeString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func setOfString() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

//
//func requiredListOfString() *schema.Schema {
//	return &schema.Schema{
//		Type:     schema.TypeList,
//		Required: true,
//		Elem: &schema.Schema{
//			Type: schema.TypeString,
//		},
//	}
//}

func resourceLinkSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"location": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func resourceRequiredLinkSchema() *schema.Schema {
	s := resourceLinkSchema()
	s.Required = true
	s.Optional = false
	return s
}

func resourcePluginConfiguration() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tables": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     resourceConfigTable(),
				},
				"fields": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem:     resourceConfigField(),
				},
				"sensitive_fields": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem:     resourceSensitiveConfigField(),
				},
			},
		},
	}
}

func resourceConfigTable() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rows": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceConfigRow(),
			},
			"inherited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceConfigRow() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fields": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceConfigField(),
			},
			"sensitive_fields": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceSensitiveConfigField(),
			},
		},
	}
}

func resourceSensitiveConfigField() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"inherited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceConfigField() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inherited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourcePasswordCredentialValidatorAttributeContract() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"inherited": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"core_attributes": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"extended_attributes": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSpAdapterAttributeContract() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"inherited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"core_attributes": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"extended_attributes": {
				Type:     schema.TypeSet,
				Optional: true,
				//MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func flattenSpAdapterAttribute(in []*pf.SpAdapterAttribute) *schema.Set {
	m := make([]interface{}, 0, len(in))
	for _, v := range in {
		m = append(m, *v.Name)
	}
	return schema.NewSet(schema.HashString, m)
}

func expandSpAdapterAttribute(in []interface{}) *[]*pf.SpAdapterAttribute {
	var contractList []*pf.SpAdapterAttribute
	for _, raw := range in {
		c := &pf.SpAdapterAttribute{
			Name: String(raw.(string)),
		}
		contractList = append(contractList, c)
	}
	return &contractList
}

func expandSpAdapterAttributeContract(in []interface{}) *pf.SpAdapterAttributeContract {
	pgc := &pf.SpAdapterAttributeContract{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["inherited"]; ok {
			pgc.Inherited = Bool(val.(bool))
		}
		if v, ok := l["extended_attributes"]; ok && len(v.(*schema.Set).List()) > 0 {
			pgc.ExtendedAttributes = expandSpAdapterAttribute(v.(*schema.Set).List())
		}
		if v, ok := l["core_attributes"]; ok && len(v.(*schema.Set).List()) > 0 {
			pgc.CoreAttributes = expandSpAdapterAttribute(v.(*schema.Set).List())
		}
	}
	return pgc
}

func flattenSpAdapterAttributeContract(in *pf.SpAdapterAttributeContract) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.Inherited != nil {
		s["inherited"] = *in.Inherited
	}
	if in.ExtendedAttributes != nil {
		s["extended_attributes"] = flattenSpAdapterAttribute(*in.ExtendedAttributes)
	}
	if in.CoreAttributes != nil {
		s["core_attributes"] = flattenSpAdapterAttribute(*in.CoreAttributes)
	}
	m = append(m, s)
	return m
}

func resourceSpAdapterTargetApplicationInfo() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"inherited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"application_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_icon_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func expandSpAdapterTargetApplicationInfo(in []interface{}) *pf.SpAdapterTargetApplicationInfo {
	pgc := &pf.SpAdapterTargetApplicationInfo{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if v, ok := l["inherited"]; ok {
			pgc.Inherited = Bool(v.(bool))
		}
		if v, ok := l["application_name"]; ok {
			pgc.ApplicationName = String(v.(string))
		}
		if v, ok := l["application_icon_url"]; ok {
			pgc.ApplicationIconUrl = String(v.(string))
		}
	}
	return pgc
}

func flattenSpAdapterTargetApplicationInfo(in *pf.SpAdapterTargetApplicationInfo) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.Inherited != nil {
		s["inherited"] = *in.Inherited
	}
	if in.ApplicationName != nil {
		s["application_name"] = *in.ApplicationName
	}
	if in.ApplicationIconUrl != nil {
		s["application_icon_url"] = *in.ApplicationIconUrl
	}
	m = append(m, s)
	return m
}

func resourceIdpAdapterAttributeContract() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"inherited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"core_attributes": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     resourceIdpAdapterAttribute(),
			},
			"mask_ognl_values": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"extended_attributes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceIdpAdapterAttribute(),
			},
		},
	}
}

func flattenIdpAdapterAttributeContract(in *pf.IdpAdapterAttributeContract) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	s["extended_attributes"] = flattenIdpAdapterAttributes(*in.ExtendedAttributes)
	if in.CoreAttributes != nil && len(*in.CoreAttributes) > 0 {
		s["core_attributes"] = flattenIdpAdapterAttributes(*in.CoreAttributes)
	}
	if in.MaskOgnlValues != nil {
		s["mask_ognl_values"] = *in.MaskOgnlValues
	}
	if in.Inherited != nil {
		s["inherited"] = *in.Inherited
	}
	m = append(m, s)
	return m
}

func expandIdpAdapterAttributeContract(in []interface{}) *pf.IdpAdapterAttributeContract {
	iac := &pf.IdpAdapterAttributeContract{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if v, ok := l["extended_attributes"]; ok && len(v.(*schema.Set).List()) > 0 {
			iac.ExtendedAttributes = expandIdpAdapterAttributes(v.(*schema.Set).List())
		}
		if v, ok := l["core_attributes"]; ok && len(v.(*schema.Set).List()) > 0 {
			iac.CoreAttributes = expandIdpAdapterAttributes(v.(*schema.Set).List())
		}
		if val, ok := l["mask_ognl_values"]; ok {
			iac.MaskOgnlValues = Bool(val.(bool))
		}
		if val, ok := l["inherited"]; ok {
			iac.Inherited = Bool(val.(bool))
		}
	}
	return iac
}

func resourceIdpAdapterAttribute() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of this attribute.",
				Required:    true,
			},
			"pseudonym": {
				Type:        schema.TypeBool,
				Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
				Optional:    true,
				Default:     false,
			},
			"masked": {
				Type:        schema.TypeBool,
				Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
				Optional:    true,
			},
		},
	}
}

func flattenIdpAdapterAttributes(in []*pf.IdpAdapterAttribute) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, len(in))
	for _, v := range in {
		s := make(map[string]interface{})
		if v.Name != nil {
			s["name"] = *v.Name
		}
		if v.Pseudonym != nil {
			s["pseudonym"] = *v.Pseudonym
		}
		if v.Masked != nil {
			s["masked"] = *v.Masked
		}
		m = append(m, s)
	}
	return m
}

func expandIdpAdapterAttributes(in []interface{}) *[]*pf.IdpAdapterAttribute {
	attributes := &[]*pf.IdpAdapterAttribute{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		c := &pf.IdpAdapterAttribute{}
		if val, ok := l["name"]; ok {
			c.Name = String(val.(string))
		}
		if val, ok := l["pseudonym"]; ok {
			c.Pseudonym = Bool(val.(bool))
		}
		if val, ok := l["masked"]; ok {
			c.Masked = Bool(val.(bool))
		}
		*attributes = append(*attributes, c)
	}
	return attributes
}

func resourceIdpAdapterAttributeMapping() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ldap_attribute_source": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceLdapAttributeSource(),
			},
			"jdbc_attribute_source": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceJdbcAttributeSource(),
			},
			"custom_attribute_source": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceCustomAttributeSource(),
			},
			"attribute_contract_fulfillment": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     resourceAttributeFulfillmentValue(),
			},
			"issuance_criteria": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceIssuanceCriteria(),
			},
			"inherited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func flattenIdpAdapterContractMapping(in *pf.IdpAdapterContractMapping) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.Inherited != nil {
		s["inherited"] = *in.Inherited
	}
	if in.AttributeContractFulfillment != nil {
		s["attribute_contract_fulfillment"] = flattenMapOfAttributeFulfillmentValue(in.AttributeContractFulfillment)
	}
	if in.IssuanceCriteria != nil && (in.IssuanceCriteria.ExpressionCriteria != nil && in.IssuanceCriteria.ConditionalCriteria != nil) {
		s["issuance_criteria"] = flattenIssuanceCriteria(in.IssuanceCriteria)
	}

	if *in.AttributeSources != nil && len(*in.AttributeSources) > 0 {
		var ldapAttributes []interface{}
		var jdbcAttributes []interface{}
		var customAttributes []interface{}
		for _, v := range *in.AttributeSources {
			switch *v.Type {
			case "LDAP":
				ldapAttributes = append(ldapAttributes, flattenLdapAttributeSource(&v.LdapAttributeSource))
			case "JDBC":
				jdbcAttributes = append(jdbcAttributes, flattenJdbcAttributeSource(v))
			case "CUSTOM":
				customAttributes = append(customAttributes, flattenCustomAttributeSource(&v.CustomAttributeSource))
			}
		}
		if len(ldapAttributes) > 0 {
			s["ldap_attribute_source"] = ldapAttributes
		}
		if len(jdbcAttributes) > 0 {
			s["jdbc_attribute_source"] = jdbcAttributes
		}
		if len(customAttributes) > 0 {
			s["custom_attribute_source"] = customAttributes
		}
	}
	m = append(m, s)
	return m
}

func expandIdpAdapterContractMapping(in []interface{}) *pf.IdpAdapterContractMapping {
	iac := &pf.IdpAdapterContractMapping{AttributeSources: &[]*pf.AttributeSource{}}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if v, ok := l["inherited"]; ok {
			iac.Inherited = Bool(v.(bool))
		}
		if v, ok := l["attribute_contract_fulfillment"]; ok {
			iac.AttributeContractFulfillment = expandMapOfAttributeFulfillmentValue(v.(*schema.Set).List())
		}
		if v, ok := l["issuance_criteria"]; ok {
			iac.IssuanceCriteria = expandIssuanceCriteria(v.([]interface{}))
		}

		if v, ok := l["ldap_attribute_source"]; ok && len(v.([]interface{})) > 0 {
			*iac.AttributeSources = append(*iac.AttributeSources, *expandLdapAttributeSource(v.([]interface{}))...)
		}
		if v, ok := l["jdbc_attribute_source"]; ok && len(v.([]interface{})) > 0 {
			*iac.AttributeSources = append(*iac.AttributeSources, *expandJdbcAttributeSource(v.([]interface{}))...)
		}
		if v, ok := l["custom_attribute_source"]; ok && len(v.([]interface{})) > 0 {
			*iac.AttributeSources = append(*iac.AttributeSources, *expandCustomAttributeSource(v.([]interface{}))...)
		}

	}
	return iac
}

func resourceLdapAttributeSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_store_ref": resourceRequiredLinkSchema(),
			"base_dn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search_scope": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search_filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attribute_contract_fulfillment": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceAttributeFulfillmentValue(),
			},
			"binary_attribute_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"member_of_nested_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func flattenLdapAttributeSource(in *pf.LdapAttributeSource) map[string]interface{} {
	s := make(map[string]interface{})
	if in.DataStoreRef != nil {
		s["data_store_ref"] = flattenResourceLink(in.DataStoreRef)
	}
	if in.BaseDn != nil {
		s["base_dn"] = *in.BaseDn
	}
	if in.Id != nil {
		s["id"] = *in.Id
	}
	if in.SearchScope != nil {
		s["search_scope"] = *in.SearchScope
	}
	if in.Description != nil {
		s["description"] = *in.Description
	}
	if in.SearchFilter != nil {
		s["search_filter"] = *in.SearchFilter
	}
	if in.AttributeContractFulfillment != nil {
		s["attribute_contract_fulfillment"] = flattenMapOfAttributeFulfillmentValue(in.AttributeContractFulfillment)
	}
	if in.BinaryAttributeSettings != nil {
		attributes := map[string]string{}
		for s2 := range in.BinaryAttributeSettings {
			attributes[s2] = *(*in.BinaryAttributeSettings[s2]).BinaryEncoding
		}
		s["binary_attribute_settings"] = attributes
	}
	if in.MemberOfNestedGroup != nil {
		s["member_of_nested_group"] = *in.MemberOfNestedGroup
	}
	return s
}

func expandLdapAttributeSource(in []interface{}) *[]*pf.AttributeSource {
	var sources []*pf.AttributeSource
	for _, raw := range in {
		l := raw.(map[string]interface{})
		src := &pf.AttributeSource{Type: String("LDAP")}
		iac := &pf.LdapAttributeSource{Type: String("LDAP")}
		if v, ok := l["data_store_ref"]; ok && len(v.([]interface{})) > 0 {
			iac.DataStoreRef = expandResourceLink(v.([]interface{}))
			src.DataStoreRef = iac.DataStoreRef
		}
		if v, ok := l["base_dn"]; ok {
			iac.BaseDn = String(v.(string))
		}
		if v, ok := l["id"]; ok {
			iac.Id = String(v.(string))
			src.Id = iac.Id
		}
		if v, ok := l["search_scope"]; ok {
			iac.SearchScope = String(v.(string))
		}
		if v, ok := l["description"]; ok {
			iac.Description = String(v.(string))
			src.Description = iac.Description
		}
		if v, ok := l["search_filter"]; ok {
			iac.SearchFilter = String(v.(string))
		}
		if v, ok := l["attribute_contract_fulfillment"]; ok {
			iac.AttributeContractFulfillment = expandMapOfAttributeFulfillmentValue(v.(*schema.Set).List())
			src.AttributeContractFulfillment = iac.AttributeContractFulfillment
		}
		if v, ok := l["binary_attribute_settings"]; ok {
			ca := map[string]*pf.BinaryLdapAttributeSettings{}
			for key, val := range v.(map[string]interface{}) {
				ca[key] = &pf.BinaryLdapAttributeSettings{BinaryEncoding: String(val.(string))}
			}
			iac.BinaryAttributeSettings = ca
		}
		if v, ok := l["member_of_nested_group"]; ok {
			iac.MemberOfNestedGroup = Bool(v.(bool))
		}
		src.LdapAttributeSource = *iac
		sources = append(sources, src)
	}
	return &sources
}

func resourceJdbcAttributeSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_store_ref": resourceRequiredLinkSchema(),
			"schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attribute_contract_fulfillment": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceAttributeFulfillmentValue(),
			},
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func flattenJdbcAttributeSource(in *pf.AttributeSource) map[string]interface{} {
	s := make(map[string]interface{})
	if in.DataStoreRef != nil {
		s["data_store_ref"] = flattenResourceLink(in.DataStoreRef)
	}
	if in.Schema != nil {
		s["schema"] = *in.Schema
	}
	if in.Id != nil {
		s["id"] = *in.Id
	}
	if in.Table != nil {
		s["table"] = *in.Table
	}
	if in.Description != nil {
		s["description"] = *in.Description
	}
	if in.AttributeContractFulfillment != nil {
		s["attribute_contract_fulfillment"] = flattenMapOfAttributeFulfillmentValue(in.AttributeContractFulfillment)
	}
	if in.Filter != nil {
		s["filter"] = *in.Filter
	}
	return s
}

func expandJdbcAttributeSource(in []interface{}) *[]*pf.AttributeSource {
	var sources []*pf.AttributeSource
	for _, raw := range in {
		l := raw.(map[string]interface{})
		src := &pf.AttributeSource{Type: String("JDBC")}
		iac := &pf.JdbcAttributeSource{Type: String("JDBC")}
		if v, ok := l["data_store_ref"]; ok && len(v.([]interface{})) > 0 {
			iac.DataStoreRef = expandResourceLink(v.([]interface{}))
			src.DataStoreRef = iac.DataStoreRef
		}
		if v, ok := l["schema"]; ok {
			iac.Schema = String(v.(string))
		}
		if v, ok := l["id"]; ok {
			iac.Id = String(v.(string))
			src.Id = iac.Id
		}
		if v, ok := l["table"]; ok {
			iac.Table = String(v.(string))
		}
		if v, ok := l["description"]; ok {
			iac.Description = String(v.(string))
			src.Description = iac.Description
		}
		if v, ok := l["filter"]; ok {
			iac.Filter = String(v.(string))
		}
		if v, ok := l["attribute_contract_fulfillment"]; ok {
			iac.AttributeContractFulfillment = expandMapOfAttributeFulfillmentValue(v.(*schema.Set).List())
			src.AttributeContractFulfillment = iac.AttributeContractFulfillment
		}
		src.JdbcAttributeSource = *iac
		sources = append(sources, src)
	}
	return &sources
}

func resourceCustomAttributeSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_store_ref": resourceRequiredLinkSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attribute_contract_fulfillment": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceAttributeFulfillmentValue(),
			},
			"filter_fields": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     resourceFieldEntry(),
			},
		},
	}
}

func flattenCustomAttributeSource(in *pf.CustomAttributeSource) map[string]interface{} {
	s := make(map[string]interface{})
	if in.DataStoreRef != nil {
		s["data_store_ref"] = flattenResourceLink(in.DataStoreRef)
	}
	if in.Id != nil {
		s["id"] = *in.Id
	}
	if in.Description != nil {
		s["description"] = *in.Description
	}
	if in.AttributeContractFulfillment != nil {
		s["attribute_contract_fulfillment"] = flattenMapOfAttributeFulfillmentValue(in.AttributeContractFulfillment)
	}
	if in.FilterFields != nil {
		s["filter_fields"] = flattenFieldEntry(in.FilterFields)
	}
	return s
}

func expandCustomAttributeSource(in []interface{}) *[]*pf.AttributeSource {
	var sources []*pf.AttributeSource
	for _, raw := range in {
		l := raw.(map[string]interface{})
		src := &pf.AttributeSource{Type: String("CUSTOM")}
		iac := &pf.CustomAttributeSource{Type: String("CUSTOM")}
		if v, ok := l["data_store_ref"]; ok && len(v.([]interface{})) > 0 {
			iac.DataStoreRef = expandResourceLink(v.([]interface{}))
			src.DataStoreRef = iac.DataStoreRef
		}
		if v, ok := l["id"]; ok {
			iac.Id = String(v.(string))
			src.Id = iac.Id
		}
		if v, ok := l["description"]; ok {
			iac.Description = String(v.(string))
			src.Description = iac.Description
		}
		if v, ok := l["filter_fields"]; ok {
			iac.FilterFields = expandFieldEntry(v.([]interface{}))
		}
		if v, ok := l["attribute_contract_fulfillment"]; ok {
			iac.AttributeContractFulfillment = expandMapOfAttributeFulfillmentValue(v.(*schema.Set).List())
			src.AttributeContractFulfillment = iac.AttributeContractFulfillment
		}
		src.CustomAttributeSource = *iac
		sources = append(sources, src)
	}
	return &sources
}

func resourceFieldEntry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func expandFieldEntry(in []interface{}) *[]*pf.FieldEntry {
	var fields []*pf.FieldEntry
	for _, raw := range in {
		l := raw.(map[string]interface{})
		f := &pf.FieldEntry{}
		if v, ok := l["name"]; ok {
			f.Name = String(v.(string))
		}
		if v, ok := l["value"]; ok {
			f.Value = String(v.(string))
		}
		fields = append(fields, f)
	}
	return &fields
}

func flattenFieldEntry(in *[]*pf.FieldEntry) []interface{} {
	var m []interface{}
	for _, v := range *in {
		s := make(map[string]interface{})
		if v.Name != nil {
			s["name"] = *v.Name
		}
		if v.Value != nil {
			s["value"] = *v.Value
		}
		m = append(m, s)
	}
	return m
}

func resourceAttributeFulfillmentValue() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": resourceSourceTypeIdKey(),
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func expandMapOfAttributeFulfillmentValue(in []interface{}) map[string]*pf.AttributeFulfillmentValue {
	ca := map[string]*pf.AttributeFulfillmentValue{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if v, ok := l["key_name"]; ok {
			ca[v.(string)] = expandAttributeFulfillmentValue(l)
		}
	}
	return ca
}

func expandAttributeFulfillmentValue(in map[string]interface{}) *pf.AttributeFulfillmentValue {
	ca := &pf.AttributeFulfillmentValue{}
	if v, ok := in["source"]; ok {
		ca.Source = expandSourceTypeIdKey(v.([]interface{}))
	}
	if v, ok := in["value"]; ok && v != "" {
		ca.Value = String(v.(string))
	}
	return ca
}

func flattenMapOfAttributeFulfillmentValue(in map[string]*pf.AttributeFulfillmentValue) *schema.Set {
	m := make([]interface{}, 0, len(in))
	for s2 := range in {
		s := flattenAttributeFulfillmentValue(in[s2])
		s["key_name"] = s2
		m = append(m, s)
	}
	return schema.NewSet(attributeFulfillmentValueHash, m)
}

func attributeFulfillmentValueHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(m["key_name"].(string))
	//if d, ok := m["value"]; ok && d.(string) != "" {
	//	buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	//}
	//if d, ok := m["source"]; ok && d.(string) != "" {
	//	buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	//}
	return hashcodeString(buf.String())
}

func flattenAttributeFulfillmentValue(in *pf.AttributeFulfillmentValue) map[string]interface{} {
	s := make(map[string]interface{})
	if in.Source != nil {
		s["source"] = flattenSourceTypeIdKey(in.Source)
	}
	if in.Value != nil {
		s["value"] = *in.Value
	}
	return s
}

func resourceSourceTypeIdKey() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					//TODO ValidateFunc:
					// ['TOKEN_EXCHANGE_PROCESSOR_POLICY' or 'ACCOUNT_LINK' or 'ADAPTER' or 'ASSERTION' or 'CONTEXT' or 'CUSTOM_DATA_STORE' or 'EXPRESSION' or 'JDBC_DATA_STORE' or 'LDAP_DATA_STORE' or 'MAPPED_ATTRIBUTES' or 'NO_MAPPING' or 'TEXT' or 'TOKEN' or 'REQUEST' or 'OAUTH_PERSISTENT_GRANT' or 'SUBJECT_TOKEN' or 'ACTOR_TOKEN' or 'PASSWORD_CREDENTIAL_VALIDATOR' or 'IDP_CONNECTION' or 'AUTHENTICATION_POLICY_CONTRACT' or 'CLAIMS' or 'LOCAL_IDENTITY_PROFILE' or 'EXTENDED_CLIENT_METADATA' or 'EXTENDED_PROPERTIES' or 'TRACKED_HTTP_PARAMS']
				},
				"id": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func expandSourceTypeIdKey(in []interface{}) *pf.SourceTypeIdKey {
	ca := &pf.SourceTypeIdKey{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["id"]; ok && val.(string) != "" { //TODO im not sure why it insists on saving the ID as empty
			ca.Id = String(val.(string))
		}
		if val, ok := l["type"]; ok {
			ca.Type = String(val.(string))
		}
	}
	return ca
}

func flattenSourceTypeIdKey(in *pf.SourceTypeIdKey) []interface{} {
	m := make([]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.Id != nil {
		s["id"] = *in.Id
	}
	if in.Type != nil {
		s["type"] = *in.Type
	}
	m = append(m, s)
	return m
}

func resourceIssuanceCriteria() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"conditional_criteria": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceConditionalIssuanceCriteriaEntry(),
			},
			"expression_criteria": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceExpressionIssuanceCriteriaEntry(),
			},
		},
	}
}

func flattenIssuanceCriteria(in *pf.IssuanceCriteria) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.ConditionalCriteria != nil && len(*in.ConditionalCriteria) > 0 {
		s["conditional_criteria"] = flattenConditionalIssuanceCriteriaEntry(*in.ConditionalCriteria)
	}
	if in.ExpressionCriteria != nil && len(*in.ExpressionCriteria) > 0 {
		s["expression_criteria"] = flattenExpressionIssuanceCriteriaEntry(*in.ExpressionCriteria)
	}
	m = append(m, s)
	return m
}

func expandIssuanceCriteria(in []interface{}) *pf.IssuanceCriteria {
	exp := &pf.IssuanceCriteria{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if v, ok := l["conditional_criteria"]; ok {
			exp.ConditionalCriteria = expandConditionalIssuanceCriteriaEntry(v.([]interface{}))
		}
		if v, ok := l["expression_criteria"]; ok {
			exp.ExpressionCriteria = expandExpressionIssuanceCriteriaEntry(v.([]interface{}))
		}
	}
	return exp
}

func resourceConditionalIssuanceCriteriaEntry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": resourceSourceTypeIdKey(),
			"attribute_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:     schema.TypeString,
				Required: true,
				//TODO ValidateFunc: //['EQUALS' or 'EQUALS_CASE_INSENSITIVE' or 'EQUALS_DN' or 'NOT_EQUAL' or 'NOT_EQUAL_CASE_INSENSITIVE' or 'NOT_EQUAL_DN' or 'MULTIVALUE_CONTAINS' or 'MULTIVALUE_CONTAINS_CASE_INSENSITIVE' or 'MULTIVALUE_CONTAINS_DN' or 'MULTIVALUE_DOES_NOT_CONTAIN' or 'MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE' or 'MULTIVALUE_DOES_NOT_CONTAIN_DN']
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"error_result": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func flattenConditionalIssuanceCriteriaEntry(in []*pf.ConditionalIssuanceCriteriaEntry) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, len(in))
	for _, v := range in {
		s := make(map[string]interface{})
		if v.Source != nil {
			s["source"] = flattenSourceTypeIdKey(v.Source)
		}
		if v.AttributeName != nil {
			s["attribute_name"] = *v.AttributeName
		}
		if v.Condition != nil {
			s["condition"] = *v.Condition
		}
		if v.Value != nil {
			s["value"] = *v.Value
		}
		if v.ErrorResult != nil {
			s["error_result"] = *v.ErrorResult
		}
		m = append(m, s)
	}
	return m
}

func expandConditionalIssuanceCriteriaEntry(in []interface{}) *[]*pf.ConditionalIssuanceCriteriaEntry {
	exps := &[]*pf.ConditionalIssuanceCriteriaEntry{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		exp := &pf.ConditionalIssuanceCriteriaEntry{}
		if v, ok := l["source"]; ok {
			exp.Source = expandSourceTypeIdKey(v.([]interface{}))
		}
		if v, ok := l["attribute_name"]; ok {
			exp.AttributeName = String(v.(string))
		}
		if v, ok := l["condition"]; ok {
			exp.Condition = String(v.(string))
		}
		if v, ok := l["value"]; ok {
			exp.Value = String(v.(string))
		}
		if v, ok := l["error_result"]; ok {
			exp.ErrorResult = String(v.(string))
		}
		*exps = append(*exps, exp)
	}
	return exps
}

func resourceExpressionIssuanceCriteriaEntry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"error_result": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func flattenExpressionIssuanceCriteriaEntry(in []*pf.ExpressionIssuanceCriteriaEntry) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, len(in))
	for _, v := range in {
		s := make(map[string]interface{})
		if v.Expression != nil {
			s["expression"] = *v.Expression
		}
		if v.ErrorResult != nil {
			s["error_result"] = *v.ErrorResult
		}
		m = append(m, s)
	}
	return m
}

func expandExpressionIssuanceCriteriaEntry(in []interface{}) *[]*pf.ExpressionIssuanceCriteriaEntry {
	exps := &[]*pf.ExpressionIssuanceCriteriaEntry{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		exp := &pf.ExpressionIssuanceCriteriaEntry{}
		if v, ok := l["expression"]; ok {
			exp.Expression = String(v.(string))
		}
		if v, ok := l["error_result"]; ok {
			exp.ErrorResult = String(v.(string))
		}
		*exps = append(*exps, exp)
	}
	return exps
}

// Takes list of pointers to strings. Expand to an array
// of raw strings and returns a []interface{}
// to keep compatibility w/ schema.NewSetschema.NewSet
func flattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

func flattenScopes(in []*pf.ScopeEntry) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, len(in))
	for _, v := range in {
		s := make(map[string]interface{})
		s["name"] = *v.Name
		s["description"] = *v.Description
		m = append(m, s)
	}
	return m
}

func expandScopes(in []interface{}) *[]*pf.ScopeEntry {
	var scopeList []*pf.ScopeEntry
	for _, raw := range in {
		l := raw.(map[string]interface{})
		s := &pf.ScopeEntry{
			Name:        String(l["name"].(string)),
			Description: String(l["description"].(string)),
		}
		scopeList = append(scopeList, s)
	}
	return &scopeList
}

func flattenScopeGroups(in []*pf.ScopeGroupEntry) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, len(in))
	for _, v := range in {
		s := make(map[string]interface{})
		s["name"] = *v.Name
		s["description"] = *v.Description
		s["scopes"] = flattenStringList(*v.Scopes)
		m = append(m, s)
	}
	return m
}

func expandScopeGroups(in []interface{}) *[]*pf.ScopeGroupEntry {
	var scopeGroupList []*pf.ScopeGroupEntry
	for _, raw := range in {
		l := raw.(map[string]interface{})
		s := &pf.ScopeGroupEntry{
			Name:        String(l["name"].(string)),
			Description: String(l["description"].(string)),
		}
		var scopes []*string
		for _, scope := range l["scopes"].([]interface{}) {
			scopes = append(scopes, String(scope.(string)))
		}
		s.Scopes = &scopes
		scopeGroupList = append(scopeGroupList, s)
	}
	return &scopeGroupList
}

func flattenPersistentGrantContract(in *pf.PersistentGrantContract) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	s["extended_attributes"] = flattenPersistentGrantAttributes(*in.ExtendedAttributes)
	m = append(m, s)
	return m
}

func expandPersistentGrantContract(in []interface{}) *pf.PersistentGrantContract {
	pgc := &pf.PersistentGrantContract{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		var atr []*pf.PersistentGrantAttribute
		for _, exAtr := range l["extended_attributes"].([]interface{}) {
			atr = append(atr, &pf.PersistentGrantAttribute{Name: String(exAtr.(string))})
		}
		pgc.ExtendedAttributes = &atr
	}
	return pgc
}

func flattenPersistentGrantAttributes(in []*pf.PersistentGrantAttribute) []interface{} {
	m := make([]interface{}, 0, len(in))
	for _, v := range in {
		m = append(m, *v.Name)
	}
	return m
}

func expandClientAuth(in []interface{}) *pf.ClientAuth {
	ca := &pf.ClientAuth{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["client_cert_issuer_dn"]; ok {
			ca.ClientCertIssuerDn = String(val.(string))
		}
		if val, ok := l["client_cert_subject_dn"]; ok {
			ca.ClientCertSubjectDn = String(val.(string))
		}
		if val, ok := l["enforce_replay_prevention"]; ok {
			ca.EnforceReplayPrevention = Bool(val.(bool))
		}
		if val, ok := l["secret"]; ok {
			ca.Secret = String(val.(string))
		}
		ca.Type = String(l["type"].(string))
	}
	return ca
}

func flattenClientAuth(orig, in *pf.ClientAuth) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.ClientCertIssuerDn != nil {
		s["client_cert_issuer_dn"] = *in.ClientCertIssuerDn
	}
	if in.ClientCertSubjectDn != nil {
		s["client_cert_subject_dn"] = *in.ClientCertSubjectDn
	}
	if in.EnforceReplayPrevention != nil {
		s["enforce_replay_prevention"] = *in.EnforceReplayPrevention
	}
	if in.Secret == nil && orig.Secret != nil {
		s["secret"] = *orig.Secret
	}
	s["type"] = *in.Type
	m = append(m, s)
	return m
}

func expandJwksSettings(in []interface{}) *pf.JwksSettings {
	ca := &pf.JwksSettings{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["jwks"]; ok {
			ca.Jwks = String(val.(string))
		}
		if val, ok := l["jwks_url"]; ok {
			ca.JwksUrl = String(val.(string))
		}
	}
	return ca
}

func flattenJwksSettings(in *pf.JwksSettings) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.Jwks != nil {
		s["jwks"] = *in.Jwks
	}
	if in.JwksUrl != nil {
		s["jwks_url"] = *in.JwksUrl
	}
	m = append(m, s)
	return m
}

func expandResourceLink(in []interface{}) *pf.ResourceLink {
	ca := &pf.ResourceLink{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["id"]; ok {
			ca.Id = String(val.(string))
		}
		// if val, ok := l["location"]; ok {
		// 	ca.Location = String(val.(string))
		// }
	}
	return ca
}

func flattenResourceLink(in *pf.ResourceLink) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.Id != nil {
		s["id"] = *in.Id
	}
	// if in.Location != nil {
	// 	s["location"] = *in.Location
	// }
	m = append(m, s)
	return m
}

func expandClientOIDCPolicy(in []interface{}) *pf.ClientOIDCPolicy {
	ca := &pf.ClientOIDCPolicy{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["grant_access_session_revocation_api"]; ok {
			ca.GrantAccessSessionRevocationApi = Bool(val.(bool))
		}
		if val, ok := l["id_token_signing_algorithm"]; ok {
			ca.IdTokenSigningAlgorithm = String(val.(string))
		}
		if val, ok := l["logout_uris"]; ok {
			str := expandStringList(val.([]interface{}))
			ca.LogoutUris = &str
		}
		if val, ok := l["ping_access_logout_capable"]; ok {
			ca.PingAccessLogoutCapable = Bool(val.(bool))
		}
		if val, ok := l["policy_group"]; ok {
			ca.PolicyGroup = expandResourceLink(val.([]interface{}))
		}
	}
	return ca
}

func flattenClientOIDCPolicy(in *pf.ClientOIDCPolicy) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.GrantAccessSessionRevocationApi != nil {
		s["grant_access_session_revocation_api"] = *in.GrantAccessSessionRevocationApi
	}
	if in.IdTokenSigningAlgorithm != nil {
		s["id_token_signing_algorithm"] = *in.IdTokenSigningAlgorithm
	}
	if in.LogoutUris != nil && len(*in.LogoutUris) > 0 {
		s["logout_uris"] = flattenStringList(*in.LogoutUris)
	}
	if in.PingAccessLogoutCapable != nil {
		s["ping_access_logout_capable"] = *in.PingAccessLogoutCapable
	}
	if in.PolicyGroup != nil {
		s["policy_group"] = flattenResourceLink(in.PolicyGroup)
	}
	m = append(m, s)
	return m
}

func flattenConfigField(in []*pf.ConfigField) *schema.Set {
	var m []interface{}
	for _, v := range in {
		if v.EncryptedValue != nil {
			continue
		}
		s := make(map[string]interface{})
		s["name"] = *v.Name
		//We check if the Encrypted value is set, if its not we can update the value as a normal password field
		//will not return the value so we need to not overwrite it, which unfortunely means we cannot track password changes
		//this is a limitation of ping federate.
		if v.Value != nil && v.EncryptedValue == nil {
			s["value"] = *v.Value
		}
		// if v.EncryptedValue != nil && *v.EncryptedValue != "" {
		// 	s["encrypted_value"] = *v.EncryptedValue
		// }
		if v.Inherited != nil {
			s["inherited"] = *v.Inherited
		}
		m = append(m, s)
	}
	return schema.NewSet(configFieldHash, m)
}

func flattenSensitiveConfigField(in []*pf.ConfigField) *schema.Set {
	var m []interface{}
	for _, v := range in {
		if v.EncryptedValue == nil {
			continue
		}
		s := make(map[string]interface{})
		s["name"] = *v.Name
		//We check if the Encrypted value is set, if its not we can update the value as a normal password field
		//will not return the value so we need to not overwrite it, which unfortunely means we cannot track password changes
		//this is a limitation of ping federate.
		//if v.Value != nil && v.EncryptedValue == nil {
		if v.Value != nil {
			s["value"] = *v.Value
		}
		// if v.EncryptedValue != nil && *v.EncryptedValue != "" {
		// 	s["encrypted_value"] = *v.EncryptedValue
		// }
		if v.Inherited != nil {
			s["inherited"] = *v.Inherited
		}
		m = append(m, s)
	}
	return schema.NewSet(configFieldHash, m)
}

func configFieldHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(m["name"].(string))
	if d, ok := m["value"]; ok && d.(string) != "" {
		buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	}
	// if d, ok := m["encrypted_value"]; ok && d.(string) != "" {
	// 	buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	// }
	if d, ok := m["inherited"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", d.(bool)))
	}
	return hashcodeString(buf.String())
}

func expandConfigFields(in []interface{}) *[]*pf.ConfigField {
	var configFields []*pf.ConfigField
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if _, ok := l["encrypted_value"]; !ok {
			s := &pf.ConfigField{
				Name: String(l["name"].(string)),
			}
			if val, ok := l["value"]; ok {
				s.Value = String(val.(string))
			}
			if val, ok := l["inherited"]; ok {
				s.Inherited = Bool(val.(bool))
			}
			configFields = append(configFields, s)
		}
	}
	return &configFields
}

func expandSensitiveConfigFields(in []interface{}) *[]*pf.ConfigField {
	var configFields []*pf.ConfigField
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["value"]; ok && val.(string) != "" {
			s := &pf.ConfigField{
				Name: String(l["name"].(string)),
			}
			if val, ok := l["value"]; ok {
				s.Value = String(val.(string))
			}
			if val, ok := l["inherited"]; ok {
				s.Inherited = Bool(val.(bool))
			}
			configFields = append(configFields, s)
		}
	}
	return &configFields
}

func flattenConfigRow(in []*pf.ConfigRow) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, len(in))
	for _, v := range in {
		s := make(map[string]interface{})
		s["fields"] = flattenConfigField(*v.Fields)
		s["sensitive_fields"] = flattenSensitiveConfigField(*v.Fields)
		m = append(m, s)
	}
	return m
}

func expandConfigRow(in []interface{}) *[]*pf.ConfigRow {
	configRows := []*pf.ConfigRow{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		row := &pf.ConfigRow{}
		if val, ok := l["fields"]; ok {
			row.Fields = expandConfigFields(val.(*schema.Set).List())
		}
		if val, ok := l["sensitive_fields"]; ok {
			fields := expandSensitiveConfigFields(val.(*schema.Set).List())
			*row.Fields = append(*row.Fields, *fields...)
		}
		configRows = append(configRows, row)
	}
	return &configRows
}

func flattenConfigTable(in []*pf.ConfigTable) []interface{} {
	var m []interface{}
	for _, v := range in {
		s := make(map[string]interface{})
		s["name"] = *v.Name
		if v.Rows != nil {
			s["rows"] = flattenConfigRow(*v.Rows)
		}
		if v.Inherited != nil {
			s["inherited"] = *v.Inherited
		}
		m = append(m, s)
	}
	return m
}

//func configTableHash(v interface{}) int {
//	var buf bytes.Buffer
//	m := v.(map[string]interface{})
//	buf.WriteString(m["name"].(string))
//	return hashcodeString(buf.String())
//}

func expandConfigTable(in []interface{}) *[]*pf.ConfigTable {
	var configTables []*pf.ConfigTable
	for _, raw := range in {
		l := raw.(map[string]interface{})
		s := &pf.ConfigTable{
			Name: String(l["name"].(string)),
		}
		if val, ok := l["rows"]; ok {
			s.Rows = expandConfigRow(val.([]interface{}))
		}
		if val, ok := l["inherited"]; ok {
			s.Inherited = Bool(val.(bool))
		}
		configTables = append(configTables, s)
	}
	return &configTables
}

func flattenPluginConfiguration(in *pf.PluginConfiguration) []interface{} {
	s := make(map[string]interface{})
	if in.Tables != nil {
		s["tables"] = flattenConfigTable(*in.Tables)
	}
	if in.Fields != nil {
		s["fields"] = flattenConfigField(*in.Fields)
	}
	if in.Fields != nil {
		s["sensitive_fields"] = flattenSensitiveConfigField(*in.Fields)
	}
	return []interface{}{s}
}

func expandPluginConfiguration(in []interface{}) *pf.PluginConfiguration {
	config := &pf.PluginConfiguration{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["tables"]; ok && len(val.([]interface{})) > 0 {
			config.Tables = expandConfigTable(val.([]interface{}))
		}
		if val, ok := l["fields"]; ok {
			config.Fields = expandConfigFields(val.(*schema.Set).List())
		}
		if val, ok := l["sensitive_fields"]; ok {
			fields := expandSensitiveConfigFields(val.(*schema.Set).List())
			*config.Fields = append(*config.Fields, *fields...)
		}
	}
	return config
}

func flattenAccessTokenAttributeContract(in *pf.AccessTokenAttributeContract) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.ExtendedAttributes != nil {
		s["extended_attributes"] = flattenAccessTokenAttributes(*in.ExtendedAttributes)
	}
	if in.CoreAttributes != nil {
		s["core_attributes"] = flattenAccessTokenAttributes(*in.CoreAttributes)
	}
	m = append(m, s)
	return m
}

func expandAccessTokenAttributeContract(in []interface{}) *pf.AccessTokenAttributeContract {
	pgc := &pf.AccessTokenAttributeContract{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		var atr []*pf.AccessTokenAttribute
		for _, exAtr := range l["extended_attributes"].([]interface{}) {
			atr = append(atr, &pf.AccessTokenAttribute{Name: String(exAtr.(string))})
		}
		pgc.ExtendedAttributes = &atr
	}
	return pgc
}

func flattenAccessTokenAttributes(in []*pf.AccessTokenAttribute) []interface{} {
	m := make([]interface{}, 0, len(in))
	for _, v := range in {
		m = append(m, *v.Name)
	}
	return m
}

func flattenAuthenticationPolicyContractAttribute(in []*pf.AuthenticationPolicyContractAttribute) []interface{} {
	m := make([]interface{}, 0, len(in))
	for _, v := range in {
		m = append(m, *v.Name)
	}
	return m
}

func expandAuthenticationPolicyContractAttribute(in []interface{}) *[]*pf.AuthenticationPolicyContractAttribute {
	var contractList []*pf.AuthenticationPolicyContractAttribute
	for _, raw := range in {
		c := &pf.AuthenticationPolicyContractAttribute{
			Name: String(raw.(string)),
		}
		contractList = append(contractList, c)
	}
	return &contractList
}

func flattenPasswordCredentialValidatorAttribute(in []*pf.PasswordCredentialValidatorAttribute) *schema.Set {
	m := make([]interface{}, 0, len(in))
	for _, v := range in {
		m = append(m, *v.Name)
	}
	return schema.NewSet(schema.HashString, m)
}

func expandPasswordCredentialValidatorAttribute(in []interface{}) *[]*pf.PasswordCredentialValidatorAttribute {
	contractList := []*pf.PasswordCredentialValidatorAttribute{}
	for _, raw := range in {
		c := &pf.PasswordCredentialValidatorAttribute{
			Name: String(raw.(string)),
		}
		contractList = append(contractList, c)
	}
	return &contractList
}

func expandPasswordCredentialValidatorAttributeContract(in []interface{}) *pf.PasswordCredentialValidatorAttributeContract {
	pgc := &pf.PasswordCredentialValidatorAttributeContract{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["inherited"]; ok {
			pgc.Inherited = Bool(val.(bool))
		}
		if v, ok := l["extended_attributes"]; ok && len(v.(*schema.Set).List()) > 0 {
			pgc.ExtendedAttributes = expandPasswordCredentialValidatorAttribute(v.(*schema.Set).List())
		}
		if v, ok := l["core_attributes"]; ok && len(v.(*schema.Set).List()) > 0 {
			pgc.CoreAttributes = expandPasswordCredentialValidatorAttribute(v.(*schema.Set).List())
		}
	}
	return pgc
}

func flattenPasswordCredentialValidatorAttributeContract(in *pf.PasswordCredentialValidatorAttributeContract) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.Inherited != nil {
		s["inherited"] = *in.Inherited
	}
	if in.ExtendedAttributes != nil {
		s["extended_attributes"] = flattenPasswordCredentialValidatorAttribute(*in.ExtendedAttributes)
	}
	if in.CoreAttributes != nil {
		s["core_attributes"] = flattenPasswordCredentialValidatorAttribute(*in.CoreAttributes)
	}
	m = append(m, s)
	return m
}

func flattenJdbcDataStore(in *pf.JdbcDataStore) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.MaskAttributeValues != nil {
		s["mask_attribute_values"] = *in.MaskAttributeValues
	}
	if in.ConnectionUrlTags != nil && len(*in.ConnectionUrlTags) != 0 {
		s["connection_url_tags"] = flattenJdbcTagConfigs(in.ConnectionUrlTags)
	}
	if in.ConnectionUrl != nil {
		s["connection_url"] = *in.ConnectionUrl
	}
	if in.Name != nil {
		s["name"] = *in.Name
	}
	if in.DriverClass != nil {
		s["driver_class"] = *in.DriverClass
	}
	if in.UserName != nil {
		s["user_name"] = *in.UserName
	}
	if in.Password != nil {
		//TODO i need to handle this not being set
		s["password"] = *in.Password
	}
	if in.EncryptedPassword != nil {
		s["encrypted_password"] = *in.EncryptedPassword
	}
	if in.ValidateConnectionSql != nil {
		s["validate_connection_sql"] = *in.ValidateConnectionSql
	}
	if in.AllowMultiValueAttributes != nil {
		s["allow_multi_value_attributes"] = *in.AllowMultiValueAttributes
	}
	if in.MinPoolSize != nil {
		s["min_pool_size"] = *in.MinPoolSize
	}
	if in.MaxPoolSize != nil {
		s["max_pool_size"] = *in.MaxPoolSize
	}
	if in.BlockingTimeout != nil {
		s["blocking_timeout"] = *in.BlockingTimeout
	}
	if in.IdleTimeout != nil {
		s["idle_timeout"] = *in.IdleTimeout
	}
	//s["type"] = String("JDBC")
	m = append(m, s)
	return m
}

func expandJdbcDataStore(in []interface{}) *pf.JdbcDataStore {
	ca := &pf.JdbcDataStore{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["mask_attribute_values"]; ok {
			ca.MaskAttributeValues = Bool(val.(bool))
		}
		if val, ok := l["connection_url_tags"]; ok && len(val.(*schema.Set).List()) > 0 {
			ca.ConnectionUrlTags = expandJdbcTagConfigs(val.(*schema.Set).List())
		}
		if val, ok := l["connection_url"]; ok {
			ca.ConnectionUrl = String(val.(string))
		}
		if val, ok := l["name"]; ok {
			ca.Name = String(val.(string))
		}
		if val, ok := l["driver_class"]; ok {
			ca.DriverClass = String(val.(string))
		}
		if val, ok := l["user_name"]; ok {
			ca.UserName = String(val.(string))
		}
		if val, ok := l["password"]; ok {
			ca.Password = String(val.(string))
		}
		if val, ok := l["encrypted_password"]; ok {
			ca.EncryptedPassword = String(val.(string))
		}
		if val, ok := l["validate_connection_sql"]; ok {
			ca.ValidateConnectionSql = String(val.(string))
		}
		if val, ok := l["allow_multi_value_attributes"]; ok {
			ca.AllowMultiValueAttributes = Bool(val.(bool))
		}
		if val, ok := l["min_pool_size"]; ok {
			ca.MinPoolSize = Int(val.(int))
		}
		if val, ok := l["max_pool_size"]; ok {
			ca.MaxPoolSize = Int(val.(int))
		}
		if val, ok := l["blocking_timeout"]; ok {
			ca.BlockingTimeout = Int(val.(int))
		}
		if val, ok := l["idle_timeout"]; ok {
			ca.IdleTimeout = Int(val.(int))
		}
		ca.Type = String("JDBC")
	}
	return ca
}

func flattenJdbcTagConfigs(in *[]*pf.JdbcTagConfig) *schema.Set {
	var m []interface{}
	for _, v := range *in {
		s := make(map[string]interface{})
		if v.ConnectionUrl != nil {
			s["connection_url"] = *v.ConnectionUrl
		}
		if v.Tags != nil {
			s["tags"] = *v.Tags
		}
		if v.DefaultSource != nil {
			s["default_source"] = *v.DefaultSource
		}
		m = append(m, s)
	}
	return schema.NewSet(jdbcTagConfigHash, m)
}

func jdbcTagConfigHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(m["connection_url"].(string))
	if d, ok := m["tags"]; ok && d.(string) != "" {
		buf.WriteString(fmt.Sprintf("%s-", d.(string)))
	}
	if d, ok := m["default_source"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", d.(bool)))
	}
	return hashcodeString(buf.String())
}

func expandJdbcTagConfigs(in []interface{}) *[]*pf.JdbcTagConfig {
	var tags []*pf.JdbcTagConfig
	for _, raw := range in {
		l := raw.(map[string]interface{})
		f := &pf.JdbcTagConfig{}
		if v, ok := l["connection_url"]; ok {
			f.ConnectionUrl = String(v.(string))
		}
		if v, ok := l["tags"]; ok {
			f.Tags = String(v.(string))
		}
		if v, ok := l["default_source"]; ok {
			f.DefaultSource = Bool(v.(bool))
		}
		tags = append(tags, f)
	}
	return &tags
}

func flattenLdapDataStore(in *pf.LdapDataStore) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.MaskAttributeValues != nil {
		s["mask_attribute_values"] = *in.MaskAttributeValues
	}
	//if in.HostnamesTags != nil && len(*in.HostnamesTags) != 0 {
	//	//TODO connection_url_tags
	//}
	if in.Hostnames != nil {
		s["hostnames"] = flattenStringList(*in.Hostnames)
	}
	if in.Name != nil {
		s["name"] = *in.Name
	}
	if in.LdapType != nil {
		s["ldap_type"] = *in.LdapType
	}
	if in.BindAnonymously != nil {
		s["bind_anonymously"] = *in.BindAnonymously
	}
	if in.UserDN != nil {
		s["user_dn"] = *in.UserDN
	}
	if in.Password != nil {
		s["password"] = *in.Password
	}
	if in.EncryptedPassword != nil {
		s["encrypted_password"] = *in.EncryptedPassword
	}
	if in.UseSsl != nil {
		s["use_ssl"] = *in.UseSsl
	}
	if in.UseDnsSrvRecords != nil {
		s["use_dns_srv_records"] = *in.UseDnsSrvRecords
	}
	if in.FollowLDAPReferrals != nil {
		s["follow_ldap_referrals"] = *in.FollowLDAPReferrals
	}
	if in.TestOnBorrow != nil {
		s["test_on_borrow"] = *in.TestOnBorrow
	}
	if in.TestOnReturn != nil {
		s["test_on_return"] = *in.TestOnReturn
	}
	if in.CreateIfNecessary != nil {
		s["create_if_necessary"] = *in.CreateIfNecessary
	}
	if in.VerifyHost != nil {
		s["verify_host"] = *in.VerifyHost
	}
	if in.MinConnections != nil {
		s["min_connections"] = *in.MinConnections
	}
	if in.MaxConnections != nil {
		s["max_connections"] = *in.MaxConnections
	}
	if in.MaxWait != nil {
		s["max_wait"] = *in.MaxWait
	}
	if in.TimeBetweenEvictions != nil {
		s["time_between_evictions"] = *in.TimeBetweenEvictions
	}
	if in.ReadTimeout != nil {
		s["read_timeout"] = *in.ReadTimeout
	}
	if in.ConnectionTimeout != nil {
		s["connection_timeout"] = *in.ConnectionTimeout
	}
	if in.DnsTtl != nil {
		s["dns_ttl"] = *in.DnsTtl
	}
	if in.LdapDnsSrvPrefix != nil {
		s["ldap_dns_srv_prefix"] = *in.LdapDnsSrvPrefix
	}
	if in.LdapsDnsSrvPrefix != nil {
		s["ldaps_dns_srv_prefix"] = *in.LdapsDnsSrvPrefix
	}
	if in.BinaryAttributes != nil {
		s["binary_attributes"] = flattenStringList(*in.BinaryAttributes)
	}
	//s["type"] = String("LDAP")
	m = append(m, s)
	return m
}

func expandLdapDataStore(in []interface{}) *pf.LdapDataStore {
	ca := &pf.LdapDataStore{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		//log.Printf("[DEBUG] HELPER expandLdapDataStore: %v", l)
		if val, ok := l["mask_attribute_values"]; ok {
			ca.MaskAttributeValues = Bool(val.(bool))
		}
		//TODO hostnames_tags
		if val, ok := l["hostnames"]; ok {
			strs := expandStringList(val.(*schema.Set).List())
			ca.Hostnames = &strs
		}
		if val, ok := l["name"]; ok {
			ca.Name = String(val.(string))
		}
		if val, ok := l["ldap_type"]; ok {
			ca.LdapType = String(val.(string))
		}
		if val, ok := l["bind_anonymously"]; ok {
			ca.BindAnonymously = Bool(val.(bool))
		}
		if val, ok := l["user_dn"]; ok {
			ca.UserDN = String(val.(string))
		}
		if val, ok := l["password"]; ok {
			ca.Password = String(val.(string))
		}
		if val, ok := l["encrypted_password"]; ok {
			ca.EncryptedPassword = String(val.(string))
		}
		if val, ok := l["use_ssl"]; ok {
			ca.UseSsl = Bool(val.(bool))
		}
		if val, ok := l["use_dns_srv_records"]; ok {
			ca.UseDnsSrvRecords = Bool(val.(bool))
		}
		if val, ok := l["follow_ldap_referrals"]; ok {
			ca.FollowLDAPReferrals = Bool(val.(bool))
		}
		if val, ok := l["test_on_borrow"]; ok {
			ca.TestOnBorrow = Bool(val.(bool))
		}
		if val, ok := l["test_on_return"]; ok {
			ca.TestOnReturn = Bool(val.(bool))
		}
		if val, ok := l["create_if_necessary"]; ok {
			ca.CreateIfNecessary = Bool(val.(bool))
		}
		if val, ok := l["verify_host"]; ok {
			ca.VerifyHost = Bool(val.(bool))
		}
		if val, ok := l["min_connections"]; ok {
			ca.MinConnections = Int(val.(int))
		}
		if val, ok := l["max_connections"]; ok {
			ca.MaxConnections = Int(val.(int))
		}
		if val, ok := l["max_wait"]; ok {
			ca.MaxWait = Int(val.(int))
		}
		if val, ok := l["time_between_evictions"]; ok {
			ca.TimeBetweenEvictions = Int(val.(int))
		}
		if val, ok := l["read_timeout"]; ok {
			ca.ReadTimeout = Int(val.(int))
		}
		if val, ok := l["connection_timeout"]; ok {
			ca.ConnectionTimeout = Int(val.(int))
		}
		if val, ok := l["dns_ttl"]; ok {
			ca.DnsTtl = Int(val.(int))
		}
		if val, ok := l["ldap_dns_srv_prefix"]; ok {
			ca.LdapDnsSrvPrefix = String(val.(string))
		}
		if val, ok := l["ldaps_dns_srv_prefix"]; ok {
			ca.LdapsDnsSrvPrefix = String(val.(string))
		}
		if val, ok := l["binary_attributes"]; ok {
			strs := expandStringList(val.(*schema.Set).List())
			ca.BinaryAttributes = &strs
		}
		ca.Type = String("LDAP")
	}
	return ca
}

func maskPluginConfigurationFromDescriptor(desc *pf.PluginConfigDescriptor, origConf, conf *pf.PluginConfiguration) []interface{} {
	//printPluginConfig("originConf",origConf)
	//printPluginConfig("conf",conf)

	//if origConf.Fields != nil {
	for _, f := range *desc.Fields {
		if *f.Type == "HASHED_TEXT" || ((*f).Encrypted != nil && *f.Encrypted) {
			for _, i := range *conf.Fields {
				if *i.Name == *f.Name {
					s, _ := getConfigFieldValueByName(*i.Name, origConf.Fields)
					i.Value = String(s)
				}
			}
		}
	}
	//}

	//if origConf.Tables != nil {
	for _, dt := range *desc.Tables {
		for _, dc := range *dt.Columns {
			if *dc.Type == "HASHED_TEXT" || ((*dc).Encrypted != nil && *dc.Encrypted) {
				for ctIndex, ct := range *conf.Tables {
					for crIndex, cr := range *ct.Rows {
						for _, f := range *cr.Fields {
							if *f.Name == *dc.Name {
								val, _ := getConfigFieldValueByName(*f.Name, (*(*origConf.Tables)[ctIndex].Rows)[crIndex].Fields)
								f.Value = &val
							}
						}
					}
				}
			}
		}
		//}
	}
	return flattenPluginConfiguration(conf)
}

func getConfigFieldValueByName(name string, fields *[]*pf.ConfigField) (string, error) {
	for _, f := range *fields {
		if *f.Name == name {
			return *f.Value, nil
		}
	}
	return "", nil
}

func resourceAuthenticationSelectorAttributeContract() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"extended_attributes": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func flattenAuthenticationSelectorAttributeContract(in *pf.AuthenticationSelectorAttributeContract) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	s["extended_attributes"] = flattenAuthenticationSelectorAttributes(*in.ExtendedAttributes)
	m = append(m, s)
	return m
}

func expandAuthenticationSelectorAttributeContract(in []interface{}) *pf.AuthenticationSelectorAttributeContract {
	pgc := &pf.AuthenticationSelectorAttributeContract{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		var atr []*pf.AuthenticationSelectorAttribute
		for _, exAtr := range l["extended_attributes"].([]interface{}) {
			atr = append(atr, &pf.AuthenticationSelectorAttribute{Name: String(exAtr.(string))})
		}
		pgc.ExtendedAttributes = &atr
	}
	return pgc
}

func flattenAuthenticationSelectorAttributes(in []*pf.AuthenticationSelectorAttribute) []interface{} {
	m := make([]interface{}, 0, len(in))
	for _, v := range in {
		m = append(m, *v.Name)
	}
	return m
}

func expandAccessTokenMappingContext(in []interface{}) *pf.AccessTokenMappingContext {
	pgc := &pf.AccessTokenMappingContext{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if val, ok := l["type"]; ok {
			pgc.Type = String(val.(string))
		}
		if val, ok := l["context_ref"]; ok {
			pgc.ContextRef = expandResourceLink(val.([]interface{}))
		}
	}
	return pgc
}

func flattenAccessTokenMappingContext(in *pf.AccessTokenMappingContext) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	s["type"] = in.Type
	if in.ContextRef != nil {
		s["context_ref"] = flattenResourceLink(in.ContextRef)
	}
	m = append(m, s)
	return m
}

func flattenAttributeSources(d *schema.ResourceData, rv *[]*pf.AttributeSource) error {
	if *rv != nil && len(*rv) > 0 {
		var ldapAttributes []interface{}
		var jdbcAttributes []interface{}
		var customAttributes []interface{}
		for _, v := range *rv {
			switch *v.Type {
			case "LDAP":
				ldapAttributes = append(ldapAttributes, flattenLdapAttributeSource(&v.LdapAttributeSource))
			case "JDBC":
				jdbcAttributes = append(jdbcAttributes, flattenJdbcAttributeSource(v))
			case "CUSTOM":
				customAttributes = append(customAttributes, flattenCustomAttributeSource(&v.CustomAttributeSource))
			}
		}
		if len(ldapAttributes) > 0 {
			if err := d.Set("ldap_attribute_source", ldapAttributes); err != nil {
				return err
			}
		}
		if len(jdbcAttributes) > 0 {
			if err := d.Set("jdbc_attribute_source", jdbcAttributes); err != nil {
				return err
			}
		}
		if len(customAttributes) > 0 {
			if err := d.Set("custom_attribute_source", customAttributes); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceOpenIdConnectAttribute() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of this attribute.",
				Required:    true,
			},
			"include_in_id_token": {
				Type:        schema.TypeBool,
				Description: "Attribute is included in the ID Token.",
				Optional:    true,
				Default:     false,
			},
			"include_in_user_info": {
				Type:        schema.TypeBool,
				Description: "Attribute is included in the User Info.",
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func expandOpenIdConnectAttributes(in []interface{}) *[]*pf.OpenIdConnectAttribute {
	attributes := &[]*pf.OpenIdConnectAttribute{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		c := &pf.OpenIdConnectAttribute{}
		if val, ok := l["name"]; ok {
			c.Name = String(val.(string))
		}
		if val, ok := l["include_in_id_token"]; ok {
			c.IncludeInIdToken = Bool(val.(bool))
		}
		if val, ok := l["include_in_user_info"]; ok {
			c.IncludeInUserInfo = Bool(val.(bool))
		}
		*attributes = append(*attributes, c)
	}
	return attributes
}

func flattenOpenIdConnectAttributes(in []*pf.OpenIdConnectAttribute) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, len(in))
	for _, v := range in {
		s := make(map[string]interface{})
		s["name"] = v.Name
		if v.IncludeInUserInfo != nil {
			s["include_in_user_info"] = v.IncludeInUserInfo
		}
		if v.IncludeInIdToken != nil {
			s["include_in_id_token"] = v.IncludeInIdToken
		}
		m = append(m, s)
	}
	return m
}

func resourceOpenIdConnectAttributeContract() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"core_attributes": {
					Type: schema.TypeSet,
					//Optional: true,
					Computed: true,
					Elem:     resourceOpenIdConnectAttribute(),
				},
				"extended_attributes": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem:     resourceOpenIdConnectAttribute(),
				},
			},
		},
	}
}

func flattenOpenIdConnectAttributeContract(in *pf.OpenIdConnectAttributeContract) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.ExtendedAttributes != nil && len(*in.ExtendedAttributes) > 0 {
		s["extended_attributes"] = flattenOpenIdConnectAttributes(*in.ExtendedAttributes)
	}
	if in.CoreAttributes != nil && len(*in.CoreAttributes) > 0 {
		s["core_attributes"] = flattenOpenIdConnectAttributes(*in.CoreAttributes)
	}
	m = append(m, s)
	return m
}

func expandOpenIdConnectAttributeContract(in []interface{}) *pf.OpenIdConnectAttributeContract {
	iac := &pf.OpenIdConnectAttributeContract{}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if v, ok := l["extended_attributes"]; ok && len(v.(*schema.Set).List()) > 0 {
			iac.ExtendedAttributes = expandOpenIdConnectAttributes(v.(*schema.Set).List())
		}
		if v, ok := l["core_attributes"]; ok && len(v.(*schema.Set).List()) > 0 {
			iac.CoreAttributes = expandOpenIdConnectAttributes(v.(*schema.Set).List())
		}
	}
	return iac
}

func resourceAttributeMapping() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ldap_attribute_source": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceLdapAttributeSource(),
			},
			"jdbc_attribute_source": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceJdbcAttributeSource(),
			},
			"custom_attribute_source": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceCustomAttributeSource(),
			},
			"attribute_contract_fulfillment": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     resourceAttributeFulfillmentValue(),
			},
			"issuance_criteria": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceIssuanceCriteria(),
			},
		},
	}
}

func flattenAttributeMapping(in *pf.AttributeMapping) []map[string]interface{} {
	m := make([]map[string]interface{}, 0, 1)
	s := make(map[string]interface{})
	if in.AttributeContractFulfillment != nil {
		s["attribute_contract_fulfillment"] = flattenMapOfAttributeFulfillmentValue(in.AttributeContractFulfillment)
	}
	if in.IssuanceCriteria != nil && (in.IssuanceCriteria.ExpressionCriteria != nil && in.IssuanceCriteria.ConditionalCriteria != nil) {
		s["issuance_criteria"] = flattenIssuanceCriteria(in.IssuanceCriteria)
	}

	if *in.AttributeSources != nil && len(*in.AttributeSources) > 0 {
		var ldapAttributes []interface{}
		var jdbcAttributes []interface{}
		var customAttributes []interface{}
		for _, v := range *in.AttributeSources {
			switch *v.Type {
			case "LDAP":
				ldapAttributes = append(ldapAttributes, flattenLdapAttributeSource(&v.LdapAttributeSource))
			case "JDBC":
				jdbcAttributes = append(jdbcAttributes, flattenJdbcAttributeSource(v))
			case "CUSTOM":
				customAttributes = append(customAttributes, flattenCustomAttributeSource(&v.CustomAttributeSource))
			}
		}
		if len(ldapAttributes) > 0 {
			s["ldap_attribute_source"] = ldapAttributes
		}
		if len(jdbcAttributes) > 0 {
			s["jdbc_attribute_source"] = jdbcAttributes
		}
		if len(customAttributes) > 0 {
			s["custom_attribute_source"] = customAttributes
		}
	}
	m = append(m, s)
	return m
}

func expandAttributeMapping(in []interface{}) *pf.AttributeMapping {
	iac := &pf.AttributeMapping{AttributeSources: &[]*pf.AttributeSource{}}
	for _, raw := range in {
		l := raw.(map[string]interface{})
		if v, ok := l["attribute_contract_fulfillment"]; ok {
			iac.AttributeContractFulfillment = expandMapOfAttributeFulfillmentValue(v.(*schema.Set).List())
		}
		if v, ok := l["issuance_criteria"]; ok {
			iac.IssuanceCriteria = expandIssuanceCriteria(v.([]interface{}))
		}

		if v, ok := l["ldap_attribute_source"]; ok && len(v.([]interface{})) > 0 {
			*iac.AttributeSources = append(*iac.AttributeSources, *expandLdapAttributeSource(v.([]interface{}))...)
		}
		if v, ok := l["jdbc_attribute_source"]; ok && len(v.([]interface{})) > 0 {
			*iac.AttributeSources = append(*iac.AttributeSources, *expandJdbcAttributeSource(v.([]interface{}))...)
		}
		if v, ok := l["custom_attribute_source"]; ok && len(v.([]interface{})) > 0 {
			*iac.AttributeSources = append(*iac.AttributeSources, *expandCustomAttributeSource(v.([]interface{}))...)
		}

	}
	return iac
}

//func flattenScopeAttributeMappings(in map[string]*pf.ParameterValues) map[string][]interface{} {
//	s := make(map[string][]interface{})
//	for key, val := range in {
//		s[key] = flattenStringList(*val.Values)
//	}
//	return s
//}
//
//func expandScopeAttributeMappings(in map[string]interface{}) map[string]*pf.ParameterValues {
//	mappings := map[string]*pf.ParameterValues{}
//	m := expandMapOfLists(in)
//	for key, val := range m {
//		mappings[key] = &pf.ParameterValues{Values: &val}
//	}
//	return mappings
//}
//
//func expandMapOfLists(in map[string]interface{}) map[string][]*string {
//	m := map[string][]*string{}
//	for s := range in {
//		i := strings.LastIndex(s, ".")
//		first := s[0:i]
//		last := s[i+1:]
//		if last != "#" {
//			m[first] = append(m[first], String(in[s].(string)))
//		}
//	}
//	return m
//}

//func expandPluginConfigurationWithDescriptor(in []interface{}, desc *pf.PluginConfigDescriptor) *pf.PluginConfiguration {
//	config := expandPluginConfiguration(in)
//	for _, descriptor := range *desc.Fields {
//		if descriptor.DefaultValue != nil {
//			if !hasField(*descriptor.Name, config) {
//				*config.Fields = append(*config.Fields, &pf.ConfigField{Name: descriptor.Name, Value: descriptor.DefaultValue})
//			}
//		}
//	}
//
//	return config
//}
//
//func hasField(name string, c *pf.PluginConfiguration) bool {
//	for _, field := range *c.Fields {
//		if *field.Name == name {
//			return true
//		}
//	}
//	return false
//}
