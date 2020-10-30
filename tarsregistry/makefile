TARGET := tarsregistry

ALL: tars2go
	cd protocol && ../tars2go --outdir=../autogen *.tars
	go build -o tarsregistry
	
tars2go:
	go build -o tars2go github.com/TarsCloud/TarsGo/tars/tools/tars2go
	
clean:
	rm -rf $(TARGET) 