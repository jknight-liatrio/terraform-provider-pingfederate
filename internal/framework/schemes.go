package framework

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iwarapter/terraform-provider-pingfederate/internal/framework/defaults"
)

func resourceApcToPersistentGrantMapping() schema.Schema {
	return schema.Schema{
		Description: `An authentication policy contract mapping into an OAuth persistent grant.`,
		Version:     1,
		Attributes: map[string]schema.Attribute{
			"attribute_contract_fulfillment": schema.MapNestedAttribute{
				Description: `A list of mappings from attribute names to their fulfillment values.`,
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: mapAttributeFulfillmentValue(),
				},
			},
			"jdbc_attribute_sources": schema.ListNestedAttribute{
				Description: `The configured settings used to look up attributes from a JDBC data store.`,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: listJdbcAttributeSource(),
				},
			},
			"ldap_attribute_sources": schema.ListNestedAttribute{
				Description: `The configured settings used to look up attributes from a LDAP data store.`,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: listLdapAttributeSource(),
				},
			},
			"custom_attribute_sources": schema.ListNestedAttribute{
				Description: `The configured settings used to look up attributes from a custom data store.`,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: listCustomAttributeSource(),
				},
			},
			"authentication_policy_contract_ref": schema.StringAttribute{
				Description: `Reference to the associated authentication policy contract. The reference cannot be changed after the mapping has been created.`,
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: `The ID of the authentication policy contract to persistent grant mapping.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"issuance_criteria": schema.SingleNestedAttribute{
				Description: `The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled.`,
				Optional:    true,
				Attributes:  singleIssuanceCriteria(),
			},
		},
	}
}

func resourceAuthenticationPolicyContract() schema.Schema {
	return schema.Schema{
		Description: `Authentication Policy Contracts carry user attributes from the identity provider to the service provider.`,
		Attributes: map[string]schema.Attribute{
			"core_attributes": schema.SetAttribute{
				ElementType: types.StringType,
				Description: `A list of read-only assertion attributes (for example, subject) that are automatically populated by PingFederate.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					defaults.NewDefaultSet(types.StringType, []attr.Value{types.StringValue("subject")}),
				},
			},
			"extended_attributes": schema.SetAttribute{
				ElementType: types.StringType,
				Description: `A list of additional attributes as needed.`,
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: `The persistent, unique ID for the authentication policy contract. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: `The Authentication Policy Contract Name. Name is unique.`,
				Required:    true,
			},
		},
	}
}

func resourceClient() schema.Schema {
	return schema.Schema{
		Description: `OAuth client.`,
		Version:     1,
		Attributes: map[string]schema.Attribute{
			"allow_authentication_api_init": schema.BoolAttribute{
				Description: `Set to true to allow this client to initiate the authentication API redirectless flow.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(false),
				},
			},
			"bypass_activation_code_confirmation_override": schema.BoolAttribute{
				Description: `Indicates if the Activation Code Confirmation page should be bypassed if 'verification_url_complete' is used by the end user to authorize a device. This overrides the 'bypassUseCodeConfirmation' value present in Authorization Server Settings.`,
				Optional:    true,
			},
			"bypass_approval_page": schema.BoolAttribute{
				Description: `Use this setting, for example, when you want to deploy a trusted application and authenticate end users via an IdP adapter or IdP connection.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(false),
				},
			},
			"ciba_delivery_mode": schema.StringAttribute{
				Description: `The token delivery mode for the client.  The default value is 'POLL'.`,
				Optional:    true,
			},
			"ciba_notification_endpoint": schema.StringAttribute{
				Description: `The endpoint the OP will call after a successful or failed end-user authentication.`,
				Optional:    true,
			},
			"ciba_polling_interval": schema.NumberAttribute{
				Description: `The minimum amount of time in seconds that the Client must wait between polling requests to the token endpoint. The default is 3 seconds.`,
				Optional:    true,
			},
			"ciba_request_object_signing_algorithm": schema.StringAttribute{
				Description: `The JSON Web Signature [JWS] algorithm that must be used to sign the CIBA Request Object. All signing algorithms are allowed if value is not present <br>RS256 - RSA using SHA-256<br>RS384 - RSA using SHA-384<br>RS512 - RSA using SHA-512<br>ES256 - ECDSA using P256 Curve and SHA-256<br>ES384 - ECDSA using P384 Curve and SHA-384<br>ES512 - ECDSA using P521 Curve and SHA-512<br>PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256<br>PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384<br>PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512<br>RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.`,
				Optional:    true,
			},
			"ciba_require_signed_requests": schema.BoolAttribute{
				Description: `Determines whether CIBA signed requests are required for this client.`,
				Optional:    true,
			},
			"ciba_user_code_supported": schema.BoolAttribute{
				Description: `Determines whether CIBA user code is supported for this client.`,
				Optional:    true,
			},
			"client_auth": schema.SingleNestedAttribute{
				Description: `Client authentication settings.  If this model is null, it indicates that no client authentication will be used.`,
				Optional:    true,
				Attributes:  singleClientAuth(),
			},
			"client_id": schema.StringAttribute{
				Description: `A unique identifier the client provides to the Resource Server to identify itself. This identifier is included with every request the client makes. For PUT requests, this field is optional and it will be overridden by the 'id' parameter of the PUT request.`,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_secret_changed_time": schema.StringAttribute{
				Description: `The time at which the client secret was last changed. This property is read only and is ignored on PUT and POST requests.`,
				Optional:    true,
				Computed:    true,
			},
			"client_secret_retention_period": schema.NumberAttribute{
				Description: `The length of time in minutes that client secrets will be retained as secondary secrets after secret change. The default value is 0, which will disable secondary client secret retention. This value will override the Client Secret Retention Period value on the Authorization Server Settings.`,
				Optional:    true,
			},
			"client_secret_retention_period_type": schema.StringAttribute{
				Description: `Use OVERRIDE_SERVER_DEFAULT to override the Client Secret Retention Period value on the Authorization Server Settings. SERVER_DEFAULT will default to the Client Secret Retention Period value on the Authorization Server Setting. Defaults to SERVER_DEFAULT.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("SERVER_DEFAULT"),
				},
			},
			"default_access_token_manager_ref": schema.StringAttribute{
				Description: `The default access token manager for this client.`,
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: `A description of what the client application does. This description appears when the user is prompted for authorization.`,
				Optional:    true,
			},
			"device_flow_setting_type": schema.StringAttribute{
				Description: `Allows an administrator to override the Device Authorization Settings set globally for the OAuth AS. Defaults to SERVER_DEFAULT.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("SERVER_DEFAULT"),
				},
			},
			"device_polling_interval_override": schema.NumberAttribute{
				Description: `The amount of time client should wait between polling requests, in seconds. This overrides the 'devicePollingInterval' value present in Authorization Server Settings.`,
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: `Specifies whether the client is enabled. The default value is true.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(true),
				},
			},
			"exclusive_scopes": schema.ListAttribute{
				Description: `The exclusive scopes available for this client.`,
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					defaults.NewDefaultList(types.StringType, []attr.Value{}),
				},
			},
			"extended_parameters": schema.MapNestedAttribute{
				Description: `OAuth Client Metadata can be extended to use custom Client Metadata Parameters. The names of these custom parameters should be defined in /extendedProperties.`,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: mapParameterValues(),
				},
			},
			"grant_types": schema.ListAttribute{
				Description: `The grant types allowed for this client. The EXTENSION grant type applies to SAML/JWT assertion grants.`,
				Required:    true,
				ElementType: types.StringType,
			},
			"id": schema.StringAttribute{
				Description: `The client_id of the oauth client.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"jwks_settings": schema.SingleNestedAttribute{
				Description: `JSON Web Key Set Settings of the OAuth client. Required if private key JWT client authentication or signed requests is enabled.`,
				Optional:    true,
				Attributes:  singleJwksSettings(),
			},
			"jwt_secured_authorization_response_mode_content_encryption_algorithm": schema.StringAttribute{
				Description: `The JSON Web Encryption [JWE] content-encryption algorithm for the JWT Secured Authorization Response.<br>AES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256<br>AES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384<br>AES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512<br>AES_128_GCM - AES-GCM-128<br>AES_192_GCM - AES-GCM-192<br>AES_256_GCM - AES-GCM-256`,
				Optional:    true,
			},
			"jwt_secured_authorization_response_mode_encryption_algorithm": schema.StringAttribute{
				Description: `The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the JWT Secured Authorization Response.<br>DIR - Direct Encryption with symmetric key<br>A128KW - AES-128 Key Wrap<br>A192KW - AES-192 Key Wrap<br>A256KW - AES-256 Key Wrap<br>A128GCMKW - AES-GCM-128 key encryption<br>A192GCMKW - AES-GCM-192 key encryption<br>A256GCMKW - AES-GCM-256 key encryption<br>ECDH_ES - ECDH-ES<br>ECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap<br>ECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap<br>ECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap<br>RSA_OAEP - RSAES OAEP<br>RSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256`,
				Optional:    true,
			},
			"jwt_secured_authorization_response_mode_signing_algorithm": schema.StringAttribute{
				Description: `The JSON Web Signature [JWS] algorithm required to sign the JWT Secured Authorization Response.<br>HS256 - HMAC using SHA-256<br>HS384 - HMAC using SHA-384<br>HS512 - HMAC using SHA-512<br>RS256 - RSA using SHA-256<br>RS384 - RSA using SHA-384<br>RS512 - RSA using SHA-512<br>ES256 - ECDSA using P256 Curve and SHA-256<br>ES384 - ECDSA using P384 Curve and SHA-384<br>ES512 - ECDSA using P521 Curve and SHA-512<br>PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256<br>PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384<br>PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512<br>A null value will represent the default algorithm which is RS256.<br>RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11`,
				Optional:    true,
			},
			"logo_url": schema.StringAttribute{
				Description: `The location of the logo used on user-facing OAuth grant authorization and revocation pages.`,
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: `A descriptive name for the client instance. This name appears when the user is prompted for authorization.`,
				Required:    true,
			},
			"oidc_policy": schema.SingleNestedAttribute{
				Description: `Open ID Connect Policy settings.  This is included in the message only when OIDC is enabled.`,
				Optional:    true,
				Attributes:  singleClientOIDCPolicy(),
			},
			"pending_authorization_timeout_override": schema.NumberAttribute{
				Description: `The 'device_code' and 'user_code' timeout, in seconds. This overrides the 'pendingAuthorizationTimeout' value present in Authorization Server Settings.`,
				Optional:    true,
			},
			"persistent_grant_expiration_time": schema.NumberAttribute{
				Description: `The persistent grant expiration time. -1 indicates an indefinite amount of time.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Number{
					defaults.NewDefaultNumber(big.NewFloat(0)),
				},
			},
			"persistent_grant_expiration_time_unit": schema.StringAttribute{
				Description: `The persistent grant expiration time unit.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("DAYS"),
				},
			},
			"persistent_grant_expiration_type": schema.StringAttribute{
				Description: `Allows an administrator to override the Persistent Grant Lifetime set globally for the OAuth AS. Defaults to SERVER_DEFAULT.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("SERVER_DEFAULT"),
				},
			},
			"persistent_grant_idle_timeout": schema.NumberAttribute{
				Description: `The persistent grant idle timeout.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Number{
					defaults.NewDefaultNumber(big.NewFloat(0)),
				},
			},
			"persistent_grant_idle_timeout_time_unit": schema.StringAttribute{
				Description: `The persistent grant idle timeout time unit.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("DAYS"),
				},
			},
			"persistent_grant_idle_timeout_type": schema.StringAttribute{
				Description: `Allows an administrator to override the Persistent Grant Idle Timeout set globally for the OAuth AS. Defaults to SERVER_DEFAULT.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("SERVER_DEFAULT"),
				},
			},
			"persistent_grant_reuse_grant_types": schema.ListAttribute{
				Description: `The grant types that the OAuth AS can reuse rather than creating a new grant for each request. This value will override the Reuse Existing Persistent Access Grants for Grant Types on the Authorization Server Settings. Only 'IMPLICIT' or 'AUTHORIZATION_CODE' or 'RESOURCE_OWNER_CREDENTIALS' are valid grant types.`,
				Optional:    true,
				ElementType: types.StringType,
			},
			"persistent_grant_reuse_type": schema.StringAttribute{
				Description: `Allows and administrator to override the Reuse Existing Persistent Access Grants for Grant Types set globally for OAuth AS. Defaults to SERVER_DEFAULT.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("SERVER_DEFAULT"),
				},
			},
			"redirect_uris": schema.ListAttribute{
				Description: `URIs to which the OAuth AS may redirect the resource owner's user agent after authorization is obtained. A redirection URI is used with the Authorization Code and Implicit grant types. Wildcards are allowed. However, for security reasons, make the URL as restrictive as possible.For example: https://*.company.com/* Important: If more than one URI is added or if a single URI uses wildcards, then Authorization Code grant and token requests must contain a specific matching redirect uri parameter.`,
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					defaults.NewDefaultList(types.StringType, []attr.Value{}),
				},
			},
			"refresh_rolling": schema.StringAttribute{
				Description: `Use ROLL or DONT_ROLL to override the Roll Refresh Token Values setting on the Authorization Server Settings. SERVER_DEFAULT will default to the Roll Refresh Token Values setting on the Authorization Server Setting screen. Defaults to SERVER_DEFAULT.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("SERVER_DEFAULT"),
				},
			},
			"refresh_token_rolling_grace_period": schema.NumberAttribute{
				Description: `The grace period that a rolled refresh token remains valid in seconds.`,
				Optional:    true,
			},
			"refresh_token_rolling_grace_period_type": schema.StringAttribute{
				Description: `When specified, it overrides the global Refresh Token Grace Period defined in the Authorization Server Settings. The default value is SERVER_DEFAULT`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("SERVER_DEFAULT"),
				},
			},
			"refresh_token_rolling_interval": schema.NumberAttribute{
				Description: `The minimum interval to roll refresh tokens, in hours. This value will override the Refresh Token Rolling Interval Value on the Authorization Server Settings.`,
				Optional:    true,
			},
			"refresh_token_rolling_interval_type": schema.StringAttribute{
				Description: `Use OVERRIDE_SERVER_DEFAULT to override the Refresh Token Rolling Interval value on the Authorization Server Settings. SERVER_DEFAULT will default to the Refresh Token Rolling Interval value on the Authorization Server Setting. Defaults to SERVER_DEFAULT.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					defaults.NewDefaultString("SERVER_DEFAULT"),
				},
			},
			"request_object_signing_algorithm": schema.StringAttribute{
				Description: `The JSON Web Signature [JWS] algorithm that must be used to sign the Request Object. All signing algorithms are allowed if value is not present <br>RS256 - RSA using SHA-256<br>RS384 - RSA using SHA-384<br>RS512 - RSA using SHA-512<br>ES256 - ECDSA using P256 Curve and SHA-256<br>ES384 - ECDSA using P384 Curve and SHA-384<br>ES512 - ECDSA using P521 Curve and SHA-512<br>PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256<br>PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384<br>PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512<br>RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.`,
				Optional:    true,
			},
			"request_policy_ref": schema.StringAttribute{
				Description: `The CIBA request policy.`,
				Optional:    true,
			},
			"require_jwt_secured_authorization_response_mode": schema.BoolAttribute{
				Description: `Determines whether JWT Secured authorization response mode is required when initiating an authorization request. The default is false.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(false),
				},
			},
			"require_proof_key_for_code_exchange": schema.BoolAttribute{
				Description: `Determines whether Proof Key for Code Exchange (PKCE) is required for this client.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(false),
				},
			},
			"require_pushed_authorization_requests": schema.BoolAttribute{
				Description: `Determines whether pushed authorization requests are required when initiating an authorization request. The default is false.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(false),
				},
			},
			"require_signed_requests": schema.BoolAttribute{
				Description: `Determines whether signed requests are required for this client`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(false),
				},
			},
			"restrict_scopes": schema.BoolAttribute{
				Description: `Restricts this client's access to specific scopes.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(false),
				},
			},
			"restrict_to_default_access_token_manager": schema.BoolAttribute{
				Description: `Determines whether the client is restricted to using only its default access token manager. The default is false.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(false),
				},
			},
			"restricted_response_types": schema.ListAttribute{
				Description: `The response types allowed for this client. If omitted all response types are available to the client.`,
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					defaults.NewDefaultList(types.StringType, []attr.Value{}),
				},
			},
			"restricted_scopes": schema.ListAttribute{
				Description: `The scopes available for this client.`,
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					defaults.NewDefaultList(types.StringType, []attr.Value{}),
				},
			},
			"token_exchange_processor_policy_ref": schema.StringAttribute{
				Description: `The Token Exchange Processor policy.`,
				Optional:    true,
			},
			"token_introspection_content_encryption_algorithm": schema.StringAttribute{
				Description: `The JSON Web Encryption [JWE] content-encryption algorithm for the Token Introspection Response.<br>AES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256<br>AES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384<br>AES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512<br>AES_128_GCM - AES-GCM-128<br>AES_192_GCM - AES-GCM-192<br>AES_256_GCM - AES-GCM-256`,
				Optional:    true,
			},
			"token_introspection_encryption_algorithm": schema.StringAttribute{
				Description: `The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content-encryption key of the Token Introspection Response.<br>DIR - Direct Encryption with symmetric key<br>A128KW - AES-128 Key Wrap<br>A192KW - AES-192 Key Wrap<br>A256KW - AES-256 Key Wrap<br>A128GCMKW - AES-GCM-128 key encryption<br>A192GCMKW - AES-GCM-192 key encryption<br>A256GCMKW - AES-GCM-256 key encryption<br>ECDH_ES - ECDH-ES<br>ECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap<br>ECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap<br>ECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap<br>RSA_OAEP - RSAES OAEP<br>RSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256`,
				Optional:    true,
			},
			"token_introspection_signing_algorithm": schema.StringAttribute{
				Description: `The JSON Web Signature [JWS] algorithm required to sign the Token Introspection Response.<br>HS256 - HMAC using SHA-256<br>HS384 - HMAC using SHA-384<br>HS512 - HMAC using SHA-512<br>RS256 - RSA using SHA-256<br>RS384 - RSA using SHA-384<br>RS512 - RSA using SHA-512<br>ES256 - ECDSA using P256 Curve and SHA-256<br>ES384 - ECDSA using P384 Curve and SHA-384<br>ES512 - ECDSA using P521 Curve and SHA-512<br>PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256<br>PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384<br>PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512<br>A null value will represent the default algorithm which is RS256.<br>RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11`,
				Optional:    true,
			},
			"user_authorization_url_override": schema.StringAttribute{
				Description: `The URL used as 'verification_url' and 'verification_url_complete' values in a Device Authorization request. This property overrides the 'userAuthorizationUrl' value present in Authorization Server Settings.`,
				Optional:    true,
			},
			"validate_using_all_eligible_atms": schema.BoolAttribute{
				Description: `Validates token using all eligible access token managers for the client. This setting is ignored if 'restrictToDefaultAccessTokenManager' is set to true.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					defaults.NewDefaultBool(false),
				},
			},
		},
	}
}

