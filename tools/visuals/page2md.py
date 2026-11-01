"""Convert TRAIL_visuals.html into README.md so the page can be read on GitHub.

Usage: python3 -I page2md.py <TRAIL_visuals.html> <gif-dir> <README.md>

Each <figure class="stage"> becomes an image of its exported GIF. Text, lists
and tables are carried over; in-page links are rewritten to GitHub heading
anchors. Fails if a GIF is missing, so README and visuals/ stay in step.
"""
import os
import re
import sys
from html.parser import HTMLParser

VOID = {"meta", "link", "br", "img", "input", "hr", "source"}
VERDICT_KINDS = {"ok", "limit", "ext"}


def gif_name(i, scene):
    """Same rule as names.js: keyTerms -> 01-key-terms.gif, scMonitor -> 14-monitor.gif."""
    s = re.sub(r"^sc(?=[A-Z])", "", scene)
    return f"{i + 1:02d}-" + re.sub(r"(?<!^)([A-Z])", r"-\1", s).lower() + ".gif"


class Node:
    def __init__(self, tag, attrs, parent):
        self.tag, self.a, self.parent, self.kids = tag, dict(attrs), parent, []


class Parser(HTMLParser):
    def __init__(self):
        super().__init__(convert_charrefs=True)
        self.root = Node("root", {}, None)
        self.cur = self.root

    def handle_starttag(self, tag, attrs):
        n = Node(tag, attrs, self.cur)
        self.cur.kids.append(n)
        if tag not in VOID:
            self.cur = n

    def handle_endtag(self, tag):
        n = self.cur
        while n is not self.root and n.tag != tag:
            n = n.parent
        if n is not self.root:
            self.cur = n.parent

    def handle_data(self, d):
        self.cur.kids.append(d)


def find(n, pred):
    if isinstance(n, str):
        return
    if pred(n):
        yield n
    for k in n.kids:
        yield from find(k, pred)


def elems(n):
    return [k for k in n.kids if not isinstance(k, str)]


def cls(n):
    return n.a.get("class") or ""


def is_verdict(n):
    return not isinstance(n, str) and cls(n).split()[:1] == ["verdict"]


def slug(t):
    """GitHub's heading anchor rule."""
    return re.sub(r"[^\w\- ]", "", t.lower().strip()).replace(" ", "-")


def htext(n):
    """Heading text without its verdict badge."""
    return "".join(k if isinstance(k, str) else htext(k) for k in n.kids if not is_verdict(k)).strip()


def main():
    src_path, gif_dir, out_path = sys.argv[1:4]
    p = Parser()
    p.feed(open(src_path, encoding="utf-8").read())
    root = next(find(p.root, lambda n: n.tag == "main"))
    out_dir = os.path.dirname(os.path.abspath(out_path))

    anchors = {}
    for h in find(root, lambda n: n.tag in ("h2", "h3") and n.a.get("id")):
        anchors[h.a["id"]] = slug(htext(h))
    for s in find(root, lambda n: n.tag == "section" and n.a.get("id")):
        h2 = next(find(s, lambda n: n.tag == "h2"), None)
        if h2:
            anchors[s.a["id"]] = slug(htext(h2))

    figs = list(find(root, lambda n: n.tag == "figure" and "stage" in cls(n)))
    gifs = {f.a["id"]: gif_name(i, f.a["data-scene"]) for i, f in enumerate(figs)}
    missing = [g for g in gifs.values() if not os.path.exists(os.path.join(out_dir, gif_dir, g))]
    if missing:
        sys.exit("missing GIFs (run tools/visuals/build.sh): " + ", ".join(missing))

    def inline(n):
        if isinstance(n, str):
            return re.sub(r"\s+", " ", n)
        if is_verdict(n):
            return ""
        inner = "".join(inline(k) for k in n.kids)
        if n.tag in ("strong", "b"):
            return f"**{inner.strip()}**"
        if n.tag in ("em", "i"):
            return f"*{inner.strip()}*"
        if n.tag == "code" or "mono" in cls(n).split():
            return f"`{inner.strip()}`"
        if n.tag == "a":
            href = n.a.get("href", "")
            if href.startswith("#"):
                href = "#" + anchors.get(href[1:], href[1:])
            return f"[{inner.strip()}]({href})"
        return inner

    out = []
    pending_verdict = [False]  # next <p> after a scenario heading is its verdict

    def block(n):
        if isinstance(n, str):
            return
        t, c = n.tag, cls(n)
        if t == "h1":
            out.append("# " + n.a.get("aria-label", htext(n)))
            return
        if t in ("h2", "h3"):
            out.append(("## " if t == "h2" else "### ") + htext(n))
            v = next((k for k in elems(n) if is_verdict(k)), None)
            pending_verdict[0] = False
            if v:
                words = cls(v).split()[1:]
                if words and not set(words) <= VERDICT_KINDS:
                    # Summary sentence was placed in the class attribute; the
                    # next <p> holds the verdict label.
                    out.append(" ".join(words))
                    pending_verdict[0] = True
                else:
                    # Well-formed badge: <span class="verdict ok">Works</span>
                    label = "".join(inline(k) for k in v.kids).strip()
                    if label:
                        out.append(f"**Verdict:** {label}")
            return
        if t == "p":
            txt = inline(n).strip()
            if not txt:
                return
            if "where" in c.split():
                txt = f"*{txt}*"
            elif pending_verdict[0]:
                txt = f"**Verdict:** {txt}"
            pending_verdict[0] = False
            out.append(txt)
            return
        if t == "figure" and "stage" in c:
            pending_verdict[0] = False
            out.append(f"![{n.a.get('aria-label', '')}]({gif_dir}/{gifs[n.a['id']]})")
            return
        if t == "ul":
            items = []
            for li in (k for k in elems(n) if k.tag == "li"):
                h = next((k for k in elems(li) if k.tag in ("h3", "h4")), None)
                if h:
                    rest = " ".join(inline(k).strip() for k in elems(li) if k is not h)
                    items.append(f"- **{htext(h)}**: {rest}")
                else:
                    items.append("- " + inline(li).strip())
            out.append("\n".join(items))
            return
        if t == "table":
            rows = [[inline(td).strip().replace("|", "\\|") for td in elems(tr) if td.tag in ("td", "th")]
                    for tr in find(n, lambda k: k.tag == "tr")]
            lines = ["| " + " | ".join(rows[0]) + " |", "|" + "---|" * len(rows[0])]
            lines += ["| " + " | ".join(r) + " |" for r in rows[1:]]
            out.append("\n".join(lines))
            return
        if t in ("script", "style", "canvas"):
            return
        for k in n.kids:
            block(k)

    block(root)
    out.insert(1, "> Generated from [`TRAIL_visuals.html`](TRAIL_visuals.html) by "
                  "`tools/visuals/build.sh`; do not edit by hand. The animations play slower than on "
                  "the interactive page. The draft itself is "
                  "[`draft-nassar-agentproto-trail.md`](draft-nassar-agentproto-trail.md).")
    with open(out_path, "w", encoding="utf-8") as f:
        f.write("\n\n".join(out) + "\n")
    print(f"wrote {out_path}: {len(figs)} animations")


if __name__ == "__main__":
    main()
