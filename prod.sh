CGO_ENABLED=0 go build -o JSE .

git rm -r app/ core/ internal/ main.go prod.sh tmp/ utils/
