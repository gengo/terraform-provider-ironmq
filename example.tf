provider "ironmq" {
}

resource "ironmq_queue" "queue" {
    name = "terraform-example"
}