func resourceRedirectValidationSettings() schema.Schema {
	return schema.Schema{
		Description: `Settings for redirect validation for SSO, SLO and IdP discovery.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: ``,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"redirect_validation_local_settings": schema.SingleNestedAttribute{
				Description: `Settings for local redirect validation.`,
				Optional:    true,
				Attributes:  singleRedirectValidationLocalSettings(),
			},
			"redirect_validation_partner_settings": schema.SingleNestedAttribute{
				Description: `Settings for redirection at a partner site.`,
				Optional:    true,
				Attributes:  singleRedirectValidationPartnerSettings(),
			},
		},
	}
}

func listJdbcAttributeSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"attribute_contract_fulfillment": schema.MapNestedAttribute{
			Description: `A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings`,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: mapAttributeFulfillmentValue(),
			},
		},
		"column_names": schema.ListAttribute{
			Description: `A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.`,
			Optional:    true,
			ElementType: types.StringType,
		},
		"data_store_ref": schema.StringAttribute{
			Description: `Reference to the associated data store.`,
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: `The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings`,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				defaults.NewDefaultString("JDBC"),
			},
		},
		"filter": schema.StringAttribute{
			Description: `The JDBC WHERE clause used to query your data store to locate a user record.`,
			Optional:    true,
		},
		"id": schema.StringAttribute{
			Description: `The ID that defines this attribute source. Only alphanumeric characters allowed.<br>Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.`,
			Optional:    true,
		},
		"schema": schema.StringAttribute{
			Description: `Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.`,
			Optional:    true,
		},
		"table": schema.StringAttribute{
			Description: `The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.`,
			Optional:    true,
		},
	}
}

func listLdapAttributeSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"attribute_contract_fulfillment": schema.MapNestedAttribute{
			Description: `A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings`,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: mapAttributeFulfillmentValue(),
			},
		},
		"base_dn": schema.StringAttribute{
			Description: `The base DN to search from. If not specified, the search will start at the LDAP's root.`,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				defaults.NewDefaultString(""),
			},
		},
		"binary_attribute_settings": schema.MapNestedAttribute{
			Description: `The advanced settings for binary LDAP attributes.`,
			Optional:    true,
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: mapBinaryLdapAttributeSettings(),
			},
			PlanModifiers: []planmodifier.Map{
				defaults.NewDefaultMap(types.ObjectType{AttrTypes: map[string]attr.Type{"binary_encoding": types.StringType}}, map[string]attr.Value{}),
			},
		},
		"data_store_ref": schema.StringAttribute{
			Description: `Reference to the associated data store.`,
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: `The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings`,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				defaults.NewDefaultString("JDBC"),
				defaults.NewDefaultString("JDBC"),
			},
		},
		"id": schema.StringAttribute{
			Description: `The ID that defines this attribute source. Only alphanumeric characters allowed.<br>Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.`,
			Optional:    true,
		},
		"member_of_nested_group": schema.BoolAttribute{
			Description: `Set this to true to return transitive group memberships for the 'memberOf' attribute.  This only applies for Active Directory data sources.  All other data sources will be set to false.`,
			Optional:    true,
		},
		"search_attributes": schema.ListAttribute{
			Description: `A list of LDAP attributes returned from search and available for mapping.`,
			Optional:    true,
			ElementType: types.StringType,
		},
		"search_filter": schema.StringAttribute{
			Description: `The LDAP filter that will be used to lookup the objects from the directory.`,
			Optional:    true,
		},
		"search_scope": schema.StringAttribute{
			Description: `Determines the node depth of the query.`,
			Optional:    true,
		},
	}
}

func listCustomAttributeSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"attribute_contract_fulfillment": schema.MapNestedAttribute{
			Description: `A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings`,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: mapAttributeFulfillmentValue(),
			},
		},
		"data_store_ref": schema.StringAttribute{
			Description: `Reference to the associated data store.`,
			Required:    true,
		},
		"description": schema.StringAttribute{
			Description: `The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings`,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				defaults.NewDefaultString("JDBC"),
				defaults.NewDefaultString("JDBC"),
				defaults.NewDefaultString("JDBC"),
			},
		},
		"filter_fields": schema.ListNestedAttribute{
			Description: `The list of fields that can be used to filter a request to the custom data store.`,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: listFieldEntry(),
			},
		},
		"id": schema.StringAttribute{
			Description: `The ID that defines this attribute source. Only alphanumeric characters allowed.<br>Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.`,
			Optional:    true,
		},
	}
}

func mapAttributeFulfillmentValue() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"source": schema.SingleNestedAttribute{
			Description: `The attribute value source.`,
			Required:    true,
			Attributes:  singleSourceTypeIdKey(),
		},
		"value": schema.StringAttribute{
			Description: `The value for this attribute.`,
			Required:    true,
		},
	}
}

func mapBinaryLdapAttributeSettings() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"binary_encoding": schema.StringAttribute{
			Description: `Get the encoding type for this attribute. If not specified, the default is BASE64.`,
			Optional:    true,
		},
	}
}

