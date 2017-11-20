rm *.log
rm *.hlt
go build mybot.go && ./halite -d "240 160" "./old/19" "./mybot" -s 10644969
