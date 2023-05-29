import SampleQuestions from './SampleQuestions';
import { useState } from 'react';

type message = {
    question: string,
    query: string,
    data_csv: string,
    error: string,
}

function Conversation({conversationid, messages, sendQuestion}: {conversationid: string, messages: message[], sendQuestion: (question: string) => void}) {
    const [newQuestion, setNewQuestion] = useState("");

    const messageList = messages.map((message, index)=>
        <div key={index}>
            <div className="uk-card uk-card-default uk-card-body uk-margin">
                {message.question}
            </div>
            <div className="uk-card uk-card-default uk-card-body uk-margin">
                <div className="query">{message.query}</div>
                <div className="data"><pre>{message.data_csv}</pre></div>
            </div>
        </div>)

    return <div className="Conversation">
        <p>{conversationid}</p>
        <div>
            <SampleQuestions conversationid={conversationid} />
        </div>
        <div>
            {messageList}
        </div>
        <hr/>
        <div>
            <input className="uk-input uk-width-4-5" type="text" placeholder="Ask a question about your data" value={newQuestion} onChange={e => setNewQuestion(e.target.value)} />
            <button className="uk-button uk-button-primary uk-width-1-5" onClick={() => { 
                sendQuestion(newQuestion);
                setNewQuestion("");
            }}>Send</button>
        </div>
    </div>
}

export default Conversation;
export type { message };