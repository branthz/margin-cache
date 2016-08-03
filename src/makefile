PROGS = marginCache
all: $(PROGS)
$(PROGS):*.go
	go build -ldflags "-o marginCache -w -s" *.go
clean:
	rm -rf $(PROGS)