func singleClientAuth() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"client_cert_issuer_dn": schema.StringAttribute{
			Description: `Client TLS Certificate Issuer DN.`,
			Optional:    true,
		},
		"client_cert_subject_dn": schema.StringAttribute{
			Description: `Client TLS Certificate Subject DN.`,
			Optional:    true,
		},
		"encrypted_secret": schema.StringAttribute{
			Description: `For GET requests, this field contains the encrypted client secret, if one exists.  For POST and PUT requests, if you wish to reuse the existing secret, this field should be passed back unchanged.`,
			Optional:    true,
			Computed:    true,
		},
		"enforce_replay_prevention": schema.BoolAttribute{
			Description: `Enforce replay prevention on JSON Web Tokens. This field is applicable only for Private Key JWT Client Authentication.`,
			Optional:    true,
		},
		"secret": schema.StringAttribute{
			Description: `Client secret for Basic Authentication.  To update the client secret, specify the plaintext value in this field.  This field will not be populated for GET requests.`,
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"token_endpoint_auth_signing_algorithm": schema.StringAttribute{
			Description: `The JSON Web Signature [JWS] algorithm that must be used to sign the JSON Web Tokens. This field is applicable only for Private Key JWT Client Authentication. All signing algorithms are allowed if value is not present <br>RS256 - RSA using SHA-256<br>RS384 - RSA using SHA-384<br>RS512 - RSA using SHA-512<br>ES256 - ECDSA using P256 Curve and SHA-256<br>ES384 - ECDSA using P384 Curve and SHA-384<br>ES512 - ECDSA using P521 Curve and SHA-512<br>PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256<br>PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384<br>PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512<br>RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11.`,
			Optional:    true,
		},
		"type": schema.StringAttribute{
			Description: `Client authentication type.<br>The required field for type SECRET is secret.<br>The required fields for type CERTIFICATE are clientCertIssuerDn and clientCertSubjectDn.<br>The required field for type PRIVATE_KEY_JWT is: either jwks or jwksUrl.`,
			Optional:    true,
		},
	}
}

