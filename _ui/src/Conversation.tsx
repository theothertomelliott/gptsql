import SampleQuestions from './SampleQuestions';
import { useState } from 'react';

type message = {
    query: string,
    data_csv: string
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
        <SampleQuestions conversationid={conversationid} />
        </div>
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

export default Conversation;