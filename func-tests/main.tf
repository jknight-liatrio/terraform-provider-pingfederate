terraform {
  required_providers {
    pingfederate = {
      source  = "iwarapter/pingfederate"
      #for functional testing
    }
  }
}

provider "pingfederate" {
  bypass_external_validation = true
  password                   = "2FederateM0re"
}

locals {
  isPF10_0 = length(regexall("10.[0]", data.pingfederate_version.instance.version)) > 0
  isPF10_1 = length(regexall("10.[1]", data.pingfederate_version.instance.version)) > 0
  isPF10_2 = length(regexall("10.[2]", data.pingfederate_version.instance.version)) > 0
}

data "pingfederate_version" "instance" {}

resource "pingfederate_server_settings" "settings" {
  federation_info {
    base_url         = "https://localhost:9031"
    saml2_entity_id  = "testing"
    saml1x_issuer_id = "testing"
    wsfed_realm      = "testing"
  }
  roles_and_protocols {
    enable_idp_discovery = true
    idp_role {
      enable                       = true
      enable_outbound_provisioning = true
      enable_saml10                = true
      enable_saml11                = true
      enable_ws_fed                = true
      enable_ws_trust              = true
      saml20_profile {
        enable = true
      }
    }
    oauth_role {
      enable_oauth          = true
      enable_openid_connect = true
    }
    sp_role {
      enable                      = true
      enable_inbound_provisioning = true
      enable_openid_connect       = true
      enable_saml10               = true
      enable_saml11               = true
      enable_ws_fed               = true
      enable_ws_trust             = true
      saml20_profile {
        enable      = true
        enable_xasp = true
      }
    }
  }
}
//
resource "pingfederate_oauth_auth_server_settings" "settings" {
  scopes {
    name        = "address"
    description = "address"
  }

  scopes {
    name        = "mail"
    description = "mail"
  }

  scopes {
    name        = "openid"
    description = "openid"
  }

  scopes {
    name        = "phone"
    description = "phone"
  }

  scopes {
    name        = "profile"
    description = "profile"
  }

  scope_groups {
    name        = "group1"
    description = "group1"

    scopes = [
      "address",
      "mail",
      "phone",
      "openid",
      "profile",
    ]
  }

  persistent_grant_contract {
    extended_attributes = ["woot"]
  }

  allowed_origins = [
    "http://localhost",
  ]

  default_scope_description  = ""
  authorization_code_timeout = 60
  authorization_code_entropy = 30
  refresh_token_length       = 42
  refresh_rolling_interval   = 0
}

resource "pingfederate_oauth_client" "woot" {
  client_id = "woot"
  name      = "woot"

  grant_types = [
    "EXTENSION",
  ]

  client_auth = {
    // type                      = "CERTIFICATE"
    // client_cert_issuer_dn     = ""
    // client_cert_subject_dn    = ""
    enforce_replay_prevention = false

    secret = "super_top_secret"
    type   = "SECRET"
  }

  // jwks_settings {
  //   jwks = "https://stuff"
  // }
  default_access_token_manager_ref = pingfederate_oauth_access_token_manager.my_atm.id

  oidc_policy = {
    grant_access_session_revocation_api = false

    logout_uris = [
      "https://logout",
      "https://foo",
    ]

    ping_access_logout_capable = true
  }

  depends_on = [
    pingfederate_server_settings.settings,
    pingfederate_oauth_auth_server_settings.settings
  ]
}

resource "pingfederate_oauth_access_token_manager" "my_atm" {
  instance_id = "myatat"
  name        = "my_atat"

  plugin_descriptor_ref {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }

  configuration {
    fields {
      name  = "Token Length"
      value = "28"
    }

    fields {
      name  = "Token Lifetime"
      value = "120"
    }

    fields {
      name  = "Lifetime Extension Policy"
      value = "ALL"
    }

    fields {
      name  = "Maximum Token Lifetime"
      value = ""
    }

    fields {
      name  = "Lifetime Extension Threshold Percentage"
      value = "30"
    }

    fields {
      name  = "Mode for Synchronous RPC"
      value = "3"
    }

    fields {
      name  = "RPC Timeout"
      value = "500"
    }

    fields {
      name  = "Expand Scope Groups"
      value = "false"
    }
  }

  attribute_contract {
    extended_attributes = [
    "sub"]
  }

  depends_on = [
    pingfederate_server_settings.settings,
    pingfederate_oauth_auth_server_settings.settings
  ]
}

