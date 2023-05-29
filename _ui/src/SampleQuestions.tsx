
import { useState } from 'react';

function SampleQuestions({conversationid}: {conversationid: string}) {
    const [questions, setQuestions] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
  
    if (questions.length === 0 && !loading) {
      setLoading(true);
      fetch("/sample-questions",{
        method: "POST",
        body: JSON.stringify({
          conversation_id: conversationid
          })
        })
      .then(res => res.json())
      .then(
        (result) => {
        console.log(result);
        setQuestions(result.questions);
        setLoading(false);
      },
      (error) => {
        setError(error);
        setLoading(false);
      });
    }
  
    if (error !== "") {
      return <div>{error}</div>;
    }
  
    if (loading) {
      return <div>Loading sample questions...</div>
    }
  
    const questionList = questions.map((question, index)=>
      <li key={index}>{question}</li>);  
  
    return <div>
      <p>Welcome to GPTSQL! Your schema has been read and you may ask questions like the below:</p>
      <ul>
        {questionList}
      </ul>
    </div>;
  }

  export default SampleQuestions;