func singleClientOIDCPolicy() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"grant_access_session_revocation_api": schema.BoolAttribute{
			Description: `Determines whether this client is allowed to access the Session Revocation API.`,
			Optional:    true,
		},
		"grant_access_session_session_management_api": schema.BoolAttribute{
			Description: `Determines whether this client is allowed to access the Session Management API.`,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				defaults.NewDefaultBool(false),
			},
		},
		"id_token_content_encryption_algorithm": schema.StringAttribute{
			Description: `The JSON Web Encryption [JWE] content encryption algorithm for the ID Token.<br>AES_128_CBC_HMAC_SHA_256 - Composite AES-CBC-128 HMAC-SHA-256<br>AES_192_CBC_HMAC_SHA_384 - Composite AES-CBC-192 HMAC-SHA-384<br>AES_256_CBC_HMAC_SHA_512 - Composite AES-CBC-256 HMAC-SHA-512<br>AES_128_GCM - AES-GCM-128<br>AES_192_GCM - AES-GCM-192<br>AES_256_GCM - AES-GCM-256`,
			Optional:    true,
		},
		"id_token_encryption_algorithm": schema.StringAttribute{
			Description: `The JSON Web Encryption [JWE] encryption algorithm used to encrypt the content encryption key for the ID Token.<br>DIR - Direct Encryption with symmetric key<br>A128KW - AES-128 Key Wrap<br>A192KW - AES-192 Key Wrap<br>A256KW - AES-256 Key Wrap<br>A128GCMKW - AES-GCM-128 key encryption<br>A192GCMKW - AES-GCM-192 key encryption<br>A256GCMKW - AES-GCM-256 key encryption<br>ECDH_ES - ECDH-ES<br>ECDH_ES_A128KW - ECDH-ES with AES-128 Key Wrap<br>ECDH_ES_A192KW - ECDH-ES with AES-192 Key Wrap<br>ECDH_ES_A256KW - ECDH-ES with AES-256 Key Wrap<br>RSA_OAEP - RSAES OAEP<br>RSA_OAEP_256 - RSAES OAEP using SHA-256 and MGF1 with SHA-256`,
			Optional:    true,
		},
		"id_token_signing_algorithm": schema.StringAttribute{
			Description: `The JSON Web Signature [JWS] algorithm required for the ID Token.<br>NONE - No signing algorithm<br>HS256 - HMAC using SHA-256<br>HS384 - HMAC using SHA-384<br>HS512 - HMAC using SHA-512<br>RS256 - RSA using SHA-256<br>RS384 - RSA using SHA-384<br>RS512 - RSA using SHA-512<br>ES256 - ECDSA using P256 Curve and SHA-256<br>ES384 - ECDSA using P384 Curve and SHA-384<br>ES512 - ECDSA using P521 Curve and SHA-512<br>PS256 - RSASSA-PSS using SHA-256 and MGF1 padding with SHA-256<br>PS384 - RSASSA-PSS using SHA-384 and MGF1 padding with SHA-384<br>PS512 - RSASSA-PSS using SHA-512 and MGF1 padding with SHA-512<br>A null value will represent the default algorithm which is RS256.<br>RSASSA-PSS is only supported with SafeNet Luna, Thales nCipher or Java 11`,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				defaults.NewDefaultString("RS256"),
			},
		},
		"logout_uris": schema.ListAttribute{
			Description: `A list of client logout URI's which will be invoked when a user logs out through one of PingFederate's SLO endpoints.`,
			Optional:    true,
			ElementType: types.StringType,
		},
		"pairwise_identifier_user_type": schema.BoolAttribute{
			Description: `Determines whether the subject identifier type is pairwise.`,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				defaults.NewDefaultBool(false),
			},
		},
		"ping_access_logout_capable": schema.BoolAttribute{
			Description: `Set this value to true if you wish to enable client application logout, and the client is PingAccess, or its logout endpoints follow the PingAccess path convention.`,
			Optional:    true,
		},
		"policy_group": schema.StringAttribute{
			Description: `The Open ID Connect policy. A null value will represent the default policy group.`,
			Optional:    true,
		},
		"sector_identifier_uri": schema.StringAttribute{
			Description: `The URI references a file with a single JSON array of Redirect URI and JWKS URL values.`,
			Optional:    true,
		},
	}
}

