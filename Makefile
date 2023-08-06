build:
	@CGO_ENABLED=0 GOOS=linux go build -o flistify

prepare: 
	@bash scripts/prepare.sh