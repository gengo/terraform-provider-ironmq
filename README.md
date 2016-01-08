# terraform-provider-ironmq

terraform-provider-ironmq is a plugin of [Terraform](https://terraform.io/).
It allows you to manage queues in [IronMQ](http://www.iron.io/mq/) with Terraform.

## Status
Experimental.

The resource configuration schema can change without notice nor migration support.

## Installation
```sh
go get github.com/gengo/terraform-provider-ironmq
```

Then, copy `$GOPATH/bin/terraform-provider-ironmq` into the same directory as `terraform` binary.
 
## Example Usage
At first, put a valid [client configuration of IronMQ](http://dev.iron.io/mq/3/libraries/) into a right place.
Then you can run `terraform` command as usual.

example.tf:
```
provider "ironmq" {
}

# Pull queues
resource "ironmq_queue" "pull-example1" {
    name = "pull-example1"
}

resource "ironmq_queue" "pull-example2" {
    name = "pull-example2"
}

# Push queue
resource "ironmq_queue" "push-example" {
    name = "push-example"
    type = "multicast"
    push {
        subscribers {
                url = "ironmq:///pull-example1"
        }
        subscribers {
                url = "ironmq:///pull-example2"
        }
    }
}
```

## Argument Reference

* `env` - (Optional) environment name defined in the IronMQ client configuration. 

## Resource Reference
This plugin provides a resource type named `ironmq_queue`.
It supports the following arguments:

* `name` - (Required, string) The name of the queue.
* `type` - (Optional, string) The type of the queue. Must be either `pull`, `unicast` or `multicast`. `pull` is the default value.
* `push` - (Optional, object) See below. It is required unless `type` is `pull`.

The `push` argument supports the following nested arguments.

* `subscribers` - (Requried, List of objects) See below.

The `push.subscribers` argument supports the following nested arguments.

* `url` - (Required, string) The URL of the subscriber endpoint.


## License
terraform-provider-ironmq is licensed under Mozilla Public License, version 2.0.
See `LICENSE.txt` for more details.
