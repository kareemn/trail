// GIF file name for the i-th animation (0-based, page order) with a given data-scene.
// keyTerms -> 01-key-terms.gif, scMonitor -> 14-monitor.gif.
// page2md.py has the same rule (gif_name); keep them in step.
const gifName = (i, scene) =>
  String(i + 1).padStart(2, "0") + "-" +
  scene.replace(/^sc(?=[A-Z])/, "").replace(/(?<!^)([A-Z])/g, "-$1").toLowerCase() + ".gif";

module.exports = { gifName };