resource "pingfederate_authentication_policy_contract" "apc_foo" {
  name                = "example"
  extended_attributes = ["foo", "bar"]
}

resource "pingfederate_oauth_access_token_mappings" "auth_policy_mapping_demo" {
  access_token_manager_ref {
    id = pingfederate_oauth_access_token_manager.my_atm.id
  }

  context {
    type = "AUTHENTICATION_POLICY_CONTRACT"
    context_ref {
      id = pingfederate_authentication_policy_contract.apc_foo.id
    }
  }
  attribute_contract_fulfillment {
    key_name = "sub"
    source {
      type = "NO_MAPPING"
    }
  }
}

resource "pingfederate_oauth_resource_owner_credentials_mappings" "demo" {
  password_validator_ref {
    id = pingfederate_password_credential_validator.demo.id
  }

  attribute_contract_fulfillment {
    key_name = "USER_KEY"
    source {
      type = "NO_MAPPING"
    }
  }

  attribute_contract_fulfillment {
    key_name = "woot"
    source {
      type = "NO_MAPPING"
    }
  }
  issuance_criteria {
    conditional_criteria {
      attribute_name = "username"
      condition      = "EQUALS"
      error_result   = "deny"
      value          = "foo"

      source {
        type = "PASSWORD_CREDENTIAL_VALIDATOR"
      }
    }
  }

  depends_on = [
    pingfederate_server_settings.settings,
    pingfederate_oauth_auth_server_settings.settings
  ]
}

resource "pingfederate_password_credential_validator" "demo" {
  name = "foo"
  plugin_descriptor_ref {
    id = "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator"
  }

  configuration {
    tables {
      name = "Users"
      rows {
        fields {
          name  = "Username"
          value = "bob"
        }

        sensitive_fields {
          name  = "Password"
          value = "demo2"
        }

        sensitive_fields {
          name  = "Confirm Password"
          value = "demo2"
        }

        fields {
          name  = "Relax Password Requirements"
          value = "true"
        }
      }
    }
  }
}

resource "pingfederate_jdbc_data_store" "demo" {
  name           = "terraform"
  driver_class   = "org.hsqldb.jdbcDriver"
  user_name      = "sa"
  password       = "secret"
  connection_url = "jdbc:hsqldb:mem:mymemdb"
  connection_url_tags {
    connection_url = "jdbc:hsqldb:mem:mymemdb"
    default_source = true
  }
}

resource "pingfederate_ldap_data_store" "demo_ldap" {
  name             = "terraform_ldap"
  ldap_type        = "PING_DIRECTORY"
  hostnames        = ["host.docker.internal:1389"]
  bind_anonymously = true
  min_connections  = 1
  max_connections  = 1
}

resource "pingfederate_idp_adapter" "demo" {
  name = "bart"
  plugin_descriptor_ref {
    id = "com.pingidentity.adapters.httpbasic.idp.HttpBasicIdpAuthnAdapter"
  }

  configuration {
    tables {
      name = "Credential Validators"
      rows {
        fields {
          name  = "Password Credential Validator Instance"
          value = pingfederate_password_credential_validator.demo.name
        }
      }
    }
    fields {
      name  = "Realm"
      value = "foo"
    }

    fields {
      name  = "Challenge Retries"
      value = "3"
    }
  }

  attribute_contract {
    core_attributes {
      name      = "username"
      pseudonym = true
    }
    extended_attributes {
      name = "sub"
    }
  }
  attribute_mapping {
    attribute_contract_fulfillment {
      key_name = "sub"
      source {
        type = "ADAPTER"
      }
      value = "sub"
    }
    attribute_contract_fulfillment {
      key_name = "username"
      source {
        type = "ADAPTER"
      }
      value = "username"
    }
    jdbc_attribute_source {
      filter      = "\"\""
      description = "foo"
      schema      = "INFORMATION_SCHEMA"
      table       = "ADMINISTRABLE_ROLE_AUTHORIZATIONS"
      data_store_ref {
        id = "ProvisionerDS"
      }
    }
  }
}

