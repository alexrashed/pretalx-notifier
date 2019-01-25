# Pretalx Notifier
Simple tool to check for new submissions on a predefined instance of [PreTalx](https://pretalx.com/). If a new submission is detected a notification on [PushOver](https://pushover.net/) is sent.
Environment variables to set:
- PRETALX_URL: Url of the PreTalx instance, pointing to the API of the event. F.e. https://pretalx.example.com/api/event/exampleEvent/
- PRETALX_AUTH: Auth token to authenticate to PreTalx
- PUSHOVER_API_TOKEN; API token to authenticate to pushover
- PUSHOVER_USER_TOKEN: User token of the pushover notification receiver
- MINUTES: Optional, number of minutes between checks