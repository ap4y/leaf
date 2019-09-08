import StatsGraph from "./stats_graph.js";

const template = `
<header>
  <h3>Deck: <span id="stats-deck"></span></h3>
  <div class="stats-select-row">
    <h4>Cards: </h4>
    <select id="stats-list"></select>
  </div>
</header>

<main class="container">
  <h3>Current Stats for <span id="stats-card"></span></h3>

  <ul>
    <li>
      <strong>Last Reviewed At: </strong>
      <span id="reviewed-at"></span>
    </li>
    <li>
      <strong>Interval: </strong>
      <span id="interval"></span>
    </li>
    <li>
      <strong>Difficulty: </strong>
      <span id="difficulty"></span>
    </li>
  </ul>
</main>
`;

export default class StatsList {
  constructor() {
    this._el = document.createElement("div");
    this._el.innerHTML = template;
    this._graph = new StatsGraph();
    this._el.querySelector("main").appendChild(this._graph.element);
  }

  get element() {
    return this._el;
  }

  set deck(deck) {
    this._el.querySelector("#stats-deck").innerHTML = deck;
  }

  set stats(stats) {
    this._stats = stats;
    this._renderItems();
  }

  _renderItems() {
    const stats = this._stats;
    if (stats.length === 0) return;

    this._populateSelect(stats);
    this._renderStats(stats[0]);

    this._el.querySelector("#stats-list").onchange = ({ target }) => {
      const stat = stats.find(({ card }) => card === target.value);
      this._renderStats(stat);
    };
  }

  _populateSelect(stats) {
    this._el.querySelector("#stats-list").innerHTML = stats
      .map(({ card }) => `<option>${card}</option>`)
      .join("");
  }

  _renderStats({ card, stats }) {
    const interval = Math.round(24 * stats["Interval"]);
    const intervalString =
      (interval >= 24 ? `${Math.floor(interval / 24)}d ` : "") +
      `${interval % 24}h`;

    this._el.querySelector("#stats-card").innerHTML = card;
    this._el.querySelector("#reviewed-at").innerHTML = new Date(
      stats["LastReviewedAt"]
    ).toLocaleString();
    this._el.querySelector("#interval").innerHTML = intervalString;
    this._el.querySelector("#difficulty").innerHTML = stats["Difficulty"];
    this._graph.stats = stats["Historical"];
  }
}
