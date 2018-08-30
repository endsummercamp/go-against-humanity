.PHONY: ita

ita:
	mkdir -p cards/json-against-humanity/src/ita-original
	python cards/italian/parser.py cards/italian/source/cah-ita-originale-federico.tex cards/json-against-humanity/src/ita-original
