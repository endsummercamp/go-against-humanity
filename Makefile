.PHONY: build-card ita-original

build-card:
	mkdir -p cards/json-against-humanity/src/${DECK}
	python cards/italian/parser.py cards/italian/source/${DECK}/main.tex cards/json-against-humanity/src/${DECK}

ita-original: 
	@$(MAKE) build-card DECK=ita-original

ita-original-sfoltita: 
	@$(MAKE) build-card DECK=ita-original-sfoltita

ita-espansione: 
	@$(MAKE) build-card DECK=ita-espansione

all: ita-original ita-original-sfoltita ita-espansione
