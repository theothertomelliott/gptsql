import './App.css';
import Conversation from './Conversation';
import { useState } from 'react';

function App() {
  const [conversationID, setConversationID] = useState("");
  const [error, setError] = useState("");
  const [messages, setMessages] = useState<any[]>([]);
  let child: JSX.Element = <div></div>;
  if (conversationID === "") {
    let newConversationHandler = function() {
      createNewConversation(setConversationID,setError);
    }
    child = <NewConversation clicked={newConversationHandler} />
  } else if (error !== "") {
    console.log(error);
    child = <p>{JSON.stringify(error)}</p>
  } else {
    let addMessage = function(message: message) {
      setMessages(
        [
          ...messages,
          message,
        ]
      )
    }
    let sendQuestion = function(question: string) {
      sendNewQuestion(question, conversationID, addMessage, setError);
    }
    child = <Conversation conversationid={conversationID} sendQuestion={sendQuestion} messages={messages} />
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>
          GPTSQL
        </h1>
        {child}
      </header>
    </div>
  );
}

function NewConversation({clicked}: {clicked: () => void}) {
  return <button onClick={clicked}>New Conversation</button>
}


function createNewConversation(
  setConversationID : (value: string) => void,
  setError: (value: string) => void) {
  fetch("/new")
  .then(res => res.json())
  .then(
    (result) => {
      console.log(result);
      setConversationID(result.conversation_id);
    },
    (error) => {
      setError(error);
    }
  )
}

function sendNewQuestion(
  question: string, 
  conversationID: string,
  addMessage: (value: message) => void,
  setError: (value: string) => void
  ) {
  console.log("Sending question: " + question);
  fetch("/ask", {
    method: "POST",
    body: JSON.stringify({
      question: question,
      conversation_id: conversationID
      })
    })
  .then(res => res.json())
  .then(
    (result) => {
      console.log(result);
      addMessage(result);
    },
    (error) => {
      setError(error);
    }
  )
}

type message = {
  query: string,
  data_csv: string
}

export default App;
