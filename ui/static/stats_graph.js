export default class StatsGraph {
  constructor() {
    this._el = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    this._el.classList.add("stats-graph");
  }

  get element() {
    return this._el;
  }

  set stats(stats) {
    this._stats = stats;
    this._renderGraph();
  }

  _renderGraph() {
    const stats = this._stats;
    const count = stats.length;
    if (count < 2) return;

    const width = this._el.clientWidth;
    const height = this._el.clientHeight;
    this._el.setAttribute("viewBox", `0 0 ${width} ${height}`);

    const margin = { top: 10, left: 15, right: 10, bottom: 15 };
    const innerWidth = width - margin.left - margin.right;
    const innerHeight = height - margin.top - margin.bottom;

    const startX = stats[0].ts;
    const endX = stats[count - 1].ts;
    const scaleX = innerWidth / (endX - startX);
    const X = ts => (ts - startX) * scaleX;

    const maxY = Math.max(...stats.map(({ interval }) => interval)) + 1;
    const scaleY = innerHeight / Math.ceil(maxY);
    const Y = interval => (maxY - interval) * scaleY;

    const maxF = Math.max(...stats.map(({ factor }) => factor)) + 1;
    const scaleF = innerHeight / Math.ceil(maxF);
    const fY = factor => (maxF - factor) * scaleF;

    const path = this._stats.map(({ ts, interval }, idx) =>
      idx === 0 ? `M${X(ts)},${Y(interval)}` : `L${X(ts)},${Y(interval)}`
    );
    const fPath = this._stats.map(({ ts, factor }, idx) =>
      idx === 0 ? `M${X(ts)},${fY(factor)}` : `L${X(ts)},${fY(factor)}`
    );

    const dots = this._stats.map(
      ({ ts, interval }) => `<circle cx="${X(ts)}" cy="${Y(interval)}" r="5"/>`
    );
    const fDots = this._stats.map(
      ({ ts, factor }) => `<circle cx="${X(ts)}" cy="${fY(factor)}" r="5"/>`
    );

    const values = this._stats.map(
      ({ ts, interval }) =>
        `<text x="${X(ts)}" y="${Y(interval) - 10}">${label(interval)}</text>`
    );
    const fValues = this._stats.map(
      ({ ts, factor }) =>
        `<text x="${X(ts)}" y="${fY(factor) - 10}">${factor.toFixed(2)}</text>`
    );

    this._el.innerHTML = `
<g transform="translate(${margin.left},${margin.top})">
  <g class="axis">
    <line x1="0" y1="${innerHeight}" x2="${innerWidth}" y2="${innerHeight}"></line>
    <text x="0" y="${innerHeight + 15}">${new Date(
      startX * 1000
    ).toLocaleDateString()}</text>
    <text x="${innerWidth}" y="${innerHeight + 15}">${new Date(
      endX * 1000
    ).toLocaleDateString()}</text>
  </g>

  <g class="graph factor">
    <path d="${fPath.join("")}"/>
    ${fDots}
    ${fValues}
  </g>

  <g class="graph interval">
    <path d="${path.join("")}"/>
    ${dots}
    ${values}
  </g>
</g>`;
  }
}

function label(interval) {
  const hours = (interval * 24).toFixed();
  const div = hours / 24;
  const mod = hours % 24;
  return (
    (hours >= 24 ? `${Math.floor(div)}d ` : "") + (mod > 0 ? `${mod}h` : "")
  );
}
