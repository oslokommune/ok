# app

Boilerplate for an application running on ECS.

## Example

In `vars-common.yml`:

```yml
AccountId: "12345679876"
Region: "eu-west-1"
Team: "pirates"
Environment: "pirates-dev"
```

In `vars-app-rain.yml`:

```yml
StackName: "app-rain"
AppName: "rain"
WithExampleImage: true
WithDailyShutdown: false
WithAlbHostRouting: true
WithOpenTelemetrySidecar: false
WithVPCEndpoints: false
```

Render the template:

```sh
boilerplate \
  --template-url . \
  --var-file vars-common.yml \
  --var-file vars-app-rain.yml \
  --output-folder app-rain \
  --non-interactive
```
