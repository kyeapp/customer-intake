windows:
	env GOOS=windows GOARCH=amd64 go build server.go

serve:
	light-server -s . -p 8080 \
		-w dist/tailwind.css \
		-w index.html

install-server-env:
	sudo apt install make
	wget "https://go.dev/dl/$$(curl 'https://go.dev/VERSION?m=text').linux-amd64.tar.gz"
	sudo rm -rf /usr/local/go
	sudo tar --directory /usr/local -xzf go*.linux-amd64.tar.gz
	sudo ln -s /usr/local/go/bin/go /usr/local/bin/go
	sudo ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt

live:
	sudo pkill server || true
	sudo -b nohup go run server.go

# install cert here:
# https://certbot.eff.org/
renew-cert:
	sudo certbot certonly --webroot

tailwind:
	npx tailwindcss -i ./src/tailwind_source.css -o ./dist/tailwind.css --watch

serve-nav2:
	light-server -s . -p 8080 -w ./nav2.html -w ./dist/style.css
