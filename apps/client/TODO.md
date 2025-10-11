# TODO

- add id manager pending messages notifier loop (remove ad-hoc notification
when message fails)
- add get random user endpoint service for the id manager conn
- add pending message sender loop
  - it should prioritize own (immune) messages, the message cache
  interface should allow for querying them

## DONE

- add message receiver loop
  - it will handle messages for oneself and for others
- messages are to be sent in secure envelopes

## CANCELLED

- add message cachers storage for the middleware __(clients will not ask for pending
messages, they will attempt to give pending messages)__
- add pending message querier (for cachers updates and for allowance to delete
non-immune pending messages) __(same as above)__