resource "pingfederate_sp_adapter" "demo" {
  name          = "bar"
  sp_adapter_id = "bar"
  plugin_descriptor_ref {
    id = "com.pingidentity.adapters.opentoken.SpAuthnAdapter"
  }

  configuration {

    #pf version specific config
    dynamic "fields" {
      for_each = local.isPF10_0 ? [] : [1]
      content {
        name  = "SameSite Cookie"
        value = "3"
      }
    }

    sensitive_fields {
      name  = "Password"
      value = "Secret123"
    }
    sensitive_fields {
      name  = "Confirm Password"
      value = "Secret123"
    }
    fields {
      name  = "Account Link Service"
      value = ""
    }
    fields {
      name  = "Authentication Service"
      value = ""
    }
    fields {
      name  = "Cipher Suite"
      value = "2"
    }
    fields {
      name  = "Cookie Domain"
      value = ""
    }
    fields {
      name  = "Cookie Path"
      value = "/"
    }
    fields {
      name  = "Force SunJCE Provider"
      value = "false"
    }
    fields {
      name  = "HTTP Only Flag"
      value = "true"
    }
    fields {
      name  = "Logout Service"
      value = ""
    }
    fields {
      name  = "Not Before Tolerance"
      value = "0"
    }
    fields {
      name  = "Obfuscate Password"
      value = "true"
    }
    fields {
      name  = "Secure Cookie"
      value = "false"
    }
    fields {
      name  = "Send Extended Attributes"
      value = ""
    }
    fields {
      name  = "Send Subject as Query Parameter"
      value = "false"
    }
    fields {
      name  = "Session Cookie"
      value = "false"
    }
    fields {
      name  = "Session Lifetime"
      value = "43200"
    }
    fields {
      name  = "Skip Trimming of Trailing Backslashes"
      value = "false"
    }
    fields {
      name  = "Subject Query Parameter                 "
      value = ""
    }
    fields {
      name  = "Token Lifetime"
      value = "300"
    }
    fields {
      name  = "Token Name"
      value = "opentoken"
    }
    fields {
      name  = "Transport Mode"
      value = "2"
    }
    fields {
      name  = "URL Encode Cookie Values"
      value = "true"
    }
    fields {
      name  = "Use Verbose Error Messages"
      value = "false"
    }
  }

  target_application_info {
    application_name     = "foo"
    application_icon_url = "https://foo"
  }

  depends_on = [
    pingfederate_server_settings.settings
  ]
}

resource "pingfederate_oauth_access_token_mappings" "demo" {
  access_token_manager_ref {
    id = pingfederate_oauth_access_token_manager.my_atm.id
  }

  context {
    type = "CLIENT_CREDENTIALS"
  }
  attribute_contract_fulfillment {
    key_name = "sub"
    source {
      type = "CONTEXT"
    }
    value = "ClientId"
  }
}

resource "pingfederate_oauth_openid_connect_policy" "demo" {
  policy_id = "foo"
  name      = "foo"
  access_token_manager_ref {
    id = pingfederate_oauth_access_token_manager.my_atm.id
  }
  attribute_contract {
    extended_attributes {
      name                 = "email"
      include_in_user_info = true
    }
    extended_attributes {
      name = "email_verified"
    }
    extended_attributes {
      name = "family_name"
    }
    extended_attributes {
      name                = "name"
      include_in_id_token = true
    }
  }
  attribute_mapping {
    attribute_contract_fulfillment {
      key_name = "sub"
      source {
        type = "NO_MAPPING"
      }
    }
    attribute_contract_fulfillment {
      key_name = "email"
      source {
        type = "NO_MAPPING"
      }
    }
    attribute_contract_fulfillment {
      key_name = "email_verified"
      source {
        type = "NO_MAPPING"
      }
    }
    attribute_contract_fulfillment {
      key_name = "family_name"
      source {
        type = "NO_MAPPING"
      }
    }
    attribute_contract_fulfillment {
      key_name = "name"
      source {
        type = "NO_MAPPING"
      }
    }
  }

  //  scope_attribute_mappings = { //TODO hoping the new TF 2.0.0 SDK will finally support sensible maps
  //    address = ["foo", "bar"]
  //  }

  depends_on = [
    pingfederate_server_settings.settings
  ]
}

