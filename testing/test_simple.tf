resource "etcd_keys" "ami" {
    key {
        name = "ami"
        path = "service/app/launch_ami"
        value = "bleh"
    }
}
