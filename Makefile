#.PHONY: build-card ita-original compile

#build-card:
#	mkdir -p cards/json-against-humanity/src/${DECK}
#	python cards/italian/parser.py cards/italian/source/${DECK}/main.tex cards/json-against-humanity/src/${DECK}
#	cp cards/italian/source/${DECK}/metadata.json cards/json-against-humanity/src/${DECK}/
#
#hack-it: 
#	@$(MAKE) build-card DECK=hack-it
#
#ita-original: 
#	@$(MAKE) build-card DECK=ita-original
#
#ita-original-sfoltita: 
#	@$(MAKE) build-card DECK=ita-original-sfoltita
#
#ita-espansione: 
#	@$(MAKE) build-card DECK=ita-espansione
#all: hack-it ita-original ita-original-sfoltita ita-espansione

.PHONY: compile
compile:
	cd cards/json-against-humanity; python compile.py