resource "pingfederate_notification_publisher" "demo" {
  name         = "bar"
  publisher_id = "foo1"
  plugin_descriptor_ref {
    id = "com.pingidentity.email.SmtpNotificationPlugin"
  }

  configuration {
    #pf version specific config
    dynamic "fields" {
      for_each = local.isPF10_0 ? [] : [1]
      content {
        name  = "UTF-8 Message Header Support"
        value = "false"
      }
    }

    fields {
      name  = "From Address"
      value = "help@foo.org"
    }
    fields {
      name  = "Email Server"
      value = "foo"
    }
    fields {
      name  = "SMTP Port"
      value = "25"
    }
    fields {
      name  = "Encryption Method"
      value = "NONE"
    }
    fields {
      name  = "SMTPS Port"
      value = "465"
    }
    fields {
      name  = "Verify Hostname"
      value = "true"
    }
    fields {
      name  = "Username"
      value = ""
    }
    fields {
      name  = "Password"
      value = ""
    }
    fields {
      name  = "Test Address"
      value = ""
    }
    fields {
      name  = "Connection Timeout"
      value = "30"
    }
    fields {
      name  = "Retry Attempt"
      value = "2"
    }
    fields {
      name  = "Retry Delay"
      value = "2"
    }
    fields {
      name  = "Enable SMTP Debugging Messages"
      value = "false"
    }
  }
}

resource "pingfederate_authentication_selector" "demo" {
  name = "one"
  plugin_descriptor_ref {
    id = "com.pingidentity.pf.selectors.cidr.CIDRAdapterSelector"
  }

  configuration {
    fields {
      name  = "Result Attribute Name"
      value = ""
    }
    tables {
      name = "Networks"
      rows {
        fields {
          name  = "Network Range (CIDR notation)"
          value = "127.0.0.1/32"
        }
      }
    }
  }
}

resource "pingfederate_authentication_policies" "demo" {
  fail_if_no_selection    = false
  tracked_http_parameters = ["foo"]
  authn_selection_trees {
    name = "one"
    root_node {
      action {
        type = "AUTHN_SELECTOR"
        authentication_selector_ref {
          id = pingfederate_authentication_selector.demo.id
        }
      }
      children {
        action {
          type    = "AUTHN_SELECTOR"
          context = "No"
          authentication_selector_ref {
            id = pingfederate_authentication_selector.demo.id
          }
        }
        children {
          action {
            type    = "CONTINUE"
            context = "No"
          }
        }
        children {
          action {
            type    = "CONTINUE"
            context = "Yes"
          }
        }
      }
      children {
        action {
          type    = "CONTINUE"
          context = "Yes"
        }
      }
    }
  }
  authn_selection_trees {
    name = "foo"
    root_node {
      action {
        type = "AUTHN_SOURCE"
        authentication_source {
          type = "IDP_ADAPTER"
          source_ref {
            id = pingfederate_idp_adapter.demo.id
          }
        }
      }
      children {
        action {
          type    = "RESTART"
          context = "Fail"
        }
      }
      children {
        action {
          type    = "DONE"
          context = "Success"
        }
      }
    }
  }
  dynamic "authn_selection_trees" {
    for_each = local.isPF10_2 ? [1] : []
    content {
      enabled = true
      name    = "frag"
      root_node {
        action {
          type = "FRAGMENT"
          fragment {
            id = pingfederate_authentication_policy_fragment.demo[0].id
          }
          fragment_mapping {
            attribute_contract_fulfillment {
              key_name = "one"

              source {
                type = "NO_MAPPING"
              }
            }
            attribute_contract_fulfillment {
              key_name = "subject"

              source {
                type = "NO_MAPPING"
              }
            }
            attribute_contract_fulfillment {
              key_name = "two"

              source {
                type = "NO_MAPPING"
              }
            }
          }
        }

        children {
          action {
            context = "Fail"
            type    = "DONE"
          }
        }
        children {
          action {
            context = "Success"
            type    = "DONE"
          }
        }
      }
    }
  }
}

