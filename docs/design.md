# Portcullis Design Doc

Portcullis is a transparent pseudo-broker that sets up security groups for
other brokers on bind calls. In order to acheive this, it acts as a
man-in the-middle between the CF Cloud Controller and the target service broker,
parsing out bind calls, and forwarding all information to the service broker
backend. To Cloud Foundry, Portcullis appears to be the target service broker,
and to the service broker, Portcullis appears to be Cloud Foundry.

Portcullis also has the added benefit of providing a permissions handle, such
that users configured by an admin can create service brokers in all orgs/spaces
without giving those users the other permissions associated with being an
admin.

Portcullis itself will need to be deployed with admin credentials to Cloud
Foundry in order to create the service brokers and security groups.

## Creating Service Brokers

### For Service Broker Admins

The service broker admin can log into Portcullis with their UAA credentials to
authenticate with Portcullis. Once authenticated, the admin can use a bearer
token to interact with HTTPS endpoints that allow making service broker
associations with Portcullis.

The creation endpoint needs to link Portcullis to the service broker API
address, create a Portcullis endpoint that routes to that service broker, and
then register the broker with Cloud Foundry using its API.

While the general use case may be to register service brokers with all
spaces, there may need to be a way implemented to create service brokers mapped to
specific spaces, or with no space bindings at all.

There will also need to be an endpoint to delete or edit these mappings.

Authentication for these admins will need to be configured with the UAA by
somebody with UAA credentials.

### For Cloud Foundry Users

Service brokers are typically created by space developers using the CF CLI
targeted against the service broker backend. Using Portcullis, the process is
largely the same, except for that the creation call needs to target a
Portcullis endpoint created by a Portcullis admin
(e.g. `https://portcullis-addr.io/redis`).

Once a service broker is created, service instance creation and binding
should continue as they normally would without Portcullis, as far as the user
is concerned.

## Security Group Creation

Portcullis is passing along all service broker calls to the service broker API,
but when it receives a bind request, it makes a security group creation
request to Cloud Foundry, opening up the address/port that the application
will contact the service on.

This information is gathered by parsing the
`credentials` JSON handed back from the service broker in response to the
bind request. While many service brokers follow the convention of including a
`port` key and host `key` or `uri` key at the top level of this JSON, there is
no enforcement of this, and there are definitely brokers which deviate in a way
that this heuristic will not work. So while Portcullis will look for this
information by default, a configuration of what to parse in the broker to get
the address and port information can be provided to Portcullis during broker
creation or by edit endpoint afterward. Also, there should be a way to specify
to open a range of ports, or not open any egress at all for this broker when it
is bound.

## Storage

User credentials will be stored in the UAA, and if you point Portcullis toward
the same UAA that your Cloud Foundry or BOSH use, you can log in with those same
credentials if you assign the Portcullis-designated permission group to that
user. To be clear, Portcullis will not do any handling of the user credentials
itself.

Service broker administrative credentials will be held by the Cloud Controller
database the way they always have been. The credentials just get passed through
in service broker calls by the Cloud Controller in order to authenticate to the
service broker API.

Service Broker names to API locations need to be able to withstand redeployments,
outages, etc. Therefore, the mappings could theoretically be stored in a
Postgres database (such as the on the Postgres node used for Cloud Foundry itself),
but a more lightweight solution involves just storing these mappings in a flat
file in a parsable format (e.g. YAML). The stored mappings are loaded in on
Portcullis's launch, and more can be added to the file as more configurations are
made. This file would also contain custom rules for pulling out port and address
information from a broker that were configured from the Portcullis API. This
may also need to be the location where mappings between created security groups
to service binding are stored.

## Transitioning to Portcullis

As it stands, your Cloud Foundry probably has service brokers mapped to the
Cloud Controller the "normal way". While the most direct path to using
Portcullis could involve deleting your current service broker (and, by extension,
all existing service instances and bindings), this may lead to a loss of service
data if your service broker discards an application's data when a service
instance is deprovisioned. This could reasonably be seen as unacceptable.

An alternative strategy is to create new service broker mappings with Cloud
Foundry that are configured through Portcullis, and to deprecate the old broker
mappings using `cf disable-service-access`. This way, all new service instances
would be made through Portcullis. The downside to this is that some "manual"
(see: not done by Portcullis) cleanup of previously existing security-groups
needs to be made as the transition to Portcullis is made over time. Portcullis
will never be responible for deleting security-groups that it did not itself
create - that is, Portcullis _opens_ traffic, and the CF admin is responsible
for keeping the default restrictions of Cloud Foundry as closed as they wish.

By far, beginning to use Portcullis after having used Cloud Foundry without it
is the biggest obstacle of this design.

## Interface

Communication with Portcullis, at least at first, will be made by HTTP API
calls. In the future, a CLI or web frontend may be created, but the focus at
first will be functionality.

## Certificates

As with every other SSL/TLS-using thing, you'll need a certificate to use it.
Portcullis will need to support this, as credentials will be sent over the wire,
and so you will probably want to be able to support providing certificates as well.