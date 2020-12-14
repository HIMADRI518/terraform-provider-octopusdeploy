resource "octopusdeploy_project" "example" {
  allow_deployments_to_no_targets      = false
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  deployment_process_id                = "deploymentprocess-Projects-123"
  description                          = "The development project."
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = "Lifecycles-123"
  name                                 = "Development Project (OK to Delete)"
  project_group_id                     = "ProjectGroups-123"
  slug                                 = "development-project"
  tenanted_deployment_participation    = "TenantedOrUntenanted"
  variable_set_id                      = "variableset-Projects-123"

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }

  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.NextPatch}"
  }
}