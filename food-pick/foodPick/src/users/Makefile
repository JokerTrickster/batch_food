.PHONY: build clean

build:
	ls
    GOOS=linux GOARCH=amd64 go build -o .aws-sam/build/GetUserFunction/bootstrap ./main.go

# SAM에서 호출하는 빌드 타겟
build-GetUserFunction:
	GOOS=linux GOARCH=amd64 go build -o .aws-sam/build/GetUserFunction/bootstrap ./main.go


clean:
	rm -rf ../.aws-sam/build/GetUserFunction/bootstrap