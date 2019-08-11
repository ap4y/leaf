const template = `
<header>
  <h3>Deck: <span id="statsDeck"></span></h3>
  <div class="stats-select-row">
    <h4>Cards: </h4>
    <select id="statsList"></select>
  </div>
</header>

<main class="container">
  <h3>Current Stats for <span id="statsCard"></span></h3>

  <ul>
    <li>
      <strong>Last Reviewed At: </strong>
      <span id="reviewedAt"></span>
    </li>
    <li>
      <strong>Interval: </strong>
      <span id="interval"></span>
    </li>
    <li>
      <strong>Difficulty: </strong>
      <span id="difficulty"></span>
    </li>
    <li>
      <strong>Historical: </strong>
      <span id="historical"></span>
    </li>
  </ul>
</main>
`;

export default class StatsList {
  constructor() {
    this._el = document.createElement("div");
    this._el.innerHTML = template;
  }

  get element() {
    return this._el;
  }

  set deck(deck) {
    this._el.querySelector("#statsDeck").innerHTML = deck;
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

    this._el.querySelector("#statsList").onchange = ({ target }) => {
      const stat = stats.find(({ card }) => card === target.value);
      this._renderStats(stat);
    };
  }

  _populateSelect(stats) {
    this._el.querySelector("#statsList").innerHTML = stats
      .map(({ card }) => `<option>${card}</option>`)
      .join("");
  }

  _renderStats({ card, stats }) {
    const interval = Math.round(24 * stats["Interval"]);
    const intervalString =
      (interval >= 24 ? `${Math.floor(interval / 24)}d ` : "") +
      `${interval % 24}h`;

    this._el.querySelector("#statsCard").innerHTML = card;
    this._el.querySelector("#reviewedAt").innerHTML = new Date(
      stats["LastReviewedAt"]
    ).toLocaleString();
    this._el.querySelector("#interval").innerHTML = intervalString;
    this._el.querySelector("#difficulty").innerHTML = stats["Difficulty"];
    this._el.querySelector("#historical").innerHTML = (
      stats["Historical"] || []
    )
      .map(({ interval }) => interval)
      .join(", ");
  }
}
