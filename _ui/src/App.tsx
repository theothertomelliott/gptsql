import Conversation from './Conversation';
import type { message } from './Conversation';
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
    let addMessage = function(question: string, message: message) {
      setMessages(
        [
          ...messages,
          {
            question: question,
            query: message.query,
            data_csv: message.data_csv,
            error: message.error,
          },
        ]
      )
    }
    let sendQuestion = function(question: string) {
      sendNewQuestion(question, conversationID, addMessage, setError);
    }
    child = <Conversation conversationid={conversationID} sendQuestion={sendQuestion} messages={messages} />
  }

  return (
    <div className="App uk-container">
      <h1>
        GPTSQL
      </h1>
      {child}
    </div>
  );
}

function NewConversation({clicked}: {clicked: () => void}) {
  return <div className="uk-card uk-card-default uk-card-body">
      <button className="uk-button uk-button-default uk-width-1-1" onClick={clicked}>New Conversation</button>
    </div>;
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
  addMessage: (question: string, msg: message) => void,
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
      addMessage(question, result);
    },
    (error) => {
      setError(error);
    }
  )
}

export default App;
