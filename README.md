# capture-events
it is meant to capture the events of any website.

1. It has a websocket which captures any raw data and pushes it into the elasticsearch.
2. The fasthttp server based websocket has been used to accomodate high scalability.
3. The data dumped into elasticsearch can be used for analysis.
