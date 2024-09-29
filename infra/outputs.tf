output "game_server" {
    value = {
        instance_id = module.game_server.instance_id
        public_ip = module.game_server.public_ip
    }
}