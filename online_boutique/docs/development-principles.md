# Development Principles

> **Note:** This document outlines guidances behind some development decisions
> behind the Online Boutique demo application.

### Minimal configuration

Running the demo locally or on GCP should require minimal to no
configuration unless absolutely necessary to run critical parts of the demo.

Configuration that takes multiple steps, especially such as creating service
accounts should be avoided.

### Microservice implementations should not be complex

Each service should provide a minimal implementation and try to avoid
unnecessary code and logic that's not executed.

Keep in mind that any service implementation is a decent example of “a GRPC
application that runs on Kubernetes”. Keeping the source code short and
navigable will serve this purpose.

It is okay to have intentional inefficiencies in the code as they help
illustrate the capabilities of profiling and diagnostics offerings.
