rm *.log
rm *.hlt
go build mybot.go && ./halite -d "240 160" "./mybot" "./old/20"
