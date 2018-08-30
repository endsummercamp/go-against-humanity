.PHONY: build-card ita-original

build-card:
	mkdir -p cards/json-against-humanity/src/${DECK}
	python cards/italian/parser.py cards/italian/source/${DECK}/main.tex cards/json-against-humanity/src/$1

ita-original: 
	@$(MAKE) build-card DECK=ita-original
