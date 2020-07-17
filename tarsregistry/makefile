TARGET := tarsregistry

ALL:
	tars2go --outdir=autogen protocol/tarsregistry.tars
	go build -o tarsregistry
	
clean:
	rm -rf $(TARGET) 