resource "pingfederate_certificates_ca" "demo" {
  count          = 20
  certificate_id = "example${count.index}"
  file_data      = base64encode(file("certificate_ca/cacert${count.index}.pem"))
}

resource "pingfederate_idp_sp_connection" "demo" {
  name         = "acc_test_foo"
  entity_id    = "foo"
  active       = true
  logging_mode = "STANDARD"
  credentials {
    certs {
      x509_file {
        file_data = file("amazon_root_ca1.pem")
      }
    }
    inbound_back_channel_auth {
      type                    = "INBOUND"
      digital_signature       = false
      require_ssl             = false
      verification_subject_dn = "cn=bar"
    }
  }
  attribute_query {
    jdbc_attribute_source {
      filter      = "*"
      description = "foo"
      schema      = "INFORMATION_SCHEMA"
      table       = "ADMINISTRABLE_ROLE_AUTHORIZATIONS"
      id          = "foo"
      data_store_ref {
        id = "ProvisionerDS"
      }
    }

    attribute_contract_fulfillment {
      key_name = "foo"
      source {
        type = "JDBC_DATA_STORE"
        id   = "foo"
      }
      value = "GRANTEE"
    }

    attributes = ["foo"]
    policy {
      sign_response                  = false
      sign_assertion                 = false
      encrypt_assertion              = false
      require_signed_attribute_query = false
      require_encrypted_name_id      = false
    }
  }
}

resource "pingfederate_keypair_signing" "demo" {
  file_data = filebase64("identity.p12")
  password  = "changeit"
}

resource "pingfederate_custom_data_store" "example" {
  name = "example"
  plugin_descriptor_ref {
    id = "com.pingidentity.pf.datastore.other.RestDataSourceDriver"
  }
  configuration {
    tables {
      name = "Base URLs and Tags"
      rows {
        fields {
          name  = "Base URL"
          value = "https://example.com"
        }
        fields {
          name  = "Tags"
          value = ""
        }
        default_row = true
      }
    }
    tables {
      name = "HTTP Request Headers"
    }
    tables {
      name = "Attributes"
      rows {
        fields {
          name  = "Local Attribute"
          value = "foo"
        }
        fields {
          name  = "JSON Response Attribute Path"
          value = "/bar"
        }
      }
    }
    fields {
      name  = "Authentication Method"
      value = "None"
    }
    fields {
      name  = "Username"
      value = ""
    }
    fields {
      name  = "Password"
      value = ""
    }
    fields {
      name  = "OAuth Token Endpoint"
      value = ""
    }
    fields {
      name  = "OAuth Scope"
      value = ""
    }
    fields {
      name  = "Client ID"
      value = ""
    }
    fields {
      name  = "Client Secret"
      value = ""
    }
    fields {
      name  = "Enable HTTPS Hostname Verification"
      value = "true"
    }
    fields {
      name  = "Read Timeout (ms)"
      value = "10000"
    }
    fields {
      name  = "Connection Timeout (ms)"
      value = "10000"
    }
    fields {
      name  = "Max Payload Size (KB)"
      value = "1024"
    }
    fields {
      name  = "Retry Request"
      value = "true"
    }
    fields {
      name  = "Maximum Retries Limit"
      value = "5"
    }
    fields {
      name  = "Retry Error Codes"
      value = "429"
    }
    fields {
      name = "Test Connection URL"
    }
  }
}

resource "pingfederate_authentication_policy_fragment" "demo" {
  count       = local.isPF10_2 ? 1 : 0
  name        = "fragtest"
  description = "functional test"
  inputs {
    id = pingfederate_authentication_policy_contract.input.id
  }
  outputs {
    id = pingfederate_authentication_policy_contract.output.id
  }

  root_node {
    action {
      type = "AUTHN_SELECTOR"
      authentication_selector_ref {
        id = pingfederate_authentication_selector.demo.id
      }
    }
    children {
      action {
        type    = "DONE"
        context = "No"
      }
    }
    children {
      action {
        type    = "DONE"
        context = "Yes"
      }
    }
  }
}

