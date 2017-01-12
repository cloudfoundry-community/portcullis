# Portcullis Release 0.2.0

Another early cut of Portcullis. Not even close to 1.0 yet, but
good progress is being made.

## For users

* Formal command line support was added to the server.
* There is now a /v1/info endpoint in the API.
* Fixed a bug where specifying an improper auth type
  would result in an unintuitive crash outside of the
  initialization phase.
* log_level can now be set in the config.

## For developers

* Pumped out some more tests in the API and config suites.
* Created a concourse pipeline for the project, located
  in the CI folder.