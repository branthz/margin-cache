PROGS = marginCache
BUILD_VERBOSE := -v

TEST_VERBOSE := -v

all: $(PROGS)

.PHONY: $(PROGS) 
$(PROGS): main.go
	go build -o $(PROGS) -ldflags "-w -s" 

clean:
	go clean
	rm -rf $(PROGS)