func listConditionalIssuanceCriteriaEntry() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"attribute_name": schema.StringAttribute{
			Description: `The name of the attribute to use in this issuance criterion.`,
			Required:    true,
		},
		"condition": schema.StringAttribute{
			Description: `The condition that will be applied to the source attribute's value and the expected value.`,
			Required:    true,
		},
		"error_result": schema.StringAttribute{
			Description: `The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.`,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				defaults.NewDefaultString(""),
			},
		},
		"source": schema.SingleNestedAttribute{
			Description: `The source of the attribute.`,
			Required:    true,
			Attributes:  singleSourceTypeIdKey(),
		},
		"value": schema.StringAttribute{
			Description: `The expected value of this issuance criterion.`,
			Required:    true,
		},
	}
}

func listExpressionIssuanceCriteriaEntry() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"error_result": schema.StringAttribute{
			Description: `The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.`,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				defaults.NewDefaultString(""),
			},
		},
		"expression": schema.StringAttribute{
			Description: `The OGNL expression to evaluate.`,
			Required:    true,
		},
	}
}

func listFieldEntry() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: `The name of this field.`,
			Required:    true,
		},
		"value": schema.StringAttribute{
			Description: `The value of this field. Whether or not the value is required will be determined by plugin validation checks.`,
			Optional:    true,
		},
	}
}

