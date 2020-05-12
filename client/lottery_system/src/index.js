import React from 'react';
import ReactDOM from 'react-dom';

class App extends React.Component {

  constructor(props) {
    super(props);
  }

  state = {
    ticketNumber: '',
    ticket: {}
  }

  handleChange = (e) => {
    this.setState({ [e.target.name]: e.target.value })
  }

  render() {
    return (
      <div>
        <h1>Lottery System</h1>
        <form onSubmit={this.checkTicket}>
          <label>Ticket Number:</label>
          <input type="text" name="ticketNumber" value={this.state.ticketNumber} onChange={this.handleChange}></input>
          <input type="submit" value="Submit"></input>
        </form>
      </div>)
  }

  checkTicket = (e) => {
    e.preventDefault();
    console.log(this.state.ticketNumber)

    const Http = new XMLHttpRequest();
    const url = 'https://localhost:8080/ticket/'//+this.state.ticketNumber;
    Http.open("GET", url);
    Http.send();
    
    Http.onreadystatechange = (e) => {
      console.log(e)
      console.log(Http.responseText)
    }
  }
}

ReactDOM.render(<App />, document.getElementById('root'))