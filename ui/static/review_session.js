export default class ReviewSession {
  constructor() {
    this.isAnswering = true;
    this.deck = null;
    this.session = null;
  }

  async render() {
    document.getElementById("inputForm").onsubmit = e => {
      e.preventDefault();
      if (this.isAnswering) {
        this._resolveAnswer();
      } else {
        this._nextQuestion();
      }
    };

    this.deck = window.history.state.deck;
    document.getElementById("deck").innerHTML = this.deck;
    this.session = await this._startSession(this.deck);
    this._updateState();
  }

  _updateState() {
    const { question, total, left, answerLen } = this.session;

    if (left === 0) {
      window.history.back();
      return;
    }

    document.getElementById("progress").innerHTML = `${total - left}/${total}`;
    document.getElementById("question").innerHTML = question;
    document.getElementById("input").style.width = `${answerLen}ch`;
    document.getElementById("answerState").innerHTML = "&nbsp";
    document.getElementById("correctAnswer").innerHTML = "&nbsp";
    document.getElementById("input").value = "";
    document.getElementById("input").focus();
  }

  async _nextQuestion() {
    this.session = await this._fetchNext();
    this.isAnswering = true;
    this._updateState();
  }

  async _resolveAnswer() {
    const answer = document.getElementById("input").value;
    const answerState = document.getElementById("answerState");
    const correctAnswer = document.getElementById("correctAnswer");

    this.isAnswering = false;
    const result = await this._submitAnswer(answer);
    if (result.is_correct) {
      answerState.innerHTML = "✓";
      answerState.style.color = "green";
      correctAnswer.innerHTML = "&nbsp";
    } else {
      answerState.innerHTML = "✕";
      answerState.style.color = "red";
      correctAnswer.innerHTML = result.correct;
    }
  }

  _startSession(deck) {
    return window
      .fetch(`start/${deck}`, {
        method: "POST"
      })
      .then(res => {
        if (res.ok) return res.json();

        return res
          .text()
          .then(text => alert(`Failed to start new session: ${text}`));
      });
  }

  _fetchNext() {
    return window.fetch("next").then(res => {
      if (res.ok) return res.json();

      return res
        .text()
        .then(text => alert(`Failed to fetch next question: ${text}`));
    });
  }

  _submitAnswer(answer) {
    return window
      .fetch("resolve", {
        method: "POST",
        body: JSON.stringify({ answer })
      })
      .then(res => {
        if (res.ok) return res.json();

        return res
          .text()
          .then(text => alert(`Failed to submit answer: ${text}`));
      });
  }
}