func singleIssuanceCriteria() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"conditional_criteria": schema.ListNestedAttribute{
			Description: `A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue.`,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: listConditionalIssuanceCriteriaEntry(),
			},
		},
		"expression_criteria": schema.ListNestedAttribute{
			Description: `A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.`,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: listExpressionIssuanceCriteriaEntry(),
			},
		},
	}
}

func singleJwksSettings() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"jwks": schema.StringAttribute{
			Description: `JSON Web Key Set (JWKS) document of the OAuth client. Either 'jwks' or 'jwksUrl' must be provided if private key JWT client authentication or signed requests is enabled.  If the client signs its JWTs using an RSASSA-PSS signing algorithm, PingFederate must either use Java 11 or be integrated with a hardware security module (HSM) to process the digital signatures.`,
			Optional:    true,
		},
		"jwks_url": schema.StringAttribute{
			Description: `JSON Web Key Set (JWKS) URL of the OAuth client. Either 'jwks' or 'jwksUrl' must be provided if private key JWT client authentication or signed requests is enabled.  If the client signs its JWTs using an RSASSA-PSS signing algorithm, PingFederate must either use Java 11 or be integrated with a hardware security module (HSM) to process the digital signatures.`,
			Optional:    true,
		},
	}
}

func mapParameterValues() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"values": schema.ListAttribute{
			Description: `A List of values`,
			Optional:    true,
			ElementType: types.StringType,
		},
	}
}

