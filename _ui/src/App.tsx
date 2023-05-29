import './App.css';
import { useState } from 'react';

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

function Conversation({conversationid, messages, sendQuestion}: {conversationid: string, messages: message[], sendQuestion: (question: string) => void}) {
  const [newQuestion, setNewQuestion] = useState("");

  const messageList = messages.map((message, index)=>
    <li key={index}>
      <div>{message.query}</div>
      <div><pre>{message.data_csv}</pre></div>
      </li>)

  return <div className="Conversation">
    <p>{conversationid}</p>
    <div>
      <ul>
      {messageList}
      </ul>
    </div>
    <div>
      <input type="text" value={newQuestion} onChange={e => setNewQuestion(e.target.value)} />
      <button onClick={() => { 
        sendQuestion(newQuestion);
        setNewQuestion("");
      }}>Send</button>
    </div>
  </div>
}

export default App;
