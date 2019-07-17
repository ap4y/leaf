export default class StatsList {
  async render() {
    this.deck = window.history.state.deck;
    document.getElementById("stats-deck").innerHTML = this.deck;

    const stats = await this._fetchStats();
    if (stats.length === 0) return;

    this._populateSelect(stats);
    this._renderStats(stats[0]);

    document.getElementById("statsList").onchange = ({ target }) => {
      const stat = stats.find(({ card }) => card === target.value);
      this._renderStats(stat);
    };
  }

  _populateSelect(stats) {
    document.getElementById("statsList").innerHTML = stats
      .map(({ card }) => `<option>${card}</option>`)
      .join("");
  }

  _renderStats({ card, stats }) {
    const interval = Math.round(24 * stats["Interval"]);
    const intervalString =
      (interval >= 24 ? `${Math.floor(interval / 24)}d ` : "") +
      `${interval % 24}h`;

    document.getElementById("stats-card").innerHTML = card;
    document.getElementById("reviewedAt").innerHTML = new Date(
      stats["LastReviewedAt"]
    ).toLocaleString();
    document.getElementById("interval").innerHTML = intervalString;
    document.getElementById("difficulty").innerHTML = stats["Difficulty"];
    document.getElementById("historical").innerHTML = (
      stats["Historical"] || []
    )
      .map(({ interval }) => interval)
      .join(", ");
  }

  async _fetchStats() {
    const res = await window.fetch(`stats/${this.deck}`);
    if (res.ok) return await res.json();

    alert(`Failed to fetch stats: ${await res.text()}`);
    return null;
  }
}
