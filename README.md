# Events

A tool for abstracting CloudWatch events so that they can be used productively as business events.

To run the tool use:

```
go run ./cli/*.go --event=someEvent --namespace=com.example
```

## Use cases

There are a number of goals for this tool that will help it become a productive and compelling offering for those using CloudWatch / EventBridge:

- [x] Query all subscribers to a specific event (including other AWS accounts)
- [ ] Query the publisher(s) of a specific event
- [ ] Validate events as they are in transit on the bus, publish metrics about event validation
- [ ] Re-raise events in a new namespace to indicate they've passed validation
- [ ] Query the relationship between events & functions
- [ ] Group events & functions by tag to create contextual diagrams
- [ ] Query built-in AWS events and see the relationship between these events and business events
- [ ] Query publisher / subscriber relationships in multi-regional setups

## Requirements

This tool makes the following assumptions about your CloudWatch events:

1. You use the `Source` field to denote the "namespace" under which the event is raised (eg. `com.mybusiness`).
2. You use the `DetailType` field to denote the event name.
3. You use the `Detail` field to store the event payload.
4. You add a lambda ARN to the `Resources` field of the event to indicate which function raised the event.
5. All your events are run over the `default` event bus (TODO: Remove this requirement)
