import SampleQuestions from './SampleQuestions';
import { useState } from 'react';

type message = {
    question: string,
    query: string,
    data_csv: string,
    err: string,
}

function Conversation({conversationid, messages}: {conversationid: string, messages: message[]}) {
    const messageList = messages.map((message, index)=> {
        let content = <>
            <div className="uk-margin-bottom">
                <h3>SQL Query</h3>
                <div>{message.query}</div>
            </div>
            <div>
                <h3>Result</h3>
                <pre className="uk-height-medium uk-overflow-auto">{message.data_csv}</pre>
            </div>
        </>;
        console.log(message)
        if (message.err && message.err !== "") {
            content = <div>
            <h3>Error</h3>
            <div>{message.err}</div>
        </div>;
        }
        
        return <div key={index}>
            <div className="uk-card uk-card-default uk-card-body uk-margin uk-text-italic">
                {message.question}
            </div>
            <div className="uk-card uk-card-default uk-card-body uk-margin">
            {content}
            </div>
        </div>
    });

    return <div className="Conversation">
        <div>
            <SampleQuestions conversationid={conversationid} />
        </div>
        <div>
            {messageList}
        </div>
    </div>
}

function QuestionInput({sendQuestion}: {sendQuestion: (question: string) => void}) {
    const [newQuestion, setNewQuestion] = useState("");
    return <div className="uk-margin-bottom">
    <a id="question-input"></a>
    <input className="uk-input uk-width-4-5" type="text" placeholder="Ask a question about your data" value={newQuestion} onChange={e => setNewQuestion(e.target.value)} />
    <button className="uk-button uk-button-primary uk-width-1-5" onClick={() => { 
        sendQuestion(newQuestion);
        setNewQuestion("");
    }}>Send</button>
    </div>
}

export default Conversation;
export { QuestionInput };
export type { message };