resource "pingfederate_authentication_policy_contract" "input" {
  name                = "fragmenttest1"
  extended_attributes = ["one", "two"]
}

resource "pingfederate_authentication_policy_contract" "output" {
  name                = "fragmenttest2"
  extended_attributes = ["three", "four"]
}

resource "pingfederate_keypairs_oauth_openid_connect" "settings" {
  static_jwks_enabled = true
  rsa_active_cert_ref {
    id = pingfederate_keypair_signing.rsa[0].id
  }
  rsa_previous_cert_ref {
    id = pingfederate_keypair_signing.rsa[1].id
  }
  rsa_decryption_active_cert_ref {
    id = pingfederate_keypair_signing.rsa[0].id
  }
  rsa_decryption_previous_cert_ref {
    id = pingfederate_keypair_signing.rsa[1].id
  }

  p256_active_cert_ref {
    id = pingfederate_keypair_signing.ec256[0].id
  }
  p256_previous_cert_ref {
    id = pingfederate_keypair_signing.ec256[1].id
  }
  p256_decryption_active_cert_ref {
    id = pingfederate_keypair_signing.ec256[0].id
  }
  p256_decryption_previous_cert_ref {
    id = pingfederate_keypair_signing.ec256[1].id
  }

  p384_active_cert_ref {
    id = pingfederate_keypair_signing.ec384[0].id
  }
  p384_previous_cert_ref {
    id = pingfederate_keypair_signing.ec384[1].id
  }
  p384_decryption_active_cert_ref {
    id = pingfederate_keypair_signing.ec384[0].id
  }
  p384_decryption_previous_cert_ref {
    id = pingfederate_keypair_signing.ec384[1].id
  }

  p521_active_cert_ref {
    id = pingfederate_keypair_signing.ec521[0].id
  }
  p521_previous_cert_ref {
    id = pingfederate_keypair_signing.ec521[1].id
  }
  p521_decryption_active_cert_ref {
    id = pingfederate_keypair_signing.ec521[0].id
  }
  p521_decryption_previous_cert_ref {
    id = pingfederate_keypair_signing.ec521[1].id
  }
  rsa_publish_x5c_parameter             = true
  p256_publish_x5c_parameter            = true
  p384_publish_x5c_parameter            = true
  p521_publish_x5c_parameter            = true
  p256_decryption_publish_x5c_parameter = true
  p384_decryption_publish_x5c_parameter = true
  p521_decryption_publish_x5c_parameter = true
  rsa_decryption_publish_x5c_parameter  = true
}

resource "pingfederate_keypair_signing" "rsa" {
  count                     = 2
  city                      = "Test"
  common_name               = "example${count.index}"
  country                   = "GB"
  key_algorithm             = "RSA"
  key_size                  = 2048
  organization              = "Test"
  organization_unit         = "Test"
  state                     = "Test"
  valid_days                = 365
  subject_alternative_names = ["foo", "bar"]
}

resource "pingfederate_keypair_signing" "ec256" {
  count                     = 2
  common_name               = "example${count.index}"
  country                   = "GB"
  key_algorithm             = "EC"
  key_size                  = 256
  organization              = "Test"
  valid_days                = 365
  subject_alternative_names = ["foo", "bar"]
}

resource "pingfederate_keypair_signing" "ec384" {
  count                     = 2
  common_name               = "example${count.index}"
  country                   = "GB"
  key_algorithm             = "EC"
  key_size                  = 384
  organization              = "Test"
  valid_days                = 365
  subject_alternative_names = ["foo", "bar"]
}

resource "pingfederate_keypair_signing" "ec521" {
  count                     = 2
  common_name               = "example${count.index}"
  country                   = "GB"
  key_algorithm             = "EC"
  key_size                  = 521
  organization              = "Test"
  valid_days                = 365
  subject_alternative_names = ["foo", "bar"]
}

resource "pingfederate_authentication_policy_contract" "foo" {
  name                = "issue117"
  extended_attributes = ["realName", "role", "mail"]
}

