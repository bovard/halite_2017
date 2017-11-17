zip -r -X Archive.zip MyBot.go src/
go build MyBot.go && hlt bot -b Archive.zip
