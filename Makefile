.PHONY: bin all clean

default:
	$(MAKE) all
bin:
	./scripts/build.sh
all:
	$(MAKE) bin

clean:
	rm build/*
