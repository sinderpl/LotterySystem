import React from 'react';
import ReactDOM from 'react-dom';

// const myfirstelement = <h1>Lottery System</h1>
// ReactDOM.render(myfirstelement, document.getElementById('root'));

// const formElement = (
//                     <div>
//                         <h1>Lottery System</h1>
//                         <form onSubmit = {checkTicket}>
//                             <label>Ticket Number:</label>
//                             <input type="text" id="lname" name="lname"></input>
//                             <input type="submit" value="Submit"></input>
//                         </form> 
//                     </div>)

// ReactDOM.render(formElement, document.getElementById('root'));



class TicketChecker extends React.Component {


    constructor(props) {
        super(props);
        let formElement = (
            <div>
                <h1>Lottery System</h1>
                <form onSubmit={this.checkTicket}>
                    <label>Ticket Number:</label>
                    <input type="text" id="lname" name="lname"></input>
                    <input type="submit" value="Submit"></input>
                </form>
            </div>)
        this.ticket = { ticketNumber: "", ticket: {} };

        ReactDOM.render(formElement, document.getElementById('root'));

        // This binding is necessary to make `this` work in the callback
        this.handleClick = this.handleClick.bind(this);
    }

  
    handleClick() {
    //   this.setState(state => ({
    //     isToggleOn: !state.isToggleOn
    //   }));
    }
  
    render() {
      return (
        <button onClick={this.handleClick}>
          {this.state.isToggleOn ? 'ON' : 'OFF'}
        </button>
      );
    }

    checkTicket(e) {
        e.preventDefault();
        console.log('The link was clicked.');
        console.log(e)
    }
  }
