# AGENTS.md

Instructions for coding agents (Claude Code, Codex, and others) working in this
repo. `CLAUDE.md` imports this file.

This repo holds the TRAIL Internet-Draft and an illustrated explainer of it.
**The draft, the explainer page, the README, and the GIFs must agree.** A change
to one is not done until the others match.

## What is here

| Path | What it is | Edit it? |
|---|---|---|
| `draft-nassar-agentproto-trail.md` | The draft, in kramdown-rfc Markdown. **Source of truth for the protocol.** | Yes |
| `draft-nassar-agentproto-trail.html` | The draft rendered by xml2rfc | No: rebuild it |
| `TRAIL_visuals.html` | Self-contained explainer page: text plus 21 canvas animations. **Source of truth for the pictures.** | Yes |
| `README.md` | `TRAIL_visuals.html` converted to Markdown so it renders on GitHub | No: run `tools/visuals/build.sh` |
| `visuals/NN-name.gif` | Each animation exported as a GIF | No: run `tools/visuals/build.sh` |
| `tools/visuals/` | The export and check scripts | Yes |
| `Makefile` | Builds the draft with i-d-template from `../moq-transport/lib` | Rarely |

The draft says what TRAIL is. The page explains it. When they disagree, the
draft wins, unless the change is something the authors are deliberately
proposing on the page first. Those proposals are listed in the page's "Where the
page and the draft differ" table.

## Workflow: after changing the draft

1. **Edit `draft-nassar-agentproto-trail.md`.**
2. **Rebuild the draft HTML.** Use `make` if the i-d-template lib is at
   `../moq-transport/lib`. Otherwise:
   ```sh
   kramdown-rfc draft-nassar-agentproto-trail.md > draft-nassar-agentproto-trail.xml
   xml2rfc --html draft-nassar-agentproto-trail.xml -o draft-nassar-agentproto-trail.html
   rm -rf draft-nassar-agentproto-trail.xml .refcache
   ```
3. **Check section citations:**
   `python3 -I tools/visuals/check_refs.py draft-nassar-agentproto-trail.md TRAIL_visuals.html`
   - `MISSING` means the page cites a section number that no longer exists. Fix it.
   - `check` means the number exists, but its title shares no words with what
     the page calls it. Read that section and confirm the citation is still right.
   - Adding, removing, or reordering a section renumbers everything after it.
     Always rerun this check after doing so.
4. **Update the page wherever it describes what you changed.** Search
   `TRAIL_visuals.html` for the term, directive, field, or behavior. Check all of
   these:
   - the paragraph above each animation;
   - the `In the draft: name (N.N); ...` line under each animation;
   - the glossary table ("The draft's terms, in plain words");
   - the scenario `Uses:` and `Missing:` notes;
   - the "What TRAIL doesn't cover" list;
   - the "Where the page and the draft differ" table: update the Status column,
     and remove rows the draft now covers;
   - the animation code: labels, `term()` text, `captions()`, and `missing()` boxes.
     If the draft now defines something an animation drew as a dashed red
     "not in draft yet" box, redraw it as a normal element.
5. **Regenerate the GIFs and README** for every animation whose drawing changed:
   ```sh
   tools/visuals/build.sh                  # all 21 (about 2 minutes)
   tools/visuals/build.sh records,tamper   # only these data-scene names
   ```
   Text-only page changes outside the canvases need only the README. Running
   `build.sh` with a single scene is the fast way to get it.
6. **Review before committing:** `git status` should show only the GIFs you
   meant to change. Read the `README.md` diff. Look at one frame of each changed
   GIF; ffmpeg can extract a PNG.
7. **Commit the source change and the generated files together.**

The same steps apply in reverse after changing the page: if the page now says
something the draft does not, either change the draft or add a row to the
differences table.

## Writing style: the draft

- **Follow kramdown-rfc conventions:**
  - anchors are `{#kebab-case}`;
  - references are `{{anchor}}` and `{{RFC9110}}`;
  - every table ends with `{: #tab-name title="..."}`;
  - terms are defined as definition lists (`Term:` on one line, `: definition` on the next);
  - figures use `~~~ ascii-art` and `~~~ cddl`;
  - prose wraps at 72 columns.
