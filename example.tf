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
