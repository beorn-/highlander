.PHONY: test bin docker debug
default:
	$(MAKE) all
bin:
	./scripts/build.sh
all:
	$(MAKE) bin

clean:
	rm build/*
