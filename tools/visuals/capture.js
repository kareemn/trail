// Steps every canvas animation in TRAIL_visuals.html through its slider and
// writes PNG frames, one directory per animation, named after its GIF.
//
// Usage: node capture.js <TRAIL_visuals.html> <outdir> <fps> <slowdown> [scene,scene...]
//   fps       playback rate of the GIF
//   slowdown  how much slower than the page (1.8 = a 10 s scene lasts 18 s)
//   scenes    optional data-scene names to export; default is all of them
//
// Frames are taken at slider positions, so the result does not depend on
// machine speed. Needs Google Chrome installed (uses channel "chrome").
const { chromium } = require("playwright-core");
const fs = require("fs");
const path = require("path");
const { gifName } = require("./names");

const [, , htmlPath, outDir, fpsArg, slowArg, onlyArg] = process.argv;
if (!htmlPath || !outDir || !fpsArg || !slowArg) {
  console.error("usage: node capture.js <html> <outdir> <fps> <slowdown> [scenes]");
  process.exit(2);
}
const FPS = Number(fpsArg), SLOW = Number(slowArg);
const ONLY = onlyArg ? new Set(onlyArg.split(",")) : null;

(async () => {
  const src = fs.readFileSync(htmlPath, "utf8");
  // scene durations come from the source: `const <name> = ... dur: N`
  const names = src.match(/const scenes = \{([^}]*)\}/)[1].split(",").map(s => s.trim()).filter(Boolean);
  const durs = {};
  for (const n of names) {
    const m = src.match(new RegExp("const " + n + " = [\\s\\S]*?dur: (\\d+(?:\\.\\d+)?)"));
    if (!m) throw new Error("no dur for scene " + n);
    durs[n] = Number(m[1]);
  }

  const browser = await chromium.launch({ channel: "chrome", headless: true });
  // 2x pixels, downscaled when encoding, for smoother text.
  // Reduced motion stops the page from autoplaying while we step it.
  const page = await browser.newPage({ viewport: { width: 1200, height: 900 }, deviceScaleFactor: 2, reducedMotion: "reduce" });
  await page.goto("file://" + path.resolve(htmlPath));
  await page.evaluate(() => document.fonts.ready);
  await page.waitForTimeout(500);

  const figs = await page.evaluate(() =>
    [...document.querySelectorAll("figure.stage")].map(f => ({ id: f.id, scene: f.dataset.scene })));

  fs.mkdirSync(outDir, { recursive: true });
  for (let i = 0; i < figs.length; i++) {
    const f = figs[i];
    if (ONLY && !ONLY.has(f.scene)) continue;
    const dur = durs[f.scene];
    if (!dur) throw new Error("no dur for " + f.scene);
    const n = Math.round(dur * SLOW * FPS);
    const dir = path.join(outDir, gifName(i, f.scene).replace(/\.gif$/, ""));
    fs.mkdirSync(dir, { recursive: true });
    for (let j = 0; j <= n; j++) {
      const v = Math.round((j / n) * 1000);
      const data = await page.evaluate(([id, v]) => {
        const fig = document.getElementById(id);
        const r = fig.querySelector("input");
        r.value = String(v);
        r.dispatchEvent(new Event("input"));
        return fig.querySelector("canvas").toDataURL("image/png");
      }, [f.id, v]);
      fs.writeFileSync(path.join(dir, String(j).padStart(4, "0") + ".png"), Buffer.from(data.split(",")[1], "base64"));
    }
    console.log(path.basename(dir), `${dur}s scene, ${n + 1} frames`);
  }
  await browser.close();
})().catch(e => { console.error(e); process.exit(1); });
