PROGS = fusionCache
all: $(PROGS)
$(PROGS):*.go
	go build -ldflags "-o fusionCache -w -s" *.go
clean:
	rm -rf $(PROGS)