- **Use BCP 14 keywords (MUST, SHOULD, MAY) in uppercase only for normative
  requirements.** Explanatory text uses ordinary lowercase words.
- **Write in plain, short, declarative sentences:** one idea per sentence,
  present tense. Say what the edge, authority, or witness does. Avoid hedges
  ("generally", "typically") and selling ("powerful", "seamless").
- **Don't use em dashes.** Neither the draft nor the page uses any. Use a
  colon, a semicolon, parentheses, or a new sentence.
- **Wrap angle-bracket placeholders in backticks**, for example
  `` `(mcp, <session-id>, ...)` ``. Bare `<session-id>` is parsed as an HTML tag
  and breaks the xml2rfc build.
- **Use straight ASCII quotes** in the draft.
- **Use the draft's terms exactly as defined in Terminology (§1.6):**
  - lowercase: sandbox, cell, edge, trail, checkpoint, evidence plane, state,
    directive, directed record, grant, authority, witness, admission issuer,
    admission token, verifier;
  - capitalized: Action Record, Attestation Result;
  - attestation levels in backticks: `software`, `tee`, `isolated`;
  - directives lowercase and hyphenated: policy, hold, freshness-budget,
    require-gate, require-witness, require-level, require-admission,
    require-record.

  A new term goes into Terminology and into the page's glossary table.
- **"MOQT" in the draft, "MoQ" on the page.** The draft names the protocol
  precisely; the page uses the name readers know. Don't mix them within one file.
- **List the authors in this order:** Kareem Nassar, Suhas Nandakumar, Cullen
  Jennings. They appear in three places: the draft front matter, the page byline
  (`.byline`), and the page footer.
- **Keep the draft name consistent:**
  - docname `draft-nassar-agentproto-trail-latest`;
  - file name `draft-nassar-agentproto-trail.md`;
  - `DRAFT` in `Makefile`.

  Rename all three together.

## Writing style: the page and README

