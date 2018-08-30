import re, sys

def usage():
    print("Usage:\n%s [input file] [output file]" % (sys.argv[0]))

def getcard(card_snippet):
    cards = []
    matcher = re.compile("\\\carta\{(.*?)\}", re.S | re.I)
    m = re.findall(matcher, card_snippet)
    for match in m:
        text = match.replace('\n', ' ').replace('\puntini', '_').replace('{\small', '').strip()
        text = text.replace('\%', '%')
        text = text.replace('``', '\'')
        text = text.replace("''", "'")
        cards.append(text)
    return cards

if len(sys.argv) != 3:
    usage()
    exit()

with open(sys.argv[1], 'r') as texfile:
    document = texfile.read()

matcher_pages = re.compile(
    r'\\begin\{longtable\}\{.*?\}(.*?)\\end\{longtable\}.*?\\begin\{longtable\}\{.*?\}(.*?)\\end{longtable}', re.S | re.I)

m = re.search(matcher_pages, document)


whites = getcard(m[1])
blacks = getcard(m[2])

with open(sys.argv[2] + "/white.md.txt", "w") as f:
    for l in whites:
        f.write(l + "\n")  

with open(sys.argv[2] + "/black.md.txt", "w") as f:
    for l in blacks:
        f.write(l + "\n")
