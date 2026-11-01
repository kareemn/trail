"""Check that section numbers cited in TRAIL_visuals.html exist in the draft.

Usage: python3 -I check_refs.py <draft.md> <TRAIL_visuals.html>

Numbers sections the way kramdown-rfc/xml2rfc will: '#'/'##'/'###' headings
in the middle are 1, 1.1, 1.1.1; References is inserted after the last middle
section; back-matter headings are appendices A, B, ... Headings followed by
{:numbered="false"} are skipped.

Prints every citation with the draft's current title for that number.
Exits 1 if a cited number does not exist. Lines marked "check" cite a number
whose title shares no words with the citation; read those by hand.
"""
import html
import re
import sys

STOP = {"the", "and", "for", "with", "what", "how", "its", "per", "a", "an", "of", "to", "in", "on", "as"}


def sections(md):
    lines = md.split("\n")
    nums, part, fence = {}, None, False
    count = [0, 0, 0]
    app = [0, 0, 0]
    for i, line in enumerate(lines):
        if line.startswith("~~~") or line.startswith("```"):
            fence = not fence
            continue
        if fence:
            continue
        if line.startswith("--- "):
            if part == "middle":
                count[0] += 1
                nums[str(count[0])] = "References"
            part = line[4:].strip()
            continue
        m = re.match(r"^(#{1,3}) (.+?)\s*(\{#[^}]*\})?\s*$", line)
        if not m or part not in ("middle", "back"):
            continue
        nxt = lines[i + 1] if i + 1 < len(lines) else ""
        if "numbered=\"false\"" in nxt:
            continue
        lvl, title = len(m.group(1)) - 1, m.group(2)
        c = count if part == "middle" else app
        c[lvl] += 1
        for j in range(lvl + 1, 3):
            c[j] = 0
        head = str(c[0]) if part == "middle" else chr(ord("A") + c[0] - 1)
        nums[".".join([head] + [str(x) for x in c[1:lvl + 1]])] = title
    return nums


def words(s):
    return {w.rstrip("s") for w in re.findall(r"[a-z]+", s.lower()) if w not in STOP and len(w) > 2}


def main():
    md_path, page_path = sys.argv[1:3]
    nums = sections(open(md_path, encoding="utf-8").read())
    page = open(page_path, encoding="utf-8").read()
    text = html.unescape(re.sub(r"<[^>]+>", "", page.split("<script")[0]))

    cites = []  # (cited-as, number, context)
    for line in text.split("\n"):
        line = line.strip()
        if line.startswith("In the draft:"):
            for part in line[len("In the draft:"):].split(";"):
                m = re.match(r"\s*(.+?)\s*\(([^)]*)\)", part)
                if m:
                    for n in re.split(r"[,\s]+", m.group(2).strip()):
                        if n:
                            cites.append((m.group(1), n, line))
        for m in re.finditer(r"§\s?([A-Z]?\d+(?:\.\d+)*)", line):
            cites.append((line, m.group(1), line))
        for m in re.finditer(r"Appendix ([A-Z](?:\.\d+)*)", line):
            cites.append((line, m.group(1), line))

    bad = 0
    for name, n, _ in cites:
        title = nums.get(n)
        if title is None:
            bad += 1
            print(f"MISSING  {n:6} cited as: {name[:90]}")
        elif not (words(name) & words(title)) and not name.startswith(n):
            print(f"check    {n:6} draft: {title!r:45} cited as: {name[:60]}")
        else:
            print(f"ok       {n:6} draft: {title!r:45}")
    print(f"\n{len(cites)} citations, {bad} missing")
    sys.exit(1 if bad else 0)


if __name__ == "__main__":
    main()
