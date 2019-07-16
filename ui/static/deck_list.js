export default class DeckList {
  async render() {
    const decks = await this._fetchDecks();
    document.getElementById("deckList").innerHTML = decks
      .sort((a, b) => a.name > b.name)
      .map(
        ({ name, cards_ready, next_review_at }) =>
          `<li>
<a href="#${name}" onclick="app.startSession('${name}',${cards_ready}); return false;">${name}</a>
<div>
  <code>${this._reviewStats(cards_ready, new Date(next_review_at))}</code>
</div>
<a class="stats-link" href="#stats-${name}" onclick="app.showStats('${name}'); return false;">
  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
    <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z"/>
    <path d="M0 0h24v24H0z" fill="none"/>
  </svg>
</a>
</li>`
      )
      .join("");
  }

  _fetchDecks() {
    return window.fetch("decks").then(res => {
      if (res.ok) return res.json();

      return res.text().then(text => alert(`Failed to fetch decks: ${text}`));
    });
  }

  _reviewStats(ready, next) {
    if (ready > 0) return ready;

    const diff = next - new Date();
    if (diff < 1000) return "available now";

    let d = next;
    d = [
      "0" + d.getDate(),
      "0" + (d.getMonth() + 1),
      "" + d.getFullYear(),
      "0" + d.getHours(),
      "0" + d.getMinutes()
    ].map(component => component.slice(-2));

    return d.slice(0, 3).join(".") + " " + d.slice(3).join(":");
  }
}