func singleRedirectValidationLocalSettings() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enable_in_error_resource_validation": schema.BoolAttribute{
			Description: `Enable validation for error resource.`,
			Optional:    true,
		},
		"enable_target_resource_validation_for_idp_discovery": schema.BoolAttribute{
			Description: `Enable target resource validation for IdP discovery.`,
			Optional:    true,
		},
		"enable_target_resource_validation_for_slo": schema.BoolAttribute{
			Description: `Enable target resource validation for SLO.`,
			Optional:    true,
		},
		"enable_target_resource_validation_for_sso": schema.BoolAttribute{
			Description: `Enable target resource validation for SSO.`,
			Optional:    true,
		},
		"white_list": schema.ListNestedAttribute{
			Description: `List of URLs that are designated as valid target resources.`,
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: listRedirectValidationSettingsWhitelistEntry(),
			},
		},
	}
}

func singleRedirectValidationPartnerSettings() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enable_wreply_validation_slo": schema.BoolAttribute{
			Description: `Enable wreply validation for SLO.`,
			Optional:    true,
		},
	}
}

func listRedirectValidationSettingsWhitelistEntry() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"allow_query_and_fragment": schema.BoolAttribute{
			Description: `Allow any query parameters and fragment in the resource.`,
			Optional:    true,
		},
		"idp_discovery": schema.BoolAttribute{
			Description: `Enable this target resource for IdP discovery validation.`,
			Optional:    true,
		},
		"in_error_resource": schema.BoolAttribute{
			Description: `Enable this target resource for in error resource validation.`,
			Optional:    true,
		},
		"require_https": schema.BoolAttribute{
			Description: `Require HTTPS for accessing this resource.`,
			Optional:    true,
		},
		"target_resource_slo": schema.BoolAttribute{
			Description: `Enable this target resource for SLO redirect validation.`,
			Optional:    true,
		},
		"target_resource_sso": schema.BoolAttribute{
			Description: `Enable this target resource for SSO redirect validation.`,
			Optional:    true,
		},
		"valid_domain": schema.StringAttribute{
			Description: `Domain of a valid resource.`,
			Required:    true,
		},
		"valid_path": schema.StringAttribute{
			Description: `Path of a valid resource.`,
			Optional:    true,
		},
	}
}

func singleSourceTypeIdKey() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: `The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.`,
			Optional:    true,
		},
		"type": schema.StringAttribute{
			Description: `The source type of this key.`,
			Required:    true,
		},
	}
}
