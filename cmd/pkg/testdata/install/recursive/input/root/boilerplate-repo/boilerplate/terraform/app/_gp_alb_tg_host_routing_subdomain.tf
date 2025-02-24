{{- if dig "Subdomain" "Enable" false .AlbHostRouting -}}
{{ template "doNotEdit" . }}
module "alb_tg_host_routing_subdomain" {

  # https://github.com/oslokommune/golden-path-iac/tree/main/terraform/modules/alb-tg-host-routing
  source = "git@github.com:oslokommune/golden-path-iac//terraform/modules/alb-tg-host-routing-subdomain?ref=alb-tg-host-routing-subdomain-v0.8.0"

  service_name = local.main_container_name

  alb_https_listener_arn = data.aws_lb_listener.this.arn
  alb_listener_priority  = null # Set to null to automatically assign a priority

  target_group_name = "${local.environment}-${local.main_container_name}"
  target_group_port = local.main_container_port

  # https://docs.aws.amazon.com/elasticloadbalancing/latest/application/target-group-health-checks.html#health-check-settings
  target_group_health_check = {
    enabled             = true
    protocol            = "HTTP"
    path                = local.main_container_health_check_path
    port                = "traffic-port"
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    interval            = 30
    matcher             = "200"
  }

  route53_zone_name    = data.aws_route53_zone.this.name
  route53_record_name  = local.main_container_name
  route53_record_value = "" # Makes the record value the same as the zone name

{{- if dig "Subdomain" "TargetGroupTargetStickiness" false .AlbHostRouting }}

  target_group_target_stickiness = {
    enabled = true
    type    = "lb_cookie"
    cookie_duration = 86400
  }
{{- end }}

}

moved {
  from = module.alb_tg_host_routing
  to   = module.alb_tg_host_routing_subdomain
}
{{ else -}}
# x-boilerplate-delete
{{ end -}}
