provider "ironmq" {
}

resource "ironmq_queue" "queue" {
    name = "terraform-example"
    type = "multicast"
    push {
        subscribers {
                url = "ironmq:///yugui_test_subscriber1"
        }
        subscribers {
                url = "ironmq:///yugui_test_subscriber2"
        }
    }
}
