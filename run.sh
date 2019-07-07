gnome-terminal -- go run server.go
gnome-terminal -- go run server_critical_region.go

# for i in 0 1 2
# do
# 	declare -i port=8090+$i
# 	gnome-terminal -- go run client.go localhost:$port
# done