resource "pingfederate_idp_sp_connection" "foo_connection" {
  name         = "issue117"
  entity_id    = "issue117"
  active       = true
  base_url     = "https://some.url"
  logging_mode = "STANDARD"
  contact_info {}
  credentials {
    signing_settings {
      signing_key_pair_ref {
        id = pingfederate_keypair_signing.demo.id
      }
      include_cert_in_signature    = false
      include_raw_key_in_signature = false
      algorithm                    = "SHA256withRSA"
    }
    certs {
      active_verification_cert    = true
      encryption_cert             = false
      primary_verification_cert   = true
      secondary_verification_cert = false
      x509_file {
        file_data = file("amazon_root_ca1.pem")
      }
    }
  }
  sp_browser_sso {
    protocol          = "SAML20"
    incoming_bindings = ["REDIRECT"]
    enabled_profiles  = ["SP_INITIATED_SSO"]
    sso_service_endpoints {
      binding    = "POST"
      url        = "/saml/acs"
      is_default = true
      index      = 0
    }
    sign_assertions               = true
    sign_response_as_required     = true
    sp_saml_identity_mapping      = "STANDARD"
    require_signed_authn_requests = true
    assertion_lifetime {
      minutes_before = 5
      minutes_after  = 5
    }
    encryption_policy {
      encrypt_assertion             = false
      encrypt_slo_subject_name_id   = false
      slo_subject_name_id_encrypted = false
      encrypted_attributes          = []
    }
    attribute_contract {
      core_attributes {
        name        = "SAML_SUBJECT"
        name_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
      }
      extended_attributes {
        name        = "role"
        name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"
      }
      extended_attributes {
        name        = "realName"
        name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"
      }
      extended_attributes {
        name        = "mail"
        name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"
      }
    }
    authentication_policy_contract_assertion_mappings {
      attribute_contract_fulfillment {
        key_name = "SAML_SUBJECT"
        source {
          type = "AUTHENTICATION_POLICY_CONTRACT"
        }
        value = "subject"
      }
      attribute_contract_fulfillment {
        key_name = "mail"
        source {
          type = "AUTHENTICATION_POLICY_CONTRACT"
        }
        value = "mail"
      }
      attribute_contract_fulfillment {
        key_name = "realName"
        source {
          type = "AUTHENTICATION_POLICY_CONTRACT"
        }
        value = "realName"
      }
      attribute_contract_fulfillment {
        key_name = "role"
        source {
          type = "AUTHENTICATION_POLICY_CONTRACT"
        }
        value = "role"
      }
      authentication_policy_contract_ref {
        id = pingfederate_authentication_policy_contract.foo.id
      }
      restrict_virtual_entity_ids        = false
      restricted_virtual_entity_ids      = []
      abort_sso_transaction_as_fail_safe = false
    }
  }
}

