# portcullis

The goal behind Portcullis is to make it easier for CF service broker 
admins and CF admins to handle running services in Cloud Foundry. 
The observation here is that in large Cloud Foundry environments, the
people managing services are not typically the same people with the
adminstrative rights to Cloud Foundry - and that makes sense. However,
this communication can be the bottleneck in getting services up and
running. Portcullis seeks to wrap the permissions system of the UAA/Cloud
Foundry to provide a handle for which service broker devs/admins can manage
the access of their own broker. Also, security groups will be opened to
users as bindings are made to service instances.

## How Do I Run It?
The dependencies are vendored, so cloning the repo and either `go run`ing
the main.go file or `go build`ing the project will start the program. A
configuration file is required - look in assets/examples for an idea of how
to make one (docs at a later time), and set the `PORTCULLIS_CONFIG`
environment variable to the path at which you have created the configuration.

This project is still in early development, so the master branch may contain 
buggy code until a preliminary set of features has been completed.