- **Use the same plain-language rules as the draft, with no BCP 14 keywords.**
  Explain from the reader's side ("Edge B refuses it"), not the spec's ("the
  receiving edge MUST reject").
- **Call it "the draft"**, not "the spec" or "the RFC".
- **Use typographic apostrophes (’) in page text.** The canvas code uses
  straight quotes inside JS strings.
- **Name each party by role plus one capital letter, and reuse the existing
  letters:**
  - Agent A/B, Edge A/B, Machine A/B, Cell A/B, Operator X/Y;
  - Regulator R, Notary N, Verifier V, Monitor M, Security team T, Bank K,
    Officer O, Safety institute S, Token issuer Q.

  A new party gets a new, unused letter.
- **Name records by edge letter and sequence number:** A1, A2, B1, B2. A record
  that answers another "cites" it ("B1 cites A1").
- **Structure each building block like the others:** an `h3.scene` heading, one
  paragraph, then the `<figure class="stage">`, then
  `<p class="where">In the draft: name (N.N); name (N.N)</p>`.
- **Number each scenario**, with a verdict badge, a verdict line (Works /
  Works, with limits / Needs one addition / ...), the figure, and a
  `.draftcheck` block with `Uses:` (links to building blocks) and `Missing:`
  (with the open issue it maps to).
- **Build scenarios only from building blocks.** If a scenario needs a
  mechanism that has no block, add the block first.

## Animation rules (`TRAIL_visuals.html` script)

- **A scene is `{ dur, draw(c, t) }`.**
  - `draw` must be a pure function of `t`: no state carried between frames, no
    `Math.random()`, no `Date`. The GIF export steps the slider to set `t`, so
    anything else exports wrong.
  - Register new scenes in the `scenes` map.
  - Add a `<figure class="stage" data-scene="name" aria-label="Animation: ...">`
    with the same controls markup as the others.
  - Keep `dur` between 14 and 22 seconds.
- **Use the 900×600 canvas with the standard layout.** It is always dark
  (`BG #14171C`), whatever the page theme. Coordinates are in units of `U = 75`
  pixels from the center, with `y` pointing up. Operators go at the top, the
  evidence plane in the middle, and the parties that read it (authorities,
  witnesses, verifiers) at the bottom.
- **Each color means one thing, in every scene and in the legend.** Use the
  constants; never reuse a color for a different meaning:

  | Constant | Means |
  |---|---|
  | `AGENT` #7AA2F7 | agent |
  | `EDGE` #FF9E64 | edge |
  | `CELL` #9ECE6A | cell |
  | `PLANE` #7DCFFF | evidence plane |
  | `AUTH` #BB9AF7 | authority |
  | `WIT` #73DACA | witness, attestation level badge |
  | `REC` #E0C46C | Action Record; draft terms (yellow mono) |
  | `TOOL` #A9B1D6 | service outside TRAIL |
  | `OK` / `BAD` / `WARN` | allowed / refused or broken / degraded |
  | `HOST` #5A6370, `SANDBOX` #4E6A9E | machine, sandbox |

  A new kind of thing needs a new color, a constant, and a legend entry in "How
  to read the pictures".
- **Show draft terms in yellow mono** with `term()` or `termAfter()`, under or
  beside the example name.
- **Mark parts the draft doesn't define** with `missing()`, a dashed red box
  labeled "not in draft yet". Each one must also appear in that scenario's
  `Missing:` note and in the draft's Open Issues.
- **Captions** (`captions(c, t, [[t0, "text", color?]])`):
  - one sentence, at most about 90 characters, no trailing period;
  - `BAD` color for the failure or attack beat;
  - time the beats to the drawing, then check the caption is readable in the
    exported GIF.
- **Reuse the helpers** (`box`, `party`, `cell`, `edgeBox`, `trail`, `tag`,
  `move`, `travel`, `rest`, `seal`, `xmark`, `check`) instead of drawing raw
  shapes. Pictures stay consistent that way.
- **Keep the page self-contained.** The only external resource is the IBM Plex
  Sans/Mono stylesheet from Google Fonts. Page colors are CSS tokens on `:root`,
  with dark overrides under both `prefers-color-scheme: dark` and
  `[data-theme="dark"]`. Keep both in step.

## GIF and README export

- **Keep the GIF settings the authors chose:** 1.8× slower than the page, so
  captions can be read (the page felt too fast); 10 fps; the last frame held
  3 s before looping; 900 px wide; about 1 MB each. These are the defaults in
  `tools/visuals/build.sh`. Don't change them without asking.
- **Don't rename the GIFs by hand.** Names are `NN-scene.gif`, where NN is the
  figure's position on the page. Inserting a figure renumbers the ones after it,
  and a full `build.sh` run cleans up the old names.
- **Requirements:** Node 18+, Python 3, ffmpeg, and Google Chrome. Playwright
  drives the system Chrome, so it downloads no browser. `build.sh` runs
  `npm ci` on first use.
- **Keep the README GitHub-only.** The repo is private, so the README is how
  members view the page. Don't publish the page or GIFs anywhere public: GitHub
  Pages from this repo would be public.

## Repo hygiene

- **Keep these out of git:** `node_modules/`, `.refcache/`, intermediate `.xml`,
  extracted frames, and `.history` (see `.gitignore`).
- **Write commit messages as a short imperative summary**, with a body when the
  reason isn't obvious.

## Known issues (remove each line once fixed)

- **Stale citations:** the page cites §1.7 for terminology, which is now §1.6.
  It also cites Appendix D.1–D.5 for scenario extensions, which are now one
  bullet in Open Issues (Appendix D). `check_refs.py` lists them all.
- **Malformed verdict badges:** each scenario's badge is written as
  `<span class="verdict SUMMARY">kind</span>`, so the page never shows the
  summary and the badge loses its color. Intended:
  `<span class="verdict kind">label</span>` followed by the summary paragraph.
  `page2md.py` handles both forms.