resource "pingfederate_idp_sp_connection" "outbound_provisioning" {
  name         = "outbound_provision"
  entity_id    = "outbound_provision"
  active       = true
  base_url     = "https://some.url"
  logging_mode = "STANDARD"
  contact_info {}
  credentials {
    signing_settings {
      signing_key_pair_ref {
        id = pingfederate_keypair_signing.demo.id
      }
      include_cert_in_signature    = false
      include_raw_key_in_signature = false
      algorithm                    = "SHA256withRSA"
    }
    certs {
      active_verification_cert    = true
      encryption_cert             = false
      primary_verification_cert   = true
      secondary_verification_cert = false
      x509_file {
        file_data = file("amazon_root_ca1.pem")
      }
    }
  }
  sp_browser_sso {
    protocol          = "SAML20"
    incoming_bindings = ["REDIRECT"]
    enabled_profiles  = ["SP_INITIATED_SSO"]
    sso_service_endpoints {
      binding    = "POST"
      url        = "/saml/acs"
      is_default = true
      index      = 0
    }
    sign_assertions               = true
    sign_response_as_required     = true
    sp_saml_identity_mapping      = "STANDARD"
    require_signed_authn_requests = true
    assertion_lifetime {
      minutes_before = 5
      minutes_after  = 5
    }
    encryption_policy {
      encrypt_assertion             = false
      encrypt_slo_subject_name_id   = false
      slo_subject_name_id_encrypted = false
      encrypted_attributes          = []
    }
    attribute_contract {
      core_attributes {
        name        = "SAML_SUBJECT"
        name_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
      }
      extended_attributes {
        name        = "role"
        name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"
      }
      extended_attributes {
        name        = "realName"
        name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"
      }
      extended_attributes {
        name        = "mail"
        name_format = "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"
      }
    }
    authentication_policy_contract_assertion_mappings {
      attribute_contract_fulfillment {
        key_name = "SAML_SUBJECT"
        source {
          type = "AUTHENTICATION_POLICY_CONTRACT"
        }
        value = "subject"
      }
      attribute_contract_fulfillment {
        key_name = "mail"
        source {
          type = "AUTHENTICATION_POLICY_CONTRACT"
        }
        value = "mail"
      }
      attribute_contract_fulfillment {
        key_name = "realName"
        source {
          type = "AUTHENTICATION_POLICY_CONTRACT"
        }
        value = "realName"
      }
      attribute_contract_fulfillment {
        key_name = "role"
        source {
          type = "AUTHENTICATION_POLICY_CONTRACT"
        }
        value = "role"
      }
      authentication_policy_contract_ref {
        id = pingfederate_authentication_policy_contract.foo.id
      }
      restrict_virtual_entity_ids        = false
      restricted_virtual_entity_ids      = []
      abort_sso_transaction_as_fail_safe = false
    }
  }
  outbound_provision {
    type = "PingIDForWorkforce"

    channels {
      active      = false
      max_threads = 1
      name        = "bar"
      timeout     = 60

      attribute_mapping {
        field_name = "userName"

        saas_field_info {
          attribute_names = []
          character_case  = "NONE"
          create_only     = false
          default_value   = "asdasd"
          masked          = false
          parser          = "NONE"
          trim            = false
        }
      }
      attribute_mapping {
        field_name = "email"

        saas_field_info {
          attribute_names = []
          character_case  = "NONE"
          create_only     = false
          masked          = false
          parser          = "NONE"
          trim            = false
        }
      }
      attribute_mapping {
        field_name = "fName"

        saas_field_info {
          attribute_names = []
          character_case  = "NONE"
          create_only     = false
          masked          = false
          parser          = "NONE"
          trim            = false
        }
      }
      attribute_mapping {
        field_name = "lName"

        saas_field_info {
          attribute_names = []
          character_case  = "NONE"
          create_only     = false
          masked          = false
          parser          = "NONE"
          trim            = false
        }
      }

      channel_source {
        base_dn             = "cn=bar"
        guid_attribute_name = "entryUUID"
        guid_binary         = false

        account_management_settings {
          account_status_algorithm      = "ACCOUNT_STATUS_ALGORITHM_FLAG"
          account_status_attribute_name = "nsaccountlock"
          default_status                = true
          flag_comparison_status        = false
          flag_comparison_value         = "true"
        }

        change_detection_settings {
          changed_users_algorithm   = "TIMESTAMP_NO_NEGATION"
          group_object_class        = "groupOfUniqueNames"
          time_stamp_attribute_name = "modifyTimestamp"
          user_object_class         = "person"
        }

        data_source {
          id = pingfederate_ldap_data_store.demo_ldap.id
        }

        group_membership_detection {
          group_member_attribute_name = "uniqueMember"
        }

        group_source_location {
          nested_search = false
        }

        user_source_location {
          group_dn      = "cn=bar"
          nested_search = false
        }
      }
    }

    sensitive_target_settings {
      inherited = false
      name      = "base64Key"
      value     = "secret"
    }
    sensitive_target_settings {
      inherited = false
      name      = "token"
      value     = "secret"
    }

    target_settings {
      inherited = false
      name      = "Provisioning Options"
    }
    target_settings {
      inherited = false
      name      = "disableNewUsers"
      value     = "true"
    }
    target_settings {
      inherited = false
      name      = "domain"
      value     = "idpxnyl3m.pingidentity.eu"
    }
    target_settings {
      inherited = false
      name      = "orgAlias"
      value     = "foo"
    }
    target_settings {
      inherited = false
      name      = "removeAction"
      value     = "Disable"
    }
    target_settings {
      inherited = false
      name      = "updateNewUsers"
      value     = "true"
    }
  }

}
